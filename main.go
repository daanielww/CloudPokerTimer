package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database


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
	StartTime                 int64          `json:"StartTime" bson:"StartTime"`
	CurrentPausedTime         int64          `json:"CurrentPausedTime" bson:"CurrentPausedTime"`
	AccumulatedPausedDuration int64          `json:"AccumulatedPausedTime" bson:"AccumulatedPausedTime"`
	Paused                    bool           `json:"Paused" bson:"Paused"`
	GameInfo                  blindStructure `json:"GameInfo" bson:"GameInfo"`
}

type Email struct {
	Email string
}

func games(w http.ResponseWriter, r *http.Request) {

	u := Email{}

	json.NewDecoder(r.Body).Decode(&u)

	email := u.Email

	fmt.Println(email)

	game := UserGame{}

	err := db.C("gameInfo").Find(bson.M{"UserID": email}).One(&game)

	fmt.Println(err == nil)

	if err == nil {
		err := db.C("gameInfo").Remove(bson.M{"UserID": email})
		if err != nil {
			fmt.Printf("remove fail %v\n", err)
			os.Exit(1)
		}
	}

	row1 := row{
		Small:    1,
		Big:      2,
		Ante:     123,
		Level:    2222,
		Duration: 123123,
	}

	bs := blindStructure{
		Name:      "asd",
		AllLevels: []row{row1},
	}

	userGame := UserGame{
		CurrentLevel:              12312,
		UserID:                    "Sammmmmmmmmmmmmmmmmmmy",
		StartTime:                 2231,
		CurrentPausedTime:         12123,
		AccumulatedPausedDuration: 1213,
		Paused:   true,
		GameInfo: bs,
	}

	db.C("gameInfo").Insert(userGame)

	uj, error := json.Marshal(userGame)
	if error != nil {
		fmt.Println(error)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(uj)
	w.WriteHeader(http.StatusCreated) // 201

}

func existing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["id"]

	game := UserGame{}

	err := db.C("gameInfo").Find(bson.M{"UserID": email}).One(&game)
	if err != nil {
		http.Error(w, "Error: game doesn't exist ", 404)
		return
	}

	uj, error := json.Marshal(game)
	if error != nil {
		fmt.Println(error)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(uj)
	w.WriteHeader(http.StatusCreated)
}

func main() {
	session, err := mgo.Dial("localhost") // connect to server
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}

	defer session.Close() // close the connection when main returns

	db = session.DB("game")

	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/login", GetUser).Methods("POST")
	router.HandleFunc("/games", games).Methods("POST")
	router.HandleFunc("/games/{id}", existing).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
