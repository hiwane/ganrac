package main

// espresso の入力ファイルを生成する
//
// An Effective Implementation of a Special Quantifier Elimination for a Sign Definite Condition by Logical Formula Simplification
// H. Iwane, H. Higuchi, H. Anai (2013)
//
// の改良版.
//
// SDC: all([x], impl(x >= 0, f(x) > 0))
// ==>  all([x], x < 0 || f(x) > 0)
//
// のことだが，ここでは，ex() の問題として定義する．つまり，否定をとった
//
//
// negSDC: ex([x], x >= 0 && f(x) <= 0) を，DNF 形式で復帰することを考える
//
// | deg |  t |  True | False | rims2017 | casc2013 |
// |   2 |  2 |     5 |     2 |        2 |        2 |
// |   3 |  4 |    17 |     8 |        2 |        4 |
// |   4 |  6 |   117 |    78 |        6 |       10 |
// |   5 |  8 |   425 |   306 |        7 |       18 |
// |   6 | 10 |  3089 |  2660 |       31 |       57 |
// |   7 | 12 | 11897 | 10838 |       42 |      121 |
// |   8 | 14 | 87361 | 89592 |      134 |      353 |
//
//
// SH-Theorem: 	QE&CAD
// p.303 Theorem2 Sturm-Habicht Structure Theorem
//
// espresso -Dexact -epos 4.in : 5
// espresso -estrong 8.in : 5

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

const DC = 2

type sgn_t int
type val_t int

var sgns []sgn_t

func init() {
	sgns = []sgn_t{-1, 0, 1}
}

type sgn_table struct {
	deg       int
	s         []sgn_t // sign(f(+inf))
	h         []sgn_t // sign(coeff(f, j)): principal coefficient
	c         []sgn_t // sign(coeff(f, 0)): constant term
	count     []int
	fp        *bufio.Writer
	debug     bool
	val_equal val_t
	val_not   val_t
	pfalse    bool
	atom      bool
}

func (st *sgn_table) delta(n int) sgn_t {
	return sgn_t(-1 + ((n - 1) & 2))
}

func (st *sgn_table) VV() int {
	s := make([]sgn_t, 0, len(st.s))

	var last int = -1
	for i := 0; i < len(st.h); i++ {
		var d int
		var c sgn_t
		if st.h[i] != 0 {
			c = st.h[i]
			d = i
			last = i
		} else if st.s[i] == 0 {
			continue
		} else {
			d = last
			c = st.s[i]
		}
		if d%2 != 0 {
			c = -c
		}
		s = append(s, c)
	}

	return st.V(s)
}

// # of sign variation
func (st *sgn_table) V(v []sgn_t) int {
	n := 0
	var i int
	for i = 0; i < len(v); i++ {
		if v[i] != 0 {
			break
		}
	}
	s := v[i]
	for i++; i < len(v); i++ {
		if v[i] == 0 {
			continue
		} else if v[i]*s < 0 {
			s = v[i]
			n++
		}
	}
	return n
}

// # of sign variation at x = 0
func (st *sgn_table) W(debug bool) int {
	c := make([]sgn_t, 0, st.deg+1)
	for i := 0; i <= st.deg; i++ {
		if st.s[i] != 0 { // not vanish
			c = append(c, st.c[i])
		}
	}

	n := 0
	s := c[0]
	zero := 0
	for i := 1; i < len(c); i++ {
		if c[i] == 0 {
			zero++
		} else {
			switch zero {
			case 0, 1:
				if s != c[i] {
					n++
				} else if zero == 1 {
					if !debug {
						if !st.debug {
							st.debug = true
						}
						st.print(9, "+0+, -0- はないって")
						panic("why!")
					}
				}
			case 2:
				if s == c[i] {
					n += 2
				} else {
					n++
				}
			default:
				if !debug {
					panic("Why!")
				}
			}
			zero = 0
			s = c[i]
		}
	}
	return n
}

