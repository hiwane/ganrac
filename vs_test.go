package ganrac_test

import (
	. "github.com/hiwane/ganrac"

	"fmt"
	"strings"
	"testing"
)

func TestVsLin1(t *testing.T) {

	g := NewGANRAC()

	for i, ss := range []struct {
		qff    string
		expect string
	}{
		{"x*y >0", "x != 0;"},
		{"y>=0", "true;"},
		{"x+y >=0", "true;"},
		{"x+y >0", "true;"},
		{"x*y <0", "x != 0;"},
		{"x*y >=0", "true;"},
		{"(x+1)*y <0", "x != -1;"},
		{"(2-x)*y >0", "x != 2;"},
		{"x+y <0", "true;"},
		{"x+y >=0", "true;"},
		{"x+y >=0 && x + y + 2 <= 0", "false;"},
		{"(y+1)*x-3 >=0", "x!=0;"},
		{"(x+1)*y-3 >=0", "x!=-1;"},
		{"y > a && y <= b", "a < b;"},
		{"y == a && y != b && a > 0", "a != b && a > 0;"},
		{"y == x && y != b && x > 0", "x != b && x > 0;"},
	} {
		for j, s := range []struct {
			qff    string
			expect string
		}{ // 再帰表現なので，自由変数と束縛変数のレベルの大小で動きが異なる
			{ss.qff, ss.expect},
			{strings.ReplaceAll(ss.qff, "x", "z"),
				strings.ReplaceAll(ss.expect, "x", "z")}} {

			//fmt.Printf("==============[%d,%d]==================================\n%v\n", i, j, s.qff)
			input := fmt.Sprintf("ex([y], %s);", s.qff)
			_fof, err := g.Eval(strings.NewReader(input))
			if err != nil {
				t.Errorf("%d-%d: eval failed input=`%s`: err:`%s`", i, j, input, err)
				return
			}
			fof, ok := _fof.(Fof)
			if !ok {
				t.Errorf("%d: eval failed\ninput=%s\neval=%v\nerr:%s", i, input, fof, err)
				return
			}

			ans, err := g.Eval(strings.NewReader(s.expect))
			if err != nil {
				t.Errorf("%d: eval failed input=%s: err:%s", i, s.expect, err)
				return
			}

			lv := Level(1)
			// fmt.Printf(">>>>>>>> go   %v[%s]\n", fof, VarStr(lv))
			qff := QeVS(fof, lv, 1, g)
			// fmt.Printf("<<<<<<<< done %v[%s]\n", fof, VarStr(lv))
			if err = ValidFof(qff); err != nil {
				t.Errorf("%d: formula is broken input=`%s`: out=`%s`, %v", i, s.qff, qff, err)
				return
			}

			if HasVar(qff, lv) {
				t.Errorf("%d: variable %d is not eliminated input=%s: out=%s", i, lv, s.qff, qff)
				return
			}

			if !qff.Equals(ans) {
				// qff.dump(os.Stdout)
				// fmt.Printf("\n")
				//
				// fmt.Printf("----------------------------\n")
				var q Fof
				lllv := []Level{0, 2, 3, 4, 5, 6}
				q = NewQuantifier(true, lllv, NewFmlEquiv(qff, ans.(Fof)))
				for _, llv := range lllv {
					q = QeVS(q, llv, 1, g)
					if HasVar(q, llv) {
						t.Errorf("%d: variable %d is not eliminated: X1 out=%s", i, llv, q)
						return
					}
				}
				if _, ok := q.(*AtomT); !ok {
					fmt.Printf("q=%v\n", q)
					t.Errorf("%d: qe failed\ninput =%s\nexpect=%v\nactual=%v", i, input, s.expect, qff)
					return
				}
			}

			sqff := SimplBasic(qff, TrueObj, FalseObj)
			if err = ValidFof(sqff); err != nil {
				t.Errorf("%d: formula is broken input=`%v`: out=`%v`, %v", i, qff, sqff, err)
				return
			}

			if !sqff.Equals(ans) {
				var q Fof
				lllv := []Level{0, 2, 3, 4, 5, 6}
				q = NewQuantifier(true, lllv, NewFmlEquiv(sqff, ans.(Fof)))
				for _, llv := range lllv {
					q = QeVS(q, llv, 1, g)
					if HasVar(q, llv) {
						t.Errorf("%d: variable %d is not eliminated: X1 out=%s", i, llv, q)
						return
					}
				}
				if _, ok := q.(*AtomT); !ok {
					fmt.Printf("q=%v\n", q)
					t.Errorf("%d: qe failed\ninput =%s\nexpect=%v\nactual=%v", i, qff, s.expect, s.qff)
					return
				}
			}
		}
	}
}

