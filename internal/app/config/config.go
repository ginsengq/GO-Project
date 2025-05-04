package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	Database struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`
	JWT struct {
		Secret string `mapstructure:"secret"`
		TTL    string `mapstructure:"ttl"`
	} `mapstructure:"jwt"`
	App struct {
		Environment string `mapstructure:"environment"`
		LogLevel    string `mapstructure:"log_level"`
	} `mapstructure:"app"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/app/config")

	viper.SetDefault("server.port", "8000")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.log_level", "debug")
	viper.AutomaticEnv()
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("jwt.secret", "JWT_SECRET")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}

	if cfg.Database.Host == "" || cfg.Database.User == "" || cfg.Database.DBName == "" {
		log.Fatal("Database host, user and dbname are required configurations")
	}

	if cfg.JWT.Secret == "" {
		log.Fatal("JWT secret is required")
	}

	return &cfg
}
