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
	flagRunAddr := flag.String("a", "8080", "Server address")
	reportInterval := flag.Int64("r", 10, "Report interval")
	pollInterval := flag.Int64("p", 2, "Poll interval")

	flag.Parse()
	config := &Config{
		FlagRunAddr:    getEnv("ADDRESS", *flagRunAddr),
		ReportInterval: getEnvInt("REPORT_INTERVAL", *reportInterval),
		PollInterval:   getEnvInt("POLL_INTERVAL", *pollInterval),
	}
	appConfig = *config

}

func getEnv(envName, defaultVal string) string {
	if value := os.Getenv(envName); value != "" {
		return value
	}
	return defaultVal
}

func getEnvInt(envName string, defaultVal int64) int64 {
	if value := os.Getenv(envName); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return defaultVal
}