func TestVsLin2(t *testing.T) {

	gan := NewGANRAC() // log 用

	for ii, ss := range []struct {
		lv     Level
		p      Fof
		expect Fof
	}{
		{2, // ex([x], a*x == b && 3*x > 1);
			NewQuantifier(false, []Level{2}, NewFmlAnds(
				NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, -1), NewPolyCoef(0, 0, 1)), EQ),
				NewAtom(NewPolyCoef(2, -1, 3), GT))),
			NewFmlOrs( // [ a = 0 /\ 3 b - a = 0 ] \/ [ a > 0 /\ 3 b - a > 0 ] \/ [ a < 0 /\ 3 b - a < 0 ]
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), EQ),
					NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -1), 3), EQ)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), GT),
					NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -1), 3), GT)),
				NewFmlAnds(
					NewAtom(NewPolyCoef(0, 0, 1), LT),
					NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -1), 3), LT))),
		}, {2, // ex([x], a*x < b && 3*x > 1);
			NewQuantifier(false, []Level{2}, NewFmlAnds(
				NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, -1), NewPolyCoef(0, 0, 1)), LT),
				NewAtom(NewPolyCoef(2, -1, 3), GT))),
			NewFmlOrs(
				NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -1), 3), GT),
				NewAtom(NewPolyCoef(0, 0, 1), LT)),
		}, {2, // ex([x], a*x > b && 3*x > 1);
			NewQuantifier(false, []Level{2}, NewFmlAnds(
				NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, -1), NewPolyCoef(0, 0, 1)), GT),
				NewAtom(NewPolyCoef(2, -1, 3), GT))),
			NewFmlOrs(
				NewAtom(NewPolyCoef(0, 0, 1), GT),
				NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -1), 3), LT)),
		}, {2,
			// ex([x], a*x+b >= 0 && 3*x+1 > 0)
			NewQuantifier(false, []Level{2}, NewFmlAnds(
				NewAtom(NewPolyCoef(2, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(2, 1, 3), GT))),
			// b >= 0 \/ a > 0 \/ 3 b - a > 0
			NewFmlOrs(
				NewAtom(NewPolyCoef(1, 0, 1), GE),
				NewAtom(NewPolyCoef(0, 0, 1), GT),
				NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -1), 3), GT)),
		}, {4,
			// ex([x], a*x+b >= 0 && c*x+d > 0)
			NewQuantifier(false, []Level{4}, NewFmlAnds(
				NewAtom(NewPolyCoef(4, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)), GE),
				NewAtom(NewPolyCoef(4, NewPolyCoef(3, 0, 1), NewPolyCoef(2, 0, 1)), GT))),
			//  b >= 0 && d > 0  ||  a >= 0 && c > 0 && a*d - b*c <= 0  ||  a <= 0 && c < 0 && a*d - b*c >= 0  ||  a > 0 && a*d - b*c > 0  ||  a < 0 && a*d - b*c < 0
			NewFmlOrs(
				NewFmlAnds(NewAtom(NewPolyCoef(1, 0, 1), GE), NewAtom(NewPolyCoef(3, 0, 1), GT)),
				NewFmlAnds(NewAtom(NewPolyCoef(0, 0, 1), GE), NewAtom(NewPolyCoef(2, 0, 1), GT), NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, NewPolyCoef(1, 0, -1)), NewPolyCoef(0, 0, 1)), LE)),
				NewFmlAnds(NewAtom(NewPolyCoef(0, 0, 1), LE), NewAtom(NewPolyCoef(2, 0, 1), LT), NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, NewPolyCoef(1, 0, -1)), NewPolyCoef(0, 0, 1)), GE)),
				NewFmlAnds(NewAtom(NewPolyCoef(0, 0, 1), GT), NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, NewPolyCoef(1, 0, -1)), NewPolyCoef(0, 0, 1)), GT)),
				NewFmlAnds(NewAtom(NewPolyCoef(0, 0, 1), LT), NewAtom(NewPolyCoef(3, NewPolyCoef(2, 0, NewPolyCoef(1, 0, -1)), NewPolyCoef(0, 0, 1)), LT))),
		},
	} {
		for jj, sss := range []struct {
			p      Fof
			expect Fof
		}{
			{ss.p, ss.expect},
			{ss.p.Not(), ss.expect.Not()},
		} {
			f := QeVS(sss.p, ss.lv, 1, gan)
			f = SimplBasic(f, TrueObj, FalseObj)

			q := make([]Level, ss.lv)
			for i := Level(0); i < ss.lv; i++ {
				q[i] = i
			}
			g := NewQuantifier(true, q, NewFmlEquiv(f, sss.expect))
			for i := Level(0); i < ss.lv; i++ {
				g = QeVS(g, i, 1, gan)
				switch g.(type) {
				case *AtomF:
					t.Errorf("invalid %d, %d, F\n in=%v\nexp=%v\nout=%v", ii, jj, sss.p, sss.expect, f)
					return
				case FofQ:
					continue
				case *AtomT:
					break
				default:
					t.Errorf("invalid %d, %d, 2, %v\n in=%v\nexp=%v\nout=%v", ii, jj, g, sss.p, sss.expect, f)
					return
				}
				g = SimplBasic(g, TrueObj, FalseObj)
			}
			if _, ok := g.(*AtomT); !ok {
				t.Errorf("invalid %d, %d, 3, %v\n in=%v\nexp=%v\nout=%v", ii, jj, g, sss.p, sss.expect, f)
				return
			}
		}
	}
}