// atom 用
func (st *sgn_table) evalAtom() (val_t, string) {
	ks := 0
	var ask sgn_t = 0 // 非ゼロの右端
	var m int = 0
	var p int = 0
	var c int = 0
	var e int = 0
	for i := len(st.h) - 2; i >= 0; i-- {
		a := st.h[i]
		if a == 0 {
			if ks == 0 {
				ask = st.h[i+1]
			}
			ks++
		} else if ks == 0 {
			if a == st.h[i+1] {
				m++
				p++
			} else if st.h[i+1] == 0 {
				panic("why!")
			} else {
				m--
				c++
			}
		} else {
			if ks%2 == 0 {
				var u int = 1
				if (ks/2)%2 == 1 {
					u = -1
				}
				u = u * int(ask*a)
				m += u
				e += u
			}
			ks = 0
		}
	}

	mes := fmt.Sprintf(": m(%d) = p(%d) - c(%d) + e(%d)", m, p, c, e)
	if m > 0 {
		return st.val_not, mes
	} else if m < 0 {
		return 3, mes
	} else {
		return st.val_equal, mes
	}
}

func (st *sgn_table) evalSdc() (val_t, string) {
	n := st.deg
	// 主係数版固有. h[0] != 0, h[1] = ... = h[m-1] = 0 h[m] != 0 ==> h[m] = c[m]
	var ret val_t = 3

	// Lemma3, SH-Theorem (A)
	u := 0
	for u = 0; u <= n; u++ {
		if st.h[u] != 0 {
			if st.c[u] == 0 {
				return ret, "lem3"
			}
			break
		}
	}

	// Lemma2: s[i] = 0 => c[i] = 0
	for i := u + 1; i < n; i++ {
		if st.h[i+1] == 0 && st.h[i] == 0 && st.c[i] != 0 {
			return ret, "lem2"
		}
	}

	for j := n - 1; j >= u; j-- {
		if st.h[j+1] != 0 {
			for r := j; ; r-- {
				if st.h[r] != 0 {
					if st.h[j] == 0 {
						// SH-Theorem(B)
						if (j-r)%2 == 0 && st.c[r] != st.delta(j-r)*st.c[j] {
							return ret, "thmB (even c)"
						}
						if st.c[r]*st.c[j] == 0 && st.c[r] != st.c[j] {
							// 一方がゼロなら他方もゼロ
							return ret, "thmB (zero c)"
						}
						if (j-r)%2 != 0 && st.h[j+1]*st.h[r] != st.delta(j-r) {
							return ret, "thmB (odd h)"
						}
					}
					if st.c[j] == 0 && r > 0 {
						hj1 := sgn_t(1)
						if (j-r)%2 != 0 {
							hj1 = st.h[j+1]
						}

						if (j-r)%2 == 0 && hj1*st.c[r-1] != st.delta(j-r+2)*st.c[j+1] {
							return ret, "thmC"
						} else if st.c[r-1]*st.c[j+1] == 0 && st.c[r-1] != st.c[j+1] {
							return ret, "thmC2"
						}
					}
					break
				}
			}
		}
	}

	st.set_s()

	w0 := st.W(false)
	wn := st.VV()

	if wn-w0 < 0 {
		return ret, "VV-W < 0"
	}
	if wn-w0 == 0 && n%2 != 0 {
		// 奇数次なら，負の実根をもつ
		return ret, "VV-W = 0 (deg=odd)"
	}

	wp := st.V(st.s)
	if mm := w0 - wp; mm > 0 {
		return st.val_not, ""
	} else if mm == 0 {
		return st.val_equal, ""
	} else {
		return ret, "W-V < 0"
	}
}

func (st *sgn_table) set_s() {
	// h/c の符号から s の符号が確定する
	var j int
	j = st.deg
	var u int = 0
	for st.h[u] == 0 {
		st.s[u] = 0
		u++
	}

	for r := st.deg - 1; r >= u; r-- {
		if st.h[r] != 0 {
			st.s[r] = st.h[r]
			if st.h[r+1] == 0 {
				switch (j - r) % 4 {
				case 0:
					st.s[j] = st.h[r]
				case 2:
					st.s[j] = -st.h[r]
				default:
					if st.c[j] != 0 {
						st.s[j] = st.delta(j-r) * st.h[j+1] * st.c[j] * st.c[r]
					} else {
						if r == 0 {
							st.print(9, fmt.Sprintf("j=%d, r=%d", j, r))
						}
						st.s[j] = st.delta(j-r+2) * st.h[j+1] * st.c[j+1] * st.c[r-1]
					}
				}
			}
			j = r - 1
		} else if st.h[r+1] == 0 {
			st.s[r] = 0
		} else {
			st.s[r] = DC
		}
	}
	for i := 0; i < st.deg+1; i++ {
		if st.s[i] == DC {
			st.print(9, "s[i] == DC")
			panic("s[i] == DC")
		}
		if st.h[i] != 0 && st.s[i] != st.h[i] {
			st.print(9, "s[i] != h[i]")
			panic("s[i] != h[i]")
		}
	}
}

