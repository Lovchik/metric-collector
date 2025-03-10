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
	var config Config
	flag.StringVar(&config.FlagRunAddr, "a", ":8080", "Server address")
	flag.Int64Var(&config.ReportInterval, "r", 10, "Report interval")
	flag.Int64Var(&config.PollInterval, "p", 2, "Poll interval")

	flag.Parse()

	getEnv("ADDRESS", &config.FlagRunAddr)
	getEnvInt("REPORT_INTERVAL", &config.ReportInterval)
	getEnvInt("POLL_INTERVAL", &config.PollInterval)

	appConfig = config

}

func getEnv(envName string, config *string) {
	if value := os.Getenv(envName); value != "" {
		config = &value
	}
	return
}

func getEnvInt(envName string, config *int64) {
	if value := os.Getenv(envName); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			config = &parsed
		}
	}
	return
}
