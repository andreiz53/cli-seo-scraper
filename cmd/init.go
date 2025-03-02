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

	"cli-seo-scraper/colors"
	"cli-seo-scraper/config"
	"cli-seo-scraper/scraper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your config file",
	Long:  `The init command starts the process of creating your config file.`,
	Run:   handleInitCmd}

func init() {
	rootCmd.AddCommand(initCmd)
}

func handleInitCmd(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	var configName string
	for {
		cmd.Print(colors.Bold("Enter the name for your config (example: config.json): "))
		configName, _ = reader.ReadString('\n')
		configName = strings.TrimSpace(configName)

		if strings.HasSuffix(configName, ".json") {
			break
		}
		cmd.Println(colors.Error("Config name provided is not a .json file"))
	}
	var outputFilename string
	for {
		cmd.Print(colors.Bold("Enter the name for your output file (example: output.csv): "))
		outputFilename, _ = reader.ReadString('\n')
		outputFilename = strings.TrimSpace(outputFilename)

		if strings.HasSuffix(outputFilename, ".csv") {
			break
		}

		cmd.Println(colors.Error("Output file name provided is not a .csv file"))
	}

	websites := []string{}
	cmd.Println(colors.Info("Great! Now let's add some websites to scrape (example: https://google.com)"))
	for {
		cmd.Print(colors.Bold("Enter website URL (or 'done' to finish): "))
		website, _ := reader.ReadString('\n')
		website = strings.TrimSpace(website)

		if website == "done" {
			break
		}
		if !strings.Contains(website, "http") {
			cmd.Println(colors.Error("Please add http:// or https://"))
			continue
		}

		websites = append(websites, website)
	}
	scraperConfig := scraper.NewScraperConfig(websites, outputFilename)
	cfgFile, err := os.Create(configName)
	if err != nil {
		cmd.Println(colors.Error("Could not create config with name", configName))
		return
	}
	defer cfgFile.Close()

	encoder := json.NewEncoder(cfgFile)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(scraperConfig)
	if err != nil {
		cmd.Println(colors.Error("Could not write to config file:", err))
		return
	}

	cmd.Println(colors.Success("Configuration saved successfully to ", configName))
	// save the the config filename into some preset file
	appConfig := config.NewAppConfig(configName)
	err = appConfig.GenerateConfig()
	if err != nil {
		cmd.Println(colors.Error("Could not generate application config file"))
		cmd.Println(err)
	}
}
