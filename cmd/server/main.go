package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"metric-collector/internal/server/config"
	"metric-collector/internal/server/handlers"
	"metric-collector/internal/server/services"
	"metric-collector/internal/server/storage"
)

func main() {
	config.InitConfig()
	Serve()
}

func Serve() {
	s := &services.Service{}
	s.WebServer = gin.Default()
	s.Store = storage.NewMemStorage()

	if config.GetConfig().DatabaseDNS == "" {
		s.Store = storage.NewMemStorage()
	} else {
		ctx := context.WithoutCancel(context.Background())

		pgStorage, err := storage.NewPgStorage(ctx, config.GetConfig().DatabaseDNS)
		if err != nil {
			return
		}
		s.Store = pgStorage
	}

	if config.GetConfig().Restore {
		err := s.Store.LoadMetricsInMemory(config.GetConfig().FileStoragePath)
		if err != nil {
			log.Error("Error loading metrics: ", err)
		}
	}
	go s.SaveMetricsToMemory()
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
