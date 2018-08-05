package res

import (
	"fmt"
	"net"
	"sync"
)

// Unix represents a unix domain socket.
type Unix struct {
	isOpen     bool
	isOpenLock *sync.Mutex

	path string
	conn net.Conn
}

// NewUnix allocates and initializes a unix instance.
func NewUnix(path *string) *Unix {
	return &Unix{
		isOpen:     false,
		isOpenLock: &sync.Mutex{},

		path: *path,
		conn: nil,
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
	tmp := fmt.Sprintf("%s:%s", TypeUnix, res.path)
	return &tmp
}

func (res *Unix) Read(b []byte) (n int, err error) {
	return res.conn.Read(b)
}

func (res *Unix) Write(b []byte) (n int, err error) {
	return res.conn.Write(b)
}
