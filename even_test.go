package ganrac_test

import (
	"fmt"
	. "github.com/hiwane/ganrac"
	"testing"
)

func TestEvenSimpl(t *testing.T) { // @TEST @TODO
	funcname := "EvenSimpl"

	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip TestEvenSimpl... (no cas)\n")
		return
	}
	defer g.Close()
	vars := []Level{0, 1, 2, 3, 4, 5}

	for ii, ss := range []struct {
		input  Fof
		expect Fof
	}{
		{
			// ex([x], a*x^2+b >= 0);
			NewQuantifier(false, []Level{2}, NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE)),
			NewFmlOrs(NewAtom(NewPolyCoef(0, 0, 1), GT), NewAtom(NewPolyCoef(1, 0, 1), GE)),
		},
		{
			// ex([x], a*x^2+b >= 0 && c*x^2+d >= 0);
			// <==>
			// [ a > 0 /\ c > 0 ] \/ [ a > 0 /\ d >= 0 /\ a d - b c >= 0 ] \/ [ c > 0 /\ d < 0 /\ a d - b c <= 0 ] \/ [ b >= 0 /\ d >= 0 ]
			NewQuantifier(false, []Level{4}, NewFmlAnds(
				NewAtom(NewPolyCoef(4, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 1), 0, NewPolyCoef(2, 0, 1)), GE))),
			NewFmlOrs(
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(2, 0, 1), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(3, 0, 1), GE),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, NewPolyCoef(1, 0, -1)), NewPolyCoef(0, 0, 1)), GE)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(2, 0, 1), GT),
					NewAtom(NewPolyCoef(3, 0, 1), LT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, NewPolyCoef(1, 0, -1)), NewPolyCoef(0, 0, 1)), LE)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(1, 0, 1), GE),
					NewAtom(NewPolyCoef(3, 0, 1), GE))),
		},
		{
			// ex([x], a*x^2+b >= 0 && x > 2);
			// <==>
			// a > 0 \/ b + 4 a > 0 \/ [ a = 0 /\ b + 4 a = 0 ]
			NewQuantifier(false, []Level{2}, NewFmlAnds(
				NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(2, -2, 1), GT))),
			NewFmlOrs(
				NewAtom(NewPolyCoef(0, 0, 1), GT),
				NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 4), 1), GT),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), EQ),
					NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 4), 1), EQ))),
		},
		{
			// ex([x], a*x^2+b >= 0 && x > -2);
			// <==>
			// a > 0 \/ b >= 0
			NewQuantifier(false, []Level{2}, NewFmlAnds(
				NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(2, +2, 1), GT))),
			NewFmlOrs(
				NewAtom(NewPolyCoef(0, 0, 1), GT),
				NewAtom(NewPolyCoef(1, 0, 1), GE)),
		},
		{
			// ex([x], a*x+b >= 0 && 3*x>2);
			// <==>
			// a > 0 \/ 9 b + 4 a > 0 \/ [ a = 0 /\ 9 b + 4 a = 0 ]
			// a > 0 || 9*b+4*a > 0 || (a == 0 && 9*b+4*a==0);
			NewQuantifier(false, []Level{3}, NewFmlAnds(
				NewAtom(NewPolyCoef(3, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(3, -2, 3), GT))),
			NewFmlOrs(
				NewAtom(NewPolyCoef(0, 0, 1), GT),
				NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 4), 9), GT), NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), EQ),
					NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, 4), 9), EQ))),
		},
		{
			// (E x) [ a x^2 + b >= 0 /\ 3 x > c ].
			// <==>
			// a > 0 \/ a c^2 + 9 b > 0 \/ [ a = 0 /\ a c^2 + 9 b = 0 ] \/ [ b >= 0 /\ c < 0 ]
			// a > 0 || a*c^2+9*b > 0 || (a == 0 && a*c^2+9*b == 0) || (b >= 0 && c < 0);
			NewQuantifier(false, []Level{3}, NewFmlAnds(
				NewAtom(NewPolyCoef(3, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, -1), 3), GT))),
			NewFmlOrs(NewAtom(NewPolyCoef(0, 0, 1), GT), NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 9), 0, NewPolyCoef(0, 0, 1)), GT), NewFmlAnds(NewAtom(NewPolyCoef(0, 0, 1), EQ), NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 9), 0, NewPolyCoef(0, 0, 1)), EQ)), NewFmlAnds(NewAtom(NewPolyCoef(1, 0, 1), GE), NewAtom(NewPolyCoef(2, 0, 1), LT))),
		},
		{
			// (E x) [ a x^2 + b >= 0 /\ c x > 7 ].
			// <==>
			// c /= 0 /\ [ a > 0 \/ b c^2 + 49 a > 0 \/ [ b = 0 /\ b c^2 + 49 a = 0 ] ]
			// c != 0 && (a > 0 || b*c^2+49*a > 0 || (b == 0 && b*c^2+49*a == 0));
			NewQuantifier(false, []Level{3}, NewFmlAnds(NewAtom(NewPolyCoef(3, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE), NewAtom(NewPolyCoef(3, -7, NewPolyCoef(2, 0, 1)), GT))),
			NewFmlAnds(NewAtom(NewPolyCoef(2, 0, 1), NE), NewFmlOrs(NewAtom(NewPolyCoef(0, 0, 1), GT), NewAtom(NewPolyCoef(2, NewPolyCoef(0, 0, 49), 0, NewPolyCoef(1, 0, 1)), GT), NewFmlAnds(NewAtom(NewPolyCoef(1, 0, 1), EQ), NewAtom(NewPolyCoef(2, NewPolyCoef(0, 0, 49), 0, NewPolyCoef(1, 0, 1)), EQ)))),
		},
		{
			// (E x) [ a x^2 + b >= 0 /\ c  x + d > 0 ].
			// <==>
			// [ a >= 0 /\ c > 0 /\ a d^2 + b c^2 = 0 ] \/ [ a >= 0 /\ c < 0 /\ a d^2 + b c^2 = 0 ] \/ [ d > 0 /\ a d^2 + b c^2 > 0 ] \/ [ a > 0 /\ a d^2 + b c^2 < 0 ] \/ [ c > 0 /\ a d^2 + b c^2 > 0 ] \/ [ b >= 0 /\ d > 0 ] \/ [ c < 0 /\ a d^2 + b c^2 > 0 ]

			NewQuantifier(false, []Level{4}, NewFmlAnds(
				NewAtom(NewPolyCoef(4, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1)), GT))),
			NewFmlOrs(
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GE),
					NewAtom(NewPolyCoef(2, 0, 1), GT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), EQ)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GE),
					NewAtom(NewPolyCoef(2, 0, 1), LT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), EQ)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(3, 0, 1), GT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), LT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(2, 0, 1), GT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(1, 0, 1), GE),
					NewAtom(NewPolyCoef(3, 0, 1), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(2, 0, 1), LT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), GT))),
		},
		{
			// (E x) [ a x^2 + b >= 0 /\ c  x + d >= 0 ].
			// <==>
			// [ d >= 0 /\ a d^2 + b c^2 > 0 ] \/ [ a > 0 /\ a d^2 + b c^2 <= 0 ] \/ [ c > 0 /\ a d^2 + b c^2 >= 0 ] \/ [ b >= 0 /\ d >= 0 ] \/ [ c < 0 /\ a d^2 + b c^2 >= 0 ]
			NewQuantifier(false, []Level{4}, NewFmlAnds(
				NewAtom(NewPolyCoef(4, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1)), GE))),
			NewFmlOrs(
				NewFmlAnds(
					NewAtom(NewPolyCoef(3, 0, 1), GE),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), LE)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(2, 0, 1), GT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), GE)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(1, 0, 1), GE),
					NewAtom(NewPolyCoef(3, 0, 1), GE)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(2, 0, 1), LT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), GE))),
		},
		{
			// (E x) [ a x^2 + b >= 0 /\ c  x + d < 0 ].
			// <==>
			// [ a >= 0 /\ c > 0 /\ a d^2 + b c^2 = 0 ] \/ [ a >= 0 /\ c < 0 /\ a d^2 + b c^2 = 0 ] \/ [ a > 0 /\ a d^2 + b c^2 < 0 ] \/ [ a > 0 /\ d < 0 ] \/ [ c > 0 /\ a d^2 + b c^2 > 0 ] \/ [ c < 0 /\ a d^2 + b c^2 > 0 ] \/ [ b >= 0 /\ d < 0 ]
			NewQuantifier(false, []Level{4}, NewFmlAnds(
				NewAtom(NewPolyCoef(4, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1)), LT))),
			NewFmlOrs(
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GE),
					NewAtom(NewPolyCoef(2, 0, 1), GT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), EQ)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GE),
					NewAtom(NewPolyCoef(2, 0, 1), LT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), EQ)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), LT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(3, 0, 1), LT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(2, 0, 1), GT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(2, 0, 1), LT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(1, 0, 1), GE),
					NewAtom(NewPolyCoef(3, 0, 1), LT))),
		},
		{
			// (E x) [ a x^2 + b >= 0 /\ c  x + d <= 0 ].
			// <==>
			// [ a > 0 /\ a d^2 + b c^2 < 0 ] \/ [ a > 0 /\ d <= 0 ] \/ [ c > 0 /\ a d^2 + b c^2 >= 0 ] \/ [ c < 0 /\ a d^2 + b c^2 >= 0 ] \/ [ b >= 0 /\ d <= 0 ]
			NewQuantifier(false, []Level{4}, NewFmlAnds(
				NewAtom(NewPolyCoef(4, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1)), LE))),
			NewFmlOrs(
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), LT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(3, 0, 1), LE)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(2, 0, 1), GT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), GE)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(2, 0, 1), LT),
					NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), GE)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(1, 0, 1), GE),
					NewAtom(NewPolyCoef(3, 0, 1), LE))),
		},
		{
			// (E x) [ a x^2 + b >= 0 /\ c  x + d = 0 ].
			// <==>
			// a d^2 + b c^2 >= 0 /\ [ c /= 0 \/ [ b >= 0 /\ d = 0 ] \/ [ a > 0 /\ a d^2 + b c^2 = 0 ] ]

			NewQuantifier(false, []Level{4}, NewFmlAnds(
				NewAtom(NewPolyCoef(4, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1)), EQ))),
			NewFmlAnds(
				NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), GE),
				NewFmlOrs(
					NewAtom(NewPolyCoef(2, 0, 1), NE),
					NewFmlAnds(
						NewAtom(NewPolyCoef(1, 0, 1), GE),
						NewAtom(NewPolyCoef(3, 0, 1), EQ)),
					NewFmlAnds(
						NewAtom(NewPolyCoef(0, 0, 1), GT),
						NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, 0, NewPolyCoef(1, 0, 1)), 0, NewPolyCoef(0, 0, 1)), EQ)))),
		},
		{
			// (E x) [ a x^2 + b >= 0 /\ c  x + d /= 0 ].
			// <==>
			// [ a = 0 /\ b >= 0 /\ c > 0 ] \/ [ a = 0 /\ b >= 0 /\ c < 0 ] \/ [ a > 0 /\ c > 0 ] \/ [ a > 0 /\ c < 0 ] \/ [ b > 0 /\ c > 0 ] \/ [ b > 0 /\ c < 0 ] \/ [ a > 0 /\ d > 0 ] \/ [ a > 0 /\ d < 0 ] \/ [ b >= 0 /\ d > 0 ] \/ [ b >= 0 /\ d < 0 ]
			NewQuantifier(false, []Level{4}, NewFmlAnds(
				NewAtom(NewPolyCoef(4, NewPolyCoef(1, 0, 1), 0, NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1)), NE))),
			NewFmlOrs(
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), EQ),
					NewAtom(NewPolyCoef(1, 0, 1), GE),
					NewAtom(NewPolyCoef(2, 0, 1), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), EQ),
					NewAtom(NewPolyCoef(1, 0, 1), GE),
					NewAtom(NewPolyCoef(2, 0, 1), LT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(2, 0, 1), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(2, 0, 1), LT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(1, 0, 1), GT),
					NewAtom(NewPolyCoef(2, 0, 1), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(1, 0, 1), GT),
					NewAtom(NewPolyCoef(2, 0, 1), LT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(3, 0, 1), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(3, 0, 1), LT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(1, 0, 1), GE),
					NewAtom(NewPolyCoef(3, 0, 1), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(1, 0, 1), GE),
					NewAtom(NewPolyCoef(3, 0, 1), LT))),
		},
	} {
		f := ss.input.(FofQ)

		opt := NewQEopt()
		cond := NewQeCond()
		opt.Qe_init(g, f)

		h := opt.Qe_evenq(f, *cond, 2)
		if h == nil {
			if ss.expect != nil {
				t.Errorf("ii=%d, %s not worked: %v", ii, funcname, ss.input)
				continue
			}
		} else if ss.expect == nil {
			t.Errorf("ii=%d, %s WORKED: %v", ii, funcname, ss.input)
			continue
		} else {

			opt2 := NewQEopt()
			opt2.DelAlgo(QEALGO_SMPL_EVEN) // 使わない
			u := NewQuantifier(true, vars, NewFmlEquiv(ss.expect, h))
			// fmt.Printf("u=%v\n", u)
			if _, ok := g.QE(u, opt2).(*AtomT); !ok {
				t.Errorf("ii=%d %s\ninput= %v.\nexpect= %v.\nactual= %v.\n", ii, funcname, ss.input, ss.expect, h)
				break
			}
		}

		fnot := f.Not()
		hnot := opt.Qe_evenq(fnot, *cond, 2)
		if hnot == nil {
			if ss.expect != nil {
				t.Errorf("ii=%d, %s.not not worked: %v", ii, funcname, ss.input)
				continue
			}
		} else if ss.expect == nil {
			t.Errorf("ii=%d, %s.not WORKED: %v", ii, funcname, ss.input)
			continue
		} else {
			opt2 := NewQEopt()
			opt2.DelAlgo(QEALGO_SMPL_EVEN) // 使わない
			u := NewQuantifier(true, vars, NewFmlEquiv(ss.expect.Not(), hnot))
			// fmt.Printf("u=%v\n", u)
			if _, ok := g.QE(u, opt2).(*AtomT); ok {
				continue
			}
			t.Errorf("ii=%d %s.not\ninput= %v.\nexpect= %v.\nactual= %v.\n", ii, funcname, fnot, ss.expect.Not(), h)
		}
	}
}
