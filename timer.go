package main

import "time"

func newTimer() {

}

func getLevelAndLevelTime(game int) (int, int) {
	startTime := game.Start
	currentTime := time.Now()
	currentPauseStartTime := game.currentPauseState
	accumulatePauseDuration := game.accumulatePauseDuration
	paused := game.Paused
	gameElapsedTime := (currentTime - startTime)
	gameBlindStructureLevels := game.BlindStucture.AllLevels[]
	 
	if paused = true {
		accumulatePauseDuration += (currentTime - currentPausedStartTime)
	}

	gameDuration := gameElapsedTime - accumulatePauseDuration

	var currentLevel int

	for i := 0; i <= len(); i++ {
		currentLevel = i;
		accumulatedLevelTime += gameBlindStructureLevels[i]
		if accumulatedLevelTime > gameDuration {
			break
		}
	}

	return 

}
