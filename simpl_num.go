package ganrac

// a symbolic-numeric method for formula simplification
// 数値数式手法による論理式の簡単化, 岩根秀直, JSSAC 2017

// @TODO y^2-x<=0 && z^2+2*x*z+(-x+1)*y^2-2*x*y-x<=0 && x-1>=0 で -1<y<1 は真

import (
	"fmt"
	"sort"
)

type gray_intv []*Interval

type grayRegion struct {
	prec uint
	r    map[Level]gray_intv
	x    map[Level]*Interval
	t    *NumRegion // true/whiteRegion
	f    *NumRegion // false/blackRegion
}

// open/closed interval
type ninterval struct {
	// inf=nil => -infinity
	inf NObj
	// sup=nil => +infinity
	sup NObj
}

// nil は empty set を表す.
type NumRegion struct {
	// len(r[lv]) == 0 は  -inf <= x[lv] <= +inf を表す.
	// もし，要素があれば, その要素の union であるが，
	// 昇順にソートされ，重複はない
	r map[Level][]*ninterval
}

////////////////////////////////////////////////////////////////////////
// ninterval
////////////////////////////////////////////////////////////////////////

func (x *ninterval) Format(s fmt.State, format rune) {
	fmt.Fprintf(s, "[")
	if x.inf == nil {
		fmt.Fprintf(s, "-inf")
	} else {
		x.inf.Format(s, format)
	}
	fmt.Fprintf(s, ",")
	if x.sup == nil {
		fmt.Fprintf(s, "+inf")
	} else {
		x.sup.Format(s, format)
	}
	fmt.Fprintf(s, "]")
}

////////////////////////////////////////////////////////////////////////
// NumRegion
////////////////////////////////////////////////////////////////////////

func newGrayResion(lv Level, prec uint) *grayRegion {
	p := new(grayRegion)
	p.r = make(map[Level]gray_intv, lv)
	p.x = make(map[Level]*Interval, lv)
	p.prec = prec
	return p
}

func (gray *grayRegion) setTF(t, f *NumRegion) {
	gray.t = t
	gray.f = f
}

func (gray *grayRegion) upGray(lv Level) {
	gg := gray.t.getU(gray.f, lv)
	gray.r[lv] = gg
	x := newInterval(gray.prec)
	x.inf = gg[0].inf
	x.sup = gg[len(gg)-1].sup
	gray.x[lv] = x
}

func NewNumRegion() *NumRegion {
	f := new(NumRegion)
	f.r = make(map[Level][]*ninterval, 0)
	return f
}

func (m *NumRegion) add(inf, sup NObj, lv Level) {
	m.r[lv] = append(m.r[lv], &ninterval{inf, sup})
}

func (m *NumRegion) Format(s fmt.State, format rune) {
	if m == nil {
		fmt.Fprintf(s, "{empty}")
		return
	}
	fmt.Fprintf(s, "{")
	sep := ""
	for lv, vv := range m.r {
		fmt.Fprintf(s, "%s%s:[", sep, VarStr(lv))
		for i, v := range vv {
			if i != 0 {
				fmt.Fprintf(s, ",")
			}
			v.Format(s, format)
		}
		fmt.Fprintf(s, "]")
		sep = ", "
	}
	fmt.Fprintf(s, "}")
}

// 左端点が小さい方，または
// 左端点が同じかつ，右端点が小さい方
func (_ *NumRegion) less(n, m *ninterval) bool {
	if n.inf == nil {
		if m.inf != nil {
			return true
		}
	} else if m.inf == nil {
		return false
	} else {
		c := n.inf.Cmp(m.inf)
		if c < 0 {
			return true
		} else if c > 0 {
			return false
		}
	}
	// n.inf == m.inf

	if n.sup == nil {
		return false
	}
	if m.sup == nil {
		return true
	}
	return n.sup.Cmp(m.sup) < 0
}

func (_ *NumRegion) maxInt(n, m int) int {
	if n < m {
		return m
	} else {
		return n
	}
}

