package res

import (
	"fmt"
	"os"
	"sync"
)

// FIFO represents a FIFO (Named pipe).
type FIFO struct {
	fp   *os.File
	path string

	isOpenLock *sync.Mutex
	isOpen     bool

	mode byte
}

// NewFIFO allocates and initializes a FIFO instance.
func NewFIFO(path *string, mode byte) *FIFO {
	return &FIFO{
		fp:   nil,
		path: *path,

		isOpenLock: &sync.Mutex{},
		isOpen:     false,

		mode: mode,
	}
}

// Open connects to the FIFO socket.
func (res *FIFO) Open() error {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	if res.isOpen == true {
		return ErrALO
	}

	fp, errOpen := os.OpenFile(res.path, os.O_RDWR, os.ModeNamedPipe)
	if errOpen != nil {
		return errOpen
	}

	res.isOpen = true
	res.fp = fp
	return nil
}

// Close closes the FIFO domain socket.
func (res *FIFO) Close() error {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	if res.isOpen == false {
		return ErrALC
	}

	res.isOpen = false
	return res.fp.Close()
}

// GetInfo get tcp resource's info.
func (res *FIFO) GetInfo() *string {
	tmp := fmt.Sprintf("%s:%s", TypeFIFO, res.path)
	return &tmp
}

func (res *FIFO) Read(b []byte) (n int, err error) {
	return res.fp.Read(b)
}

func (res *FIFO) Write(b []byte) (n int, err error) {
	return res.fp.Write(b)
}

// IsOpen checks open of the resource.
func (res *FIFO) IsOpen() bool {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	return res.isOpen
}

// IsRable checks resource is readable.
func (res *FIFO) IsRable() bool {
	if res.mode&(1<<ModeR) == (1 << ModeR) {
		return true
	} else {
		return false
	}
}

// IsWable check resource is writeable
func (res *FIFO) IsWable() bool {
	if res.mode&(1<<ModeW) == (1 << ModeW) {
		return true
	} else {
		return false
	}
}
