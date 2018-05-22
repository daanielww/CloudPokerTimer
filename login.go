package main

import (
	"gopkg.in/mgo.v2"
	"net/http"
	"encoding/json"
	"fmt"
)

type UserController struct {
	session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	u := user{}

	json.NewDecoder(r.Body).Decode(&u)

	// create bson ID
//	u.Id = bson.NewObjectId()

	// store the user in mongodb
	uc.session.DB("game").C("users").Insert(u)

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	fmt.Fprintf(w, "%s\n", uj)
}


