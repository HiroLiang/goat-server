package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Database map[string]struct {
		Driver string `mapstructure:"driver"`
		Dsn    string `mapstructure:"dsn"`
	} `mapstructure:"databases"`
}

var Cfg AppConfig

func LoadConfig(path string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config error: %w", err)
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		return fmt.Errorf("unmarshal config error: %w", err)
	}

	fmt.Printf("Config loaded: %+v\n", Cfg)
	return nil
}
