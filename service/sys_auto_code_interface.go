package service

import (
	"go-admin/global"
	"go-admin/model"
)

type DatabaseStruct interface {
	GetDB(businessDB string) (data []model.Db, err error)
	GetTables(businessDB string, dbName string) (data []model.Table, err error)
	GetColumn(businessDB string, tableName string, dbName string) (data []model.Column, err error)
}

func Database(businessDB string) DatabaseStruct {

	if businessDB == "" {
		switch global.GVA_CONFIG.System.DbType {
		case "mysql":
			return AutoCodeMysql
		// case "pgsql":
		// 	return AutoCodePgsql
		default:
			return AutoCodeMysql
		}
	} else {
		for _, info := range global.GVA_CONFIG.DBList {
			if info.AliasName == businessDB {

				switch info.Type {
				case "mysql":
					return AutoCodeMysql
				// case "mssql":
				// 	return AutoCodeMssql
				// case "pgsql":
				// 	return AutoCodePgsql
				// case "oracle":
				// 	return AutoCodeOracle
				default:
					return AutoCodeMysql
				}
			}
		}
		return AutoCodeMysql
	}

}
