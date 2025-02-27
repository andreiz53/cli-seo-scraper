package cfg

import "github.com/spf13/viper"

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