func TestVsLin3(t *testing.T) {
	funcname := "TestVsLin3"
	print_log := false
	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip %s... (no cas)\n", funcname)
		return
	}
	defer g.Close()

	optbl := []string{"==", "!=", ">=", "<=", ">", "<"}

	for i, ss := range []string{
		"a*x-b %s 0 && x %s c",
		"a*x-b %s 0 && c*x %s -1",
	} {
		for j, vars := range []struct {
			v  string // 変数順序
			lv Level  // 消去する変数のレベル
		}{
			{"b,x,c,a", 1},
			{"a,b,c,x", 3},
			{"a,b,x,c", 2},
			{"x,c,b,a", 0},
		} {
			vstr := fmt.Sprintf("vars(%s);", vars.v)
			_, err := g.Eval(strings.NewReader(vstr))
			if err != nil {
				t.Errorf("[%d] %s failed: %s", i, vstr, err)
				return
			}
			for ops := 0; ops < 6*6; ops++ {
				op1 := optbl[ops/6]
				op2 := optbl[ops%6]
				input := fmt.Sprintf("ex([x], %s);", fmt.Sprintf(ss, op1, op2))
				if print_log {
					fmt.Printf(" == [%d,%d,%d] ==== %s ==[%s:%d]======================\n", i, j, ops, input, vars.v, vars.lv)
				}
				_fof, err := g.Eval(strings.NewReader(input))
				if err != nil {
					t.Errorf("[%d] parse error: %s, %s, %s\nerr=%s\ninput=%s", i, vstr, op1, op2, err, input)
					return
				}
				fof, ok := _fof.(Fof)
				if !ok {
					t.Errorf("[%d] eval failed\ninput=%s\neval=%v", i, ss, _fof)
					return
				}
				if err := ValidFof(fof); err != nil {
					t.Errorf("[%d] invalid fvs\ninput=%s\nactual=%v", i, ss, fof)
					return
				}

				fvs := QeVS(fof, vars.lv, 1, g)
				if print_log {
					fmt.Printf(">>>>>>>>>>>>\n")
				}
				if err := ValidFof(fvs); err != nil {
					t.Errorf("[%d] invalid fvs\ninput=%s\nactual=%v", i, ss, fvs)
					return
				}
				if HasVar(fvs, vars.lv) {
					t.Errorf("[%d,%d,%s] vs2 failed: has `%s`[%d]\ninput=%s\nactual=%s", i, j, vars.v, VarStr(vars.lv), vars.lv, input, fvs)
					return
				}

				// VS2次を使わないモードで
				opt := NewQEopt()
				opt.SetAlgo(QEALGO_VSQUAD|QEALGO_VSLIN, false)
				fqe := g.QE(fof, opt)
				// fmt.Printf(">>>>>>>>>>>>EQUIV\n")

				eq := g.QE(NewForAll([]Level{0, 1, 2, 3}, NewFmlEquiv(fvs, fqe)), opt)
				if eq != TrueObj {
					t.Errorf("[%d,%d,%d,'%s'] invalid\ninputt=%v\nexpect=%v\nactual=%v\nactua2=%v\nequall=%v", i, j, ops, vars.v, fof, fqe, fvs, g.SimplFof(fvs), eq)

					opt.SetAlgo(QEALGO_VSLIN, false)
					for _, lvs := range [][]Level{
						{},
						{0, 1},
						{0, 2},
						{0, 3},
						{1, 2},
						{1, 3},
						{2, 3},
					} {
						if len(lvs) == 0 {
							break
						}
						vv := g.QE(NewForAll(lvs, NewFmlEquiv(fvs, fqe)), opt)
						im := g.QE(NewForAll(lvs, NewFmlImpl(fvs, fqe)), opt)
						re := g.QE(NewForAll(lvs, NewFmlImpl(fqe, fvs)), opt)
						t.Errorf("%v: %v, im=%v, re=%v", lvs, vv, im, re)
					}
					return
				}
			}
		}
	}
}

