package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"./handlers"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/so2/clima/temperatura", handlers.HandleGetTemperature).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", r))
}
