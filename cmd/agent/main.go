package main

import (
	"fmt"
	"github.com/dranikpg/dto-mapper"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"metric-collector/cmd/agent/metric"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

var Stats metric.Metric

func main() {
	ticker := time.NewTicker(2 * time.Second)
	ticker2 := time.NewTicker(10 * time.Second)
	Stats.PollCount = 0
	go func() {
		for range ticker.C {
			UpdateMemStats()
			log.Info("Update MemStats")
		}
	}()
	go func() {
		for range ticker2.C {

			v := reflect.ValueOf(Stats)
			t := reflect.TypeOf(Stats)

			for i := 0; i < v.NumField(); i++ {
				field := t.Field(i)
				value := v.Field(i)

				switch field.Type.Kind() {
				case reflect.Int64, reflect.Int32:
					sendHTTPRequest("http://127.0.0.1:8080/update/", field.Name, "counter", value.Int())
				case reflect.Float64:
					sendHTTPRequest("http://127.0.0.1:8080/update/", "gauge", field.Name, value.Float())
				default:
					fmt.Printf("%s имеет неизвестный тип: %s\n", field.Name, field.Type)
				}

			}
		}
	}()
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
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(resp.Body)

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Response Status: ", resp.Status, " Response Body: ", responseBody)
	}

}
