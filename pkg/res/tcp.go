package res

import (
	"fmt"
	"net"
	"sync"
)

// TCP represents a TCP socket.
type TCP struct {
	isOpen     bool
	isOpenLock *sync.Mutex

	ip   string
	port int
	conn net.Conn
}

// NewTCP allocates and initializes a TCP instance.
func NewTCP(ip *string, port int) *TCP {
	return &TCP{
		isOpen:     false,
		isOpenLock: &sync.Mutex{},

		ip:   *ip,
		port: port,
		conn: nil,
	}
}

// Open connects to the TCP socket.
func (res *TCP) Open() error {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	if res.isOpen == true {
		return ErrALO
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", res.ip, res.port))
	if err != nil {
		return err
	}

	res.isOpen = true
	res.conn = conn
	return nil
}

// Close closes the TCP domain socket.
func (res *TCP) Close() error {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	if res.isOpen == false {
		return ErrALC
	}

	res.isOpen = false
	return res.conn.Close()
}

// IsOpen Checks open of the resource.
func (res *TCP) IsOpen() bool {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	return res.isOpen
}

// GetInfo get tcp resource's info.
func (res *TCP) GetInfo() *string {
	tmp := fmt.Sprintf("%s:%s:%d", TypeTCP, res.ip, res.port)
	return &tmp
}

func (res *TCP) Read(b []byte) (n int, err error) {
	return res.conn.Read(b)
}

func (res *TCP) Write(b []byte) (n int, err error) {
	return res.conn.Write(b)
}
