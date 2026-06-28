package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

type AppConfig struct {
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	MaxOpenConns int
	MaxIdleConns int
	ConnTimeout  time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// DSN возвращает строку подключения к PostgreSQL.
// Метод на структуре — не глобальная функция. Это важно для тестируемости.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=%d",
		d.Host, d.Port, d.User, d.Password, d.Name,
		int(d.ConnTimeout.Seconds()),
	)
}

// Load читает конфиг из переменных окружения.
// Паника при отсутствии обязательных значений — это осознанное решение:
// приложение не должно стартовать с неполным конфигом.
func Load() (*Config, error) {
	dbMaxOpen, err := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_OPEN_CONNS: %w", err)
	}

	dbMaxIdle, err := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE_CONNS: %w", err)
	}

	connTimeout, err := time.ParseDuration(getEnv("DB_CONN_TIMEOUT", "5s"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_CONN_TIMEOUT: %w", err)
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:         mustGetEnv("DB_HOST"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         mustGetEnv("DB_USER"),
			Password:     mustGetEnv("DB_PASSWORD"),
			Name:         mustGetEnv("DB_NAME"),
			MaxOpenConns: dbMaxOpen,
			MaxIdleConns: dbMaxIdle,
			ConnTimeout:  connTimeout,
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return v
}
