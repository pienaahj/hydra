package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// root handler
	http.HandleFunc("/", sroot)
	
	fmt.Println("Server started at port 8080...")

	// create server at port 8080 with default serve Mux
	log.Fatal(http.ListenAndServe(":8080", nil))
	

}

// sroot handles the root route
func sroot(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome to the Hydra software system")
}