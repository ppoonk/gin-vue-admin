package model

import (
	"strconv"
	"strings"
	"time"

	"go-admin/global"
)

// SysAutoCodeHistory 自动迁移代码记录,用于回滚,重放使用
type SysAutoCodeHistory struct {
	global.GVA_MODEL
	Package       string `json:"package"`
	BusinessDB    string `json:"businessDB"`
	TableName     string `json:"tableName"`
	RequestMeta   string `gorm:"type:text" json:"requestMeta,omitempty"`   // 前端传入的结构化信息
	AutoCodePath  string `gorm:"type:text" json:"autoCodePath,omitempty"`  // 其他meta信息 path;path
	InjectionMeta string `gorm:"type:text" json:"injectionMeta,omitempty"` // 注入的内容 RouterPath@functionName@RouterString;
	StructName    string `json:"structName"`
	StructCNName  string `json:"structCNName"`
	ApiIDs        string `json:"apiIDs,omitempty"` // api表注册内容
	Flag          int    `json:"flag"`             // 表示对应状态 0 代表创建, 1 代表回滚 ...
}

// ToRequestIds ApiIDs 转换 IdsReq
// Author [SliverHorn](https://github.com/SliverHorn)
func (m *SysAutoCodeHistory) ToRequestIds() IdsReq {
	if m.ApiIDs == "" {
		return IdsReq{}
	}
	slice := strings.Split(m.ApiIDs, ";")
	ids := make([]int, 0, len(slice))
	length := len(slice)
	for i := 0; i < length; i++ {
		id, _ := strconv.ParseInt(slice[i], 10, 32)
		ids = append(ids, int(id))
	}
	return IdsReq{Ids: ids}
}

type SysAutoHistory struct {
	PageInfo
}

// GetById Find by id structure
type RollBack struct {
	ID          int  `json:"id" form:"id"`                   // 主键ID
	DeleteTable bool `json:"deleteTable" form:"deleteTable"` // 是否删除表
}
type AutoCodeHistory struct {
	ID           uint      `json:"ID" gorm:"column:id"`
	CreatedAt    time.Time `json:"CreatedAt" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"UpdatedAt" gorm:"column:updated_at"`
	BusinessDB   string    `json:"businessDB" gorm:"column:business_db"`
	TableName    string    `json:"tableName" gorm:"column:table_name"`
	StructName   string    `json:"structName" gorm:"column:struct_name"`
	StructCNName string    `json:"structCNName" gorm:"column:struct_cn_name"`
	Flag         int       `json:"flag" gorm:"column:flag"`
}
