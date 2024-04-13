package contoller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type TaskIdResponse struct {
	Id string `json:"id"`
}

func task(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		taskPost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func taskPost(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	task, err = checkTask(&task)
	if err != nil {
		responseError := ErrorResponse{Error: err.Error()}
		jsonResponse, err := json.Marshal(responseError)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	// Проверяем ошибку при вызове connectDB()
	db, err := connectDB()
	if err != nil {
		// Если произошла ошибка, обрабатываем ее
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Если ошибок нет, продолжаем выполнение кода
	// Вызываем функцию добавления задачи, передавая полученное соединение с базой данных
	lastId, err := addingTask(db, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseTaskId := TaskIdResponse{Id: strconv.FormatInt(lastId, 10)}
	jsonResponse, err := json.Marshal(responseTaskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

const (
	TITLE_NOT_SET = "Заголовок не может быть пустым!"
)

func checkTask(m *Task) (Task, error) {
	if strings.TrimSpace(m.Title) == "" {
		return Task{}, fmt.Errorf(TITLE_NOT_SET)
	}

	now := time.Now()
	if m.Date == "" {
		m.Date = now.Format("20060102")
	}

	_, err := time.Parse("20060102", m.Date)
	if err != nil {
		return Task{}, fmt.Errorf("Не могу преобразовать дату!")
	}

	if m.Date < now.Format("20060102") {
		if m.Repeat == "" {
			m.Date = now.Format("20060102")
		} else {
			m.Date, err = NextDate(now, m.Date, m.Repeat)
			if err != nil {
				return Task{}, err
			}
		}
	}

	return *m, nil
}

func addingTask(db *sql.DB, task Task) (int64, error) {
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
