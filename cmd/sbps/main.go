package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ssup2/sbps/pkg/log"
	"github.com/ssup2/sbps/pkg/res"
	"github.com/ssup2/sbps/pkg/server"
)

// Version, Build info for sbps and cmake
var (
	Version string
	Build   string
)

// SResOpt represents a pair of server resource option
type SResOpt struct {
	sResType string
	sResInfo []string
}

// SplitSRes splits server resource option
func SplitSRes(optSResLoc *string) *[]*SResOpt {
	var sRess []*SResOpt

	for _, sRes := range strings.Split(*optSResLoc, ",") {
		rSplit := strings.Split(sRes, ":")
		length := len(rSplit)

		if length < 2 || length > 4 {
			log.Critf("Wrong server resource option - %s", sRes)
			os.Exit(1)
		}

		rType := rSplit[0]
		rInfo := append(rSplit[:0], rSplit[1:]...)
		sRess = append(sRess, &SResOpt{sResType: rType, sResInfo: rInfo})
	}

	return &sRess
}

func main() {
	// Options
	optVersion := flag.Bool("v", false,
		"Print version")
	optMode := flag.String("mode", server.TypeTCP+":6060",
		"sbps proxy server mode (option TCP:port, UNIX:path)")
	optSResLoc := flag.String("resource", "",
		"Server resources (option TCP:ip:port[:RW], UDP:ip:port[:RW], UNIX:path[:RW], FIFO:path[:RW])")
	optSResInter := flag.Int("interval", 2,
		"Seconds of retry interval for closed server resources")
	optLogPath := flag.String("logpath", "./sbps.log",
		"Log path")
	optLogLevel := flag.String("loglevel", "INFO",
		"Log level (option DEBUG, INFO, WARN, ERROR, CRIT)")
	flag.Parse()

	if *optVersion {
		fmt.Printf("sbps version %s, build %s\n", Version, Build)
		return
	}

	if (len(os.Args) < 2) || strings.Compare(*optSResLoc, "") == 0 {
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
	server, serverError := server.New(optMode, *optSResInter)
	if serverError != nil {
		log.Critf("Allocation of a server failed - %s", serverError.Error())
		os.Exit(1)
	}
	defer server.Close()

	sRess := SplitSRes(optSResLoc)
	for _, sRes := range *sRess {
		r, resError := res.New(&sRes.sResType, sRes.sResInfo)
		if resError != nil {
			resStr := sRes.sResType
			for _, info := range sRes.sResInfo {
				resStr += (":" + info)
			}

			log.Critf("Allocation a server resource (%s) error - %s",
				resStr, resError.Error())
			os.Exit(1)
		}

		h := res.NewHandler(r, server.GetSResHNoti())
		openError := r.Open()
		if openError != nil {
			log.Warnf("Open of a server resource error - %s", openError.Error())

			if *optSResInter > 0 {
				server.AddSResHandler(h)
				server.AddSResClosedHandler(h)
			}
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
