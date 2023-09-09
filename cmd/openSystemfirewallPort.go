package cmd

import (
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model/request"
	"TTPanel/internal/service"
	"fmt"
	"github.com/spf13/cobra"
)

var systemFirewallPort int
var systemFirewallProtocol string

// openSystemFirewallPortCmd represents the setPort command
var openSystemFirewallPortCmd = &cobra.Command{
	Use:   "openSystemFirewallPort",
	Short: "开放系统防火墙端口",
	Long:  `使用该命令开放系统防火墙端口，如：panel tools openSystemFirewallPort -p 8888 -t tcp`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("开放系统防火墙端口：%d \n", systemFirewallPort)
		if systemFirewallPort < 1 || systemFirewallPort > 65535 {
			fmt.Printf("错误格式,端口范围1-65535\n")
			return
		}
		if systemFirewallProtocol != "tcp" && systemFirewallProtocol != "udp" && systemFirewallProtocol != "tcp/udp" {
			fmt.Printf("错误格式,协议只能为tcp或udp或tcp/udp\n")
			return
		}
		var ServiceGroupApp = service.GroupApp
		firewall, err := ServiceGroupApp.SystemFirewallServiceApp.New()
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		param := &request.CreatePortRuleR{
			Port:     systemFirewallPort,
			Strategy: constant.SystemFirewallStrategyAllow,
			SourceIp: "",
			Ps:       "命令行添加",
			Protocol: systemFirewallProtocol,
		}
		err = firewall.BatchCreatePortRule([]*request.CreatePortRuleR{param})
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		fmt.Printf("开放系统防火墙端口成功\n")
	},
}

func init() {
	toolsCmd.AddCommand(openSystemFirewallPortCmd)
	openSystemFirewallPortCmd.Flags().IntVarP(&systemFirewallPort, "port", "p", 8888, "需要设置的端口")
	openSystemFirewallPortCmd.Flags().StringVarP(&systemFirewallProtocol, "systemFirewallProtocol", "t", "tcp", "协议")
	_ = openSystemFirewallPortCmd.MarkFlagRequired("port")
}
