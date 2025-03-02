/*
Copyright © 2025 Andrei Zamfira <andrei.zamfira53@gmail.com>
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

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate SEO stats based on the generated config.",
	Long:  ``,
	Run:   generateStats}

func init() {
	rootCmd.AddCommand(generateCmd)
}

func generateStats(cmd *cobra.Command, args []string) {
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
	outFile, err := os.Create(scraperCfg.OutputFilename)
	if err != nil {
		cmd.Println(colors.Error("Could not create output file"))
		cmd.Println(err)
		return
	}
	writer := csv.NewWriter(outFile)
	err = writer.Write(seo.CSVHeader())
	if err != nil {
		cmd.Println(colors.Error("Could not write into the output file"))
		cmd.Println(err)
		return
	}
	c := scraper.NewCollector()
	scr := scraper.NewScraper(c, scraperCfg)
	t := time.Now()
	settings := scr.ScrapeSEO()
	endT := time.Since(t)
	for _, s := range settings {
		err = writer.Write(s.ToCSVLine())
		if err != nil {
			cmd.Println(colors.Error("Could not write line into output file"))
			cmd.Println(err)
			return
		}
	}
	writer.Flush()
	cmd.Println(colors.Success("Websites scraped successfully in ", endT))
}
