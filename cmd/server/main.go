package main

import (
	"metric-collector/cmd/server/handlers"
	"metric-collector/cmd/server/storage"
	"net/http"
)

func main() {
	http.Handle("/update/{type}/{name}/{value}", handlers.ValidationMiddleware(http.HandlerFunc(handlers.MetricPage)))
	storage.NewMemStorage()
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}
