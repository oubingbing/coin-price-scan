package main

import (
	"admin-api/service"
	"admin-api/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	payMethod := 1	//银行卡
	amount := 10000.0
	targetPrice := 6.14



	router := gin.Default()
	//首页
	router.GET("/", func(c *gin.Context) {

	})
	//更新策略
	router.POST("/update", func(c *gin.Context) {

	})

	server := &http.Server{
		Addr:           ":8002",
		Handler:        router,
	}

	if 	err := server.ListenAndServe(); err != nil {
		util.Error(fmt.Sprintf("启动服务失败：%v\n",err))
	}

}
