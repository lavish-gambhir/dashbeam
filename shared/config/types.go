package config

import (
	"fmt"
	"net/url"
	"time"
)

type AppConfig struct {
	Env       Environment     `mapstructure:"environment"` // e.g., "development", "production", "local"
	Server    ServerConfig    `mapstructure:"server"`
	Database  DBConfig        `mapstructure:"database"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Quiz      QuizConfig      `mapstructure:"quiz"`
	Analytics AnalyticsConfig `mapstructure:"analytics"`
	Reporting ReportingConfig `mapstructure:"reporting"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DBConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"` // TODO: fetch from secretsmanager
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type AuthConfig struct {
	JWTSecretKey       string        `mapstructure:"jwt_secret_key"` // TODO: fetch from secrets manager
	AccessTokenExpiry  time.Duration `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
}

type RedisConfig struct {
	DB         int           `mapstructure:"db"`
	PoolSize   int           `mapstructure:"pool_size"`
	MaxRetries int           `mapstructure:"max_retries"`
	Addr       string        `mapstructure:"addr"`
	Password   string        `mapstructure:"password"` // TODO: fetch from secretsmanager
	Timeout    time.Duration `mapstructure:"timeout"`
}

type QuizConfig struct {
}

type AnalyticsConfig struct {
}

type ReportingConfig struct {
}

func (d *DBConfig) Address() string {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(d.User, d.Password),
		Host:   fmt.Sprintf("%s:%d", d.Host, d.Port),
		Path:   d.DBName,
	}
	q := dsn.Query()
	q.Add("sslmode", d.SSLMode)
	dsn.RawQuery = q.Encode()
	return dsn.String()
}
