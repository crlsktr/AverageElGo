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
	file, _ := ioutil.ReadFile("data.pgn") //this is probably pretty bad since we load all of the file all at once.
	filetxt := string(file)
	// w.Write(file)
	w.Header().Set("Content-Type", "application/json")
	var database pgn.DB
	database.Parse(filetxt)

	var myratings []*MyGame
	totalrating := 0
	var averagerating float64

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

		totalrating += rating
		averagerating = float64(totalrating) / float64(len(myratings))
	}
	jsongames, _ := json.Marshal(averagerating)
	w.Write(jsongames)
}
