package service

import (
	"coin-price-scan/enums"
	"coin-price-scan/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const RATE_DOMAIN = "https://api.it120.cc"

type Rate struct {
	Code int
	Data map[string]float64
	Msg string
}

func GetRate(from,to string) float64 {
	var err error
	url  := RATE_DOMAIN+fmt.Sprintf("/gooking/forex/rate?fromCode=%v&toCode=%v",from,to)
	rate := 0.0

	httpClient := &util.HttpClient{}
	err = httpClient.Get(url,nil, func(resp *http.Response) {
		body,readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			err = readErr
			util.ErrDetail(enums.HB_REQUEST_ERR,fmt.Sprintf("汇率数据获取异常,url=%v",url),readErr.Error())
			return
		}

		result := Rate{}
		if err := json.Unmarshal(body,&result);err != nil {
			err = readErr
			util.ErrDetail(enums.HB_REQUEST_DATA_DECODE_ERR,fmt.Sprintf("汇率解析数据异常,url=%v,result=%v",url,string(body)),err.Error())
			return
		}

		if result.Code != 0 {
			util.ErrDetail(enums.HB_REQUEST_DATA_DECODE_ERR,fmt.Sprintf("汇率接口异常,url=%v,result=%v",url,string(body)),result.Msg)
			return
		}

		rate = result.Data["rate"]
	})

	if err != nil {
		util.Error(fmt.Sprintf("汇率请求异常,err=%v,url=%v",err,url))
	}

	return rate
}
