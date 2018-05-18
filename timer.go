package main

import "time"

func newTimer() {

}

func calculateCurrentGameTime(game) {
	startTime := game.Start
	currentTime := time.Now()
	currentPauseStartTime := game.currentPauseState
	accumulatePauseDuration := game.accumulatePauseDuration
	paused := game.Paused
	currentGameDuration := (currentTime - startTime)
	currentLevel := (currentTime - startTime)

}
