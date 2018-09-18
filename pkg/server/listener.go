package server

import (
	"errors"
	"fmt"
	"net"
)

// Listener types
const (
	TypeTCP  = "TCP"
	TypeUDP  = "UDP"
	TypeUnix = "UNIX"
)

// Listener represents listener information
type Listener struct {
	ln net.Listener

	lType string
	lOpt  string
}

// NewListener allocates and initialize a listener instance
// depends on listener type.
func NewListener(lType *string, lOpt *string) (*Listener, error) {
	var ln net.Listener
	var err error

	switch *lType {
	case TypeTCP:
		ln, err = net.Listen("tcp", fmt.Sprintf(":%s", *lOpt))
	case TypeUnix:
		ln, err = net.Listen("unix", *lOpt)
	default:
		return nil, errors.New("Wrong listener type")
	}

	if err != nil {
		return nil, err
	}

	return &Listener{ln: ln, lType: *lType, lOpt: *lOpt}, nil
}
