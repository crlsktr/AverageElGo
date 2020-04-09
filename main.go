package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/malbrecht/chess/pgn"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/getaveragerating", getChatRequest).Methods("GET")

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}

type MyGame struct {
	Rating int
	Date   string
}

func getChatRequest(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("data.pgn")
	filetxt := string(file)
	// w.Write(file)
	w.Header().Set("Content-Type", "application/json")
	var database pgn.DB
	database.Parse(filetxt)

	var myratings []*MyGame
	for _, game := range database.Games {

		var side string
		if game.Tags["BlackElo"] == "crlsktr" {
			side = "BlackElo"
		} else {
			side = "WhiteElo"
		}

		rating, _ := strconv.Atoi(game.Tags[side])
		myratings = append(myratings, &MyGame{
			Rating: rating,
			Date:   game.Tags["UTCDate"],
		})
	}
	jsongames, _ := json.Marshal(myratings)
	w.Write(jsongames)
}