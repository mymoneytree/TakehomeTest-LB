package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Router struct {
	Router *mux.Router
}

type Event struct {
	Encoded string `json:"encoded,omitempty"`
	Meta map[string]interface{} `json:"meta"`
	Decoded string `json:"decoded,omitempty"`
}

func main() {
	fmt.Println("Server listening on port: 5555")

	router := InitRouter()

	err := http.ListenAndServe(fmt.Sprintf(":%d", 5555), handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router.Router))
	if err != nil {
		fmt.Printf("Error from http server: %s\n", err.Error())
	}
}

func InitRouter() Router {
	// init router objs
	newRouter := mux.NewRouter()
	router := Router{Router: newRouter}

	// register base routes
	router.RegisterRoute("/", StatusHandler, "GET")

	router.RegisterRoute("/event", EventHandler, "POST")

	return router
}

func (router *Router) RegisterRoute(url string, f func(http.ResponseWriter, *http.Request), httpMethod string) error {
	route := router.Router.HandleFunc(url, f).Methods(httpMethod)
	return route.GetError()
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func EventHandler(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Sprintf("Error reading request body: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		errString := []byte(fmt.Sprintf("{\"error\": \"Error parsing request: %s\"}", err.Error()))
		w.Write(errString)
		return
	}


  var event Event

	if err := json.Unmarshal(payload, &event); err != nil {
		fmt.Printf("⚠️ Error unmarshaling JSON paylaod: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		errString := []byte(fmt.Sprintf("{\"error\": \"Error decoding JSON payload: %s\"}", err.Error()))
		w.Write(errString)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(event.Encoded)
	if err != nil {
		fmt.Printf("⚠️ Error decoding encoded string: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		errString := []byte(fmt.Sprintf("{\"error\": \"Error decoding encoded string: %s\"}", err.Error()))
		w.Write(errString)
		return
	}

	event.Decoded = string(decoded)
	event.Encoded = ""

	go DispatchEvent(event)

	w.WriteHeader(http.StatusOK)

}


func DispatchEvent(event Event) {
	waitTime := rand.Intn(10)
	time.Sleep(time.Duration(waitTime) * time.Second)
	eventBytes, _ := json.Marshal(event)
	
	_, err := http.Post("http://localhost:5556/event", "application/json",
        bytes.NewBuffer(eventBytes))
	if err != nil {
		fmt.Printf("{\"error\": \"Error sending decoded payload: %s\"}\n", err.Error())
	}
}