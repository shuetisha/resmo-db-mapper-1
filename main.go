package main

import (
	"context"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"log"
	"resmo-db-mapper/pkg"
	"resmo-db-mapper/pkg/config"
	"strings"
	"time"
)

var (
	version string
)

func main() {
	if err := realMain(); err != nil {
		log.Fatalln("failed to run resmo-resmo-db-mapper:", err)
	}
}

func realMain() error {
	var config config.Config
	err := config.ReadConfig(version)
	fmt.Println(config)
	if err != nil {
		return fmt.Errorf("error while reading config : %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("config read successfully, will start to run queries")

	dbType, err := getDatabaseType(config.DSN)
	if err != nil {
		return fmt.Errorf("error while getting database type from DSN: %w", err)
	}
	if config.Schedule == "" {
		runQueries(ctx, config, dbType)
		log.Println("database resources ingested successfully")
		return nil
	}
	dur, err := time.ParseDuration(config.Schedule)
	if err != nil {
		return fmt.Errorf("could not parse schedule from config: %w", err)
	}

	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	ctxDur, err := time.ParseDuration(config.ContextTime)
	if err != nil {
		return fmt.Errorf("could not parse context duration from config: %w", err)
	}
	for range ticker.C {
		log.Println("running queries again with schedule: ", config.Schedule)
		ctx, cancel := context.WithTimeout(context.Background(), ctxDur)
		runQueries(ctx, config, dbType)
		cancel()
	}

	log.Println("database resources ingested successfully")
	return nil
}

//func readConfig(config *config.Config) error {
//	viper.SetConfigName("config")
//	viper.AddConfigPath(".")
//
//	if err := viper.ReadInConfig(); err != nil {
//		log.Println("failed to read config:", err)
//		log.Println("will try to read config from command line arguments")
//	}
//
//	flag.StringVar(&config.Schedule, "schedule", "", "schedule for running queries")
//	flag.StringVar(&config.DSN, "datasourceName", "", "database datasource name")
//	flag.StringVar(&config.IngestKey, "ingestKey", "", "ingestKey of the integration")
//	flag.StringVar(&config.URL, "url", "", "url for data to send to")
//	flag.Parse()
//
//	config.Version = version
//
//	if err := viper.Unmarshal(&config); err != nil {
//		return fmt.Errorf("failed to unmarshal config: %w", err)
//	}
//	if err := config.Validate(); err != nil {
//		return fmt.Errorf("config is not valid: %w", err)
//
//	}
//	return nil
//}

func runQueries(ctx context.Context, config config.Config, dbType string) {
	switch dbType {
	case "mongo":
		pkg.RunMongoQueries(ctx, config, dbType)
	default:
		pkg.RunSQLDatabaseQueries(ctx, config, dbType)
	}
}

func getDatabaseType(connectionString string) (string, error) {
	var dbType string

	if strings.HasPrefix(connectionString, "postgres://") || strings.HasPrefix(connectionString, "postgresql://") {
		dbType = "postgresql"
	} else if strings.HasPrefix(connectionString, "mongodb://") {
		dbType = "mongo"
	} else if strings.HasPrefix(connectionString, "clickhouse://") {
		dbType = "clickhouse"
	} else if strings.HasPrefix(connectionString, "mysql://") {
		dbType = "mysql"
	} else {
		return "", fmt.Errorf("unsupported database type for connection string: %s", connectionString)
	}

	return dbType, nil
}
