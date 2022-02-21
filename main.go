package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pienaahj/hydra/hlogger"
)

func main() {
	logger := hlogger.GetInstance()
	logger.Println("Starting Hydra web service...")


	// handle the root route
	http.HandleFunc("/", sroot)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func sroot(w http.ResponseWriter, r *http.Request) {
	logger := hlogger.GetInstance()
	fmt.Fprint(w, "Welcome to the Hydra software system")

	logger.Println("Received an http request on root url")
}