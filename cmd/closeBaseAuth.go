package cmd

import (
	"TTPanel/internal/global"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// closeBaseAuthCmd represents the closeBaseAuth command
var closeBaseAuthCmd = &cobra.Command{
	Use:   "closeBaseAuth",
	Short: "关闭BaseAuth认证",
	Long:  `使用该命令关闭BaseAuth认证,如：panel tools closeBaseAuth`,
	Run: func(cmd *cobra.Command, args []string) {
		//开始写入配置文件
		global.Config.System.BasicAuth.Status = false
		newConfig := global.Config.System
		global.Vp.Set("system", newConfig)
		err := global.Vp.WriteConfig() // 保存配置文件
		if err != nil {
			fmt.Printf("保存配置文件失败：%s\n", err)
			os.Exit(1)
		}
		fmt.Println("已关闭BaseAuth认证")
	},
}

func init() {
	toolsCmd.AddCommand(closeBaseAuthCmd)
}
