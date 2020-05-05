package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func dbHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("url path is " + r.URL.Path[1:] + "\n"))

	var dbConnectQuery string
	dbConnectQuery = "root:" + os.Getenv("DB_PASSWORD") + "@tcp(" +
		os.Getenv("DB_HOST") + ":3306)/" + os.Getenv("DB_NAME")
	db, err := sql.Open("mysql", dbConnectQuery)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM user") //
	if err != nil {
		panic(err.Error())
	}

	columns, err := rows.Columns() // カラム名を取得
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
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

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/db", dbHandler)

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
