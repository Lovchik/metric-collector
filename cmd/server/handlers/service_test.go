package handlers

import (
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
		url      string
		expected int
	}{
		{"Valid gauge metric", "/update/gauge/cpu/99.9", http.StatusOK},
		{"Valid counter metric", "/update/counter/requests/10", http.StatusOK},
		{"Invalid type", "/update/unknown/cpu/99.9", http.StatusBadRequest},
		{"Missing name", "/update/gauge//99.9", http.StatusNotFound},
		{"Invalid gauge value", "/update/gauge/cpu/abc", http.StatusBadRequest},
		{"Invalid counter value", "/update/counter/requests/1.5", http.StatusBadRequest},
	}

	router := gin.Default()
	router.GET("/update/:type/:name/:value", validateMetricsToUpdate)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.url, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expected, w.Code)
		})
	}
}
