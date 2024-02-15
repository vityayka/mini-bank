package utils

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DBDriver            string        `mapstructure:"DB _DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig reads configuration from environment file or variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// maybe the config envVarExists in the Environment, try it:
		var envVarExists = false
		if config.DBDriver, envVarExists = os.LookupEnv("DB_DRIVER"); !envVarExists {
			return
		}
		log.Println("dbdriver: ", config.DBDriver)
		if config.DBSource, envVarExists = os.LookupEnv("DB_SOURCE"); !envVarExists {
			return
		}
		log.Println("db src: ", config.DBSource)
		if config.ServerAddress, envVarExists = os.LookupEnv("SERVER_ADDRESS"); !envVarExists {
			return
		}
		log.Println("srv addr: ", config.ServerAddress)
		if config.TokenSymmetricKey, envVarExists = os.LookupEnv("TOKEN_SYMMETRIC_KEY"); !envVarExists {
			return
		}
		log.Println("sym key: ", config.TokenSymmetricKey)
		if accessTokenDuration, envVarExists := os.LookupEnv("ACCESS_TOKEN_DURATION"); envVarExists {
			config.AccessTokenDuration, err = time.ParseDuration(accessTokenDuration)
			return
		}
		return
	}

	err = viper.Unmarshal(&config)
	return
}
