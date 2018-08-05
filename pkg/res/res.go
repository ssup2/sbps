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

	ModeR = 0
	ModeW = 1
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

	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)

	IsOpen() bool
	IsRable() bool
	IsWable() bool
}

// New allocates and initializes a res instance
// depends on resource Type. TypeConn is not supported.
func New(rType *string, rInfo []string) (Res, error) {
	mode := (byte)((1 << ModeR) | (1 << ModeW))

	switch *rType {
	case TypeTCP, TypeUDP:
		// Check rInfo
		ip := rInfo[0]
		port, err := strconv.Atoi(rInfo[1])
		if net.ParseIP(ip) == nil || err != nil ||
			!(port >= 0 && port <= 65535) {
			return nil, ErrInfo
		}

		if len(rInfo) >= 3 {
			tmpMode, err := MapMode(&rInfo[2])
			if err != nil {
				return nil, ErrInfo
			}

			mode = tmpMode
		}

		// Allocate a resource
		if strings.Compare(TypeTCP, *rType) == 0 {
			return NewTCP(&ip, port, mode), nil
		} else {
			return NewUDP(&ip, port, mode), nil
		}

	case TypeUnix, TypeFIFO:
		// Check rInfo
		path := rInfo[0]
		if path[0] != '/' && path[0] != '.' {
			return nil, ErrInfo
		}

		if len(rInfo) >= 2 {
			tmpMode, err := MapMode(&rInfo[1])
			if err != nil {
				return nil, ErrInfo
			}
			mode = tmpMode
		}

		// Allocate a resource
		if strings.Compare(TypeUnix, *rType) == 0 {
			return NewUnix(&path, mode), nil
		} else {
			return NewFIFO(&path, mode), nil
		}

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

// MapMode mapping a mode option to a byte.
func MapMode(rType *string) (byte, error) {
	if strings.Compare("RW", *rType) == 0 ||
		strings.Compare("WR", *rType) == 0 {
		return (byte)((1 << ModeR) | (1 << ModeW)), nil
	} else if strings.Compare("R", *rType) == 0 {
		return (byte)(1 << ModeR), nil
	} else if strings.Compare("W", *rType) == 0 {
		return (byte)(1 << ModeW), nil
	} else {
		return 0, errors.New("Wrong mode")
	}
}
