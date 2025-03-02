package main

import (
	"flag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"metric-collector/cmd/server/handlers"
	"metric-collector/cmd/server/storage"
)

func main() {
	parseFlags()
	Serve()
}

var flagRunAddr string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
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
	err := s.WebServer.Run(flagRunAddr)
	if err != nil {
		return
	}
}
