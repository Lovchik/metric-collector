package config

import (
	"flag"
	"os"
)

var appConfig Config

type Config struct {
	FlagRunAddr string
}

func GetConfig() Config {
	return appConfig
}

func InitConfig() {
	config := &Config{
		FlagRunAddr: getVariable("a", "ADDRESS", ":8080"),
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
