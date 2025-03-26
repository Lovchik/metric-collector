package services

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"metric-collector/internal/server/metric"
	"metric-collector/internal/server/storage"
	"net/http"
	"strconv"
	"strings"
)

type Service struct {
	WebServer *gin.Engine
	Store     *storage.MemStorage
}

func (s *Service) UpdateMetric(c *gin.Context) {
	validateMetricsToUpdate(c)
	newMetric, err := metric.NewMetric(
		c.Param("name"),
		c.Param("type"),
		c.Param("value"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	metr, err := s.Store.Update(newMetric)
	log.Info(metr)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return

	}
	c.JSON(http.StatusOK, nil)
}

func (s *Service) UpdateMetricViaJSON(ctx *gin.Context) {
	var metrics metric.Metrics

	if err := ctx.ShouldBindJSON(&metrics); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	validateMetricsToUpdateViaJSON(ctx, metrics)
	newMetric, err := metric.NewMetricFromJSON(metrics)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	newMetric, err = s.Store.Update(newMetric)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return

	}
	value := newMetric.GetValue()

	if metrics.MType == "gauge" {
		if v, ok := value.(float64); ok {
			metrics.Value = &v
		}
	} else {
		if v, ok := value.(int64); ok {
			metrics.Delta = &v
		}
	}
	responseData, err := json.Marshal(metrics)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	if strings.Contains(ctx.GetHeader("Accept-Encoding"), "gzip") {
		ctx.Header("Content-Encoding", "gzip")
		ctx.Header("Content-Type", "application/json")
		ctx.Header("Content-Type", "text/html")
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		_, err := gz.Write(responseData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
		gz.Close()
		ctx.Data(http.StatusOK, "application/json", buf.Bytes())
	} else {
		ctx.JSON(http.StatusOK, metrics)
	}
}

func validateMetricsToUpdateViaJSON(ctx *gin.Context, metrics metric.Metrics) {
	if validateType(ctx, metrics) {
		return
	}

	if metrics.ID == "" || (metrics.Value == nil && metrics.Delta == nil) {
		ctx.JSON(http.StatusNotFound, nil)
		return

	}
	if metrics.MType == "gauge" {

		if metrics.Value == nil {
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}

	} else {
		if metrics.Delta == nil {
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}
	}
}

func validateType(ctx *gin.Context, metrics metric.Metrics) bool {
	if metrics.MType == "" || (metrics.MType != "counter" && metrics.MType != "gauge") {
		ctx.JSON(http.StatusBadRequest, nil)
		return true
	}
	return false
}

func (s *Service) GetMetric(ctx *gin.Context) {
	var metrics metric.Metrics

	if err := ctx.ShouldBindJSON(&metrics); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	if validateType(ctx, metrics) {
		return
	}

	if metrics.ID == "" {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}
	value, exists := s.Store.GetValueByName(metrics.ID)
	if !exists {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}
	if metrics.MType == "counter" {
		metrics.Delta = new(int64)
		log.Info("metric counter value: ", value)
		counterValue, ok := value.(int64)
		if !ok {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
		metrics.Delta = &counterValue

	}
	if metrics.MType == "gauge" {
		metrics.Value = new(float64)
		log.Info("metric gauge value: ", value)

		gaugeValue, ok := value.(float64)
		if !ok {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
		metrics.Value = &gaugeValue

	}
	ctx.Header("Content-Type", "text/html")
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, metrics)

}
func (s *Service) GetAllMetrics(context *gin.Context) {
	all := s.Store

	acceptEncoding := context.GetHeader("Accept-Encoding")
	accept := context.GetHeader("Accept")

	if accept == "text/html" {
		context.Header("Content-Type", "text/html")
	}

	if acceptEncoding == "gzip" {
		context.Header("Content-Encoding", "gzip")
		writer := gzip.NewWriter(context.Writer)
		defer writer.Close()
		context.Writer = &gzipResponseWriter{Writer: writer, ResponseWriter: context.Writer}
	}

	context.JSON(http.StatusOK, all)
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.Writer.Write(data)
}

func (s *Service) GetGauge(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	value, exists := s.Store.GetValueByName(name)
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
	value, exists := s.Store.GetValueByName(name)
	if !exists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	log.Info(value)
	c.JSON(http.StatusOK, value)

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
