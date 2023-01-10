package cas

/*
 * test for psc()
 */
import (
	"github.com/hiwane/ganrac"
	"testing"
)

func pscTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T, ii int,
	pstr, qstr, xstr string, expects []string) {
	p, err := str2poly(g, pstr)
	if err != nil {
		t.Errorf("ii=%d, %s input error [%s]", ii, pstr, err.Error())
		return
	}
	q, err := str2poly(g, qstr)
	if err != nil {
		t.Errorf("ii=%d, %s input error [%s]", ii, qstr, err.Error())
		return
	}
	x, err := str2poly(g, xstr)
	if err != nil {
		t.Errorf("ii=%d, %s input error [%s]", ii, xstr, err.Error())
		return
	}

	for jj, expstr := range expects {
		exp, err := str2robj(g, expstr)
		if err != nil {
			t.Errorf("ii=(%d,%d) %s expect error [%s]", ii, jj, expstr, err.Error())
			return
		}

		act := cas.Psc(p, q, x.Level(), int32(jj))
		if act == nil {
			t.Errorf("ii=(%d,%d) %s psc error [%s]", ii, jj, expstr, err.Error())
			return
		}

		if !exp.Equals(act) {
			t.Errorf("ii=(%d,%d) failed expect=%s, actual=%s", ii, jj, exp, act)
			return
		}
	}
}

func PscTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T) {
	for ii, ss := range []struct {
		f, g, x string
		expect1 []string
		expect2 []string
	}{
		{
			"a*x^2+b*x+c",
			"z*x+w",
			"x",
			[]string{
				"z^2*c-z*w*b+w^2*a",
				"z",
			},
			[]string{
				"z^2*c-z*w*b+w^2*a",
				"z",
			},
		}, {
			"a*x^2+b*x+c",
			"y*x^2+z*x+w",
			"x",
			[]string{
				"y^2*c^2+(-y*z*b+(-2*y*w+z^2)*a)*c+y*w*b^2-z*w*a*b+w^2*a^2",
				"-y*b+z*a",
				"1",
			},
			[]string{
				"y^2*c^2+(-y*z*b+(-2*y*w+z^2)*a)*c+y*w*b^2-z*w*a*b+w^2*a^2",
				"y*b-z*a",
				"1",
			},
		}, {
			"a*x^3+b*x^2+c*x+d",
			"z*x+w",
			"x",
			[]string{
				"-z^3*d+z^2*w*c-z*w^2*b+w^3*a",
				"z^2",
			},
			[]string{
				"z^3*d-z^2*w*c+z*w^2*b-w^3*a",
				"z^2",
			},
		}, {
			"a*x^3+b*x^2+c*x+d",
			"y*x^2+z*x+w",
			"x",
			[]string{
				"y^3*d^2+(-y^2*z*c+(-2*y^2*w+y*z^2)*b+(3*y*z*w-z^3)*a)*d+y^2*w*c^2+(-y*z*w*b+(-2*y*w^2+z^2*w)*a)*c+y*w^2*b^2-z*w^2*a*b+w^3*a^2",
				"y^2*c-y*z*b+(-y*w+z^2)*a",
				"y",
			},
			[]string{
				"y^3*d^2+(-y^2*z*c+(-2*y^2*w+y*z^2)*b+(3*y*z*w-z^3)*a)*d+y^2*w*c^2+(-y*z*w*b+(-2*y*w^2+z^2*w)*a)*c+y*w^2*b^2-z*w^2*a*b+w^3*a^2",
				"y^2*c-y*z*b+(-y*w+z^2)*a",
				"y",
			},
		}, {
			// 計算機代数の基礎理論 p86 例題4-31
			"(2*x^3-1)*(x+3)*(3*x-2)^2",
			"(x-5)*(3*x-2)^2",
			"x",
			[]string{
				"0",
				"0",
				"2^3*3^9*83*3^2",
				"3^4*3^2",
			},
			[]string{},
		},
	} {
		pscTest(g, cas, t, ii, ss.f, ss.g, ss.x, ss.expect1)
		pscTest(g, cas, t, ii+1000, ss.g, ss.f, ss.x, ss.expect2)
	}
}
