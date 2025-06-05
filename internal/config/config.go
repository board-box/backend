package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App        AppConfig
	Postgres   PostgresConfig
	HTTP       HTTPConfig
	JWT        JWTConfig
	ChatApiKey string
}

type AppConfig struct {
	Name    string
	Host    string
	Version string
	Env     string // dev/stage/prod
}

type PostgresConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxConns        int
	MinConns        int
	MaxConnLifetime time.Duration
}

type HTTPConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

func New() (*Config, error) {
	var cfg Config

	_ = godotenv.Load()

	cfg.App = AppConfig{
		Name:    getEnv("APP_NAME", "myapp"),
		Version: getEnv("APP_VERSION", "1.0.0"),
		Env:     getEnv("APP_ENV", "dev"),
		Host:    getEnv("APP_HOST", "localhost"),
	}

	pgPort, err := strconv.Atoi(getEnv("PG_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid PG_PORT: %w", err)
	}

	pgMaxConns, err := strconv.Atoi(getEnv("PG_MAX_CONNS", "10"))
	if err != nil {
		return nil, fmt.Errorf("invalid PG_MAX_CONNS: %w", err)
	}

	pgMinConns, err := strconv.Atoi(getEnv("PG_MIN_CONNS", "2"))
	if err != nil {
		return nil, fmt.Errorf("invalid PG_MIN_CONNS: %w", err)
	}

	pgMaxConnLifetime, err := time.ParseDuration(getEnv("PG_MAX_CONN_LIFETIME", "1h"))
	if err != nil {
		return nil, fmt.Errorf("invalid PG_MAX_CONN_LIFETIME: %w", err)
	}

	cfg.Postgres = PostgresConfig{
		Host:            getEnv("PG_HOST", "localhost"),
		Port:            pgPort,
		User:            getEnv("PG_USER", "user"),
		Password:        getEnv("PG_PASSWORD", "password"),
		DBName:          getEnv("PG_DBNAME", "boardbox"),
		SSLMode:         getEnv("PG_SSLMODE", "disable"),
		MaxConns:        pgMaxConns,
		MinConns:        pgMinConns,
		MaxConnLifetime: pgMaxConnLifetime,
	}

	httpPort, err := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_PORT: %w", err)
	}

	readTimeout, err := time.ParseDuration(getEnv("HTTP_READ_TIMEOUT", "5s"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_READ_TIMEOUT: %w", err)
	}

	writeTimeout, err := time.ParseDuration(getEnv("HTTP_WRITE_TIMEOUT", "10s"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_WRITE_TIMEOUT: %w", err)
	}

	idleTimeout, err := time.ParseDuration(getEnv("HTTP_IDLE_TIMEOUT", "60s"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_IDLE_TIMEOUT: %w", err)
	}

	cfg.HTTP = HTTPConfig{
		Host:         getEnv("HTTP_HOST", "0.0.0.0"),
		Port:         httpPort,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	cfg.JWT = JWTConfig{
		SecretKey:     getEnv("JWT_SECRET", "secret"),
		TokenDuration: time.Minute * 5,
	}

	cfg.ChatApiKey = getEnv("CHAT_API_KEY", "")

	return &cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.DBName,
		c.Postgres.SSLMode,
	)
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.HTTP.Host, c.HTTP.Port)
}

func (c *Config) ExternalAddr() string {
	return fmt.Sprintf("%s:%d", c.App.Host, c.HTTP.Port)
}
