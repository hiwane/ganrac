package cas

import (
	"fmt"
	"github.com/hiwane/ganrac"
	"strings"
	"testing"
)

func sresTest(gan *ganrac.Ganrac, cas ganrac.CAS, t *testing.T, eyec string, x ganrac.Level, fi, gi string, expect []string) {

	f, err := str2poly(gan, fi)
	if err != nil {
		t.Errorf("[%s] f=%s input error [%s]", eyec, fi, err)
		return
	}
	g, err := str2poly(gan, gi)
	if err != nil {
		t.Errorf("[%s] g=%s input error [%s]", eyec, gi, err)
		return
	}

	exp := make([]*ganrac.Poly, len(expect))
	for i, es := range expect {
		exp[i], err = str2poly(gan, es)
		if err != nil {
			t.Errorf("[%s] expect[%d]=%s input error [%s]", eyec, i, es, err)
			return
		}
	}

	for cc := int32(0); cc < int32(4); cc++ {
		v := cas.Sres(f, g, x, cc)
		if v == nil {
			t.Errorf("[%s,%d] return nil, f=%s, g=%s", eyec, cc, fi, gi)
			return
		}

		if v.Len() != len(expect) {
			t.Errorf("[%s,%d,%v] length expect=%d, actual=%d:%v, f=%v, g=%v", eyec, cc, x, len(expect), v.Len(), v, f, g)
			return
		}

		for i, es := range expect {
			w, _ := v.Get(ganrac.NewInt(int64(i)))
			var p ganrac.RObj
			if cc == 0 {
				p = exp[i]
			} else if cc == 1 {
				p = exp[i].Coef(x, uint(i))
			} else if cc == 2 {
				p = exp[i].Coef(x, 0)
			} else if i == 0 {
				p = exp[i]
			} else {
				p = ganrac.Add(exp[i].Coef(x, 0), ganrac.NewPolyVarn(x, i).Mul(exp[i].Coef(x, uint(i))))
			}

			if !p.Equals(w) {
				t.Errorf("[%s,%d] expect[%d]=%s input error [%s]", eyec, cc, i, es, err)
				return
			}
		}
	}

}

func SresTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T) {
	for jj, vv := range []struct {
		v string
		x ganrac.Level
	}{
		{"x,d,c,b,a", 0},
		{"a,b,c,d,x", 4},
		{"d,b,x,c,a", 2},
	} {
		vstr := fmt.Sprintf("vars(%s);", vv.v)
		_, err := g.Eval(strings.NewReader(vstr))
		if err != nil {
			t.Errorf("[%d] %s failed: %s", jj, vstr, err)
			return
		}

		for ii, ss := range []struct {
			fi     string
			gi     string
			expect []string
		}{
			{
				"a*x^3+b*x^2+c*x+d",
				"3*a*x^2+2*b*x+c",
				[]string{
					"27*d^2*a^3+(-18*d*c*b+4*c^3)*a^2+(4*d*b^3-c^2*b^2)*a",
					"9*a^2*d+(-a*b+6*x*a^2)*c-2*x*a*b^2",
				},
			}, {
				"x^4+a*x^3+b*x^2+c*x+d",
				"4*x^3+3*a*x^2+2*b*x+c",
				[]string{
					"-27*d^2*a^4+(18*d*c*b-4*c^3)*a^3+(-4*d*b^3+c^2*b^2+144*d^2*b-6*d*c^2)*a^2+(-80*d*c*b^2+18*c^3*b-192*d^2*c)*a+16*d*b^4-4*c^2*b^3-128*d^2*b^2+144*d*c^2*b-27*c^4+256*d^3",
					"(48*c+(-32*a-32*x)*b+9*a^3+12*x*a^2)*d+(-3*a+36*x)*c^2+(4*b^2+(-a^2-28*x*a)*b+6*x*a^3)*c+8*x*b^3-2*x*a^2*b^2",
					"16*d+(-a+12*x)*c+(-2*x*a+8*x^2)*b-3*x^2*a^2",
				},
			},
		} {
			sresTest(g, cas, t, fmt.Sprintf("[%d,%d]", ii, jj), vv.x, ss.fi, ss.gi, ss.expect)
		}
	}
}
