package services

import (
	"bytes"
	"encoding/json"
	"metric-collector/internal/server/metric"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateMetricsToUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		metrics  metric.Metrics
		expected int
	}{
		{
			name: "Valid gauge",
			metrics: metric.Metrics{
				ID:    "cpu_load",
				MType: "gauge",
				Value: floatPtr(1.23),
			},
			expected: http.StatusOK,
		},
		{
			name: "Valid counter",
			metrics: metric.Metrics{
				ID:    "requests_count",
				MType: "counter",
				Delta: int64Ptr(10),
			},
			expected: http.StatusOK,
		},
		{
			name: "Missing ID",
			metrics: metric.Metrics{
				MType: "gauge",
				Value: floatPtr(1.23),
			},
			expected: http.StatusNotFound,
		},
		{
			name: "Invalid type",
			metrics: metric.Metrics{
				ID:    "invalid",
				MType: "unknown",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Gauge with Delta",
			metrics: metric.Metrics{
				ID:    "cpu_load",
				MType: "gauge",
				Delta: int64Ptr(5),
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			r.POST("/update", func(ctx *gin.Context) {
				validateMetricsToUpdate(ctx, tt.metrics)
				ctx.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			body, _ := json.Marshal(tt.metrics)
			req, _ := http.NewRequest("POST", "/update", bytes.NewBuffer(body))
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expected, w.Code)
		})
	}
}

func floatPtr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
