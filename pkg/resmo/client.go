package resmo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"resmo-db-mapper/pkg/config"
)

func Ingest(ctx context.Context, config config.Config, driverType string, resourceKey, results interface{}) error {
	url := "http://host.docker.internal:9090/integration/%s/ingest/%s"

	data, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("error marshaling %s results: %w", resourceKey, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf(url, driverType, resourceKey), bytes.NewBufferString(string(data)))

	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Ingest-Key", config.IngestKey)
	req.Header.Set("Resmo-Database-Agent", config.Version)
	req.Header.Set("DB-URL", config.DSN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}
