package ganrac

import (
	"fmt"
)

func testConnectOx(g *Ganrac) CAS {
	return nil
}

func testConnectCAS(g *Ganrac) CAS {
	ox := new(testCAS)
	g.SetCAS(ox)
	return ox
}

type testCAS struct {
}

func (*testCAS) Gcd(p, q *Poly) RObj {
	fname := "Gcd"
	fmt.Printf("%s p=%v\n", fname, p)
	fmt.Printf("%s q=%v\n", fname, q)
	return p
}

func (*testCAS) Factor(p *Poly) *List {
	fname := "Factor"
	for _, q := range []*Poly{
		NewPolyCoef(1, 0, 1),
		NewPolyCoef(0, 0, 1),
		NewPolyCoef(2, 0, 1),
		NewPolyCoef(3, 0, 1),
		NewPolyCoef(3, NewPolyCoef(2, 0, NewPolyCoef(1, 0, -1)), NewPolyCoef(0, 0, 1)),
		NewPolyCoef(2, 0, NewPolyCoef(1, 0, 1)),
		NewPolyCoef(1, NewPolyCoef(0, 0, 4), 1),
	} {
		if p.Equals(q) { // 既約
			return NewList(NewList(one, one), NewList(p, one))
		}
	}

	fmt.Printf("%s p=%v\n", fname, p)
	fmt.Printf("%s p=%S\n", fname, p)
	return NewList()
}

