package sage

import (
	"fmt"
	"testing"

	"github.com/hiwane/ganrac"
	"github.com/hiwane/ganrac/cas"
)

func TestFactor(t *testing.T) {
	funcname := "TestFactor"
	g := ganrac.NewGANRAC()
	sage, err := NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	cas.FactorTest(g, sage, t)
}

func TestResultant(t *testing.T) {
	funcname := "TestResultant"
	g := ganrac.NewGANRAC()
	sage, err := NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	cas.ResultantTest(g, sage, t)
}

func TestDiscrim(t *testing.T) {
	funcname := "TestDiscrim"
	g := ganrac.NewGANRAC()
	sage, err := NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	cas.DiscrimTest(g, sage, t)
}

func TestGBReduce(t *testing.T) {
	funcname := "TestGBReduce"
	g := ganrac.NewGANRAC()
	sage, err := NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	cas.GBRedTest(g, sage, t)
}

func TestPsc(t *testing.T) {
	funcname := "TestPsc"
	g := ganrac.NewGANRAC()
	sage, err := NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	cas.PscTest(g, sage, t)
}

func TestSres(t *testing.T) {
	funcname := "TestSres"
	g := ganrac.NewGANRAC()
	sage, err := NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		t.Errorf("stop")
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	cas.SresTest(g, sage, t)
}
