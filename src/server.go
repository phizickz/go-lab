package main

import (
	"fmt"
	"log"
	"net/http"
	"math/rand"
	"io"
	"os"
	"io/ioutil"
	"time"
)

var pingCounter int
var killSwitch int

func healthHandler(w http.ResponseWriter, r *http.Request) {
	path := "/health"

	if r.URL.Path != path  {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	if killSwitch == 0 {
		http.Error(w, "Killswitch is not set correctly.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Healthy!")
}

func createSeedNumberFromTime() int {
	return time.Now().Nanosecond() / 10000
}

func getRandomNumber() int {
	return rand.Intn(createSeedNumberFromTime())
}

func incrementPingCounter() {
	pingCounter++
}

func pingNeighbor() {
	var deploymentName string
	if os.Getenv("DEPLOYMENT_NAME") == "" {
		deploymentName = "localhost:8080"
	} else {
		deploymentName = os.Getenv("DEPLOYMENT_NAME") 
	}
	fmt.Printf("Pinging neighbor.\n")
	resp, err := http.Get(fmt.Sprintf("http://%v/ping", deploymentName))
	if err != nil {
		fmt.Println(err)
	}

    body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(body))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	io.WriteString(w, "Pong from " + hostname +"\n")
	incrementPingCounter()

	if pingCounter >= killSwitch {
		os.Exit(1)
	}
	
	for i := 0; i >= 5; i++ {
		go pingNeighbor()
	}
}

func main() {
	pingCounter = 0
	killSwitch = getRandomNumber()

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	// http.HandleFunc("/", healthHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ping", pingHandler)

	fmt.Printf("Starting server at port 8080.\n")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}