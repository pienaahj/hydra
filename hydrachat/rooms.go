package hydrachat

import (
	"fmt"
	"net"
	"sync"
)

// Room creates a chat room struct
type room struct {
	name          string
	Msgch         chan string
	clients       map[chan<- string]struct{} // create a set datatype; when we send a message here it will go to a specific client
	Quit          chan struct{}
	*sync.RWMutex // to protext the clients map
}

// CreateRoom creates the chat room
func CreateRoom(name string) *room {
	r := &room{
		name:    name,
		Msgch:   make(chan string),
		RWMutex: new(sync.RWMutex),
		clients: make(map[chan<- string]struct{}),
		Quit:    make(chan struct{}),
	}
	r.Run() //  method to initialise a room to start waiting for clients to connect
	fmt.Println("Room created successfully")
	return r
}

// AddClient adds a client to the chat room
func (r *room) AddClient(c net.Conn) {
	logger.Println("Adding client", c.RemoteAddr())
	r.Lock()
	wc, done := StartClient(r.Msgch, c, r.Quit)
	r.clients[wc] = struct{}{}
	r.Unlock()

	// remove client when done is signalled
	go func() {
		<-done
		r.RemoveClient(wc)
	}()
}

// ClClient counts the clients in the chat room
func (r *room) ClCount() int {
	return len(r.clients)
}

// RemoveClient removes a client from the chat room clean resources
func (r *room) RemoveClient(wc chan<- string) {
	logger.Println("Removing client")
	r.Lock()
	close(wc)
	delete(r.clients, wc)
	r.Unlock()
	select {
	case <-r.Quit: // clean the chat room
		if len(r.clients) == 0 { // if no more clients
			close(r.Msgch) // close the room
		}
	default: // don't block
	}
}

// Run method starts a logger and spawns a gouroutine to handle the messages
func (r *room) Run() {
	logger.Println("Starting chat room", r.name)
	go func() { // not blocking the run method
		for msg := range r.Msgch { // will exit if the message channel is closed
			r.broadcastMsg(msg)
			fmt.Println("I'm still ranging over the message ch") // second stage of pipeline starts
		}
	}()
}

// broadcastMsgmethod  broadcasts the messages to all the room clients
// this represents second part of the pipeline second stage
func (r *room) broadcastMsg(msg string) {
	r.RLock()
	defer r.RUnlock()
	fmt.Println("Received message: ", msg)
	for wc := range r.clients { // wc for write only channels
		go func(wc chan<- string) { // handle all clients on their own goroutine - if a problem occures the other clients aren't affected
			wc <- msg // this will block until the channel is read
		}(wc) // send an argument to ensure a copy is sent - value might change while in the loop
	}
}
