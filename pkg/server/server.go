package server

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/ssup2/sbps/pkg/log"
	"github.com/ssup2/sbps/pkg/res"
)

// Server manages a server resource handler and a listen goroutine.
type Server struct {
	listener *Listener
	ticker   *time.Ticker

	mQuit chan struct{}
	lQuit chan struct{}

	sResHs       map[*res.Handler]struct{}
	sResClosedHs map[*res.Handler]struct{}
	sResHNoti    chan *res.Handler
	cResHs       map[*res.Handler]struct{}
	cResHNoti    chan *res.Handler
	resHLock     *sync.Mutex

	isRun     bool
	isRunLock *sync.Mutex
}

// New allocates and initialize a server instance.
func New(opt *string) (s *Server, err error) {
	log.Infof("Allocate a server")

	opts := strings.Split(*opt, ":")
	if len(opts) != 2 {
		return nil, errors.New("Wrong server options")
	}

	listener, err := NewListener(&opts[0], &opts[1])
	if err != nil {
		return nil, err
	}

	return &Server{
		listener: listener,
		ticker:   nil,

		mQuit: make(chan struct{}, 1),
		lQuit: make(chan struct{}, 1),

		sResHs:       make(map[*res.Handler]struct{}),
		sResClosedHs: make(map[*res.Handler]struct{}),
		sResHNoti:    make(chan *res.Handler, 1),
		cResHs:       make(map[*res.Handler]struct{}),
		cResHNoti:    make(chan *res.Handler, 1),
		resHLock:     &sync.Mutex{},

		isRun:     false,
		isRunLock: &sync.Mutex{},
	}, nil
}

// Close deinit and clean the server.
func (s *Server) Close() {
	// Stop goroutines
	s.isRunLock.Lock()
	if s.isRun == true {
		s.lQuit <- struct{}{}
		s.mQuit <- struct{}{}
		close(s.lQuit)
		close(s.mQuit)
	}
	s.isRun = false
	s.isRunLock.Unlock()

	// Deinit
	s.ticker.Stop()

	// Close server, client resource handlers
	s.resHLock.Lock()
	close(s.cResHNoti)
	for sResH := range s.sResHs {
		sResH.Stop()
		sResH.Close()
	}
	close(s.sResHNoti)
	for cResH := range s.cResHs {
		cResH.Stop()
		cResH.Close()
	}
	s.sResHs = nil
	s.cResHs = nil
	s.resHLock.Unlock()
}

// AddSResHandler append a server resource handler.
func (s *Server) AddSResHandler(sResH *res.Handler) {
	log.Infof("Add the server resource - %s", *sResH.GetRes().GetInfo())
	s.resHLock.Lock()
	defer s.resHLock.Unlock()

	_, exist := s.sResHs[sResH]
	if exist {
		return
	}
	s.sResHs[sResH] = struct{}{}

	// Set write target handler for each handlers.
	for cResH := range s.cResHs {
		sResH.AddWriteTarget(cResH)
		cResH.AddWriteTarget(sResH)
	}
}

// RemoveSResHandler remove the server resource handler.
func (s *Server) RemoveSResHandler(sResH *res.Handler) {
	log.Infof("Remove the server resource - %s", *sResH.GetRes().GetInfo())
	s.resHLock.Lock()
	defer s.resHLock.Unlock()

	_, exist := s.sResHs[sResH]
	if !exist {
		return
	}
	delete(s.sResHs, sResH)

	_, existClosed := s.sResClosedHs[sResH]
	if !existClosed {
		return
	}
	delete(s.sResClosedHs, sResH)
}

// AddSResClosedHandler append server resource handler to closed handler map.
func (s *Server) AddSResClosedHandler(sResH *res.Handler) {
	s.resHLock.Lock()
	defer s.resHLock.Unlock()

	_, exist := s.sResHs[sResH]
	if !exist {
		return
	}

	_, existClosed := s.sResClosedHs[sResH]
	if existClosed {
		return
	}
	s.sResClosedHs[sResH] = struct{}{}
}

