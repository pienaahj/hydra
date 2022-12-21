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
	//  this loop will onloy start if the data on the channel becomes available
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
			{
				Id:           1,
				Name:         "Kevin",
				SecClearance: 5,
				Position:     "Pilot",
			},
			{
				Id:           2,
				Name:         "Jade",
				SecClearance: 4,
				Position:     "Tech",
			},
			{
				Id:           3,
				Name:         "Wally",
				SecClearance: 3,
				Position:     "Engineer",
			},
		},
	}
	if err := c.EncodeAndSend(ship, dest); err != nil {
		log.Println("Error occured while sending message", err)
	}
}

/* to run
server
go run testclient.go -type s -addr :8484
client
go run testclient.go -type c -addr :8484

*/
