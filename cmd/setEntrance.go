package cmd

import (
	"TTPanel/internal/global"
	"TTPanel/pkg/util"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var entrance string

// setEntranceCmd represents the setEntrance command
var setEntranceCmd = &cobra.Command{
	Use:   "setEntrance",
	Short: "设置面板登陆入口",
	Long:  `使用该命令修改面板的登陆入口，以 '/' 开头，如：TTPanel tools setEntrance -e /yourEntrance`,
	Run: func(cmd *cobra.Command, args []string) {
		if util.StrIsEmpty(entrance) {
			fmt.Printf("-e参数为空\n")
			return
		}
		if strings.HasPrefix(entrance, "/") == false {
			fmt.Printf("错误格式,入口必须以 '/' 开头\n")
			return
		}
		//entrance验证是否是数字、字母、下划线组成
		for index, v := range entrance {
			if index == 0 {
				continue
			}
			if (v < '0' || v > '9') && (v < 'a' || v > 'z') && (v < 'A' || v > 'Z') && v != '_' {
				fmt.Printf("错误格式,入口由数字、字母、下划线组成\n")
				return
			}
		}
		//entrance长度大于4且小于16
		if len(entrance) < 4 || len(entrance) > 16 {
			fmt.Printf("长度必须4-16之间\n")
			return
		}
		//开始写入配置文件
		global.Config.System.Entrance = entrance
		newConfig := global.Config.System
		global.Vp.Set("system", newConfig)
		err := global.Vp.WriteConfig() // 保存配置文件
		if err != nil {
			fmt.Printf("保存配置文件失败：%s\n", err)
			os.Exit(1)
		}
		fmt.Printf("入口设置成功，当前入口为：%s\n", entrance)
		_, err = util.ExecShell("tt restart")
		if err != nil {
			fmt.Printf("重启面板失败：%s\n", err)
			return
		}
	},
}

func init() {
	toolsCmd.AddCommand(setEntranceCmd)
	setEntranceCmd.Flags().StringVarP(&entrance, "entrance", "e", "", "需要设置的入口,以 '/' 开头")
}
