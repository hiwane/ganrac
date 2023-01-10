package sage

import (
	"fmt"
	"testing"

	"github.com/hiwane/ganrac"
	castest "github.com/hiwane/ganrac/cas/test"
)

func TestFactor(t *testing.T) {
	funcname := "TestFactor"
	g := ganrac.NewGANRAC()
	sage, err := NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	castest.FactorTest(g, sage, t)
}

func TestResultant(t *testing.T) {
	funcname := "TestResultant"
	g := ganrac.NewGANRAC()
	sage, err := NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	castest.ResultantTest(g, sage, t)
}

func TestDiscrim(t *testing.T) {
	funcname := "TestDiscrim"
	g := ganrac.NewGANRAC()
	sage, err := NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	castest.DiscrimTest(g, sage, t)
}

func TestGBReduce(t *testing.T) {
	funcname := "TestGBReduce"
	g := ganrac.NewGANRAC()
	sage, err := NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	castest.GBRedTest(g, sage, t)
}

func TestPsc(t *testing.T) {
	funcname := "TestPsc"
	g := ganrac.NewGANRAC()
	sage, err := NewSage(g, "/tmp/ganrac.log")
	if err != nil {
		fmt.Printf("skip %s... init sage failed\n", funcname)
		return
	}

	castest.PscTest(g, sage, t)
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

	castest.SresTest(g, sage, t)
}
