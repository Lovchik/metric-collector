package handlers

import (
	"github.com/gin-gonic/gin"
)

func MetricRouter(router *gin.RouterGroup, s *Service) {
	router.POST("update/gauge/:name/:value", s.UpdateGauge)
	router.POST("update/counter/:name/:value", s.UpdateCounter)
	router.GET("value/counter/:name", s.GetCounter)
	router.GET("value/gauge/:name", s.GetGauge)

}
