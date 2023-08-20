package ganrac_test

import (
	"fmt"
	. "github.com/hiwane/ganrac"
	"strings"
	"testing"
)

func TestSimplReduce1(t *testing.T) {
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
			inf := NewReduceInfo(g, s.input, TrueObj, FalseObj)
			o := SimplReduce(s.input, g, inf)
			if TTestSameFormAndOr(o, s.expect) {
				continue
			}

			u := NewFmlEquiv(o, s.expect)
			ux := g.QE(u, opt)
			switch uqe := ux.(type) {
			case *AtomT:
				continue
			default:
				t.Errorf("<%d,%d>\n input=%v\nexpect=%v\noutput=%v\ncmp=%v", ii, jj, s.input, s.expect, o, uqe)
			}
		}
	}
}

func TestSimplReduce2(t *testing.T) {
	funcname := "TestSimplReduce2"
	//print_log := false
	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip %s... (no cas)\n", funcname)
		return
	}
	defer g.Close()

	lvs := []Level{0, 1, 2, 3, 4, 5, 6, 7}

	qeopt := NewQEopt()
	qeopt.SetG(g)

	cadopt := NewDict()
	cadopt.Set("var", NewInt(1))

	for i, ss := range []struct {
		eqpoly string
		other  string
		expect string
	}{
		{
			"a",
			"a*x^2-b*x+c > 0",
			"-b*x+c > 0",
		}, {
			"a",
			//			"ex([a], a != 0 && b == 0)",
			"a*x^2-b*x+a >= 0",
			"b*x <= 0",
		},
	} {
		for j, vars := range []struct {
			v string // 変数順序
		}{
			{"x,a,b,c"},
			{"a,c,x,b"},
			{"a,b,x,c"},
			{"c,b,a,x"},
		} {
			vstr := fmt.Sprintf("vars(%s);", vars.v)
			_, err := g.Eval(strings.NewReader(vstr))
			if err != nil {
				t.Errorf("[%d] %s failed: %s", i, vstr, err)
				return
			}

			eqpol, err := str2poly(g, ss.eqpoly)
			if err != nil {
				t.Errorf("[%d,%d,%s] eval(input) failed: %s, %s", i, j, vstr, err, ss.eqpoly)
				return
			}
			other, err := str2fof(g, ss.other)
			if err != nil {
				t.Errorf("[%d,%d,%s] eval(other) failed: %s, %s", i, j, vstr, err, ss.other)
				return

			}
			expect, err := str2fof(g, ss.expect)
			if err != nil {
				t.Errorf("[%d,%d,%s] eval(expect) failed: %s, %s", i, j, vstr, err, ss.expect)
				return
			}

			eqfof := NewAtom(eqpol, EQ)
			nefof := NewAtom(eqpol, NE)
			for _, tt := range []struct {
				input  Fof
				neccon Fof
				sufcon Fof
				expect Fof
			}{
				{
					NewFmlAnd(other, eqfof),
					TrueObj,
					FalseObj,
					NewFmlAnd(expect, eqfof),
				}, {
					other,
					eqfof,
					FalseObj,
					expect,
				}, {
					other,
					TrueObj,
					nefof,
					expect,
				},
			} {
				inf := NewReduceInfo(g, tt.input, tt.neccon, tt.sufcon)
				ret := SimplReduce(tt.input, g, inf)
				//				fmt.Printf("ret=%v\n", ret)

				if err := ValidFof(ret); err != nil {
					t.Errorf("[%d,%d] invalid simpl_reduce\ninput=%s\nactual=%v\nactual=%V", i, j, ss, ret, ret)
					return
				}

				expect2 := tt.expect
				ret2 := ret

				if !ret.IsQff() {
					ret2 = g.QE(ret, qeopt)
				}
				if !expect.IsQff() {
					expect2 = g.QE(expect, qeopt)
				}

				c, err := FuncCAD(g, "CAD", []interface{}{
					NewForAll(lvs, NewFmlEquiv(expect2, ret2)), cadopt})
				if err != nil {
					t.Errorf("[%d,%d] invalid simpl_reduce; cad\ninput=%s\nactual=%v\nexpect=%v\nerr=%v", i, j, ss, ret, expect, err)
					return
				}
				if c != TrueObj {
					t.Errorf("[%d,%d] invalid simpl_reduce; cad\ninput=%s\nneccon=%v\nsufcon=%v\nactual=%v\na _ _ =%v\nexpect=%v", i, j, tt.input, tt.neccon, tt.sufcon, ret, ret2, expect)
					return
				}
			}
		}
	}
}
