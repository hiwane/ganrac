package cmd

import (
	"github.com/hiwane/ganrac"

	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type CmdParam struct {
	Verbose    int
	CadVerbose int
	Color      bool
	Quiet      bool
}

func (cp CmdParam) NewGanracLogger(cas, revision string) (*ganrac.Ganrac, *log.Logger) {
	if !cp.Quiet {
		if revision == "" {
			fmt.Printf("GaNRAC [cas=%s] see help();\n", cas)
		} else {
			fmt.Printf("GaNRAC [cas=%s, revision=%s] see help();\n", cas, revision)
		}
	}
	g := ganrac.NewGANRAC()
	logger := log.New(os.Stderr, "", log.LstdFlags)
	if cp.Quiet {
		logger.SetOutput(ioutil.Discard)
	}
	if cp.Color {
		ganrac.SetColordFml(true)
	}
	g.Eval(strings.NewReader(fmt.Sprintf("verbose(%d,%d);", cp.Verbose, cp.CadVerbose)))
	return g, logger
}

/*
 * 1文を取得.
 * 入力エラーリカバリが面倒だから１文ずつ処理する
 */
func get_line(in *bufio.Reader) (string, error) {
	//	line, err := in.ReadBytes(';')
	line := make([]rune, 0, 100)
	in_str := false  // 文字列内
	in_com := false  // コメント内
	depth_curly := 0 // 波括弧の深さ
	for {
		c, _, err := in.ReadRune()
		if err != nil {
			return "", err
		}
		line = append(line, c)
		if in_com {
			if c == '\n' {
				in_com = false
			}
			continue
		}
		if c == '"' {
			in_str = !in_str
		} else if in_str {
			//
		} else if c == '{' {
			depth_curly++
		} else if c == '}' && depth_curly > 0 {
			depth_curly--
		} else if c == ';' { // eol
			break
		} else if c == ':' && depth_curly <= 0 { // eolq
			break
		} else if c == '#' {
			// 改行まで skip
			in_com = true
		}
	}
	return string(line), nil
}

func Interpreter(g *ganrac.Ganrac) {
	in := bufio.NewReader(os.Stdin)
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
			continue
		}
		if p != nil {
			fmt.Println(p)
		}
	}
}
