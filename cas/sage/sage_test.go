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
