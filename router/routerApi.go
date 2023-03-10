package router

import (
	"go-admin/api/v1"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

func InitApiRouter(Router *gin.RouterGroup) {
	apiRouter := Router.Group("api").Use(middleware.OperationRecord())
	apiRouterWithoutRecord := Router.Group("api")
	{
		apiRouter.POST("createApi", api.CreateApi)               // 创建Api
		apiRouter.POST("deleteApi", api.DeleteApi)               // 删除Api
		apiRouter.POST("getApiById", api.GetApiById)             // 获取单条Api消息
		apiRouter.POST("updateApi", api.UpdateApi)               // 更新api
		apiRouter.DELETE("deleteApisByIds", api.DeleteApisByIds) // 删除选中api
	}
	{
		apiRouterWithoutRecord.POST("getAllApis", api.GetAllApis) // 获取所有api
		apiRouterWithoutRecord.POST("getApiList", api.GetApiList) // 获取Api列表
	}
}
