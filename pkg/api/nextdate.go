package api

import (
	"log"
	"net/http"
	"time"

	"github.com/EgorO1999/go_final_project/pkg/rule"
)

const DateFormat = "20060102"

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	if nowStr == "" {
		nowStr = time.Now().Format(DateFormat)
	}

	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if dateStr == "" || repeat == "" {
		http.Error(w, "missing date or repeat parameter", http.StatusBadRequest)
		return
	}

	now, err := time.Parse(DateFormat, nowStr)
	if err != nil {
		http.Error(w, "invalid now date format", http.StatusBadRequest)
		return
	}

	next, err := rule.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(next)); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
