package main

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

type user struct {
	Email string `json:"email" bson:"email"`
	Username string `json: "name" bson:"name"`
	Password string `json: "pass" bson:"pass"`
}

type blindStructure struct {
	Name      string
	AllLevels []row
}

type row struct {
	Small    int64
	Big      int64
	Ante     int64
	Level    int64
	Duration int64
}

type UserGame struct {
	Level                     int64
	UserID                    string
	StartTime                 int64
	CurrentPausedTime         int64
	AccumulatedPausedDuration int64
	Paused                    bool
	GameInfo                  blindStructure
}

func main() {

	mux := http.NewServeMux()


	session, err := mgo.Dial("localhost") // connect to server
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.ListenAndServe(":3000", nil)

	defer session.Close() // close the connection when main returns

	collection := session.DB("game").C("userGame") //make the collection

	uc := NewUserController(session)

	// Login
	mux.HandleFunc("/login", uc.CreateUser)

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

		Level:                     123,
		UserID:                    "asdasdasd",
		StartTime:                 2222,
		CurrentPausedTime:         333,
		AccumulatedPausedDuration: 9229,
		Paused:   true,
		GameInfo: bs,
	}

	err = collection.Insert(user)

}
