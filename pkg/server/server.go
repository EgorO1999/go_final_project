package server

import (
	"net/http"
	"os"

	"github.com/EgorO1999/go_final_project/pkg/api"
)

func Run() error {
	port := "7540"

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	webDir := "web"

	api.Init()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	return http.ListenAndServe(":"+port, nil)
}
