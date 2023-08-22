package ganrac_test

import (
	"fmt"
	"testing"

	. "github.com/hiwane/ganrac"
)

func TestSdcQEmain(t *testing.T) {
	funcname := "TestSdcQEmain"

	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip %s... (no cas)\n", funcname)
		return
	}
	defer g.Close()

	qeopt := NewQEopt()
	qeopt.SetG(g)

	lvs := []Level{0, 1, 2, 3, 4, 5}

	for ii, tt := range []struct {
		varstr string
		lv     Level // x のレベル
	}{
		{"x,a,b,c", 0},
		{"x,b,a,c", 0},
		{"a,x,b,c", 1},
		{"b,x,a,c", 1},
		{"a,b,x,c", 2},
		{"b,a,x,c", 2},
	} {
		vstr := fmt.Sprintf("vars(%s);", tt.varstr)
		_, err := evalstr(g, vstr)
		if err != nil {
			t.Errorf("[%d] `%s` failed: %s", ii, vstr, err)
			return
		}

		lv := tt.lv

		for jj, ss := range []struct {
			input  string
			expect string
		}{
			{
				"a*x^2-b*x+1",
				"(a==0 && b>0) || (b^2-4*a>=0 && (a<0 || (a!=0 && b>=0)))",
			}, {
				"3*x^2+a*x+b",
				"(12*b-a^2<=0 && a<=0) || b<=0",
			}, {
				"a*x^2-5*x+b",
				"4*a*b-25<=0 || a<=0",
			}, {
				"a*x^3+b*x^2-5*x+1",
				"a<0 || 4*b^3-25*b^2+90*a*b+27*a^2-500*a<=0",
			}, {
				"a*x^3+b*x^2-c*x+1",
				"a<0 || (c>=0 && 4*a*c^3+b^2*c^2-18*a*b*c-4*b^3-27*a^2>0) || (4*a*c^3+b^2*c^2-18*a*b*c-4*b^3-27*a^2==0 && c>0) || (4*a*c^3+b^2*c^2-18*a*b*c-4*b^3-27*a^2>=0 && b<0)",
			}, {
				"a*x^3+b*x^2-7*x+c",
				"a<0 || 27*a^2*c^2+(4*b^3+126*a*b)*c-49*b^2-1372*a<=0 || c<=0",
			}, {
				"a*x^3+b*x^2+7*x+c",
				"(27*a^2*c^2+(4*b^3-126*a*b)*c-49*b^2+1372*a==0 && b<0) || (b<=0 && 27*a^2*c^2+(4*b^3-126*a*b)*c-49*b^2+1372*a<0) || c<=0 || a<0",
			}, {
				"x^4+5*x^3+a*x^2+b*x+1",
				"(4*a-17<=0 && 27*b^4+(-90*a+500)*b^3+(4*a^3-25*a^2-144*a+150)*b^2+(400*a^2-2250*a+960)*b-16*a^4+100*a^3+128*a^2-3600*a+16619<=0) || (54*b-45*a+250<=0 && 27*b^4+(-90*a+500)*b^3+(4*a^3-25*a^2-144*a+150)*b^2+(400*a^2-2250*a+960)*b-16*a^4+100*a^3+128*a^2-3600*a+16619>=0 && 54*b^3+(-135*a+750)*b^2+(4*a^3-25*a^2-144*a+150)*b+200*a^2-1125*a+480<=0) ",
			}, {
				"x^5-3*x^4+5*x^3+a*x^2+b*x+1",
				"(256*b^5+(576*a+1093)*b^4+(666*a^2+1970*a+4072)*b^3+(-27*a^4-162*a^3-3335*a^2-7906*a+4600)*b^2+(-2934*a^3-13650*a^2+1680*a+46740)*b+108*a^5+648*a^4+3800*a^3+35457*a^2+138300*a+173792<=0 && 5832*a^5+94041*a^4+666036*a^3+2384694*a^2+4356236*a+3531833<0) || (256*b^5+(576*a+1093)*b^4+(666*a^2+1970*a+4072)*b^3+(-27*a^4-162*a^3-3335*a^2-7906*a+4600)*b^2+(-2934*a^3-13650*a^2+1680*a+46740)*b+108*a^5+648*a^4+3800*a^3+35457*a^2+138300*a+173792<=0 && 2560*b^3+(3456*a+6558)*b^2+(1998*a^2+5910*a+12216)*b-27*a^4-162*a^3-3335*a^2-7906*a+4600<=0 && 640*b^4+(1152*a+2186)*b^3+(999*a^2+2955*a+6108)*b^2+(-27*a^4-162*a^3-3335*a^2-7906*a+4600)*b-1467*a^3-6825*a^2+840*a+23370>=0)",
			}, {
				"x^6-3*x^4+5*x^3+a*x^2+b",
				"(46656*b^4+(186624*a-1248048)*b^3+(13824*a^3+155520*a^2-454896*a+381996)*b^2+(27648*a^4+98496*a^3-1722384*a^2-7946100*a-9568125)*b+1024*a^6-4608*a^5-38016*a^4-56700*a^3>0 && 16*b+16*a-107<=0) || b<=0 || (46656*b^4+(186624*a-1248048)*b^3+(13824*a^3+155520*a^2-454896*a+381996)*b^2+(27648*a^4+98496*a^3-1722384*a^2-7946100*a-9568125)*b+1024*a^6-4608*a^5-38016*a^4-56700*a^3>=0 && 6912*b^3+(20736*a-138672)*b^2+(1024*a^3+11520*a^2-33696*a+28296)*b+1024*a^4+3648*a^3-63792*a^2-294300*a-354375<0)",
			}, {
				"x^7-3*x^4+5*x^3+a*x^2+1",
				"12500*a^7+772500*a^5+5204399*a^4+16242800*a^3+151312576*a^2+663464448*a+913030144<=0",
			}, {
				"x^7+a*x^2+b*x+1",
				"(46656*b^7-381024*a*b^5+926100*a^2*b^3-3125*a^6*b^2-600250*a^3*b+12500*a^7+823543==0 && 432*b^4-1008*a*b^2+245*a^2!=0 && 3125*a^7-7411887<0) || (46656*b^7-381024*a*b^5+926100*a^2*b^3-3125*a^6*b^2-600250*a^3*b+12500*a^7+823543==0 && b<=0 && 3125*a^7-7411887>=0) || (46656*b^7-381024*a*b^5+926100*a^2*b^3-3125*a^6*b^2-600250*a^3*b+12500*a^7+823543<0 && b<=0) || (46656*b^7-381024*a*b^5+926100*a^2*b^3-3125*a^6*b^2-600250*a^3*b+12500*a^7+823543<=0 && a<=0)",
			},
		} {

			p, err := str2poly(g, ss.input)
			if err != nil {
				t.Errorf("[%d,%d] parse error: `%s`\nerr=%s\ninput=%s.", ii, jj, vstr, err, ss.input)
				return
			}

			expect, err := str2fof(g, ss.expect)
			if err != nil {
				t.Errorf("[%d,%d] parse error: %s\nerr=%s\nexpect=%s", ii, jj, vstr, err, ss.expect)
				return
			}

			qeopt.SetAlgo(QEALGO_SDC, true)
			ret := SdcQEmain(p, lv, TrueObj, *qeopt)
			if err := ValidFof(ret); err != nil {
				t.Errorf("[%d,%d] ret invalid\ninput=%s\nerr=%v", ii, jj, ss.input, err)
				return
			}

			c0 := p.Coef(lv, 0)

			ret_c0le := NewFmlOr(ret, NewAtom(c0, LE))
			v := NewForAll(lvs, NewFmlEquiv(ret_c0le, expect))
			qeopt.SetAlgo(QEALGO_SDC|QEALGO_ATOM, false)
			v = g.QE(v, qeopt)
			if err := ValidFof(v); err != nil {
				t.Errorf("[%d,%d] v invalid\ninput=%s\nerr=%v", ii, jj, ss.input, err)
				return
			}
			if v != TrueObj {
				dict := NewDict()
				dict.Set("var", NewInt(1))
				impl, _ := FuncCAD(g, "CAD", []interface{}{
					NewFmlImpl(ret_c0le, expect), dict,
				})
				repl, _ := FuncCAD(g, "CAD", []interface{}{
					NewFmlImpl(expect, ret_c0le), dict,
				})
				t.Errorf("[%d,%d] failed\ninput =%v\nexpect=%v\nactual=%v OR %v <= 0\nimpl=%v\nrepl=%v", ii, jj, ss.input, expect, ret, c0, impl, repl)
				return
			}
		}
	}
}

