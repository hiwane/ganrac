package ganrac_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/hiwane/ganrac"
)

func TestPolyNegX(t *testing.T) {
	funcname := "TestPolyNegX"
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

func TestPolySetZero(t *testing.T) {
	funcname := "TestPolySetZero"
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
		{"z,y,x", 2},
		{"x,y,z", 0},
		{"x,z,y", 0},
		{"y,x,z", 1},
		{"z,x,y", 1},
		{"y,z,x", 2},
	} {
		for i, ss := range []struct {
			input  string
			expect []string
		}{
			{
				"(3*y-5*z)*x^3+(4*z^2-11*y^2)*x-13*y*z",
				[]string{
					"(3*y-5*z)*x^3+(4*z^2-11*y^2)*x",
					"(3*y-5*z)*x^3+(4*z^2-11*y^2)*0-13*y*z",
					"",
					"0*(3*y-5*z)*x^3+(4*z^2-11*y^2)*x-13*y*z",
					"",
					"",
				},
			}, {
				"(-3)*x^4+(5*y^2-3*z-4)*x^3+(4*z^5-7*y^2)*x^2-5",
				[]string{
					"(-3)*x^4+(5*y^2-3*z-4)*x^3+(4*z^5-7*y^2)*x^2-0",
					"",
					"(-3)*x^4+(5*y^2-3*z-4)*x^3+0*(4*z^5-7*y^2)*x^2-5",
					"(-3)*x^4+0*(5*y^2-3*z-4)*x^3+(4*z^5-7*y^2)*x^2-5",
					"0*(-3)*x^4+(5*y^2-3*z-4)*x^3+(4*z^5-7*y^2)*x^2-5",
					"",
					"",
				},
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
				t.Errorf("[%d,%d] eval failed\ninput=%s\neval =%v", i, j, ss, _p)
				return
			}

			for deg, expect := range ss.expect {
				if expect == "" {
					expect = ss.input
				}

				q, err := g.Eval(strings.NewReader(expect + ";"))
				if err != nil {
					t.Errorf("[%d,%d,%d] parse error: %s\nerr=%s\nexpect=%s", i, j, deg, vstr, err, expect)
					return
				}
				ret := p.SetZero(tt.lv, deg)
				if !ret.Equals(q) {
					t.Errorf("[%d,%d,%d] %s() failed\ninput =%v = %s\nexpect=%v = %s\nactual=%v", i, j, deg, funcname, p, ss.input, q, expect, ret)
					return
				}
			}
		}
	}
}
func TestSubstXinvLv(t *testing.T) {
	funcname := "TestSubstXinvLv"
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
		{"z,y,x", 2},
		{"x,y,z", 0},
		{"x,z,y", 0},
		{"y,x,z", 1},
		{"z,x,y", 1},
		{"y,z,x", 2},
	} {
		vstr := fmt.Sprintf("vars(%s);", tt.varstr)
		_, err := g.Eval(strings.NewReader(vstr))
		if err != nil {
			t.Errorf("[%d] `%s` failed: %s", j, vstr, err)
			return
		}

		for i, ss := range []struct {
			input  string
			expect string
		}{
			{
				"5*x^2-3*x-1",
				"-x^2-3*x+5",
			}, {
				"-7*x",
				"-7",
			}, {
				"7*y*(z-11)*x",
				"7*y*(z-11)",
			}, {
				"(3*y-5*z)*x^3+(4*z^2-11*y^2)*x-13*y*z",
				"(3*y-5*z)*1^3+(4*z^2-11*y^2)*x^2-13*y*z*x^3",
			}, {
				"(-3)*x^4+(5*y^2-3*z-4)*x^3+(4*z^5-7*y^2)*x^2-5",
				"(-3)*1^4+(5*y^2-3*z-4)*x^1+(4*z^5-7*y^2)*x^2-5*x^4",
			},
		} {

			_p, err := g.Eval(strings.NewReader(ss.input + ";"))
			if err != nil {
				t.Errorf("[%d,%d] p parse error: `%s`\nerr=%s\ninput=%s.", i, j, vstr, err, ss.input)
				return
			}
			p, ok := _p.(*Poly)
			if !ok {
				t.Errorf("[%d,%d] p eval failed\ninput=%s\neval =%v", i, j, ss, _p)
				return
			}

			_q, err := g.Eval(strings.NewReader(ss.expect + ";"))
			if err != nil {
				t.Errorf("[%d,%d] q parse error: %s\nerr=%s\nexpect=%s", i, j, vstr, err, ss.expect)
				return
			}
			q, ok := _q.(RObj)
			if !ok {
				t.Errorf("[%d,%d] q eval failed\ninput=%s\neval =%v", i, j, ss, _q)
				return
			}

			n := p.Deg(tt.lv)
			for k := 0; k < 3; k++ {
				ret := p.SubstXinvLv(tt.lv, n)
				if err := ValidRObj(ret); err != nil {
					t.Errorf("[%d,%d,%d,%d] %s() invalid\ninput =%v = %s\nexpect=%v = (%s)*x^%d\nactual=%v", i, j, k, n, funcname, p, ss.input, q, ss.expect, k, ret)
					return
				}

				if !ret.Equals(q) {
					t.Errorf("[%d,%d,%d,%d] %s() failed\ninput =%v = %s\nexpect=%v = (%s)*x^%d\nactual=%v", i, j, k, n, funcname, p, ss.input, q, ss.expect, k, ret)
					return
				}
				n++
				q = Mul(q, NewPolyVar(tt.lv))
			}
		}
	}
}
