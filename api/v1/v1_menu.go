package api

import (
	"go-admin/model"
	"go-admin/service"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
)

// 获取用户动态路由
func GetMenu(c *gin.Context) {
	menus, err := service.GetMenuTree(utils.GetUserAuthorityId(c))
	if err != nil {
		model.Result(7, "获取失败", nil, c)
	}
	if menus == nil {
		menus = []model.SysMenu{}
	}
	model.Result(0, "获取成功", gin.H{
		"menus": menus,
	}, c)
}

// 分页获取基础menu列表
func GetMenuList(c *gin.Context) {
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
	menuList, total, err := service.GetInfoList()
	if err != nil {
		model.Result(7, "获取失败", nil, c)
		return
	}

	model.Result(0, "获取成功", gin.H{
		"list":     menuList,
		"total":    total,
		"page":     pageInfo.Page,
		"pageSize": pageInfo.PageSize,
	}, c)
}

// 获取用户动态路由
func GetBaseMenuTree(c *gin.Context) {
	menus, err := service.GetBaseMenuTree()
	if err != nil {
		model.Result(7, "获取失败", nil, c)
		return
	}
	model.Result(0, "获取成功", gin.H{
		"menus": menus,
	}, c)
}

// 获取指定角色menu
func GetMenuAuthority(c *gin.Context) {
	var param model.GetAuthorityId
	err := c.ShouldBindJSON(&param)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(param, utils.AuthorityIdVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	menus, err := service.GetMenuAuthority(&param)
	if err != nil {
		model.Result(7, "获取失败", nil, c)
		//response.FailWithDetailed(modelRes.SysMenusResponse{Menus: menus}, "获取失败", c)
		return
	}
	model.Result(0, "获取成功", gin.H{
		"menus": menus,
	}, c)
}

// 根据id获取菜单
func GetBaseMenuById(c *gin.Context) {
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
	menu, err := service.GetBaseMenuById(idInfo.ID)
	if err != nil {

		model.Result(7, "获取失败", nil, c)
		return
	}
	model.Result(0, "获取成功", gin.H{
		"menu": menu,
	}, c)
}

// 新增菜单
func AddBaseMenu(c *gin.Context) {
	var menu model.SysBaseMenu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(menu, utils.MenuVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(menu.Meta, utils.MenuMetaVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = service.AddBaseMenu(menu)
	if err != nil {
		model.Result(7, "添加失败", nil, c)
		return
	}
	model.Result(0, "添加成功", nil, c)

}

// 增加menu和角色关联关系
func AddMenuAuthority(c *gin.Context) {
	var authorityMenu model.AddMenuAuthorityInfo
	err := c.ShouldBindJSON(&authorityMenu)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	if err := utils.Verify(authorityMenu, utils.AuthorityIdVerify); err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	if err := service.AddMenuAuthority(authorityMenu.Menus, authorityMenu.AuthorityId); err != nil {
		model.Result(7, "添加失败!", nil, c)
	} else {
		model.Result(0, "添加成功", nil, c)
	}
}

// 删除菜单
func DeleteBaseMenu(c *gin.Context) {
	var menu model.GetById
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(menu, utils.IdVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = service.DeleteBaseMenu(menu.ID)
	if err != nil {
		model.Result(7, "删除失败", nil, c)
		return
	}
	model.Result(0, "删除成功", nil, c)
}

// 更新菜单

func UpdateBaseMenu(c *gin.Context) {
	var menu model.SysBaseMenu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(menu, utils.MenuVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = utils.Verify(menu.Meta, utils.MenuMetaVerify)
	if err != nil {
		model.Result(7, err.Error(), nil, c)
		return
	}
	err = service.UpdateBaseMenu(menu)
	if err != nil {
		model.Result(7, "更新失败", nil, c)
		return
	}
	model.Result(0, "更新成功", nil, c)
}
