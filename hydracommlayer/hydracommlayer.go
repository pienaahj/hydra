package hydracommlayer

import "Hydra/hydracommlayer/hydraproto"

// Communication messages
const (
	Protobuf uint8 = iota
)

// abstract the protocol to impliment any other protocol
type HydraConnection interface {
	EncodeAndSend(obj interface{}, destination string) error
	ListenAndDecode(listenaddress string) (chan interface{}, error) // impliment the channel generator pattern
	// so we can listen to this channel to intercept any new communication
}

// add factory pattern to generate the desired protocol
func NewConnection(connType uint8) HydraConnection {
	switch connType {
	case Protobuf:
		return hydraproto.NewProtoHandler() // call relevant constructor
	}
	return nil
}
