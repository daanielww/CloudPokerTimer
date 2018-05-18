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
	gameElapsedTime := (currentTime - startTime)
	gameBlindStructureLevels := game.BlindStucture.AllLevels[]
	 
	if paused = true {
		accumulatePauseDuration += (currentTime - currentPausedStartTime)
	}

	gameDuration := gameElapsedTime - accumulatePauseDuration

	func main() {  
		for i := 1; i <= 10; i++ {
			fmt.Printf(" %d",i)
		}
	}

}
