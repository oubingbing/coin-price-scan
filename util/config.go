package util

import (
	"admin-api/enums"
	"fmt"
	"gopkg.in/ini.v1"
)

/**
 * 读取配置信息
 */
func GetConfig() (map[string]string,*enums.ErrorInfo) {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		Error(fmt.Sprintf("读取配置文件失败：%v",err.Error()))
		return nil,&enums.ErrorInfo{err,enums.GET_INI_ERR}
	}

	mp := make(map[string]string)
	mp["app_mode"] = cfg.Section("env").Key("app_mode").String()

	mp["DING_ACCESS_TOKEN"] = cfg.Section("dingding").Key("DING_ACCESS_TOKEN").String()
	mp["DING_SECRET"] 	  	= cfg.Section("dingding").Key("DING_SECRET").String()
	return mp,nil
}
