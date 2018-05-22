package main

import (
	"fmt"
	"errors"
    "strings"
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
	"encoding/json"
	"time"
	"strconv"
)

func LoadExistingGame(email string, collection *mgo.Collection) (string, error){
   if collection == nil {
   	  return "",errors.New("Collection is nil") 
   }
   	   
   if len(strings.TrimSpace(email)) == 0 {
      return "",errors.New("email is empty")
   }
  
    // retrieve user to see if they exist
    result := UserGame{}
    if collection.Find(bson.M{"UserID": email}).One(&result) == nil {
	    return "",errors.New("Could not find game")
    }   
   
   
   levelInfos := result.GameInfo.AllLevels
   var levelMinutes []int   
    
	for _, levelInfo := range levelInfos {
		levelMinutes = append(levelMinutes,levelInfo.Duration)
	}
 
   duration, _ := time.ParseDuration(strconv.FormatFloat(result.AccumulatedPausedDuration, 'g', 1, 64)+"s")
   
   level, levelTime := GetLevelAndLevelTime(result.StartTime, duration, result.Paused, result.CurrentPausedTime, levelMinutes)

   currGame := CurrentGame{
	   User:                      result.UserID,  
	   StartTime:                 result.StartTime,
	   Paused:                    result.Paused,
	   CurrentPausedStartTime:    result.CurrentPausedTime,
	   CurrentLevelTime:          levelTime,
	   CurrentLevel: 	         level,
	   GameInfo:                  result.GameInfo, 
   }

   resultString, err := json.Marshal(currGame)
   if err != nil {
	   fmt.Printf("Could not marshal to json")
   }
	   
   return string(resultString),nil
}
