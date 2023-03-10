package main

import (
	"go-admin/global"
	"go-admin/initialize"
)

func main() {
	// conf.Init()
	// dao.SystemInit()
	// router.NewRouter()
	global.GVA_VP = initialize.InitViper() // 初始化Viper
	global.GVA_DB = initialize.Gorm()      // gorm连接数据库
	initialize.OtherInit()                 //初始global.BlackCache
	initialize.InitRouter()                //初始总路由
	// initialize.Timer()
	// initialize.DBList()
	// if global.GVA_DB != nil {
	// 	initialize.RegisterTables(global.GVA_DB) // 初始化表
	// 	// 程序结束前关闭数据库链接
	// 	db, _ := global.GVA_DB.DB()
	// 	defer db.Close()
	// }
	// core.RunWindowsServer()
}
