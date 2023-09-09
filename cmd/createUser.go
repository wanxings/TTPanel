package cmd

import (
	"TTPanel/internal/global"
	"TTPanel/internal/model"
	"TTPanel/internal/service"
	"fmt"
	"github.com/spf13/cobra"
)

var createUserName string
var createUserPwd string

// createUserCmd represents the setUser command
var createUserCmd = &cobra.Command{
	Use:   "createUser",
	Short: "创建用户",
	Long:  `使用该命令创建面板管理员账户，如：panel tools createUser -u yourUserName -p yourPwd`,
	Run: func(cmd *cobra.Command, args []string) {
		//判断长度在3-20位
		fmt.Printf("创建的用户名: %s \n", createUserName)
		fmt.Printf("创建的密码: %s \n", createUserPwd)
		if len(createUserName) < 3 || len(createUserName) > 20 {
			fmt.Printf("错误格式,用户名长度必须在3-20位\n")
			return
		}
		if len(createUserPwd) < 6 || len(createUserPwd) > 32 {
			fmt.Printf("错误格式,密码长度必须是6-32位\n")
			return
		}
		sPwd, salt := service.EncryptPasswordAndSalt(createUserPwd)
		create, err := (&model.User{
			ID:        1,
			Username:  createUserName,
			Password:  sPwd,
			LoginIp:   "127.0.0.1",
			LoginTime: 0,
			Email:     createUserName + "@admin.com",
			Salt:      salt,
		}).Create(global.PanelDB)
		if err != nil || create.ID == 0 {
			fmt.Printf("创建管理员账户失败.Error:%s\n", err)
			return
		}
		fmt.Printf("创建管理员账户成功\n")
	},
}

func init() {
	toolsCmd.AddCommand(createUserCmd)
	createUserCmd.Flags().StringVarP(&createUserName, "createUserName", "u", "", "需要创建的面板管理员用户名")
	createUserCmd.Flags().StringVarP(&createUserPwd, "createUserPwd", "p", "", "需要创建的面板密码")
}