func (*testCAS) Discrim(p *Poly, lv Level) RObj {
	fname := "Discrim"
	fmt.Printf("%s p[%d]=%v\n", fname, lv, p)
	fmt.Printf("%s p=(%d, %S)\n", fname, lv, p)
	return p
}
func (*testCAS) Resultant(p *Poly, q *Poly, lv Level) RObj {
	fname := "Resultant"
	fmt.Printf("%s p[%d]=%v\n", fname, lv, p)
	fmt.Printf("%s q[%d]=%v\n", fname, lv, q)
	fmt.Printf("%s p=(%d, %S, %s)\n", fname, lv, p, q)
	return p
}
func (*testCAS) Psc(p *Poly, q *Poly, lv Level, j int32) RObj {
	fname := "Psc"
	fmt.Printf("%s p[%d]=%v\n", fname, lv, p)
	fmt.Printf("%s q[%d]=%v\n", fname, lv, q)
	fmt.Printf("%s j=%d\n", fname, j)
	return p
}
func (*testCAS) Sres(p *Poly, q *Poly, lv Level, k int32) RObj {
	fname := "Sres"
	fmt.Printf("%s p[%d]=%v\n", fname, lv, p)
	fmt.Printf("%s q[%d]=%v\n", fname, lv, q)
	fmt.Printf("%s k=%d\n", fname, k)
	return p
}
func (*testCAS) GB(p *List, vars *List, n int) *List {

	for _, q := range []struct {
		n   int
		p   *List
		v   *List
		ret *List
	}{
		{
			0,
			NewList(NewPolyCoef(0, 0, 1)),
			NewList(NewPolyCoef(0, 0, 1)),
			nil,
		}, {
			0,
			NewList(NewPolyCoef(1, 0, 1)),
			NewList(NewPolyCoef(1, 0, 1)),
			nil,
		}, {
			0,
			NewList(NewPolyCoef(2, 0, 1)),
			NewList(NewPolyCoef(2, 0, 1)),
			nil,
		}, {
			0,
			NewList(NewPolyCoef(3, 0, 1)),
			NewList(NewPolyCoef(3, 0, 1)),
			nil,
		}, {
			0,
			NewList(NewPolyCoef(2, 0, NewPolyCoef(1, 0, 1))),
			NewList(NewPolyCoef(1, 0, 1), NewPolyCoef(2, 0, 1)),
			nil,
		}, {
			0,
			NewList(NewPolyCoef(2, 0, NewPolyCoef(1, 0, 1)), NewPolyCoef(3, 0, 1)),
			NewList(NewPolyCoef(1, 0, 1), NewPolyCoef(2, 0, 1), NewPolyCoef(3, 0, 1)),
			nil,
		}, {
			0,
			NewList(NewPolyCoef(1, NewPolyCoef(0, 0, 4), 1)),
			NewList(NewPolyCoef(0, 0, 1), NewPolyCoef(1, 0, 1)),
			nil,
		}, {
			0,
			NewList(NewPolyCoef(1, NewPolyCoef(0, 0, 4), 1), NewPolyCoef(0, 0, 1)),
			NewList(NewPolyCoef(0, 0, 1), NewPolyCoef(1, 0, 1)),
			NewList(NewPolyCoef(0, 0, 1), NewPolyCoef(1, 0, 1)),
		}, {
			0,
			NewList(NewPolyCoef(0, 0, 1), NewPolyCoef(1, NewPolyCoef(0, 0, 4), 1)),
			NewList(NewPolyCoef(0, 0, 1), NewPolyCoef(1, 0, 1)),
			NewList(NewPolyCoef(0, 0, 1), NewPolyCoef(1, 0, 1)),
		},
	} {
		if n == q.n && p.Equals(q.p) && vars.Equals(q.v) {
			return p
		}
	}

	fname := "GB"
	fmt.Printf("%s p=%v\n", fname, p)
	fmt.Printf("%s v[%d]=%v\n", fname, n, vars)
	if n == 0 {
		fmt.Printf("nd_gr(%v, %v, 0, 0)\n", p, vars)
	} else {
		fmt.Printf("nd_gr(%v, %v, 0, [[0, %d], [0, %d]]);\n", p, vars, vars.Len()-n, n)
	}
	fmt.Printf("}, {\n%d,\n%S,\n%S,\nnil,\n", n, p, vars)
	return NewList()
}
func (*testCAS) Reduce(p *Poly, gb *List, vars *List, n int) (RObj, bool) {
	f := false
	for i := Level(0); i <= p.lv && !f; i++ {
		if p.hasVar(i) {
			for j := 0; j < gb.Len(); j++ {
				g := gb.geti(j)
				if g.(*Poly).hasVar(i) {
					f = true
					break
				}
			}
		}
	}
	if !false {
		return p, false
	}

	for _, q := range []struct {
		n   int
		p   *Poly
		gb  *List
		v   *List
		ret RObj
		b   bool
	}{
		{
			0,
			NewPolyCoef(2, 0, 1),
			NewList(NewPolyCoef(2, 0, 1)),
			NewList(NewPolyCoef(2, 0, 1)),
			zero, false,
		}, {
			2,
			NewPolyCoef(2, 0, NewPolyCoef(1, 0, 1)),
			NewList(NewPolyCoef(3, 0, 1)),
			NewList(NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1), NewPolyCoef(1, 0, 1)),
			NewPolyCoef(2, 0, NewPolyCoef(1, 0, 1)),
			false,
		}, {
			0,
			NewPolyCoef(3, 0, 1),
			NewList(NewPolyCoef(3, 0, 1)),
			NewList(NewPolyCoef(3, 0, 1)),
			zero, false,
		}, {
			1,
			NewPolyCoef(3, 0, 1),
			NewList(NewPolyCoef(2, 0, NewPolyCoef(1, 0, 1))),
			NewList(NewPolyCoef(1, 0, 1), NewPolyCoef(2, 0, 1), NewPolyCoef(3, 0, 1)),
			NewPolyCoef(3, 0, 1), false,
		},
	} {
		if n == q.n && p.Equals(q.p) && gb.Equals(q.gb) && vars.Equals(q.v) {
			return q.ret, q.b
		}
	}

	fname := "Reduce"
	fmt.Printf("%s p=%v\n", fname, p)
	fmt.Printf("%s gb=%v\n", fname, gb)
	fmt.Printf("%s vars=%v\n", fname, vars)
	fmt.Printf("%s n=%v\n", fname, n)
	fmt.Printf("p_true_nf(%v, %v, %v, 0);\n", p, gb, vars)
	fmt.Printf("}, {\n%d,\n%S,\n%S,\n%S,\n", n, p, gb, vars)
	panic("bo")
}
func (*testCAS) Eval(p string) (GObj, error) {
	fname := "Eval"
	fmt.Printf("%s p=%v\n", fname, p)
	return zero, nil
}

func (*testCAS) Close() error {
	return nil
}
