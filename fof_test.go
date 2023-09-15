package ganrac

import (
	"testing"
)

func TestAtom(t *testing.T) {
	for _, s := range []struct {
		op1, op2 OP
	}{
		{LE, GT}, // 3, 4
		{LT, GE}, // 1, 6
		{EQ, NE}, // 2, 5
	} {
		pp := NewAtom(NewPolyCoef(0, 0, 1), s.op1)
		p, ok := pp.(*Atom)
		if !ok {
			t.Errorf("invalid atom %v", pp)
			return
		}

		if _, ok = pp.(Fof); !ok {
			t.Errorf("not Fof %v", pp)
			return
		}

		qq := pp.Not()
		q, ok := qq.(*Atom)
		if !ok {
			t.Errorf("invalid atom not(%v)=%v", pp, qq)
			return
		}

		if len(p.p) != len(q.p) || q.op != s.op2 {
			t.Errorf("invalid atom not(%v)=%v", pp, qq)
			return
		}

		for i := 0; i < len(p.p); i++ {
			if p.p[i] != q.p[i] {
				t.Errorf("invalid atom [%d] not(%v)=%v", i, pp, qq)
				return
			}
		}

		rr := q.Not()
		r, ok := rr.(*Atom)
		if !ok {
			t.Errorf("invalid atom not(%v)=%v", qq, rr)
			return
		}

		if len(p.p) != len(r.p) || r.op != p.op {
			t.Errorf("invalid atom not(%v)=%v", qq, rr)
			return
		}

		for i := 0; i < len(p.p); i++ {
			if p.p[i] != r.p[i] {
				t.Errorf("invalid atom [%d] not(%v)=%v", i, pp, rr)
				return
			}
		}
	}
}

func TestFmlAnd(t *testing.T) {
	fmls := []Fof{
		NewAtom(NewPolyCoef(0, 1, 2, 3), GE),
		NewAtom(NewPolyCoef(0, 2, 3, 4), NE),
		NewBool(true),
		NewAtom(NewPolyCoef(1, 5, 1, 2), LT),
		NewAtom(NewPolyCoef(2, 1, 1, 2), LT),
	}

	var f Fof = NewFmlAnd(fmls[0], fmls[1])
	var ans Fof = NewFmlAnd(f, fmls[3])
	ans = NewFmlAnd(ans, fmls[4])
	if err := ans.valid(); err != nil {
		t.Errorf("ans %s", ans.valid())
	}

	f = NewFmlAnd(fmls[0], fmls[1])
	g := NewFmlAnd(fmls[2], fmls[3])
	g = NewFmlAnd(g, fmls[4])
	h := NewFmlAnd(f, g)
	if !ans.Equals(h) {
		t.Errorf("(0 && 1) && ((2 && 3) && 4)\n%v\n%v", h, ans)
	}

	f = fmls[1]
	f = NewFmlAnd(f, fmls[2])
	f = NewFmlAnd(f, fmls[3])
	f = NewFmlAnd(f, fmls[4])
	h = NewFmlAnd(fmls[0], f)
	if !ans.Equals(h) {
		t.Errorf("0 && (((1 && 2) && 3) && 4)\nh=%v\na=%v", h, ans)
	}

	h = NewFmlAnd(h, NewBool(false))
	if !h.Equals(falseObj) {
		t.Errorf("not false: %v", h)
	}
}

func TestFmlOr(t *testing.T) {
	fmls := []Fof{
		NewAtom(NewPolyCoef(0, 1, 2, 3), GE),
		NewAtom(NewPolyCoef(0, 2, 3, 4), NE),
		NewBool(false),
		NewAtom(NewPolyCoef(1, 5, 1, 2), LT),
		NewAtom(NewPolyCoef(2, 1, 1, 2), LE),
	}

	var f Fof = NewFmlOr(fmls[0], fmls[1])
	var ans Fof = NewFmlOr(f, fmls[3])
	ans = NewFmlOr(ans, fmls[4])
	if err := ans.valid(); err != nil {
		t.Errorf("ans %s", ans.valid())
	}

	f = NewFmlOr(fmls[0], fmls[1])
	g := NewFmlOr(fmls[2], fmls[3])
	g = NewFmlOr(g, fmls[4])
	h := NewFmlOr(f, g)
	if !ans.Equals(h) {
		t.Errorf("(0 || 1) || ((2 || 3) || 4)\nactual=%v\nexpect=%v", h, ans)
	}

	f = fmls[1]
	f = NewFmlOr(f, fmls[2])
	f = NewFmlOr(f, fmls[3])
	f = NewFmlOr(f, fmls[4])
	h = NewFmlOr(fmls[0], f)
	if !ans.Equals(h) {
		t.Errorf("0 || (((1 || 2) || 3) || 4)\nactual=%v\nexpect=%v", h, ans)
	}

	h = NewFmlOr(h, NewBool(true))
	if !h.Equals(trueObj) {
		t.Errorf("not false: %v", h)
	}
}

func TestOp(t *testing.T) {
	for ii, ss := range []struct {
		input  OP
		not    OP
		neg    OP
		strict OP
	}{
		{LT, GE, GT, LT},
		{LE, GT, GE, LT},
		{EQ, NE, EQ, OP_FALSE},
		{GT, LE, LT, GT},
		{GE, LT, LE, GT},
		{NE, EQ, NE, NE},
	} {
		oo := ss.input.not()
		if oo != ss.not {
			t.Errorf("[%d] input=%d, actual=%d, expect=%d, %s", ii, ss.input, oo, ss.not, "not")
		}

		oo = ss.input.neg()
		if oo != ss.neg {
			t.Errorf("[%d] input=%d, actual=%d, expect=%d, %s", ii, ss.input, oo, ss.neg, "neg")
		}

		oo = ss.input.strict()
		if oo != ss.strict {
			t.Errorf("[%d] input=%d, actual=%d, expect=%d, %s", ii, ss.input, oo, ss.strict, "strict")
		}

	}
}
