package res

import (
	"errors"
	"io"
	"sync"

	"github.com/ssup2/sbps/pkg/log"
)

// constants for handler
const (
	ReadBufSize      = 4096
	WriteChannelSize = 16
)

// ErrNR is error instance when res handler is not running
var ErrNR = errors.New("Handler is not running")

// WriteResult represents result of a write function.
type WriteResult struct {
	n   int
	err *error
}

// Handler manages goroutines to read from a resource or write to resource.
type Handler struct {
	res Res

	rQuit chan struct{}
	wQuit chan struct{}

	isRunLock   *sync.RWMutex
	wChanData   chan []byte
	wChanResult chan *WriteResult
	isRun       bool

	wTargetsLock *sync.Mutex
	wTargets     map[*Handler]struct{}

	closeNoti chan *Handler
}

// NewHandler allocates and initializes a handler instance.
func NewHandler(res Res, closeNoti chan *Handler) *Handler {
	return &Handler{
		res: res,

		rQuit: make(chan struct{}, 1),
		wQuit: make(chan struct{}, 1),

		isRunLock:   &sync.RWMutex{},
		wChanData:   make(chan []byte),
		wChanResult: make(chan *WriteResult),
		isRun:       false,

		wTargetsLock: &sync.Mutex{},
		wTargets:     make(map[*Handler]struct{}),

		closeNoti: closeNoti,
	}
}

// Close deinit and clean the handler.
func (h *Handler) Close() {
	// Stop goroutines and close write channel
	h.isRunLock.Lock()
	if h.isRun == true {
		h.rQuit <- struct{}{}
		h.wQuit <- struct{}{}
		close(h.rQuit)
		close(h.wQuit)

		close(h.wChanData)
		close(h.wChanResult)
	}
	h.isRun = false
	h.isRunLock.Unlock()

	// Clear write targets
	h.wTargetsLock.Lock()
	h.wTargets = nil
	h.wTargetsLock.Unlock()
}

// GetRes returns handler's resource
func (h *Handler) GetRes() Res {
	return h.res
}

// AddWriteTarget adds a write target handler.
func (h *Handler) AddWriteTarget(target *Handler) {
	log.Infof("Add the write target - %s", *h.res.GetInfo())

	h.wTargetsLock.Lock()
	defer h.wTargetsLock.Unlock()

	_, exist := h.wTargets[target]
	if exist {
		return
	}
	h.wTargets[target] = struct{}{}
}

// RemoveWriteTarget remove a write target handler.
func (h *Handler) RemoveWriteTarget(target *Handler) {
	log.Infof("Remove the write target - %s", *h.res.GetInfo())

	h.wTargetsLock.Lock()
	defer h.wTargetsLock.Unlock()

	_, exist := h.wTargets[target]
	if !exist {
		return
	}
	delete(h.wTargets, target)
}

// Write send data to write goroutine through the white channel.
func (h *Handler) Write(b []byte) (n int, err error) {
	h.isRunLock.RLock()
	defer h.isRunLock.RUnlock()

	if !h.isRun {
		return 0, ErrNR
	}

	h.wChanData <- b
	result := <-h.wChanResult
	return (*result).n, *(*result).err
}

// Run runs handler.
func (h *Handler) Run() {
	log.Infof("Run the res handler - %s", *h.res.GetInfo())
	h.isRunLock.Lock()
	defer h.isRunLock.Unlock()

	if h.isRun == true {
		return
	}
	h.isRun = true

	// Read goroutine
	go func() {
		for {
			select {
			case <-h.rQuit:
				log.Infof("Res handler - %s - read goroutine - close",
					*h.res.GetInfo())
				return

			default:
				// Read from resource
				b := make([]byte, ReadBufSize)
				n, err := h.res.Read(b)
				if err != nil {
					if err == io.EOF {
						// Resource (connection) is closed
						log.Infof("Res handler - %s - resource is closed",
							*h.res.GetInfo())
						h.res.Close()
						h.Stop()

						// Send close event
						if h.closeNoti != nil {
							h.closeNoti <- h
						}
					} else {
						// Error
						log.Errorf("Res handler - %s - read goroutine - "+
							"read from resource error - %s",
							*h.res.GetInfo(), err.Error())
					}
					continue
				}

				// Write to all write targets
				for target := range h.wTargets {
					_, err := target.Write(b[:n])
					if err != nil {
						if err == ErrNR {
							log.Infof("Res handler - %s - write target (%s) is closed",
								*h.res.GetInfo(), *target.res.GetInfo())
							h.RemoveWriteTarget(target)
						} else {
							log.Errorf("Res handler - %s - read goroutine - "+
								"write to write target error - %s",
								*h.res.GetInfo(), err.Error())
						}
					}
				}
			}
		}
	}()

	// Write goroutine
	go func() {
		for {
			select {
			case <-h.wQuit:
				log.Infof("Res handler - %s - write goroutine - close",
					*h.res.GetInfo())
				return

			case data := <-h.wChanData:
				// Write to resource
				n, err := h.res.Write(data)
				if err != nil {
					log.Errorf("Res handler - %s - write goroutine - "+
						"write to resource error - %s",
						*h.res.GetInfo(), err.Error())
				} else if n != len(data) {
					log.Errorf("Res handler - %s - write goroutine - "+
						"size of write is diff - request %d - result %d",
						*h.res.GetInfo(), len(data), n)
				}

				h.wChanResult <- &WriteResult{n: n, err: &err}
			}
		}
	}()
}

// Stop stops the handler.
func (h *Handler) Stop() {
	log.Infof("Stop the res handler - %s", *h.res.GetInfo())
	h.isRunLock.Lock()
	defer h.isRunLock.Unlock()

	if h.isRun == false {
		return
	}

	h.rQuit <- struct{}{}
	h.wQuit <- struct{}{}
	h.isRun = false
}
