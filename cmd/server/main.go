package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"metric-collector/internal/server/config"
	"metric-collector/internal/server/handlers"
	"metric-collector/internal/server/services"
	"metric-collector/internal/server/storage"
	"time"
)

func main() {
	config.InitConfig()
	Serve()
}

func Serve() {
	s := &services.Service{}
	s.WebServer = gin.Default()
	s.Store = storage.NewMemStorage()
	if config.GetConfig().Restore {
		err := s.Store.LoadMetricsInMemory(config.GetConfig().FileStoragePath)
		if err != nil {
			log.Error("Error loading metrics: ", err)

		}
	}
	go func() {
		for {
			if config.GetConfig().StoreInterval > 0 {
				time.Sleep(time.Duration(config.GetConfig().StoreInterval) * time.Second)
			}
			err := s.Store.SaveMemoryInfo(config.GetConfig().FileStoragePath)
			if err != nil {
				log.Error(err)
			}
		}
	}()
	ginConfig := cors.DefaultConfig()
	ginConfig.AllowAllOrigins = true
	s.WebServer.Use(cors.New(ginConfig))
	api := s.WebServer.Group("/")
	handlers.MetricRouter(api.Group(""), s)
	err := s.WebServer.Run(config.GetConfig().FlagRunAddr)
	if err != nil {
		log.Error(err)
		return
	}
}
