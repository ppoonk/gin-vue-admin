package router

import (
	"go-admin/api/v1"
	"go-admin/middleware"

	"github.com/gin-gonic/gin"
)

func InitMenuRouter(Router *gin.RouterGroup) {
	menuRouter := Router.Group("menu").Use(middleware.OperationRecord())
	menuRouterWithoutRecord := Router.Group("menu")

	{
		menuRouter.POST("addBaseMenu", api.AddBaseMenu)           // 新增菜单
		menuRouter.POST("addMenuAuthority", api.AddMenuAuthority) //	增加menu和角色关联关系
		menuRouter.POST("deleteBaseMenu", api.DeleteBaseMenu)     // 删除菜单
		menuRouter.POST("updateBaseMenu", api.UpdateBaseMenu)     // 更新菜单
	}
	{
		menuRouterWithoutRecord.POST("getMenu", api.GetMenu)                   // 获取菜单树
		menuRouterWithoutRecord.POST("getMenuList", api.GetMenuList)           // 分页获取基础menu列表
		menuRouterWithoutRecord.POST("getBaseMenuTree", api.GetBaseMenuTree)   // 获取用户动态路由
		menuRouterWithoutRecord.POST("getMenuAuthority", api.GetMenuAuthority) // 获取指定角色menu
		menuRouterWithoutRecord.POST("getBaseMenuById", api.GetBaseMenuById)   // 根据id获取菜单
	}

}
