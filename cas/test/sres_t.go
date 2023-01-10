package cas

/*
 * test for sres()
 */
import (
	"github.com/hiwane/ganrac"
	"testing"
)

func sresTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T, ii int,
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

		act := cas.Sres(p, q, x.Level(), int32(jj))
		if act == nil {
			t.Errorf("ii=(%d,%d) %s sres error [%s]", ii, jj, expstr, err.Error())
			return
		}

		if !exp.Equals(act) {
			t.Errorf("ii=(%d,%d) failed expect=%s, actual=%s\np=%s\nq=%s\nx=%s", ii, jj, exp, act, pstr, qstr, xstr)
			return
		}
	}
}

func SresTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T) {
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
				"-b*z+2*a*w",
				"z",
			},
			[]string{
				"c*z^2-b*w*z+a*w^2",
				"0",
			},
		}, {
			"a*x^2+b*x+c",
			"y*x^2+z*x+w",
			"x",
			[]string{
				"(-2*c*a+b^2)*y-b*a*z+2*a^2*w",
				"-b*y+a*z",
			},
			[]string{
				"2*c*y^2+(-b*z-2*a*w)*y+a*z^2",
				"b*y-a*z",
			},
		}, {
			// hong93
			"a*x^2+b*x+1",
			"b*x^3+c*x+a",
			"x",
			[]string{
				"2*a^4-c*b*a^2+3*b^2*a-b^4", // T_B
				"c*a^2-b*a+b^3",             // S_B
				"0",
			},
			[]string{
				"-2*c*b*a+3*b^2",
				"3*b^2",
			},
		}, {
			// hong93
			"a*x^2+b*x+1",
			"c*x^2+b*x+a",
			"x",
			[]string{
				"2*a^3+(-b^2-2*c)*a+c*b^2", // T_C
				"b*a-c*b",                  // S_C
			},
			[]string{
				"-2*c*a^2+b^2*a-c*b^2+2*c^2",
				"-b*a+c*b",
			},
		},
	} {
		sresTest(g, cas, t, ii+5000, ss.f, ss.g, ss.x, ss.expect1)
		sresTest(g, cas, t, ii+1000, ss.g, ss.f, ss.x, ss.expect2)
	}
}
