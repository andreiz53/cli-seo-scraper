package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	CONFIG_DIR               = "cli-seo-scraper"
	CONFIG_SETTINGS_FILENAME = "settings.json"
	CONFIG_IDENTIFIER        = "config_file"
)

type AppConfig struct {
	ScraperConfigFilename string
}

func NewAppConfig(scraperConfigFilename string) *AppConfig {
	return &AppConfig{
		ScraperConfigFilename: scraperConfigFilename,
	}
}

func (cfg *AppConfig) GenerateConfig() error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appConfigPath := filepath.Join(userConfigDir, CONFIG_DIR)
	err = os.MkdirAll(appConfigPath, os.ModePerm)
	if err != nil {
		return err
	}

	appConfigFullpath := filepath.Join(appConfigPath, CONFIG_SETTINGS_FILENAME)

	viper.Set(CONFIG_IDENTIFIER, cfg.ScraperConfigFilename)
	viper.SetConfigFile(appConfigFullpath)
	err = viper.WriteConfigAs(appConfigFullpath)
	if err != nil {
		return err
	}

	return nil
}

func GetAppConfig() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	viper.SetConfigFile(filepath.Join(userConfigDir, CONFIG_DIR, CONFIG_SETTINGS_FILENAME))

	err = viper.ReadInConfig()
	if err != nil {
		return "", err
	}

	configFilepath := viper.GetString(CONFIG_IDENTIFIER)
	if configFilepath == "" {
		return "", errors.New("could not get the application config file")
	}

	return configFilepath, nil
}
