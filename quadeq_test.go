package ganrac_test

import (
	"fmt"
	. "github.com/hiwane/ganrac"
	"strings"
	"testing"
)

func TestLinEq(t *testing.T) {
	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip TestLinEq... (no cas)\n")
		return
	}
	defer g.Close()

	p1 := NewPolyCoef(3, -3, NewPolyCoef(2, 0, 1))                   // z*w == 3
	p2 := NewPolyCoef(3, NewPolyCoef(1, 0, 1), NewPolyCoef(0, 0, 1)) // x*w+y
	z := NewPolyCoef(2, 0, 1)

	tbl := NewFofQuadEq(g, p1, 3)

	opt := NewQEopt()
	opt.Algo &= ^(QEALGO_EQLIN | QEALGO_EQQUAD)

	if (opt.Algo & (QEALGO_EQLIN | QEALGO_EQQUAD)) != 0 {
		t.Errorf("algo=%x", opt.Algo)
		return
	}

	for ii, ss := range []struct {
		op     OP
		expect Fof
	}{
		// ex([w], z*w==3 && x*w+y op 0)
		{
			EQ,
			NewFmlAnds(NewAtom(NewPolyCoef(2, 0, 1), NE), NewAtom(NewPolyCoef(2, NewPolyCoef(0, 0, 3), NewPolyCoef(1, 0, 1)), EQ)),
		}, {
			NE,
			NewFmlAnds(NewAtom(NewPolyCoef(2, 0, 1), NE), NewAtom(NewPolyCoef(2, NewPolyCoef(0, 0, 3), NewPolyCoef(1, 0, 1)), NE)),
		}, {
			LT,
			NewAtom(NewPolyCoef(2, 0, NewPolyCoef(0, 0, 3), NewPolyCoef(1, 0, 1)), LT),
		}, {
			GT,
			NewAtom(NewPolyCoef(2, 0, NewPolyCoef(0, 0, 3), NewPolyCoef(1, 0, 1)), GT),
		},
	} {
		a := NewAtom(p2, ss.op)
		tbl.SetSgnLcp(1)
		opos := NewFmlAnd(QeLinEq(a, tbl), NewAtom(z, GT))

		tbl.SetSgnLcp(-1)
		oneg := NewFmlAnd(QeLinEq(a, tbl), NewAtom(z, LT))

		o := NewFmlOr(opos, oneg)

		fof := NewQuantifier(true, []Level{0, 1, 2}, NewFmlEquiv(o, ss.expect))
		cad, _ := FuncCAD(g, "cad", []interface{}{fof})
		switch cmp := cad.(type) {
		case *AtomT:
			break
		default:
			t.Errorf("ii=%d, op=%d:`%s`\nexpect= %v.\nactual= %v OR %v.\ncmp=%v", ii, ss.op, ss.op, ss.expect, opos, oneg, cmp)
			return
		}

		switch cmp := g.QE(fof, opt).(type) {
		case *AtomT:
			continue
		default:
			t.Errorf("ii=%d, op=%d\nexpect= %v.\nactual= %v OR %v.\ncmp=%v", ii, ss.op, ss.expect, opos, oneg, cmp)
			return
		}
	}
}

