package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	HTTPServer HTTPServer `mapstructure:"httpServer"`
	GRPCServer GRPCServer `mapstructure:"grpcServer"`
}

type HTTPServer struct {
	Port       string `mapstructure:"port"`
	JWTSignKey string `mapstructure:"jwtSignKey"`
}

type GRPCServer struct {
	Port string `mapstructure:"port"`
}

func Load() (*config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return &cfg, nil
}
