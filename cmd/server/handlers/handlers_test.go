package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_metricPage(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		url        string
		wantStatus int
	}{
		{
			name:       "Test gauge status ok",
			method:     "POST",
			url:        "/update/gauge/Asa/122",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Test counter status ok",
			method:     "POST",
			url:        "/update/counter/Asa/122",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_ = httptest.NewRequest(tt.method, tt.url, nil)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
