package service

import (
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
	Stats      metric.Metric
	StatsMutex sync.Mutex
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
		a.StatsMutex.Lock()
		defer a.StatsMutex.Unlock()
		for range reporter.C {
			v := reflect.ValueOf(a.Stats)
			t := reflect.TypeOf(a.Stats)

			for i := 0; i < v.NumField(); i++ {
				field := t.Field(i)
				value := v.Field(i)

				switch field.Type.Kind() {
				case reflect.Int64, reflect.Int32:
					sendHTTPRequest("http://"+config.GetConfig().FlagRunAddr+"/update/", field.Name, "counter", value.Int(), client)
				case reflect.Float64:
					sendHTTPRequest("http://"+config.GetConfig().FlagRunAddr+"/update/", "gauge", field.Name, value.Float(), client)
				default:
					fmt.Printf("%s имеет неизвестный тип: %s\n", field.Name, field.Type)
				}

			}
		}

	}()

	wg.Wait()
}

func (a *Agent) updateMemStats() {
	var runtimeStats runtime.MemStats
	runtime.ReadMemStats(&runtimeStats)
	a.StatsMutex.Lock()
	defer a.StatsMutex.Unlock()
	err := dto.Map(&a.Stats, runtimeStats)
	if err != nil {
		log.Fatal(err)
	}
	a.Stats.PollCount = a.Stats.PollCount + 1
}

func sendHTTPRequest(baseURL, nameValue string, typeValue string, value interface{}, client *http.Client) {
	var stringValue string
	switch v := value.(type) {
	case float64:
		stringValue = fmt.Sprintf("%.2f", v)
	case int64:
		stringValue = fmt.Sprintf("%d", v)
	}

	url := baseURL + typeValue + "/" + nameValue + "/" + stringValue
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error(err)
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Response Status: ", resp.Status, " Response Body: ", responseBody)
	log.Info(url)
}
