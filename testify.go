package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/leonj1/testify/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ErrorResponse struct {
	ErrorMessage string `json:"err_msg"`
}

type Route struct{}

func (m *Route) getAllConfessionsHandler(w http.ResponseWriter, r *http.Request) {
	c := models.Confession{}
	confessions, err := c.FindAll()
	if err != nil {
		log.Printf("Problem fetching all confessions: %s\n", spew.Sdump(err))
		response := &ErrorResponse{ErrorMessage: "Problem getting all confessions"}
		respondWithJSON(w, 404, response)
		return
	}
	respondWithJSON(w, 200, confessions)
}

func (m *Route) getConfessionByNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	confessionName := vars["name"]
	c := models.Confession{}
	confession, err := c.FindByName(confessionName)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem getting confession by name"}
		respondWithJSON(w, 404, response)
		return
	}
	respondWithJSON(w, 200, confession)
}

func (m *Route) addConfessionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem reading body of request"}
		respondWithJSON(w, 404, response)
		return
	}
	log.Printf("Payload received %s\n", body)
	var confession models.Confession
	err = json.Unmarshal(body, &confession)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
		response := &ErrorResponse{ErrorMessage: "Problem serializing confession from request payload"}
		respondWithJSON(w, 404, response)
		return
	}
	log.Printf("Marshalled: %s\n", spew.Sdump(confession))
	savedConfession, err := confession.Save()
	if err != nil {
		if err.Error() == "host already exists" {
			respondWithTEXT(w, 403, err.Error())
		} else {
			respondWithTEXT(w, 500, err.Error())
		}
		return
	}
	saved, _ := confession.FindByName(savedConfession.Name)
	respondWithJSON(w, 201, saved)
}

func (m *Route) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithTEXT(w http.ResponseWriter, code int, response string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(response))
}

func main() {
	var userName = flag.String("user", "", "db username")
	var password = flag.String("pass", "", "db password")
	var databaseName = flag.String("db-name", "", "db name")
	var databaseHost = flag.String("db-host", "", "db host")
	var databasePort = flag.String("db-port", "", "db port")
	var serverPort = flag.String("http-port", "", "server port")
	flag.Parse()

	// open connection to db:
	// <username>:<pw>@tcp(<HOST>:<port>)/<dbname>
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", *userName, *password, *databaseHost, *databasePort, *databaseName)
	fmt.Printf("Connection string: %s\n", connectionString)
	models.InitDB(connectionString)

	log.SetOutput(os.Stdout)
	log.SetOutput(os.Stderr)

	handlers := &Route{}
	s := mux.NewRouter()
	s.HandleFunc("/confessions", handlers.addConfessionHandler).Methods("POST")
	s.HandleFunc("/confessions", handlers.getAllConfessionsHandler).Methods("GET")
	s.HandleFunc("/confessions/{name}", handlers.getConfessionByNameHandler).Methods("GET")

	// HealthCheck
	s.HandleFunc("/public/health", handlers.healthCheckHandler).Methods("GET")

	log.Printf("Staring HTTPS service on %s ...\n", *serverPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", *serverPort), s); err != nil {
		panic(err)
	}
}
