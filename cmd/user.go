package cmd

import (
	"TTPanel/internal/service"
	"fmt"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "获取面板用户名",
	Long:  `使用该命令获取面板用户名，如：TTPanel tools user`,
	Run: func(cmd *cobra.Command, args []string) {
		user, err := service.GroupApp.UserServiceApp.Info(1)
		if err != nil {
			fmt.Printf("获取用户信息失败：%s\n", err)
			return
		}
		fmt.Printf("%s", user.Username)
	},
}

func init() {
	toolsCmd.AddCommand(userCmd)
}
