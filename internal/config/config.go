package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Keycloak KeycloakConfig
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type KeycloakConfig struct {
	URL          string
	Realm        string
	ClientID     string
	ClientSecret string
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("RESELL")

	// Defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "resell")
	viper.SetDefault("database.sslmode", "disable")

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("keycloak.url", "http://localhost:8180")
	viper.SetDefault("keycloak.realm", "resell")
	viper.SetDefault("keycloak.client_id", "resell-api")
	viper.SetDefault("keycloak.client_secret", "")

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	cfg.Server.Port = viper.GetInt("server.port")
	cfg.Server.ReadTimeout = viper.GetDuration("server.read_timeout")
	cfg.Server.WriteTimeout = viper.GetDuration("server.write_timeout")

	cfg.Database.Host = viper.GetString("database.host")
	cfg.Database.Port = viper.GetInt("database.port")
	cfg.Database.User = viper.GetString("database.user")
	cfg.Database.Password = viper.GetString("database.password")
	cfg.Database.DBName = viper.GetString("database.dbname")
	cfg.Database.SSLMode = viper.GetString("database.sslmode")

	cfg.Redis.Host = viper.GetString("redis.host")
	cfg.Redis.Port = viper.GetInt("redis.port")
	cfg.Redis.Password = viper.GetString("redis.password")
	cfg.Redis.DB = viper.GetInt("redis.db")

	cfg.Keycloak.URL = viper.GetString("keycloak.url")
	cfg.Keycloak.Realm = viper.GetString("keycloak.realm")
	cfg.Keycloak.ClientID = viper.GetString("keycloak.client_id")
	cfg.Keycloak.ClientSecret = viper.GetString("keycloak.client_secret")

	return &cfg, nil
}
