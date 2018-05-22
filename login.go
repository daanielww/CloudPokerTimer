package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	u := user{}

	json.NewDecoder(r.Body).Decode(&u)
	fmt.Println(u)

	// create bson ID
	//	u.Id = bson.NewObjectId()

	_, err := findUser(u.Email)

	if err == nil {
		fmt.Println("Error: User already exists ", err)
		http.Error(w, "Error: User already exists ", 404)
		return
	}

	// store the user in mongodb
	db.C("users").Insert(u)

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(uj)
	w.WriteHeader(http.StatusCreated) // 201
	fmt.Fprintf(w, "%s\n", uj)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Grab name
	r.ParseForm()
	email := r.Form.Get("email")

	/*
		//necessary??
		// Verify name is ObjectId hex representation, otherwise return status not found
		if !bson.IsObjectIdHex(name) {
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}

		// ObjectIdHex returns an ObjectId from the provided hex representation.
		oid := bson.ObjectIdHex(name)
	*/

	// composite literal
	u, err := findUser(email)

	// Fetch user
	if err != nil {
		fmt.Println("Error: user could not be found ", err)
		http.Error(w, "Error: user could not be found ", 404)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(uj)
	w.WriteHeader(http.StatusOK) // 200
	//redirect to game page
	http.Redirect(w, r, "/game", http.StatusPermanentRedirect)

	fmt.Fprintf(w, "%s\n", uj)
}

func findUser(email string) (user, error) {

	u := user{}

	// Fetch user
	err := db.C("users").Find(bson.M{"email": email}).One(&u)

	return u, err

}
