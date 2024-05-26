package configs

import "github.com/spf13/viper"

type Config struct {
	DBProtocol     string `mapstructure:"DB_PROTOCOL"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBDatabase     string `mapstructure:"DB_DATABASE"`
	LimitByToken   bool   `mapstructure:"LIMIT_BY_TOKEN"`
	RateLimit      int    `mapstructure:"RATE_LIMIT"`
	ExpirationTime int    `mapstructure:"EXPIRATION_TIME"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("dev")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
