package service

import (
	"coin-price-scan/enums"
	"coin-price-scan/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	HB_DOMAIN   = "https://otc-api-hk.eiijo.cn"
	BA_DOMAIN   = "https://p2p.binance.com"
)

type Strategy struct {
	PayMethod 		int
	Amount 			float64 `form:"amount" json:"amount" binding:"required"`
	TargetPrice 	float64 `form:"target_price" json:"target_price" binding:"required"`
	TargetProfit 	float64 `form:"target_profit" json:"target_profit" binding:"required"`
}

type HB struct {
	Code 		uint
	Message 	string
	TotalCount 	uint
	PageSize 	uint
	TotalPage 	uint
	CurrPage 	uint
	Data 		[]HBItem
}

type HBItem struct {
	Id					int
	UserName 			string
	IsOnline 			bool	//是否在线
	Price 				string  //单价
	MinTradeLimit 		string 	//限额
	MaxTradeLimit 		string 	//限额
	OrderCompleteRate 	string 	//订单完成率
	TradeMonthTimes 	uint    //月单数
}

type BA struct {
	Code 		string
	Message 	string
	Data 		[]BaItem
}

type BaItem struct {
	Adv 	   BAAdv
	Advertiser BAAdvertiser
}

type BAAdv struct {
	Price string
}

type BAAdvertiser struct {
	NickName string
}

//扫描火币
func ScanHb(payMethod,pageNum int,amount float64) ([]HBItem,error) {
	url := HB_DOMAIN+fmt.Sprintf("/v1/data/trade-market?coinId=2&currency=1&tradeType=sell&currPage=%v&payMethod=%v&acceptOrder=-1&country=&blockType=general&online=1&range=0&amount=%v",pageNum,payMethod,amount)
	var minPrice float64
	var scanResult []HBItem
	var err error

	httpClient := &util.HttpClient{}
	err = httpClient.Get(url,nil, func(resp *http.Response) {
		body,readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			err = readErr
			util.ErrDetail(enums.HB_REQUEST_ERR,fmt.Sprintf("火币数据获取异常,url=%v",url),readErr.Error())
			return
		}

		result := HB{}
		if err := json.Unmarshal(body,&result);err != nil {
			err = readErr
			util.ErrDetail(enums.HB_REQUEST_DATA_DECODE_ERR,fmt.Sprintf("火币解析数据异常,url=%v,result=%v",url,string(body)),err.Error())
			return
		}

		if result.Code != 200 {
			util.ErrDetail(enums.HB_REQUEST_DATA_DECODE_ERR,fmt.Sprintf("火币接口异常,url=%v,result=%v",url,string(body)),result.Message)
			return
		}

		for index,item := range result.Data{
			price,_ := strconv.ParseFloat(item.Price,64)
			if index == 0 {
				minPrice = price
			}
			scanResult = append(scanResult,item)
		}
	})

	if err != nil {
		util.Error(fmt.Sprintf("火币请求异常,err=%v,url=%v",err,url))
	}

	return scanResult,err
}

//扫描币安
func ScanBA() []BaItem {
	var err error
	var scanResult []BaItem

	data := make(map[string]interface{})
	data["asset"] 			= "USDT"
	data["fiat"] 			= "HKD"
	data["merchantCheck"] 	= false
	data["page"]			= 1
	data["payTypes"]		= []interface{}{}
	data["publisherType"]	= nil
	data["rows"]			= 10
	data["tradeType"]		= "SELL"
	data["transAmount"]		= "10000"

	byteData,encodeErr := json.Marshal(&data)
	if encodeErr != nil {
	}

	httpClient := &util.HttpClient{}
	err = httpClient.Post(BA_DOMAIN+"/bapi/c2c/v2/friendly/c2c/adv/search",string(byteData),nil, func(resp *http.Response) {
		body,readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			err = readErr
			util.ErrDetail(enums.HB_REQUEST_ERR,fmt.Sprintf("币安数据获取异常,err=%v"),readErr.Error())
			return
		}

		result := BA{}
		if err = json.Unmarshal(body,&result);err != nil {
			fmt.Println(err.Error())
			util.ErrDetail(enums.HB_REQUEST_DATA_DECODE_ERR,fmt.Sprintf("币安解析数据异常,result=%v",string(body)),err.Error())
			return
		}

		if result.Code != "000000" {
			util.ErrDetail(enums.HB_REQUEST_DATA_DECODE_ERR,fmt.Sprintf("币安接口异常,result=%v",string(body)),result.Message)
			return
		}

		for _,item := range result.Data{
			scanResult = append(scanResult,item)
		}
	})

	return scanResult
}

