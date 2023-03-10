package initialize

import (
	"fmt"
	"go-admin/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 读取配置文件，并转换成 struct结构
func InitViper() *viper.Viper {
	//    //获取项目的执行路径
	//    path, err := os.Getwd()
	//    if err != nil {
	// 	   panic(err)
	//    }
	v := viper.New()
	//v.AddConfigPath(path + "/config")//设置读取的文件路径
	//v.SetConfigName("application") //设置读取的文件名

	v.SetConfigFile("config.yaml") //config路径
	v.SetConfigType("yaml")        //设置文件的类型
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig() //监听

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.GVA_CONFIG); err != nil {
		fmt.Println(err)
	}
	fmt.Println("参数测试：", global.GVA_CONFIG.Mysql.Path)
	fmt.Println("参数测试：", global.GVA_CONFIG.Mysql.Dbname)
	fmt.Println("参数测试：", global.GVA_CONFIG.Mysql)
	return v
}
