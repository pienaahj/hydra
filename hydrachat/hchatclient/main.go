package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("Anonymous%d", rand.Intn(400))
	fmt.Println("Starting HydraChatClient...")
	fmt.Println("What is your name?")
	fmt.Scanln(&name)

	fmt.Printf("Hello %s, connecting to the hydra chat system.... \n", name)
	conn, err := net.Dial("tcp", "192.168.0.143:2300")
	if err != nil {
		log.Fatal("Could not connect to hydra chat system", err)
	}
	fmt.Println("Connected to the hydra chat system")
	name += ":"
	defer conn.Close()
	//  Listen on the connection for a message
	go func() { //  don't block the rest of the code
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	//  send a message to the chat server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() && err == nil {
		msg := scanner.Text()
		_, err = fmt.Fprintf(conn, name+msg+"\n")
	}

}
