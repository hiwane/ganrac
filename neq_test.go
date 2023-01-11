package ganrac_test

import (
	"fmt"
	. "github.com/hiwane/ganrac"
	"testing"
)

func TestNeqQE(t *testing.T) {
	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip TestNeqQE... (no CAS)\n")
		return
	}
	defer g.Close()

	for ii, ss := range []struct {
		input  Fof
		expect Fof
	}{
		{
			// ex([x], a*x^2+b*x+c != 0)
			NewQuantifier(false, []Level{3}, NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 1), NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), NE)),
			NewFmlOrs(NewAtom(NewPolyCoef(0, 0, 1), NE), NewAtom(NewPolyCoef(2, 0, 1), NE), NewAtom(NewPolyCoef(1, 0, 1), NE)),
		}, {
			// ex([x], a*x+b != 0 && c*x+d < 0);
			// <==>
			// (a != 0 || b != 0) && (c != 0 || d < 0)
			NewQuantifier(false, []Level{4}, NewFmlAnds(
				NewAtom(NewPolyCoef(4, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), NE),
				NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1)), LT))),
			NewFmlAnds(
				NewFmlOrs(
					NewAtom(NewPolyCoef(0, 0, 1), NE),
					NewAtom(NewPolyCoef(1, 0, 1), NE)),
				NewFmlOrs(
					NewAtom(NewPolyCoef(2, 0, 1), NE),
					NewAtom(NewPolyCoef(3, 0, 1), LT))),
		}, {
			// ex([x], a*x+b != 0 && c*x+d <= 0);
			// <==>
			// (a != 0 || b != 0) && (c != 0 || d <= 0)
			NewQuantifier(false, []Level{4}, NewFmlAnds(
				NewAtom(NewPolyCoef(4, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), NE),
				NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1)), LE))),
			NewFmlAnds(
				NewFmlOrs(
					NewAtom(NewPolyCoef(0, 0, 1), NE),
					NewAtom(NewPolyCoef(1, 0, 1), NE)),
				NewFmlOrs(
					NewAtom(NewPolyCoef(2, 0, 1), NE),
					NewAtom(NewPolyCoef(3, 0, 1), LE))),
		}, {
			//   ex([x], a*x+b != 0 && c^2*d*x+e < 0);
			// <==>
			//   (a != 0 || b != 0) && (c < 0 || d^4-4ec > 0 || e < 0)
			NewQuantifier(false, []Level{5}, NewFmlAnds(
				NewAtom(NewPolyCoef(5, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), NE),
				NewAtom(NewPolyCoef(5, NewPolyCoef(4, 0, 1), NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1)), LT))),
			NewFmlAnds(
				NewFmlOrs(
					NewAtom(NewPolyCoef(0, 0, 1), NE),
					NewAtom(NewPolyCoef(1, 0, 1), NE)),
				NewFmlOrs(
					NewAtom(NewPolyCoef(2, 0, 1), LT),
					NewAtom(NewPolyCoef(4, 0, 1), LT),
					NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 0, -1), NewPolyCoef(2, 0, 4)), LT))),
		}, {
			//   ex([x], a*x+b != 0 && c^2*d*x+e <= 0);
			// <==>
			//   (a != 0 || b != 0) && (c < 0 || d^4-4ec > 0 || e < 0)
			NewQuantifier(false, []Level{5}, NewFmlAnds(
				NewAtom(NewPolyCoef(5, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), NE),
				NewAtom(NewPolyCoef(5, NewPolyCoef(4, 0, 1), NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1)), LE))),
			NewFmlOrs(
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), NE),
					NewAtom(NewPolyCoef(2, 0, 1), LE),
					NewAtom(NewPolyCoef(4, 0, 1), LE)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), NE),
					NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 0, -1), NewPolyCoef(2, 0, 4)), LT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(1, 0, 1), NE),
					NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 0, -1), NewPolyCoef(2, 0, 4)), LT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(1, 0, 1), NE),
					NewAtom(NewPolyCoef(4, 0, 1), LE)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, NewPolyCoef(1, 0, -2)), NewPolyCoef(0, 0, 1)), NE),
					NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 0, -1), NewPolyCoef(2, 0, 4)), EQ))),
		}, {
			//      ex([x], a*x+b != 0 && s*x^4 + 4*x^3 - 8*x+4 <= 0);
			// <==> s <= 1 && (a != 0 || b != 0)
			NewQuantifier(false, []Level{3}, NewFmlAnds(
				NewAtom(NewPolyCoef(3, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), NE),
				NewAtom(NewPolyCoef(3, 4, -8, 0, 4, NewPolyCoef(2, 0, 1)), LE))),
			nil,
			// NewFmlAnds(
			// 	NewAtom(NewPolyCoef(2, -1, 1), LE),
			// 	NewFmlOrs(
			// 		NewAtom(NewPolyCoef(0, 0, 1), NE),
			// 		NewAtom(NewPolyCoef(1, 0, 1), NE))),
		}, {
			//      ex([x], a*x+b != 0 && s*x^5 + (x^2+2*x-2) <= 0);
			// <==> (a != 0 || b != 0)
			NewQuantifier(false, []Level{3}, NewFmlAnds(
				NewAtom(NewPolyCoef(3, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), NE),
				NewAtom(NewPolyCoef(3, -2, 2, 1, 0, 0, NewPolyCoef(2, 0, 1)), LE))),
			NewFmlOrs(NewAtom(NewPolyCoef(0, 0, 1), NE), NewAtom(NewPolyCoef(1, 0, 1), NE)),
		}, {
			//      ex([x], a*x+b != 0 && s*x^5 + (x^2+2*x-2)^2 <= 0);
			// <==> (a != 0 || b != 0)
			NewQuantifier(false, []Level{3}, NewFmlAnds(
				NewAtom(NewPolyCoef(3, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), NE),
				NewAtom(NewPolyCoef(3, 4, -8, 0, 4, 1, NewPolyCoef(2, 0, 1)), LE))),
			NewFmlOrs(NewAtom(NewPolyCoef(0, 0, 1), NE), NewAtom(NewPolyCoef(1, 0, 1), NE)),
		},
	} {
		opt := NewQEopt()

		f := ss.input.(FofQ)
		cond := NewQeCond()
		opt.Qe_init(g, f)

		// fmt.Printf("ii=%d: %s\n", ii, f)
		h := opt.Qe_neq(f, *cond)
		// fmt.Printf("h=%v\n", h)
		if h == nil {
			if ss.expect == nil {
				continue
			}
			t.Errorf("ii=%d, neqQE not worked: %v", ii, ss.input)
			continue
		} else if ss.expect == nil {
			t.Errorf("ii=%d, neqQE WORKED: %v", ii, ss.input)
			continue
		}

		vars := []Level{0, 1, 2, 3, 4, 5}
		opt2 := NewQEopt()
		opt2.Algo &= ^QEALGO_NEQ // NEQ は使わない
		u := NewQuantifier(true, vars, NewFmlEquiv(ss.expect, h))
		// fmt.Printf("u=%v\n", u)
		if _, ok := g.QE(u, opt2).(*AtomT); ok {
			continue
		}
		t.Errorf("ii=%d\ninput= %v.\nexpect= %v.\nactual= %v.\n", ii, ss.input, ss.expect, h)
		return
	}
}
