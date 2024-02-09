package sympy

import (
	"fmt"
	"testing"

	"github.com/hiwane/ganrac"
	castest "github.com/hiwane/ganrac/cas/test"
)

var sympy_tmpfname = "/tmp/ganrac_sympy.log"

func TestFactor(t *testing.T) {
	funcname := "TestFactor"
	g := ganrac.NewGANRAC()
	sage, err := NewSymPy(g, sympy_tmpfname)
	if err != nil {
		fmt.Printf("skip %s... init sympy failed\n", funcname)
		return
	}

	castest.FactorTest(g, sage, t)
}

func TestResultant(t *testing.T) {
	funcname := "TestResultant"
	g := ganrac.NewGANRAC()
	sage, err := NewSymPy(g, sympy_tmpfname)
	if err != nil {
		fmt.Printf("skip %s... init sympy failed\n", funcname)
		return
	}

	castest.ResultantTest(g, sage, t)
}

func TestDiscrim(t *testing.T) {
	funcname := "TestDiscrim"
	g := ganrac.NewGANRAC()
	sage, err := NewSymPy(g, sympy_tmpfname)
	if err != nil {
		fmt.Printf("skip %s... init sympy failed\n", funcname)
		return
	}

	castest.DiscrimTest(g, sage, t)
}

func TestGBReduce(t *testing.T) {
	funcname := "TestGBReduce"
	g := ganrac.NewGANRAC()
	sage, err := NewSymPy(g, sympy_tmpfname)
	if err != nil {
		fmt.Printf("skip %s... init sympy failed\n", funcname)
		return
	}

	castest.GBRedTest(g, sage, t)
}

func TestPsc(t *testing.T) {
	funcname := "TestPsc"
	g := ganrac.NewGANRAC()
	sage, err := NewSymPy(g, sympy_tmpfname)
	if err != nil {
		fmt.Printf("skip %s... init sympy failed\n", funcname)
		return
	}

	castest.PscTest(g, sage, t)
}

func TestSlope(t *testing.T) {
	funcname := "TestSlope"
	g := ganrac.NewGANRAC()
	sage, err := NewSymPy(g, sympy_tmpfname)
	if err != nil {
		t.Errorf("stop")
		fmt.Printf("skip %s... init sympy failed\n", funcname)
		return
	}

	castest.SlopeTest(g, sage, t)
}

func TestSymPySres(t *testing.T) {
	funcname := "TestSymPySres"
	g := ganrac.NewGANRAC()
	sage, err := NewSymPy(g, sympy_tmpfname)
	if err != nil {
		t.Errorf("stop")
		fmt.Printf("skip %s... init sympy failed\n", funcname)
		return
	}
	castest.SresTest(g, sage, t)
}

func TestSymPyGB(t *testing.T) {
	funcname := "TestSymPyGB"
	g := ganrac.NewGANRAC()
	sage, err := NewSymPy(g, sympy_tmpfname)
	if err != nil {
		t.Errorf("stop")
		fmt.Printf("skip %s... init sympy failed\n", funcname)
		return
	}
	castest.GBTest(g, sage, t)
}
