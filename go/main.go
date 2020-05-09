package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/gorilla/mux"
	"github.com/t-9/jingu-character-service/go/handler"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.RootHandler)

	r.HandleFunc("/v1/list", handler.ListHandler).Methods("GET")
	r.HandleFunc("/v1/store", handler.StoreHandler).Methods("POST")
	r.HandleFunc(
		"/v1/destroy/{id:[0-9]+}", handler.DestroyHandler).Methods("DELETE")
	r.HandleFunc(
		"/v1/show/{id:[0-9]+}", handler.ShowHandler).Methods("GET")
	r.HandleFunc(
		"/v1/update/{id:[0-9]+}", handler.UpdateHandler).Methods("PUT")

	r.HandleFunc(
		"/v1/admin/db/create", handler.CreateDBHandler).Methods("POST")
	r.HandleFunc(
		"/v1/admin/db/migrate/up", handler.MigrateUpDBHandler).Methods("POST")
	r.HandleFunc(
		"/v1/admin/db/migrate/down",
		handler.MigrateDownDBHandler).Methods("POST")
	r.HandleFunc(
		"/v1/admin/db/migrate/force/{id:[0-9]+}",
		handler.MigrateForceDBHandler).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
