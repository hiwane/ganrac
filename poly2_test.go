package ganrac_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/hiwane/ganrac"
)

func TestNegX(t *testing.T) {

	funcname := "TestNegX"
	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip %s... (no cas)\n", funcname)
		return
	}
	defer g.Close()

	for j, tt := range []struct {
		varstr string
		lv     Level // x のLevel
	}{
		{"x,y,z", 0},
		{"x,z,y", 0},
		{"y,x,z", 1},
		{"z,x,y", 1},
		{"y,z,x", 2},
		{"z,y,x", 2},
	} {
		for i, ss := range []struct {
			input  string
			expect string
		}{
			{
				"+3*y*x+4*x^2-5*y-11*x^3*y^2+13*x^5*y^3+z^2-y*z-23*x*y*z",
				"-3*y*x+4*x^2-5*y+11*x^3*y^2-13*x^5*y^3+z^2-y*z+23*x*y*z",
			},
		} {
			vstr := fmt.Sprintf("vars(%s);", tt.varstr)
			_, err := g.Eval(strings.NewReader(vstr))
			if err != nil {
				t.Errorf("[%d,%d] `%s` failed: %s", i, j, vstr, err)
				return
			}

			_p, err := g.Eval(strings.NewReader(ss.input + ";"))
			if err != nil {
				t.Errorf("[%d,%d] parse error: `%s`\nerr=%s\ninput=%s.", i, j, vstr, err, ss.input)
				return
			}
			p, ok := _p.(*Poly)
			if !ok {
				t.Errorf("[%d,%d] eval failed\ninput=%s\neval=%v", i, j, ss, _p)
				return
			}
			_q, err := g.Eval(strings.NewReader(ss.expect + ";"))
			if err != nil {
				t.Errorf("[%d,%d] parse error: %s\nerr=%s\nexpect=%s", i, j, vstr, err, ss.expect)
				return
			}
			q, ok := _q.(*Poly)
			if !ok {
				t.Errorf("[%d,%d] eval failed\nexpect=%s\neval=%v", i, j, ss.expect, _q)
				return
			}
			ret := p.NegX(tt.lv)
			if !ret.Equals(q) {
				t.Errorf("[%d,%d] NegX() failed\ninput=%s\nexpect=%s\nactual=%v", i, j, ss.input, ss.expect, ret)
				return
			}
		}
	}
}
