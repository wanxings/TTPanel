package cmd

import (
	"TTPanel/internal/global"
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "版本号",
	Long:  `当前面板的版本号`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("TT-Panel Version: %s-%s\n", global.Version, global.Config.System.PreReleaseVersion)
	},
}

func init() {
	toolsCmd.AddCommand(versionCmd)
}
