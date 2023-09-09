package cmd

import (
	"TTPanel/internal/global"
	"TTPanel/pkg/util"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var panelPort int

// setPortCmd represents the setPort command
var setPortCmd = &cobra.Command{
	Use:   "setPort",
	Short: "修改面板端口",
	Long: `使用该命令修改面板的端口，如：TTPanel tools setPort -p 8080
警告：设置过程不检查端口使用情况，请勿使用常用端口，如21、22、80、443、3306、8080等，如果是阿里云、腾讯云等自带防火墙/安全组的服务商，修改前确认已开启对应端口，否则可能导致面板无法访问`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("设置端口：%d \n", panelPort)
		if panelPort < 1 || panelPort > 65535 {
			fmt.Printf("错误格式,端口范围1-65535\n")
			return
		}
		//开始写入配置文件
		global.Config.System.PanelPort = panelPort
		newConfig := global.Config.System
		global.Vp.Set("system", newConfig)
		err := global.Vp.WriteConfig() // 保存配置文件
		if err != nil {
			fmt.Printf("保存配置文件失败：%s\n", err)
			os.Exit(1)
		}
		fmt.Printf("端口修改成功,当前端口：%d\n", panelPort)
		_, err = util.ExecShell("tt restart")
		if err != nil {
			fmt.Printf("重启面板失败：%s\n", err)
			return
		}
	},
}

func init() {
	toolsCmd.AddCommand(setPortCmd)
	setPortCmd.Flags().IntVarP(&panelPort, "port", "p", 8888, "需要设置的端口")
}
