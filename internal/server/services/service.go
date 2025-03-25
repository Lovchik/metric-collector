package services

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"metric-collector/internal/server/metric"
	"metric-collector/internal/server/storage"
	"net/http"
)

type Service struct {
	WebServer *gin.Engine
	Store     *storage.MemStorage
}

func (s *Service) UpdateMetric(ctx *gin.Context) {
	var metrics metric.Metrics

	if err := ctx.ShouldBindJSON(&metrics); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	validateMetricsToUpdate(ctx, metrics)
	newMetric, err := metric.NewMetric(metrics)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	newMetric, err = s.Store.Update(newMetric)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return

	}
	ctx.JSON(http.StatusOK, metrics)
}

func validateMetricsToUpdate(ctx *gin.Context, metrics metric.Metrics) {
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
		counterValue := value.(int64)
		metrics.Delta = &counterValue

	}
	if metrics.MType == "gauge" {
		metrics.Value = new(float64)
		log.Info("metric gauge value: ", value)

		gaugeValue := value.(float64)
		metrics.Value = &gaugeValue

	}
	ctx.JSON(http.StatusOK, metrics)

}

func (s *Service) GetAllMetrics(context *gin.Context) {
	all := s.Store.GetAll()
	context.JSON(http.StatusOK, all)
}
