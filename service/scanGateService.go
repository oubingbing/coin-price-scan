package service

import (
	"coin-price-scan/enums"
	"coin-price-scan/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	GATE_DOMAIN = "https://www.gate.io"
)

type Gate struct {
	Result bool
	PushOrder []GateItem
}

type GateItem struct {
	Rate string `json:"rate"`
	LimitTotal string `json:"limit_total"`
	Username string `json:"username"`
}

//扫描币安
func ScanGate(amount string) []GateItem {
	var err error
	var scanResult []GateItem

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

		fmt.Println(string(body))
		result := Gate{}
		fmt.Println(json.Unmarshal(body,&result))

		if err = json.Unmarshal(body,&result);err != nil {
			fmt.Println(err.Error())
			util.ErrDetail(enums.HB_REQUEST_DATA_DECODE_ERR,fmt.Sprintf("gate解析数据异常,result=%v",string(body)),err.Error())
			return
		}

		/*if !result.Result {
			util.ErrDetail(enums.HB_REQUEST_DATA_DECODE_ERR,fmt.Sprintf("gate接口异常,result=%v",string(body)),nil)
			return
		}*/

		fmt.Println(result)

		for _,item := range result.PushOrder{
			fmt.Println(item)
			//scanResult = append(scanResult,item)
		}
	})

	return scanResult
}