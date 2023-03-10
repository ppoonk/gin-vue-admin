package service

import (
	"errors"
	"go-admin/global"
	"go-admin/model"
	"strconv"

	"gorm.io/gorm"
)

// UserAuthorityDefaultRouter 用户角色默认路由检查
func UserAuthorityDefaultRouter(user *model.SysUser) {
	var menuIds []string
	err := global.GVA_DB.Model(&model.SysAuthorityMenu{}).Where("sys_authority_authority_id = ?", user.AuthorityId).Pluck("sys_base_menu_id", &menuIds).Error
	if err != nil {
		return
	}
	var am model.SysBaseMenu
	err = global.GVA_DB.First(&am, "name = ? and id in (?)", user.Authority.DefaultRouter, menuIds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user.Authority.DefaultRouter = "404"
	}
}

// 获取子菜单
func getChildrenList(menu *model.SysMenu, treeMap map[string][]model.SysMenu) (err error) {
	menu.Children = treeMap[menu.MenuId]
	for i := 0; i < len(menu.Children); i++ {
		err = getChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

// 获取动态菜单树
func GetMenuTree(authorityId uint) (menus []model.SysMenu, err error) {
	menuTree, err := getMenuTreeMap(authorityId)
	menus = menuTree["0"]
	for i := 0; i < len(menus); i++ {
		err = getChildrenList(&menus[i], menuTree)
	}
	return menus, err
}

func getMenuTreeMap(authorityId uint) (treeMap map[string][]model.SysMenu, err error) {
	var allMenus []model.SysMenu
	var baseMenu []model.SysBaseMenu
	var btns []model.SysAuthorityBtn
	treeMap = make(map[string][]model.SysMenu)
	var SysAuthorityMenus []model.SysAuthorityMenu

	err = global.GVA_DB.Debug().Where("sys_authority_authority_id = ?", authorityId).Find(&SysAuthorityMenus).Error
	if err != nil {
		return
	}
	//以上，通过authorityId从数据库获取该用户的menu列表
	var MenuIds []string
	//menu数组
	for i := range SysAuthorityMenus {
		MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
	}
	//获取menu
	//此时该用户的MenuIds： [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29]
	//从完整menu中获取该用户的menu
	err = global.GVA_DB.Debug().Where("id in (?)", MenuIds).Order("sort").Preload("Parameters").Find(&baseMenu).Error
	//SELECT * FROM `sys_base_menu_parameters` WHERE `sys_base_menu_parameters`.`sys_base_menu_id` IN (26,23,22,21,27,4,9,20,17,1,5,28,18,29,19,6,3,11,7,16,13,8,14,10,25,15,12,24,2) AND `sys_base_menu_parameters`.`deleted_at` IS NULL
	//SELECT * FROM `sys_base_menus` WHERE id in ('1','2','3','4','5','6','7','8','9','10','11','12','13','14','15','16','17','18','19','20','21','22','23','24','25','26','27','28','29') AND `sys_base_menus`.`deleted_at` IS NULL ORDER BY sort
	if err != nil {
		return
	}
	//所有的menuid拼接
	for i := range baseMenu {
		allMenus = append(allMenus, model.SysMenu{
			SysBaseMenu: baseMenu[i],
			AuthorityId: authorityId,
			MenuId:      strconv.Itoa(int(baseMenu[i].ID)),
			Parameters:  baseMenu[i].Parameters,
		})
	}

	err = global.GVA_DB.Debug().Where("authority_id = ?", authorityId).Preload("SysBaseMenuBtn").Find(&btns).Error
	//SELECT * FROM `sys_authority_btns` WHERE authority_id = 888
	if err != nil {
		return
	}
	var btnMap = make(map[uint]map[string]uint)
	for _, v := range btns {
		if btnMap[v.SysMenuID] == nil {
			btnMap[v.SysMenuID] = make(map[string]uint)
		}
		btnMap[v.SysMenuID][v.SysBaseMenuBtn.Name] = authorityId
	}
	for _, v := range allMenus {
		v.Btns = btnMap[v.ID]
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}

	return treeMap, err
}

// 分页获取基础menu列表
func GetInfoList() (list interface{}, total int64, err error) {
	var menuList []model.SysBaseMenu
	treeMap, err := getBaseMenuTreeMap()
	menuList = treeMap["0"]
	for i := 0; i < len(menuList); i++ {
		err = getBaseChildrenList(&menuList[i], treeMap)
	}
	return menuList, total, err
}

// 获取路由总树map
func getBaseMenuTreeMap() (treeMap map[string][]model.SysBaseMenu, err error) {
	var allMenus []model.SysBaseMenu
	treeMap = make(map[string][]model.SysBaseMenu)
	err = global.GVA_DB.Order("sort").Preload("MenuBtn").Preload("Parameters").Find(&allMenus).Error
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return treeMap, err
}

// 获取菜单的子菜单
func getBaseChildrenList(menu *model.SysBaseMenu, treeMap map[string][]model.SysBaseMenu) (err error) {
	menu.Children = treeMap[strconv.Itoa(int(menu.ID))]
	for i := 0; i < len(menu.Children); i++ {
		err = getBaseChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

// 获取基础路由树
func GetBaseMenuTree() (menus []model.SysBaseMenu, err error) {
	treeMap, err := getBaseMenuTreeMap()
	menus = treeMap["0"]
	for i := 0; i < len(menus); i++ {
		err = getBaseChildrenList(&menus[i], treeMap)
	}
	return menus, err
}

// 查看当前角色树
func GetMenuAuthority(info *model.GetAuthorityId) (menus []model.SysMenu, err error) {
	var baseMenu []model.SysBaseMenu
	var SysAuthorityMenus []model.SysAuthorityMenu
	err = global.GVA_DB.Where("sys_authority_authority_id = ?", info.AuthorityId).Find(&SysAuthorityMenus).Error
	if err != nil {
		return
	}

	var MenuIds []string

	for i := range SysAuthorityMenus {
		MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
	}

	err = global.GVA_DB.Where("id in (?) ", MenuIds).Order("sort").Find(&baseMenu).Error

	for i := range baseMenu {
		menus = append(menus, model.SysMenu{
			SysBaseMenu: baseMenu[i],
			AuthorityId: info.AuthorityId,
			MenuId:      strconv.Itoa(int(baseMenu[i].ID)),
			Parameters:  baseMenu[i].Parameters,
		})
	}
	// sql := "SELECT authority_menu.keep_alive,authority_menu.default_menu,authority_menu.created_at,authority_menu.updated_at,authority_menu.deleted_at,authority_menu.menu_level,authority_menu.parent_id,authority_menu.path,authority_menu.`name`,authority_menu.hidden,authority_menu.component,authority_menu.title,authority_menu.icon,authority_menu.sort,authority_menu.menu_id,authority_menu.authority_id FROM authority_menu WHERE authority_menu.authority_id = ? ORDER BY authority_menu.sort ASC"
	// err = global.GVA_DB.Raw(sql, authorityId).Scan(&menus).Error
	return menus, err
}

// 根据id获取菜单
func GetBaseMenuById(id int) (menu model.SysBaseMenu, err error) {
	err = global.GVA_DB.Preload("MenuBtn").Preload("Parameters").Where("id = ?", id).First(&menu).Error
	return
}

// 添加基础路由
func AddBaseMenu(menu model.SysBaseMenu) error {
	if !errors.Is(global.GVA_DB.Where("name = ?", menu.Name).First(&model.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	return global.GVA_DB.Create(&menu).Error
}

// 为角色增加menu树,增加menu和角色关联关系
func AddMenuAuthority(menus []model.SysBaseMenu, authorityId uint) (err error) {
	var auth model.SysAuthority
	auth.AuthorityId = authorityId
	auth.SysBaseMenus = menus
	err = SetMenuAuthority(&auth)
	return err
}

// 删除基础路由
func DeleteBaseMenu(id int) (err error) {
	err = global.GVA_DB.Preload("MenuBtn").Preload("Parameters").Where("parent_id = ?", id).First(&model.SysBaseMenu{}).Error
	if err != nil {
		var menu model.SysBaseMenu
		db := global.GVA_DB.Preload("SysAuthoritys").Where("id = ?", id).First(&menu).Delete(&menu)
		err = global.GVA_DB.Delete(&model.SysBaseMenuParameter{}, "sys_base_menu_id = ?", id).Error
		err = global.GVA_DB.Delete(&model.SysBaseMenuBtn{}, "sys_base_menu_id = ?", id).Error
		err = global.GVA_DB.Delete(&model.SysAuthorityBtn{}, "sys_menu_id = ?", id).Error
		if err != nil {
			return err
		}
		if len(menu.SysAuthoritys) > 0 {
			err = global.GVA_DB.Model(&menu).Association("SysAuthoritys").Delete(&menu.SysAuthoritys)
		} else {
			err = db.Error
			if err != nil {
				return
			}
		}
	} else {
		return errors.New("此菜单存在子菜单不可删除")
	}
	return err
}

//更新路由

func UpdateBaseMenu(menu model.SysBaseMenu) (err error) {
	var oldMenu model.SysBaseMenu
	upDateMap := make(map[string]interface{})
	upDateMap["keep_alive"] = menu.KeepAlive
	upDateMap["close_tab"] = menu.CloseTab
	upDateMap["default_menu"] = menu.DefaultMenu
	upDateMap["parent_id"] = menu.ParentId
	upDateMap["path"] = menu.Path
	upDateMap["name"] = menu.Name
	upDateMap["hidden"] = menu.Hidden
	upDateMap["component"] = menu.Component
	upDateMap["title"] = menu.Title
	upDateMap["active_name"] = menu.ActiveName
	upDateMap["icon"] = menu.Icon
	upDateMap["sort"] = menu.Sort

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", menu.ID).Find(&oldMenu)
		if oldMenu.Name != menu.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", menu.ID, menu.Name).First(&model.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
				//global.GVA_LOG.Debug("存在相同name修改失败")
				return errors.New("存在相同name修改失败")
			}
		}
		txErr := tx.Unscoped().Delete(&model.SysBaseMenuParameter{}, "sys_base_menu_id = ?", menu.ID).Error
		if txErr != nil {
			//global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		txErr = tx.Unscoped().Delete(&model.SysBaseMenuBtn{}, "sys_base_menu_id = ?", menu.ID).Error
		if txErr != nil {
			//global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		if len(menu.Parameters) > 0 {
			for k := range menu.Parameters {
				menu.Parameters[k].SysBaseMenuID = menu.ID
			}
			txErr = tx.Create(&menu.Parameters).Error
			if txErr != nil {
				//global.GVA_LOG.Debug(txErr.Error())
				return txErr
			}
		}

		if len(menu.MenuBtn) > 0 {
			for k := range menu.MenuBtn {
				menu.MenuBtn[k].SysBaseMenuID = menu.ID
			}
			txErr = tx.Create(&menu.MenuBtn).Error
			if txErr != nil {
				//global.GVA_LOG.Debug(txErr.Error())
				return txErr
			}
		}

		txErr = db.Updates(upDateMap).Error
		if txErr != nil {
			//global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}
