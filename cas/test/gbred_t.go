package cas

/*
 * test for GB(), Reduce()
 */
import (
	"github.com/hiwane/ganrac"
	"testing"
)

func gbRedTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T,
	ii int, input, varstr string, n int) {
	polys, err := str2list(g, input)
	if err != nil {
		t.Errorf("ii=%d, %s input error [%s]", ii, input, err.Error())
		return
	}

	vars, err := str2list(g, varstr)
	if err != nil {
		t.Errorf("ii=%d, %s vars error [%s]", ii, vars, err.Error())
		return
	}

	gb := cas.GB(polys, vars, n)
	if gb == nil {
		t.Errorf("ii=%d, GB returns nil [%s]", ii, vars)
		return
	}

	fs := make([]*ganrac.Poly, 3)
	for i := 0; i < len(fs); i++ {
		fi, err := gb.Geti(i)
		if err != nil {
			fi = ganrac.NewInt(0)
		}
		fs[i] = fi.(*ganrac.Poly)
	}

	for jj, pstr := range []string{
		"x+2*y-3",
		"x+y-3*z-7",
		"-4*z"} {
		p, err := str2poly(g, pstr)
		if err != nil {
			t.Errorf("ii=%d, jj=%d, invalid pstr=%s, err=%s", ii, jj, pstr, err.Error())
			return
		}

		for kk, tbl := range [][]int64{
			{1, 2, 3},
			{0, -3, 1},
			{-4, -5, 3},
		} {
			var q ganrac.RObj = p
			for j, s := range tbl {
				q = ganrac.Add(q, ganrac.Mul(ganrac.NewInt(s), fs[j]))
			}

			u, _ := cas.Reduce(q.(*ganrac.Poly), gb, vars, n)
			if up, ok := u.(*ganrac.Poly); ok {
				u, _ = up.PPC()
			}
			px, _ := p.PPC()

			if !u.Equals(px) {
				t.Errorf("ii=(%d,%d,%d)\np=%s\ngb=%s\nexpect=%s\nactual=%s", ii, jj, kk, px, gb, pstr, u)
			}
		}
	}
}

func GBRedTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T) {
	for ii, ss := range []struct {
		input string
		vars  string
		n     int
	}{
		{
			"[x^5 + y^4 + z^3 - 1,  x^3 + y^3 + z^2 - 1]",
			"[x, y, z]",
			0,
		},
	} {
		gbRedTest(g, cas, t, ii, ss.input, ss.vars, ss.n)
	}
}
