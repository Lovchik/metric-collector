package config

import (
	"flag"
	"os"
	"strconv"
)

var appConfig Config

type Config struct {
	FlagRunAddr    string
	ReportInterval int64
	PollInterval   int64
}

func GetConfig() Config {
	return appConfig
}

func InitConfig() {
	config := Config{}

	getEnv("ADDRESS", "a", "localhost:8080", "Server address", &config.FlagRunAddr)
	getEnvInt("REPORT_INTERVAL", "r", 3, "Report interval", &config.ReportInterval)
	getEnvInt("POLL_INTERVAL", "p", 1, "Poll interval", &config.PollInterval)
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
