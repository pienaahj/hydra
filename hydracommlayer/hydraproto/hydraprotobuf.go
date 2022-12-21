package hydraproto

import (
	"errors"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/protobuf/proto"
)

type ProtoHandler struct{}

// Constructor for the protobuf sender
func NewProtoHandler() *ProtoHandler {
	return new(ProtoHandler)
}

func (pSender *ProtoHandler) EncodeAndSend(obj interface{}, destination string) error {
	v, ok := obj.(*Ship)
	if !ok {
		return errors.New("proto: unknown message type")
	}
	data, err := proto.Marshal(v)
	if err != nil {
		return err
	}
	return sendmessage(data, destination)
}

func (pSender *ProtoHandler) DecodeProto(buffer []byte) (*Ship, error) {
	pb := new(Ship)
	return pb, proto.Unmarshal(buffer, pb)
}

// Take note that this method will return and empty channel until there is processing of the data
func (pSender *ProtoHandler) ListenAndDecode(listenaddress string) (chan interface{}, error) {
	outChan := make(chan interface{})
	l, err := net.Listen("tcp", listenaddress)
	if err != nil {
		return outChan, err
	}
	log.Println("listening to ", listenaddress)
	go func() { // accept the new connections and handle the listener
		defer l.Close()
		// this for loop will block the method forever, so use go routines to unblock the flow
		for {
			c, err := l.Accept()
			if err != nil {
				break
			}
			log.Println("Accepted connection from ", c.RemoteAddr())
			//  handle the connection separately for every connection otherwise the acceptance of new connections will wait
			go func(c net.Conn) { //  c changes for every connection NB! you need to divorce the caller from the goroutine
				//  otherwise the memory will be changed as well.
				defer c.Close()
				for {
					buffer, err := ioutil.ReadAll(c)
					if err != nil {
						break
					}
					if len(buffer) == 0 {
						continue //  no data in connection, wait for new connection to process
					}
					obj, err := pSender.DecodeProto(buffer)
					if err != nil {
						continue //  error in data decoding,  wait for new connection to process
					}
					//  use select not to wait if no-one is listening for this channel
					select {
					case outChan <- obj: // handle the data return
					default: // don't wait if nobody is listening for this channel
					}
					// or use a timeout to wait for 1 second
					// select {
					// 	case outChan <- obj: // handle the data return
					// 	<-time.After(1 * time.Second): // don't wait if nobody is listening for this channel
					// 	}
				}
			}(c)
		}
	}()
	return outChan, nil
}
func EncodeProto(obj interface{}) ([]byte, error) {
	if v, ok := obj.(*Ship); ok {
		return proto.Marshal(v)
	}
	return nil, errors.New("proto: unknown message type")
}

func sendmessage(buffer []byte, destination string) error {
	conn, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("Sending %d bytes to %s \n", len(buffer), destination)
	_, err = conn.Write(buffer)
	return err
}
