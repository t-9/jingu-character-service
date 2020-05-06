package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func dbHandler(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT * FROM %s",
			os.Getenv("DB_NAME"),
		),
	)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			w.Write([]byte(columns[i] + ": " + value + "\n"))
		}
		fmt.Println("-----------------------------------")
	}
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
	r.HandleFunc("/db", dbHandler)
	r.HandleFunc("/db/create", dbCreateHandler)
	r.HandleFunc("/db/migrate/up", dbMigrateUpHandler)
	r.HandleFunc("/db/migrate/down", dbMigrateDownHandler)
	r.HandleFunc("/db/add", dbAddHandler)

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
