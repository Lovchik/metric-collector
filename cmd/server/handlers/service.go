package handlers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"metric-collector/cmd/server/metric"
	"metric-collector/cmd/server/storage"
	"net/http"
	"strconv"
)

type Service struct {
	WebServer *gin.Engine
}

func (s *Service) UpdateMetric(c *gin.Context) {
	validateMetricsToUpdate(c)
	newMetric := metric.NewMetric(
		c.Param("name"),
		c.Param("type"),
		c.Param("value"))
	err := newMetric.Update()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return

	}
	c.JSON(http.StatusOK, nil)
}

func validateMetricsToUpdate(c *gin.Context) {
	metricType := c.Param("type")

	if metricType == "" || (metricType != "counter" && metricType != "gauge") {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	value := c.Param("value")
	if c.Param("name") == "" || value == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return

	}
	if metricType == "gauge" {
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

	} else {
		_, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
}

func (s *Service) GetGauge(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	value, exists := storage.Store.GetValueByName(name)
	if !exists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, value)
}

func (s *Service) GetCounter(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	value, exists := storage.Store.GetValueByName(name)
	if !exists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	log.Info(value)
	c.JSON(http.StatusOK, value)

}

func (s *Service) GetAllMetrics(context *gin.Context) {
	all := storage.Store.GetAll()
	context.JSON(http.StatusOK, all)
}
