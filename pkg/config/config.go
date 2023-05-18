package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	URL            string
	IngestKey      string
	Version        string
	DSN            string
	Schedule       string
	ContextTime    string
	ConfigPath     string
	DomainOverride string
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
	flag.StringVar(&c.DomainOverride, "domainOverride", "", "domain url for ingesting")

	flag.Parse()

	if c.DSN != "" && c.IngestKey != "" {
		if err := c.Validate(); err != nil {
			return fmt.Errorf("config is not valid: %w", err)
		}

		return nil
	}

	c.DSN = os.Getenv("DSN")
	c.IngestKey = os.Getenv("INGEST_KEY")

	if err := c.Validate(); err != nil {
		return fmt.Errorf("config is not valid: %w", err)
	}

	return nil
}