func TestSdcQEcont(t *testing.T) {
	funcname := "TestSdcQEcont"

	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip %s... (no cas)\n", funcname)
		return
	}
	defer g.Close()

	qeopt := NewQEopt()
	qeopt.SetG(g)

	lvs := []Level{0, 1, 2, 3, 4, 5}

	for ii, tt := range []struct {
		varstr string
		lv     Level // x のレベル
	}{
		{"x,a,b,c", 0},
		{"x,b,a,c", 0},
		{"a,x,b,c", 1},
		{"b,x,a,c", 1},
		{"a,b,c,x", 3},
		{"b,a,c,x", 3},
	} {
		vstr := fmt.Sprintf("vars(%s);", tt.varstr)
		_, err := evalstr(g, vstr)
		if err != nil {
			t.Errorf("[%d] `%s` failed: %s", ii, vstr, err)
			return
		}

		lv := tt.lv

		for jj, ss := range []struct {
			input     string
			rmin      []string
			rmax      []string
			expect_ge string
			expect_le string
			expect_eq string
		}{
			{
				"a*x^2-3*x+1",
				[]string{"x-5"},
				[]string{},
				"a>0",
				"25*a-14<=0",
				"a>0 && 25*a-14<=0",
			}, {
				"a*x^2-5*x-3",
				[]string{},
				[]string{"x+7"},
				"49*a+32>=0",
				"a<0",
				"49*a+32>=0 && a<0",
			}, {
				"a*x^2-5*x-3",
				[]string{"3*x-7"},
				[]string{},
				"a>0",
				"49*a-132<=0",
				"a > 0 && 49*a-132<=0",
			}, {
				"a*x^2-3*x+1",
				[]string{"x-b"},
				[]string{},
				"a*b^2-3*b+1>=0 || a>0 || 2*a*b-3>=0",
				"a<=0 || (4*a-9<=0 && 2*a*b-3<=0) || (a*b^2-3*b+1<=0 && 4*a-9<=0)",
				"(a*b^2-3*b+1>=0 && a<=0) || (a>0 && 4*a-9<=0 && 2*a*b-3<=0) || (2*a*b-3>=0 && a<=0) || (a*b^2-3*b+1<=0 && 4*a-9<=0 && a>0)",
			}, {
				"a*x^2-5*x-3",
				[]string{},
				[]string{"x-b"},
				"(12*a+25>=0 && a*b^2-5*b-3>=0) || a>=0 || (12*a+25>=0 && 2*a*b-5<=0)",
				"a*b^2-5*b-3<=0 || 2*a*b-5>=0 || a<0",
				"(12*a+25>=0 && a*b^2-5*b-3>=0 && a<0) || (2*a*b-5>=0 && a>=0) || (12*a+25>=0 && 2*a*b-5<=0 && a<0) || (a*b^2-5*b-3<=0 && a>=0)",
			}, {
				"x^2+a*x+3",
				[]string{"5*x-8"},
				[]string{"3*x-13"},
				"39*a+196>=0",
				"a^2-12>=0 && a<=0",
				"(39*a+196>=0 && 40*a+139<=0) || (a^2-12>=0 && a<=0 && 40*a+139>=0)",
			}, {
				"x^2+a*x+3",
				[]string{"3*x-2"},
				[]string{"7*x-11"},
				"6*a+31>=0",
				"77*a+268<=0",
				"6*a+31>=0 && 77*a+268<=0",
			}, {
				"a*x^2+b*x+3",
				[]string{"x-1"},
				[]string{"x-2"},
				"b+a+3>=0 || 2*b+4*a+3>=0",
				"2*b+4*a+3<=0 || b+a+3<=0 || (4*a-3>=0 && b<=0 && b^2-12*a>=0 && a-3<=0)",
				"(b+a+3>=0 && 2*b+4*a+3<=0) || (2*b+4*a+3>=0 && b+a+3<=0) || (b<=0 && 4*a-3>=0 && b^2-12*a>=0 && b+a+3>=0 && a-3<=0)",
			}, {
				// 区間に変数1
				"x^2+b*x+3",
				[]string{"x-1"},
				[]string{"x-a"},
				"(a-1>=0 && b+4>=0) || (a*b+a^2+3>=0 && a-1>=0)",
				"(a*b+a^2+3<=0 && a-1>=0) || (b<=0 && a^2-3>=0 && b^2-12>=0 && a-1>=0)",
				"(a*b+a^2+3>=0 && b+4<=0 && a-1>=0) || (b<=0 && a^2-3>=0 && b^2-12>=0 && b+4>=0 && a-1>=0) || (a*b+a^2+3<=0 && b+4>=0 && a-1>=0)",
			}, {
				// 区間に変数2
				"x^2-4*x+3",
				[]string{"3*x-a"},
				[]string{"5*x-7*b"},
				"(21*b-5*a>=0 && a-3<=0) || (21*b-5*a>=0 && a-9>=0) || (7*b-15>=0 && a-9<=0)",
				"(21*b-5*a>=0 && a-3>=0 && a-9<=0) || (7*b-5>=0 && a-3<=0)",
				"(7*b-5>=0 && a-3<=0) || (a-9<=0 && 7*b-15>=0)",
			}, {
				// 左端点がゼロ点
				"(3*x-a)*(5*x-7*b+1)",
				[]string{"3*x-a"},
				[]string{"5*x-7*b"},
				"21*b-5*a>=0",
				"21*b-5*a>=0",
				"21*b-5*a>=0",
			}, {
				// 右端点がゼロ点
				"(3*x-a+1)*(5*x-7*b)",
				[]string{"3*x-a"},
				[]string{"5*x-7*b"},
				"21*b-5*a>=0",
				"21*b-5*a>=0",
				"21*b-5*a>=0",
			},
		} {

			p, err := str2poly(g, ss.input)
			if err != nil {
				t.Errorf("[%d,%d] parse error: `%s`\nerr=%s\ninput=%s.", ii, jj, vstr, err, ss.input)
				return
			}

			rmin := make([]*Atom, len(ss.rmin))
			for kk, rr := range ss.rmin {
				var err error
				rmin[kk], err = str2atom(g, rr+" >= 0")
				if err != nil {
					t.Errorf("[%d,%d] rmin parse error: %s\nerr=%s\neval=%s", ii, jj, vstr, err, rr)
					return
				}
			}

			rmax := make([]*Atom, len(ss.rmax))
			for kk, rr := range ss.rmax {
				var err error
				rmax[kk], err = str2atom(g, rr+" <= 0")
				if err != nil {
					t.Errorf("[%d,%d] rmax parse error: %s\nerr=%s\neval=%s", ii, jj, vstr, err, rr)
					return
				}
			}

			for _, vv := range []struct {
				op     OP
				expect string
			}{
				{GE, ss.expect_ge},
				{LE, ss.expect_le},
				{EQ, ss.expect_eq},
			} {
				qeopt.SetAlgo(QEALGO_SDC|QEALGO_ATOM, true)

				expect, err := str2fof(g, vv.expect)
				if err != nil {
					t.Errorf("[%d,%d,%s] parse error: %s\nerr=%s\nexpect=%s", ii, jj, vv.op, vstr, err, vv.expect)
					return
				}

				atm := NewAtom(p, vv.op).(*Atom)

				ret := SdcQEcont(atm, rmin, rmax, lv, TrueObj, *qeopt)
				if ret == nil {
					t.Errorf("[%d,%d,%s] %s return nil.", ii, jj, vv.op, atm)
					return
				}
				if err := ValidFof(ret); err != nil {
					t.Errorf("[%d,%d,%s] ret invalid\ninput=%s\nerr=%v", ii, jj, vv.op, ss.input, err)
					return
				}

				v := NewForAll(lvs, NewFmlEquiv(ret, expect))
				qeopt.SetAlgo(QEALGO_SDC|QEALGO_ATOM, false)
				v = g.QE(v, qeopt)
				if v != TrueObj {
					dict := NewDict()
					dict.Set("var", NewInt(1))
					impl, _ := FuncCAD(g, "CAD", []interface{}{
						NewFmlImpl(ret, expect), dict,
					})
					repl, _ := FuncCAD(g, "CAD", []interface{}{
						NewFmlImpl(expect, ret), dict,
					})

					smpl, _ := FuncCAD(g, "CAD", []interface{}{
						ret, dict,
					})
					t.Errorf("[%d,%d] failed\ninput =%v\nrmin=%v\nrmax=%v\nexpect=%v\nactual=%v\n      =%v\nimpl=%v\nrepl=%v", ii, jj, atm, rmin, rmax, expect, ret, smpl, impl, repl)
					return
				}
			}
		}
	}
}

