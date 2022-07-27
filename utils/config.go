package utils

import "github.com/spf13/viper"

type Conf struct {
	DataSourceName string `mapstructure:"DATA_SOURCE_NAME"`
	ServerPort     string `mapstructure:"SERVER_PORT"`
	ApiSource      string `mapstructure:"API_KEY"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config *Conf, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	config = &Conf{}
	err = viper.Unmarshal(config)
	return
}
