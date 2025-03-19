package handlers

import (
	"github.com/gin-gonic/gin"
	"metric-collector/internal/server/middleware"
	"metric-collector/internal/server/services"
)

func MetricRouter(router *gin.RouterGroup, s *services.Service) {
	router.POST("update/:type/:name/:value", middleware.LoggerMiddleware(), s.UpdateMetric)
	router.GET("value/counter/:name", middleware.LoggerMiddleware(), s.GetCounter)
	router.GET("value/gauge/:name", middleware.LoggerMiddleware(), s.GetGauge)
	router.GET("", middleware.LoggerMiddleware(), s.GetAllMetrics)
}
