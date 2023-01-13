package ganrac_test

import (
	"fmt"
	"math/big"
	"testing"

	. "github.com/hiwane/ganrac"
)

// func TestSymSqfr2(t *testing.T) {
// 	p := NewPolyCoef(2,
// 			NewPolyCoef(1, ParseInt("-4722201209543931678482839880788832023380506408746984383", 10), ParseInt("7037081672937279457168302335815099996936095844515446784", 10), ParseInt("-7370402196928139846235864997969418656470294922939858944", 10)),
// 			NewPolyCoef(1, ParseInt("-937089791168204322416009923114557457813156603499642880", 10), ParseInt("695112178837912870216157370077032489892039416212357120", 10), ParseInt("-59907021258380514474371077442594431218361575250329600", 10)),
// 			NewPolyCoef(1, ParseInt("-46489827399607932944452513320345263121254652785459200", 10), ParseInt("1412477194662381313033576520811520736137630973952000", 10), ParseInt("30100594890350423012826167438885134580526627998924800", 10)))
// 	d := NewPolyCoef(1,
// 			ParseInt("189708409272553183839474655", 10),
// 			ParseInt("-1436744321078070999056602079", 10),
// 			ParseInt("-5184691718935390049781415936", 10),
// 			ParseInt("5196819548962240103538753536", 10))
// }

