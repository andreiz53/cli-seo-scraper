/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Adds a config",
	Long: `This command helps you set the config for the application.
The config is used in reading the websites that need to be scraped as well as different configuration settings.`,
	Run: setAppConfig}

func init() {
	rootCmd.AddCommand(configCmd)

}

type AppConfig struct {
	Websites []string `mapstructure:"websites"`
}

func LoadAppConfig(path string) (*AppConfig, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var config AppConfig
	err = viper.Unmarshal(&config)
	return &config, err
}

func setAppConfig(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.PrintErr(`The config command accepts one argument, which is the path to the config file. 

Example: cli-go-scraper config ./config.json
`)
		return
	}
	config, err := LoadAppConfig(args[0])
	if err != nil {
		cmd.PrintErr(`Could not load config at path `, args[0], "\n")
		return
	}
	fmt.Println("got config", config)
}
