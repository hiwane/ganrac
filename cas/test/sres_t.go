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

	exp := make([]ganrac.RObj, len(expect))
	for i, es := range expect {
		exp[i], err = str2robj(gan, es)
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
			t.Errorf("[%s,cc=%d,x=%v] length expect=%d, actual=%d:%v, f=%v, g=%v\nv=%v\ne=%v", eyec, cc, x, len(expect), v.Len(), v, f, g, v, expect)
			return
		}

		for i, es := range expect {
			w, _ := v.Get(ganrac.NewInt(int64(i)))
			var p ganrac.RObj
			if cc == 0 {
				p = exp[i]
			} else if cc == 1 {
				p = ganrac.Coeff(exp[i], x, uint(i))
			} else if cc == 2 {
				p = ganrac.Coeff(exp[i], x, 0)
			} else if i == 0 {
				p = exp[i]
			} else {
				p = ganrac.Add(ganrac.Coeff(exp[i], x, 0), ganrac.NewPolyVarn(x, i).Mul(ganrac.Coeff(exp[i], x, uint(i))))
			}

			if !p.Equals(w) {
				t.Errorf("[%s,cc=%d,i=%d] not equal. \nexpect=%v\nactual=%v\nes=%v\nv=%v\ne=%v\nf=%v\ng=%v", eyec, cc, i, p, w, es, v, expect, f, g)
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
		// 変数順序による影響確認
		{"x,d,c,b,a", 0}, // 最初
		{"a,b,c,d,x", 4}, // 最後
		{"d,b,x,c,a", 2}, // 中
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
			eyec   string
			expect []string
		}{
			{
				"a*x^3+b*x^2+c*x+d",
				"3*a*x^2+2*b*x+c",
				"cubic",
				[]string{
					"27*d^2*a^3+(-18*d*c*b+4*c^3)*a^2+(4*d*b^3-c^2*b^2)*a",
					"9*a^2*d+(-a*b+6*x*a^2)*c-2*x*a*b^2",
				},
			}, {
				"x^4+a*x^3+b*x^2+c*x+d",
				"4*x^3+3*a*x^2+2*b*x+c",
				"quart",
				[]string{
					"-27*d^2*a^4+(18*d*c*b-4*c^3)*a^3+(-4*d*b^3+c^2*b^2+144*d^2*b-6*d*c^2)*a^2+(-80*d*c*b^2+18*c^3*b-192*d^2*c)*a+16*d*b^4-4*c^2*b^3-128*d^2*b^2+144*d*c^2*b-27*c^4+256*d^3",
					"(48*c+(-32*a-32*x)*b+9*a^3+12*x*a^2)*d+(-3*a+36*x)*c^2+(4*b^2+(-a^2-28*x*a)*b+6*x*a^3)*c+8*x*b^3-2*x*a^2*b^2",
					"16*d+(-a+12*x)*c+(-2*x*a+8*x^2)*b-3*x^2*a^2",
				},
			}, {
				"a*x^6+b*x^2+c*x+1",
				"6*a*x^5+2*b*x+c",
				"6deg-non-regular",
				[]string{
					"-3125*a^5*c^6+22500*a^5*b*c^4+(-256*a^4*b^5-43200*a^5*b^2)*c^2+1024*a^4*b^6+13824*a^5*b^3+46656*a^6",
					"(3750*a^5*c^4-10800*a^5*b*c^2+512*a^4*b^5+3456*a^5*b^2)*x+4500*a^5*c^3+(256*a^4*b^4-8640*a^5*b)*c",
					"384*a^4*b^3*x^2+480*a^4*b^2*c*x+576*a^4*b^2",
					"0",
					"24*a^2*b*x^2+30*a^2*c*x+36*a^2",
				},
			}, {
				"(x^2+b*x+c) * (a*x-5)^2",
				"(2*x*a^2-10*a)*c+(3*x^2*a^2-20*x*a+25)*b+4*x^3*a^2-30*x^2*a+50*x",
				"6deg-multi",
				[]string{
					"0",
					"(8*x*a^10-40*a^9)*c^3+((-2*x*a^10+10*a^9)*b^2+(80*x*a^9-400*a^8)*b+400*x*a^8-2000*a^7)*c^2+((-20*x*a^9+100*a^8)*b^3+(100*x*a^8-500*a^7)*b^2+(2000*x*a^7-10000*a^6)*b+5000*x*a^6-25000*a^5)*c+(-50*x*a^8+250*a^7)*b^4+(-500*x*a^7+2500*a^6)*b^3+(-1250*x*a^6+6250*a^5)*b^2",
					"((-2*x*a^6+10*a^5)*b+8*x^2*a^6-100*x*a^5+300*a^4)*c+(-3*x^2*a^6+20*x*a^5-25*a^4)*b^2+(-20*x^2*a^5+50*x*a^4+250*a^3)*b-100*x^2*a^4+500*x*a^3",
				},
			}, {
				"1875*5*x^6-22500*a*x^5+22500*a^2*x^4-12000*a^3*x^3+3600*a^4*x^2-576*a^5*x+67500*a",
				"-576*a^5+7200*x*a^4-36000*x^2*a^3+90000*x^3*a^2-112500*x^4*a+56250*x^5",
				"6deg-com",
				[]string{
					"-2644790538240000000000000000000000000*a^30+23245229340000000000000000000000000000000*a^25-81721509398437500000000000000000000000000000*a^20+143651090739440917968750000000000000000000000000*a^15-126255841470211744308471679687500000000000000000000*a^10+44386819266871316358447074890136718750000000000000000*a^5",
					"0", "0", "0",
					"-121500000000*a^6+213574218750000*a",
				},
			},
		} {
			sresTest(g, cas, t, fmt.Sprintf("[%s,%d,%d]", ss.eyec, ii, jj), vv.x, ss.fi, ss.gi, ss.expect)
		}
	}
}
