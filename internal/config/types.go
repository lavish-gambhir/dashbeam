package config

import "time"

type AppConfig struct {
	Environment string          `mapstructure:"environment"` // e.g., "development", "production", "local"
	Server      ServerConfig    `mapstructure:"server"`
	Database    DBConfig        `mapstructure:"database"`
	Auth        AuthConfig      `mapstructure:"auth"`
	Redis       RedisConfig     `mapstructure:"redis"`
	Quiz        QuizConfig      `mapstructure:"quiz"`
	Analytics   AnalyticsConfig `mapstructure:"analytics"`
	Reporting   ReportingConfig `mapstructure:"reporting"`
}

type ServerConfig struct {
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
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"` // TODO: fetch from secretsmanager
	DB       int    `mapstructure:"db"`
}

type QuizConfig struct {
}

type AnalyticsConfig struct {
}

type ReportingConfig struct {
}
