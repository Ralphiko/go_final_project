package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Ralphiko/go_final_project/model"
	_ "modernc.org/sqlite"
)

var database *sql.DB

func getDBFilePath() string {
	dbFilePath := os.Getenv("TODO_DBFILE")
	return dbFilePath
}

func InsertTask(task model.Task) (int, error) {
	result, err := database.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func connectDB() (*sql.DB, error) {
	dbFilePath := getDBFilePath()
	if dbFilePath == "" {
		dbFilePath = "scheduler.db"
	}

	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTable(db *sql.DB) error {
	createStmt := `
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			title TEXT,
			comment TEXT,
			repeat TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	`
	fmt.Println("Создание таблицы перед тестированием...")
	_, err := db.Exec(createStmt)
	if err != nil {
		return err
	}
	fmt.Println("Таблица успешно создана перед тестированием.")
	return nil
}

func InitializeDB() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	newFunction1(err, dbFile)

	dbFilePath := os.Getenv("TODO_DBFILE")
	if dbFilePath != "" {
		fmt.Printf("Используется пользовательский путь к базе данных: %s\n", dbFilePath)
	}

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createTable(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("База данных и таблица успешно созданы.")
}

func newFunction1(err error, dbFile string) {
	_, err = os.Stat(dbFile)
}

func SearchTasks(search string) ([]model.Task, error) {
	var tasks []model.Task

	search = fmt.Sprintf("%%%s%%", search)
	rows, err := database.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date",
		sql.Named("search", search))
	if err != nil {
		return []model.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []model.Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []model.Task{}, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, nil
}

func SearchTasksByDate(date string) ([]model.Task, error) {
	var tasks []model.Task

	rows, err := database.Query("SELECT * FROM scheduler WHERE date = :date",
		sql.Named("date", date))
	if err != nil {
		return []model.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []model.Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []model.Task{}, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, nil
}

func ReadTask(id int) (model.Task, error) {
	var task model.Task

	row := database.QueryRow("SELECT * FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func UpdateTask(task model.Task) (model.Task, error) {
	result, err := database.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.Id))
	if err != nil {
		return model.Task{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Task{}, err
	}

	if rowsAffected == 0 {
		return model.Task{}, errors.New("failed to update")
	}

	return task, nil
}

func DeleteTaskDb(id string) error {
	result, err := database.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("failed to delete")
	}

	return err
}

func ReadTasks() ([]model.Task, error) {
	rows, err := database.Query("SELECT * FROM scheduler")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
func InstallDb() {
	dbFilePath := getDbFilePath()
	_, err := os.Stat(dbFilePath)

	if err != nil {
		database, err = createDbFile(dbFilePath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		database, err = sql.Open("sqlite3", dbFilePath)
	}

	if err != nil {
		log.Fatal(err)
	}
	createTable(database)
}
func getDbFilePath() string {
	dbFilePath := "scheduler.db"

	envDbFilePath := os.Getenv("TODO_DBFILE")
	if len(envDbFilePath) > 0 {
		dbFilePath = envDbFilePath
	}

	return dbFilePath
}
func createDbFile(dbFilePath string) (*sql.DB, error) {
	_, err := os.Create(dbFilePath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	return db, nil
}
