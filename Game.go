package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func LoadExistingGame(email string, collection *mgo.Collection) (string, error) {
	if collection == nil {
		return "", errors.New("Collection is nil")
	}

	if len(strings.TrimSpace(email)) == 0 {
		return "", errors.New("email is empty")
	}

	// retrieve user to see if they exist
	result := UserGame{}
	if collection.Find(bson.M{"UserID": email}).One(&result) == nil {
		return "", errors.New("Could not find game")
	}

	resultString, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("Could not marshal to json")
	}

	return string(resultString), nil
}
