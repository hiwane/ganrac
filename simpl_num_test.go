package ganrac_test

import (
	"fmt"
	. "github.com/hiwane/ganrac"
	"testing"
)

func TestSimplNum(t *testing.T) {

	SetColordFml(true)
	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip TestSimplNum... (no cas)\n")
		return
	}
	defer g.Close()

	for ii, ss := range []struct {
		a      Fof
		expect Fof
	}{
		{
			// x^2-12<0 || (x+4!=0 && (x+4<0 || x+2>=0 || (x^2-12<=0 && x+2>=0))) || (x^2-12<=0 && x+4==0) || (x+2<=0 && x+4>=0)
			NewFmlOrs(
				NewAtom(NewPolyCoef(0, -12, 0, 1), LT),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 4, 1), NE),
					NewFmlOrs(
						NewAtom(NewPolyCoef(0, 4, 1), LT),
						NewAtom(NewPolyCoef(0, 2, 1), GE),
						NewFmlAnds(
							NewAtom(NewPolyCoef(0, -12, 0, 1), LE),
							NewAtom(NewPolyCoef(0, 2, 1), GE)))),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, -12, 0, 1), LE),
					NewAtom(NewPolyCoef(0, 4, 1), EQ)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 2, 1), LE),
					NewAtom(NewPolyCoef(0, 4, 1), GE))),
			// x^2-12<0 || (x+4!=0 && (x+4<0 || x+2>=0)) || (x+2<=0 && x+4>=0)
			NewFmlOrs(
				NewAtom(NewPolyCoef(0, -12, 0, 1), LT),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 4, 1), NE),
					NewFmlOrs(
						NewAtom(NewPolyCoef(0, 4, 1), LT),
						NewAtom(NewPolyCoef(0, 2, 1), GE))),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 2, 1), LE),
					NewAtom(NewPolyCoef(0, 4, 1), GE))),
		}, {
			// 20*z+y^2>=0
			NewAtom(NewPolyCoef(2, NewPolyVarn(1, 2), 20), GE),
			NewAtom(NewPolyCoef(2, NewPolyVarn(1, 2), 20), GE),
		}, {
			// adam2-1 && projection
			// 4*y^2+4*x^2-1<=0 && y^2-y+x^2-x==2
			NewFmlAnds(
				NewAtom(NewPolyCoef(1, NewPolyCoef(0, -1, 0, 100), 0, 100), LE),
				NewAtom(NewPolyCoef(1, NewPolyCoef(0, -2, -1, 1), -1, 1), EQ)),
			FalseObj,
		}, {
			// adam2-1 && projection
			// x>0 && y>0 && 4*y^2+4*x^2-1<=0 && y^2-y+x^2-x==2
			NewFmlAnds(
				NewAtom(NewPolyCoef(0, 0, 1), GT),
				NewAtom(NewPolyCoef(1, 0, 1), GT),
				NewAtom(NewPolyCoef(1, NewPolyCoef(0, -1, 0, 4), 0, 4), LE), NewAtom(NewPolyCoef(1, NewPolyCoef(0, -2, -1, 1), -1, 1), EQ)),
			FalseObj,
		}, {
			// (x-1)*(x-3)>=0 && x*(x-4) >= 0
			NewFmlAnds(NewAtom(NewPolyCoef(0, 3, -4, 1), GE), NewAtom(NewPolyCoef(0, 0, -4, 1), GE)),
			// x*(x-4) >= 0
			NewAtom(NewPolyCoef(0, 0, -4, 1), GE),
		}, {
			// all([x], (2*x^2-1>=0 && (x>0 || x+1<0)) || (x^2-1<=0 && (x<0 || 2*x^2-1<0)))
			NewFmlOrs(
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, -1, 0, 2), GE),
					NewFmlOrs(
						NewAtom(NewPolyCoef(0, 0, 1), GT),
						NewAtom(NewPolyCoef(0, 1, 1), LT)))),
			nil,
		}, {
			// all([x], (2*x^2-1>=0 && (x>0 || x+1<0)) || (x^2-1<=0 && (x<0 || 2*x^2-1<0)))
			NewQuantifier(true, []Level{0}, NewFmlOrs(NewFmlAnds(NewAtom(NewPolyCoef(0, -1, 0, 2), GE), NewFmlOrs(NewAtom(NewPolyCoef(0, 0, 1), GT), NewAtom(NewPolyCoef(0, 1, 1), LT))), NewFmlAnds(NewAtom(NewPolyCoef(0, -1, 0, 1), LE), NewFmlOrs(NewAtom(NewPolyCoef(0, 0, 1), LT), NewAtom(NewPolyCoef(0, -1, 0, 2), LT))))),
			nil,
		}, {
			NewAtom(NewPolyCoef(0, 1, 0, 1), GT),
			TrueObj,
		}, {
			NewAtom(NewPolyCoef(0, 1, 0, 1), LE),
			FalseObj,
		}, {
			// (x-2)^2 >= 10 && (x-1)^2 >= 2
			NewFmlAnds(NewAtom(NewPolyCoef(0, -6, -4, 1), GE), NewAtom(NewPolyCoef(0, -1, -2, 1), GE)),
			// (x-2)^2 >= 10
			NewAtom(NewPolyCoef(0, -6, -4, 1), GE),
		}, {
			// x>2 && x^2+y^2 > 1
			NewFmlAnds(NewAtom(NewPolyCoef(0, -2, 1), GT), NewAtom(NewPolyCoef(1, NewPolyCoef(0, -1, 0, 1), 0, 1), GT)),
			NewAtom(NewPolyCoef(0, -2, 1), GT), // x>2
		}, {
			NewFmlAnds(NewAtom(NewPolyCoef(1, -1, 1), GT), NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -1), 0, 1), LT)),
			NewFmlAnds(NewAtom(NewPolyCoef(1, -1, 1), GT), NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -1), 0, 1), LT)),
		}, {
			// (3*x+1!=0 && x<0) || (3*x-1>=0 && 3*x-2<0)
			NewFmlOrs(
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 1, 3), NE),
					NewAtom(NewPolyCoef(0, 0, 1), LT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, -1, 3), GE),
					NewAtom(NewPolyCoef(0, -2, 3), LT))),
			nil,
		}, {
			// (3*x+1<0 && 3*x-1<=0) || ((3*x-2)*(3*x+1)<0 && (x<0 || 3*x-1>=0));
			NewFmlOrs(
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 1, 3), LT),
					NewAtom(NewPolyCoef(0, -1, 3), LE)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, -2, -3, 9), LT),
					NewFmlOrs(
						NewAtom(NewPolyCoef(0, 0, 1), LT),
						NewAtom(NewPolyCoef(0, -1, 3), GE)))),
			// 3*x+1<0 || ((3*x-2)*(3*x+1)<0 && (x<0 || 3*x-1>=0))
			NewFmlOrs(
				NewAtom(NewPolyCoef(0, 1, 3), LT),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, -2, -3, 9), LT),
					NewFmlOrs(
						NewAtom(NewPolyCoef(0, 0, 1), LT),
						NewAtom(NewPolyCoef(0, -1, 3), GE)))),
		}, {
			// x*(729*x^3+8856*x^2+43632*x+33856)<=0 && 162*x^3+1971*x^2+8772*x+4048>0;
			NewFmlAnds(NewAtom(NewPolyCoef(0, 0, 33856, 43632, 8856, 729), LE), NewAtom(NewPolyCoef(0, 4048, 8772, 1971, 162), GT)),
			nil,
		},
	} {
		if ss.expect == nil {
			ss.expect = ss.a
		}

		for jj, s := range []struct {
			a      Fof
			expect Fof
		}{
			{ss.a, ss.expect},
			{ss.a.Not(), ss.expect.Not()},
		} {
			// fmt.Printf("===== in=%v\n", s.a)
			o, tf, ff := SimplNum(s.a, g, nil, nil)
			// fmt.Printf("i=%v\n", s.a)
			// fmt.Printf("o=%v\n", o)
			if !o.Equals(s.expect) {
				t.Errorf("<%d,%d>\ninput =%v\nexpect=%v\noutput=%v\nt=%v\nf=%v\n", ii, jj, s.a, s.expect, o, tf, ff)
				return
			}
		}
	}
}

func TestSimplNumUniPoly(test *testing.T) {
	lv := Level(0)
	for _, s := range []struct {
		a  *Poly
		t  []NObj
		f  []NObj
		op OP
	}{
		{
			NewPolyCoef(lv, 2, 0, -1),
			[]NObj{NewInt(0), nil},
			[]NObj{},
			OP_TRUE,
		},
	} {
		t := NewNumRegion()
		for i := 0; i+1 < len(s.t); i += 2 {
			t.Append(lv, s.t[i], s.t[i+1])
		}
		f := NewNumRegion()
		for i := 0; i+1 < len(s.f); i += 2 {
			f.Append(lv, s.f[i], s.f[i+1])
		}

		op, pos, neg := s.a.SimplNumUniPoly(t, f)
		if op != s.op {
			test.Errorf("a=%v, t=%v, f=%v\nexpect=%v\nactual=%v\npos=%v\nneg=%v\n",
				s.a, t, f, s.op, op, pos, neg)
		}
	}
}
