package env

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv         string `mapstructure:"APP_ENV"`
	DBDriver       string `mapstructure:"DB_DRIVER"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         int    `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	DBSSLMode      string `mapstructure:"DB_SSL_MODE"`
	DBMaxOpenConns int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBMaxIdleTime  string `mapstructure:"DB_MAX_IDLE_TIME"`
	AppPort        string `mapstructure:"PORT"`
}

func LoadConfig() (config Config, err error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // Default environment
	}

	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName(env)

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
	viper.WatchConfig()

	return
}
