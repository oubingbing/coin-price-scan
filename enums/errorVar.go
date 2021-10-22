package enums

import "errors"

//公共错误
var (
	DecodeErr						= errors.New("数据解析失败")
	SystemErr						= errors.New("系统繁忙，请稍后重试")
	NetErr							= errors.New("网络繁忙")
	JsonEncodeErr					= errors.New("json encode err")
	JsonDecodeErr					= errors.New("json decode err")
	TransFail						= errors.New("数据转换失败")
)
