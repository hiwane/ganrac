package ganrac_test

import (
	"fmt"
	. "github.com/hiwane/ganrac"
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
