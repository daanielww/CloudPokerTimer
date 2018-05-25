package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

//Runs asychronously as a goroutine. levelTimeRemaining and currentLevel are channels to return the values back to the caller.
//Currently requires four game related values in order to calculate levelTimeRemaining and currentLevel.
//Once the data format of the Game struct has been established, consider passing in the single struct instead.
func GetLevelAndLevelTime(startTime time.Time, accumulatedPauseDuration time.Duration, paused bool,
	currentPauseStartTime time.Time, levelMinutes []int) (int, float64) {
	currentTime := time.Now()
	gameElapsedTime := currentTime.Sub(startTime)
	fmt.Printf("Game Elapsed Time is %v\n", gameElapsedTime)

	if paused == true {
		accumulatedPauseDuration += currentPauseStartTime.Sub(currentTime)
	}

	gameDuration := gameElapsedTime - accumulatedPauseDuration
	fmt.Printf("Game Duration is %v\n", gameDuration)

	var curLevel = 1
	accumulatedLevelTime, _ := time.ParseDuration("0m")

	for i := 0; i <= len(levelMinutes)-1; i++ {
		NextLevelMinutesDuration, _ := time.ParseDuration(strconv.Itoa(levelMinutes[i+1]) + "m")
		curLevel = i + 1
		if accumulatedLevelTime+NextLevelMinutesDuration > gameDuration {
			break
		}
		levelMinutesDuration, _ := time.ParseDuration(strconv.Itoa(levelMinutes[i]) + "m")
		accumulatedLevelTime += levelMinutesDuration
	}

	fmt.Printf("Acc Level Time = :%v\n", accumulatedLevelTime)

	levelTimeRemaining := gameDuration - accumulatedLevelTime

	return curLevel, levelTimeRemaining.Seconds()
}

// function to return the email value given from the query parameters
func getEmail(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["id"]
}

// function to trigger pause and play action
func updateGamePauseState(w http.ResponseWriter, newGameState bool, r *http.Request) {
	fmt.Println("working")
	vars := mux.Vars(r)
	email := vars["id"]
	level := vars["level"]
	levelTime := vars["levelTime"]

	update := bson.M{"User": email}

	if newGameState == false {
		change := bson.M{"$set": bson.M{"Paused": newGameState}}
		err := db.C("gameInfo").Update(update, change)
		if err != nil {
			log.Fatal("Pause Update Error: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println(email)
		fmt.Println(level)
		fmt.Println(levelTime)
		new_level, _ := strconv.ParseInt(level, 10, 64)
		new_levelTime, _ := strconv.ParseInt(levelTime, 10, 64)
		change := bson.M{"$set": bson.M{"Paused": newGameState, "currentPausedTime": time.Now(), "currentLevel": new_level, "currentLevelTime": new_levelTime}}
		err := db.C("gameInfo").Update(update, change)
		if err != nil {
			log.Fatal("Pause Update Error: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

/*

	var result UserGame // used to store the usergame return data from mongo
	// Update where
	update := bson.M{"User": userID}
	db.C("gameInfo").Find(bson.M{"User": userID}).One(&result)



	if newGameState == false {
		// if false, set pause as false and calculate the cumulative paused time in nanoseconds.
		// value to be set along with pause bool
		accPauseTime := result.AccumulatedPausedDuration + (time.Now().Sub(result.CurrentPausedTime).Minutes())
		change := bson.M{"$set": bson.M{"Paused": newGameState, "AccumulatedPausedTime": (accPauseTime)}}
		err := db.C("gameInfo").Update(update, change)
		if err != nil {
			log.Fatal("Pause Update Error: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		// if pause is false then set pause as true and current pause time as the current time
		change := bson.M{"$set": bson.M{"Paused": newGameState, "CurrentPausedTime": time.Now(), "CurrentLevel": }}
		err := db.C("gameInfo").Update(update, change)
		if err != nil {
			log.Fatal("Pause Update Error: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
*/
