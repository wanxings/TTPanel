package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// toolsCmd represents the tools command
var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tools called")
	},
}

func init() {
	rootCmd.AddCommand(toolsCmd)

	//err := doc.GenMarkdownTree(rootCmd, "./docs")
	//if err != nil {
	//	log.Fatal(err)
	//}
}
