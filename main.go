package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

var db *mgo.Database

type user struct {
	Email    string  `json: "email" bson: "email"`
	Username string  `json: "username" bson: "username"`
	Password string  `json: "pass" bson: "pass"`
	PHash    *[]byte `json:"-", omitempty`
}

type blindStructure struct {
	Name      string `json:"Name" bson:"Name"`
	AllLevels []row  `json:"AllLevels" bson:"AllLevels"`
}

type row struct {
	Small    int64 `json:"Small" bson:"Small"`
	Big      int64 `json:"Big" bson:"Big"`
	Ante     int64 `json:"Ante" bson:"Ante"`
	Level    int64 `json:"Level" bson:"Level"`
	Duration int64 `json:"Duration" bson:"Duration"`
}

type UserGame struct {
	CurrentLevel              int64          `jsons:"CurrentLevel" bson:"CurrentLevel"`
	UserID                    string         `json:"UserID" bson:"UserID"`
	StartTime                 time.Time      `json:"StartTime" bson:"StartTime"`
	CurrentPausedTime         time.Time      `json:"CurrentPausedTime" bson:"CurrentPausedTime"`
	AccumulatedPausedDuration float64        `json:"AccumulatedPausedTime" bson:"AccumulatedPausedTime"`
	CurrentPausedStartTime    time.Time      `json:"CurrentPausedStartTime" bson:"CurrentPausedStartTime"`
	CurrentLevelTime          float64        `json:"CurrentLevelTime" bson:"CurrentLevelTime"`
	Paused                    bool           `json:"Paused" bson:"Paused"`
	GameInfo                  blindStructure `json:"GameInfo" bson:"GameInfo"`
}

func main() {
	session, err := mgo.Dial("localhost") // connect to server
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}

	defer session.Close() // close the connection when main returns

	db = session.DB("game")

	router := mux.NewRouter()
	router.HandleFunc("/", CreateUser).Methods("POST")
	router.HandleFunc("/login", GetUser).Methods("GET")
	router.HandleFunc("/games", games).Methods("POST")
	router.HandleFunc("/games/{id}", existing).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