func TestQuadEq1(t *testing.T) {
	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip TestQuadEq1... (no cas)\n")
		return
	}
	defer g.Close()

	z := NewPolyCoef(2, 0, 1)                         // 主係数
	p1 := NewPolyCoef(3, -5, NewPolyCoef(1, 0, 1), z) // z*w^2+y*w-5
	p2 := NewPolyCoef(3, -3, NewPolyCoef(0, 0, 1))    // x*w-3;

	tbl := NewFofQuadEq(g, p1, 3)

	opt := NewQEopt()
	opt.Algo &= ^(QEALGO_EQLIN | QEALGO_EQQUAD)

	if (opt.Algo & (QEALGO_EQLIN | QEALGO_EQQUAD)) != 0 {
		t.Errorf("algo=%x", opt.Algo)
		return
	}

	// discrim(p1)=y^2+20z >= 0: necessary condition
	d := NewPolyCoef(2, NewPolyCoef(1, 0, 0, 1), 20)
	dge := NewAtom(d, GE)
	// res(p1. p2)= -5*x^2+3*y*x+9*z
	r := NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, -5), NewPolyCoef(0, 0, 3)), 9)

	for ii, ss := range []struct {
		op     OP
		expect Fof
	}{
		// ex([w], p1 = 0 && p2 op 0)
		{EQ,
			NewFmlAnds(
				NewAtom(NewPolyCoef(0, 0, 1), NE),
				NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -5), 3), NE),
				NewAtom(NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 0, -5), NewPolyCoef(0, 0, 3)), 9), EQ)),
		}, {NE, // 1
			NewFmlAnds(NewAtom(z, NE), NewAtom(d, GE),
				NewFmlOr(NewAtom(d, GT), NewAtom(NewPolyCoef(1, NewPolyCoef(0, 0, -10), 3), NE))),
		}, {GT, // 2
			NewFmlAnd(dge, NewFmlOrs(
				NewFmlAnds(
					NewAtom(z, LT),
					NewAtom(NewPolyCoef(1, 0, NewPolyCoef(0, 0, -10), 3), LT)),
				NewAtom(Mul(r, z), LT))),
		}, {LE, // 3
			NewFmlAnds(dge, NewAtom(z, NE),
				NewFmlOrs(
					NewAtom(z, GT),
					NewAtom(r, GE),
					NewAtom(NewPolyCoef(1, 0, NewPolyCoef(0, 0, -10), 3), GT))),
		}, {LT, // 4
			NewFmlAnds(dge, NewAtom(z, NE), NewFmlOrs(
				NewAtom(z, GT),
				NewAtom(r, GT),
				NewAtom(NewPolyCoef(1, 0, NewPolyCoef(0, 0, -10), 3), GT))),
		}, {GE, // 5
			NewFmlAnd(dge, NewFmlOrs(
				NewFmlAnds(
					NewAtom(z, LT),
					NewAtom(NewPolyCoef(1, 0, NewPolyCoef(0, 0, -10), 3), LT)),
				NewFmlAnd(
					NewAtom(z, NE), NewAtom(Mul(r, z), LE)))),
		},
	} {
		a := NewAtom(p2, ss.op)
		var o Fof = FalseObj
		for _, sgns := range [][]int{
			{+1, +1},
			{+1, -1},
			{-1, +1},
			{-1, -1}} {
			tbl.SetSgnLcp(sgns[0])
			tbl.SetSgnS(sgns[1])
			op := GT
			if sgns[0] < 0 {
				op = LT
			}
			opp := NewFmlAnd(QeQuadEq(a, tbl), NewAtom(z, op))
			fmt.Printf("<%d,%2d,%2d> %v\n", ii, sgns[0], sgns[1], opp)
			o = NewFmlOr(o, opp)
		}

		fof := NewQuantifier(true, []Level{0, 1, 2}, NewFmlEquiv(NewFmlAnd(o, dge), ss.expect))
		switch cmp := g.QE(fof, opt).(type) {
		case *AtomT:
			continue
		default:
			t.Errorf("ii=%d, op=%d\ninput=ex([%v], (%v != 0) && %v == 0 && %v %s 0).\nexpect= %v.\nactual= (%v)\n   AND  (%v).\ncmp=%v\nfof=%v", ii, ss.op, NewPolyVar(3),
				z, p1, p2, ss.op, ss.expect, o, dge, cmp, fof)
			return
		}
	}
}

