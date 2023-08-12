package ganrac

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var DebugCounter map[string]int = make(map[string]int)

func IncDebugCounter(f string, n int) int {
	if _, ok := DebugCounter[f]; !ok {
		DebugCounter[f] = 0
	}
	DebugCounter[f] += n
	return DebugCounter[f]
}

var init_var_funcname string = "vars"

type func_table struct {
	name     string
	min, max int
	f        func(g *Ganrac, name string, args []interface{}) (interface{}, error)
	ox       bool
	args     string // 引数情報
	descript string // 1行説明
	help     string // 詳細説明
}

type Ganrac struct {
	varmap             map[string]interface{}
	sones, sfuns       []token
	history            []interface{}
	builtin_func_table []func_table
	ox                 CAS
	logger             *log.Logger
	logcnt             int
	verbose            int
	verbose_cad        int
	paranum            int
}

func NewGANRAC() *Ganrac {
	g := new(Ganrac)
	g.varmap = make(map[string]interface{}, 100)
	g.InitVarList([]string{
		"x", "y", "z", "w", "a", "b", "c", "d", "e", "f", "g", "h",
	})
	g.setBuiltinFuncTable()
	g.logger = log.New(ioutil.Discard, "", 0)
	g.sones = []token{
		{"+", plus},
		{"-", minus},
		{"*", mult},
		{"/", div},
		{"^", pow},
		{"[", lb},
		{"]", rb},
		{"{", lc},
		{"}", rc},
		{"(", lp},
		{")", rp},
		{",", comma},
		{";", eol},
		{":", eolq},
		{"==", eqop},
		{"=", assign},
		{"!=", neop},
		{"<=", leop},
		{"<", ltop},
		{">=", geop},
		{">", gtop},
		{"&&", and},
		{"||", or},
	}

	g.sfuns = []token{
		// {"impl", impl},
		// {"repl", repl},
		// {"equiv", equiv},
		// {"not", not},
		// {"all", all},
		// {"ex", ex},
		{init_var_funcname, initvar},
		{"time", f_time},
		{"true", f_true},
		{"false", f_false},
	}

	return g
}

func (g *Ganrac) SetParaNum(n int) {
	g.paranum = n
}

func (g *Ganrac) Close() {
	if g.ox != nil {
		g.ox.Close()
	}
}

func (g *Ganrac) addHisto(o interface{}) {
	g.history = append(g.history, o)
	if len(g.history) > 10 {
		g.history = g.history[len(g.history)-10:]
	}
}

func (g *Ganrac) genLexer(r io.Reader) *pLexer {
	lexer := newLexer(false)
	// yyErrorVerbose = true
	// yyDebug = 5
	lexer.Init(r)
	lexer.sones = g.sones
	lexer.sfuns = g.sfuns
	return lexer
}

func (g *Ganrac) parse(r io.Reader) (*pStack, error) {
	lexer := g.genLexer(r)
	yyParse(lexer)
	if lexer.err != nil {
		return nil, lexer.err
	}
	return lexer.stack, nil
}

func (g *Ganrac) Eval(r io.Reader) (interface{}, error) {
	stack, err := g.parse(r)
	if err != nil {
		return nil, err
	}
	pp, err := g.evalStack(stack)
	if err != nil {
		return nil, err
	}
	if pp == nil {
		return nil, nil
	}
	switch p := pp.(type) {
	case Fof:
		err = p.valid()
		if err != nil {
			return nil, err
		}
	case RObj:
		err = p.valid()
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return pp, nil
}

func (g *Ganrac) SetLogger(logger *log.Logger) {
	g.logger = logger
}

func (g *Ganrac) Logger() *log.Logger {
	return g.logger
}

func (g *Ganrac) SetCAS(cas CAS) {
	g.ox = cas
}

func (g *Ganrac) log(lv, caller int, format string, a ...interface{}) {
	if lv <= g.verbose {
		_, file, line, _ := runtime.Caller(caller)
		// runtime.FuncForPC(pc).Name()  (pc := 1st element of Caller())
		g.logcnt++
		_, fname := filepath.Split(file)
		fname = fname[:len(fname)-3] // .go は不要
		v := make([]interface{}, 4+len(a))
		v[0] = g.logcnt
		v[1] = lv
		v[2] = fname
		v[3] = line
		for i, au := range a {
			v[i+4] = au
		}
		g.logger.Printf("[%d,%d] %12.12s:%4d: "+format, v...)
	}
}

func maxint(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func gpanic(s string) {
	fmt.Fprintf(os.Stderr, "Please report the bug to https://github.com/hiwane/ganrac/issues\n")
	panic(s)
}
