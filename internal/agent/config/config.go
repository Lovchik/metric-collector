package config

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

var appConfig Config

type Config struct {
	FlagRunAddr    string
	Key            string
	ReportInterval int64
	PollInterval   int64
}

func GetConfig() Config {
	return appConfig
}

func InitConfig() {
	config := Config{}

	getEnv("ADDRESS", "a", "localhost:8080", "Server address", &config.FlagRunAddr)
	getEnv("KEY", "k", "123123", "KEY", &config.Key)
	getEnvInt("REPORT_INTERVAL", "r", 3, "Report interval", &config.ReportInterval)
	getEnvInt("POLL_INTERVAL", "p", 1, "Poll interval", &config.PollInterval)
	flag.Parse()
	appConfig = config
	log.Info("Agent config : ", config)

}

func getEnv(envName, flagName, defaultValue, usage string, config *string) {
	flag.StringVar(config, flagName, defaultValue, usage)

	if value := os.Getenv(envName); value != "" {
		log.Info("Using environment variable "+envName, "value "+value)
		*config = value
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
