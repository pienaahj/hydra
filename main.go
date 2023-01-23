package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/pienaahj/hydra/hlogger"
	"github.com/pienaahj/hydra/hydrachat"
	"github.com/pienaahj/hydra/hydraweb/hydraportal"
)

func main() {
	logger := hlogger.GetInstance()
	logger.Println("Starting Hydra web service...")
	fmt.Println("Started Hydralogger...")
	//  -o w =>  hydra is running as a webserver
	//  -o c => hydra is running as a tcp chat server
	operation := flag.String("o", "w", "Operation: w for web \n c for chat")
	flag.Parse()
	switch strings.ToLower(*operation) {
	case "c":
		fmt.Println("running chat server...")
		exit := make(chan bool)
		err := hydrachat.Run(":2100")
		if err != nil {
			logger.Println("could not run hydra chat", err)
		}
		<-exit
	case "w":
		fmt.Println("running web server...")
		err := hydraportal.Run()
		if err != nil {
			logger.Println("Could not run hydra portal", err)
		}
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
