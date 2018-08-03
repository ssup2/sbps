package res

import (
	"fmt"
	"net"
	"sync"
)

// UDP represents a UDP socket.
type UDP struct {
	isOpen     bool
	isOpenLock *sync.Mutex

	ip   string
	port int
	conn net.Conn
}

// NewUDP allocates and initializes a UDP instance.
func NewUDP(ip *string, port int) *UDP {
	return &UDP{
		isOpen:     false,
		isOpenLock: &sync.Mutex{},

		ip:   *ip,
		port: port,
		conn: nil,
	}
}

// Open connects to the UDP socket.
func (res *UDP) Open() error {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	if res.isOpen == true {
		return ErrALO
	}

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", res.ip, res.port))
	if err != nil {
		return err
	}

	res.isOpen = true
	res.conn = conn
	return nil
}

// Close closes the UDP domain socket.
func (res *UDP) Close() error {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	if res.isOpen == false {
		return ErrALC
	}

	res.isOpen = false
	return res.conn.Close()
}

// IsOpen Checks open of the resource.
func (res *UDP) IsOpen() bool {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	return res.isOpen
}

// GetInfo get udp resource's info.
func (res *UDP) GetInfo() *string {
	tmp := fmt.Sprintf("%s:%s:%d", TypeUDP, res.ip, res.port)
	return &tmp
}

func (res *UDP) Read(b []byte) (n int, err error) {
	return res.conn.Read(b)
}

func (res *UDP) Write(b []byte) (n int, err error) {
	return res.conn.Write(b)
}
