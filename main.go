package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func handlePing(rw http.ResponseWriter, req *http.Request) {
	writeJSON(rw, Response{
		Message: "pong",
	})
}

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/ping", handlePing).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func writeJSON(rw http.ResponseWriter, resp interface{}) {
	j, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, "unable to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(j)
}
