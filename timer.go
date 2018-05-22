package main

import (
	"fmt"
	"strconv"
	"time"
)

//Runs asychronously as a goroutine. levelTimeRemaining and currentLevel are channels to return the values back to the caller.
//Currently requires four game related values in order to calculate levelTimeRemaining and currentLevel.
//Once the data format of the Game struct has been established, consider passing in the single struct instead.
func getLevelAndLevelTime(startTime time.Time, accumulatedPauseDuration time.Duration, paused bool,
	currentPauseStartTime time.Time, levelMinutes []int, levelTimeRemaining chan time.Duration, currentLevel chan int) {
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

	levelTimeRemaining <- gameDuration - accumulatedLevelTime
	currentLevel <- curLevel

}

func mainHenry() {
	now := time.Now()
	startTime := now.Add(-89 * time.Minute)
	levelMinutes := []int{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7}
	accumulatePauseDuration, _ := time.ParseDuration("1m30s")
	var currentPausedStartTime time.Time
	levelTimeRemaining := make(chan time.Duration)
	currentLevel := make(chan int)

	go getLevelAndLevelTime(startTime, accumulatePauseDuration, false, currentPausedStartTime, levelMinutes[:], levelTimeRemaining, currentLevel)

	ltr, cl := <-levelTimeRemaining, <-currentLevel

	fmt.Println("Timer should be", ltr, "Level is", cl)
}
