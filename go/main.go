package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
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
	Surname   string    `gorm:"column:surname"`
	GivenName string    `gorm:"column:given_name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the character service.")
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openGorm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	var characters []Character
	db.Find(&characters)

	j, err := json.Marshal(characters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	surname := r.FormValue("surname")
	givenName := r.FormValue("given_name")

	db, err := openGorm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	c := Character{
		Surname:   surname,
		GivenName: givenName,
	}

	if !db.NewRecord(&c) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not create new record"))
		return
	}

	j, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := db.Create(&c).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func destroyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}

	db, err := openGorm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	c := Character{
		ID: uint(id),
	}

	if err := db.Delete(&c).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func dbCreateHandler(w http.ResponseWriter, r *http.Request) {
	db, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully created database.."))
}

func dbMigrateUpHandler(w http.ResponseWriter, r *http.Request) {
	m, err := makeMigrate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = m.Up()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully up migration version"))
}

func dbMigrateDownHandler(w http.ResponseWriter, r *http.Request) {
	m, err := makeMigrate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = m.Down()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully up migration version"))
}

func dbMigrateForceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	v, err := strconv.Atoi(vars["version"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid version"))
		return
	}

	m, err := makeMigrate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = m.Force(int(v))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully up migration version"))
}

func makeMigrate() (*migrate.Migrate, error) {
	return migrate.New(
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
}

func openDB() (*sql.DB, error) {
	return sql.Open(
		os.Getenv("DB_DRIVER"),
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/mysql",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
		),
	)
}

func openGorm() (*gorm.DB, error) {
	return gorm.Open(
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
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)

	r.HandleFunc("/v1/list", listHandler).Methods("GET")
	r.HandleFunc("/v1/store", storeHandler).Methods("POST")
	r.HandleFunc(
		"/v1/destroy/{id:[0-9]+}", destroyHandler).Methods("DELETE")

	r.HandleFunc(
		"/v1/admin/db/create", dbCreateHandler).Methods("POST")
	r.HandleFunc(
		"/v1/admin/db/migrate/up", dbMigrateUpHandler).Methods("POST")
	r.HandleFunc(
		"/v1/admin/db/migrate/down", dbMigrateDownHandler).Methods("POST")
	r.HandleFunc(
		"/v1/admin/db/migrate/force/{id:[0-9]+}",
		dbMigrateDownHandler).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
