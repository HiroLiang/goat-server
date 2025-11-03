package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/viper"
)

// AppConfig - struct of app config.
// Add more config here if any new config been added in config.yaml.
type AppConfig struct {
	Database map[string]struct {
		Driver string `mapstructure:"driver"`
		Dsn    string `mapstructure:"dsn"`
		Config *struct {
			MaxOpenConns    int `destructure:"max_open_conns"`
			MaxIdleConns    int `destructure:"max_idle_conns"`
			ConnMaxLifetime int `destructure:"conn_max_lifetime"`
			ConnMaxIdleTime int `destructure:"conn_max_idle_time"`
		} `mapstructure:"config"`
	} `mapstructure:"databases"`
	Redis struct {
		Addr     string `destructure:"addr"`
		Password string `destructure:"password"`
		DB       int    `destructure:"db"`
	}
}

// Cfg - global config
var Cfg AppConfig

// LoadConfig - Load config YAML in deliver path.
// path: config folder path
func LoadConfig(path string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config error: %w", err)
	}

	for _, key := range viper.AllKeys() {
		val := os.ExpandEnv(viper.GetString(key))
		viper.Set(key, val)
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		return fmt.Errorf("unmarshal config error: %w", err)
	}

	fmt.Printf("Config loaded: %+v\n", Cfg)
	return nil
}

func expandEnvWithDefault(s string) string {
	re := regexp.MustCompile(`\$\{([^:}]+)(?::([^}]+))?}`)
	return re.ReplaceAllStringFunc(s, func(sub string) string {
		matches := re.FindStringSubmatch(sub)
		if len(matches) >= 2 {
			key := matches[1]
			def := ""
			if len(matches) >= 3 {
				def = matches[2]
			}
			if val, ok := os.LookupEnv(key); ok && val != "" {
				return val
			}
			return def
		}
		return sub
	})
}

func containsBraceEnv(s string) bool {
	return regexp.MustCompile(`\$\{[^}]+}`).MatchString(s)
}
