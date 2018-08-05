package res

import (
	"fmt"
	"os"
	"sync"
)

// FIFO represents a FIFO (Named pipe).
type FIFO struct {
	isOpen     bool
	isOpenLock *sync.Mutex

	path string
	fp   *os.File
}

// NewFIFO allocates and initializes a FIFO instance.
func NewFIFO(path *string) *FIFO {
	return &FIFO{
		isOpen:     false,
		isOpenLock: &sync.Mutex{},

		path: *path,
		fp:   nil,
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

// IsOpen Checks open of the resource.
func (res *FIFO) IsOpen() bool {
	res.isOpenLock.Lock()
	defer res.isOpenLock.Unlock()

	return res.isOpen
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
