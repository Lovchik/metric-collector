package main

import (
	"metric-collector/internal/agent/config"
	"metric-collector/internal/agent/metric"
	"metric-collector/internal/agent/service"
	"sync"
)

func main() {
	config.InitConfig()
	agent := service.Agent{
		Stats:      metric.Metric{},
		StatsMutex: sync.Mutex{},
	}
	agent.Start()
}
