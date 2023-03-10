package router

import (
	"go-admin/api/v1"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(r *gin.RouterGroup) {

	baseRouter := r.Group("base")
	{
		baseRouter.POST("login", api.Login)
		baseRouter.POST("captcha", api.Captcha)
		//ping
		baseRouter.GET("/ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})
	}
}
