package ganrac_test

import (
	"github.com/hiwane/ganrac"
	"github.com/hiwane/ganrac/cas/sage"
	"testing"

	"log"
	"os"
)

func makeCAS(t *testing.T) *ganrac.Ganrac {
	m := 0
	g := ganrac.NewGANRAC()
	var c ganrac.CAS
	logger := log.New(os.Stderr, "", log.LstdFlags)
	g.SetLogger(logger)
	// g.Eval(strings.NewReader(fmt.Sprintf("verbose(%d,%d);", 9, 9)))
	if m == 0 {
		// ox がインストールされていても
		// python.so がないと実行すらできない状態になってしまう
		// 健全でない気がする
		c = makeSage(t, g)
	} else {
		return nil
	}
	g.SetCAS(c)
	return g
}

func makeSage(t *testing.T, g *ganrac.Ganrac) ganrac.CAS {
	s, err := sage.NewSage(g, "/tmp/ganrac_cas_test.tmp")
	if err != nil {
		if t != nil {
			t.Errorf("generate sage failed: %s", err.Error())
		}
	}
	return s
}
