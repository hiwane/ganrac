package ganrac_test

import (
	"fmt"
	. "github.com/hiwane/ganrac"
	"testing"
)

func benchmarkQE(b *testing.B, name string, cad bool) {
	input := GetExampleFof(name).Input
	g := makeCAS(nil)
	if g == nil {
		fmt.Printf("skip benchmarkQE... (no cas)\n")
		return
	}
	defer g.Close()
	for i := 0; i < b.N; i++ {
		if cad {
			FuncCAD(g, "cad", []interface{}{input})
		} else {
			opt := NewQEopt()
			g.QE(input, opt)
		}
	}
}

func BenchmarkCADAdam1(b *testing.B) {
	benchmarkQE(b, "adam1", true)
}

func BenchmarkQEAdam1(b *testing.B) {
	benchmarkQE(b, "adam1", true)
}

func BenchmarkQEAdam3(b *testing.B) {
	benchmarkQE(b, "adam3", true)
}

func TestBench(t *testing.T) {
}
