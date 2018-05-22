package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2"
)

var db *mgo.Database

type user struct {
	Email string `json: "email" bson: "email"`
	Username string `json: "username" bson: "username"`
	Password string `json: "pass" bson: "pass"`
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
	CurrentLevel              int64          `json:"CurrentLevel" bson:"CurrentLevel"`
	UserID                    string         `json:"UserID" bson:"UserID"`
	StartTime                 int64          `json:"StartTime" bson:"StartTime"`
	CurrentPausedTime         int64          `json:"CurrentPausedTime" bson:"CurrentPausedTime"`
	AccumulatedPausedDuration int64          `json:"AccumulatedPausedTime" bson:"AccumulatedPausedTime"`
	Paused                    bool           `json:"Paused" bson:"Paused"`
	GameInfo                  blindStructure `json:"GameInfo" bson:"GameInfo"`
}

type CurrentGame struct {
	User                      string  
	StartTime                 int64
	Paused                    bool
	CurrentPausedStartTime    int64
	CurrentLevelTime          int64
	CurrentLevel 	          int64
	GameInfo                  blindStructure
}


func main() {
	session, err := mgo.Dial("localhost") // connect to server
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}

	defer session.Close() // close the connection when main returns

	db = session.DB("game")

	router := mux.NewRouter()
	router.HandleFunc("/", CreateUser).Methods("GET")
	router.HandleFunc("/get", GetUser).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
