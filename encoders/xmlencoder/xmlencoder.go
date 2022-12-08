package main

import (
	"encoding/xml"
	"log"
	"os"
)

type CrewMemeber struct {
	XMLName           xml.Name `xml:"member"`
	ID                int      `xml:"id,omitempty"`
	Name              string   `xml:"name,attr"`
	SecurityClearance int      `xml:"clearance,attr"`
	AccessCodes       []string `xml:"codes>code"`
}

type ShipInfo struct {
	XMLName   xml.Name `xml:"ship"`
	ShipID    int      `xml:"shipInfo>ShipID"`
	ShipClass string   `xml:"shipInfo>ShipClass"`
	Captain   CrewMemeber
}

func main() {
	file, err := os.Create("xmlfile.xml")
	if err != nil {
		log.Fatal("Could not create file: ", err)
	}
	defer file.Close()
	// fill some data
	cm := CrewMemeber{Name: "Jaro", SecurityClearance: 10, AccessCodes: []string{"ADA", "LAL"}}
	si := ShipInfo{ShipID: 1, ShipClass: "Fighter", Captain: cm}
	enc := xml.NewEncoder(file)
	enc.Indent(" ", "	")
	enc.Encode(si)
	if err != nil {
		log.Fatal("Could not encode xml file", err)
	}
}
