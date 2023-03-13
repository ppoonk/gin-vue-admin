package service

import (
	"errors"
	"go-admin/global"
	"go-admin/model"

	"gorm.io/gorm"
)

func GetAuthorityBtn(req model.SysAuthorityBtnReq) (res model.SysAuthorityBtnRes, err error) {
	var authorityBtn []model.SysAuthorityBtn
	err = global.GVA_DB.Find(&authorityBtn, "authority_id = ? and sys_menu_id = ?", req.AuthorityId, req.MenuID).Error
	if err != nil {
		return
	}
	var selected []uint
	for _, v := range authorityBtn {
		selected = append(selected, v.SysBaseMenuBtnID)
	}
	res.Selected = selected
	return res, err
}

func SetAuthorityBtn(req model.SysAuthorityBtnReq) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var authorityBtn []model.SysAuthorityBtn
		err = tx.Delete(&[]model.SysAuthorityBtn{}, "authority_id = ? and sys_menu_id = ?", req.AuthorityId, req.MenuID).Error
		if err != nil {
			return err
		}
		for _, v := range req.Selected {
			authorityBtn = append(authorityBtn, model.SysAuthorityBtn{
				AuthorityId:      req.AuthorityId,
				SysMenuID:        req.MenuID,
				SysBaseMenuBtnID: v,
			})
		}
		if len(authorityBtn) > 0 {
			err = tx.Create(&authorityBtn).Error
		}
		if err != nil {
			return err
		}
		return err
	})
}

func CanRemoveAuthorityBtn(ID string) (err error) {
	fErr := global.GVA_DB.First(&model.SysAuthorityBtn{}, "sys_base_menu_btn_id = ?", ID).Error
	if errors.Is(fErr, gorm.ErrRecordNotFound) {
		return nil
	}
	return errors.New("此按钮正在被使用无法删除")
}
