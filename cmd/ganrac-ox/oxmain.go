package main

import (
	openxm "github.com/hiwane/ganrac/cas/ox"
	"github.com/hiwane/ganrac/cmd"

	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var gitCommit string

func main() {
	var (
		cport      = flag.String("control", "localhost:1234", "ox-asir, control port")
		dport      = flag.String("data", "localhost:4321", "ox-asir, data port")
		ox         = flag.Bool("ox", false, "use ox-asir")
		ox_verbose = flag.Bool("ox_verbose", false, "ox_verbose")

		cp cmd.CmdParam
	)

	flag.IntVar(&cp.Verbose, "verbose", 0, "verbose")
	flag.IntVar(&cp.CadVerbose, "cad_verbose", 0, "cad_verbose")
	flag.BoolVar(&cp.Color, "color", false, "colored")
	flag.BoolVar(&cp.Quiet, "q", false, "quiet mode")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-ox][-data host:port][-control host:port]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	g, logger := cp.NewGanracLogger("OX", gitCommit)
	if *ox {
		logger.Printf("connect OX!!!!")
		connc, err := net.Dial("tcp", *cport)
		if err != nil {
			fmt.Fprintf(os.Stderr, "connect control [%s] failed: %s\n", *cport, err.Error())
			os.Exit(1)
		}

		time.Sleep(time.Second * 1)

		connd, err := net.Dial("tcp", *dport)
		if err != nil {
			connc.Close()
			fmt.Fprintf(os.Stderr, "connect data [%s] failed: %s\n", *dport, err.Error())
			os.Exit(1)
		}

		ox, err := openxm.NewOpenXM(connc, connd, g.Logger())
		if err != nil {
			fmt.Fprintf(os.Stderr, "connect ox failed: %s", err.Error())
			os.Exit(1)
		}

		defer ox.Close()
		g.SetCAS(ox)
		if *ox_verbose {
			ox.SetLogger(logger)
		}
	}

	logger.Printf("START!!!!")
	cmd.Interpreter(g)
}
