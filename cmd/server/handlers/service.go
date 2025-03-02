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

func (s *Service) UpdateCounter(c *gin.Context) {
	metricType := "counter"
	validateMetricsToUpdate(c, metricType)

	newMetric := metric.NewMetric(
		c.Param("name"),
		metricType,
		c.Param("value"))
	err := newMetric.Update()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return

	}
	c.JSON(http.StatusOK, nil)
}

func validateMetricsToUpdate(c *gin.Context, metricType string) {
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

	} else if metricType == "counter" {
		_, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

func (s *Service) UpdateGauge(c *gin.Context) {
	metricType := "gauge"
	validateMetricsToUpdate(c, metricType)

	newMetric := metric.NewMetric(
		c.Param("name"),
		metricType,
		c.Param("value"))
	err := newMetric.Update()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, nil)

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
