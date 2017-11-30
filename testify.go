package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/leonj1/enchilada/models"
	"github.com/orcaman/concurrent-map"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
)

type PostResponse struct {
	Digest string `json:"digest"`
}

type QueryResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"err_msg"`
}

type MyConcurrentMap struct {
	cMap *cmap.ConcurrentMap
}

func (m *MyConcurrentMap) getAllHardwareHandler(w http.ResponseWriter, r *http.Request) {
	hw := models.Hardware{}
	hardware, err := hw.AllHardware()
	if err != nil {
		log.Printf("Problem fetching all hardware: %s\n", spew.Sdump(err))
		response := &ErrorResponse{ErrorMessage: "Problem getting all hardware"}
		respondWithJSON(w, 404, response)
		return
	}
	respondWithJSON(w, 200, hardware)
}

func (m *MyConcurrentMap) getHardwareHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host"]
	hw := models.Hardware{}
	hardware, err := hw.FindByHostName(host)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem getting host by hostname"}
		respondWithJSON(w, 404, response)
		return
	}
	respondWithJSON(w, 200, hardware)
}

func (m *MyConcurrentMap) getHardwareByServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["service"]
	s := models.Service{}
	hosts, err := s.FindHostsByServiceShortName(name)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem getting hosts by service"}
		respondWithJSON(w, 404, response)
		return
	}
	respondWithJSON(w, 200, hosts)
}

func (m *MyConcurrentMap) getHardwareByTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tagName := vars["tag"]
	t := models.Tag{}
	hosts, err := t.FindHostsByTag(tagName)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem getting hosts by tag"}
		respondWithJSON(w, 404, response)
		return
	}
	respondWithJSON(w, 200, hosts)
}

func (m *MyConcurrentMap) addHardwareHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem reading body of request"}
		respondWithJSON(w, 404, response)
		return
	}
	log.Printf("Payload received %s\n", body)
	var hardware models.Hardware
	err = json.Unmarshal(body, &hardware)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
		response := &ErrorResponse{ErrorMessage: "Problem serializing hardware from request payload"}
		respondWithJSON(w, 404, response)
		return
	}
	log.Printf("Marshalled: %s\n", spew.Sdump(hardware))
	savedHardware, err := hardware.Save()
	if err != nil {
		if err.Error() == "host already exists" {
			respondWithTEXT(w, 403, err.Error())
		} else {
			respondWithTEXT(w, 500, err.Error())
		}
		return
	}
	saved, _ := hardware.FindByHostName(savedHardware.Host)
	respondWithJSON(w, 201, saved)
}

func (m *MyConcurrentMap) updateServiceVersionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem reading request body"}
		respondWithJSON(w, 404, response)
		return
	}
	log.Printf("Request body received: %s\n", spew.Sdump(body))
	hw := models.Hardware{}
	hardware, err := hw.FindByHostName(host)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem fetching host"}
		respondWithJSON(w, 404, response)
		return
	} else if hardware == nil {
		response := &ErrorResponse{ErrorMessage: "Host does not exist"}
		respondWithJSON(w, 404, response)
		return
	}

	var f models.Service
	err = json.Unmarshal(body, &f)
	if err != nil {
		log.Printf("Problem serializing service from body: %s\n", spew.Sdump(err))
		response := &ErrorResponse{ErrorMessage: "Problem serializing service from request body"}
		respondWithJSON(w, 404, response)
		return
	}
	f.Host = host
	log.Printf("Unmarshalled Service from body: %s\n", spew.Sdump(f))

	if f.Version == "" || f.ShortName == "" {
		log.Printf("Nothing to update since version or shortname not provided\n")
		response := &ErrorResponse{ErrorMessage: "Nothing to update since version or shortname not provided"}
		respondWithJSON(w, 404, response)
		return
	}

	log.Printf("Before update service: %s\n", spew.Sdump(f))

	// fetch the current service by name and host, then update that object with whats provided
	actualService, err := f.FindByShortNameAndHost(f.ShortName, f.Host)
	if err != nil {
		log.Printf("Problem fetching actual service by shortname and host: %s\n", spew.Sdump(err))
		response := &ErrorResponse{ErrorMessage: "Problem fetching actual service"}
		respondWithJSON(w, 404, response)
		return
	}

	if actualService == nil {
		log.Printf("Service %s not found for host %s\n", f.ShortName, f.Host)
		response := &ErrorResponse{ErrorMessage: "How not found"}
		respondWithJSON(w, 403, response)
		return
	}

	log.Printf("Unmarshalled payload: %s\n", spew.Sdump(f))

	actualService.Version = f.Version

	_, err = actualService.Save()
	if err != nil {
		log.Printf("Problem saving service: %s\n", spew.Sdump(err))
		response := &ErrorResponse{ErrorMessage: "Problem saving service"}
		respondWithJSON(w, 404, response)
		return
	}
	hardware, _ = hw.FindByHostName(host)
	respondWithJSON(w, 200, hardware)
}

func (m *MyConcurrentMap) addServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem reading request body"}
		respondWithJSON(w, 404, response)
		return
	}
	log.Printf("Request body received: %s\n", spew.Sdump(body))
	hw := models.Hardware{}
	hardware, err := hw.FindByHostName(host)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem fetching host"}
		respondWithJSON(w, 404, response)
		return
	} else if hardware == nil {
		response := &ErrorResponse{ErrorMessage: "Host does not exist"}
		respondWithJSON(w, 404, response)
		return
	}

	var f models.Service
	err = json.Unmarshal(body, &f)
	if err != nil {
		log.Printf("Problem serializing service from body: %s\n", spew.Sdump(err))
		response := &ErrorResponse{ErrorMessage: "Problem serializing service from request body"}
		respondWithJSON(w, 404, response)
		return
	}
	f.Host = host
	log.Printf("Unmarshalled Service from body: %s\n", spew.Sdump(f))
	_, err = f.Save()
	if err != nil {
		log.Printf("Problem saving service: %s\n", spew.Sdump(err))
		response := &ErrorResponse{ErrorMessage: "Problem saving service"}
		respondWithJSON(w, 404, response)
		return
	}
	hardware, _ = hw.FindByHostName(host)
	respondWithJSON(w, 201, hardware)
}

func (m *MyConcurrentMap) addTagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := vars["host"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem reading request body"}
		respondWithJSON(w, 404, response)
		return
	}
	hw := models.Hardware{}
	hardware, err := hw.FindByHostName(host)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem fetching host"}
		respondWithJSON(w, 404, response)
		return
	} else if hardware == nil {
		response := &ErrorResponse{ErrorMessage: "Host does not exist"}
		respondWithJSON(w, 404, response)
		return
	}

	var f map[string]string
	err = json.Unmarshal(body, &f)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Problem serializing tag from request body"}
		respondWithJSON(w, 404, response)
		return
	}
	keys := reflect.ValueOf(f).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
		t := models.Tag{
			Key:   keys[i].String(),
			Value: f[keys[i].String()],
			Host:  host,
		}
		_, err = t.Save()
		if err != nil {
			response := &ErrorResponse{ErrorMessage: "Problem saving tag"}
			respondWithJSON(w, 404, response)
			return
		}
	}
	hardware, _ = hw.FindByHostName(host)
	respondWithJSON(w, 201, hardware)
}

func (m *MyConcurrentMap) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
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

	bar := cmap.New()
	handlers := &MyConcurrentMap{cMap: &bar}
	s := mux.NewRouter()
	s.HandleFunc("/hardware", handlers.addHardwareHandler).Methods("POST")
	s.HandleFunc("/hardware", handlers.getAllHardwareHandler).Methods("GET")

	// Tags
	s.HandleFunc("/hardware/{host}", handlers.getHardwareHandler).Methods("GET")
	s.HandleFunc("/tags/{host}", handlers.addTagHandler).Methods("POST")
	s.HandleFunc("/tags/{tag}", handlers.getHardwareByTagHandler).Methods("GET")

	// Services
	s.HandleFunc("/services/{host}", handlers.addServiceHandler).Methods("POST")
	s.HandleFunc("/services/{host}", handlers.updateServiceVersionHandler).Methods("PUT")
	s.HandleFunc("/services/{service}", handlers.getHardwareByServiceHandler).Methods("GET")

	// HealthCheck
	s.HandleFunc("/public/health", handlers.healthCheckHandler).Methods("GET")

	log.Printf("Staring HTTPS service on %s ...\n", *serverPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", *serverPort), s); err != nil {
		panic(err)
	}
}
