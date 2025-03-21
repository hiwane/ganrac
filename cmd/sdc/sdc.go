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
// deg ... f(x) の次数
// t   ... 符号を考慮する要素の数
//         @  SDC の場合, h[deg-2], ..., h[0], c[deg-1], ..., c[1]
//         @ ROOT の場合, h[deg-2], ..., h[0]
//                 h[i]: PSC, c[i]: constant term
// DC  ... Don't Care の割合
//
// SDC: ex([x], x >= 0 && f(x) <= 0)
// | deg |  t |  True | False | rims2017 | casc2013 |    DC |
// |   2 |  2 |     2 |     5 |        2 |        2 | 0.222 |
// |   3 |  4 |     8 |    17 |        2 |        4 | 0.691 |
// |   4 |  6 |    78 |   117 |        6 |       10 | 0.733 |
// |   5 |  8 |   306 |   425 |        7 |       18 | 0.889 |
// |   6 | 10 |  2660 |  3089 |       31 |       57 | 0.903 |
// |   7 | 12 | 10838 | 11897 |       42 |     *121 | 0.957 |
// |   8 | 14 | 89592 | 87361 |     *134 |     *353 | 0.963 |
//
// ROOT (Atom):  ex([x], f(x) <= 0)
// | deg |  t |  True | False |     # |    DC |
// |   2 |  1 |     2 |     1 |     1 | 0.000 |
// |   4 |  3 |    16 |     9 |     4 | 0.074 |
// |   6 |  5 |   138 |    73 |    16 | 0.132 |
// |   8 |  7 |  1216 |   593 |    68 | 0.173 |
// |  10 |  9 | 10802 |  4881 |  *302 | 0.203 |
// |  12 | 11 | 96336 | 40665 | *1338 | 0.227 |
//
// CGS-QE (等式制約のみの場合)
// | deg |  t |  True | False |     # |    DC |
// |   2 |  2 |     2 |     3 |     1 | 0.444 |
// |   4 |  4 |    16 |    16 |     2 | 0.803 |
// |   6 |  6 |   104 |    79 |     6 | 0.749 |
// |   8 |  8 |   640 |   400 |    20 | 0.842 |
// |  10 | 10 |  3858 |  2083 |   *41 | 0.899 |
// |  12 | 12 | 23024 | 11072 |  *143 | 0.936 |
//
// SH-Theorem: 	QE&CAD L., Gonzalez-Vega
// p.303 Theorem2 Sturm-Habicht Structure Theorem
//
// espresso -Dexact -epos 4.in : 5
// espresso -estrong 8.in : 5
//
// see cmd/esp/

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
	deg      int
	s        []sgn_t // sign(f(+inf))
	h        []sgn_t // sign(coeff(f, j)): principal coefficient
	c        []sgn_t // sign(coeff(f, 0)): constant term
	count    []int
	fp       *bufio.Writer
	debug    bool
	valEqual val_t
	valNot   val_t
	pfalse   bool
	atom     bool
}

func (st *sgn_table) delta(n int) sgn_t {
	return sgn_t(-1 + ((n - 1) & 2))
}

func (st *sgn_table) VV() int {
	s := make([]sgn_t, 0, len(st.s))

	var last = -1
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

/** atom 用
 *
 * psc の符号から，実根の数を求める
 * QE&CAD L., Gonzalez-Vega p369 Proposition 1
 */
func (st *sgn_table) evalAtom() (val_t, string) {
	ks := 0
	var ask sgn_t // 非ゼロの右端
	var m int
	var p int
	var c int
	var e int
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
				var u = 1
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

	mes := fmt.Sprintf(": m(%2d) = p(%d) - c(%d) + e(%2d)", m, p, c, e)
	if m > 0 {
		return st.valNot, mes
	} else if m < 0 {
		return 3, mes
	} else {
		return st.valEqual, mes
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

	st.setS()

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
		return st.valNot, ""
	} else if mm == 0 {
		return st.valEqual, ""
	} else {
		return ret, "W-V < 0"
	}
}

func (st *sgn_table) setS() {
	// h/c の符号から s の符号が確定する
	var j int
	j = st.deg
	var u int
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

/** looph -> loops -> loopc **/
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

/** looph -> loops -> loopc **/
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

/**
 * psc の符号を決定する
 *
 * gen() から looph(0) で呼び出される
 *
 * looph -> loops -> loopc
 **/
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
			if t > 0 || st.pfalse {
				st.print(t, mes)
			}
			return
		}
		st.loops()
	} else {
		for _, i := range sgns { // 符号
			st.h[n] = i
			st.looph(n + 1)
		}
	}
}

func (st *sgn_table) setDC() {
	for i := 0; i < st.deg-1; i++ {
		st.s[i] = DC
		st.h[i] = DC
		st.c[i+1] = DC
	}
	st.c[0] = DC
}

func (st *sgn_table) printDC() {
	st.setDC()
	for i := 0; i < st.deg-1; i++ {
		st.h[i] = 3
		st.print(DC, "DC 11h")
		st.h[i] = DC
	}
	if st.atom {
		return
	}

	st.setDC()
	for i := 0; i < st.deg-1; i++ {
		st.c[i+1] = 3
		st.print(DC, "DC 11c")
		st.c[i+1] = DC
	}

	// @DC1
	st.setDC()
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
	st.setDC()
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
	st.printDC()
	st.footer()
}

func main() {
	var (
		deg        = flag.Int("d", 2, "degree")
		debug      = flag.Bool("debug", false, "debug print mode")
		falsePrint = flag.Bool("false", false, "print false")
		all        = flag.Bool("all", false, "all([x], x >= 0 impl f(x) > 0)")
		atom       = flag.Bool("atom", false, "ex([x], x^n + ai x^i + ... + a1 x + a0 <= 0)")
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
	if *all { // forall 用の論理式の場合
		s.valEqual = 1
		s.valNot = 0
	} else {
		s.valEqual = 0
		s.valNot = 1
	}
	s.pfalse = *falsePrint
	s.fp = bufio.NewWriter(os.Stdout)

	(&s).gen()
	if s.debug {
		fmt.Fprintf(s.fp, "T=%d, F=%d\n", s.count[1], s.count[0])
	}
	s.fp.Flush()

}
