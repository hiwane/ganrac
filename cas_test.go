package ganrac_test

import (
	"github.com/hiwane/ganrac"
	//	"github.com/hiwane/ganrac/cas/sage"
	openxm "github.com/hiwane/ganrac/cas/ox"
	"testing"

	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func evalstr(g *ganrac.Ganrac, s string) (interface{}, error) {
	if !strings.HasSuffix(s, ";") {
		s += ";"
	}
	return g.Eval(strings.NewReader(s))
}

func str2poly(g *ganrac.Ganrac, s string) (*ganrac.Poly, error) {
	x, err := evalstr(g, s)
	if err != nil {
		return nil, err
	}

	y, ok := x.(*ganrac.Poly)
	if !ok {
		return nil, fmt.Errorf("not a polynomial: %v", s)
	}
	return y, nil
}

func str2fofq(g *ganrac.Ganrac, s string) (ganrac.FofQ, error) {
	x, err := evalstr(g, s)
	if err != nil {
		return nil, err
	}

	y, ok := x.(ganrac.FofQ)
	if !ok {
		return nil, fmt.Errorf("not a quantified formula: %v", s)
	}
	return y, nil
}

func str2atom(g *ganrac.Ganrac, s string) (*ganrac.Atom, error) {
	x, err := evalstr(g, s)
	if err != nil {
		return nil, err
	}

	y, ok := x.(*ganrac.Atom)
	if !ok {
		return nil, fmt.Errorf("not an atomic formula: %v", s)
	}
	return y, nil
}

func str2fof(g *ganrac.Ganrac, s string) (ganrac.Fof, error) {
	x, err := evalstr(g, s)
	if err != nil {
		return nil, err
	}

	y, ok := x.(ganrac.Fof)
	if !ok {
		return nil, fmt.Errorf("not a formula: %v", s)
	}
	return y, nil
}

func makeCAS(t *testing.T) *ganrac.Ganrac {
	m := 1
	g := ganrac.NewGANRAC()
	var c ganrac.CAS
	logger := log.New(os.Stderr, "", log.LstdFlags)
	g.SetLogger(logger)
	// g.Eval(strings.NewReader(fmt.Sprintf("verbose(%d,%d);", 9, 9)))
	if m == 0 {
		// NOTE: sage版で例外が発生すると，python上で発生したようになってデバッグが大変
		// ox がインストールされていても
		// python.so がないと実行すらできない状態になってしまう
		// 健全でない気がする
		c = makeSage(t, g)
	} else if m == 1 {
		c = makeOX(t, g)
	} else {
		return nil
	}
	if c == nil {
		return nil
	}
	g.SetCAS(c)
	// c.SetLogger(logger)
	return g
}

func makeSage(t *testing.T, g *ganrac.Ganrac) ganrac.CAS {
	return nil
	// s, err := sage.NewSage(g, "/tmp/ganrac_cas_test.tmp")
	// if err != nil {
	// 	if t != nil {
	// 		t.Errorf("generate sage failed: %s", err.Error())
	// 	}
	// }
	// return s
}

func makeOX(t *testing.T, g *ganrac.Ganrac) ganrac.CAS {
	time.Sleep(time.Second / 2)
	for i := 1; i <= 10; i++ {
		ox := _makeOX(t, g)
		if ox != nil {
			return ox
		}
		fmt.Printf("waiting for openxm server... %d sec\n", i*3)
		time.Sleep(time.Second * 3)
	}
	return nil
}

func _makeOX(t *testing.T, g *ganrac.Ganrac) ganrac.CAS {
	cport := "localhost:1234"
	dport := "localhost:4321"
	connc, err := net.Dial("tcp", cport)
	if err != nil {
		return nil
	}

	time.Sleep(time.Second / 20)

	connd, err := net.Dial("tcp", dport)
	if err != nil {
		connc.Close()
		return nil
	}

	ox, err := openxm.NewOpenXM(connc, connd, g.Logger())
	if err != nil {
		connc.Close()
		connd.Close()
		return nil
	}

	return ox
}

func checkEquivalentQff(g *ganrac.Ganrac, p, q ganrac.Fof) string {
	if !p.IsQff() {
		return "checkEquivalent(p,q): p is not a QFF formula"
	}
	if !q.IsQff() {
		return "checkEquivalent(p,q): q is not a QFF formula"
	}

	fof := ganrac.NewForAll([]ganrac.Level{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, ganrac.NewFmlEquiv(p, q))
	qeopt := ganrac.NewQEopt()
	qeopt.Qe_init(g, fof)
	cond := ganrac.NewQeCond()

	ret := qeopt.QE(fof, *cond)
	if _, ok := ret.(*ganrac.AtomT); ok {
		return ""
	} else {
		return "checkEquivalent(p,q): q is not equivalent"
	}
}
