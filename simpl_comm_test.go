package ganrac

import (
	"testing"
)

func TestSimplComm1(t *testing.T) {
	x := NewPolyCoef(0, 1, 2)
	y := NewPolyCoef(1, 2, 3)
	z := NewPolyCoef(2, 4, -3)
	w := NewPolyCoef(3, 4, -1)
	X := NewAtom(x, LT)
	Y := NewAtom(y, LE)
	Z := NewAtom(z, EQ)
	W := NewAtom(w, GE)

	for ii, ss := range []struct {
		input  Fof
		expect Fof // simplified input
	}{
		{ // 0
			NewFmlAnd(NewFmlOr(X, Y), NewFmlOr(X, Y)),
			NewFmlOr(X, Y),
		}, { // 1
			NewFmlAnd(NewFmlOr(X, Y), NewFmlOrs(X, Y, Z)),
			NewFmlOrs(X, Y),
		}, { // 2
			NewFmlAnd(NewFmlOrs(X, Y, W, Z), NewFmlOrs(X, Y, Z)),
			NewFmlOrs(X, Y, Z),
		}, { // 3
			NewFmlAnd(NewFmlOrs(W, X, Y, Z), NewFmlOrs(X, Y, Z)),
			NewFmlOrs(X, Y, Z),
		}, { // 4
			NewFmlAnd(NewFmlOrs(X, Y, NewAtom(z, GT)), NewFmlOrs(X, Y, NewAtom(z, GE))),
			NewFmlOrs(X, Y, NewAtom(z, GT)),
		}, {
			NewFmlAnd(X, NewFmlOrs(Z, NewFmlAnd(X, Y))),
			NewFmlAnd(X, NewFmlOrs(Z, Y)),
		}, {
			NewFmlAnd(X, NewFmlOr(X, Y)),
			X,
		}, {
			NewFmlAnd(X, NewFmlOrs(X, Y, Z)),
			X,
		},
	} {
		if ss.expect == nil {
			ss.expect = ss.input
		}
		for i, s := range []struct {
			input  Fof
			expect Fof // simplified input
		}{
			{ss.input, ss.expect},
			{ss.input.Not(), ss.expect.Not()},
		} {
			output := s.input.simplComm()
			if TTestSameFormAndOr(output, s.expect) {
				continue
			}

			out2 := output.simplBasic(trueObj, falseObj)
			if TTestSameFormAndOr(out2, s.expect) {
				continue
			}

			t.Errorf("%d/%d: not same form:\ninput =`%v`\noutput=`%v`\noutput=`%v`\nexpect=`%v`", ii, i, s.input, output, out2, s.expect)
			return
		}
	}
}
