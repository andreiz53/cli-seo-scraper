/*
Copyright © 2025 Andrei Zamfira <andrei.zamfira53@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli-seo-scraper",
	Short: "CLI SEO Scraper allows you to easily scrape websites and find SEO issues",
	Long: `This CLI tool's purpose is to help people that have to mantain multiple websites.

Based on the given config file, it scrapes the websites and gives you all SEO related information you need.
Run the help command to learn more.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli-seo-scraper.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
