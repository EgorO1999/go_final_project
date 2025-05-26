package main

import (
	"net/http"
	"log"
	"fmt"
)

func main() {
	// Определение порта, если переменная окружения не установлена, то используем 7540
	port := "7540"

	// Директория с файлами фронтенда
	webDir := "./web"

	// Обработчик для статических файлов
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Запуск сервера на порту 7540
	fmt.Printf("Сервер запущен на http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}