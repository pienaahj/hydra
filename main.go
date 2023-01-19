package main

import (
	"fmt"

	"github.com/pienaahj/hydra/hlogger"
	"github.com/pienaahj/hydra/hydraweb/hydraportal"
)

func main() {
	logger := hlogger.GetInstance()
	logger.Println("Starting Hydra web service...")
	fmt.Println("Started Hydralogger...")

	err := hydraportal.Run()
	if err != nil {
		fmt.Println("Failed to start", err.Error())
	}
}

/*
	// handle the root route
	http.HandleFunc("/", sroot)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func sroot(w http.ResponseWriter, r *http.Request) {
	logger := hlogger.GetInstance()
	fmt.Fprint(w, "Welcome to the Hydra software system")

	logger.Println("Received an http request on root url")
}
*/
