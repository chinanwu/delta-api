package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/chinanwu/solver"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Response struct {
	Message string `json:"message"`
}

type WordsResponse struct {
	Words []string `json:"words"`
}

type GameResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type BoolResponse struct {
	Success bool `json:"success"`
}

type SolveResponse struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	Solution []string `json:"solution"`
}

type HintResponse struct {
	Hint    string `json:"hint"`
	NumLeft int    `json:"numLeft"`
}

// /ping
func handlePing(rw http.ResponseWriter, req *http.Request) {
	writeJSON(rw, Response{
		Message: "pong",
	})
}

// /allWords
func handleAllWords(rw http.ResponseWriter, req *http.Request) {
	arr, err := getWordsFromFile()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	writeJSON(rw, WordsResponse{
		Words: arr,
	})
}

// /words
func handleWords(rw http.ResponseWriter, req *http.Request) {
	rand.Seed(time.Now().UnixNano())

	arr, err := getWordsFromFile()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	max := len(arr)
	from := arr[rand.Intn(max)]
	for {
		to := arr[rand.Intn(max)]
		if to != from {
			writeJSON(rw, GameResponse{
				From: from,
				To:   to,
			})
			break
		}
	}
}

func handleValidate(rw http.ResponseWriter, req *http.Request) {
	words, _ := req.URL.Query()["word"]
	arr, err := getWordsFromFile()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	numFound := 0
	expected := len(words)
	for _, elem := range arr {
		for _, word := range words {
			if word == elem {
				numFound++
			}
		}

		if numFound == expected {
			writeJSON(rw, BoolResponse{
				Success: true,
			})
			return
		}
	}

	writeJSON(rw, BoolResponse{
		Success: false,
	})
}

func handleScore(rw http.ResponseWriter, req *http.Request) {
	from := req.FormValue("from")
	to := req.FormValue("to")

	// Coming soon

	fmt.Println(from, to)
}

func handleSolve(rw http.ResponseWriter, req *http.Request) {
	from := req.FormValue("from")
	to := req.FormValue("to")

	arr, err := getWordsFromFile()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	path, _, err := solver.Solve(from, to, arr, 4)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	writeJSON(rw, SolveResponse{
		From:     from,
		To:       to,
		Solution: path,
	})
}

func handleHint(rw http.ResponseWriter, req *http.Request) {
	from := req.FormValue("from")
	to := req.FormValue("to")

	arr, err := getWordsFromFile()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	path, _, err := solver.Solve(from, to, arr, 4)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	writeJSON(rw, HintResponse{
		Hint:    path[1],
		NumLeft: len(path) - 2,
	})
}

func main() {
	log.Printf("[STARTUP] Server starting...")

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/ping", handlePing).Methods(http.MethodGet)
	api.HandleFunc("/allWords", handleAllWords).Methods(http.MethodGet)
	api.HandleFunc("/words", handleWords).Methods(http.MethodGet)
	api.HandleFunc("/validate", handleValidate).Methods(http.MethodGet)
	api.HandleFunc("/score", handleScore).Methods(http.MethodGet)
	api.HandleFunc("/solve", handleSolve).Methods(http.MethodGet)
	api.HandleFunc("/hint", handleHint).Methods(http.MethodGet)

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

func writeJSON(rw http.ResponseWriter, resp interface{}) {
	j, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, "unable to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, writeErr := rw.Write(j)

	if writeErr != nil {
		// Do http.Error here? Or log?
		log.Printf("unable to write response json: " + writeErr.Error())
	}
}

func getWordsFromFile() ([]string, error) {
	words, err := ioutil.ReadFile("./assets/words.txt")

	if err != nil {
		return nil, err
	}

	arr := strings.Split(string(words), " ")
	return arr, nil
}
