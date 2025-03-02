/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"cli-seo-scraper/cfg"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your config file",
	Long:  `The init command starts the process of creating your config file.`,
	Run:   handleInitConfig}

func init() {
	rootCmd.AddCommand(initCmd)
}

func handleInitConfig(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	var configName string
	for {
		cmd.Print("Enter the name for your config (example: config.json): ")
		configName, _ = reader.ReadString('\n')
		configName = strings.TrimSpace(configName)

		if strings.HasSuffix(configName, ".json") {
			break
		}
		cmd.Println("Config name provided is not a .json file")
	}
	var outputFilename string
	for {
		cmd.Print("Enter the name for your output file (example: output.csv): ")
		outputFilename, _ = reader.ReadString('\n')
		outputFilename = strings.TrimSpace(outputFilename)

		if strings.HasSuffix(outputFilename, ".csv") {
			break
		}

		cmd.Println("Output file name provided is not a .csv file")
	}

	websites := []string{}
	cmd.Println("Great! Now let's add some websites to scrape (example: https://google.com)")
	for {
		cmd.Print("Enter website URL (or 'done' to finish):")
		website, _ := reader.ReadString('\n')
		website = strings.TrimSpace(website)

		if website == "done" {
			break
		}
		if !strings.Contains(website, "http") {
			cmd.Println("Please add http:// or https://")
			continue
		}

		websites = append(websites, website)
	}
	appConfig := cfg.NewAppConfig(websites, configName, outputFilename)
	cfgFile, err := os.Create(appConfig.ConfigFilename)
	if err != nil {
		cmd.Println("Could not create config with name", appConfig.ConfigFilename)
		return
	}
	defer cfgFile.Close()

	encoder := json.NewEncoder(cfgFile)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(appConfig)
	if err != nil {
		cmd.Println("Could not write to config file:", err)
		return
	}

	cmd.Println("Configuration saved successfully to", appConfig.ConfigFilename)
}
