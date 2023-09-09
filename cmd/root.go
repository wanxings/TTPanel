package cmd

import (
	"TTPanel/internal/global"
	"TTPanel/internal/initialize"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "TTPanel",
	Short: "Linux管理面板",
	Long:  `方便快捷的Linux管理面板`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	global.Version = "0.1.0"
	//初始化Viper
	global.Vp = initialize.InitViper() // 初始化Viper
	//初始化调试日志
	global.Log = initialize.InitLogger()
	//初始化i18n
	global.I18n = initialize.InitI18n()
	//初始化go-cache缓存
	global.GoCache = initialize.InitGoCache()
	//初始化面板数据库
	global.PanelDB = initialize.InitPanelDB(global.Config.Sqlite) // gorm连接数据库
	//初始化TTWaf数据库
	global.TTWafDB = initialize.InitTTWafDB(global.Config.Sqlite) // gorm连接数据库
	//数据库Migrate
	initialize.InitMigrate() // gormAutoMigrate
	//if global.Config.System.RunMode != gin.DebugMode {
	//	initialize.InitMigrate() // gormAutoMigrate
	//}
}
