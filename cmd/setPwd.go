package cmd

import (
	"TTPanel/internal/service"
	"fmt"
	"github.com/spf13/cobra"
)

var panelPwd string

// setPwdCmd represents the setPwd command
var setPwdCmd = &cobra.Command{
	Use:   "setPwd",
	Short: "修改面板密码",
	Long:  `使用该命令修改面板的登录密码，密码长度6-32位，如：TTPanel tools setPwd -p yourPwd`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("设置密码: %s \n", panelPwd)
		//判断长度是否是6-32位
		if len(panelPwd) < 6 || len(panelPwd) > 32 {
			fmt.Printf("错误格式,密码长度必须是6-32位\n")
			return
		}
		err := service.GroupApp.SettingsServiceApp.SetUser("", panelPwd)
		if err != nil {
			fmt.Printf("修改失败：%s\n", err)
			return
		}
		fmt.Printf("密码修改成功,当前密码：%s\n", panelPwd)
	},
}

func init() {
	toolsCmd.AddCommand(setPwdCmd)
	setPwdCmd.Flags().StringVarP(&panelPwd, "panelPwd", "p", "", "需要设置的面板密码")
	_ = setPwdCmd.MarkFlagRequired("panelPwd")
}
