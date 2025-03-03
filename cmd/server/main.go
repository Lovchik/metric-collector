package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"metric-collector/cmd/server/config"
	"metric-collector/cmd/server/handlers"
	"metric-collector/cmd/server/storage"
)

func main() {
	config.InitConfig()
	Serve()
}

func Serve() {

	storage.NewMemStorage()
	s := &handlers.Service{}
	s.WebServer = gin.Default()
	ginConfig := cors.DefaultConfig()
	ginConfig.AllowAllOrigins = true
	s.WebServer.Use(cors.New(ginConfig))
	api := s.WebServer.Group("/")
	handlers.MetricRouter(api.Group(""), s)
	err := s.WebServer.Run(config.GetConfig().FlagRunAddr)
	if err != nil {
		return
	}
}