func TestQuadEq2(t *testing.T) {
	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip TestQuadEq2... (no cas)\n")
		return
	}
	defer g.Close()

	z := NewPolyCoef(2, 0, 1)                                             // 主係数
	p1 := NewPolyCoef(3, -3, -2, z)                                       // z*w^2-2*w-3
	p2 := NewPolyCoef(3, NewPolyCoef(1, -3, -1), 0, NewPolyCoef(0, 0, 1)) // x*w^2-y-3

	tbl := NewFofQuadEq(g, p1, 3)

	// (y^2+6*y+9)*z^2+(-6*x*y-18*x)*z-4*x*y+9*x^2-12*x
	r := NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, -12, 9), NewPolyCoef(0, 0, -4)), NewPolyCoef(1, NewPolyCoef(0, 0, -18), NewPolyCoef(0, 0, -6)), NewPolyCoef(1, 9, 6, 1))
	// 3z+1
	d := NewPolyCoef(2, 1, 3)
	dge := NewAtom(d, GE)

	opt := NewQEopt()
	opt.Algo &= ^(QEALGO_EQLIN | QEALGO_EQQUAD)

	if (opt.Algo & (QEALGO_EQLIN | QEALGO_EQQUAD)) != 0 {
		t.Errorf("algo=%x", opt.Algo)
		return
	}

	for ii, ss := range []struct {
		op     OP
		expect Fof
	}{
		// ex([w], z*w^2-2*w==3 && x*w+y op 3)
		{EQ,
			NewFmlAnds(dge, NewAtom(z, NE), NewAtom(r, EQ)),
		}, {GT,
			NewFmlAnds(dge, NewAtom(z, NE),
				NewFmlOrs(NewAtom(r, LT),
					NewFmlAnd(
						NewAtom(z, GT),
						NewAtom(NewPolyCoef(2, NewPolyCoef(0, 0, -3), NewPolyCoef(1, 3, 1)), LT)),
					NewFmlAnd(
						NewAtom(z, LT),
						NewAtom(NewPolyCoef(1, NewPolyCoef(0, 3, -9), 1), LT)))),
		},
	} {
		a := NewAtom(p2, ss.op)

		var o Fof = FalseObj
		for _, sgns := range [][]int{
			{+1, +1},
			{+1, -1},
			{-1, +1},
			{-1, -1}} {
			tbl.SetSgnLcp(sgns[0])
			tbl.SetSgnS(sgns[1])
			op := GT
			if sgns[0] < 0 {
				op = LT
			}
			opp := NewFmlAnd(QeQuadEq(a, tbl), NewAtom(z, op))
			o = NewFmlOr(o, opp)
		}

		fof := NewQuantifier(true, []Level{0, 1, 2}, NewFmlEquiv(NewFmlAnd(o, dge), ss.expect))
		switch cmp := g.QE(fof, opt).(type) {
		case *AtomT:
			continue
		default:
			fff := g.SimplFof(NewFmlAnd(o, dge))
			t.Errorf("ii=%d, op=%d\ninput =(%v == 0) && %v\nexpect= %v.\nactual= (%v) AND %v.\n      =%v\ncmp=%v", ii, ss.op, p1, a, ss.expect, o, dge, fff, cmp)
			return
		}
	}
}

