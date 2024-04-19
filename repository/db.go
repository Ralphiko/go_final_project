package repository

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	EnvDbfile             = "TODO_DBFILE"
	DbNameFromEnvironment = "Получаем имя БД из окружения..."
	DbNameSet             = "Имя БД задано -- %v"
	NotSetUsingDefault    = "Имя БД не задано. Будем использовать из конфига -- %v"
	SqlDriver             = "sqlite3"
	FailedToOpenDatabase  = "Не удалось открыть БД: "
	TableCreationError    = "Упс!.. Ошбика при создании таблицы: "
	IndexCreationError    = "Упс!.. Ошбика при создании индекса: "
	TaskTable             = "scheduler"
)

type Config struct {
	SQLDriver string
	DBFile    string
}

func GetDB() *sqlx.DB {
	dbname, err := CheckDb()
	if err != nil {
		logrus.Fatal(err)
	}
	return sqlx.MustConnect(viper.Get("DB.SQLDriver").(string), dbname)
}

func CheckDb() (string, error) {
	dbName := EnvDBFILE(EnvDbfile)

	appPath, err := os.Executable()
	if err != nil {
		logrus.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbName)
	_, err = os.Stat(dbFile)
	if err != nil {
		installDB(dbName)
	}
	return dbName, nil
}

func EnvDBFILE(key string) string {
	logrus.Println(DbNameFromEnvironment)
	dbName := os.Getenv(key)
	if len(dbName) == 0 {
		dbName = viper.Get("DB.DBFile").(string)
		logrus.Warnf(NotSetUsingDefault, dbName)
	} else {
		logrus.Printf(DbNameSet, dbName)
	}
	return dbName
}

func installDB(dbName string) {
	db, err := sql.Open(SqlDriver, dbName)
	if err != nil {
		logrus.Fatal(FailedToOpenDatabase, err)
	}
	defer db.Close()

	createTableSQL := viper.Get("DB.SQLCreateTables").(string)
	_, err = db.Exec(createTableSQL)
	if err != nil {
		logrus.Fatal(TableCreationError, err)
	}

	createIndexSQL := viper.Get("DB.SQLCreateIndexes").(string)
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		logrus.Fatal(IndexCreationError, err)
	}
}
