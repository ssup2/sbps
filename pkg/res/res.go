package res

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

// Resource type.
const (
	TypeTCP  = "TCP"
	TypeUDP  = "UDP"
	TypeUnix = "UNIX"
	TypeConn = "CONN"
	TypeFIFO = "FIFO"
)

// ErrType is error instance for wrong resource type.
var ErrType = errors.New("Wrong resource type")

// ErrInfo is error instance for wrong resource info.
var ErrInfo = errors.New("Wrong resource info")

// ErrALO is error instance when res is already opened.
var ErrALO = errors.New("Already opened")

// ErrALC is error instance when res is alread closed.
var ErrALC = errors.New("Already closed")

// Res defines resource read,write interface.
type Res interface {
	GetInfo() *string

	Open() error
	Close() error
	IsOpen() bool

	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

// New allocates and initializes a res instance
// depends on resource Type.
func New(rType *string, rInfo1 interface{}, rInfo2 interface{}) (Res, error) {
	switch *rType {
	case TypeTCP, TypeUDP:
		ip := net.ParseIP(*rInfo1.(*string))
		port, errTcp := strconv.Atoi(*rInfo2.(*string))
		if ip == nil || errTcp != nil || !(port >= 0 && port <= 65535) {
			return nil, ErrInfo
		}

		if strings.Compare(TypeTCP, *rType) == 0 {
			return NewTCP(rInfo1.(*string), port), nil
		} else {
			return NewUDP(rInfo1.(*string), port), nil
		}
	case TypeUnix:
		return NewUnix(rInfo1.(*string)), nil
	case TypeConn:
		return NewConn(rInfo1.(*net.Conn)), nil
	case TypeFIFO:
		return NewFIFO(rInfo1.(*string)), nil
	default:
	}

	return nil, ErrType
}

// CheckType checks resource type
func CheckType(rType string) bool {
	switch rType {
	case TypeTCP:
	case TypeUDP:
	case TypeUnix:
	case TypeConn:
	case TypeFIFO:
		return true
	default:
	}

	return false
}
