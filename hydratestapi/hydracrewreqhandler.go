package hydratestapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pienaahj/hydra/hydradblayer"
)

type HydraCrewReqHandler struct {
	dbConn hydradblayer.DBLayer
}

// new constructor for database interface struct
func NewhydraCrewReqHandler() *HydraCrewReqHandler {
	return new(HydraCrewReqHandler)
}

// connects to o - database type and conn - connection string
// point to the handler struct and populate the field inside to produce an active handler object
func (hcwreq *HydraCrewReqHandler) connect(o, conn string) error {
	dblayer, err := hydradblayer.ConnectDatabase(o, conn)
	if err != nil {
		return err
	}
	hcwreq.dbConn = dblayer
	return nil
}

//  there are two ways to impliment handlers for routes
//  first - create serveHttp handler
//  second create a handlerfunc

// incoming request that needs handling eg. http.HandleFunc("/hydracrew/", h.handleHydraCrewRequests)
func (hcwreq *HydraCrewReqHandler) handleHydraCrewRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ids := r.RequestURI[len("/hydracrew/"):] // url - /hydracrew/3
		id, err := strconv.Atoi(ids)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "id %s provided is not of valid number. \n", ids)
			return
		}
		cm, err := hcwreq.dbConn.FindMember(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error %s occured while searching for id %d \n", err.Error(), id)
			return
		}
		//  write the json encoded crewmember to the response
		json.NewEncoder(w).Encode(&cm)
	case "POST":
		cm := new(hydradblayer.CrewMember)
		err := json.NewDecoder(r.Body).Decode(cm)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error %s occored", err.Error())
			return
		}
		err = hcwreq.dbConn.AddMember(cm)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error %s occured while adding a crew member to the Hydra database", err.Error())
		}
		fmt.Fprintf(w, "Successfully inserted id %d \n", cm.ID)
	}
}
