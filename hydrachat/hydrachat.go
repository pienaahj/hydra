package hydrachat

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/pienaahj/hydra/hlogger"
)

var logger = hlogger.GetInstance()

// Run implements the tcp hydra chat server - this routine will block while the listening continues
func Run(connection string) error {
	fmt.Println("Listening on port ", connection)
	l, err := net.Listen("tcp", connection)
	if err != nil {
		logger.Println("Error connecting to chat client", err)
		return err
	}
	r := CreateRoom("HydraChat")
	fmt.Println("Chatroom created")
	go func(l net.Listener) {
		fmt.Println("invoked listener")
		for {
			fmt.Println("Waiting on connection")
			fmt.Printf("Listener %T listening\n", l)
			conn, err := l.Accept() // this blocks till a new connection or error
			if err != nil {
				logger.Println("Error accepting connection form client", err)
				break
			}
			go handleConnection(r, conn) // does not disturb the loop
		}
	}(l)
	go func() {
		//do not block the Run function excution and check if the microservice is shutting down
		//  Handle ctrl C
		fmt.Println("Checking interupt")
		ch := make(chan os.Signal, 1)
		// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM) // old way depricated
		signal.Notify(ch, os.Interrupt)

		s := <-ch //  block until the os signals exiting till os signal is received
		fmt.Println("got signal", s)
		//  then cleanup the resources
		l.Close() // close the connection
		fmt.Println("\nClosing tcp connection")
		close(r.Quit) // close the chatroom
		// ensure that all the message channels are closed before exiting
		if r.ClCount() > 0 {
			<-r.Msgch
		}
		os.Exit(0)
	}()
	return nil
}

// handleConnection logs the ip of the client and add it to the chat room
func handleConnection(r *room, c net.Conn) {
	logger.Println("Received request from client", c.RemoteAddr())
	r.AddClient(c)
}
