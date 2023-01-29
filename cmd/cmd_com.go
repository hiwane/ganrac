package cmd

// cmd/*/*main.go から呼び出されることを想定

import (
	"github.com/chzyer/readline"
	"github.com/hiwane/ganrac"

	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type CmdParam struct {
	Verbose    int
	CadVerbose int
	Color      bool
	Quiet      bool

	CmdHistory string
}

type CmdLine struct {
	line        string
	pos         int
	in_string   bool
	depth_curly int

	rl *readline.Instance
}

/* コマンドパラメタ情報からの Ganrac の初期化 */
func (cp CmdParam) NewGanracLogger(cas, revision string) (*ganrac.Ganrac, *log.Logger) {
	if !cp.Quiet {
		if revision == "" {
			fmt.Printf("GaNRAC [cas=%s] see help();\n", cas)
		} else {
			fmt.Printf("GaNRAC [cas=%s, revision=%s] see help();\n", cas, revision)
		}
	}
	g := ganrac.NewGANRAC()
	logger := log.New(os.Stderr, "", log.Ltime)
	if cp.Color {
		ganrac.SetColordFml(true)
	}
	g.Eval(strings.NewReader(fmt.Sprintf("verbose(%d,%d);", cp.Verbose, cp.CadVerbose)))
	g.SetLogger(logger)
	return g, logger
}

func (cl *CmdLine) get_line() (string, error) {
	var err error
	ret := ""
	for {
		if cl.pos >= len(cl.line) {
			cl.line, err = cl.rl.Readline()
			if err != nil {
				return "", err
			}
			cl.pos = 0
		}

		for cl.pos < len(cl.line) {
			c := cl.line[cl.pos]
			cl.pos++
			if c == '"' {
				cl.in_string = !cl.in_string
			} else if cl.in_string {
				// do nothing.
			} else if c == '{' {
				cl.depth_curly++
			} else if c == '}' && cl.depth_curly > 0 {
				cl.depth_curly--
			} else if c == ';' { // eol
				goto _RETURN
			} else if c == ':' && cl.depth_curly <= 0 { // eolq
				goto _RETURN
			} else if c == '#' {
				ret += cl.line + "\n"
				cl.line = ""
				break
			}
		}
		ret += cl.line + "\n"
		cl.line = ""
	}
_RETURN:
	ret += cl.line[:cl.pos]
	cl.line = cl.line[cl.pos:]
	cl.pos = 0
	return ret, nil
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

type ganCompleter struct {
	funcs    []string
	examples []string
	qeopts   []string
	re_var   *regexp.Regexp
	re_func  *regexp.Regexp
	re_dict  *regexp.Regexp
}

func NewGanCompleter(g *ganrac.Ganrac) *ganCompleter {
	gc := new(ganCompleter)
	gc.re_var = regexp.MustCompile(`[a-z][a-z_]*$`)
	gc.re_func = regexp.MustCompile(`(cad|qe|help)\s*\(\s*"([a-z_]*)$`)
	gc.re_dict = regexp.MustCompile(`[,{]\s*(")?([a-z_]*)$`)
	gc.funcs = g.FunctionNames()
	gc.examples = ganrac.ExampleNames()
	gc.qeopts = ganrac.QEOptionNames()
	return gc
}

func (compl *ganCompleter) matchFunction(line, suf string, cand []string) (newLine [][]rune, length int) {
	newLine = make([][]rune, 0, 10)
	n := len(line)
	for _, s := range cand {
		if strings.HasPrefix(s, line) || line == "" {
			newLine = append(newLine, []rune(s[n:]+suf))
		}
	}

	return newLine, n
}

func (compl *ganCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	//    |      カーソル位置
	// 1234567
	// で TAB すると， ([49, 50, 51, 52, 53, 54, 55], 3) が復帰される

	// Readline will pass the whole line and current offset to it
	// Completer need to pass all the candidates, and how long they shared the same characters in line
	// Example:
	//   [go, git, git-shell, grep]
	//   Do("g", 1) => ["o", "it", "it-shell", "rep"], 1
	//   Do("gi", 2) => ["t", "t-shell"], 2
	//   Do("git", 3) => ["", "-shell"], 3
	matches := compl.re_func.FindStringSubmatch(string(line[:pos]))
	if len(matches) > 0 {
		if matches[1] == "help" {
			return compl.matchFunction(matches[2], "\");", compl.funcs)
		} else {
			return compl.matchFunction(matches[2], "\"", compl.examples)
		}
	}
	matches = compl.re_dict.FindStringSubmatch(string(line[:pos]))
	if len(matches) > 0 {
		// いまのところ qe コマンドくらいしかない
		return compl.matchFunction(matches[2], matches[1]+":", compl.qeopts)
	}
	match := compl.re_var.FindString(string(line[:pos]))
	if match != "" {
		return compl.matchFunction(match, "(", compl.funcs)
	}
	return nil, 0
}

func (cp CmdParam) Interpreter(g *ganrac.Ganrac) {

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "» ",
		HistoryFile:     cp.CmdHistory,
		AutoComplete:    NewGanCompleter(g),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
		//FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	rl.CaptureExitSignal()

	var cl CmdLine
	cl.rl = rl

	for {
		line, err := cl.get_line()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		p, err := g.Eval(strings.NewReader(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			continue
		}
		if p != nil {
			fmt.Println(p)
		}
	}
}
