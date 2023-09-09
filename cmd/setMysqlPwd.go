package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setMysqlPwdCmd represents the setMysqlPwd command
var setMysqlPwdCmd = &cobra.Command{
	Use:   "setMysqlPwd",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setMysqlPwd called")
	},
}

func init() {
	toolsCmd.AddCommand(setMysqlPwdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setMysqlPwdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setMysqlPwdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
