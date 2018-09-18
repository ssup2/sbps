package res

import (
	"fmt"
	"net"
	"sync"
)

// Unix represents a unix domain socket.
type Unix struct {
	conn net.Conn
	path string

	isOpenLock *sync.Mutex
	isOpen     bool

	mode byte
}

// NewUnix allocates and initializes a unix instance.
func NewUnix(path *string, mode byte) *Unix {
	return &Unix{
		conn: nil,
		path: *path,

		isOpenLock: &sync.Mutex{},
		isOpen:     false,

		mode: mode,
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

// GetInfo get unix resource's info.
func (res *Unix) GetInfo() *string {
	tmp := fmt.Sprintf("%s:%s", TypeUnix, res.path)
	return &tmp
}

func (res *Unix) Read(b []byte) (n int, err error) {
	return res.conn.Read(b)
}

func (res *Unix) Write(b []byte) (n int, err error) {
	return res.conn.Write(b)
}

// IsOpen checks open of the resource.
func (res *Unix) IsOpen() bool {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	return res.isOpen
}

// IsRable checks resource is readable.
func (res *Unix) IsRable() bool {
	if res.mode&(1<<ModeR) == (1 << ModeR) {
		return true
	}
	return false
}

// IsWable check resource is writeable
func (res *Unix) IsWable() bool {
	if res.mode&(1<<ModeW) == (1 << ModeW) {
		return true
	}
	return false
}
