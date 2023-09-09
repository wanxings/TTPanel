package cmd

import (
	"TTPanel/internal/service"
	"TTPanel/pkg/util"
	"fmt"

	"github.com/spf13/cobra"
)

var rootPwd string

// setRootPwdCmd represents the setRootPwd command
var setRootPwdCmd = &cobra.Command{
	Use:   "setRootPwd",
	Short: "设置mysql root密码",
	Long:  `使用该命令修改mysql root密码，如：TTPanel mysql setRootPwd -p yourPassword`,
	Run: func(cmd *cobra.Command, args []string) {
		if util.StrIsEmpty(rootPwd) {
			fmt.Printf("参数不能为空\n")
			return
		}
		err := service.GroupApp.DatabaseMysqlServiceApp.ChangeLocalRootPassword(rootPwd)
		if err != nil {
			fmt.Printf("修改失败：%s\n", err)
			return
		}
		return
	},
}

func init() {
	mysqlCmd.AddCommand(setRootPwdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setRootPwdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setRootPwdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setRootPwdCmd.Flags().StringVarP(&rootPwd, "rootPwd", "p", "", "需要设置的mysql root密码")
	_ = setRootPwdCmd.MarkFlagRequired("rootPwd")
}
