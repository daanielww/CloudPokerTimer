package main

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type blindStucture struct {
	Small    int64
	Big      int64
	Ante     int64
	Level    int64
	Name     string
	Duration int64
}

type UserGame struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Level     int64
	UserID    string
	Start     int64
	State     bool
	Remain    int64
	structure blindStucture
}

func main() {
	session, err := mgo.Dial("localhost") // connect to server
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}
	defer session.Close() // close the connection when main returns

	collection := session.DB("game").C("userGame") //make the collection

	bs := blindStucture{
		Small:    10,
		Big:      11,
		Ante:     12,
		Level:    13,
		Name:     "Daniel",
		Duration: 24,
	}

	user := UserGame{

		Level:     123,
		UserID:    "asdasdasd",
		Start:     2222,
		State:     true,
		Remain:    123123,
		structure: bs,
	}

	err = collection.Insert(user)

}
