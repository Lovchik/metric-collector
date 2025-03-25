package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dranikpg/dto-mapper"
	log "github.com/sirupsen/logrus"
	"io"
	"metric-collector/internal/agent/config"
	"metric-collector/internal/agent/metric"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type Agent struct {
	Stats   metric.Metric
	StatsMu sync.Mutex
}

func (a *Agent) Start() {
	poller := time.NewTicker(time.Duration(config.GetConfig().PollInterval) * time.Second)
	reporter := time.NewTicker(time.Duration(config.GetConfig().ReportInterval) * time.Second)
	defer poller.Stop()
	defer reporter.Stop()
	a.Stats.PollCount = 0
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for range poller.C {
			a.updateMemStats()
			log.Info("Update MemStats")
		}
	}()
	go func() {
		defer wg.Done()
		client := &http.Client{}
		a.StatsMu.Lock()
		defer a.StatsMu.Unlock()
		for range reporter.C {
			v := reflect.ValueOf(a.Stats)
			t := reflect.TypeOf(a.Stats)

			for i := 0; i < v.NumField(); i++ {
				field := t.Field(i)
				metricToUpload := metric.MetricsToUpload{
					ID: field.Name,
				}
				value := v.Field(i)

				switch field.Type.Kind() {
				case reflect.Int64, reflect.Int32:
					metricToUpload.MType = "counter"
					metricToUpload.Delta = new(int64)
					*metricToUpload.Delta = value.Int()

				case reflect.Float64:
					metricToUpload.MType = "gauge"
					metricToUpload.Value = new(float64)
					*metricToUpload.Value = value.Float()
				default:
					fmt.Printf("%s имеет неизвестный тип: %s\n", field.Name, field.Type)
				}
				sendHTTPRequest("http://"+config.GetConfig().FlagRunAddr+"/update", metricToUpload, client)

			}
		}

	}()

	wg.Wait()
}

func (a *Agent) updateMemStats() {
	var runtimeStats runtime.MemStats
	runtime.ReadMemStats(&runtimeStats)
	a.StatsMu.Lock()
	defer a.StatsMu.Unlock()
	err := dto.Map(&a.Stats, runtimeStats)
	if err != nil {
		log.Fatal(err)
	}
	a.Stats.PollCount = a.Stats.PollCount + 1
}

func sendHTTPRequest(baseURL string, metricToUpload metric.MetricsToUpload, client *http.Client) {

	jsonData, err := json.Marshal(metricToUpload)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Response Status: ", resp.Status, " Response Body: ", responseBody)
	log.Info(baseURL)
}
