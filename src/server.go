package main

import (
	"fmt"
	"log"
	"net/http"
	"math/rand"
	"io"
	"os"
	"io/ioutil"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	path := "/hello"

	if r.URL.Path != path  {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w,"Hello again!")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err  != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful.\n")
	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Address = %s\n", address)
}

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
	fmt.Fprintf(w, "Healthy!")
}

func getRandomNumber(x int) int {
	return rand.Intn(x)
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

func Ping(w http.ResponseWriter, r *http.Request) {
	gambleNumber := getRandomNumber(100)
	hostname, err := os.Hostname()
	
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	io.WriteString(w, "Pong from " + hostname +"\n")
	if gambleNumber > 0 {
		pingNeighbor()
	}
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/ping", Ping)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}