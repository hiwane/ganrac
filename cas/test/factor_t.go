package cas

/*
 * test for factor()
 */
import (
	"github.com/hiwane/ganrac"
	"testing"
)

func factorTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T, ii int, input, expect string) {
	one := ganrac.NewInt(1)
	r, err := str2poly(g, input)
	if err != nil {
		t.Errorf("ii=%d, %s input error [%s]", ii, input, err.Error())
		return
	}
	exp, err := str2list(g, expect)
	if err != nil {
		t.Errorf("ii=%d, %s expect error [%s]", ii, expect, err.Error())
		return
	}

	act := cas.Factor(r)
	if act == nil {
		t.Errorf("ii=%d, execute failed %s.", ii, input)
		return
	}
	// fmt.Printf("factor(%s) = %s\n", input, act)

	if exp.Len() != act.Len() {
		t.Errorf("ii=%d, invalid len expect=%d, actual=%d\nexpect=%s\nactual=%s\n",
			ii, exp.Len(), act.Len(), exp, act)
		return
	}
	var mul ganrac.RObj
	mul = one
	e := ganrac.NewInt(0) // 指数
	for i := 0; i < act.Len(); i++ {
		_ui, _ := act.Geti(i)
		ui := _ui.(*ganrac.List)
		if ui == nil || ui.Len() != 2 {
			t.Errorf("ii=%d, %dth element is not a list or length is not 2 [%s]", ii, i, act)
			return
		}
		_ui1, _ := ui.Geti(1)
		ui1 := _ui1.(*ganrac.Int)
		if ui1 == nil || ui1.Cmp(one) < 0 {
			// 指数部が不正
			t.Errorf("ii=%d, %dth element1 is invalid: %s", ii, i, act)
			return
		}
		e = e.Mul(ui1).(*ganrac.Int)

		_ui0, _ := ui.Geti(0)
		if _, ok := _ui0.(*ganrac.Int); ok {
			if i != 0 {
				// [0,0] 番目の要素のみ整数
				t.Errorf("ii=%d, %dth element0-i is invalid: %s", ii, i, act)
				return
			}
		} else if ui0p, ok := _ui0.(*ganrac.Poly); ok {
			if i == 0 {
				// [0,0] 番目の要素のみ整数
				t.Errorf("ii=%d, %dth element0-p is invalid: %s", ii, i, act)
				return
			}
			_, cont := ui0p.PPC()
			if cont.CmpAbs(one) != 0 {
				// 因数の容量が 1 ではなかった
				t.Errorf("ii=%d, %dth element0-c is invalid: %s", ii, i, act)
				return
			}
		} else {
			t.Errorf("ii=%d, %dth element0-a is invalid: %s", ii, i, act)
			return
		}
		ff := _ui0.(ganrac.RObj).Pow(ui1)
		mul = ganrac.Mul(mul, ff)
	}
	/* 復帰値をかけたらもとに戻るか */
	if !mul.Equals(r) {
		t.Errorf("ii=%d, mul=%s,input=%s", ii, mul, input)
		return
	}

	/* 因子は一致? */
	sign := 1
	found := make([]bool, exp.Len())
	for i := 0; i < act.Len(); i++ {
		var aineg ganrac.RObj = nil
		ook := false
		_aix, _ := act.Geti(i, 1)
		aix := _aix.(*ganrac.Int)
		for j := 0; j < exp.Len(); j++ {
			_ai, _ := act.Geti(i, 0)
			_ej, _ := exp.Geti(j, 0)
			ai := _ai.(ganrac.RObj)
			ej := _ej.(ganrac.RObj)
			if !ai.Equals(ej) {
				if aineg == nil {
					aineg = ai.Neg()
				}
				if !aineg.Equals(ej) {
					return
				}
				if aix.Bit(0) != 0 {
					// 指数が奇数の場合のみ符号が入れ替わる
					sign *= -1
				}
			}
			if found[j] {
				// 同じ因子があった
				t.Errorf("ii=%d, invalid f1 exp=%s, act=%s", ii, exp, act)
				return
			}
			found[j] = true

			// 指数が一致するか
			_ejx, _ := act.Geti(i, 1)
			if !_aix.(*ganrac.Int).Equals(_ejx.(*ganrac.Int)) {
				t.Errorf("ii=%d, invalid f1 exp=%s, act=%s", ii, exp, act)
				return
			}
		}
		if !ook {
			// あるはずの因子がみつからなかった
			t.Errorf("ii=%d, i=%d, invalid fx exp=%s, act=%s", ii, i, exp, act)
			return
		}
	}

	// 容量
	_ce, _ := exp.Geti(0, 0)
	_ca, _ := act.Geti(0, 0)
	ce := _ce.(*ganrac.Int)
	ca := _ca.(*ganrac.Int)
	if ce.CmpAbs(ca) != 0 {
		t.Errorf("ii=%d, invalid cont exp=%s, act=%s", ii, ce, ca)
		return
	}

	if ca.Cmp(ce) == 0 && sign < 0 || ca.Cmp(ce) != 0 && sign > 0 {
		t.Errorf("ii=%d, invalid cont sgn=%d, exp=%s, act=%s", ii, sign, ce, ca)
		return
	}
}

func FactorTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T) {
	for ii, ss := range []struct {
		input  string
		expect string
	}{
		{
			"x-1",
			"[[1, 1], [x-1, 1]]",
		}, {
			"2*x-1",
			"[[1, 1], [2*x-1, 1]]",
		}, {
			"2*x-2",
			"[[2, 1], [x-1, 1]]",
		}, {
			"-2*x+2",
			"[[-2, 1], [x-1, 1]]",
		}, {
			"x^2-1",
			"[[1, 1], [x-1, 1], [x+1, 1]]",
		}, {
			"x^2-2*x+1",
			"[[1, 1], [x-1, 2]]",
		}, {
			"-7*x^2+14*x-7",
			"[[-7, 1], [x-1, 2]]",
		}, {
			"-7*x^2+14*x-7",
			"[[-7, 1], [1-x, 2]]",
		}, {
			"x*y - x - y + 1",
			"[[1, 1], [x-1, 1], [y-1, 1]]",
		}, {
			"5*x*y - 5",
			"[[5, 1], [x*y-1, 1]]",
		}, {
			"(5*z^3-45*z^2+135*z-135)*y^2*x^2+(-10*z^3+90*z^2-270*z+270)*y*x+5*z^3-45*z^2+135*z-135",
			"[[5, 1], [x*y-1, 2], [z-3, 3]]",
		},
	} {
		factorTest(g, cas, t, ii, ss.input, ss.expect)
	}
}