func (st *sgn_table) header() {
	n := st.deg
	if st.debug {
		fmt.Fprintf(st.fp, "h")
		for i := n; i >= 0; i-- {
			fmt.Fprintf(st.fp, "%d", i)
		}
		if !st.atom {
			fmt.Fprintf(st.fp, " s")
			for i := n; i >= 0; i-- {
				fmt.Fprintf(st.fp, "%d", i)
			}
			fmt.Fprintf(st.fp, " c")
			for i := n; i >= 0; i-- {
				fmt.Fprintf(st.fp, "%d", i)
			}
		}
		fmt.Fprintf(st.fp, " i")
		for i := n - 2; i >= 0; i-- {
			fmt.Fprintf(st.fp, "%d", i)
		}
		for i := n - 1; i > 0; i-- {
			fmt.Fprintf(st.fp, "%d", i)
		}

		fmt.Fprintf(st.fp, "  t 0 +inf")
		fmt.Fprintf(st.fp, "\n")
	} else {
		if st.atom {
			fmt.Fprintf(st.fp, ".i %d\n", 2*(n-1))
		} else {
			fmt.Fprintf(st.fp, ".i %d\n", 4*(n-1))
		}
		fmt.Fprintf(st.fp, ".o 1\n")
		fmt.Fprintf(st.fp, ".lib")
		for i := n - 2; i >= 0; i-- {
			fmt.Fprintf(st.fp, " h%d i%d", i, i)
		}
		if !st.atom {
			for i := n - 1; i > 0; i-- {
				fmt.Fprintf(st.fp, " c%d d%d", i, i)
			}
		}
		fmt.Fprintf(st.fp, "\n.ob f0\n")
	}
}

func (st *sgn_table) footer() {
	fmt.Fprintf(st.fp, ".e\n")
}

func (st *sgn_table) print(t val_t, mes string) {
	n := st.deg
	if st.debug {
		const str string = "-0+*@"
		fmt.Fprintf(st.fp, " ")
		for i := n; i >= 0; i-- {
			fmt.Fprintf(st.fp, "%c", str[st.h[i]+1])
		}
		if !st.atom {
			fmt.Fprintf(st.fp, "  ")
			for i := n; i >= 0; i-- {
				fmt.Fprintf(st.fp, "%c", str[st.s[i]+1])
			}
			fmt.Fprintf(st.fp, "  ")
			for i := n; i >= 0; i-- {
				fmt.Fprintf(st.fp, "%c", str[st.c[i]+1])
			}
		}
		var tc string
		switch t {
		case 0:
			tc = "0"
		case 1:
			tc = "1"
		case 2:
			tc = "_"
			tc = "2"
		case 3:
			tc = " "
			tc = "2"
		default:
			tc = "?"
		}

		if !st.atom {
			fmt.Fprintf(st.fp, " %s %d %d  %s\n", tc, st.W(true), st.V(st.s), mes)
		} else {
			fmt.Fprintf(st.fp, " %s %s\n", tc, mes)
		}
	} else {
		ss := []string{"01", "00", "10", "--", "11"}

		for i := n - 2; i >= 0; i-- {
			fmt.Fprintf(st.fp, "%s", ss[1+st.h[i]])
		}
		if !st.atom {
			for i := n - 1; i > 0; i-- {
				fmt.Fprintf(st.fp, "%s", ss[1+st.c[i]])
			}
		}
		if t > 2 {
			t = 2
		}
		fmt.Fprintf(st.fp, " %d\n", t)
	}
}

func (st *sgn_table) loopc(n int) {
	if n == st.deg {
		t, mes := st.evalSdc()
		if t == 0 || t == 1 {
			st.count[t]++
		}
		if 0 < t || st.pfalse {
			st.print(t, mes)
		}
		return
	}
	if st.h[n] == 0 && st.h[n+1] == 0 { // @DC2
		st.c[n] = 0
		st.loopc(n + 1)
		return
	}

	for _, i := range sgns { // 符号
		st.c[n] = i
		st.loopc(n + 1)
	}
}

