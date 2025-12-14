package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
    Port              string `mapstructure:"PORT"`
    RedisAddr         string `mapstructure:"REDIS_ADDR"`
    JWTSecret         string `mapstructure:"JWT_SECRET"`
    AuthServiceURL    string `mapstructure:"AUTH_SERVICE_URL"`
    OrderServiceURL   string `mapstructure:"ORDER_SERVICE_URL"`
    PaymentServiceURL string `mapstructure:"PAYMENT_SERVICE_URL"`
    EventServiceURL   string `mapstructure:"EVENT_SERVICE_URL"`  
    MonolithURL       string `mapstructure:"MONOLITH_URL"`
    RateLimitReq      int    `mapstructure:"RATE_LIMIT_REQ"`
    RateLimitWindow   int    `mapstructure:"RATE_LIMIT_WINDOW"`
}

func LoadConfig() (cfg Config, err error) {
    viper.AddConfigPath(".")
    viper.SetConfigFile(".env")
    viper.AutomaticEnv()

    // Explicitly bind env vars for Docker
    viper.BindEnv("PORT")
    viper.BindEnv("REDIS_ADDR")
    viper.BindEnv("JWT_SECRET")
    viper.BindEnv("AUTH_SERVICE_URL")
    viper.BindEnv("ORDER_SERVICE_URL")
    viper.BindEnv("PAYMENT_SERVICE_URL")
    viper.BindEnv("EVENT_SERVICE_URL")  
    viper.BindEnv("MONOLITH_URL")
    viper.BindEnv("RATE_LIMIT_REQ")
    viper.BindEnv("RATE_LIMIT_WINDOW")

    err = viper.ReadInConfig()
    if err != nil {
        log.Println("No .env file found, using system environment variables")
    }
    
    err = viper.Unmarshal(&cfg)
    return
}