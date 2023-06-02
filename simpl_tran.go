package ganrac

/////////////////////////////////////
// Formula Simplification for Real Quantifier Elimination
// Using Geometric Invariance
// H. Iwane, H. Anai, ISSAC 2017
// https://doi.org/10.1145/3087604.3087627
//
// translation formula ver. (scale invariant)
// related: simpl_homo
// related: simpl_rot
/////////////////////////////////////

import (
	"fmt"
)

func (p *Poly) get_tran_cond_coef(conds *List, n Level, xn []RObj) *List {
	if p.lv >= 0 {
		for _, cc := range p.c {
			if c, ok := cc.(*Poly); ok {
				c.get_tran_cond_coef(conds, n, xn)
			} else if !cc.IsZero() {
				panic("why!")
			}
		}
		return conds
	}
	if len(p.c) == 2 && p.c[0].IsZero() {
		if _, ok := p.c[1].(NObj); ok {
			xn[p.lv+2*n] = zero
			return conds
		}
	}
	conds.Append(p)
	return conds
}

// xn は破壊される
func (p *Poly) get_tran_cond(conds *List, n Level, xn []RObj) *List {
	var q RObj = p
	for i := Level(0); i < n; i++ {
		q = q.Subst(Sub(xn[i], Mul(xn[i+n], xn[2*n])), i)
	}
	qqq := Sub(p, q)
	if qq, ok := qqq.(*Poly); ok {
		conds = qq.get_tran_cond_coef(conds, n, xn)
	} else if !qqq.IsZero() {
		fmt.Printf("p=%v\n", p)
		fmt.Printf("q=%v\n", q)
		gpanic("(p - q) is not zero")
	}
	return conds
}

func (qeopt QEopt) get_tran_cond(fof Fof, cond qeCond, n Level, xn []RObj) *List {
	conds := NewList()
	fmls := make([]Fof, 1)
	fmls[0] = fof
	m := 0
	for m >= 0 {
		if m+1 != len(fmls) {
			gpanic(fmt.Sprintf("nooo invalid m, (%d, %d)", m, len(fmls)))
		}
		fof = fmls[m]

		switch p := fof.(type) {
		case FofQ:
			fmls[m] = p.Fml()
		case FofAO:
			fs := p.Fmls()
			fmls[m] = fs[0]
			m += len(fs) - 1
			for i := len(fs) - 1; i > 0; i-- {
				fmls = append(fmls, fs[i])
			}
		case *Atom:
			fmls = fmls[:m]
			m--
			for _, p := range p.p {
				conds = p.get_tran_cond(conds, n, xn)
			}
		default:
			fmls = fmls[:m]
			m--
		}
	}
	return conds
}

/*
 * GB に渡すとき負のレベルの変数では困るので，
 * Level に n 加える破壊的操作
 *
 * related: varShift()
 */
func (p *Poly) shift_var_tran(n Level) *Poly {
	if p.lv < 0 {
		p.lv += n
	}
	for _, cc := range p.c {
		if c, ok := cc.(*Poly); ok {
			c.shift_var_tran(n)
		}
	}

	return p
}

type poly_root struct {
	p *Poly
	r NObj
}

func (qeopt QEopt) gbred_tran(gb2, vars *List, lv Level) (*List, map[Level]poly_root) {
	ki := make(map[Level]poly_root, vars.Len())
	gb2 = gb2.Subst(one, lv)
	loop := true
	x0 := vars.getiPoly(int(lv))
	for i := 0; loop; i++ {
		loop = false
		gb := NewListN(gb2.Len())
		for _, gx := range gb2.Iter() {
			if g, ok := gx.(*Poly); ok {
				if g.isUnivariate() {
					// 値が確定した. 有理数でない場合ある???
					if len(g.c) != 2 {
						fmt.Printf("lv=%d\n", lv)
						fmt.Printf("gx %v => %v\n", gx, g)
						gpanic("unexpected")
					}
					if g.c[0].IsZero() {
						continue
					}
					// xi = xi - k * x0
					xi := vars.getiPoly(int(g.lv))
					p := (xi.Mul(g.c[1]).Add(x0.Mul(g.c[0]))).(*Poly)
					if g.c[1].Sign() < 0 {
						p = p.Neg().(*Poly)
					}
					ki[g.lv] = poly_root{
						p: p,
						r: (g.c[0].Neg().Div(g.c[1].(NObj))).(NObj),
					}
					// fmt.Printf("g[%d] = %v => [%v, %v]\n", g.lv, g, ki[g.lv].p, ki[g.lv].r)
					loop = true
				} else {
					gb.Append(g)
				}
			}
		}
		is_q := false
		for k, v := range ki {
			if _, ok := v.r.(*Int); !ok {
				is_q = true
			}
			gb = gb.Subst(v.r, k)
		}
		if is_q {
			for i, v := range gb.Iter() {
				p, _ := v.(*Poly).pp()
				gb.Seti(i, p)
			}
		}
		gb2 = gb
	}
	return gb2, ki
}

/*
 * 自由変数では平行移動可能な変数の組がなかった
 * fofq は prenex と思って，最初の連続する限量変数に対して消去可能か確かめる
 */