// lv 限定で intersect 計算
func (m *NumRegion) intersect_lv(n *NumRegion, lv Level) []*ninterval {
	mx, ok := m.r[lv]
	if !ok {
		return n.r[lv]
	}
	nx, ok := n.r[lv]
	if !ok {
		return mx
	}

	ret := make([]*ninterval, 0, m.maxInt(len(mx), len(nx)))
	i := 0
	j := 0
	for i < len(mx) && j < len(nx) {
		if m.less(mx[i], nx[j]) {
			// [m .... m]
			//                 [n ..... n]
			if mx[i].sup != nil && nx[j].inf != nil && mx[i].sup.Cmp(nx[j].inf) <= 0 {
				// 重複なし.
				i++
				continue
			}

			x := new(ninterval)
			x.inf = nx[j].inf
			if mx[i].sup == nil || (nx[j].sup != nil) && mx[i].sup.Cmp(nx[j].sup) > 0 {
				// [m ....... m]
				//   [n...n]
				x.sup = nx[j].sup
				j++
			} else {
				// [m .... m]
				//     [n.......n]
				x.sup = mx[i].sup
				i++
			}
			ret = append(ret, x)
		} else {
			if nx[j].sup != nil && mx[i].inf != nil && nx[j].sup.Cmp(mx[i].inf) <= 0 {
				//              [m....m]
				//  [n .... n]
				// 重複なし.
				j++
				continue
			}

			x := new(ninterval)
			x.inf = mx[i].inf
			if nx[j].sup == nil || mx[i].sup != nil && nx[j].sup.Cmp(mx[i].sup) > 0 {
				//   [m...m]
				// [n ....... n]
				x.sup = mx[i].sup
				i++
			} else {
				//     [m.......m]
				// [n .... n]
				x.sup = nx[j].sup
				j++
			}
			ret = append(ret, x)
		}
	}
	return ret
}

// 2つの領域の intersect を計算する
func (m *NumRegion) intersect(n *NumRegion) *NumRegion {
	if n == nil || m == nil {
		return nil
	}
	u := NewNumRegion()
	for lv, _ := range m.r {
		v := m.intersect_lv(n, lv)
		if v != nil {
			u.r[lv] = v
		}
	}
	for lv, nx := range n.r {
		if _, ok := u.r[lv]; !ok {
			u.r[lv] = nx
		}
	}
	return u
}

// lv 限定で union 計算
func (m *NumRegion) union_lv(n *NumRegion, lv Level) []*ninterval {
	// assume: m != nil && n != nil
	mx, okm := m.r[lv]
	nx, okn := n.r[lv]
	if !okn && !okm {
		return nil
	} else if !okn || len(nx) == 0 {
		return mx
	} else if !okm || len(mx) == 0 {
		return nx
	}

	nm := make([]*ninterval, len(mx)+len(nx))
	copy(nm, mx)
	copy(nm[len(mx):], nx)
	sort.Slice(nm, func(i, j int) bool {
		return m.less(nm[i], nm[j])
	})

	r := nm[0]
	ret := make([]*ninterval, 0, len(nm))
	for _, a := range nm[1:] {
		// 重なりがあるか.
		if a.inf != nil && r.sup != nil && r.sup.Cmp(a.inf) <= 0 {
			// [r ..... r]
			//              [a  ..... a]
			ret = append(ret, r)
			r = a
		} else {
			// 合体.
			// [r ..... r]
			//       [a  ..... a]
			v := new(ninterval)
			v.inf = r.inf
			v.sup = a.sup
			r = v
		}
	}
	return append(ret, r)
}

func (m *NumRegion) union(n *NumRegion) *NumRegion {
	if n == nil {
		return m
	}
	if m == nil {
		return n
	}
	u := NewNumRegion()
	for lv := range m.r {
		u.r[lv] = m.union_lv(n, lv)
	}
	for lv, nx := range n.r {
		if _, ok := u.r[lv]; !ok {
			u.r[lv] = nx
		}
	}

	return u
}

func (m *NumRegion) del(qq []Level) *NumRegion {
	if m == nil {
		return m
	}
	for _, q := range qq {
		if _, ok := m.r[q]; ok {
			delete(m.r, q)
		}
	}
	return m
}

