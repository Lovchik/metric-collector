package config

import (
	"flag"
	"log"
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

	config.FlagRunAddr = getEnv("ADDRESS")
	config.ReportInterval = getEnvInt("REPORT_INTERVAL")
	config.PollInterval = getEnvInt("POLL_INTERVAL")

	appConfig = config

}

func getEnv(envName string) string {
	if value := os.Getenv(envName); value != "" {
		return value
	}
	log.Fatalf("Environment variable %s not set", envName)
	return ""
}

func getEnvInt(envName string) int64 {
	if value := os.Getenv(envName); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	log.Fatalf("Environment variable %s not set", envName)
	return 0
}