func Scan(strategy *Strategy,pushLog map[int]int)  {
	var scanErr error
	var HBResult []HBItem
	var BaResult []BaItem
	var CHRate float64
	baMaxPrice := 0.00
	baMinPrice := 0.00
	var wg sync.WaitGroup
	wg.Add(3)

	config,err := util.GetConfig()
	if err != nil {
		util.Error(fmt.Sprintf("err_code:%v，err_msg：%v",err.Code,err.Err.Error()))
		return
	}

	//获取火数据
	go func (){
		defer wg.Done()
		HBResult,scanErr = ScanHb(strategy.PayMethod,1,strategy.Amount)
		if scanErr != nil {

		}
	}()

	//获取币安数据
	go func() {
		defer wg.Done()
		BaResult = ScanBA();
		if len(BaResult) > 0 {
			baMaxPrice,_ = strconv.ParseFloat(BaResult[0].Adv.Price,64)
			baMinPrice,_ = strconv.ParseFloat(BaResult[len(BaResult)-1].Adv.Price,64)
		}
	}()

	//获取美元与人民币汇率 ,美元与港币汇率
	go func() {
		defer wg.Done()
		CHRate = GetRate("CNY","HKD")
	}()

	wg.Wait()

	BaCNYMaxPrice := baMaxPrice*CHRate
	BaCNYMinPrice := baMinPrice*CHRate

	for i := 0; i <= len(HBResult) - 1; i++  {
		hbFloatPrice,_ := strconv.ParseFloat(HBResult[i].Price,64)
		priceMaxDiff   := Decimal(BaCNYMaxPrice - hbFloatPrice)
		priceMinDiff   := Decimal(BaCNYMinPrice - hbFloatPrice)

		profitMax := Decimal(strategy.Amount/hbFloatPrice*priceMaxDiff)
		profitMin := Decimal(strategy.Amount/hbFloatPrice*priceMinDiff)

		_,exists := pushLog[HBResult[i].Id]
		if !exists && (profitMax >= strategy.TargetProfit || priceMaxDiff >= strategy.TargetPrice) {
			message := fmt.Sprintf(" 平台：%v\n 昵称：%v\n 价格：%v ￥\n 限额：%v ~ %v ￥\n 差价：%v ~ %v ￥\n 利润：%v ~ %v ￥\n 币安价格：%v ~ %v ￥\n","火币",HBResult[i].UserName,HBResult[i].Price,HBResult[i].MinTradeLimit,HBResult[i].MaxTradeLimit,priceMaxDiff,priceMinDiff,profitMax,profitMin,Decimal(BaCNYMaxPrice),Decimal(BaCNYMinPrice))
			SendDingDing(config["DING_ACCESS_TOKEN"], config["DING_SECRET"], DingDingReq{
				MsgType: "text",
				Text: DingDingText{
					Content: message,
				},
			})
			pushLog[HBResult[i].Id] = HBResult[i].Id
			return
		}
	}
}

func Decimal(v float64) float64 {
	v,_ = strconv.ParseFloat(fmt.Sprintf("%.2f",v),4)
	return v
}

//监听火币交易价格
func HBScanning(strategy *Strategy)  {
	hbPushLog := make(map[int]int)
	for {
		go Scan(strategy,hbPushLog)
		time.Sleep(time.Second*20)
	}
}
