package main

import (
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
