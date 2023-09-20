package ganrac

import (
	"testing"
)

func TestRatOp2(t *testing.T) {
	for _, s := range []struct {
		a1, a2, b1, b2     int64
		add, sub, mul, div RObj
	}{
		{1, 2, 4, 3,
			NewRatInt64(11, 6), NewRatInt64(3-8, 6), NewRatInt64(2, 3), NewRatInt64(3, 8)},
		{4, 3, 2, 3,
			NewInt(2), NewRatInt64(2, 3), NewRatInt64(8, 9), NewInt(2)},
		{4, 3, 1, 3,
			NewRatInt64(5, 3), NewInt(1), NewRatInt64(4, 9), NewInt(4)},
	} {
		a := NewRatInt64(s.a1, s.a2)
		b := NewRatInt64(s.b1, s.b2)

		c := a.Add(b)
		if !c.Equals(s.add) {
			t.Errorf("invalid add a=%v, b=%v, expect=%v, actual=%v\n", a, b, s.add, c)
		}

		c = b.Add(a)
		if !c.Equals(s.add) {
			t.Errorf("invalid add b=%v, a=%v, expect=%v, actual=%v\n", b, a, s.add, c)
		}

		c = a.Sub(b)
		if !c.Equals(s.sub) {
			t.Errorf("invalid sub a=%v, b=%v, expect=%v, actual=%v\n", a, b, s.sub, c)
		}

		c = b.Sub(a)
		if !c.Equals(s.sub.Neg()) {
			t.Errorf("invalid sub b=%v, a=%v, expect=%v, actual=%v\n", b, a, s.sub, c)
		}

		c = a.Mul(b)
		if !c.Equals(s.mul) {
			t.Errorf("invalid mul a=%v, b=%v, expect=%v, actual=%v\n", a, b, s.mul, c)
		}

		c = b.Mul(a)
		if !c.Equals(s.mul) {
			t.Errorf("invalid mul b=%v, a=%v, expect=%v, actual=%v\n", b, a, s.mul, c)
		}

		c = a.Div(b)
		if !c.Equals(s.div) {
			t.Errorf("invalid mul a=%v, b=%v, expect=%v, actual=%v\n", a, b, s.div, c)
		}
	}
}

func TestRatIntOp2(t *testing.T) {
	for _, s := range []struct {
		a1, a2, b          int64
		add, sub, mul, div RObj
	}{
		{1, 3, 2,
			NewRatInt64(7, 3), NewRatInt64(-5, 3), NewRatInt64(2, 3), NewRatInt64(1, 6)},
		{4, 3, 2,
			NewRatInt64(10, 3), NewRatInt64(-2, 3), NewRatInt64(8, 3), NewRatInt64(2, 3)},
		{4, 3, 6,
			NewRatInt64(18+4, 3), NewRatInt64(4-18, 3), NewInt(8), NewRatInt64(2, 9)},
	} {
		a := NewRatInt64(s.a1, s.a2)
		b := NewInt(s.b)

		c := a.Add(b)
		if !c.Equals(s.add) {
			t.Errorf("invalid add a=%v, b=%v, expect=%v, actual=%v\n", a, b, s.add, c)
		}

		c = b.Add(a)
		if !c.Equals(s.add) {
			t.Errorf("invalid add b=%v, a=%v, expect=%v, actual=%v\n", b, a, s.add, c)
		}

		c = a.Sub(b)
		if !c.Equals(s.sub) {
			t.Errorf("invalid sub a=%v, b=%v, expect=%v, actual=%v\n", a, b, s.sub, c)
		}

		c = b.Sub(a)
		if !c.Equals(s.sub.Neg()) {
			t.Errorf("invalid sub b=%v, a=%v, expect=%v, actual=%v\n", b, a, s.sub, c)
		}

		c = a.Mul(b)
		if !c.Equals(s.mul) {
			t.Errorf("invalid mul a=%v, b=%v, expect=%v, actual=%v\n", a, b, s.mul, c)
		}

		c = b.Mul(a)
		if !c.Equals(s.mul) {
			t.Errorf("invalid mul b=%v, a=%v, expect=%v, actual=%v\n", b, a, s.mul, c)
		}

		c = a.Div(b)
		if !c.Equals(s.div) {
			t.Errorf("invalid mul a=%v, b=%v, expect=%v, actual=%v\n", a, b, s.div, c)
		}
	}
}

func TestRatToInt(t *testing.T) {
	for i, s := range []struct {
		num int64
		den int64
	}{
		{5, 4},
		{3, 5},
		{1, 1},
		{7, 1},
		{19, 2},
		{19, 3},
		{19, 4},
		{19, 5},
		{19, 11},
		{19, 31},
		{19, 51},
		{19, 91},
		{675, 2},
		{675, 7},
		{675, 8},
		{675, 9},
		{675, 11},
		{675, 13},
		{675, 14},
		{675, 16},
		{675, 17},
		{675, 19},
		{675, 23},
		{675, 31},
		{675, 53},
		{675, 103},
	} {
		for sgn := int64(-1); sgn <= 1; sgn += 2 {
			b := NewRatInt64(s.num*sgn, s.den)
			actual := b.ToInt()

			// should be actual <= b < actual+1
			if b.Cmp(actual) < 0 {
				t.Errorf("A:i=%d, b=%v, actual=%v", i, b, actual)
				continue
			}

			actual_plus1 := actual.Add(one).(*Int)
			if b.Cmp(actual_plus1) >= 0 {
				t.Errorf("B:i=%d, b=%v, actual=%v", i, b, actual)
				continue
			}
		}
	}
}
