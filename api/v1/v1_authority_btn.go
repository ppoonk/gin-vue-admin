package api

import (
	"go-admin/model"
	"go-admin/service"

	"github.com/gin-gonic/gin"
)

// GetAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   获取权限按钮
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SysAuthorityBtnReq                                      true  "菜单id, 角色id, 选中的按钮id"
// @Success   200   {object}  model.Response{data=model.SysAuthorityBtnRes,msg=string}  "返回列表成功"
// @Router    /authorityBtn/getAuthorityBtn [post]
func GetAuthorityBtn(c *gin.Context) {
	var req model.SysAuthorityBtnReq
	err := c.ShouldBindJSON(&req)
	res, err := service.GetAuthorityBtn(req)
	if err != nil {
		model.Result(7, "查询失败", nil, c)
		return
	}
	model.Result(0, "查询成功", res, c)
}

// SetAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   设置权限按钮
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SysAuthorityBtnReq     true  "菜单id, 角色id, 选中的按钮id"
// @Success   200   {object}  model.Response{msg=string}  "返回列表成功"
// @Router    /authorityBtn/setAuthorityBtn [post]
func SetAuthorityBtn(c *gin.Context) {
	var req model.SysAuthorityBtnReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = service.SetAuthorityBtn(req)
	if err != nil {
		model.Result(7, "分配失败", nil, c)
		return
	}
	model.Result(0, "分配成功", nil, c)
}

// CanRemoveAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   设置权限按钮
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  model.Response{msg=string}  "删除成功"
// @Router    /authorityBtn/canRemoveAuthorityBtn [post]
func CanRemoveAuthorityBtn(c *gin.Context) {
	id := c.Query("id")
	err := service.CanRemoveAuthorityBtn(id)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	model.Result(0, "删除成功", nil, c)

}
