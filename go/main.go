package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
