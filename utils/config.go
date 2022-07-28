package utils

import "github.com/spf13/viper"

type Conf struct {
	DataSourceName string `mapstructure:"DATA_SOURCE_NAME"`
	ServerPort     string `mapstructure:"SERVER_PORT"`
	ApiKey         string `mapstructure:"API_KEY"`
}

// 'constructor'
func newConfig() (config *Conf) {
	config = &Conf{}

	config.DataSourceName = "dataSource"
	config.ServerPort = "port"
	config.ApiKey = "key"

	return config
}

// LoadConfig reads configuration from file or environment variables.
// path specifies the path to "app.env"
func LoadConfig(path string) (config *Conf, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	config = newConfig()
	err = viper.Unmarshal(config)
	return
}
