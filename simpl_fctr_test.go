package ganrac_test

import (
	"fmt"
	. "github.com/hiwane/ganrac"
	"testing"
)

func testSameFormAndOrFctr(output, expect Fof) bool {
	// 形として同じか.
	// QE によるチェックでは等価性は確認できるが簡単化はわからない
	if FofTag(output) != FofTag(expect) {
		return false
	}
	if !output.IsQff() {
		return false
	}
	if !expect.IsQff() {
		return false
	}

	var oofmls, eefmls []Fof
	switch oo := output.(type) {
	case *FmlAnd:
		NormalizeFof(oo)
		oofmls = oo.Fmls()
		ee := expect.(*FmlAnd)
		NormalizeFof(ee)
		eefmls = ee.Fmls()
	case *FmlOr:
		NormalizeFof(oo)
		oofmls = oo.Fmls()
		ee := expect.(*FmlOr)
		NormalizeFof(ee)
		eefmls = ee.Fmls()
	case *Atom:
		if oo.Equals(expect) {
			return true
		}
		NormalizeFof(oo)
		NormalizeFof(expect)
		return oo.Equals(expect)
	default:
		return oo.Equals(expect)
	}

	if len(oofmls) != len(eefmls) {
		fmt.Printf("len %d, %d\n", len(oofmls), len(eefmls))
		return false
	}

	for i := 0; i < len(oofmls); i++ {
		vars := make([]bool, 5)
		oofmls[i].Indets(vars)
		lv := Level(0)
		for j := 0; j < len(vars); j++ {
			if vars[j] {
				lv = Level(j)
			}
		}

		m := 0
		for j := 0; j < len(eefmls); j++ {
			vars = make([]bool, 5)
			eefmls[j].Indets(vars)
			if vars[lv] {
				m += 1
				if !oofmls[i].Equals(eefmls[j]) {
					return false
				}
			}
		}
		if m != 1 {
			return false
		}
	}
	return true
}

func TestSimplFctr(t *testing.T) {

	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip TestSimplFctr... (no cas)\n")
		return
	}
	defer g.Close()

	x := NewPolyVar(0)
	y := NewPolyVar(1)
	z := NewPolyVar(2)

	for i, s := range []struct {
		input  Fof
		expect Fof
	}{
		{
			NewAtom(x.Powi(3), EQ),
			NewAtom(x, EQ),
		}, {
			NewAtom(x.Powi(2).Mul(NewInt(5)), LE),
			NewAtom(x, EQ),
		}, {
			NewAtom(x.Powi(2).Mul(NewInt(5)), LT),
			FalseObj,
		}, {
			NewAtom(x.Powi(2).Mul(NewInt(5)), GE),
			TrueObj,
		}, {
			NewAtom(x.Powi(2).Mul(NewInt(5)), GT),
			NewAtom(x, NE),
		}, {
			NewAtom(x.Powi(2).Mul(NewInt(5)), NE),
			NewAtom(x, NE),
		}, {
			NewAtom(x.Powi(2).Mul(NewInt(5)), EQ),
			NewAtom(x, EQ),
		}, {
			NewAtom(x.Powi(3).Mul(NewInt(5)), LE),
			NewAtom(x, LE),
		}, {
			NewAtom(x.Powi(3).Mul(NewInt(5)), LT),
			NewAtom(x, LT),
		}, {
			NewAtom(x.Powi(3).Mul(NewInt(5)), GE),
			NewAtom(x, GE),
		}, {
			NewAtom(x.Powi(3).Mul(NewInt(5)), GT),
			NewAtom(x, GT),
		}, {
			NewAtom(x.Powi(3).Mul(NewInt(5)), NE),
			NewAtom(x, NE),
		}, {
			NewAtom(x.Powi(3).Mul(NewInt(5)), EQ),
			NewAtom(x, EQ),
		}, { // 13
			NewAtom(x.Powi(3).Mul(y.Powi(4)), LE),
			NewFmlOr(NewAtom(y, EQ), NewAtom(x, LE)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)), GE),
			NewFmlOr(NewAtom(y, EQ), NewAtom(x, GE)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)), LT),
			NewFmlAnd(NewAtom(y, NE), NewAtom(x, LT)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)), GT),
			NewFmlAnd(NewAtom(y, NE), NewAtom(x, GT)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)), EQ),
			NewFmlOr(NewAtom(y, EQ), NewAtom(x, EQ)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)), NE),
			NewFmlAnd(NewAtom(y, NE), NewAtom(x, NE)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)).Mul(z).Mul(NewInt(5)), LE),
			NewFmlOr(NewAtom(y, EQ), NewAtoms([]RObj{x, z}, LE)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)).Mul(z).Mul(NewInt(5)), GE),
			NewFmlOr(NewAtom(y, EQ), NewAtoms([]RObj{x, z}, GE)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)).Mul(z).Mul(NewInt(5)), LT),
			NewFmlAnd(NewAtom(y, NE), NewAtoms([]RObj{x, z}, LT)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)).Mul(z).Mul(NewInt(5)), GT),
			NewFmlAnd(NewAtom(y, NE), NewAtoms([]RObj{x, z}, GT)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)).Mul(z).Mul(NewInt(5)), EQ),
			NewFmlOrs(NewAtom(y, EQ), NewAtom(z, EQ), NewAtom(x, EQ)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)).Mul(z).Mul(NewInt(5)), NE),
			NewFmlAnds(NewAtom(y, NE), NewAtom(z, NE), NewAtom(x, NE)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)).Mul(z.Powi(4)).Mul(NewInt(5)), LE),
			NewFmlOrs(NewAtom(y, EQ), NewAtom(z, EQ), NewAtom(x, LE)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)).Mul(z.Powi(4)).Mul(NewInt(5)), GE),
			NewFmlOrs(NewAtom(y, EQ), NewAtom(z, EQ), NewAtom(x, GE)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)).Mul(z.Powi(4)).Mul(NewInt(5)), LT),
			NewFmlAnds(NewAtom(y, NE), NewAtom(z, NE), NewAtom(x, LT)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)).Mul(z.Powi(4)).Mul(NewInt(5)), GT),
			NewFmlAnds(NewAtom(y, NE), NewAtom(z, NE), NewAtom(x, GT)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)).Mul(z.Powi(4)).Mul(NewInt(5)), EQ),
			NewFmlOrs(NewAtom(y, EQ), NewAtom(z, EQ), NewAtom(x, EQ)),
		}, {
			NewAtom(x.Powi(3).Mul(y.Powi(4)).Mul(z.Powi(4)).Mul(NewInt(5)), NE),
			NewFmlAnds(NewAtom(y, NE), NewAtom(z, NE), NewAtom(x, NE)),
		},
	} {
		var output Fof
		output = SimplFctr(s.input, g)

		if !testSameFormAndOrFctr(output, s.expect) {
			fmt.Printf("expect %V\n", s.expect)
			fmt.Printf("output %V\n", output)
			t.Errorf("i=%d\ninput =%v\nexpect=%v\nactual=%v", i, s.input, s.expect, output)
			continue
		}
	}
}