func TestSymSqfr(t *testing.T) {
	// ox は必要ないのだけど．
	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip TestSymSqfr... (no cas)\n")
		return
	}
	defer g.Close()
	one := NewInt(1)

	fof := NewQuantifier(false, []Level{3}, NewAtom(NewPolyCoef(3, NewPolyCoef(2, NewPolyCoef(1, NewPolyCoef(0, 0, 1), 1), 1), 1), GT))
	cad, err := NewCAD(fof, g)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	cad.InitProj(0)
	cell := NewCell(cad, cad.Root(), 1)
	cell.SetDefPoly(NewPolyCoef(0, -2, 0, 1)) // x^2-2

	// all([w], (y^2-x)*w^4+(z^2+2*x*z+(-x+1)*y^2-2*x*y-x)*w^2+(-x+1)*z^2<=0);
	// lift(12, 7, 1)
	cell_adam3y := NewCell(cad, cad.Root(), 7)
	cell_adam3y.SetDefPoly(NewPolyCoef(0, 5, 10, 4))
	cell_adam3y.SetIntv(
		NewBinInt(big.NewInt(-741937353), -30),
		NewBinInt(big.NewInt(-741937352), -30))

	cell_adam3z := NewCell(cad, cell_adam3y, 1)
	cell_adam3z.SetDefPoly(NewPolyCoef(1, 0, 0, NewPolyCoef(0, 0, -40), 20, 1))
	nintv := NewIntervalInt64(3, 53)
	nintv.Inf().SetFloat64(-18.506509)
	nintv.Sup().SetFloat64(-18.506507)
	cell_adam3z.SetNintv(nintv)

	for ii, s := range []struct {
		p    *Poly
		r    []Mult_t
		d    []int
		cell *Cell
	}{
		{ // 0
			NewPolyCoef(3, -3, 7, -5, 1),
			[]Mult_t{1, 2},
			[]int{1, 1},
			cell,
		}, { // 1
			NewPolyCoef(3, -5, 15, -16, 8, -3, 1),
			[]Mult_t{1, 3},
			[]int{2, 1},
			cell,
		}, { // 2
			NewPolyCoef(3, 2, NewPolyCoef(0, 0, -2), 1),
			[]Mult_t{2},
			[]int{1},
			cell,
		}, { // 3
			NewPolyCoef(3, NewPolyCoef(0, 0, -2), 6, NewPolyCoef(0, 0, -3), 1),
			[]Mult_t{3},
			[]int{1},
			cell,
		}, { // 4
			// x^4-a*x^3-3*a^2*x^2+5*a^3*x-2*a^4
			NewPolyCoef(1, -8, NewPolyCoef(0, 0, 0, 0, 5), NewPolyCoef(0, 0, 0, -3), NewPolyCoef(0, 0, -1), 1),
			[]Mult_t{1, 3},
			[]int{1, 1},
			cell,
		}, { // 5
			NewPolyCoef(3, NewPolyCoef(0, 0, 0, 0, 0, -1), NewPolyCoef(0, 0, 10), -12, NewPolyCoef(0, 0, -4), 8),
			[]Mult_t{1, 3},
			[]int{1, 1},
			cell,
		}, { // 6
			NewPolyCoef(2, NewPolyCoef(1, 0, 0, -4), 0, NewPolyCoef(1, NewPolyCoef(0, -5, -10, -4), 10, 1), 0, NewPolyCoef(0, -5, 0, 1)),
			[]Mult_t{2},
			[]int{2},
			cell_adam3z,
		},
	} {
		for jj, ppp := range []*Poly{
			s.p, s.p.Neg().(*Poly),
			s.p.SubsXinv(), s.p.SubsXinv().Neg().(*Poly)} {

			// fmt.Printf("TestSymSqfr(%d,%d) ppp=%v ===========================\n", ii, jj, ppp)

			sqfr := cad.Sym_sqfr2(ppp, s.cell)
			var q RObj = one
			hasErr := false
			for i, sq := range sqfr {
				//fmt.Printf("<%d> [%v]^%d\n", i, sq.p, sq.r)
				if Mult_t(sq.Multi()) != s.r[i] {
					t.Errorf("<%d,%d,%d,r>\nexpect=%v\nactual=%d\nret=%v", ii, jj, i, s.r[i], sq.Multi(), sq.Poly())
					for i, sqx := range sqfr {
						t.Errorf("<%d>: (%v)^(%d)", i, sqx.Poly(), sqx.Multi())
					}
					hasErr = true
					return
				}
				if sq.Poly().Deg(sq.Poly().Level()) != s.d[i] {
					t.Errorf("<%d,%d,%d,degree>\nexpect=%v\nactual=%d\nret=%v", ii, jj, i, s.d[i], sq.Multi(), sq.Poly())

					for i, sqx := range sqfr {
						t.Errorf("<%d>: (%v)^(%d)", i, sqx.Poly(), sqx.Multi())
					}

					hasErr = true
					break
				}

				qq := sq.Poly().Pow(NewInt(int64(sq.Multi())))
				q = Mul(qq, q)
			}
			if hasErr {
				break
			}
			qq := Sub(Mul(q, ppp.LC()), Mul(ppp, q.(*Poly).LC()))
			if !qq.IsZero() {
				if qx, ok := qq.(*Poly); ok {
					flag := true
					if qx.Level() == ppp.Level() {
						lv := qx.Level()
						for _i := uint(0); _i <= uint(qx.Deg(lv)); _i++ {
							qc := qx.Coef(lv, _i)

							switch qcc := qc.(type) {
							case *Poly:
								if !cad.Sym_zero_chk(qcc, s.cell) {
									flag = false
								}
							case NObj:
								if !qcc.IsZero() {
									flag = false
								}
							}
						}
						if flag {
							continue
						}
					} else if cad.Sym_zero_chk(qx, s.cell) {
						continue
					}
				}
				t.Errorf("<%d>\ninput =%v\noutput=%v\nret=%v\nqq=%v", ii, ppp, q, sqfr, qq)
			}
		}
	}
}

func testNewCADCell() (*CAD, *Cell, *Cell) {
	cad := new(CAD)
	cad.InitRoot()
	cell0 := NewCell(cad, nil, 1) // projection していないため nil 指定
	cell0.SetLevel(0)
	cell0.SetParent(cad.Root())
	cell0.SetDefPoly(NewPolyCoef(0, -2, 0, 1)) // x^2-2
	cell1 := NewCell(cad, nil, 1)              // projection していないため nil 指定
	cell1.SetLevel(1)
	cell1.SetParent(cell0)
	cell1.SetDefPoly(NewPolyCoef(1, -1, -2, 1)) // y^2-2*y-1

	return cad, cell0, cell1
}