func TestAtomQE(t *testing.T) {
	funcname := "TestAtomQE"

	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip %s... (no cas)\n", funcname)
		return
	}
	defer g.Close()

	qeopt := NewQEopt()
	qeopt.SetG(g)
	if (qeopt.Algo & QEALGO_ATOM) == 0 {
		t.Errorf("no QEALGO_ATOM")
		return
	}

	lvs := []Level{0, 1, 2, 3, 4, 5}

	for ii, tt := range []struct {
		varstr string
		lv     Level // x のレベル
	}{
		{"x,a,b,c", 0},
		{"x,b,a,c", 0},
		{"a,x,b,c", 1},
		{"b,x,a,c", 1},
		{"a,b,c,x", 3},
		{"b,a,c,x", 3},
	} {
		vstr := fmt.Sprintf("vars(%s);", tt.varstr)
		_, err := evalstr(g, vstr)
		if err != nil {
			t.Errorf("[%d] `%s` failed: %s", ii, vstr, err)
			return
		}

		lv := tt.lv

		for jj, ss := range []struct {
			input     string
			expect_ge string // ex([x], input >= 0)
			expect_eq string // ex([x], input == 0)
			expect_le string // ex([x], input <= 0)
			expect_ne string // ex([x], input != 0)
		}{
			{
				"a*x^3+b*x^2-5*x+1",
				"true",
				"a!=0 || 4*b^3-25*b^2+90*a*b+27*a^2-500*a<=0",
				"a!=0 || 4*b^3-25*b^2+90*a*b+27*a^2-500*a<=0",
				"true",
			}, {
				"a*x^3+b*x^2-5*x+c",
				"a!=0 || b>=0 || 27*a^2*c^2+(4*b^3+90*a*b)*c-25*b^2-500*a<=0",
				"a!=0 || 27*a^2*c^2+(4*b^3+90*a*b)*c-25*b^2-500*a<=0",
				"a!=0 || 27*a^2*c^2+(4*b^3+90*a*b)*c-25*b^2-500*a<=0 || b<=0",
				"true",
			}, {
				"a*x^3+b*x^2+(a-b)*x+2*a-b^2",
				"a!=0 || 4*b^5+(-27*a^2+18*a+1)*b^4+(-18*a^2-6*a)*b^3+(108*a^3-47*a^2)*b^2+48*a^3*b-112*a^4>=0",
				"a!=0 || 4*b^5+(-27*a^2+18*a+1)*b^4+(-18*a^2-6*a)*b^3+(108*a^3-47*a^2)*b^2+48*a^3*b-112*a^4>=0",
				"true",
				"b!=0 || a!=0",
			},
		} {
			p, err := str2poly(g, ss.input)
			if err != nil {
				t.Errorf("[%d,%d] parse error: `%s`\nerr=%s\ninput=%s.", ii, jj, vstr, err, ss.input)
				return
			}

			for _, vv := range []struct {
				op     OP
				expect string
			}{
				{GE, ss.expect_ge},
				{EQ, ss.expect_eq},
				{LE, ss.expect_le},
				{NE, ss.expect_ne},
			} {
				expect, err := str2fof(g, vv.expect)
				if err != nil {
					t.Errorf("[%d,%d,%s] parse error: %s\nerr=%s\nexpect=%s", ii, jj, vv.op, vstr, err, vv.expect)
					return
				}

				atom := NewAtom(p, vv.op).(*Atom)

				qeopt.SetAlgo(QEALGO_SDC|QEALGO_ATOM, true)
				ret := AtomQE(atom, lv, TrueObj, *qeopt)
				if ret == nil {
					t.Errorf("[%d,%d,%s] %s return nil.", ii, jj, vv.op, atom)
					return
				}
				if err := ValidFof(ret); err != nil {
					t.Errorf("[%d,%d,%s] ret invalid\ninput=%s\nerr=%v", ii, jj, vv.op, ss.input, err)
					return
				}

				v := NewForAll(lvs, NewFmlEquiv(ret, expect))
				qeopt.SetAlgo(QEALGO_SDC|QEALGO_ATOM, false)
				v = g.QE(v, qeopt)
				if v != TrueObj {
					dict := NewDict()
					dict.Set("var", NewInt(1))
					impl, _ := FuncCAD(g, "CAD", []interface{}{
						NewFmlImpl(ret, expect), dict,
					})
					repl, _ := FuncCAD(g, "CAD", []interface{}{
						NewFmlImpl(expect, ret), dict,
					})

					smpl, _ := FuncCAD(g, "CAD", []interface{}{
						ret, dict,
					})
					t.Errorf("[%d,%d] failed\ninput =%v\nexpect=%v\nactual=%v\n      =%v\nimpl=%v\nrepl=%v", ii, jj, atom, expect, ret, smpl, impl, repl)
					return
				}
			}
		}
	}
}