// trueRegion と falseRegion 以外の領域を返す.
func (m *NumRegion) getU(n *NumRegion, lv Level) []*Interval {
	prec := uint(53)

	var mn []*ninterval
	if m == nil && n == nil {
		// goto _R させてくれない
		mn = nil
	} else if n == nil {
		mn = m.r[lv]
	} else if m == nil {
		mn = n.r[lv]
	} else {
		// 補集合を広くとるので， union は狭く.
		// このときの union は境界は残さないといけない. (m.sup=n.inf)
		// [m ... m]
		//         [n .... n] なら
		// [m ....m] [n ...n] と別区間扱い
		mn = m.union_lv(n, lv)
	}
	if mn == nil {
		mn = []*ninterval{}
	}

	// mn の補集合
	var inf NObj
	inf = nil
	xs := make([]*Interval, 0, len(mn)+2)
	for _, y := range mn {
		if y.inf != nil {
			x := newInterval(prec)
			if inf == nil {
				x.inf.SetInf(true)
			} else {
				xi := inf.toIntv(prec).(*Interval)
				x.inf.Set(xi.inf)
			}
			xi := y.inf.toIntv(prec).(*Interval)
			x.sup = xi.sup
			xs = append(xs, x)
		}
		inf = y.sup
	}
	if inf != nil {
		x := newInterval(prec)
		xi := inf.toIntv(prec).(*Interval)
		x.inf.Set(xi.inf)
		x.sup.SetInf(false)
		xs = append(xs, x)
	} else if len(xs) == 0 {
		x := newInterval(prec)
		x.inf.SetInf(true)
		x.sup.SetInf(false)
		xs = append(xs, x)
	}
	return xs
}

////////////////////////////////////////////////////////////////////////
// simplNumPoly
////////////////////////////////////////////////////////////////////////

