package main

import (
	"os"

	"github.com/EgorO1999/go_final_project/pkg/db"
	"github.com/EgorO1999/go_final_project/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		panic(err)
	}

	if err := server.Run(); err != nil {
		panic(err)
	}
}