func TestSdcQEpoly(t *testing.T) {
	funcname := "TestSdcQEpoly"
	print_log := false

	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip %s... (no cas)\n", funcname)
		return
	}
	defer g.Close()

	qeopt := NewQEopt()
	qeopt.SetG(g)

	lvs := []Level{0, 1, 2, 3, 4, 5}

	for ii, tt := range []struct {
		varstr string
		lv     Level // x のレベル
	}{
		{"x,a,b,c", 0},
		{"x,b,a,c", 0},
		{"a,x,b,c", 1},
		{"b,x,a,c", 1},
		{"a,b,c,x", 3},
		{"b,a,c,x", 3},
	} {
		vstr := fmt.Sprintf("vars(%s);", tt.varstr)
		_, err := evalstr(g, vstr)
		if err != nil {
			t.Errorf("[%d] `%s` failed: %s", ii, vstr, err)
			return
		}

		lv := tt.lv

		for jj, ss := range []struct {
			input     string
			rng       string
			expect_gg string // ex([x], input >= 0 && rng >= 0)
			expect_gl string // ex([x], input >= 0 && rng <= 0)
			expect_lg string // ex([x], input <= 0 && rng >= 0)
			expect_ll string // ex([x], input <= 0 && rng <= 0)
		}{
			{
				"x^2-a*x+1",
				"x+b",
				"true",
				"true",
				"(2*b+a>=0 && a+2<=0) || (b^2+a*b+1<=0 && a+2<=0) || (a-2>=0 && 2*b+a>=0) || (a-2>=0 && b^2+a*b+1<=0)",
				"(a-2>=0 && 2*b+a<=0) || (a+2<=0 && 2*b+a<=0) || (a-2>=0 && b^2+a*b+1<=0) || (b^2+a*b+1<=0 && a+2<=0)",
			}, {
				"x^2-3*x+1",
				"a*x+b",
				"a!=0 || b>=0",
				"a!=0 || b<=0",
				"b^2+3*a*b+a^2<=0 || 2*b+3*a>=0",
				"2*b+3*a<=0 || b^2+3*a*b+a^2<=0",
			}, {
				"b*x^3-6*x^2+11*x-6",
				"a*x+1",
				"(b+6*a^3+11*a^2+6*a>=0 && a<0) || (b+6*a^3+11*a^2+6*a<=0 && a>0) || (b>0 && a>=0) || (b<0 && a<=0) || (243*b^2-451*b+207<=0 && 9*a^2+11*a+3<=0) || (18*a+11>=0 && 243*b^2-451*b+207<=0)",
				"a*b < 0",
				"true",
				"(b+6*a^3+11*a^2+6*a<=0 && a<0) || (b+6*a^3+11*a^2+6*a>=0 && a>0) || (243*b^2-451*b+207<=0 && 9*a^2+11*a+3<=0) || (243*b^2-451*b+207<=0 && 18*a+11<=0)",
			},
		} {

			p, err := str2poly(g, ss.input)
			if err != nil {
				t.Errorf("[%d,%d] parse error: `%s`\nerr=%s\ninput=%s.", ii, jj, vstr, err, ss.input)
				return
			}

			rng, err := str2poly(g, ss.rng)
			if err != nil {
				t.Errorf("[%d,%d] rng parse error: %s\nerr=%s\neval=%s", ii, jj, vstr, err, ss.rng)
				return
			}

			for _, vv := range []struct {
				opi    OP
				opr    OP
				expect string
			}{
				{GE, GE, ss.expect_gg},
				{GE, LE, ss.expect_gl},
				{LE, GE, ss.expect_lg},
				{LE, LE, ss.expect_ll},
			} {
				qeopt.SetAlgo(QEALGO_SDC|QEALGO_ATOM, true)

				expect, err := str2fof(g, vv.expect)
				if err != nil {
					t.Errorf("[%d,%d,%s,%s] parse error: %s\nerr=%s\nexpect=%s", ii, jj, vv.opi, vv.opr, vstr, err, vv.expect)
					return
				}

				atm := NewAtom(p, vv.opi).(*Atom)
				atr := NewAtom(rng, vv.opr).(*Atom)

				if print_log {
					fmt.Printf("\n#### [%d,%d,%s,%s] %s.\n", ii, jj, vv.opi, vv.opr, tt.varstr)
				}
				ret := SdcQEpoly(atm, atr, lv, TrueObj, *qeopt)
				if ret == nil {
					t.Errorf("[%d,%d,%s,%s] poly returns nil %s\n\nexpect=%s", ii, jj, vv.opi, vv.opr, vstr, vv.expect)
					return
				}
				if err := ValidFof(ret); err != nil {
					t.Errorf("[%d,%d,%s,%s] poly returns nil %s\n\nexpect=%s\nret=%V", ii, jj, vv.opi, vv.opr, vstr, vv.expect, ret)
					return
				}

				v := NewForAll(lvs, NewFmlEquiv(ret, expect))
				qeopt.SetAlgo(QEALGO_SDC|QEALGO_ATOM, false)
				v = g.QE(v, qeopt)
				if v != TrueObj {
					dict := NewDict()
					dict.Set("var", NewInt(1))
					impl, _ := FuncCAD(g, "CAD", []interface{}{
						NewFmlImpl(ret, expect), dict,
					})
					repl, _ := FuncCAD(g, "CAD", []interface{}{
						NewFmlImpl(expect, ret), dict,
					})

					smpl, _ := FuncCAD(g, "CAD", []interface{}{
						ret, dict,
					})

					t.Errorf("[%d,%d] failed\ninput =%v := %v %v 0\nrng=%v := %v %v 0\nexpect=%v\nactual=%v\n      =%v\nimpl=%v\nrepl=%v", ii, jj,
						atm, p, vv.opi,
						atr, rng, vv.opr,
						expect, ret, smpl, impl, repl)
					return
				}
			}
		}
	}
}

