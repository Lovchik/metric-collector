package main

import (
	"metric-collector/internal/agent/config"
	"metric-collector/internal/agent/service"
)

func main() {
	config.InitConfig()
	service.Start()
}
