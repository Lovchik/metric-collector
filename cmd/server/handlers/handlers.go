package handlers

import (
	log "github.com/sirupsen/logrus"
	"metric-collector/cmd/server/metric"
	"net/http"
	"strconv"
)

func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.PathValue("name") == "" || r.PathValue("type") == "" || r.PathValue("value") == "" {
			w.WriteHeader(http.StatusNotFound)
		}
		metricType := r.PathValue("type")
		value := r.PathValue("value")
		if metricType == "gauge" {
			_, err := strconv.ParseFloat(value, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

		} else if metricType == "counter" {
			_, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func MetricPage(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)
	newMetric := metric.NewMetric(
		req.PathValue("name"),
		req.PathValue("type"),
		req.PathValue("value"))
	err := newMetric.Update()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)

}
