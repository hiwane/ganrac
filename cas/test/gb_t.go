package cas

/*
 * test for
 *    GB(p *List, vars *List, n int) *List
 */
import (
	"fmt"
	"github.com/hiwane/ganrac"
	"testing"
)

func containsLvs(lvs []ganrac.Level, v ganrac.Level) bool {
	for _, l := range lvs {
		if l == v {
			return true
		}
	}
	return false
}

func selectInSubring(gb *ganrac.List, vars *ganrac.List, n int) []*ganrac.Poly {
	ret := make([]*ganrac.Poly, 0)

	lvs := make([]ganrac.Level, 0)
	for i, _v := range vars.Iter() {
		if i >= n {
			v := _v.(*ganrac.Poly)
			lvs = append(lvs, v.Level())
		}
	}
	fmt.Printf("selectInSubring lvs=%v\n", lvs)

	for _, _g := range gb.Iter() {
		g := _g.(*ganrac.Poly)
		vok := true
		for lv := g.Level(); lv >= 0; lv-- {
			if !containsLvs(lvs, lv) && g.Deg(lv) > 0 {
				vok = false
				fmt.Printf("selectInSubring g=%v, lv=%d, deg=%d, lvs=%v\n", g, lv, g.Deg(lv), lvs)
				break
			}
		}
		if vok {
			ret = append(ret, g)
			fmt.Printf("selectInSubring append(%v)\n", g)
			ret = append(ret, g.Neg().(*ganrac.Poly))
		}
	}

	return ret
}

func gbTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T, ii int,
	_polys, _vars string, n int, _expect string) {

	polys, err := str2list(g, _polys)
	if err != nil {
		t.Errorf("ii=%d, %s input.polys error [%s]", ii, _polys, err.Error())
		return
	}
	vars, err := str2list(g, _vars)
	if err != nil {
		t.Errorf("ii=%d, %s input.vars error [%s]", ii, _vars, err.Error())
		return
	}
	expect, err := str2list(g, _expect)
	if err != nil {
		t.Errorf("ii=%d, %s input.expect error [%s]", ii, _expect, err.Error())
		return
	}

	fmt.Printf("[%d] gbtest start! @@@@@@@@@@@@@@@@@@@@@@ %v[%d]\n", ii, vars, n)
	gb := cas.GB(polys, vars, n)
	fmt.Printf("gb=%v\n", gb)
	gbsub := selectInSubring(gb, vars, n)
	fmt.Printf("gbsub=%v\n", gbsub)
	for jj, _exp := range expect.Iter() {
		valid := false
		exp := _exp.(*ganrac.Poly)
		for _, g := range gbsub {
			if g.Equals(exp) {
				valid = true
				break
			}
		}
		if !valid {
			t.Errorf("[%d,%d] %s is not found: ret=%s", ii, jj, exp, gbsub)
			return
		}
	}
}

func GBTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T) {
	for ii, ss := range []struct {
		polys  string
		vars   string
		expect []string
	}{
		{
			"[x^2-z, x*y-1, x^3-x^2*y-x^2-1]",
			"[x, y, z]",
			[]string{
				"[-x-y+z-1,2*y+z^2-4*z+1,(z+1)*y-z+1,y^2+2*y-z+2]",
				"[2*y+z^2-4*z+1,(-z-1)*y+z-1,-y^2-2*y+z-2]",
				"[z^3-3*z^2-z-1]",
			},
		},
	} {
		for jj, tt := range ss.expect {
			if tt == "" {
				continue
			}
			gbTest(g, cas, t, 1000*(ii+1)+jj, ss.polys, ss.vars, jj, tt)
		}
	}
}
