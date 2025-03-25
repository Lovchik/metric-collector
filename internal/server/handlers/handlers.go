package handlers

import (
	"github.com/gin-gonic/gin"
	"metric-collector/internal/server/middleware"
	"metric-collector/internal/server/services"
)

func MetricRouter(router *gin.RouterGroup, s *services.Service) {
	router.POST("/update", middleware.LoggerMiddleware(), s.UpdateMetric)
	router.POST("/value", middleware.LoggerMiddleware(), s.GetMetric)
	router.GET("", middleware.LoggerMiddleware(), s.GetAllMetrics)
}
