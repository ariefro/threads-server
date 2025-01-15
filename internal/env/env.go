package env

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	}
	Database struct {
		Driver      string `mapstructure:"driver"`
		Host        string `mapstructure:"host"`
		Port        int    `mapstructure:"port"`
		User        string `mapstructure:"user"`
		Password    string `mapstructure:"password"`
		DBName      string `mapstructure:"dbname"`
		SSLMode     string `mapstructure:"sslmode"`
		MaxOpenConn int    `mapstructure:"maxopenconn"`
		MaxIdleConn int    `mapstructure:"maxidleconn"`
		MaxIdleTime string `mapstructure:"maxidletime"`
	}
}

func LoadConfig() (config Config, err error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // Default environment
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
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
