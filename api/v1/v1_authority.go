package api

import (
	"go-admin/model"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// CreateAuthority
// @Tags      Authority
// @Summary   创建角色
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SysAuthority                                                true  "权限id, 权限名, 父角色id"
// @Success   200   {object}  response.Response{data=systemRes.SysAuthorityResponse,msg=string}  "创建角色,返回包括系统角色详情"
// @Router    /authority/createAuthority [post]
func CreateAuthority(c *gin.Context) {
	var authority model.SysAuthority
	err := c.ShouldBindJSON(&authority)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}

	err = utils.Verify(authority, utils.AuthorityVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	if authBack, err := service.CreateAuthority(authority); err != nil {
		model.Result(7, "创建失败"+err.Error(), nil, c)
	} else {
		_ = service.AddMenuAuthority(model.DefaultMenu(), authority.AuthorityId)
		_ = service.UpdateCasbin(authority.AuthorityId, model.DefaultCasbin())
		model.Result(0, "创建成功", gin.H{
			"authority": authBack,
		}, c)
	}
}

// CopyAuthority
// @Tags      Authority
// @Summary   拷贝角色
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      response.SysAuthorityCopyResponse                                  true  "旧角色id, 新权限id, 新权限名, 新父角色id"
// @Success   200   {object}  response.Response{data=systemRes.SysAuthorityResponse,msg=string}  "拷贝角色,返回包括系统角色详情"
// @Router    /authority/copyAuthority [post]
func CopyAuthority(c *gin.Context) {
	var copyInfo model.SysAuthorityCopyResponse
	err := c.ShouldBindJSON(&copyInfo)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(copyInfo, utils.OldAuthorityVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(copyInfo.Authority, utils.AuthorityVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	authBack, err := service.CopyAuthority(copyInfo)
	if err != nil {
		model.Result(7, "拷贝失败"+err.Error(), nil, c)
		return
	}
	model.Result(0, "拷贝成功", gin.H{
		"authority": authBack,
	}, c)

}

// DeleteAuthority
// @Tags      Authority
// @Summary   删除角色
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SysAuthority            true  "删除角色"
// @Success   200   {object}  response.Response{msg=string}  "删除角色"
// @Router    /authority/deleteAuthority [post]
func DeleteAuthority(c *gin.Context) {
	var authority model.SysAuthority
	err := c.ShouldBindJSON(&authority)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(authority, utils.AuthorityIdVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = service.DeleteAuthority(&authority)
	if err != nil { // 删除角色之前需要判断是否有用户正在使用此角色
		model.Result(7, "删除失败"+err.Error(), nil, c)
		return
	}
	model.Result(0, "删除成功", nil, c)
}

// UpdateAuthority
// @Tags      Authority
// @Summary   更新角色信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SysAuthority                                                true  "权限id, 权限名, 父角色id"
// @Success   200   {object}  response.Response{data=systemRes.SysAuthorityResponse,msg=string}  "更新角色信息,返回包括系统角色详情"
// @Router    /authority/updateAuthority [post]
func UpdateAuthority(c *gin.Context) {
	var auth model.SysAuthority
	err := c.ShouldBindJSON(&auth)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(auth, utils.AuthorityVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	authority, err := service.UpdateAuthority(auth)
	if err != nil {
		model.Result(7, "更新失败"+err.Error(), nil, c)
		return
	}
	model.Result(0, "更新成功", gin.H{
		"authority": authority,
	}, c)
}

// GetAuthorityList
// @Tags      Authority
// @Summary   分页获取角色列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取角色列表,返回包括列表,总数,页码,每页数量"
// @Router    /authority/getAuthorityList [post]
func GetAuthorityList(c *gin.Context) {
	var pageInfo model.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	list, total, err := service.GetAuthorityInfoList(pageInfo)
	if err != nil {
		model.Result(7, "获取失败"+err.Error(), nil, c)
		return
	}
	model.Result(0, "更新成功", gin.H{
		"list":     list,
		"total":    total,
		"page":     pageInfo.Page,
		"pageSize": pageInfo.PageSize,
	}, c)
}

// SetDataAuthority
// @Tags      Authority
// @Summary   设置角色资源权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      model.SysAuthority            true  "设置角色资源权限"
// @Success   200   {object}  response.Response{msg=string}  "设置角色资源权限"
// @Router    /authority/setDataAuthority [post]
func SetDataAuthority(c *gin.Context) {
	var auth model.SysAuthority
	err := c.ShouldBindJSON(&auth)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(auth, utils.AuthorityIdVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = service.SetDataAuthority(auth)
	if err != nil {
		model.Result(7, "设置失败"+err.Error(), nil, c)
		return
	}
	model.Result(0, "设置成功", nil, c)
}
