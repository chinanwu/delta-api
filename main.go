package main

import (
	"flag"
	"github.com/chinanwu/delta-api/routes"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Printf("[STARTUP] Server starting...")

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/ping", routes.HandlePing).Methods(http.MethodGet)
	api.HandleFunc("/allWords", routes.HandleAllWords).Methods(http.MethodGet)
	api.HandleFunc("/words", routes.HandleWords).Methods(http.MethodGet)
	api.HandleFunc("/validate", routes.HandleValidate).Methods(http.MethodGet)
	api.HandleFunc("/score", routes.HandleScore).Methods(http.MethodGet)
	api.HandleFunc("/solve", routes.HandleSolve).Methods(http.MethodGet)
	api.HandleFunc("/hint", routes.HandleHint).Methods(http.MethodGet)

	// Not comfortable with this. Need to figure out how to best set allowed origins
	// Will need to do before I deploy
	var allowedOrigin string
	flag.StringVar(
		&allowedOrigin,
		"allowedOrigin",
		os.Getenv("ALLOWED_ORIGIN"),
		"Allowed origin for CORS")
	flag.Parse()

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins: []string{allowedOrigin},
		//AllowedOrigins: []string{os.Getenv("ALLOWED_ORIGIN")},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
		MaxAge:         1728000, // 20 days
	})

	log.Fatal(http.ListenAndServe(":8080", corsWrapper.Handler(r)))
}
