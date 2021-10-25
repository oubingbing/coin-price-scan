package service

import (
	"coin-price-scan/enums"
	"coin-price-scan/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	GATE_DOMAIN = "https://www.gate.io"
)

type Gate struct {
	Result bool
	Push_Order []GateItem
}

type GateItem struct {
	Oid string
	Type string
	Rate string
	Amount string
	Limit_Total string
	Username string
}

//通用数据
type CommonItem struct {
	Exchange string
	Id int
	Username string
	Price float64
	LimitAmount string
}

//扫描芝麻开门
func ScanGate(amount string) []CommonItem {
	var err error
	var scanResult []CommonItem

	data := url.Values{}
	data.Set("type","push_order_list")
	data.Set("symbol","USDT_CNY")
	data.Set("big_trade","0")
	data.Set("amount",amount)
	data.Set("pay_type","")
	data.Set("is_blue","1")

	httpClient := &util.HttpClient{}
	err = httpClient.PostForm(GATE_DOMAIN+"/json_svr/query_push/?u=21&c=416424",data,nil, func(resp *http.Response) {
		body,readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			err = readErr
			util.ErrDetail(enums.HB_REQUEST_ERR,fmt.Sprintf("gate数据获取异常,err=%v"),readErr.Error())
			return
		}

		result := Gate{}
		if err = json.Unmarshal(body,&result);err != nil {
			util.ErrDetail(enums.HB_REQUEST_DATA_DECODE_ERR,fmt.Sprintf("gate解析数据异常,result=%v",string(body)),err.Error())
			return
		}

		if !result.Result {
			util.ErrDetail(enums.HB_REQUEST_DATA_DECODE_ERR,fmt.Sprintf("gate接口异常,result=%v",string(body)),nil)
			return
		}

		for _,item := range result.Push_Order{
			if item.Type == "sell" {
				price,_ := strconv.ParseFloat(item.Rate,64)
				id,_ := strconv.Atoi(item.Oid)
				commonItem := CommonItem{"芝麻开门",id,item.Username,price,item.Limit_Total}
				scanResult = append(scanResult,commonItem)
			}
		}
	})

	return scanResult
}