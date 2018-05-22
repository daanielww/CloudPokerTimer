package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2"
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
	CurrentLevel              int64          `json:"CurrentLevel" bson:"CurrentLevel"`
	UserID                    string         `json:"UserID" bson:"UserID"`
	StartTime                 int64          `json:"StartTime" bson:"StartTime"`
	CurrentPausedTime         int64          `json:"CurrentPausedTime" bson:"CurrentPausedTime"`
	AccumulatedPausedDuration int64          `json:"AccumulatedPausedTime" bson:"AccumulatedPausedTime"`
	Paused                    bool           `json:"Paused" bson:"Paused"`
	GameInfo                  blindStructure `json:"GameInfo" bson:"GameInfo"`
}

//asd
func GetPerson(w http.ResponseWriter, req *http.Request) {

}

//asdasd
func CreatePerson(w http.ResponseWriter, req *http.Request) {
	row1 := row{
		Small:    10,
		Big:      11,
		Ante:     12,
		Level:    13,
		Duration: 24,
	}

	row2 := row{
		Small:    123,
		Big:      1202021,
		Ante:     12132,
		Level:    12223,
		Duration: 20,
	}

	bs := blindStructure{
		Name:      "Texas holdem",
		AllLevels: []row{row1, row2},
	}

	user := UserGame{

		CurrentLevel:              123,
		UserID:                    "asdasdasd",
		StartTime:                 2222,
		CurrentPausedTime:         333,
		AccumulatedPausedDuration: 9229,
		Paused:   true,
		GameInfo: bs,
	}

	err := db.C("userGame").Insert(user)
	if err != nil {
		log.Fatal("blah", err)
	}
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	session, err := mgo.Dial("localhost") // connect to server
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}

	defer session.Close() // close the connection when main returns

	db = session.DB("game")

	router := mux.NewRouter()
	router.HandleFunc("/", CreatePerson).Methods("GET")
	router.HandleFunc("/get", GetPerson).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
