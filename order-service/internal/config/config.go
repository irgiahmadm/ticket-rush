package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
    Port        string `mapstructure:"PORT"`
    DatabaseURL string `mapstructure:"DATABASE_URL"`
}

func LoadConfig() (cfg Config, err error) {
    viper.AddConfigPath(".")
    viper.SetConfigFile(".env")
    viper.AutomaticEnv()

    viper.BindEnv("PORT")
    viper.BindEnv("DATABASE_URL")

    err = viper.ReadInConfig()
    if err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    err = viper.Unmarshal(&cfg)
    return
}