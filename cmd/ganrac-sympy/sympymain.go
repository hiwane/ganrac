package main

import (
	"github.com/hiwane/ganrac/cas/para"
	"github.com/hiwane/ganrac/cas/sympy"
	"github.com/hiwane/ganrac/cmd"

	"flag"
	"fmt"
	"os"
)

var gitCommit string

func main() {
	var (
		// SageMath版固有オプション
		tmpfile = flag.String("tmpfile", "/tmp/ganrac.tmp", "temporary file")

		// 共通オプション
		cp cmd.CmdParam
	)

	cp.FlagVars()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	g, logger := cp.NewGanracLogger("SymPy", gitCommit)
	casSP, err := sympy.NewSymPy(g, *tmpfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "initialize sage failed: %s\n", err.Error())
		os.Exit(1)
	}
	g.SetParaNum(cp.ConcurrentNum)
	if cp.ConcurrentNum > 0 {
		paracas := para.NewParaCAS(casSP)
		defer paracas.Close()
		g.SetCAS(paracas)
	} else {
		defer casSP.Close()
		g.SetCAS(casSP)
	}

	logger.Printf("START!!!!")
	cp.Interpreter(g)
}
