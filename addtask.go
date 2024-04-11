package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"./nextday"
	_ "github.com/go-sql-driver/mysql"
)

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	ID string `json:"id"`
}

func addTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		respondWithError(w, "Ошибка при декодировании JSON")
		return
	}

	if task.Title == "" {
		respondWithError(w, "Не указан заголовок задачи")
		return
	}

	// Проверка формата даты
	if _, err := time.Parse("20060102", task.Date); err != nil {
		respondWithError(w, "Неверный формат даты. Используйте формат 20060102")
		return
	}

	// Обработка даты, если она не указана или меньше сегодняшней
	currentDate := time.Now().Format("20060102")
	if task.Date == "" || task.Date < currentDate {
		if task.Repeat != "" {
			nextDate, err := nextday.NextDate(currentDate, task.Repeat)
			if err != nil {
				respondWithError(w, err.Error())
				return
			}
			task.Date = nextDate
		} else {
			task.Date = currentDate
		}
	}

	// Добавление задачи в базу данных
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname")
	if err != nil {
		respondWithError(w, "Ошибка при подключении к базе данных")
		return
	}
	defer db.Close()

	query := "INSERT INTO tasks (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		respondWithError(w, "Ошибка при добавлении задачи в базу данных")
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		respondWithError(w, "Ошибка при получении идентификатора задачи")
		return
	}

	respondWithSuccess(w, strconv.FormatInt(id, 10))
}

func respondWithError(w http.ResponseWriter, errMsg string) {
	errResp := ErrorResponse{Error: errMsg}
	respondJSON(w, errResp, http.StatusBadRequest)
}

func respondWithSuccess(w http.ResponseWriter, id string) {
	successResp := SuccessResponse{ID: id}
	respondJSON(w, successResp, http.StatusOK)
}

func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
