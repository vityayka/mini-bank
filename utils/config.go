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
			log.Println("driver missing")
			return
		}
		log.Println("dbdriver: ", config.DBDriver)
		if config.DBSource, envVarExists = os.LookupEnv("DB_SOURCE"); !envVarExists {
			log.Println("source missing")
			return
		}
		log.Println("db src: ", config.DBSource)
		if config.ServerAddress, envVarExists = os.LookupEnv("SERVER_ADDRESS"); !envVarExists {
			log.Println("srv addr missing")
			return
		}
		log.Println("srv addr: ", config.ServerAddress)
		if config.TokenSymmetricKey, envVarExists = os.LookupEnv("TOKEN_SYMMETRIC_KEY"); !envVarExists {
			log.Println("token key missing")
			return
		}
		log.Println("sym key: ", config.TokenSymmetricKey)
		if accessTokenDuration, envVarExists := os.LookupEnv("ACCESS_TOKEN_DURATION"); envVarExists {
			log.Println("token dur: ", config.AccessTokenDuration)
			config.AccessTokenDuration, err = time.ParseDuration(accessTokenDuration)
			return
		}
		log.Println("token duration missing")
		return
	}

	err = viper.Unmarshal(&config)
	return
}
