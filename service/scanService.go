package service

import (
	"admin-api/enums"
	"admin-api/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	HB_DOMAIN = "https://otc-api-hk.eiijo.cn"
)

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
	UserName 			string
	IsOnline 			bool	//是否在线
	Price 				string  //单价
	MinTradeLimit 		string 	//限额
	MaxTradeLimit 		string 	//限额
	OrderCompleteRate 	string 	//订单完成率
	TradeMonthTimes 	uint    //月单数
}

func ScanHb(payMethod,pageNum int,amount float64,targetPrice float64) (float64,[]HBItem,error) {
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
			if price <= targetPrice {
				scanResult = append(scanResult,item)
			}
		}
	})

	if err != nil {
		util.Error(fmt.Sprintf("火币请求异常,err=%v,url=%v",err,url))
	}

	return minPrice,scanResult,err
}

func Scan(payMethod,pageNum int,amount float64,targetPrice float64)  {
	ScanHb(payMethod,1,amount,targetPrice)

	config,err := util.GetConfig()
	if err != nil {
		util.Error(fmt.Sprintf("err_code:%v，err_msg：%v",err.Code,err.Err.Error()))
		return
	}

	SendDingDing(config["DING_ACCESS_TOKEN"], config["DING_SECRET"], DingDingReq{
		MsgType: "text",
		Text: DingDingText{
			Content: fmt.Sprintf(""),
		},
	})
}
