package env

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv                  string  `mapstructure:"APP_ENV"`
	AppPort                 string  `mapstructure:"PORT"`
	AuthBasicUser           string  `mapstructure:"AUTH_BASIC_USER"`
	AuthBasicPass           string  `mapstructure:"AUTH_BASIC_PASS"`
	AuthTokenSecret         string  `mapstructure:"AUTH_TOKEN_SECRET"`
	AuthTokenExpired        int     `mapstructure:"AUTH_TOKEN_EXPIRED"`
	CorsAllowedOrigin       string  `mapstructure:"CORS_ALLOWED_ORIGIN"`
	DBDriver                string  `mapstructure:"DB_DRIVER"`
	DBHost                  string  `mapstructure:"DB_HOST"`
	DBPort                  int     `mapstructure:"DB_PORT"`
	DBUser                  string  `mapstructure:"DB_USER"`
	DBPassword              string  `mapstructure:"DB_PASSWORD"`
	DBName                  string  `mapstructure:"DB_NAME"`
	DBSSLMode               string  `mapstructure:"DB_SSL_MODE"`
	DBMaxOpenConns          int     `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns          int     `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBMaxIdleTime           string  `mapstructure:"DB_MAX_IDLE_TIME"`
	FrontendURL             string  `mapstructure:"FRONTEND_URL"`
	EmailSender             string  `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword     string  `mapstructure:"EMAIL_SENDER_PASSWORD"`
	RateLimiterRequestCount int     `mapstructure:"RATELIMITER_REQUESTS_COUNT"`
	RateLimiterEnabled      bool    `mapstructure:"RATE_LIMITER_ENABLED"`
	RedisAddress            string  `mapstructure:"REDIS_ADDRESS"`
	RedisPassword           string  `mapstructure:"REDIS_PASSWORD"`
	RedisDB                 int     `mapstructure:"REDIS_DB"`
	RedisEnabled            bool    `mapstructure:"REDIS_ENABLED"`
	SenderName              string  `mapstructure:"EMAIL_SENDER_NAME"`
	SentryDsn               string  `mapstructure:"SENTRY_DSN"`
	SentrySampleRate        float64 `mapstructure:"SENTRY_SAMPLE_RATE"`
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