// assume: poly is univariate
// returns (OP, pos, neg)
// OP = (t,f) 以外で取りうる符号
func (poly *Poly) simplNumUniPoly(gray *grayRegion) (OP, *NumRegion, *NumRegion) {
	// 重根を持っていたら...?
	roots := poly.realRootIsolation(-30)
	// fmt.Printf("   simplNumUniPoly(%v) t=%v, f=%v, #root=%v\n", poly, t, f, len(roots))
	if len(roots) == 0 { // 符号一定
		if poly.Sign() > 0 {
			return GT, NewNumRegion(), nil
		} else {
			return LT, nil, NewNumRegion()
		}
	}
	xs := gray.r[poly.lv]
	prec := uint(53)

	// 根が， unknown 領域に含まれていない，かつ
	// 同じ領域に入るなら符号一定が確定する
	if len(xs) > 0 { // T union F に対してやったほうが楽か.
		idx := -1
		rooti := roots[0].toIntv(prec)
		if xs[0].inf != nil && rooti.sup.Cmp(xs[0].inf) < 0 {
			idx = 0
		} else {
			for i := 1; i < len(xs); i++ {
				if xs[i-1].sup.Cmp(rooti.inf) < 0 && rooti.sup.Cmp(xs[i].inf) < 0 {
					idx = i
					break
				}
			}
			if idx < 0 && xs[len(xs)-1].sup.Cmp(rooti.inf) < 0 {
				idx = len(xs)
			}
		}

		if idx >= 0 {
			for _, root := range roots {
				rooti := root.toIntv(prec)
				if idx == 0 {
					//  root  ... xs[0]
					if rooti.sup.Cmp(xs[idx].inf) >= 0 {
						idx = -1
						break
					}
				} else if idx == len(xs) {
					// xs[-1] ... root
					if xs[idx-1].sup.Cmp(rooti.inf) >= 0 {
						idx = -1
						break
					}
				} else {
					// xs[idx-1] ... root ... xs[idx]
					if !(xs[idx-1].sup.Cmp(rooti.inf) < 0 && rooti.sup.Cmp(xs[idx].inf) < 0) {
						idx = -1
						break
					}
				}
			}

			if idx >= 0 {
				// 根が連結した known 領域のみに含まれていて,
				// unknown 領域での符号が一定であることが確定した
				// 重根がないと仮定しているので，idx で符号が確定できる.
				if 0 < idx && idx < len(xs) {
					if poly.deg()%2 == 0 {
						// 偶数次
						if poly.Sign() > 0 {
							return GT, NewNumRegion(), nil
						} else {
							return LT, nil, NewNumRegion()
						}
					} else {
						// 奇数次
						pinf := NewNumRegion()
						pinf.r[poly.lv] = append(pinf.r[poly.lv], &ninterval{roots[len(roots)-1].low.upperBound(), nil})
						ninf := NewNumRegion()
						ninf.r[poly.lv] = append(ninf.r[poly.lv], &ninterval{nil, roots[0].low})
						if poly.Sign() > 0 {
							return OP_TRUE, pinf, ninf
						} else {
							return OP_TRUE, ninf, pinf
						}
					}
				}

				sgn := poly.Sign()
				if idx != 0 {
					// x = -inf
					sgn *= 2*(len(poly.c)%2) - 1
				}
				if sgn > 0 {
					return GT, NewNumRegion(), nil
				} else if sgn < 0 {
					return LT, nil, NewNumRegion()
				}
				panic("?")
			}
		}
	}

	// gray region を代入
	x := gray.x[poly.lv]
	p := poly.toIntv(prec).(*Poly)
	pp := p.SubstIntv(x, p.lv, prec).(*Interval)
	if ss := pp.Sign(); ss > 0 {
		return GT, NewNumRegion(), nil
	} else if ss < 0 {
		return LT, nil, NewNumRegion()
	}

	nr := []*NumRegion{NewNumRegion(), NewNumRegion()}

	var inf NObj
	for i, intv := range roots {
		nr[i%2].add(inf, intv.low, poly.lv)
		if intv.point {
			inf = intv.low
		} else {
			inf = intv.low.upperBound()
		}
	}
	nr[len(roots)%len(nr)].add(inf, nil, poly.lv)

	sgn := poly.Sign()
	sgn *= (len(poly.c)%2)*2 - 1 // x=-inf での poly の符号
	if sgn > 0 {
		return OP_TRUE, nr[0], nr[1]
	} else {
		return OP_TRUE, nr[1], nr[0]
	}
}

func (poly *Poly) simplNumNvar(g *Ganrac, gray *grayRegion, dv Level) (OP, *NumRegion, *NumRegion) {
	prec := gray.prec
	p := poly.toIntv(prec).(*Poly)

	for lv := poly.lv; lv >= 0; lv-- {
		if lv == dv {
			continue
		}
		x := gray.x[lv]
		switch pp := p.SubstIntv(x, lv, prec).(type) {
		case *Poly:
			p = pp
		case *Interval:
			// 区間 u での符号がきまった
			if ss := pp.Sign(); ss > 0 {
				goto _GT
			} else if ss < 0 {
				goto _LT
			} else if pp.inf.Sign() == 0 {
				goto _GE
			} else if pp.sup.Sign() == 0 {
				goto _LE
			}
			return OP_TRUE, gray.t, gray.f
		}
	}

	// fmt.Printf("   simplNumNvar() p[%d]=%f\n", dv, p)
	if len(p.c) == 2 { // linear
		if p.c[1].(*Interval).ContainsZero() {
			return OP_TRUE, gray.t, gray.f
		}
		// x := p.c[0].Div(p.c[1].Neg().(NObj)).(*Interval)
		// if ss := x.Sign(); ss > 0 {
		// } else if ss < 0 {
		// } else if x.inf.Sign() == 0 {
		// } else if x.sup.Sign() == 0 {
		// } else {
		// }
		return OP_TRUE, gray.t, gray.f
	}

	return OP_TRUE, gray.t, gray.f
_GT:
	return GT, gray.t, gray.f
_GE:
	return GE, gray.t, gray.f
_LE:
	return LE, gray.t, gray.f
_LT:
	return LT, gray.t, gray.f
}

