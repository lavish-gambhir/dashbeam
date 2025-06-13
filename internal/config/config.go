package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	_DefaultConfigFileName = "local"
	// EnvPrefix is the prefix for environment variables, e.g., APP_SERVER_PORT.
	_EnvPrefix = "APP"
)

var appConfig *AppConfig

func LoadConfig() (*AppConfig, error) {
	v := viper.New()

	v.SetConfigName(_DefaultConfigFileName) // Name of config file (without extension)
	v.SetConfigType("yaml")                 // Type of config file
	v.AddConfigPath("./configs")            // Search in the 'configs' directory
	v.AddConfigPath(".")                    // Search in the current directory

	env := os.Getenv(_EnvPrefix + "_ENV")
	if env == "" {
		env = "development"
	}
	v.Set("environment", env)
	v.SetConfigFile(fmt.Sprintf("configs/%s.yaml", env))
	if _, err := os.Stat(fmt.Sprintf("configs/%s.yaml", env)); err == nil {
		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("failed to merge environment-specific config file '%s.yaml': %w", env, err)
		}
		log.Printf("Successfully merged environment-specific config: %s.yaml", env)
	}

	// Read the default config.yaml (or whatever is set by SetConfigName)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("No default config file found (e.g., local.yaml), relying on environment variables or specific config.")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		log.Printf("Successfully loaded config file: %s", v.ConfigFileUsed())
	}

	v.SetEnvPrefix(_EnvPrefix)                         // All environment variables should start with APP_
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Replaces `.` in config keys with `_` for env vars
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.Unmarshal(&appConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	log.Printf("Running in %s environment", appConfig.Environment)
	log.Println("Configuration loaded successfully!")
	return appConfig, nil
}

func GetConfig() (*AppConfig, error) {
	if appConfig == nil {
		return nil, fmt.Errorf("Configuration not loaded. Call config.LoadConfig() first.")
	}
	return appConfig, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", "8080")
}
