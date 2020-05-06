package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Character represents a character.
type Character struct {
	ID        uint      `gorm:"primary_key"`
	Surname   *string   `gorm:"column:surname"`
	GivenName string    `gorm:"column:given_name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(
		os.Getenv("DB_DRIVER"),
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		),
	)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			w.Write([]byte(err.Error()))
		}
	}()

	var characters []Character
	db.Find(&characters)

	j, err := json.Marshal(characters)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(j)
}

func dbCreateHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(
		os.Getenv("DB_DRIVER"),
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/mysql",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
		),
	)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			w.Write([]byte(err.Error()))
		}
	}()

	_, err = db.Exec(
		fmt.Sprintf(
			"CREATE DATABASE IF NOT EXISTS %s",
			os.Getenv("DB_NAME"),
		),
	)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("Successfully created database.."))
}

func dbMigrateUpHandler(w http.ResponseWriter, r *http.Request) {
	m, err := migrate.New(
		"file://./db/migrations/",
		fmt.Sprintf(
			"%s://%s:%s@tcp(%s:%s)/%s",
			os.Getenv("DB_DRIVER"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		),
	)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	err = m.Up()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("Successfully up migration version"))
}

func dbMigrateDownHandler(w http.ResponseWriter, r *http.Request) {
	m, err := migrate.New(
		"file://./db/migrations/",
		fmt.Sprintf(
			"%s://%s:%s@tcp(%s:%s)/%s",
			os.Getenv("DB_DRIVER"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		),
	)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	err = m.Down()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("Successfully up migration version"))
}

func dbAddHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(
		os.Getenv("DB_DRIVER"),
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		),
	)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			w.Write([]byte(err.Error()))
		}
	}()

	_, err = db.Exec(
		fmt.Sprintf(
			"INSERT INTO characters(surname, given_name) VALUES('%s', '%s')",
			"津島",
			"善子",
		),
	)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("Successfully insert a record!"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/db/create", dbCreateHandler)
	r.HandleFunc("/db/migrate/up", dbMigrateUpHandler)
	r.HandleFunc("/db/migrate/down", dbMigrateDownHandler)
	r.HandleFunc("/db/add", dbAddHandler)

	r.HandleFunc("/v1/list", listHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
