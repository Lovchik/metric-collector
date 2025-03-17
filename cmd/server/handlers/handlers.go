package handlers

import (
	"github.com/gin-gonic/gin"
)

func MetricRouter(router *gin.RouterGroup, s *Service) {
	router.POST("update/:type/:name/:value", s.UpdateMetric)
	router.GET("value/counter/:name", s.GetCounter)
	router.GET("value/gauge/:name", s.GetGauge)
	router.GET("", s.GetAllMetrics)
}