/*
 * 因数分解された多項式のリストを受け取り，それらを単純化する．
 * start: 0 / 1 がくることを想定. ox.fctr からの復帰で，0番目の要素は定数で skip する必要があるため
 *
 * 復帰値の []RObj は fctr.get(*, 1) % 2 == 0 の場合があると正しくない
 *  simplNumPolyFctr() から呼び出す場合には復帰を利用しない
 *  それ以外の場合には，simplFctr() により重複は取り除かれていることを期待する
 */
func (fctr *List) simplNumPolys(g *Ganrac, start int, t, f *NumRegion) ([]RObj, int, OP, *NumRegion, *NumRegion) {

	// start > 0 の場合は，0番目の要素の符号を見る必要があるが，
	// 呼び出し元にまかせる
	ret_op := GT
	ret_sgn := 1 // 因子の一部を削ったときの，削った部分の符号

	var pret *NumRegion
	var nret *NumRegion

	ps := make([]RObj, 0, fctr.Len()-start) // 復帰用
	es := make([]int, 0, fctr.Len()-start)  // 次数
	var pos, neg *NumRegion
	n := fctr.Len() - 1

	var mul_even RObj = one
	var mul_odd RObj = one
	for i := n; i >= start; i-- {
		ei, _ := fctr.Geti(i, 1)
		e := ei.(*Int).Int64()

		pi, _ := fctr.Geti(i, 0)
		p := pi.(*Poly)

		var op OP
		op, pos, neg = p.simplNumPoly(g, t, f, p.lv)
		if op != OP_TRUE && false {
			fmt.Printf("   simplNumPolys(%v) :: %s\n", p, op)
		}
		if op == GT {
			continue // 因子を削れた
		} else if op == LT {
			if e%2 != 0 {
				ret_op = ret_op.neg()
				ret_sgn *= -1
			}
			continue // 因子を削れた
		} else if op == GE {
			ret_op |= EQ
		} else if op == LE {
			ret_op |= EQ
			if e%2 != 0 {
				ret_op = ret_op.neg()
			}
		} else if e%2 == 0 {
			if op&EQ != 0 {
				ret_op |= EQ
			}
		} else {
			ret_op = OP_TRUE
		}
		ps = append(ps, p)
		es = append(es, int(e))
		if n != start {
			if e%2 == 0 {
				mul_even = Mul(mul_even, p)
			} else {
				mul_odd = Mul(mul_odd, p)
			}
		}
	}

	if n <= start || len(ps) == 0 {
		// 因子が全部なくなった
		return ps, ret_sgn, ret_op, pos, neg
	}

	// 一変数の場合，実根の分離を行うので，平方因子があると困る
	if mul_even == one {
		if mul, ok := mul_odd.(*Poly); ok {
			_, pos, neg = mul.simplNumPoly(g, t, f, mul.lv)
			pret = pret.intersect(pos)
			nret = nret.intersect(nret)
			return ps, ret_sgn, ret_op, pos, neg
		}
	}
	// fmt.Printf("p^%d=%v, t=%v, f=%v => op=%v, pos=%v, neg=%v\n", p, e, t, f, op, pos, neg)
	return ps, ret_sgn, ret_op, pret, nret
}

/*
 * atom の外の多項式の符号評価
 *
 *   2 次の場合に，主係数，判別式の符号を判定する
 */
func (poly *Poly) simplNumPolyFctr(g *Ganrac, t, f *NumRegion) (OP, *NumRegion, *NumRegion) {
	fctr := g.ox.Factor(poly)
	cont, _ := fctr.Geti(0, 0)

	_, _, ret_op, pret, nret := fctr.simplNumPolys(g, 1, t, f)
	if cont.(RObj).Sign() > 0 {
		return ret_op, pret, nret
	} else {
		// ret_op は OP_TRUE/OP_FALSE の場合がある
		return ret_op.neg(), nret, pret
	}
}

