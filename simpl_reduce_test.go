package ganrac_test

import (
	"fmt"
	. "github.com/hiwane/ganrac"
	"testing"
)

func TestSimplReduce(t *testing.T) {
	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip TestSimplReduce... (no cas)\n")
		return
	}
	defer g.Close()

	opt := NewQEopt()
	opt.Algo = 0 // CAD で評価する

	x := NewPolyVar(0)
	y := NewPolyVar(1)
	z := NewPolyVar(2)

	for ii, ss := range []struct {
		input  Fof
		expect Fof // simplified input
	}{
		{
			NewFmlAnds(NewAtom(x, EQ), NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, x), 1), GT)), // x==0 && z+x*y>0
			NewFmlAnds(NewAtom(x, EQ), NewAtom(z, GT)),                                       // x==0 && z>0
		}, {
			NewFmlAnds(NewAtom(y, EQ), NewAtom(NewPolyCoef(1, 1, 1), GT)), // y==0 && y+1 > 0
			NewFmlAnds(NewAtom(y, EQ)),                                    // y==0
		}, {
			NewFmlAnds(
				NewAtom(NewPolyCoef(1, -1, 0, 2), EQ),
				NewAtom(NewPolyCoef(1, 1, 1), EQ)), // 2*y^2+1=0 && y+1 = 0
			FalseObj,
		}, {
			NewFmlAnds(NewAtom(NewPolyCoef(0, 0, 1), EQ),
				NewQuantifier(false, []Level{0}, NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, NewPolyCoef(0, 0, 1)), 1), EQ))), // x==0 && ex([x], z+x*y==0)
			nil,
		}, {
			NewFmlAnds(NewAtom(NewPolyCoef(0, 0, 1), EQ), NewQuantifier(false, []Level{0}, NewAtom(NewPolyCoef(0, 1, 1), GT))),
			NewAtom(NewPolyCoef(0, 0, 1), EQ),
		}, {
			NewFmlAnds( // x==0 && z==0 && ex([x], z+y+x==0)
				NewAtom(NewPolyCoef(0, 0, 1), EQ),
				NewAtom(NewPolyCoef(2, 2, 1), EQ),
				NewQuantifier(false, []Level{0},
					NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 1), 1), 1), EQ))),
			NewFmlAnds(
				NewAtom(NewPolyCoef(0, 0, 1), EQ),
				NewAtom(NewPolyCoef(2, 2, 1), EQ),
				NewQuantifier(false, []Level{0},
					NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 1), 1), EQ))),
		}, {
			// x==0 && ex([w], x*w^2+y*w+1==0 && y*w^3+z*w+x==0 && w-1<=0)
			NewFmlAnds(
				NewAtom(NewPolyCoef(0, 0, 1), EQ),
				NewQuantifier(false, []Level{3}, NewFmlAnds(
					NewAtom(NewPolyCoef(3, 1, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), EQ),
					NewAtom(NewPolyCoef(3, NewPolyCoef(0, 0, 1), NewPolyCoef(2, 0, 1), 0, NewPolyCoef(1, 0, 1)), EQ),
					NewAtom(NewPolyCoef(3, -1, 1), LE)))),
			// x==0 && ex([w], y*w+1==0 && y*w^3+z*w==0 && w-1<=0)
			NewFmlAnds(
				NewAtom(NewPolyCoef(0, 0, 1), EQ),
				NewQuantifier(false, []Level{3}, NewFmlAnds(
					NewAtom(NewPolyCoef(3, 1, NewPolyCoef(1, 0, 1)), EQ),
					NewAtom(NewPolyCoef(3, 0, NewPolyCoef(2, 0, 1), 0, NewPolyCoef(1, 0, 1)), EQ),
					NewAtom(NewPolyCoef(3, -1, 1), LE)))),
		},
	} {
		if ss.expect == nil {
			ss.expect = ss.input
		}
		for jj, s := range []struct {
			input  Fof
			expect Fof
		}{
			{ss.input, ss.expect},
			{ss.input.Not(), ss.expect.Not()},
		} {
			// fmt.Printf("[%d,%d] s=%v\n", ii, jj, s.input)
			inf := NewReduceInfo()
			o := SimplReduce(s.input, g, inf)
			if TTestSameFormAndOr(o, s.expect) {
				continue
			}

			u := NewFmlEquiv(o, s.expect)
			switch uqe := g.QE(u, opt).(type) {
			case *AtomT:
				continue
			default:
				t.Errorf("<%d,%d>\n input=%v\nexpect=%v\noutput=%v\ncmp=%v", ii, jj, s.input, s.expect, o, uqe)
			}
		}
	}
}
