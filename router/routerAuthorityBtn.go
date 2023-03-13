package router

import (
	"go-admin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitAuthorityBtnRouterRouter(Router *gin.RouterGroup) {
	//authorityRouter := Router.Group("authorityBtn").Use(middleware.OperationRecord())
	authorityRouterWithoutRecord := Router.Group("authorityBtn")
	{
		authorityRouterWithoutRecord.POST("getAuthorityBtn", api.GetAuthorityBtn)
		authorityRouterWithoutRecord.POST("setAuthorityBtn", api.SetAuthorityBtn)
		authorityRouterWithoutRecord.POST("canRemoveAuthorityBtn", api.CanRemoveAuthorityBtn)
	}
}