func TestSdcAtomQE(t *testing.T) {
	funcname := "TestSdcAtomQE"

	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip %s... (no cas)\n", funcname)
		return
	}
	defer g.Close()

	qeopt := NewQEopt()
	qeopt.SetG(g)

	qecond := NewQeCond()

	lvs := []Level{0, 1, 2, 3, 4, 5}

	for ii, tt := range []struct {
		varstr string
		lv     Level // x のレベル
	}{
		{"x,a,b,c", 0},
		{"x,b,a,c", 0},
		{"a,x,b,c", 1},
		{"b,x,a,c", 1},
		{"a,b,c,x", 3},
		{"b,a,c,x", 3},
	} {
		vstr := fmt.Sprintf("vars(%s);", tt.varstr)
		_, err := evalstr(g, vstr)
		if err != nil {
			t.Errorf("[%d] `%s` failed: %s", ii, vstr, err)
			return
		}

		lv := tt.lv

		for jj, ss := range []struct {
			input     string
			rng       string
			expect_ge string // ex([x], input >= 0 && rng)
			expect_le string // ex([x], input <= 0 && rng)
			expect_eq string // ex([x], input == 0 && rng)
		}{
			{
				"8*a*x-7*b",
				"true",
				"a != 0 || b <= 0",
				"a != 0 || b >= 0",
				"a != 0 || b == 0",
			}, {
				"a*x^11-b*x^8-c*x^4+1",
				"true",
				"true",
				"a != 0 || b > 0 || c > 0 && b >= 0 || c > 0 && c^2+4*b >= 0",
				"a != 0 || b > 0 || c > 0 && b >= 0 || c > 0 && c^2+4*b >= 0",
			}, {
				"a*x^8+b*x^4+1",
				"x >= 0",
				"true",
				"a<0 || (b<0 && (a<=0 || b^2-4*a>=0))",
				"a<0 || (b<0 && (a<=0 || b^2-4*a>=0))",
			},
		} {
			p, err := str2poly(g, ss.input)
			if err != nil {
				t.Errorf("[%d,%d] parse error: `%s`\nerr=%s\ninput=%s.", ii, jj, vstr, err, ss.input)
				return
			}

			rng, err := str2fof(g, ss.rng)
			if err != nil {
				t.Errorf("[%d,%d] rng parse error: %s\nerr=%s\neval=%s", ii, jj, vstr, err, ss.rng)
				return
			}

			for _, vv := range []struct {
				op     OP
				expect string
			}{
				{GE, ss.expect_ge},
				{LE, ss.expect_le},
				{EQ, ss.expect_eq},
			} {
				qeopt.SetAlgo(QEALGO_SDC|QEALGO_ATOM, true)

				var expect Fof
				if vv.expect != "" {
					expect, err = str2fof(g, vv.expect)
					if err != nil {
						t.Errorf("[%d,%d,%s] parse error: %s\nerr=%s\nexpect=%s", ii, jj, vv.op, vstr, err, vv.expect)
						return
					}
				}

				atom := NewAtom(p, vv.op).(*Atom)
				input := NewFmlAnd(atom, rng)

				ret := SdcAtomQE(input, lv, *qeopt, *qecond)
				if ret == nil {
					if expect == nil {
						continue
					}

					t.Errorf("[%d,%d,%s] SdcAtomQE returns nil %s\n\nexpect=%s", ii, jj, vv.op, vstr, vv.expect)
					return
				}
				if expect == nil {
					t.Errorf("[%d,%d,%s] SdcAtomQE returns NON-nil %s\n\nexpect=%s\nactual=%v", ii, jj, vv.op, vstr, vv.expect, ret)
					return
				}
				if err := ValidFof(ret); err != nil {
					t.Errorf("[%d,%d,%s] poly returns nil %s\n\nexpect=%s\nret=%V", ii, jj, vv.op, vstr, vv.expect, ret)
					return
				}

				v := NewForAll(lvs, NewFmlEquiv(ret, expect))
				qeopt.SetAlgo(QEALGO_SDC|QEALGO_ATOM, false)
				v = g.QE(v, qeopt)
				if v != TrueObj {
					dict := NewDict()
					dict.Set("var", NewInt(1))
					impl, _ := FuncCAD(g, "CAD", []interface{}{
						NewFmlImpl(ret, expect), dict,
					})
					repl, _ := FuncCAD(g, "CAD", []interface{}{
						NewFmlImpl(expect, ret), dict,
					})

					smpl, _ := FuncCAD(g, "CAD", []interface{}{
						ret, dict,
					})

					t.Errorf("[%d,%d] failed\ninput =%v := %v %v 0\nrng=%v\nexpect=%v\nactual=%v\n      =%v\nimpl=%v\nrepl=%v", ii, jj,
						atom, p, vv.op, rng,
						expect, ret, smpl, impl, repl)
					return
				}
			}
		}
	}
}

