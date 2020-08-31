package routes

import (
	"encoding/json"
	"fmt"
	"github.com/chinanwu/solver"
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

// ROUTE: /ping
// METHOD: GET
// RETURNS: A response with "pong" if the server is up and running
// {
//		"message": "pong",
// }
func HandlePing(rw http.ResponseWriter, req *http.Request) {
	writeJSON(rw, Response{
		Message: "pong",
	})
}

type WordsResponse struct {
	Words []string `json:"words"`
}

// ROUTE: /allWords
// METHOD: GET
// RETURNS:
//		- The master list of possible game words
// 		{
//			"words": [ "word", "boop", "foop", .... ]
// 		}
//		- Or an error that occurred during the reading of the words list file
func HandleAllWords(rw http.ResponseWriter, req *http.Request) {
	arr, err := getWordsFromFile()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(rw, WordsResponse{
		Words: arr,
	})
}

type GameResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// ROUTE: /words
// METHOD: GET
// RETURNS:
//		- Two words for the game, "from" and "to".
// 		{
//			"from": "toot",
//			"to": "foot"
// 		}
//		- Or an error that occurred during the reading of the words list file
func HandleWords(rw http.ResponseWriter, req *http.Request) {
	rand.Seed(time.Now().UnixNano())

	arr, err := getWordsFromFile()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	max := len(arr)
	from := arr[rand.Intn(max)]
	for {
		to := arr[rand.Intn(max)]

		// Prevent games that are just... Way too easy
		if to != from {
			writeJSON(rw, GameResponse{
				From: from,
				To:   to,
			})
			break
		}
	}
}

type BoolResponse struct {
	Success bool `json:"success"`
}

// ROUTE: /validate
// METHOD: GET
// QUERY PARAMS: /validate?word=from[&word=toot&word=unau&...]
// RETURNS:
//		- A bool response with if the word(s) provided
// 		{
//			"success": [true/false]
// 		}
//		- Or an error that occurred during the reading of the words list file
func HandleValidate(rw http.ResponseWriter, req *http.Request) {
	words, _ := req.URL.Query()["word"]
	arr, err := getWordsFromFile()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
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

	// Question about this: Should I just return a Response with 200 OR
	// reply with a 404 Not Found?
	writeJSON(rw, BoolResponse{
		Success: false,
	})
}

// ROUTE: /score
// METHOD: GET
// QUERY PARAMS: /score?from=from&to=unau
// RETURNS: TBA
func HandleScore(rw http.ResponseWriter, req *http.Request) {
	from := req.FormValue("from")
	to := req.FormValue("to")

	// How to calculate score?
	// - The solution given by player vs shortest path length from solver.
	// - Time
	// - If word begins with z, add to score, "z bonus"
	// - If word begins with q, add to score, "q bonus"
	// - Number of different letters - Small, base score?
	// - If last two letters of a word are both consonants and the other isn't.
	// 		- ru[ns] -> fr[om]

	fmt.Println(from, to)
}

type SolveResponse struct {
	Solution []string `json:"solution"`
}

// ROUTE: /solve
// METHOD: GET
// QUERY PARAMS: /solve?from=from&to=unau
// RETURNS:
//		- A array of strings that is the solution to getting from "from" to "to.
// 		{
//			"solution": ["heat", "meat", "mead", "meld", "mold", "cold"]
// 		}
//		- Or an error when from or to was not provided
//		- Or an error that occurred during the reading of the words list file
// 		- Or an error that occured during the solving of the words.
func HandleSolve(rw http.ResponseWriter, req *http.Request) {
	from := req.FormValue("from")
	to := req.FormValue("to")
	if from == "" || to == "" {
		http.Error(rw, "Invalid words provided", http.StatusBadRequest)
		return
	}

	arr, err := getWordsFromFile()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	path, _, err := solver.Solve(from, to, arr, 4)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(rw, SolveResponse{
		Solution: path,
	})
}

type HintResponse struct {
	Hint    string `json:"hint"`
	NumLeft int    `json:"numLeft"`
}

// ROUTE: /hint
// METHOD: GET
// QUERY PARAMS: /solve?from=from&to=unau
// RETURNS:
//		- The next suggested word, and the number of remaining steps (including the "to")
// 		{
//			"hint": "meat",
//			"numLeft": 4
// 		}
//		- Or an error when from or to was not provided
//		- Or an error that occurred during the reading of the words list file
// 		- Or an error that occurred during the solving of the words.
func HandleHint(rw http.ResponseWriter, req *http.Request) {
	from := req.FormValue("from")
	to := req.FormValue("to")
	if from == "" || to == "" {
		http.Error(rw, "Invalid words provided", http.StatusBadRequest)
		return
	}

	arr, err := getWordsFromFile()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	path, _, err := solver.Solve(from, to, arr, 4)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(rw, HintResponse{
		Hint:    path[1],
		NumLeft: len(path) - 2,
	})
}

// HELPERS

func getWordsFromFile() ([]string, error) {
	words, err := ioutil.ReadFile("./assets/words.txt")

	if err != nil {
		return nil, err
	}

	arr := strings.Split(string(words), " ")
	return arr, nil
}

func writeJSON(rw http.ResponseWriter, resp interface{}) {
	// Marshalling allows us to convert Go structs into JSON
	j, err := json.Marshal(resp)
	if err != nil {
		http.Error(rw, "unable to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, writeErr := rw.Write(j)

	if writeErr != nil {
		// Do http.Error here? Or log?
		log.Printf("unable to write response json: %v" + writeErr.Error())
	}
}
