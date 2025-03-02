package main

import (
	"flag"
	"fmt"
	"github.com/dranikpg/dto-mapper"
	log "github.com/sirupsen/logrus"
	"io"
	"metric-collector/cmd/agent/metric"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

var reportInterval, pollInterval int
var flagRunAddr string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "report interval")
	flag.IntVar(&pollInterval, "p", 2, "poll interval")
	flag.Parse()
}

var Stats metric.Metric

func main() {
	parseFlags()
	poller := time.NewTicker(time.Duration(pollInterval) * time.Second)
	reporter := time.NewTicker(time.Duration(reportInterval) * time.Second)
	Stats.PollCount = 0
	go func() {
		for range poller.C {
			UpdateMemStats()
			log.Info("Update MemStats")
		}
	}()
	go func() {
		for range reporter.C {

			v := reflect.ValueOf(Stats)
			t := reflect.TypeOf(Stats)

			for i := 0; i < v.NumField(); i++ {
				field := t.Field(i)
				value := v.Field(i)

				switch field.Type.Kind() {
				case reflect.Int64, reflect.Int32:
					sendHTTPRequest("http://"+flagRunAddr+"/update/", field.Name, "counter", value.Int())
				case reflect.Float64:
					sendHTTPRequest("http://"+flagRunAddr+"/update/", "gauge", field.Name, value.Float())
				default:
					fmt.Printf("%s имеет неизвестный тип: %s\n", field.Name, field.Type)
				}

			}
		}
	}()
	select {}
}

func UpdateMemStats() {
	var runtimeStats runtime.MemStats
	runtime.ReadMemStats(&runtimeStats)
	err := dto.Map(&Stats, runtimeStats)
	if err != nil {
		log.Fatal(err)
	}
	Stats.PollCount = Stats.PollCount + 1

}

func sendHTTPRequest(baseURL, nameValue string, typeValue string, value interface{}) {
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
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
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
