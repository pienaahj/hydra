package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
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

func runClient(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("What message would you like to send?")
	for scanner.Scan() {
		fmt.Println("Writing ", scanner.Text())
		conn.Write(append(scanner.Bytes(), '\r')) // '\r' represents a return value

		fmt.Println("What message would you like to send?")
		buffer := make([]byte, 1024)
		// conn.SetReadDeadline(time.Now.Add(5 * time.second))
		_, err := conn.Read(buffer)

		if err != nil && err != io.EOF {
			log.Fatal(err)
		} else if err == io.EOF {
			log.Println("Connection is closed.")
		}
		fmt.Println(string(buffer))
	}
	return scanner.Err()
}

func runServer(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.Println("Listening....")
	defer l.Close()
	for {
		c, err := l.Accept() // this is a blocking call
		if err != nil {
			return err
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)
	writer := bufio.NewWriter(c)
	for {
		//buffer := make([]byte, 1024)
		// this is to protect against the connection getting stuck indefinitely - stale connections
		// you can also set this read or write only
		c.SetDeadline(time.Now().Add(5 * time.Second))
		line, err := reader.ReadString('\r') // this char is used to terminate a message
		//_, r := c.Read(buffer)
		if err != nil && err != io.EOF {
			log.Println(err)
			return
		} else if err == io.EOF {
			log.Println("Connection closed")
			return
		}
		fmt.Printf("Received %s from address %s \n", line[:len(line)-1], c.RemoteAddr())
		writer.WriteString("Message received...")
		writer.Flush()
	}
}

/* to run
server: go run tcpcommunications.go -type s -addr :8787
client: go run tcpcommunications.go -type c -addr :8787
*/
