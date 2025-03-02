package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"metric-collector/cmd/server/handlers"
	"metric-collector/cmd/server/storage"
)

func main() {
	Serve()
}

func Serve() {
	storage.NewMemStorage()
	s := &handlers.Service{}
	s.WebServer = gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	s.WebServer.Use(cors.New(config))
	api := s.WebServer.Group("/")
	handlers.MetricRouter(api.Group(""), s)
	err := s.WebServer.Run(":8080")
	if err != nil {
		return
	}
}
