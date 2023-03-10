package router

import (
	"go-admin/api/v1"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user").Use(middleware.OperationRecord())
	userRouterWithoutRecord := Router.Group("user")

	{
		userRouter.POST("admin_register", api.Register)       // 管理员注册账号
		userRouter.POST("changePassword", api.ChangePassword) // 用户修改密码
		userRouter.POST("resetPassword", api.ResetPassword)   // 重制密码

		userRouter.POST("setUserAuthority", api.SetUserAuthority)     // 设置用户权限
		userRouter.POST("setUserAuthorities", api.SetUserAuthorities) // 设置用户权限组

		userRouter.DELETE("deleteUser", api.DeleteUser) // 删除用户
		userRouter.PUT("setUserInfo", api.SetUserInfo)  // 设置用户信息
		userRouter.PUT("setSelfInfo", api.SetSelfInfo)  // 设置自身信息

	}
	{
		userRouterWithoutRecord.POST("getUserList", api.GetUserList) // 分页获取用户列表
		userRouterWithoutRecord.GET("getUserInfo", api.GetUserInfo)  // 获取自身信息
	}
}
