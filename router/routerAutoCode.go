package router

import (
	"go-admin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitAutoCodeRouter(Router *gin.RouterGroup) {
	autoCodeRouter := Router.Group("autoCode")

	{
		autoCodeRouter.GET("getDB", api.GetDB)                  // 获取数据库
		autoCodeRouter.GET("getTables", api.GetTables)          // 获取对应数据库的表
		autoCodeRouter.GET("getColumn", api.GetColumn)          // 获取指定表所有字段信息
		autoCodeRouter.POST("preview", api.PreviewTemp)         // 获取自动创建代码预览
		autoCodeRouter.POST("createTemp", api.CreateTemp)       // 创建自动化代码
		autoCodeRouter.POST("createPackage", api.CreatePackage) // 创建package包
		autoCodeRouter.POST("getPackage", api.GetPackage)       // 获取package包
		autoCodeRouter.POST("delPackage", api.DelPackage)       // 删除package包
		autoCodeRouter.POST("createPlug", api.AutoPlug)         // 自动插件包模板
		autoCodeRouter.POST("installPlugin", api.InstallPlugin) // 自动安装插件
	}
}
