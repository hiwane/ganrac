package main

import (
	"testing"
)

func TestSdcDelta(t *testing.T) {
	var st *sgn_table
	st = new(sgn_table)

	for i, expect := range []sgn_t{
		1, -1, -1, 1,
		1, -1, -1, 1,
		1, -1, -1, 1,
	} {
		if actual := st.delta(i); actual != expect {
			t.Errorf("delta(%d) = %d, want %d", i, actual, expect)
		}
	}
}

func TestSdcW(t *testing.T) {
	var st *sgn_table
	st = new(sgn_table)
	for i, ss := range []struct {
		s []sgn_t
		c []sgn_t
		n int
	}{
		{
			[]sgn_t{1, 1, 1, 1},
			[]sgn_t{1, 0, 0, 1},
			2,
		}, {
			[]sgn_t{1, 0, 0, 1},
			[]sgn_t{1, 0, 0, 1},
			0,
		}, {
			[]sgn_t{1, 1, 0, 1},
			[]sgn_t{1, -1, 0, 1},
			2,
		}, {
			[]sgn_t{1, 1, 1, 1},
			[]sgn_t{1, -1, 1, -1},
			3,
		}, {
			[]sgn_t{1, 1, 0, 1},
			[]sgn_t{1, -1, 0, 1},
			2,
		}, {
			[]sgn_t{1, 1, -1},
			[]sgn_t{1, 0, -1},
			1,
		}, {
			[]sgn_t{+1, +1, -1, 0},
			[]sgn_t{+1, +1, +1, 0},
			0,
		},
	} {
		st.s = ss.s
		st.c = ss.c
		if len(ss.s) != len(ss.c) {
			t.Errorf("i=%d test setting is invalid\ns=%v\nc=%v\n", i, ss.s, ss.c)
		}
		st.deg = len(ss.s) - 1

		if st.W(true) != ss.n {
			t.Errorf("i=%d\ns=%v\nc=%v\nW() = %d, want %d", i, st.s, st.c, st.W(true), ss.n)
		}
	}
}

func TestSdcV(t *testing.T) {
	var st *sgn_table
	st = new(sgn_table)
	for i, ss := range []struct {
		v []sgn_t
		n int
	}{
		{
			[]sgn_t{1, 1, 1, 1}, 0,
		}, {
			[]sgn_t{1, 1, 1, -1}, 1,
		}, {
			[]sgn_t{1, 1, -1, 1}, 2,
		}, {
			[]sgn_t{1, 0, 0, 1, 0, 1}, 0,
		}, {
			[]sgn_t{1, 0, 0, -1, 0, 1}, 2,
		}, {
			[]sgn_t{1, 0, 1, -1, 1, 1}, 2,
		}, {
			[]sgn_t{1, 1, -1, 0}, 1,
		}, {
			[]sgn_t{0, -1, 1, 1}, 1,
		},
	} {
		if st.V(ss.v) != ss.n {
			t.Errorf("i=%d\nv=%v\nV() = %d, want %d", i, ss.v, st.V(ss.v), ss.n)
		}
	}
}
