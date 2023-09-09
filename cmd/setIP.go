package cmd

import (
	"TTPanel/internal/service"
	"TTPanel/pkg/util"
	"fmt"
	"github.com/spf13/cobra"
)

var panelIP string

// setIPCmd represents the setIP command
var setIPCmd = &cobra.Command{
	Use:   "setIP",
	Short: "修改面板IP",
	Long: `使用该命令修改面板的IP，如：TTPanel tools setIP -i 127.0.0.1
该IP不涉及其他功能，仅用作首页显示以及某些扩展程序的链接使用，不影响面板的正常使用`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("设置面板IP：%s \n", panelIP)
		if util.GetIPType(panelIP) == 0 {
			fmt.Printf("IP格式错误\n")
			return
		}
		err := (&service.SettingsService{}).SetPanelIP(panelIP)
		if err != nil {
			fmt.Printf("修改失败,err:%s\n", err.Error())
			return
		}
		fmt.Printf("IP修改成功,当前面板IP：%s\n", panelIP)
	},
}

func init() {
	toolsCmd.AddCommand(setIPCmd)
	setIPCmd.Flags().StringVarP(&panelIP, "port", "i", "127.0.0.1", "需要设置的面板IP")
}
