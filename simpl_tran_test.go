package ganrac_test

import (
	"fmt"
	. "github.com/hiwane/ganrac"
	"testing"
)

func TestSimplTran(t *testing.T) {
	funcname := "TestSimplTran"

	g := makeCAS(t)
	if g == nil {
		fmt.Printf("skip %s... (no cas)\n", funcname)
		return
	}
	defer g.Close()

	// evalstr(g, "verbose(3, 0)")

	for ii, ss := range []struct {
		input  string
		expect string
	}{
		{
			"a*x+b > 0",
			"",
		}, {
			"ex([x], all([y], x > y))",
			"false",
		}, {
			"all([x], ex([y], x > y))",
			"true",
		}, {
			"ex([y], x > y)",
			"true",
		}, {
			"ex([z], x > y && y > z && x > z)",
			"x > y",
		}, {
			"ex([z], x > 5*y && y > 3*z && x > 15*z)",
			"x > 5*y",
		}, {
			"all([x,y,z,w,a,b], impl((x != z || y != w) && " +
				"(x - z)^2 + (y - w)^2 == (z - a)^2 + (w - b)^2 && " +
				"(z - a)^2 + (w - b)^2 == (a - x)^2 + (b - y)^2, " +
				"(z - a)^2 + (w - b)^2 == (a - x)^2 + (b - y)^2 + " +
				"(x - z)^2 + (y - w)^2 - 2 * c * " +
				"((a - x)^2 + (b - y)^2)))",
			"2 * c == 1",
		},
	} {
		for jj, vars := range []string{
			//    0 1 2 3 4 5 6 7 8 9 0
			"vars(x,y,z,w,a,b,c,d,e,f,g)",
			"vars(g,f,e,d,c,b,a,w,z,y,x,k1,k2,k3)",
			"vars(z,w,a,x,b,y,c,d,e,f,g)",
		} {

			evalstr(g, vars)

			fof, err := str2fof(g, ss.input)
			if err != nil {
				t.Errorf("[%d,%d] invalid input %s\n%s", ii, jj, err.Error(), ss.input)
				return
			}
			expect, err := str2fof(g, ss.expect)
			if ss.expect != "" && err != nil {
				t.Errorf("[%d,%d] invalid expect %s\n%s", ii, jj, err.Error(), ss.expect)
				return
			}

			qeopt := NewQEopt()
			qeopt.Qe_init(g, fof)
			cond := NewQeCond()
			// fmt.Printf("@@@ GO[%d,%d]! %v\n", ii, jj, fof)
			v := qeopt.Qe_tran(fof, *cond)
			if v == nil || v == fof {
				if ss.expect != "" {
					t.Errorf("[%d] invalid tran1\ninput=%s\nactual=%s\nexpect=%s", ii, ss.input, v, ss.expect)
					break
				}
				continue
			} else if ss.expect == "" {
				t.Errorf("[%d] invalid tran2\ninput=%s\nactual=%s\nexpect=%s", ii, ss.input, v, ss.expect)
				break
			}

			if msg := checkEquivalentQff(g, expect, v); msg != "" {
				t.Errorf("[%d] invalid tran3\ninput=%s\nactual=%s\nexpect=%s\nmsg=%s", ii, ss.input, v, ss.expect, msg)
				return
			}
		}
	}
}
