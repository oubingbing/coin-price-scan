package main

import (
	"coin-price-scan/enums"
	"coin-price-scan/service"
	"coin-price-scan/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	strategy := service.Strategy{
		PayMethod:    1,
		Amount:       10000.00,
		TargetPrice:  0.15,
		TargetProfit: 200.00,
	}
	go service.HBScanning(&strategy)

	router := gin.Default()
	router.StaticFS("/statics", http.Dir("./statics"))
	router.LoadHTMLGlob("html/*")

	//首页
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"amount": strategy.Amount,
			"target_price": strategy.TargetPrice,
			"target_profit": strategy.TargetProfit,
		})
	})
	//更新策略
	router.POST("/update", func(ctx *gin.Context) {
		errInfo := &enums.ErrorInfo{}
		if errInfo.Err = ctx.ShouldBind(&strategy); errInfo.Err != nil {
			util.ResponseJson(ctx,enums.PARAM_ERR,errInfo.Err.Error(),nil)
			return
		}

		util.ResponseJson(ctx,enums.SUCCESS,"修改成功",nil)
		return
	})

	server := &http.Server{
		Addr:           ":8002",
		Handler:        router,
	}
	if 	err := server.ListenAndServe(); err != nil {
		util.Error(fmt.Sprintf("启动服务失败：%v\n",err))
	}
}
