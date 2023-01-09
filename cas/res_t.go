package cas

/*
 * test for Resultant()
 */
import (
	"github.com/hiwane/ganrac"
	"testing"
)

func resultantTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T, ii int,
	ps, qs, xs, expect string) {

	p, err := str2poly(g, ps)
	if err != nil {
		t.Errorf("ii=%d, %s input.p error [%s]", ii, ps, err.Error())
		return
	}

	q, err := str2poly(g, qs)
	if err != nil {
		t.Errorf("ii=%d, %s input.q error [%s]", ii, qs, err.Error())
		return
	}

	x, err := str2poly(g, xs)
	if err != nil {
		t.Errorf("ii=%d, %s input.x error [%s]", ii, xs, err.Error())
		return
	}

	exp, err := str2poly(g, expect)
	if err != nil {
		t.Errorf("ii=%d, %s expect error [%s]", ii, expect, err.Error())
		return
	}

	act := cas.Resultant(p, q, x.Level())
	if act == nil {
		t.Errorf("ii=%d, execute failed  res(%s,%s,%s).", ii, ps, qs, xs)
		return
	}

	if !exp.Equals(act) {
		t.Errorf("ii=%d, expect=%s, actual=%s.", ii, expect, act)
		return
	}
}

func ResultantTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T) {
	for ii, ss := range []struct {
		p, q, x string
		expect  string
	}{
		{
			"3*x^2+(-2*z^2*y^3+5*z)*x+2*y-3*z",
			"7*x^2+11*y^2-z^2*y",
			"x",
			"308*z^4*y^8-28*z^6*y^7-1540*z^3*y^5+(140*z^5+1089)*y^4+(-198*z^2-924)*y^3+(9*z^4+2009*z^2+1386*z+196)*y^2+(-175*z^4-126*z^3-588*z)*y+441*z^2",
		}, {
			"3*x^2+(-2*z^2*y^3+5*z)*x+2*y-3*z",
			"7*x^2+11*y^2-z^2*y",
			"y",
			"1372*z^4*x^8+(1386*z^4+4312*z^2)*x^5+(2310*z^5+11979)*x^4+(-6*z^8-28*z^6-1386*z^5+39930*z)*x^3+(-10*z^9+34001*z^2-23958*z+3388)*x^2+(6*z^9+1210*z^3-39930*z^2)*x-726*z^3+11979*z^2",
		}, {
			"3*x^2+(-2*z^2*y^3+5*z)*x+2*y-3*z",
			"7*x^2+11*y^2-z^2*y",
			"z",
			"196*y^6*x^6-84*y^4*x^5+(616*y^8+9*y^2-175*y)*x^4+(-132*y^6-56*y^5+210*y)*x^3+(484*y^10-263*y^3-63*y)*x^2+(-88*y^7+330*y^3)*x+4*y^4-99*y^3",
		}, {
			// hong93 R_B
			"a*x^2+b*x+1",
			"b*x^3+c*x+a",
			"x",
			"a^5-c*b*a^3+(3*b^2+c^2)*a^2+(-b^4-2*c*b)*a+c*b^3+b^2",
		}, {
			// hong93 R_C
			"a*x^2+b*x+1",
			"c*x^2+b*x+a",
			"x",
			"a^4+(-b^2-2*c)*a^2+(c+1)*b^2*a-c*b^2+c^2",
		}, {
			// 計算機代数の基礎理論 p77 例題4-16
			"x^2+y^2-1",
			"y-x^2+1",
			"y",
			"x^4-x^2",
		}, {
			// 計算機代数の基礎理論 p77 例題4-16
			"x^2+y^2-1",
			"2*x^2+x*y+y^2-1",
			"y",
			"2*x^4-x^2",
		},
	} {
		resultantTest(g, cas, t, ii, ss.p, ss.q, ss.x, ss.expect)
	}
}
