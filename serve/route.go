package serve

import (
	"github.com/gin-gonic/gin"
	"log"
)

var (
	route *gin.Engine
)

func init() {
	route = gin.Default()
	route.GET("/", wrap(BaseInfo))            // 获取基础数据
	route.GET("/pull", wrap(PullFund))        // 重新拉取所有基金数据
	route.GET("/show_daily", wrap(FundDaily)) // 查看某个基金
	route.GET("/max_sharpe", wrap(MaxSharpe)) // 获取最大夏普值
}

func Run() {
	err := route.Run(":8000")
	if err != nil {
		log.Fatalln(err)
	}
}

func wrap(f func(ctx *BaseContext)) func(*gin.Context) {
	return func(context *gin.Context) {
		f(&BaseContext{Context: context})
	}
}
