package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	RateLimit RateLimitConfig
	Cache     CacheConfig
}

type AppConfig struct {
	Env  string `env:"APP_ENV" envDefault:"development"`
	Port string `env:"APP_PORT" envDefault:"8080"`
}

type DatabaseConfig struct {
	Host         string        `env:"DB_HOST,required"`
	Port         string        `env:"DB_PORT" envDefault:"5432"`
	User         string        `env:"DB_USER,required"`
	Password     string        `env:"DB_PASSWORD,required"`
	Name         string        `env:"DB_NAME,required"`
	MaxOpenConns int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	ConnTimeout  time.Duration `env:"DB_CONN_TIMEOUT" envDefault:"5s"`
}

type RedisConfig struct {
	Addr         string        `env:"REDIS_ADDR,required"`
	Password     string        `env:"REDIS_PASSWORD"`
	DB           int           `env:"REDIS_DB" envDefault:"0"`
	PoolSize     int           `env:"REDIS_POOL_SIZE" envDefault:"10"`
	MinIdleConns int           `env:"REDIS_MIN_IDLE_CONNS" envDefault:"2"`
	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT" envDefault:"3s"`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT" envDefault:"3s"`
}

type RateLimitConfig struct {
	DefaultRPS   int           `env:"RATE_LIMIT_DEFAULT_RPS" envDefault:"10"`
	AuthRPS      int           `env:"RATE_LIMIT_AUTH_RPS" envDefault:"3"`
	WindowLength time.Duration `env:"RATE_LIMIT_WINDOW" envDefault:"1s"`
}

type CacheConfig struct {
	SiteStatusTTL time.Duration `env:"CACHE_SITE_STATUS_TTL" envDefault:"30s"`
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
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}
