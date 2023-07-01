package para

import (
	"fmt"
	"testing"

	"github.com/hiwane/ganrac"
	"github.com/hiwane/ganrac/cas/sage"
	castest "github.com/hiwane/ganrac/cas/test"
)

func TestFactor(t *testing.T) {
	funcname := "TestFactor"
	g := ganrac.NewGANRAC()
	sage, err := sage.NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	para := NewParaCAS(sage)
	defer para.Close()
	castest.FactorTest(g, para, t)
}

func TestResultant(t *testing.T) {
	funcname := "TestResultant"
	g := ganrac.NewGANRAC()
	sage, err := sage.NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	para := NewParaCAS(sage)
	defer para.Close()
	castest.ResultantTest(g, para, t)
}

func TestDiscrim(t *testing.T) {
	funcname := "TestDiscrim"
	g := ganrac.NewGANRAC()
	sage, err := sage.NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	para := NewParaCAS(sage)
	defer para.Close()
	castest.DiscrimTest(g, para, t)
}

func TestGBReduce(t *testing.T) {
	funcname := "TestGBReduce"
	g := ganrac.NewGANRAC()
	sage, err := sage.NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	para := NewParaCAS(sage)
	defer para.Close()
	castest.GBRedTest(g, para, t)
}

func TestPsc(t *testing.T) {
	funcname := "TestPsc"
	g := ganrac.NewGANRAC()
	sage, err := sage.NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	para := NewParaCAS(sage)
	defer para.Close()
	castest.PscTest(g, para, t)
}

func TestSlope(t *testing.T) {
	funcname := "TestSlope"
	g := ganrac.NewGANRAC()
	sage, err := sage.NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		t.Errorf("stop")
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	para := NewParaCAS(sage)
	defer para.Close()
	castest.SlopeTest(g, para, t)
}
