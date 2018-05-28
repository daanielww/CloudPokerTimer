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

	/*u := Email{}

	json.NewDecoder(r.Body).Decode(&u)

	email := u.Email

	fmt.Println(email) */

	email := "Sam"

	game := UserGame{}

	err := db.C("gameInfo").Find(bson.M{"UserID": email}).One(&game)

	if err == nil {
		err := db.C("gameInfo").Remove(bson.M{"UserID": email})
		if err != nil {
			fmt.Printf("remove fail %v\n", err)
			os.Exit(1)
		}
	}

	userGame := makeDummyData(email)

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
	w.WriteHeader(http.StatusOK)
	w.Write(uj)

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

func makeDummyData(email string) UserGame {

	smallBlindArray := []int64{5, 10, 25, 50, 75, 100, 150, 200, 300, 400, 500, 700, 1000, 1500, 2000, 3000}
	bigBlindArray := []int64{10, 20, 50, 100, 150, 200, 300, 400, 600, 800, 1000, 1400, 2000, 3000, 4000, 6000}
	anteArray := []int64{0, 0, 5, 10, 10, 25, 25, 25, 50, 50, 100, 100, 200, 300, 400, 600}
	durationArray := []int64{420, 420, 420, 420, 420, 420, 420, 420, 420, 420, 420, 420, 420, 420, 420, 420}
	rows := []row{}
	rowOb := row{}

	for i := 0; i < len(smallBlindArray); i++ {

		rowOb = row{
			SmallBlind: smallBlindArray[i],
			BigBlind:   bigBlindArray[i],
			Ante:       anteArray[i],
			Level:      int64(i + 1),
			Duration:   durationArray[i],
		}

		rows = append(rows, rowOb)
	}

	userGame := UserGame{
		User:                      email,
		StartTime:                 time.Now(),
		Paused:                    false,
		CurrentPausedStartTime:    time.Now(),
		CurrentLevelTime:          420,
		CurrentLevel:              1,
		BlindScheduleName:         "Office Turbo",
		Levels:                    rows,
		CurrentPausedTime:         time.Now(),
		AccumulatedPausedDuration: 1213,
	}

	return userGame
}
