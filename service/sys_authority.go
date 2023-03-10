package service

import (
	"go-admin/global"
	"go-admin/model"
)

// @description: 菜单与角色绑定
func SetMenuAuthority(auth *model.SysAuthority) error {
	var s model.SysAuthority
	global.GVA_DB.Preload("SysBaseMenus").First(&s, "authority_id = ?", auth.AuthorityId)
	err := global.GVA_DB.Model(&s).Association("SysBaseMenus").Replace(&auth.SysBaseMenus)
	return err
}
