package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hiwane/ganrac"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var gitCommit string

func get_line(in *bufio.Reader) (string, error) {
	//	line, err := in.ReadBytes(';')
	line := make([]rune, 0, 100)
	in_str := false
	for {
		c, _, err := in.ReadRune()
		if err != nil {
			return "", err
		}
		line = append(line, c)
		if c == '"' {
			in_str = !in_str
		} else if c == ';' && !in_str {
			break
		}
	}
	return string(line), nil
}

func main() {
	var (
		cport   = flag.String("control", "localhost:1234", "ox-asir, control port")
		dport   = flag.String("data", "localhost:4321", "ox-asir, data port")
		ox      = flag.Bool("ox", false, "use ox-asir")
		verbose = flag.Bool("verbose", false, "verbose")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-ox][-data host:port][-control host:port]", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	in := bufio.NewReader(os.Stdin)
	fmt.Printf("GANRAC version %s. see help();\n", gitCommit)
	g := ganrac.NewGANRAC()
	logger := log.New(os.Stderr, "", log.LstdFlags)
	if *verbose {
		g.SetLogger(logger)
	}
	if *ox {
		logger.Printf("connect OX!!!!")
		connc, err := net.Dial("tcp", *cport)
		if err != nil {
			fmt.Fprintf(os.Stderr, "connect control [%s] failed: %s\n", *cport, err.Error())
			os.Exit(1)
		}
		defer connc.Close()

		time.Sleep(time.Second * 1)

		connd, err := net.Dial("tcp", *dport)
		if err != nil {
			fmt.Fprintf(os.Stderr, "connect data [%s] failed: %s\n", *dport, err.Error())
			os.Exit(1)
		}
		defer connd.Close()

		dw := bufio.NewWriter(connd)
		dr := bufio.NewReader(connd)
		cw := bufio.NewWriter(connc)
		cr := bufio.NewReader(connc)

		err = g.ConnectOX(cw, dw, cr, dr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "connect ox failed: %s", err.Error())
			os.Exit(1)
		}
	}

	logger.Printf("START!!!!")
	for {
		if _, err := os.Stdout.WriteString("> "); err != nil {
			fmt.Fprintf(os.Stderr, "WriteString: %s", err)
			break
		}
		line, err := get_line(in)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ReadBytes: %s", err)
			continue
		}

		p, err := g.Eval(strings.NewReader(string(line)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			break
		}
		if p != nil {
			fmt.Println(p)
		}
	}
}