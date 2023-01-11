package hydraweb

import (
	"fmt"
	"net/http"

	"github.com/pienaahj/hydra/hlogger"
)

func Run() {

	http.HandleFunc("/", sroot)
	http.Handle("/testhandle", newHandler())
	http.HandleFunc("/testquery", queryTestHandler)
	http.ListenAndServe(":8080", nil) // Defaultservmux is default
	// http.ListenAndServe(":8080", newHandler()) // use the newly created custom servmux - newHandler

	/*
		//  Create a custom http server with acustom handler
		server := &http.Server{
			Addr:         ":8080",
			Handler:      newHandler(),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		}
		server.ListenAndServe()
	*/
}

// queryTestHandler handles the query string provided in the URL
func queryTestHandler(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()
	// log.Println("Forms", r.Form)
	q := r.URL.Query()
	message := fmt.Sprintf("Query: %v\n", q)

	//  key1=2&key2=3
	v1, v2 := q.Get("key1"), q.Get("key2")
	if v1 == v2 {
		message = message + fmt.Sprintf("V1 and V2 are equal %s \n", v1)
	} else {
		message = message + fmt.Sprintf("V1 is equal %s, V2 is equal %s \n", v1, v2)
	}
	fmt.Fprint(w, message)
}

func sroot(w http.ResponseWriter, r *http.Request) {
	logger := hlogger.GetInstance()
	fmt.Fprint(w, "Welcome to the Hydra software system")
	logger.Println("Received an http Get request on the root url")
}
