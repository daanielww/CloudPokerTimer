package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	u := user{}

	json.NewDecoder(r.Body).Decode(&u)
	fmt.Println(u)

	if u.isEmpty() {
		fmt.Println("Error: User is empty")
		http.Error(w, "Error: please enter user info ", 404)
		return
	}

	_, err := findUser(u.Email)

	if err == nil {
		fmt.Println("Error: User already exists ", err)
		http.Error(w, "Error: User already exists ", 404)
		return
	}

	// store the user in mongodb
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)

	u.PHash = &hash

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

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

	u := user{}

	json.NewDecoder(r.Body).Decode(&u)

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
	userCheck, err := findUser(u.Email)

	// Fetch user
	if err != nil {
		fmt.Println("Error: user could not be found ", err)
		http.Error(w, "Error: user could not be found ", 404)
		return
	}

	err = bcrypt.CompareHashAndPassword(*userCheck.PHash, []byte(u.Password))

	if err != nil {
		fmt.Println("Error: password is not correct ", err)
		http.Error(w, "Error: password is not correct ", 404)
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

func (u user) isEmpty() bool {
	return u.Email == "" || u.Password == "" || u.Username == ""
}
