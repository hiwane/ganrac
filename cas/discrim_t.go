package cas

/*
 * test for discrim()
 */
import (
	"github.com/hiwane/ganrac"
	"testing"
)

func discrimTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T, ii int,
	ps, xs, expect string) {

	p, err := str2poly(g, ps)
	if err != nil {
		t.Errorf("ii=%d, %s input.p error [%s]", ii, ps, err.Error())
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

	act := cas.Discrim(p, x.Level())
	if act == nil {
		t.Errorf("ii=%d, execute failed  discrim(%s,%s).", ii, ps, xs)
		return
	}

	if !exp.Equals(act) {
		t.Errorf("ii=%d, expect=%s, actual=%s.", ii, expect, act)
		return
	}
}

func DiscrimTest(g *ganrac.Ganrac, cas ganrac.CAS, t *testing.T) {
	for ii, ss := range []struct {
		p, x   string
		expect string
	}{
		{
			"a*x^2+b*x+c",
			"x",
			"b^2-4*a*c",
		},
	} {
		discrimTest(g, cas, t, ii, ss.p, ss.x, ss.expect)
	}
}
