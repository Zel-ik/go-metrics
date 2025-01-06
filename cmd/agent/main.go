package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

type Metric struct {
	typpe string `json:"type"`
	name  string
	value float64
}

type MetricSlice [29]Metric
type rtm runtime.MemStats

const metricsLen = len(MetricSlice{})

func main() {
	client := resty.New()
	//apiUrl := "http://localhost:8080/update"
	counterCycle := 0

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	for {
		time.Sleep(time.Second * 2)
		runtime.ReadMemStats(&rtm)
		counterCycle += 1
		metrics := fillSlice(rtm, float64(counterCycle))
		if counterCycle%5 == 0 {
			metricForSend := metrics[rand.Intn(metricsLen-1)]
			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetPathParam("type", metricForSend.typpe).
				SetPathParam("name", metricForSend.name).
				SetPathParam("value", strconv.FormatFloat(metricForSend.value, 'f', -1, 64)).
				Post("http://localhost:8080/update/{type}/{name}/{value}")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(resp)
		}
	}
}

func fillSlice(rtm runtime.MemStats, pullCounter float64) MetricSlice {
	metricSlice := MetricSlice{
		Metric{typpe: "gauge", name: "Alloc", value: float64(rtm.Alloc)},
		Metric{typpe: "gauge", name: "BuckHashSys", value: float64(rtm.BuckHashSys)},
		Metric{typpe: "gauge", name: "Frees", value: float64(rtm.Frees)},
		Metric{typpe: "gauge", name: "GCCPUFraction", value: float64(rtm.GCCPUFraction)},
		Metric{typpe: "gauge", name: "GCSys", value: float64(rtm.GCSys)},
		Metric{typpe: "gauge", name: "HeapAlloc", value: float64(rtm.HeapAlloc)},
		Metric{typpe: "gauge", name: "HeapIdle", value: float64(rtm.HeapIdle)},
		Metric{typpe: "gauge", name: "HeapInuse", value: float64(rtm.HeapInuse)},
		Metric{typpe: "gauge", name: "HeapObjects", value: float64(rtm.HeapObjects)},
		Metric{typpe: "gauge", name: "HeapReleased", value: float64(rtm.HeapReleased)},
		Metric{typpe: "gauge", name: "HeapSys", value: float64(rtm.HeapSys)},
		Metric{typpe: "gauge", name: "LastGC", value: float64(rtm.LastGC)},
		Metric{typpe: "gauge", name: "Lookups", value: float64(rtm.Lookups)},
		Metric{typpe: "gauge", name: "MCacheInuse", value: float64(rtm.MCacheInuse)},
		Metric{typpe: "gauge", name: "MCacheSys", value: float64(rtm.MCacheSys)},
		Metric{typpe: "gauge", name: "MSpanInuse", value: float64(rtm.MSpanInuse)},
		Metric{typpe: "gauge", name: "MSpanSys", value: float64(rtm.MSpanSys)},
		Metric{typpe: "gauge", name: "Mallocs", value: float64(rtm.Mallocs)},
		Metric{typpe: "gauge", name: "NextGC", value: float64(rtm.NextGC)},
		Metric{typpe: "gauge", name: "NumForcedGC", value: float64(rtm.NumForcedGC)},
		Metric{typpe: "gauge", name: "NumGC", value: float64(rtm.NumGC)},
		Metric{typpe: "gauge", name: "OtherSys", value: float64(rtm.OtherSys)},
		Metric{typpe: "gauge", name: "PauseTotalNs", value: float64(rtm.PauseTotalNs)},
		Metric{typpe: "gauge", name: "StackInuse", value: float64(rtm.StackInuse)},
		Metric{typpe: "gauge", name: "StackSys", value: float64(rtm.StackSys)},
		Metric{typpe: "gauge", name: "Sys", value: float64(rtm.Sys)},
		Metric{typpe: "gauge", name: "TotalAlloc", value: float64(rtm.TotalAlloc)},
		Metric{typpe: "counter", name: "PollCount", value: pullCounter},
		Metric{typpe: "gauge", name: "RandomValue ", value: float64(rand.Intn(metricsLen))},
	}
	return metricSlice
}
