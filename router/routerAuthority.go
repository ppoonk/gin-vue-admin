package router

import (
	"go-admin/api/v1"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

func InitAuthorityRouter(Router *gin.RouterGroup) {
	authorityRouter := Router.Group("authority").Use(middleware.OperationRecord())
	authorityRouterWithoutRecord := Router.Group("authority")

	{
		authorityRouter.POST("createAuthority", api.CreateAuthority)   // 创建角色
		authorityRouter.POST("deleteAuthority", api.DeleteAuthority)   // 删除角色
		authorityRouter.PUT("updateAuthority", api.UpdateAuthority)    // 更新角色
		authorityRouter.POST("copyAuthority", api.CopyAuthority)       // 拷贝角色
		authorityRouter.POST("setDataAuthority", api.SetDataAuthority) // 设置角色资源权限
	}
	{
		authorityRouterWithoutRecord.POST("getAuthorityList", api.GetAuthorityList) // 获取角色列表
	}
}
