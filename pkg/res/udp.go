package res

import (
	"fmt"
	"net"
	"sync"
)

// UDP represents a UDP socket.
type UDP struct {
	conn net.Conn
	ip   string
	port int

	isOpenLock *sync.Mutex
	isOpen     bool

	mode byte
}

// NewUDP allocates and initializes a UDP instance.
func NewUDP(ip *string, port int, mode byte) *UDP {
	return &UDP{
		conn: nil,
		ip:   *ip,
		port: port,

		isOpenLock: &sync.Mutex{},
		isOpen:     false,

		mode: mode,
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

// IsOpen checks open of the resource.
func (res *UDP) IsOpen() bool {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	return res.isOpen
}

// IsRable checks resource is readable.
func (res *UDP) IsRable() bool {
	if res.mode&(1<<ModeR) == (1 << ModeR) {
		return true
	} else {
		return false
	}
}

// IsWable check resource is writeable
func (res *UDP) IsWable() bool {
	if res.mode&(1<<ModeW) == (1 << ModeW) {
		return true
	} else {
		return false
	}
}
