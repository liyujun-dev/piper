package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "piper",
	Short: "Piper - CI/CD CLI Tool",
	Long:  `A CLI tool to manage CI/CD pipelines across multiple platforms.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