func (st *sgn_table) loops() {
	for c := 0; c <= st.deg; c++ {
		st.c[c] = st.h[c]
		if st.h[c] != 0 {
			// 初めて非ゼロとなった... GCD がみつかった
			// 初めて非ゼロになるまでは, vanish 確定なので， c[c] = 0 になるので skip @DC1
			if c == 0 {
				// SH_0 は定数なので， h[0] と符号が一致するため skip
				c++
			}
			st.loopc(c)
			return
		}
	}
}

func (st *sgn_table) looph(n int) {
	if n == 0 {
		// constant term
		for _, i := range sgns { // 符号
			st.h[n] = i
			st.c[n] = i
			st.looph(n + 1)
		}
	} else if n == st.deg-1 {
		// head term
		if st.atom {
			t, mes := st.evalAtom()
			if t == 0 || t == 1 {
				st.count[t]++
			}
			if 0 < t || st.pfalse {
				st.print(t, mes)
			}
			return
		} else {
			st.loops()
		}
	} else {
		for _, i := range sgns { // 符号
			st.h[n] = i
			st.looph(n + 1)
		}
	}
}

func (st *sgn_table) set_dc() {
	for i := 0; i < st.deg-1; i++ {
		st.s[i] = DC
		st.h[i] = DC
		st.c[i+1] = DC
	}
	st.c[0] = DC
}

func (st *sgn_table) print_dc() {
	st.set_dc()
	for i := 0; i < st.deg-1; i++ {
		st.h[i] = 3
		st.print(DC, "DC 11h")
		st.h[i] = DC
	}
	if st.atom {
		return
	}

	st.set_dc()
	for i := 0; i < st.deg-1; i++ {
		st.c[i+1] = 3
		st.print(DC, "DC 11c")
		st.c[i+1] = DC
	}

	// @DC1
	st.set_dc()
	st.h[0] = 0
	st.c[0] = 0
	for i := 1; i < st.deg-1; i++ {
		st.h[i] = 0
		st.c[i] = 1
		st.print(DC, "DC1+")
		st.c[i] = -1
		st.print(DC, "DC1-")
		st.c[i] = DC
	}

	// @DC2
	st.set_dc()
	for i := 1; i < st.deg-2; i++ {
		st.h[i] = 0
		st.h[i+1] = 0
		st.c[i] = 1
		st.print(DC, "DC2+")
		st.c[i] = -1
		st.print(DC, "DC2-")
		st.c[i] = DC
		st.h[i] = DC
	}

}

func (st *sgn_table) gen() {
	// 必要条件
	deg := st.deg
	st.s[deg] = 1
	st.s[deg-1] = 1
	st.h[deg] = 1
	st.h[deg-1] = 1
	st.c[deg] = 1

	st.header()
	st.looph(0)
	st.print_dc()
	st.footer()
}

func main() {
	var (
		deg         = flag.Int("d", 2, "degree")
		debug       = flag.Bool("debug", false, "debug print mode")
		false_print = flag.Bool("false", false, "print false")
		all         = flag.Bool("all", false, "all([x], x >= 0 impl f(x) > 0)")
		atom        = flag.Bool("atom", false, "ex([x], x^n + ai x^i + ... + a1 x + a0 <= 0)")
	)
	flag.Parse()
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "invalid argument\n")
		os.Exit(1)
	}

	var s sgn_table
	s.s = make([]sgn_t, *deg+1)
	s.h = make([]sgn_t, *deg+1)
	s.c = make([]sgn_t, *deg+1)
	s.atom = *atom
	s.count = make([]int, 2)
	s.deg = *deg
	s.debug = *debug
	if *all {
		s.val_equal = 1
		s.val_not = 0
	} else {
		s.val_equal = 0
		s.val_not = 1
	}
	s.pfalse = *false_print
	s.fp = bufio.NewWriter(os.Stdout)
	(&s).gen()
	if s.debug {
		fmt.Fprintf(s.fp, "T=%d, F=%d\n", s.count[1], s.count[0])
	}
	s.fp.Flush()

}
