package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type config struct {
	HTTPServer HTTPServer `mapstructure:"httpServer"`
	GRPCServer GRPCServer `mapstructure:"grpcServer"`
	JWT        JWT        `mapstructure:"jwt"`
	Mongo      Mongo      `mapstructure:"mongo"`
}

type HTTPServer struct {
	Port string `mapstructure:"port"`
}

type GRPCServer struct {
	Port string `mapstructure:"port"`
}

type JWT struct {
	SignKey string        `mapstructure:"signKey"`
	Expiry  time.Duration `mapstructure:"expiry"`
}

type Mongo struct {
	URI string `mapstructure:"uri"`
}

func Load() (*config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
