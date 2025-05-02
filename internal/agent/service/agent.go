package service

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dranikpg/dto-mapper"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"metric-collector/internal/agent/config"
	"metric-collector/internal/agent/metric"
	"metric-collector/internal/retry"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
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
			log.Info("UpdateMetric MemStats")
		}
	}()
	updatesURL := url.URL{
		Scheme: "http",
		Host:   config.GetConfig().FlagRunAddr,
		Path:   "/updates",
	}
	updateURL := url.URL{
		Scheme: "http",
		Host:   config.GetConfig().FlagRunAddr,
		Path:   "/update",
	}
	go func() {
		defer wg.Done()
		client := &http.Client{}
		for range reporter.C {
			var toUpload []metric.MetricsToUpload
			v := reflect.ValueOf(a.Stats)
			t := reflect.TypeOf(a.Stats)

			for i := 0; i < v.NumField(); i++ {
				field := t.Field(i)
				value := v.Field(i)
				metricToUpload := metric.MetricsToUpload{
					ID: field.Name,
				}

				switch field.Type.Kind() {
				case reflect.Int64, reflect.Int32:
					metricToUpload.MType = "counter"
					metricToUpload.Delta = new(int64)
					i2 := value.Int()
					metricToUpload.Delta = &i2

				case reflect.Float64:
					metricToUpload.MType = "gauge"
					metricToUpload.Value = new(float64)
					*metricToUpload.Value = value.Float()
				default:
					fmt.Printf("%s имеет неизвестный тип: %s\n", field.Name, field.Type)
				}
				err := sendHTTPRequest(updateURL.String(), metricToUpload, client)
				if err != nil {
					log.Error(err)
					continue
				}

				toUpload = append(toUpload, metricToUpload)

			}
			if len(toUpload) > 0 {
				err := sendHTTPRequest(updatesURL.String(), toUpload, client)
				if err != nil {
					log.Error(err)
					continue
				}

			}
		}

	}()

	wg.Wait()
}

func (a *Agent) updateMemStats() {
	var runtimeStats runtime.MemStats
	runtime.ReadMemStats(&runtimeStats)
	err := dto.Map(&a.Stats, runtimeStats)
	if err != nil {
		log.Fatal(err)
	}
	a.Stats.PollCount = a.Stats.PollCount + 1
	a.Stats.RandomValue = rand.Float64()
}

func sendHTTPRequest(baseURL string, metricToUpload interface{}, client *http.Client) error {
	jsonData, err := json.Marshal(metricToUpload)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err = gz.Write(jsonData)
	if err != nil {
		log.Error(err)
		return err
	}
	gz.Close()

	req, err := http.NewRequest("POST", baseURL, &buf)
	if err != nil {
		log.Error(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	addHashedHeader(req)
	resp, err := retry.Retry(3, 1, func() (*http.Response, error) {
		return client.Do(req)
	})
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	var responseBody []byte
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		gr, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Error(err)
			return err
		}
		defer gr.Close()
		responseBody, err = io.ReadAll(gr)
		if err != nil {
			log.Error(err)
		}
	} else {
		responseBody, err = io.ReadAll(resp.Body)
	}

	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("Response Status: ", resp.Status, " Response Body: ", string(responseBody))
	log.Info(baseURL)
	return nil
}

func addHashedHeader(req *http.Request) {
	if config.GetConfig().Key != "" {
		jsonData, err := io.ReadAll(req.Body)
		if err != nil {
			log.Error(err)
			return
		}
		// Restore the Body so it can be read again by the HTTP client
		req.Body = io.NopCloser(bytes.NewReader(jsonData))
		h := hmac.New(sha256.New, []byte(config.GetConfig().Key))
		h.Write(jsonData)
		result := h.Sum(nil)
		hashStr := hex.EncodeToString(result[:])
		req.Header.Set("HashSHA256", hashStr)
	}
}
