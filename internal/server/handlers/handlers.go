package handlers

import (
	"github.com/gin-gonic/gin"
	"metric-collector/internal/server/middleware"
	"metric-collector/internal/server/services"
)

func MetricRouter(router *gin.RouterGroup, s *services.Service) {
	router.POST("/update", middleware.AuthMiddleware(), middleware.LoggerMiddleware(), middleware.GzipMiddleware(), s.UpdateMetricViaJSON)
	router.POST("/updates", middleware.AuthMiddleware(), middleware.LoggerMiddleware(), middleware.GzipMiddleware(), s.UpdateMetrics)
	router.POST("/value", middleware.AuthMiddleware(), middleware.LoggerMiddleware(), middleware.GzipMiddleware(), s.GetMetric)
	router.GET("", middleware.LoggerMiddleware(), s.GetAllMetrics)
	router.POST("update/:type/:name/:value", middleware.AuthMiddleware(), middleware.LoggerMiddleware(), s.UpdateMetric)
	router.GET("value/counter/:name", middleware.LoggerMiddleware(), s.GetCounter)
	router.GET("value/gauge/:name", middleware.LoggerMiddleware(), s.GetGauge)
	router.GET("/ping", s.HealthCheck)
}