func TestSdcAtomQE2(t *testing.T) {
	funcname := "TestSdcAtomQE2"

	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip %s... (no cas)\n", funcname)
		return
	}
	defer g.Close()

	qeopt := NewQEopt()
	qeopt.SetG(g)

	qecond := NewQeCond()

	lvs := []Level{0, 1, 2, 3, 4, 5}

	for ii, tt := range []struct {
		varstr string
		lv     Level // x のレベル
	}{
		{"x,a,b,c", 0},
		{"x,b,a,c", 0},
		{"a,x,b,c", 1},
		{"b,x,a,c", 1},
		{"a,b,c,x", 3},
		{"b,a,c,x", 3},
	} {
		vstr := fmt.Sprintf("vars(%s);", tt.varstr)
		_, err := evalstr(g, vstr)
		if err != nil {
			t.Errorf("[%d] `%s` failed: %s", ii, vstr, err)
			return
		}

		lv := tt.lv

		for jj, ss := range []struct {
			input  string
			expect string
		}{
			{
				"8*a*x-7*b > 0",
				"a != 0 || b < 0",
			},
		} {
			// fmt.Printf("[%d,%d] %s, input=%v\n", ii, jj, vstr, ss.input)
			input, err := str2fof(g, ss.input)
			if err != nil {
				t.Errorf("[%d,%d] parse error: `%s`\nerr=%s\ninput=%s.", ii, jj, vstr, err, ss.input)
				return
			}

			qeopt.SetAlgo(QEALGO_SDC|QEALGO_ATOM, true)

			var expect Fof
			expect, err = str2fof(g, ss.expect)
			if err != nil {
				t.Errorf("[%d,%d] parse error: %s\nerr=%s\nexpect=%s", ii, jj, vstr, err, ss.expect)
				return
			}

			ret := SdcAtomQE(input, lv, *qeopt, *qecond)
			if ret == nil {
				if expect == nil {
					continue
				}

				t.Errorf("[%d,%d] SdcAtomQE returns nil %s\n\nexpect=%s", ii, jj, vstr, ss.expect)
				return
			}
			if expect == nil {
				t.Errorf("[%d,%d] SdcAtomQE returns NON-nil %s\n\nexpect=%s\nactual=%v", ii, jj, vstr, ss.expect, ret)
				return
			}
			if err := ValidFof(ret); err != nil {
				t.Errorf("[%d,%d] poly returns nil %s\n\nexpect=%s\nret=%V", ii, jj, vstr, ss.expect, ret)
				return
			}

			v := NewForAll(lvs, NewFmlEquiv(ret, expect))
			qeopt.SetAlgo(QEALGO_SDC|QEALGO_ATOM, false)
			v = g.QE(v, qeopt)
			if v != TrueObj {
				dict := NewDict()
				dict.Set("var", NewInt(1))
				impl, _ := FuncCAD(g, "CAD", []interface{}{
					NewFmlImpl(ret, expect), dict,
				})
				repl, _ := FuncCAD(g, "CAD", []interface{}{
					NewFmlImpl(expect, ret), dict,
				})

				smpl, _ := FuncCAD(g, "CAD", []interface{}{
					ret, dict,
				})

				t.Errorf("[%d,%d] failed\ninput =%v\nexpect=%v\nactual=%v\n      =%v\nimpl=%v\nrepl=%v", ii, jj,
					input, expect, ret, smpl, impl, repl)
				return
			}
		}
	}
}
