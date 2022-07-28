package utils

import "github.com/spf13/viper"

// Conf holds the app configuration.
type Conf struct {
	DataSourceName string `mapstructure:"DATA_SOURCE_NAME"`
	ServerPort     string `mapstructure:"SERVER_PORT"`
	ApiKey         string `mapstructure:"API_KEY"`
}

// newConfig creates and initializes the configuration structure.
func newConfig() (config *Conf) {
	config = &Conf{}

	config.DataSourceName = "dataSource"
	config.ServerPort = "port"
	config.ApiKey = "key"

	return config
}

// LoadConfig initializes the configuration from a file or
// using the environment variables. The path specifies the
// path to the config file with the name "measurement.env".
func LoadConfig(path string) (config *Conf, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("measurement")
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