func TestVsQuad(t *testing.T) {
	funcname := "TestVsQuad"
	print_log := false
	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip %s... (no cas)\n", funcname)
		return
	}
	defer g.Close()

	optbl := []string{"<", ">", ">=", "!=", "==", "<="}

	for i, ss := range []string{
		"a*x^2-b*x+1 %s 0 && c %s 3",
		"a*x^2-3*b*x-5*c %s 0 && x %s 0",
		"7*x^2-6*b*x+9*b^2 %s 0 && a %s 0",
	} {
		for j, vars := range []struct {
			v  string // 変数順序
			lv Level  // 消去する変数のレベル
		}{
			{"x,a,b,c", 0},
			{"a,b,c,x", 3},
			{"a,b,x,c", 2},
			{"a,c,x,b", 2},
			{"c,b,a,x", 3},
		} {
			vstr := fmt.Sprintf("vars(%s);", vars.v)
			_, err := g.Eval(strings.NewReader(vstr))
			if err != nil {
				t.Errorf("[%d] %s failed: %s", i, vstr, err)
				return
			}
			for ops := 0; ops < 6*6; ops++ {
				op1 := optbl[ops/6]
				op2 := optbl[ops%6]
				input := fmt.Sprintf("ex([x], %s);", fmt.Sprintf(ss, op1, op2))
				if print_log {
					fmt.Printf(" == [%d,%d,%d,%v,%v] ==== %s ==[%s:%d]======================\n", i, j, ops, op1, op2, input, vars.v, vars.lv)
				}
				_fof, err := g.Eval(strings.NewReader(input))
				if err != nil {
					t.Errorf("[%d] parse error: %s, %s, %s\nerr=%s\ninput=%s", i, vstr, op1, op2, err, input)
					return
				}
				fof, ok := _fof.(Fof)
				if !ok {
					t.Errorf("[%d] eval failed\ninput=%s\neval=%v", i, ss, _fof)
					return
				}
				if err := ValidFof(fof); err != nil {
					t.Errorf("[%d] invalid fvs\ninput=%s\nactual=%v", i, ss, fof)
					return
				}

				fvs := QeVS(fof, vars.lv, 2, g)
				if print_log {
					fmt.Printf(">>>>>>>>>>>>\n")
				}
				if err := ValidFof(fvs); err != nil {
					t.Errorf("[%d] invalid fvs\ninput=%s\nactual=%v", i, ss, fvs)
					return
				}
				if HasVar(fvs, vars.lv) {
					t.Errorf("[%d,%d,%s] vs2 failed: has `%s`[%d]\ninput=%s\nactual=%s", i, j, vars.v, VarStr(vars.lv), vars.lv, input, fvs)
					return
				}

				// VS2次を使わないモードで
				opt := NewQEopt()
				opt.SetAlgo(QEALGO_VSQUAD, false)
				fqe := g.QE(fof, opt)
				if print_log {
					fmt.Printf(">>>>>>>>>>>>EQUIV\n")
				}

				eq := g.QE(NewForAll([]Level{0, 1, 2, 3}, NewFmlEquiv(fvs, fqe)), opt)
				if eq != TrueObj {
					t.Errorf("[%d,%d,%d,%v,%v,'%s'] invalid\ninputt=%v\nexpect=%v\nactual=%v\nactua2=%v\nequall=%v", i, j, ops, op1, op2, vars.v, fof, fqe, fvs, g.SimplFof(fvs), eq)

					opt.SetAlgo(QEALGO_VSLIN, false)
					for _, lvs := range [][]Level{
						{},
						{0, 1},
						{0, 2},
						{0, 3},
						{1, 2},
						{1, 3},
						{2, 3},
					} {
						if len(lvs) == 0 {
							break
						}
						vv := g.QE(NewForAll(lvs, NewFmlEquiv(fvs, fqe)), opt)
						im := g.QE(NewForAll(lvs, NewFmlImpl(fvs, fqe)), opt)
						re := g.QE(NewForAll(lvs, NewFmlImpl(fqe, fvs)), opt)
						t.Errorf("%v: %v, im=%v, re=%v", lvs, vv, im, re)
					}
					return
				}
			}
		}
	}
}
