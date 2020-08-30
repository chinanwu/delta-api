package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
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
	fmt.Println(words[0])

	// Possible for multiple "word" to be passed in
	//words := strings.Split(wordStr, ",")

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

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/ping", handlePing).Methods(http.MethodGet)
	api.HandleFunc("/allWords", handleAllWords).Methods(http.MethodGet)
	api.HandleFunc("/words", handleWords).Methods(http.MethodGet)
	api.HandleFunc("/validate", handleValidate).Methods(http.MethodGet)

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

func getWordsFromFile() ([]string, error) {
	words, err := ioutil.ReadFile("./assets/words.txt")

	if err != nil {
		return nil, err
	}

	arr := strings.Split(string(words), " ")
	return arr, nil
}
