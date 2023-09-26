package main

import (
	"github.com/hiwane/ganrac/cas/para"
	"github.com/hiwane/ganrac/cas/sage"
	"github.com/hiwane/ganrac/cmd"

	"flag"
	"fmt"
	"os"
)

var gitCommit string

func main() {
	var (
		tmpfile = flag.String("tmpfile", "/tmp/ganrac.tmp", "temporary file")
		paran   = flag.Int("para", 0, "number of concurrent processes")

		cp cmd.CmdParam
	)

	flag.IntVar(&cp.Verbose, "verbose", 0, "verbose")
	flag.IntVar(&cp.CadVerbose, "cad-verbose", 0, "cad verbose")
	flag.BoolVar(&cp.Color, "color", false, "colored")
	flag.BoolVar(&cp.Quiet, "q", false, "quiet mode")
	flag.StringVar(&cp.CmdHistory, "history", "", "command history")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	g, logger := cp.NewGanracLogger("SageMath", gitCommit)
	Sage, err := sage.NewSage(g, *tmpfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "initialize sage failed: %s\n", err.Error())
		os.Exit(1)
	}
	g.SetParaNum(*paran)
	if *paran > 0 {
		paracas := para.NewParaCAS(Sage)
		defer paracas.Close()
		g.SetCAS(paracas)
	} else {
		defer Sage.Close()
		g.SetCAS(Sage)
	}

	logger.Printf("START!!!!")
	cp.Interpreter(g)
}
