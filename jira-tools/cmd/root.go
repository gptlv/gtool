package cmd

import (
	"fmt"
	ad "main/jira-tools/cmd/active-directory"
	"main/jira-tools/cmd/asset"
	"main/jira-tools/cmd/issue"
	writeoff "main/jira-tools/cmd/write-off"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(issue.IssueCmd)
	rootCmd.AddCommand(ad.ActiveDirectoryCmd)
	rootCmd.AddCommand(writeoff.WriteOffCmd)
	rootCmd.AddCommand(asset.AssetCmd)
}

var rootCmd = &cobra.Command{
	Use:   "jira-tools",
	Short: "CLI application for mundane tasks in jira",
	Long:  "A CLI application built with Cobra for various operations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
