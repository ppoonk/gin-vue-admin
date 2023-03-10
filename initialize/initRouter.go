package initialize

import (
	"go-admin/middleware"
	"go-admin/router"

	"github.com/gin-gonic/gin"
)

// 初始化总路由
func InitRouter() {
	Router := gin.Default()
	Router.Use(middleware.Cors()) //不开启跨域验证码出错

	//	Router.Static("/dist", "../web/dist")
	//Router.LoadHTMLGlob("templates/*")

	RouterV1Group := Router.Group("/")

	router.InitBaseRouter(RouterV1Group) //初始化base路由
	router.InitMenuRouter(RouterV1Group) //初始化menu路由
	router.InitUserRouter(RouterV1Group) //初始化user路由

	Router.Run(":8888")
}
