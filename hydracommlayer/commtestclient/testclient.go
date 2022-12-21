package main

import (
	"flag"
	"log"
	"strings"

	"github.com/pienaahj/hydra/hydracommlayer"
	"github.com/pienaahj/hydra/hydracommlayer/hydraproto"
)

func main() {
	op := flag.String("type", "", "Server (s) or client (c) ?")
	address := flag.String("addr", ":8000", "address? host:port ")
	flag.Parse()

	switch strings.ToUpper(*op) {
	case "S":
		runServer(*address)
	case "C":
		runClient(*address)
	}
}

func runServer(dest string) {
	c := hydracommlayer.NewConnection(hydracommlayer.Protobuf)
	recvChan, err := c.ListenAndDecode(dest)
	if err != nil {
		log.Fatal(err)
	}
	for mesg := range recvChan {
		log.Println("Received: ", mesg)
	}
}

func runClient(dest string) {
	c := hydracommlayer.NewConnection(hydracommlayer.Protobuf)
	ship := &hydraproto.Ship{
		Shipname:    "Hydra",
		CaptainName: "Jala",
		Crew: []*hydraproto.Ship_CrewMember{
			&hydraproto.Ship_CrewMember{1, "Kevin", 5, "Pilot"},
			&hydraproto.Ship_CrewMember{2, "Jade", 4, "Tech"},
			&hydraproto.Ship_CrewMember{3, "Wally", 3, "Engineer"},
		},
	}
	if err := c.EncodeAndSend(ship, dest); err != nil {
		log.Println("Error occured while sending message", err)
	}
}
