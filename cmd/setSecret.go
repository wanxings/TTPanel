package cmd

import (
	"TTPanel/internal/global"
	"TTPanel/pkg/util"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var secretStr string

// setSecretCmd represents the setSecret command
var setSecretCmd = &cobra.Command{
	Use:   "setSecret",
	Short: "设置密钥",
	Long:  `使用该命令修改面板的jwt、session密钥，秘钥长度32位，如：TTPanel tools setSecret -s yourSecret ,如果不设置秘钥，面板会自动生成秘钥`,
	Run: func(cmd *cobra.Command, args []string) {
		//判断长度是否是32位
		if len(secretStr) != 32 {
			secretStr = string(util.RandStr(32, util.ALL))
		}

		//开始写入配置文件
		global.Config.System.JwtSecret = secretStr
		global.Config.System.SessionSecret = secretStr
		newConfig := global.Config.System
		global.Vp.Set("system", newConfig)
		err := global.Vp.WriteConfig() // 保存配置文件
		if err != nil {
			fmt.Printf("保存配置文件失败：%s\n", err)
			os.Exit(1)
		}
		fmt.Printf("秘钥修改成功,当前秘钥：%s\n", secretStr)
		_, err = util.ExecShell("tt restart")
		if err != nil {
			fmt.Printf("重启面板失败：%s\n", err)
			return
		}
	},
}

func init() {
	toolsCmd.AddCommand(setSecretCmd)
	setSecretCmd.Flags().StringVarP(&secretStr, "secretStr", "s", "", "需要设置的秘钥")
}
