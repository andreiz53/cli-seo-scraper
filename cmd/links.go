/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"time"

	"github.com/spf13/cobra"

	"cli-seo-scraper/colors"
	"cli-seo-scraper/config"
	"cli-seo-scraper/scraper"
	"cli-seo-scraper/seo"
)

var linksCmd = &cobra.Command{
	Use:   "links",
	Short: "The 'links' command generates a report of broken links",
	Long:  ``,
	Run:   handleLinksCmd}

func init() {
	rootCmd.AddCommand(linksCmd)
}

func handleLinksCmd(cmd *cobra.Command, args []string) {

	cfgPath, err := config.GetAppConfig()
	if err != nil {
		cmd.Println(colors.Warning("No config file found. Please generate one using the 'init' command"))
		return
	}
	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		cmd.Println(colors.Error("Could not open the application config file"))
		cmd.Println(err)
		return
	}
	var scraperCfg scraper.ScraperConfig
	err = json.NewDecoder(cfgFile).Decode(&scraperCfg)
	if err != nil {
		cmd.Println(colors.Error("Scraper config file has invalid format"))
		cmd.Println(err)
		return
	}

	outFile, err := os.Create(scraperCfg.OutputSEOMetas)
	if err != nil {
		cmd.Println(colors.Error("Could not create output file"))
		cmd.Println(err)
		return
	}

	writer := csv.NewWriter(outFile)
	err = writer.Write(seo.CSVHeaderLinks())
	if err != nil {
		cmd.Println(colors.Error("Could not write into the output file"))
		cmd.Println(err)
		return
	}

	c := scraper.NewCollector()
	scr := scraper.NewScraper(c, scraperCfg)
	t := time.Now()
	brokenLinks := scr.ScrapeLinks()
	endT := time.Since(t)

	for _, websiteLinks := range brokenLinks {
		for _, link := range websiteLinks.Links {
			err = writer.Write(link.ToCSVLine())
			if err != nil {
				cmd.Println(colors.Error("Could not write into output file"))
				cmd.Println(err)
				return
			}
		}
	}

	writer.Flush()
	cmd.Println(colors.Success("Broken links scraped successfully in ", endT))
}
