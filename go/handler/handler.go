package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate"
	"github.com/gorilla/mux"
	"github.com/iancoleman/strcase"
	"github.com/t-9/jingu-character-service/go/entity"
	"github.com/t-9/jingu-character-service/go/repository"
)

// RootHandler handles access to the root directory.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")
}

// ListHandler handles access to the list API.
func ListHandler(w http.ResponseWriter, r *http.Request) {
	chars, err := repository.RetrieveAll()
	j, err := json.Marshal(chars)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// StoreHandler handles access to the Store API.
func StoreHandler(w http.ResponseWriter, r *http.Request) {
	p := make(map[string]string, len(entity.Fillable))
	for _, v := range entity.Fillable {
		p[v] = r.FormValue(strcase.ToSnake(v))
	}
	c, err := entity.Create(p)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("location", fmt.Sprintf("/v1/show/%d", c.ID))
	w.WriteHeader(http.StatusCreated)

	j, err := json.Marshal(c)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"ID\":%d}", c.ID)))
		return
	}

	w.Write(j)
}

// DestroyHandler handles to access to the destroy API.
func DestroyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}

	c := entity.Character{
		ID: uint(id),
	}

	if err := c.Destroy(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ShowHandler handles to access to the show API.
func ShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}

	c, err := repository.FindByID(uint(id))
	if c.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	j, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdateHandler handles to access to the update API.
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	}

	c, err := repository.FindByID(uint(id))
	if c.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, v := range entity.Fillable {
		c.SetFieldByName(v, r.FormValue(strcase.ToSnake(v)))
	}

	j, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := c.Update(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// CreateDBHandler handles to access to the create DB API.
func CreateDBHandler(w http.ResponseWriter, r *http.Request) {
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

// MigrateUpDBHandler handles to access to the migrate up DB API.
func MigrateUpDBHandler(w http.ResponseWriter, r *http.Request) {
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

// MigrateDownDBHandler handles to access to the migrate down DB API.
func MigrateDownDBHandler(w http.ResponseWriter, r *http.Request) {
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

// MigrateForceDBHandler handles to access to the migrate force DB API.
func MigrateForceDBHandler(w http.ResponseWriter, r *http.Request) {
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
