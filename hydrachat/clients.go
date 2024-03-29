package hydrachat

import (
	"bufio"
	"net"
)

// client a struct to represent a chat client
type client struct {
	*bufio.Reader
	*bufio.Writer
	wc chan string
}

// StartClient starts the chat client arguments message channel, connection(network), quit signal
// This represents the first stage in the pipeline
func StartClient(msgCh chan<- string, cn net.Conn, quit chan struct{}) (chan<- string, <-chan struct{}) {
	c := new(client)
	c.Reader = bufio.NewReader(cn) //  obtain buffered read/writers
	c.Writer = bufio.NewWriter(cn)
	c.wc = make(chan string)
	done := make(chan struct{})

	// setup the reader
	go func() {
		scanner := bufio.NewScanner(c.Reader) // scan the client for a message,
		// when the connection terminates the reader is closed
		for scanner.Scan() {
			logger.Println(scanner.Text())
			msgCh <- scanner.Text() // place the input text from the tcp conection onto the message channel
		} // ends the pipeline first stage
		done <- struct{}{} // send quit signal - client stop typing
	}()

	//setup the writer
	c.writeMonitor()

	//  wait for the client to terminate
	go func() {
		select {
		case <-quit:
			cn.Close() // close the connection(tcp) channel and clean up resources this will lead to the client closure
		case <-done: // exit the goroutine
		}
	}()

	return c.wc, done
}

// writeMonitor  method monitors the writing to the channel and writes to the client
func (c *client) writeMonitor() {
	go func() {
		for s := range c.wc {
			logger.Println("Sending", s)
			c.WriteString(s + "\n") // represents last part of the second stage fo the pipeline
			c.Flush()               // completes the writing to the buffer
		}
	}()
}
