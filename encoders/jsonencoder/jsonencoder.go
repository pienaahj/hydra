package main

import (
	"encoding/json"
	"log"
	"os"
)

type CrewMemeber struct {
	ID                int      `json:"id,omitempty"`
	Name              string   `json:"name"`
	SecurityClearance int      `json:"security"`
	AccessCodes       []string `json:"accesscodes"`
}

type ShipInfo struct {
	ShipID    int
	ShipClass string
	Captain   CrewMemeber
}

func main() {
	f, err := os.Create("jfile.json")
	PrintFatalError(err)
	defer f.Close()

	// Create some data to encode
	cm := CrewMemeber{Name: "Jaro", SecurityClearance: 10, AccessCodes: []string{"ADA", "LAL"}}
	si := ShipInfo{1, "Fighter", cm}

	err = json.NewEncoder(f).Encode(&si)
	PrintFatalError(err)

}
func PrintFatalError(err error) {
	if err != nil {
		log.Fatal("Error happened while processing file: ", err)
	}
}
