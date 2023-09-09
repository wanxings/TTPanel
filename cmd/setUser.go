package cmd

import (
	"TTPanel/internal/service"
	"fmt"
	"github.com/spf13/cobra"
)

var userName string

// setUserCmd represents the setUser command
var setUserCmd = &cobra.Command{
	Use:   "setUser",
	Short: "修改用户名",
	Long:  `使用该命令修改面板的用户名，如：TTPanel tools setUser -u yourUserName`,
	Run: func(cmd *cobra.Command, args []string) {
		//判断长度在3-20位
		fmt.Printf("设置用户名: %s \n", userName)
		if len(userName) < 3 || len(userName) > 20 {
			fmt.Printf("错误格式,用户名长度必须在3-20位\n")
			return
		}
		err := service.GroupApp.SettingsServiceApp.SetUser(userName, "")
		if err != nil {
			fmt.Printf("修改失败：%s\n", err)
			return
		}
		fmt.Printf("用户名修改成功,当前用户名：%s\n", userName)
	},
}

func init() {
	toolsCmd.AddCommand(setUserCmd)
	setUserCmd.Flags().StringVarP(&userName, "userName", "u", "", "需要设置的面板用户名")
}
