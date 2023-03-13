package api

import (
	"go-admin/model"
	"go-admin/service"

	"github.com/gin-gonic/gin"
)

// jwt加入黑名单
func JsonInBlacklist(c *gin.Context) {
	token := c.Request.Header.Get("x-token")
	jwt := model.JwtBlacklist{Jwt: token}
	err := service.JsonInBlacklist(jwt)
	if err != nil {
		model.Result(7, "jwt作废失败", nil, c)
		return
	}
	model.Result(0, "jwt作废成功", nil, c)
}
