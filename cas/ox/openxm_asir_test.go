package ganrac

import (
	"fmt"
	"github.com/hiwane/ganrac"
	castest "github.com/hiwane/ganrac/cas/test"
	"testing"
)

func TestAsirFactor(t *testing.T) {
	funcname := "TestAsirFactor"

	g := ganrac.NewGANRAC()
	ox := testConnectOx(g)
	if ox == nil {
		fmt.Printf("skip %s... (no ox)\n", funcname)
		return
	}
	defer ox.Close()
	castest.FactorTest(g, ox, t)
}

func TestAsirResultant(t *testing.T) {
	funcname := "TestAsirResultant"
	g := ganrac.NewGANRAC()
	ox := testConnectOx(g)
	if ox == nil {
		fmt.Printf("skip %s... (no ox)\n", funcname)
		return
	}
	defer ox.Close()
	castest.ResultantTest(g, ox, t)
}

func TestAsirDiscrim(t *testing.T) {
	funcname := "TestAsirDiscrim"
	g := ganrac.NewGANRAC()
	ox := testConnectOx(g)
	if ox == nil {
		fmt.Printf("skip %s... (no ox)\n", funcname)
		return
	}
	defer ox.Close()
	castest.DiscrimTest(g, ox, t)
}

func TestAsirGBReduce(t *testing.T) {
	funcname := "TestAsirGBReduce"
	g := ganrac.NewGANRAC()
	ox := testConnectOx(g)
	if ox == nil {
		fmt.Printf("skip %s... (no ox)\n", funcname)
		return
	}
	defer ox.Close()
	castest.GBRedTest(g, ox, t)
}

func TestAsirPsc(t *testing.T) {
	funcname := "TestAsirPsc"
	g := ganrac.NewGANRAC()
	ox := testConnectOx(g)
	if ox == nil {
		fmt.Printf("skip %s... (no ox)\n", funcname)
		return
	}
	defer ox.Close()
	castest.PscTest(g, ox, t)
}

func TestAsirSlope(t *testing.T) {
	funcname := "TestAsirSlope"
	g := ganrac.NewGANRAC()
	ox := testConnectOx(g)
	if ox == nil {
		fmt.Printf("skip %s... (no ox)\n", funcname)
		return
	}
	defer ox.Close()
	castest.SlopeTest(g, ox, t)
}
