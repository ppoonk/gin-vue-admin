package router

import (
	"go-admin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitJwtRouter(r *gin.RouterGroup) {
	jwtRouter := r.Group("jwt")
	{
		jwtRouter.POST("jsonInBlacklist", api.JsonInBlacklist) // jwt加入黑名单
	}
}