func (poly *Poly) simplNumPoly(g *Ganrac, t, f *NumRegion, dv Level) (OP, *NumRegion, *NumRegion) {

	// とりま全部 gray 区間を代入してみる
	prec := uint(53)
	pi := poly.toIntv(prec)
	gray := newGrayResion(poly.lv+1, prec)
	gray.setTF(t, f)
	for lv := poly.lv; lv >= 0; lv-- {
		if !poly.hasVar(lv) {
			continue
		}
		gray.upGray(lv)
		x := gray.x[lv]
		if pp, ok := pi.(*Poly); ok {
			pi = pp.SubstIntv(x, lv, prec)
		}
	}
	if pp, ok := pi.(*Interval); ok {
		if !pp.ContainsZero() {
			// 符号が確定した
			s := pp.Sign()
			if s > 0 {
				return GT, NewNumRegion(), nil
			} else if s < 0 {
				return LT, nil, NewNumRegion()
			}
		}
	} else {
		panic("bug")
	}
	if poly.IsUnivariate() {
		return poly.simplNumUniPoly(gray)
	}
	var pret, nret *NumRegion
	op_ret := OP_TRUE
	for v := poly.lv; v >= 0; v-- {
		deg := poly.Deg(v)
		if deg == 0 {
			continue
		}
		s, pos, neg := poly.simplNumNvar(g, gray, v)
		op_ret &= s
		if op_ret == GT || op_ret == LT {
			return op_ret, pos, neg
		}
		pret = pret.union(pos)
		nret = nret.union(neg)
		if deg != 2 || v > dv {
			// \sum_i x_i^2 = 1 のようなケースで, 判別式爆発が起こる.
			// v > dv は，多少でも削減するために限定する
			continue
		}

		/////////////////////////////////////////////////////////////////////////
		// 2次であれば，判別式が負なら符号が主係数の符号と一致することを利用する.
		/////////////////////////////////////////////////////////////////////////

		// 主係数の符号を考える
		c2 := poly.Coef(v, 2)
		var pc, nc *NumRegion // sgn of leading Coefficient
		if cp, ok := c2.(*Poly); ok {
			s, pc, nc = cp.simplNumPolyFctr(g, t, f)
			if s != GT && s != LT {
				continue
			}
		} else if c2.Sign() > 0 {
			pc = NewNumRegion()
			s = GT
		} else {
			nc = NewNumRegion()
			s = LT
		}
		discrim := poly.discrim2(v)
		// 判別式の符号を考える
		var sgn_op OP = 0
		var nd *NumRegion = nil // sgn of Discrimination
		switch dd := discrim.(type) {
		case NObj:
			sgn := dd.Sign()
			if sgn > 0 {
				sgn_op = GT
			} else if sgn < 0 {
				sgn_op = LT
				nd = NewNumRegion()
			} else {
				sgn_op = EQ
			}
		case *Poly:
			sgn_op, _, nd = dd.simplNumPolyFctr(g, t, f)
		default:
			panic(fmt.Sprintf("unknown type. dd=%v", dd))
		}

		pret = pret.union(pc.intersect(nd))
		nret = nret.union(nc.intersect(nd))

		if sgn_op&GT != 0 { // GT, GE, NE
			continue
		} else if sgn_op == LT {
			op_ret &= s
		} else { // LE, EQ
			op_ret &= (s | EQ)
		}
		if err := op_ret.valid(); err != nil {
			panic(fmt.Sprintf("invalid op_ret=%v", op_ret))
		}

		if op_ret == GT || op_ret == LT {
			break
		}
	}
	return op_ret, pret, nret
}

////////////////////////////////////////////////////////////////////////
// simplNum
////////////////////////////////////////////////////////////////////////

func (p *AtomT) simplNum(g *Ganrac, t, f *NumRegion) (Fof, *NumRegion, *NumRegion) {
	return p, t, f
}

func (p *AtomF) simplNum(g *Ganrac, t, f *NumRegion) (Fof, *NumRegion, *NumRegion) {
	return p, t, f
}

