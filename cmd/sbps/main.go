package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ssup2/sbps/pkg/log"
	"github.com/ssup2/sbps/pkg/res"
	"github.com/ssup2/sbps/pkg/server"
)

// SResOpt represents a pair of server resource option
type SResOpt struct {
	sResType string
	sResLoc  string
}

// CheckSRes checks server resource option
func CheckSRes(resOpt *string) *[]*SResOpt {
	var sRess []*SResOpt

	for _, sRes := range strings.Split(*resOpt, ",") {
		resSplit := strings.Split(sRes, ":")
		if len(resSplit) != 2 {
			log.Critf("Wrong server resource option - %s", sRes)
			os.Exit(1)
		}

		resType := resSplit[0]
		resLoc := resSplit[1]
		sRess = append(sRess, &SResOpt{sResType: resType, sResLoc: resLoc})
	}

	return &sRess
}

func main() {
	// Options
	optMode := flag.String("mode", server.TypeTCP+":6060", "sbps mode - [TCP. UDP, UNIX]:Opt")
	optSRes := flag.String("sres", "", "Server resource list - [TCP, UDP, UNIX, PIPE]:Opt...")
	optLogPath := flag.String("logpath", "sbsp_log", "Log path")
	optLogLevel := flag.String("loglevel", "INFO", "Log level - [DEBUG, INFO, WARN, ERROR, CRIT]")
	flag.Parse()

	if len(os.Args) < 2 {
		flag.PrintDefaults()
		return
	}

	// Logger
	logError := log.Init(optLogPath, optLogLevel)
	if logError != nil {
		log.Critf("Init file logger failed - %s", logError.Error())
		os.Exit(1)
	}
	defer log.Clean()

	// Server
	server, serverError := server.New(optMode)
	if serverError != nil {
		log.Critf("Allocation of a server failed - %s", serverError.Error())
		os.Exit(1)
	}
	defer server.Close()

	sRess := CheckSRes(optSRes)
	for _, sRes := range *sRess {
		r := res.New(sRes.sResType, sRes.sResLoc)
		if r == nil {
			log.Error("Allocation a server resource error - %s:%s", sRes.sResLoc, sRes.sResLoc)
			continue
		}

		h := res.NewHandler(r, server.GetSResHNoti())
		openError := r.Open()
		if openError != nil {
			log.Warnf("Open of a server resource error - %s", openError.Error())
			server.AddSResHandler(h)
			server.AddSResClosedHandler(h)
		} else {
			server.AddSResHandler(h)
			h.Run()
		}
	}
	server.Run()

	// Set signal and block main goroutine
	sigs := make(chan os.Signal)
	block := make(chan struct{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)
	go func() {
		sig := <-sigs
		log.Infof("Get signal - %s", sig)
		block <- struct{}{}
	}()

	log.Infof("Block main goroutine")
	<-block
	log.Infof("Unblock main goroutine and exit main")
	return
}
