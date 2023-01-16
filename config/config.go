package config

import (
	"github.com/kelseyhightower/envconfig"
)

type (
	// Config provides the system configuration.
	Config struct {
		ServerAddress string `envconfig:"SERVER_ADDRESS"`
		NumWorkers    int    `envconfig:"MANAGER_NUM_WORKERS" default:"1"`
		Database      Database
		Debug         bool
	}

	// Database provides the database configuration.
	Database struct {
		Datasource     string `envconfig:"DATABASE_DATASOURCE"      default:"root:1@tcp(localhost:3306)/test?parseTime=true"`
		MaxConnections int    `envconfig:"DATABASE_MAX_CONNECTIONS" default:"0"`
	}
)

// Environ returns the settings from the environment.
func Environ() (Config, error) {
	cfg := Config{}
	defaultAddress(&cfg)
	err := envconfig.Process("", &cfg)
	return cfg, err
}

func defaultAddress(c *Config) {
	if c.ServerAddress == "" {
		c.ServerAddress = ":8080"
	}
}