func (atom *Atom) simplNum(g *Ganrac, t, f *NumRegion) (Fof, *NumRegion, *NumRegion) {
	// simplFctr 通過済みと仮定したいところだが.
	if err := atom.valid(); err != nil {
		fmt.Printf("simplNum input invalid atom=%V: %v\n", atom, atom)
		panic("stop")
	}

	fctr := NewListN(len(atom.p))
	for _, p := range atom.p {
		fctr.Append(NewList(p, one))
	}

	rs, rsgn, s, pp, nn := fctr.simplNumPolys(g, 0, t, f)
	if len(rs) < len(atom.p) { // 因子が減ったということ
		// fmt.Printf("          rs=%v\n", rs)
		// fmt.Printf("           t=%v\n", t)
		// fmt.Printf("           f=%v\n", f)
		var atom2 Fof
		if len(rs) == 0 {
			if atom.op&s == 0 {
				atom2 = falseObj
			} else if atom.op|s == atom.op {
				atom2 = trueObj
			} else {
				atom2 = NewAtoms(rs, atom.op&s)
			}
		} else {
			if rsgn > 0 {
				atom2 = NewAtoms(rs, atom.op)
			} else {
				atom2 = NewAtoms(rs, atom.op.neg())
			}
		}
		if err := atom2.valid(); err != nil {
			fmt.Printf("simplNum(%v) => %v\n", atom, rs)
			for iii, ppp := range rs {
				fmt.Printf("  %d: %v\n", iii, ppp)
			}
			fmt.Printf("simplNum newAtoms=%V: %v\n", atom, atom)
			panic("stop, bug")
		}
		// fmt.Printf("simplNum(%v) => %v, t=%v, f=%v, s=%d, pp=%v, nn=%v, rs=%v\n", atom, atom2, t, f, s, pp, nn, rs)
		switch atomx := atom2.(type) {
		case *AtomT, *AtomF:
			return atom2, t, f
		case *Atom:
			atom = atomx
		default:
			panic(fmt.Sprintf("unknown type, bug, atom=%v => %v, ps=%v", atom, atom2, rs))
		}
	}

	if s&atom.op == 0 {
		return falseObj, nil, NewNumRegion()
	} else if s|atom.op == atom.op {
		if err := atom.valid(); err != nil {
			panic("bug")
		}
		g.log(3, 1, "smplnum: {%V:%v} ==> true: t=%v, f=%v\n", atom, atom, t, f)
		return trueObj, NewNumRegion(), nil
	} else if atom.op == LT || atom.op == LE {
		return atom, nn, pp
	} else if atom.op == GT || atom.op == GE {
		return atom, pp, nn
	}
	pn := pp.union(nn)
	if atom.op == EQ {
		return atom, nil, pn
	} else {
		return atom, pn, nil
	}
}

func sliceIn[T int](ary []T, val T) bool {
	for _, v := range ary {
		if v == val {
			return true
		}
	}
	return false
}

