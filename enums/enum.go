package enums

//时间格式化
const TIME_FORMAT = "2006-01-02 15:04:05"

//业务不相关错误
const (
	SUCCESS					= 0
	FAIL 					= 1
	DB_CONNECT_ERR 			= 2
	PARAM_ERR				= 3
	SYSTEM_ERR				= 4	//系统异常
	GET_INI_ERR				= 5 //读取配置文件错误
	JSON_ENCODE_ERR			= 6
	JSON_DECODE_ERR			= 7
	PARAMS_ERROR 			= 8	//参数错误
)

//扫描火币 2500 ~ 2999
const (
	HB_REQUEST_ERR = 2500 //获取数据异常
	HB_REQUEST_DATA_DECODE_ERR = 2501 //数据解析失败
)