func (qeopt QEopt) qe_tran_quan(fofq FofQ, gb2 *List, vars *List, bb []bool, cond qeCond) Fof {
	f := fofq
	var fof2 Fof = fofq
	for {
		for _, q := range f.Qs() {
			if bb[q] {
				qeopt.log(cond, 3, "qetran", "Q<%s> %v\n", VarStr(q), fofq)
				gb2, _ = qeopt.gbred_tran(gb2, vars, q)
				fof2 = fof2.Subst(zero, q)        // 変数 q を 0 に固定
				gb := qeopt.g.ox.GB(gb2, vars, 0) // もっと平行移動できるかもしれない
				gb2, bb = qeopt.make_bb_tran(gb, Level(len(bb)))
				if gb2.Len() == 0 {
					fof2 = qeopt.qe(fof2, cond)
					return fof2
				}
			}
		}
		if g, ok := f.Fml().(FofQ); ok {
			f = g
		} else if fofq != fof2 {
			return qeopt.qe(fof2, cond)
		} else {
			return fof2
		}
	}
}

/*
 * 自由変数 lv が平行移動可能な変数の組である
 */
func (qeopt QEopt) qe_tran_free(fof Fof, gb2 *List, vars *List, lv Level, bb []bool, cond qeCond) Fof {
	qeopt.log(cond, 3, "qetran", "F<%s> %v\n", VarStr(lv), fof)

	gb2, ki := qeopt.gbred_tran(gb2, vars, lv)

	gb := qeopt.g.ox.GB(gb2, vars, 0)
	gb2, bb = qeopt.make_bb_tran(gb, Level(len(bb)))
	fof2 := fof.Subst(zero, lv) // 変数を 0 に固定
	if gb2.Len() > 0 {
		fof2 = qeopt.qe_tran_body(fof2, gb2, bb, vars, cond)
	} else {
		fof2 = qeopt.qe(fof2, cond)
	}

	// 自由変数をもとに戻す
	for k, v := range ki {
		fof2 = fof2.Subst(v.p, k)
	}

	return fof2
}

/*
 * GB の結果 (*List) を []*Poly に変換する
 * その際，0 確定したものは不要なので除去する
 *
 * bb は，gb2 に含まれる変数が真となる
 */
func (qeopt QEopt) make_bb_tran(gb *List, n Level) (gb2 *List, bb []bool) {
	bb = make([]bool, n)
	gb2 = NewListN(gb.Len())
	for _, pp := range gb.Iter() {
		if p, ok := pp.(*Poly); ok {
			if len(p.c) == 2 && p.c[0].IsZero() {
				if _, ok := p.c[1].(NObj); ok {
					// 定数確定
					continue
				}
			}
			gb2.Append(p)
			for lv, b := range bb {
				if !b && p.hasVar(Level(lv)) {
					bb[lv] = true
				}
			}
		}
	}
	return gb2, bb
}

func (qeopt QEopt) qe_tran(fof Fof, cond qeCond) Fof {
	n := qeopt.varn
	xn := make([]RObj, 2*n+1)

	for i := Level(0); i < n; i++ {
		xn[i] = NewPolyVar(i)
		xn[i+n] = NewPolyVar(i - n) // これが係数になるので一番小さな Level
	}
	xn[2*n] = NewPolyVar(n)

	vars := NewListN(int(n))
	for i := Level(0); i < n; i++ {
		vars.Append(xn[i])
	}

	// xn は破壊される
	conds := qeopt.get_tran_cond(fof, cond, n, xn)
	conds2 := NewListN(conds.Len())
	for _, pp := range conds.Iter() {
		pp.(*Poly).shift_var_tran(n)
		for i := Level(0); i < n; i++ {
			if xn[i+n] == zero {
				pp = pp.(*Poly).Subst(zero, Level(i))
				if _, ok := pp.(*Poly); !ok {
					break
				}
			}
		}
		if !pp.(RObj).IsZero() {
			conds2.Append(pp)
		}
	}
	if conds2.Len() == 0 {
		return fof
	}

	// 自明な解として 全部 0 があるから
	// gb = [1] になることはない
	gb := qeopt.g.ox.GB(conds2, vars, 0)

	// *List を []*Poly に変換
	gb2, bb := qeopt.make_bb_tran(gb, n)
	if gb2.Len() == 0 {
		// 自明な解しかなかった
		return nil
	}

	return qeopt.qe_tran_body(fof, gb2, bb, vars, cond)
}

/*
NOTE: F(x, y) and ex([y], G(x,y))

	      みたいな場合，
		  G(x,y) の y は F(x,y) の y とは別の変数として扱う必要があるがやっていない
*/
func (qeopt QEopt) qe_tran_body(fof Fof, gb2 *List, bb []bool, vars *List, cond qeCond) Fof {
	// 1変数消去できることが分かったから多少時間冗長な処理でもいいんじゃない?
	// prenex と仮定していいのか?
	for lv, b := range bb {
		if b && fof.hasFreeVar(Level(lv)) {
			return qeopt.qe_tran_free(fof, gb2, vars, Level(lv), bb, cond)
		}
	}

	if fofq, ok := fof.(FofQ); ok {
		return qeopt.qe_tran_quan(fofq, gb2, vars, bb, cond)
	}

	return fof
}