// RemoveSResClosedHandler remove the server resource handler from closed handler map.
func (s *Server) RemoveSResClosedHandler(sResH *res.Handler) {
	s.resHLock.Lock()
	defer s.resHLock.Unlock()

	_, exist := s.sResClosedHs[sResH]
	if !exist {
		return
	}
	delete(s.sResClosedHs, sResH)
}

// AddCResHandler append a client resource handler.
func (s *Server) AddCResHandler(cResH *res.Handler) {
	log.Infof("Add the client resource")
	s.resHLock.Lock()
	defer s.resHLock.Unlock()

	_, exist := s.cResHs[cResH]
	if exist {
		return
	}
	s.cResHs[cResH] = struct{}{}

	// Set write target handler for each handlers.
	for sResH := range s.sResHs {
		sResH.AddWriteTarget(cResH)
		cResH.AddWriteTarget(sResH)
	}
}

// RemoveCResHandler remove the client resource handler.
func (s *Server) RemoveCResHandler(cResH *res.Handler) {
	log.Infof("Remove the client resource")
	s.resHLock.Lock()
	defer s.resHLock.Unlock()

	_, exist := s.cResHs[cResH]
	if !exist {
		return
	}
	delete(s.cResHs, cResH)
}

// ReopenSResH try to reopen for closed SResHs
func (s *Server) ReopenSResH() {
	s.resHLock.Lock()
	defer s.resHLock.Unlock()

	for sResH := range s.sResClosedHs {
		err := sResH.GetRes().Open()
		if err == nil || err == res.ErrALO {
			log.Infof("Reopen server resource - %s", *sResH.GetRes().GetInfo())

			// Set write target handler for each handlers.
			for cResH := range s.cResHs {
				sResH.AddWriteTarget(cResH)
				cResH.AddWriteTarget(sResH)
			}

			_, exist := s.sResClosedHs[sResH]
			if !exist {
				return
			}
			delete(s.sResClosedHs, sResH)

			sResH.Run()
		} else {
			log.Debugf("Reopen server failed - %s", err.Error())
		}
	}
}

// AcceptCResH accept clients to commuicate SResHs
func (s *Server) AcceptCResH() {
	conn, err := s.listener.l.Accept()
	if err != nil {
		if s.isRun == false {
			log.Infof("Accept client failed - Close listener")
		} else {
			log.Errorf("Accept client failed - %s", err.Error())
		}
		return
	}

	log.Infof("Accept the new client")
	cResH := res.NewHandler(res.New(res.TypeConn, conn), s.cResHNoti)
	s.AddCResHandler(cResH)
	cResH.Run()
}

// SResH noti channel
func (s *Server) GetSResHNoti() chan *res.Handler {
	return s.sResHNoti
}

// Run start a listen goroutine.
func (s *Server) Run() {
	log.Infof("Run the server")
	s.isRunLock.Lock()
	defer s.isRunLock.Unlock()

	if s.isRun == true {
		return
	}
	s.isRun = true
	s.ticker = time.NewTicker(time.Second * 2)

	// Main goroutine
	go func() {
		for {
			select {
			case <-s.mQuit:
				return

			case sResH := <-s.sResHNoti:
				s.AddSResClosedHandler(sResH)

			case cResH := <-s.cResHNoti:
				s.RemoveCResHandler(cResH)

			case <-s.ticker.C:
				s.ReopenSResH()
			}
		}
	}()

	// Listen goroutine
	go func() {
		for {
			select {
			case <-s.lQuit:
				return

			default:
				s.AcceptCResH()
			}
		}
	}()
}

// Stop stops the server.
func (s *Server) Stop() {
	log.Infof("Stop the server")
	s.isRunLock.Lock()
	defer s.isRunLock.Unlock()

	if s.isRun == false {
		return
	}

	s.mQuit <- struct{}{}
	s.lQuit <- struct{}{}
	s.isRun = false

	// Close for stop listen goroutine
	// TODO remove this
	s.listener.l.Close()
}
