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
	config := &Config{
		FlagRunAddr:    getVariable("a", "ADDRESS", ":8080"),
		ReportInterval: getIntVariable("r", "REPORT_INTERVAL", 10),
		PollInterval:   getIntVariable("p", "POLL_INTERVAL", 2),
	}
	appConfig = *config

}

func getVariable(flagName, envName, defaultVal string) string {
	var value string
	flag.StringVar(&value, flagName, defaultVal, "")
	flag.Parse()

	if envRunAddr := os.Getenv(envName); envRunAddr != "" {
		value = envRunAddr
	}
	return value
}

func getIntVariable(flagName, envName string, defaultVal int64) int64 {
	var value int64
	flag.Int64Var(&value, flagName, defaultVal, "")
	flag.Parse()

	if env := os.Getenv(envName); env != "" {
		value, _ = strconv.ParseInt(env, 10, 64)
	}
	return value
}
