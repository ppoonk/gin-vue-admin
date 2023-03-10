package api

import (
	"go-admin/model"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

type SystemApiApi struct{}

// CreateApi
// @Tags      SysApi
// @Summary   创建基础api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SysApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /api/createApi [post]
func CreateApi(c *gin.Context) {
	var api model.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(api, utils.ApiVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = service.CreateApi(api)
	if err != nil {

		model.Result(7, "创建失败", nil, c)
		return
	}
	model.Result(0, "创建成功", nil, c)
}

// DeleteApi
// @Tags      SysApi
// @Summary   删除api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SysApi                  true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "删除api"
// @Router    /api/deleteApi [post]
func DeleteApi(c *gin.Context) {
	var api model.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(api.GVA_MODEL, utils.IdVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = service.DeleteApi(api)
	if err != nil {

		model.Result(7, "删除失败", nil, c)
		return
	}
	model.Result(0, "删除成功", nil, c)
}

// GetApiList
// @Tags      SysApi
// @Summary   分页获取API列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SearchApiParams                               true  "分页获取API列表"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/getApiList [post]
func GetApiList(c *gin.Context) {
	var pageInfo model.SearchApiParams
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	list, total, err := service.GetAPIInfoList(pageInfo.SysApi, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {

		model.Result(7, "获取失败", nil, c)
		return
	}
	model.Result(0, "获取成功", gin.H{
		"list":     list,
		"total":    total,
		"page":     pageInfo.Page,
		"pageSize": pageInfo.PageSize,
	}, c)
}

// GetApiById
// @Tags      SysApi
// @Summary   根据id获取api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.GetById                                   true  "根据id获取api"
// @Success   200   {object}  response.Response{data=systemRes.SysAPIResponse}  "根据id获取api,返回包括api详情"
// @Router    /api/getApiById [post]
func GetApiById(c *gin.Context) {
	var idInfo model.GetById
	err := c.ShouldBindJSON(&idInfo)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(idInfo, utils.IdVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	api, err := service.GetApiById(idInfo.ID)
	if err != nil {

		model.Result(7, "获取失败", nil, c)
		return
	}
	model.Result(0, "获取成功", gin.H{
		"api": api,
	}, c)
}

// UpdateApi
// @Tags      SysApi
// @Summary   修改基础api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SysApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "修改基础api"
// @Router    /api/updateApi [post]
func UpdateApi(c *gin.Context) {
	var api model.SysApi
	err := c.ShouldBindJSON(&api)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(api, utils.ApiVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = service.UpdateApi(api)
	if err != nil {

		model.Result(7, "修改失败", nil, c)
		return
	}
	model.Result(0, "修改成功", nil, c)
}

// GetAllApis
// @Tags      SysApi
// @Summary   获取所有的Api 不分页
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=systemRes.SysAPIListResponse,msg=string}  "获取所有的Api 不分页,返回包括api列表"
// @Router    /api/getAllApis [post]
func GetAllApis(c *gin.Context) {
	apis, err := service.GetAllApis()
	if err != nil {

		model.Result(7, "获取失败", nil, c)
		return
	}
	model.Result(0, "获取成功", gin.H{
		"apis": apis,
	}, c)
}

// DeleteApisByIds
// @Tags      SysApi
// @Summary   删除选中Api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.IdsReq                 true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "删除选中Api"
// @Router    /api/deleteApisByIds [delete]
func DeleteApisByIds(c *gin.Context) {
	var ids model.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = service.DeleteApisByIds(ids)
	if err != nil {
		model.Result(7, "删除失败", nil, c)
		return
	}
	model.Result(0, "删除成功", nil, c)
}
