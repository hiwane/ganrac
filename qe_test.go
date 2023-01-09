package ganrac

import (
	"fmt"
	"testing"
)

func benchmarkQE(b *testing.B, name string, cad bool) {
	input := GetExampleFof(name).Input
	g := NewGANRAC()
	ox := testConnectOx(g)
	if ox == nil {
		fmt.Printf("skip TestNeqQE... (no ox)\n")
		return
	}
	defer ox.Close()
	for i := 0; i < b.N; i++ {
		if cad {
			funcCAD(g, "cad", []interface{}{input})
		} else {
			opt := NewQEopt()
			g.QE(input, opt)
		}
	}
}

func BenchmarkAdam1(b *testing.B) {
	benchmarkQE(b, "adam1", true)
}

func TestBench(t *testing.T) {
}
