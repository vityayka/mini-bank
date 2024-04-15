package utils

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	DBURI                string        `mapstructure:"DB_URI"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RedisAddr            string        `mapstructure:"REDIS_ADDR"`
	GmailName            string        `mapstructure:"GMAIL_NAME"`
	GmailFrom            string        `mapstructure:"GMAIL_FROM"`
	GmailAccPassword     string        `mapstructure:"GMAIL_APP_PASSWORD"`
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
		if config.DBSource, envVarExists = os.LookupEnv("DB_SOURCE"); !envVarExists {
			return
		}
		if config.HTTPServerAddress, envVarExists = os.LookupEnv("SERVER_ADDRESS"); !envVarExists {
			return
		}
		if config.TokenSymmetricKey, envVarExists = os.LookupEnv("TOKEN_SYMMETRIC_KEY"); !envVarExists {
			return
		}
		if accessTokenDuration, envVarExists := os.LookupEnv("ACCESS_TOKEN_DURATION"); envVarExists {
			config.AccessTokenDuration, err = time.ParseDuration(accessTokenDuration)
			return
		}
		if refreshTokenDuration, envVarExists := os.LookupEnv("REFRESH_TOKEN_DURATION"); envVarExists {
			config.RefreshTokenDuration, err = time.ParseDuration(refreshTokenDuration)
			return
		}
		return
	}

	err = viper.Unmarshal(&config)
	return
}
