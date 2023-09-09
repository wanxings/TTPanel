package cmd

import (
	"TTPanel/internal/initialize"
	"TTPanel/internal/service"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run Panel",
	Long:  `Run Panel`,
	Run: func(cmd *cobra.Command, args []string) {
		//监听配置文件
		//initialize.WatchConfig()
		//启动队列协程
		initialize.TaskInit()
		//启动检测项目过期协程
		service.CheckProjectExpirationInit()
		//启动系统监控协程
		service.MonitorInit()
		//initialize.MonitorInit()
		//初始化插件
		//initialize.InitPlugin()
		//启动
		initialize.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
