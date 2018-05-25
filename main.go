package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

var db *mgo.Database

type user struct {
	Email    string  `json: "email" bson: "email"`
	Username string  `json: "username" bson: "username"`
	Password string  `json: "pass" bson: "pass"`
	PHash    *[]byte `json:"-", omitempty`
}

type blindStructure struct {
	Name      string `json:"Name" bson:"Name"`
	AllLevels []row  `json:"AllLevels" bson:"AllLevels"`
}

type row struct {
	Level      int64 `json:"level" bson:"level"`
	SmallBlind int64 `json:"smallBlind" bson:"smallBlind"`
	BigBlind   int64 `json:"bigBlind" bson:"bigBlind"`
	Ante       int64 `json:"ante" bson:"ante"`
	Duration   int64 `json:"duration" bson:"duration"`
}

type UserGame struct {
	User                      string    `json:"User" bson:"User"`
	StartTime                 time.Time `json:"startTime" bson:"startTime"`
	Paused                    bool      `json:"paused" bson:"paused"`
	CurrentPausedStartTime    time.Time `json:"currentPausedStartTime" bson:"currentPausedStartTime"`
	CurrentLevelTime          float64   `json:"currentLevelTime" bson:"currentLevelTime"`
	CurrentLevel              int64     `jsons:"currentLevel" bson:"currentLevel"`
	BlindScheduleName         string    `jsons:"blindScheduleName" bson:"blindScheduleName"`
	Levels                    []row     `json:"levels" bson:"levels"`
	CurrentPausedTime         time.Time `json:"CurrentPausedTime" bson:"CurrentPausedTime"`
	AccumulatedPausedDuration float64   `json:"AccumulatedPausedTime" bson:"AccumulatedPausedTime"`
}

type CurrentGame struct {
	User                   string
	StartTime              time.Time
	Paused                 bool
	CurrentPausedStartTime time.Time
	CurrentLevelTime       float64
	CurrentLevel           int
	GameInfo               blindStructure
}

var tpl *template.Template
var router *mux.Router

func init() {
	tpl = template.Must(template.New("random").Delims("&&&", "&&&").ParseGlob("optui/public/*.html"))
}

func login(res http.ResponseWriter, _ *http.Request) {
	tpl.ExecuteTemplate(res, "index.html", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "timer.html", nil)
}

func main() {
	session, err := mgo.Dial("localhost") // connect to server
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}

	defer session.Close() // close the connection when main returns

	db = session.DB("game")

	router = mux.NewRouter()
	// Configure websocket route
	router.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages

	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/login", GetUser).Methods("POST")
	//router.HandleFunc("/create/{id}", games).Methods("POST")
	router.HandleFunc("/games/{id}", existing).Methods("GET")
	//testing button click action 2 cases pause and play
	router.HandleFunc("/games/{id}/pause/{level}/{levelTime}", func(w http.ResponseWriter, r *http.Request) {
		updateGamePauseState(w, true, r)
	}).Methods("PUT")
	router.HandleFunc("/games/{id}/play/{level}/{levelTime}", func(w http.ResponseWriter, r *http.Request) {
		updateGamePauseState(w, false, r)
	}).Methods("PUT")

	router.HandleFunc("/main", index).Methods("GET")
	router.HandleFunc("/", login).Methods("GET")

	router.PathPrefix("/css/").Handler(
		http.StripPrefix("/css", http.FileServer(http.Dir("./optui/public/css/"))))

	router.PathPrefix("/img/").Handler(
		http.StripPrefix("/img", http.FileServer(http.Dir("./optui/public/img/"))))

	router.PathPrefix("/js/").Handler(
		http.StripPrefix("/js", http.FileServer(http.Dir("./optui/public/js/"))))

	router.PathPrefix("/sounds/").Handler(
		http.StripPrefix("/sounds", http.FileServer(http.Dir("./optui/public/sounds/"))))
	// routine to control the messages in web socket
	go handleMessages()
	log.Fatal(http.ListenAndServe(":3000", router))

}
