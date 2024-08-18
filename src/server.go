package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/justinas/alice"
)

var pingCounter int
var killSwitch int

func getRandomNumber() int {
	return rand.Intn(time.Now().Nanosecond() / 10000)
}

func lifecycleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pingCounter++
		if pingCounter >= killSwitch {
			fmt.Println("Pingcount limit reached!")
			os.Exit(0)
		}
		next.ServeHTTP(w, r)
	})
}

func pingNeighborMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

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
	})
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG"))
}

func main() {
	pingCounter = 0
	killSwitch = getRandomNumber()
	fmt.Println("Killswitch is: " + fmt.Sprint(killSwitch))
	mux := http.NewServeMux()

	pingHandlerFunc := http.HandlerFunc(pingHandler)
	mwChain := alice.New(lifecycleMiddleware, pingNeighborMiddleware)

	mux.Handle("/ping", mwChain.Then(pingHandlerFunc))

	fmt.Printf("Starting server at port 8080.\n")
	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		log.Fatal(err)
	}
}
