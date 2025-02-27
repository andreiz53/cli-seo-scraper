/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"os"

	"github.com/spf13/cobra"

	"cli-seo-scraper/cfg"
	"cli-seo-scraper/scraper"
	"cli-seo-scraper/seo"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate SEO stats based on the given config.",
	Long:  ``,
	Run:   generateStats}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("config", "c", "cli-seo-scraper-config.json", "Path to the config file. Defaults to cli-seo-scraper-config.json")
	generateCmd.Flags().StringP("output", "o", "cli-seo-scraper-output.csv", "Output file for SEO stats. Defaults to cli-seo-scraper-output.csv")

}

func generateStats(cmd *cobra.Command, args []string) {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		cmd.PrintErr("Invalid value for config flag \n")
		return
	}
	outputPath, err := cmd.Flags().GetString("output")
	if err != nil {
		cmd.PrintErr("Invalid value for output flag \n")
		return
	}
	config, err := cfg.LoadAppConfig(configPath)
	if err != nil {
		cmd.PrintErr("Could not load config file from given path\n")
		return
	}
	_ = outputPath
	file, err := os.Create(outputPath)
	if err != nil {
		cmd.PrintErr("Could not create output file\n")
		return
	}
	defer file.Close()
	w := csv.NewWriter(file)
	err = w.Write(seo.CSVHeader())
	if err != nil {
		cmd.PrintErr("Could not write to output file \n")
		return
	}

	scr := scraper.NewScraper()
	for _, website := range config.Websites {
		settings := scr.WithSEOSettings(website)
		err = w.Write(settings.ToCSVLine())
		if err != nil {
			cmd.PrintErr("Could not write SEO settings to output file \n")
			return
		}
	}
	w.Flush()
}
