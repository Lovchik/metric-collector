package config

import (
	"flag"
	"os"
	"strconv"
)

var appConfig Config

type Config struct {
	FlagRunAddr     string
	StoreInterval   int64
	FileStoragePath string
	Restore         bool
	DatabaseDNS     string
}

func GetConfig() Config {
	return appConfig
}

func InitConfig() {
	config := Config{}

	getEnv("ADDRESS", "a", ":8080", "Server address", &config.FlagRunAddr)
	getEnvInt("STORE_INTERVAL", "i", 300, "Report interval", &config.StoreInterval)
	getEnv("FILE_STORAGE_PATH", "f", "file.txt", "file storage path ", &config.FileStoragePath)
	getEnv("DATABASE_DSN", "d", "localhost:5001", "file storage path ", &config.DatabaseDNS)
	getEnvBool("RESTORE", "p", false, "Poll interval", &config.Restore)
	flag.Parse()

	appConfig = config

}

func getEnv(envName, flagName, defaultValue, usage string, config *string) {
	if value := os.Getenv(envName); value != "" {
		*config = value
	} else {
		flag.StringVar(config, flagName, defaultValue, usage)
	}

}

func getEnvInt(envName string, flagName string, defaultValue int64, usage string, config *int64) {
	if value := os.Getenv(envName); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			*config = parsed
		}
	} else {
		flag.Int64Var(config, flagName, defaultValue, usage)
	}
}

func getEnvBool(envName string, flagName string, defaultValue bool, usage string, config *bool) {
	if value := os.Getenv(envName); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			*config = parsed
		}
	} else {
		flag.BoolVar(config, flagName, defaultValue, usage)
	}
}
