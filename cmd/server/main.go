package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
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

func MetricRouter() chi.Router {
	chiRouter := chi.NewRouter()
	memCashe := NewMemoryStore()

	chiRouter.Route("/update", func(r chi.Router) {
		r.Route("/{type}/{name}/{value}", func(r chi.Router) {
			r.Post("/", memCashe.postMetricsHandler)
		})
	})
	chiRouter.Get("/list", memCashe.getMetricsHandler)
	return chiRouter
}

func main() {

	err := http.ListenAndServe(`:8080`, MetricRouter())
	if err != nil {
		fmt.Print(err)
	}
}

func (m *MemStorage) postMetricsHandler(res http.ResponseWriter, req *http.Request) {
	if chi.URLParam(req, "type") == "" || chi.URLParam(req, "name") == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	switch chi.URLParam(req, "type") {
	case "gauge":
		if _, err := strconv.ParseFloat(chi.URLParam(req, "value"), 64); err == nil {
			res.WriteHeader(http.StatusOK)
			m.change(chi.URLParam(req, "name"), chi.URLParam(req, "value"))
			fmt.Println(chi.URLParam(req, "name") + ": " + chi.URLParam(req, "value"))
			return
		}
	case "counter":
		if _, err := strconv.Atoi(chi.URLParam(req, "value")); err == nil && !strings.Contains(chi.URLParam(req, "value"), ".") {
			res.WriteHeader(http.StatusOK)
			m.add(chi.URLParam(req, "name"), chi.URLParam(req, "value"))
			fmt.Println(chi.URLParam(req, "name") + ": " + chi.URLParam(req, "value"))
			return
		}
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusBadRequest)
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
