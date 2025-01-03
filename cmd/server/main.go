package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	metrics map[string]string
}

func NewMemoryStore() *MemStorage {
	m := MemStorage{
		make(map[string]string),
	}
	return &m
}

func (m *MemStorage) add(key, value string) {
	if m.metrics[key] == "" {
		m.metrics[key] = value
		return
	}
	val, err := strconv.Atoi(m.metrics[key])
	if err != nil {
		fmt.Print(err)
	}
	valToAdd, err := strconv.Atoi(value)
	if err != nil {
		fmt.Print(err)
	}
	val += valToAdd
	m.metrics[key] = strconv.Itoa(val)
}

func (m *MemStorage) change(key, value string) {
	m.metrics[key] = value
}

func (m *MemStorage) getList() map[string]string {
	return m.metrics
}

func main() {
	memCashe := NewMemoryStore()
	mux := http.NewServeMux()
	mux.HandleFunc(`/list`, memCashe.getMetricsHandler)
	mux.HandleFunc(`/update/`, memCashe.postMetricsHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		fmt.Print(err)
	}
}

func (m *MemStorage) getMetricsHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := json.MarshalIndent(m.getList(), "", "    ")
	if err != nil {
		fmt.Print(err)
	}
	res.Write(data)
	fmt.Println(data)
}

func (m *MemStorage) postMetricsHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	uri := req.RequestURI
	uri = strings.Replace(uri, "/update/", "", 1)
	values := strings.Split(uri, "/")

	if len(values) < 3 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	switch values[0] {
	case "gauge":
		if _, err := strconv.ParseFloat(values[2], 64); err == nil {
			res.WriteHeader(http.StatusOK)
			m.change(values[1], values[2])
			fmt.Println(values[1] + ": " + values[2])
			return
		}
	case "counter":
		if _, err := strconv.Atoi(values[2]); err == nil && !strings.Contains(values[2], ".") {
			res.WriteHeader(http.StatusOK)
			m.add(values[1], values[2])
			fmt.Println(values[1] + ": " + values[2])
			return
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
	}
	res.WriteHeader(http.StatusBadRequest)
}
