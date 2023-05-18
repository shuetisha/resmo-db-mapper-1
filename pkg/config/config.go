package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	URL         string
	IngestKey   string
	Version     string
	DSN         string
	Schedule    string
	ContextTime string
	ConfigPath  string
}

func (c *Config) Validate() error {
	if c.IngestKey == "" {
		return errors.New("ingest key is not set")
	}

	if c.DSN == "" {
		return errors.New("datasource name is not set")
	}

	return nil
}

func (c *Config) ReadConfig(ver string) error {
	c.Version = ver

	flag.StringVar(&c.Schedule, "schedule", "", "schedule for running queries")
	flag.StringVar(&c.ContextTime, "contextTime", "", "timeout duration for connections")
	flag.StringVar(&c.DSN, "datasourceName", "", "database datasource name")
	flag.StringVar(&c.IngestKey, "ingestKey", "", "ingestKey of the integration")
	flag.StringVar(&c.ConfigPath, "configPath", "", "absolute path of the config")

	flag.Parse()

	if c.DSN != "" && c.IngestKey != "" {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("config is not valid: %w", err)
		}

		return nil
	}

	viper.SetConfigName("resmo-db-mapper-config")
	viper.AddConfigPath(".")

	if c.ConfigPath != "" {
		viper.AddConfigPath(c.ConfigPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := c.Validate(); err != nil {
		return fmt.Errorf("config is not valid: %w", err)
	}

	return nil
}
