package cmd

import (
	"TTPanel/internal/global"
	"fmt"

	"github.com/spf13/cobra"
)

// portCmd represents the port command
var portCmd = &cobra.Command{
	Use:   "port",
	Short: "获取面板端口",
	Long:  `使用该命令获取面板端口，如：panel tools port`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(global.Config.System.PanelPort)
	},
}

func init() {
	toolsCmd.AddCommand(portCmd)
}
