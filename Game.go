package main

import (
	"net/http"
	"os"
	"time"

	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type Email struct {
	Email string
}

type Status struct {
	Status bool
	UserID string
}

/*
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

	levelInfos := result.GameInfo.AllLevels
	var levelMinutes []int

	for _, levelInfo := range levelInfos {
		levelMinutes = append(levelMinutes, levelInfo.Duration)
	}

	duration, _ := time.ParseDuration(strconv.FormatFloat(result.AccumulatedPausedDuration, 'g', 1, 64) + "s")

	level, levelTime := GetLevelAndLevelTime(result.StartTime, duration, result.Paused, result.CurrentPausedTime, levelMinutes)

	currGame := CurrentGame{
		User:                   result.UserID,
		StartTime:              result.StartTime,
		Paused:                 result.Paused,
		CurrentPausedStartTime: result.CurrentPausedTime,
		CurrentLevelTime:       levelTime,
		CurrentLevel:           level,
		GameInfo:               result.GameInfo,
	}

	resultString, err := json.Marshal(currGame)
	if err != nil {
		fmt.Printf("Could not marshal to json")
	}

	return string(resultString), nil
}
*/

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
		SmallBlind: 5,
		BigBlind:   10,
		Ante:       0,
		Level:      1,
		Duration:   11,
	}

	row2 := row{
		SmallBlind: 10,
		BigBlind:   20,
		Ante:       0,
		Level:      2,
		Duration:   12,
	}

	userGame := UserGame{
		User:                      "Sam",
		StartTime:                 time.Now(),
		Paused:                    false,
		CurrentPausedStartTime:    time.Now(),
		CurrentLevelTime:          555,
		CurrentLevel:              6,
		BlindScheduleName:         "Office Turbo",
		Levels:                    []row{row1, row2},
		CurrentPausedTime:         time.Now(),
		AccumulatedPausedDuration: 1213,
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

	err := db.C("gameInfo").Find(bson.M{"User": email}).One(&game)
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

func update(w http.ResponseWriter, r *http.Request) {
	u := Status{}
	json.NewDecoder(r.Body).Decode(&u)

	//find the record
	game := UserGame{}
	err := db.C("gameInfo").Find(bson.M{"UserID": u.UserID}).One(&game)
	if err != nil {
		fmt.Printf("find fail %v\n", err)
		os.Exit(1)
	}

	db.C("gameInfo").Update(bson.M{"User": game.User}, bson.M{"$set": bson.M{"Paused": u.Status}})
}
