package cmd

import (
	"fmt"
	"main/jira-tools/cmd/issue"
	"main/jira-tools/cmd/version"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(issue.IssueCmd)
	rootCmd.AddCommand(version.VersionCmd)
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
