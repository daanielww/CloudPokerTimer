package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

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

// function to trigger pause and play action
func updateGamePauseState(newGameState bool, userID string) {
	var result UserGame // used to store the usergame return data from mongo
	// Update where
	update := bson.M{"UserID": userID}
	db.C("gameInfo").Find(bson.M{"UserID": userID}).One(&result)
	if newGameState == false {
		// if false, set pause as false and calculate the cumulative paused time in nanoseconds.
		// value to be set along with pause bool
		accPauseTime := result.AccumulatedPausedDuration + (time.Now().Sub(result.CurrentPausedTime).Minutes())
		change := bson.M{"$set": bson.M{"Paused": newGameState, "AccumulatedPausedTime": (accPauseTime)}}
		err := db.C("gameInfo").Update(update, change)
		if err != nil {
			log.Fatal("Pause Update Error: ", err)
		}
	} else {
		// if pause is false then set pause as true and current pause time as the current time
		change := bson.M{"$set": bson.M{"Paused": newGameState, "CurrentPausedTime": time.Now()}}
		err := db.C("gameInfo").Update(update, change)
		if err != nil {
			log.Fatal("Pause Update Error: ", err)
		}
	}
}
