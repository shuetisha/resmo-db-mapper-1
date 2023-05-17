package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	URL         string
	IngestKey   string
	Version     string
	DSN         string
	Schedule    string
	ContextTime string
}

func (c *Config) Validate() error {
	if c.IngestKey != "" {
		return errors.New("ingest key is not set")
	}

	if c.DSN == "" {
		return errors.New("datasource name is not set")
	}

	return nil
}

func (c *Config) ReadConfig(ver string) error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("failed to read config:", err)
		log.Println("will try to read config from command line arguments")
	}

	// Check if config variables are set (for command line runner)
	flag.StringVar(&c.Schedule, "schedule", "", "schedule for running queries")
	flag.StringVar(&c.ContextTime, "contextTime", "", "timeout duration for connections")
	flag.StringVar(&c.DSN, "datasourceName", "", "database datasource name")
	flag.StringVar(&c.IngestKey, "ingestKey", "", "ingestKey of the integration")
	flag.Parse()

	c.Version = ver

	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := c.Validate(); err != nil {
		return fmt.Errorf("config is not valid: %w", err)
	}

	return nil
}
