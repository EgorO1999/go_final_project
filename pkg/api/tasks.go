package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/EgorO1999/go_final_project/pkg/db"
	"github.com/EgorO1999/go_final_project/pkg/rule"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")
	var tasks []*db.Task
	var err error

	if search != "" {
		tasks, err = db.SearchTasks(search, 50)
	} else {
		tasks, err = db.Tasks(50)
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err := rule.CheckDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("ошибка добавления задачи: %v", err)})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный формат данных"})
		return
	}

	if task.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Заголовок задачи обязателен"})
		return
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный формат даты"})
		return
	}

	if task.Repeat != "" {
		if _, err := rule.NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверное правило повторения"})
			return
		}
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{})
		return
	}

	now := time.Now()
	nextDate, err := rule.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверное правило повторения"})
		return
	}

	if err := db.UpdateDate(nextDate, task.ID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{})
}
