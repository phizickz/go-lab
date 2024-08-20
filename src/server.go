package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/goombaio/namegenerator"
)

type config struct {
	killSwitch      int
	name            string
	pingCounter     int
	pingRepetitions int
}

var this config

func getRandomNumber() int {
	number := time.Now().Second() * 100
	return rand.Intn(number)
}

func getRandomName() string {
	return namegenerator.NewNameGenerator(time.Now().UTC().UnixNano()).Generate()
}

func pingNeighbor() {
	var deploymentName string
	if os.Getenv("DEPLOYMENT_NAME") == "" {
		deploymentName = "localhost:8080"
	} else {
		deploymentName = os.Getenv("DEPLOYMENT_NAME")
	}

	resp, err := http.Get(fmt.Sprintf("http://%v/ping", deploymentName))
	if err != nil {
		fmt.Println(err)
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	this.pingCounter++
	if this.pingCounter >= this.killSwitch {
		fmt.Printf("Killswitch limit reached! %v shutting down...", this.name)
		os.Exit(0)
	}

	for i := 0; i < this.pingRepetitions; i++ {
		go pingNeighbor()
	}
	// w.Write([]byte("PONG"))
}

func main() {
	this = config{
		getRandomNumber(),
		getRandomName(),
		0,
		2,
	}

	http.HandleFunc("GET /ping", pingHandler)

	fmt.Printf("Starting %v at port 8080.\nKillswitch is: %v\n", this.name, this.killSwitch)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
