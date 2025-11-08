package config

import (
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/spf13/viper"
)

// AppConfig - struct of app config.
// Add more config here if any new config been added in config.yaml.
type AppConfig struct {
	Secrets struct {
		HmacSecret string `destructure:"HMAC_SECRET"`
	} `destructure:"secrets"`

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
	} `destructure:"redis"`
}

// global singleton config
var (
	config     *AppConfig
	configOnce sync.Once
	configErr  error
)

// App - safe accessor that lazy load default config if not loaded yet
func App() *AppConfig {
	configOnce.Do(func() {
		configErr = load("./config")
	})

	if configErr != nil {
		panic(fmt.Sprintf("config not loaded: %v", configErr))
	}

	return config
}

// LoadConfig - Load config YAML in deliver path.
// path: config folder path
func LoadConfig(path string) error {
	configOnce.Do(func() {
		configErr = load(path)
	})

	return configErr
}

func load(path string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config error: %w", err)
	}

	for _, key := range viper.AllKeys() {
		val := viper.GetString(key)
		if containsBraceEnv(val) {
			val = expandEnvWithDefault(val)
		} else {
			val = os.ExpandEnv(val)
		}
		viper.Set(key, val)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("unmarshal config error: %w", err)
	}

	fmt.Printf("Config loaded: %+v\n", config)
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