func TestQuadEq3(t *testing.T) {
	funcname := "TestQuadEq3"
	print_log := false
	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip %s... (no cas)\n", funcname)
		return
	}
	defer g.Close()

	qeopt := NewQEopt()
	qeopt.SetG(g)

	qecond := NewQeCond()
	optbl := []OP{LE, GE, EQ, NE, GT, LT}

	for i, ss := range []struct {
		// ex([x], eq = 0 && cond op 0 && fof)
		eq   string
		cond string
		fof  string
	}{
		{
			"x^2-x-5",
			"b*x + 3",
			"true",
		}, {
			"x^2-10*x+1",
			"b*x + 3",
			"true",
		}, {
			"a*x^2-5*x-2",
			"b*x + 1",
			"true",
		}, {
			"x^2-a*x-5",
			"b*x^2 -5*x-3",
			"true",
		}, {
			"x^2-a*x-b",
			"b*x + 1",
			"true",
		},
	} {
		for j, vars := range []struct {
			v  string // 変数順序
			lv Level  // 消去する変数のレベル
		}{
			{"x,a,b,c", 0},
			{"a,c,x,b", 2},
			{"a,b,x,c", 2},
			{"a,b,c,x", 3},
			{"c,b,a,x", 3},
		} {
			vstr := fmt.Sprintf("vars(%s);", vars.v)
			_, err := g.Eval(strings.NewReader(vstr))
			if err != nil {
				t.Errorf("[%d] %s failed: %s", i, vstr, err)
				return
			}

			eqp, err := str2poly(g, ss.eq)
			if err != nil {
				t.Errorf("[%d,%d,%s] eval(eq) failed: %s, %s", i, j, vstr, err, ss.eq)
				return
			}
			condp, err := str2poly(g, ss.cond)
			if err != nil {
				t.Errorf("[%d,%d,%s] eval(cond) failed: %s, %s", i, j, vstr, err, ss.cond)
				return
			}
			fof, err := str2fof(g, ss.fof)
			if err != nil {
				t.Errorf("[%d,%d,%s] eval(fof) failed: %s, %s", i, j, vstr, err, ss.fof)
				return
			}

			fofeq := NewFmlAnd(NewAtom(eqp, EQ), fof)
			for ops := 0; ops < 6; ops++ {
				var cond Fof

				var op OP
				if ops < 0 {
					cond = TrueObj
					op = OP_TRUE
				} else {
					op = optbl[ops%6]
					cond = NewAtom(condp, op)
				}

				input := NewExists([]Level{vars.lv}, NewFmlAnd(cond, fofeq)).(FofQ)
				if print_log {
					fmt.Printf(" == [%d,%d,%s] ==== [%s== 0 && %s %v 0 && %s]======================\n", i, j, vstr, ss.eq, ss.cond, op, ss.fof)
				}

				qeopt.SetAlgo(QEALGO_EQLIN|QEALGO_EQQUAD, true)
				fquadeq := QeOptQuadEq(input, qeopt, qecond)
				if print_log {
					fmt.Printf(">>>>>>>>>>>> %v\n", fquadeq)
				}
				if fquadeq == nil {
					t.Errorf("[%d,%d,%s,%v] invalid fvs\ninput=%s\nactual=%v", i, j, vstr, op, ss, fquadeq)
					break
				}
				if err := ValidFof(fquadeq); err != nil {
					t.Errorf("[%d,%d,%s,%v] invalid fvs\ninput=%s\nactual=%v\nactual=%V", i, j, vstr, op, ss, fquadeq, fquadeq)
					return
				}
				if HasVar(fquadeq, vars.lv) {
					if false {
						t.Errorf("[%d,%d,%s,%v] invalid fquadeq. has `%s`\ninput=%s\nactual=%v\nactual=%V", i, j, vstr, op, VarStr(vars.lv), ss, fquadeq, fquadeq)
						return
					} else {
						//						fmt.Printf("eqq=%v\n", fquadeq)
						fquadeq = g.QE(NewExists([]Level{vars.lv}, fquadeq), qeopt)
					}
				}

				qeopt.SetAlgo(QEALGO_EQLIN|QEALGO_EQQUAD, false)
				fqe := g.QE(input, qeopt)
				if print_log {
					fmt.Printf(">>>>>>>>>>>>EQUIV\n")
				}

				eq := g.QE(NewForAll([]Level{0, 1, 2, 3}, NewFmlEquiv(fquadeq, fqe)), qeopt)
				if eq != TrueObj {
					t.Errorf("[%d,%d,%s,%v] invalid\ninput=%s\nexpect=%v\nactual=%v", i, j, vstr, op, input, fqe, fquadeq)

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
						vv := g.QE(NewForAll(lvs, NewFmlEquiv(fquadeq, fqe)), qeopt)
						im := g.QE(NewForAll(lvs, NewFmlImpl(fquadeq, fqe)), qeopt)
						re := g.QE(NewForAll(lvs, NewFmlImpl(fqe, fquadeq)), qeopt)
						t.Errorf("%v: %v, im=%v, re=%v", lvs, vv, im, re)
					}

					dict := NewDict()
					dict.Set("var", NewInt(1))
					impl, _ := FuncCAD(g, "CAD", []interface{}{
						NewFmlImpl(fquadeq, fqe), dict,
					})
					repl, _ := FuncCAD(g, "CAD", []interface{}{
						NewFmlImpl(fqe, fquadeq), dict,
					})

					smpl, _ := FuncCAD(g, "CAD", []interface{}{
						fquadeq, dict,
					})

					t.Errorf("\nactual=%v\n  impl=%v\n  repl=%v", smpl, impl, repl)
					return
				}
			}
		}
	}
}