func (p *FmlAnd) simplNum(g *Ganrac, t, f *NumRegion) (Fof, *NumRegion, *NumRegion) {
	ts := make([]*NumRegion, 0, len(p.fml))
	fs := make([]*NumRegion, 0, len(p.fml))
	fmls := make([]Fof, 0, len(p.fml))
	print_log := false
	if print_log {
		fmt.Printf("   And.simplNum input And=%v, t=%v, f=%v\n", p, t, f)
	}
	for i := range p.fml {
		if print_log {
			fmt.Printf("   And.simplNum case1>%d: f=%v, t=%v, f=%v\n", i, p.fml[i], t, f)
		}
		fml, tt, ff := p.fml[i].simplNum(g, t, f)
		if print_log {
			fmt.Printf("   And.simplNum case1<%d: f=%v, t=%v, f=%v\n", i, p.fml[i], t, f)
		}
		if print_log && !p.fml[i].Equals(fml) {
			fmt.Printf("@@ And.simplNum @AND@ %v\n", p)
			fmt.Printf("@@ And.simplNum[1st,%d/%d] %v -> %v, false=%v\n", i+1, len(p.fml), p.fml[i], fml, ff)
			fmt.Printf("@@ t=%v\n", t)
			fmt.Printf("@@ f=%v\n", f)
		}
		if _, ok := fml.(*AtomF); ok {
			return falseObj, nil, nil
		}
		if _, ok := fml.(*AtomT); ok {
			continue
		}
		fmls = append(fmls, fml)
		ts = append(ts, tt)
		fs = append(fs, ff)
	}
	if print_log {
		fmt.Printf("   And.simplNum case2 fmls=%v, t=%v\n", fmls, t)
	}
	if len(fmls) <= 1 {
		if len(fmls) == 0 {
			return trueObj, nil, nil
		}
		return fmls[0], ts[0], fs[0]
	}

	var tret *NumRegion
	fret := f
	except := make([]int, 0, len(fmls))
	for i, fml := range fmls {
		ff := f
		// i 番以外の情報から white/black region を構築する
		var tt *NumRegion
		for j := 0; j < len(fmls); j++ {
			if j != i && !sliceIn(except, j) {
				ff = ff.union(fs[j])
				//		tt = tt.intersect(ts[j])	// intersect とってもほぼ偽になるからやらない
			}
		}

		if print_log {
			fmt.Printf("   And.simplNum case3>%d: f=%v, t=%v, f=%v\n", i, fmls[i], t, ff)
		}
		fmls[i], tt, ff = fml.simplNum(g, t, ff)
		if print_log {
			fmt.Printf("   And.simplNum case3<%d: f=%v, t=%v, f=%v\n", i, fmls[i], t, ff)
		}
		if err := fmls[i].valid(); err != nil {
			fmt.Printf("simplNum(%v) => %v; %V\n", fml, fmls[i], fmls[i])
			panic("stop, bug")
		}
		if !fml.Equals(fmls[i]) {
			except = append(except, i)
			if print_log {
				fmt.Printf("@@ And.simplNum @AND@ %v, t=%v\n", p, t)
				fmt.Printf("   And.simplNum[2nd,%d/%d] %v => %v, ff=%v\n", i+1, len(p.fml), fml, p.fml[i], ff)
				// fmt.Printf("@@ And.simplNum[2ND,%d/%d] %V -> %V\n", i+1, len(p.fml), fml, p.fml[i])
				fmt.Printf("               [2nd,%d/%d] t=%v\n", i+1, len(p.fml), t)
				fmt.Printf("               [2nd,%d/%d] f=%v\n", i+1, len(p.fml), f)
			}
		}
		if _, ok := fmls[i].(*AtomF); ok {
			return falseObj, nil, NewNumRegion()
		}
		tret = tret.intersect(tt)
		fret = fret.union(ff)
		// fmt.Printf("@@ And.simplNum[2nd,%d/%d] %v, Fi=%v, Fret=%v\n", i+1, len(fmls), fmls[i], ff, fret)
	}
	tret = tret.union(t)
	fml := NewFmlAnds(fmls...)
	if print_log {
		fmt.Printf("## And.simplNum[%v,end,#=%d] %v T=%v, F=%v\n", p, len(fmls), fml, tret, fret)
		fmt.Printf("               [%v,end,#=%d] %v t=%v, f=%v\n", p, len(fmls), fml, t, f)
	}
	return fml, tret, fret
}

func (p *FmlOr) simplNum(g *Ganrac, t, f *NumRegion) (Fof, *NumRegion, *NumRegion) {
	// @TODO サボり
	q := p.Not()
	q, f, t = q.simplNum(g, f, t)
	return q.Not(), t, f
}

func (p *ForAll) simplNum(g *Ganrac, t, f *NumRegion) (Fof, *NumRegion, *NumRegion) {
	fml, t, f := p.fml.simplNum(g, t, f)
	return NewQuantifier(true, p.q, fml), t.del(p.q), f.del(p.q)
}

func (p *Exists) simplNum(g *Ganrac, t, f *NumRegion) (Fof, *NumRegion, *NumRegion) {
	fml, t, f := p.fml.simplNum(g, t, f)
	return NewQuantifier(false, p.q, fml), t.del(p.q), f.del(p.q)
}