func TestSymGcdMod(t *testing.T) {

	cad, cell0, cell1 := testNewCADCell()

	for ii, s := range []struct {
		f, g   *Poly
		expect *Poly
		celf   bool
		cell   *Cell
		p      Uint
	}{
		{
			NewPolyCoef(2, -5, 0, 1),
			NewPolyCoef(2, NewPolyVar(0), -1),
			nil, false,
			cell0, 151,
		}, {
			NewPolyCoef(2, 1, 0, -5),
			NewPolyCoef(2, NewPolyVar(0), -1),
			nil, false,
			cell0, 151,
		}, {
			NewPolyCoef(2, 1, 0, -5),
			NewPolyCoef(2, NewPolyVar(0), -1),
			nil, false,
			cell1, 151,
		}, {
			NewPolyCoef(2, -2, 0, 1),
			NewPolyCoef(2, NewPolyVar(0), -1),
			NewPolyCoef(2, NewPolyVar(0).Neg(), +1), false,
			cell0, 151,
		}, {
			NewPolyCoef(2, 5, NewPolyCoef(0, 1, 1), 1),
			NewPolyCoef(2, 5, NewPolyVar(1), 1),
			nil, true,
			cell0, 151,
		}, {
			NewPolyCoef(2, 5, NewPolyCoef(0, 1, 1), 3),
			NewPolyCoef(2, 5, NewPolyVar(1), 3),
			nil, true,
			cell0, 151,
		}, {
			NewPolyCoef(2, -2, 0, NewPolyCoef(0, 1, 1)),
			NewPolyCoef(2, -2, 0, NewPolyVar(1)),
			nil, true,
			cell0, 151,
		},
	} {
		fp := s.f.Mod(s.p).(*Poly)
		gp := s.g.Mod(s.p).(*Poly)

		// fmt.Printf("<%d>===TestSymGcdMod() ======================================\nf=%v\ng=%v\n", ii, s.f, s.g)
		cellp, ok := cell1.Mod(cad, s.p)
		if !ok {
			t.Errorf("not ok ii=%d", ii)
			continue
		}

		gcd, a, b := cad.Symde_gcd_mod(fp, gp, cellp, s.p, true)
		if (gcd == nil) != (s.expect == nil) {
			t.Errorf("invalid gcd <a1, %d, %d>\nf=%v, g=%v\nexpect=%v\nactual=%v\na=%v\nb=%v\n",
				ii, s.p, s.f, s.g, s.expect, gcd, a, b)
			continue
		}
		if gcd != nil {
			if gcd.Level() != fp.Level() || gcd.Deg(gcd.Level()) != s.expect.Deg(gcd.Level()) {
				t.Errorf("invalid gcd <a2, %d, %d>\nf=%v, g=%v\nexpect=%v\nactual=%v\na=%v\nb=%v\n",
					ii, s.p, s.f, s.g, s.expect, gcd, a, b)
				continue
			}

			dega := DegModer(a)
			degb := DegModer(b)

			if degb > DegModer(fp)-DegModer(gcd) || dega > DegModer(gp)-DegModer(gcd) {
				t.Errorf("invalid gcd <a3, %d, %d>: %d, %d\nf=%v -> %v\ng=%v -> %v\nexpect=%v\nactual=%v\na=%v\nb=%v\n",
					ii, s.p, dega, degb, s.f, fp, s.g, gp, s.expect, gcd, a, b)
				return
			}

			gg := Add_mod(Mul_mod(fp, a, s.p), Mul_mod(gp, b, s.p), s.p)
			if !gcd.Equals(gg) {
				t.Errorf("invalid gcd <a4, %d, %d>: %d, %d\nf=%v, g=%v\nfp=%v\ngp=%v\nexpect=%v\nactual=%v\ngg    =%v\nfpa=%v\ngpb=%v -> %v\na=%v\nb=%v\n",
					ii, s.p, dega, degb, s.f, s.g, fp, gp, s.expect, gcd, gg,
					Mul_mod(fp, a, s.p),
					Mul_mod(gp, b, s.p),
					gg,
					a, b)
				return
			}

		}
		if (a == nil) != s.celf {
			t.Errorf("invalid gcd <a4, %d, %d>\nf=%v, g=%v\nexpect=%v\nactual=%v\na=%v\nb=%v\n",
				ii, s.p, s.f, s.g, s.expect, gcd, a, b)
			continue
		}
		if a == nil {
			if cellp.Factor1() == nil {
				t.Errorf("invalid gcd <a5, %d, %d>\nf=%v, g=%v\nexpect=%v\nactual=%v\na=%v\nb=%v\n",
					ii, s.p, s.f, s.g, s.expect, gcd, a, b)
				continue
			}
			if cellp.Factor2() != nil && (cellp.Factor1().Level() != cellp.Factor2().Level() || DegModer(cellp.Factor1()) != DegModer(cellp.Factor2())) {
				t.Errorf("invalid gcd <a6, %d, %d>\nf=%v, g=%v\nexpect=%v\nactual=%v\na=%v\nb=%v\n",
					ii, s.p, s.f, s.g, s.expect, gcd, a, b)
				continue
			}
		}

		gcd3, a3, b3 := cad.Symde_gcd_mod(gp, fp, cellp, s.p, true)
		if gcd3 == nil {
			if gcd != nil {
				t.Errorf("invalid gcd <b1, %d, %d>\nf=%v, g=%v\nexpect=%v\nactual=%v :: %v\na=%v\nb=%v\n",
					ii, s.p, s.f, s.g, s.expect, gcd3, gcd, a3, b3)
				continue

			}

		} else if (gcd == nil) != (gcd3 != nil) && !gcd3.Equals(gcd) && !Add_mod(gcd3, gcd, s.p).IsZero() {
			t.Errorf("invalid gcd <b2, %d, %d>\nf=%v, g=%v\nexpect=%v\nactual=%v :: %v\na=%v\nb=%v\n",
				ii, s.p, s.f, s.g, s.expect, gcd3, gcd, a3, b3)
			continue
		}
		if a3 != nil && (b3 == nil || DegModer(a3) != DegModer(b) || DegModer(b3) != DegModer(a)) {
			t.Errorf("invalid gcd <b3, %d, %d>\nf=%v, g=%v\nexpect=%v\nactual=%v :: %v\na=%v :: %v\nb=%v :: %v\n",
				ii, s.p, s.f, s.g, s.expect, gcd3, gcd, a3, b, b3, a)
			continue

		}

		gcd2, a2, _ := cad.Symde_gcd_mod(fp, gp, cellp, s.p, false)
		if (gcd == nil) != (gcd2 == nil) {
			t.Errorf("invalid gcd <c1, %d, %d>\nf=%v, g=%v\nexpect=%v\nactual=%v\n",
				ii, s.p, s.f, s.g, s.expect, gcd)
			continue
		}
		if gcd != nil && !gcd.Equals(gcd2) {
			t.Errorf("invalid gcd <c2, %d, %d>\nf=%v, g=%v\nexpect=%v\nactual=%v\n",
				ii, s.p, s.f, s.g, s.expect, gcd)
			continue
		}
		if a != nil && !a.Equals(a2) {
			t.Errorf("invalid gcd <c3, %d, %d>\nf=%v, g=%v\nexpect=%v\nactual=%v\n",
				ii, s.p, s.f, s.g, s.expect, gcd)
			continue
		}
	}
}

func TestSymGcd(t *testing.T) {

	for ii, s := range []struct {
		f, g   *Poly
		expect *Poly
	}{
		{
			NewPolyCoef(2, NewPolyCoef(0, -2, 1), -1, 1), // z^2-z+x-2
			NewPolyCoef(2, NewPolyCoef(0, -2, 3), -3, 1), // z^2-3*z+3*x-2
			NewPolyCoef(2, NewPolyCoef(0, 0, -1), 1),
		},
	} {
		fmt.Printf("<%d>===TestSymGcd() ======================================\nf=%v\ng=%v\n", ii, s.f, s.g)
		cad, _, cell1 := testNewCADCell()
		gcd, _ := cad.Symde_gcd2(s.f, s.g, cell1, 0)

		if !gcd.Equals(s.expect) {
			t.Errorf("i=%d\nf  =%v\ng  =%v\nexp=%v\nact=%v", ii, s.f, s.g, s.expect, gcd)
		}
	}
}
