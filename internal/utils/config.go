package utils

import "github.com/spf13/viper"

type Config struct {
	Env          string `mapstructure:"ENV"`
	ServerPort   string `mapstructure:"SERVER_PORT"`
	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
