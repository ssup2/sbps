package res

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Resource type.
const (
	TypeTCP  = "TCP"
	TypeUDP  = "UDP"
	TypeUnix = "UNIX"
	TypeConn = "CONN"
)

// ErrType is error instance for wrong resource type.
var ErrType = errors.New("Wrong resource type")

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
func New(rType string, rOpt interface{}) Res {
	switch rType {
	case TypeTCP:
	case TypeUDP:
	case TypeUnix:
		return NewUnix(rOpt.(string))
	case TypeConn:
		return NewConn(rOpt.(net.Conn))
	default:
	}

	return nil
}

// CheckType checks resource type
func CheckType(rType string) bool {
	switch rType {
	case TypeTCP:
	case TypeUDP:
	case TypeUnix:
	case TypeConn:
		return true
	default:
	}

	return false
}

// Conn represents a connection from a client.
type Conn struct {
	isOpen     bool
	isOpenLock *sync.Mutex

	conn net.Conn
}

// NewConn allocates and initializes a conn instance.
func NewConn(conn net.Conn) *Conn {
	return &Conn{
		isOpen:     true,
		isOpenLock: &sync.Mutex{},

		conn: conn,
	}
}

// Open opens the resource. Do not use this
func (res *Conn) Open() error {
	return errors.New("Do not support Open()")
}

// Close closes the connection.
func (res *Conn) Close() error {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	if res.isOpen == false {
		return ErrALC
	}

	res.isOpen = false
	return res.conn.Close()
}

// IsOpen Checks open of the resource.
func (res *Conn) IsOpen() bool {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	return res.isOpen
}

// GetInfo get conn resource's info
func (res *Conn) GetInfo() *string {
	var tmp string
	p := res.conn.RemoteAddr().Network()

	if strings.Contains(p, "tcp") {
		addr := res.conn.RemoteAddr().(*net.TCPAddr)
		tmp = fmt.Sprintf("%s:%s:%s:%d", TypeConn, "TCP", addr.IP.String(), addr.Port)
	} else if strings.Contains(p, "udp") {
		addr := res.conn.RemoteAddr().(*net.UDPAddr)
		tmp = fmt.Sprintf("%s:%s:%s:%d", TypeConn, "UDP", addr.IP.String(), addr.Port)
	} else if strings.Contains(p, "unix") {
		addr := res.conn.RemoteAddr().(*net.UnixAddr)
		tmp = fmt.Sprintf("%s:%s:%s", TypeConn, "UNIX", addr.String())
	} else {
		return nil
	}

	return &tmp
}

func (res *Conn) Read(b []byte) (n int, err error) {
	return res.conn.Read(b)
}

func (res *Conn) Write(b []byte) (n int, err error) {
	return res.conn.Write(b)
}

// Unix represents a unix domain socket.
type Unix struct {
	isOpen     bool
	isOpenLock *sync.Mutex

	path string
	conn net.Conn
}

// NewUnix allocates and initializes a unix instance.
func NewUnix(path string) *Unix {
	return &Unix{
		isOpen:     false,
		isOpenLock: &sync.Mutex{},

		conn: nil,
		path: path,
	}
}

// Open connects to the unix socket.
func (res *Unix) Open() error {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	if res.isOpen == true {
		return ErrALO
	}

	conn, err := net.Dial("unix", res.path)
	if err != nil {
		return err
	}

	res.isOpen = true
	res.conn = conn
	return nil
}

// Close closes the unix domain socket.
func (res *Unix) Close() error {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	if res.isOpen == false {
		return ErrALC
	}

	res.isOpen = false
	return res.conn.Close()
}

// IsOpen Checks open of the resource.
func (res *Unix) IsOpen() bool {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	return res.isOpen
}

// GetInfo get unix resource's info.
func (res *Unix) GetInfo() *string {
	tmp := fmt.Sprintf("%s:%s", TypeUDP, res.path)
	return &tmp
}

func (res *Unix) Read(b []byte) (n int, err error) {
	return res.conn.Read(b)
}

func (res *Unix) Write(b []byte) (n int, err error) {
	return res.conn.Write(b)
}
