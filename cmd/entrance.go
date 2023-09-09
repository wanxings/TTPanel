package cmd

import (
	"TTPanel/internal/global"
	"TTPanel/internal/service"
	"fmt"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "获取面板默认信息",
	Long:  `使用该命令获取面板默认信息，如：panel tools info`,
	Run: func(cmd *cobra.Command, args []string) {
		userInfo, err := service.GroupApp.UserServiceApp.Info(1)
		if err != nil {
			fmt.Printf("获取用户信息失败：%s\n", err)
			return
		}
		pool := "http"
		address := "IP"
		port := global.Config.System.PanelPort
		entrance := global.Config.System.Entrance
		user := userInfo.Username
		fmt.Printf("==================================================================\n")
		fmt.Printf("TT-Panel 默认信息!\n")
		fmt.Printf("==================================================================\n")
		fmt.Printf("外网面板地址: %s://%s:%d%s\n", pool, address, port, entrance)
		fmt.Printf("安装过程默认关闭Linux系统防火墙（可在面板-安全-系统防火墙中开启），若无法访问面板，请检查服务器商的 防火墙/安全组 是否有放行[%d]端口\n", port)
		fmt.Printf(" 如果未正常显示IP地址，在面板地址中手动加入IP地址即可\n")
		fmt.Printf("面板密码仅显示一次，请做好保存，后续无法获得密码，只能通过 tt 命令修改密码\n")
		fmt.Printf("==================================================================\n")

		fmt.Printf("==================================================================\n")
		fmt.Printf("TT-Panel default info!\n")
		fmt.Printf("==================================================================\n")
		fmt.Printf("External panel address: %s://%s:%d%s\n", pool, address, port, entrance)
		fmt.Printf("username: %s\n", user)
		fmt.Printf("If you cannot access the panel,\n")
		fmt.Printf("During the installation process, the Linux firewall is disabled by default. If you are unable to access the panel, please check if the server provider's firewall/security group has allowed access to the [%d] port.\n", port)
		fmt.Printf("If the IP address is not displayed correctly, you can manually add the IP address to the panel address.\n")
		fmt.Printf("The panel password is displayed only once. Please make sure to save it. You will not be able to obtain the password later. You can only modify the password using the tt command.\n")
		fmt.Printf("==================================================================\n")
	},
}

func init() {
	toolsCmd.AddCommand(infoCmd)
}
