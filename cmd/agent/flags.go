package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string
var reportInterval, pollInterval int

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {

	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "time postpone sending")
	flag.IntVar(&pollInterval, "p", 2, "time postpone metrics saving")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		pollInterval, _ = strconv.Atoi(envPollInterval)
	}

	if envReportIternal := os.Getenv("REPORT_INTERVAL"); envReportIternal != "" {
		reportInterval, _ = strconv.Atoi(envReportIternal)
	}
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	fmt.Printf("%s, %d, %d", flagRunAddr, pollInterval, reportInterval)
}
