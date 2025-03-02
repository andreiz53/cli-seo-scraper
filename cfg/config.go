package cfg

import "github.com/spf13/viper"

type AppConfig struct {
	Websites       []string `mapstructure:"websites" json:"websites"`
	ConfigFilename string   `mapstructure:"config_filename" json:"config_filename"`
	OutputFilename string   `mapstructure:"output_filename" json:"output_filename"`
}

func NewAppConfig(websites []string, config string, output string) *AppConfig {
	return &AppConfig{
		Websites:       websites,
		ConfigFilename: config,
		OutputFilename: output,
	}
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
