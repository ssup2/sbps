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

// Listener represents listener infomation
type Listener struct {
	l net.Listener

	lType string
	lOpt  string
}

// NewListener allocates and initialize a listener instance
// depends on listener type.
func NewListener(lType *string, lOpt *string) (l *Listener, err error) {
	var listener net.Listener

	switch *lType {
	case TypeTCP:
		ln, err := net.Listen("tcp", fmt.Sprintf(":%s", *lOpt))
		if err != nil {
			return nil, err
		}
		listener = ln
		break
	case TypeUDP:
		break
	default:
		return nil, errors.New("Wrong listener type")
	}

	return &Listener{l: listener, lType: *lType, lOpt: *lOpt}, nil
}
