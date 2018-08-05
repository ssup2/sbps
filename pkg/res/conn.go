package res

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Conn represents a connection from a client.
type Conn struct {
	isOpen     bool
	isOpenLock *sync.Mutex

	conn net.Conn
}

// NewConn allocates and initializes a conn instance.
func NewConn(conn *net.Conn) *Conn {
	return &Conn{
		isOpen:     true,
		isOpenLock: &sync.Mutex{},

		conn: *conn,
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
