package config

import (
	"flag"
	log "github.com/sirupsen/logrus"
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
	Key             string
}

func GetConfig() Config {
	return appConfig
}

func InitConfig() {
	config := Config{}

	getEnv("KEY", "k", "123123", "KEY", &config.Key)
	getEnv("ADDRESS", "a", ":8080", "Server address", &config.FlagRunAddr)
	getEnvInt("STORE_INTERVAL", "i", 300, "Report interval", &config.StoreInterval)
	getEnv("FILE_STORAGE_PATH", "f", "file.json", "file storage path ", &config.FileStoragePath)
	getEnv("DATABASE_DSN", "d", "", "file storage path ", &config.DatabaseDNS)
	getEnvBool("RESTORE", "r", false, "Poll interval", &config.Restore)
	flag.Parse()
	log.Info("Config: ", config)
	appConfig = config

}

func getEnv(envName, flagName, defaultValue, usage string, config *string) {
	flag.StringVar(config, flagName, defaultValue, usage)

	if value := os.Getenv(envName); value != "" {
		log.Info("Using environment variable "+envName, "value "+value)
		*config = value
	}
}

func getEnvInt(envName string, flagName string, defaultValue int64, usage string, config *int64) {
	flag.Int64Var(config, flagName, defaultValue, usage)

	if value := os.Getenv(envName); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			*config = parsed
		}
	}
}

func getEnvBool(envName string, flagName string, defaultValue bool, usage string, config *bool) {
	flag.BoolVar(config, flagName, defaultValue, usage)

	if value := os.Getenv(envName); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			*config = parsed
		}
	}
}
