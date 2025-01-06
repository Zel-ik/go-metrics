package main

import (
	"flag"
	"fmt"
	"os"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {

	if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	} else {
		// регистрируем переменную flagRunAddr
		// как аргумент -a со значением :8080 по умолчанию
		flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
		// парсим переданные серверу аргументы в зарегистрированные переменные
		flag.Parse()
	}

	fmt.Print(flagRunAddr)

}
