package ganrac

import (
	"fmt"
)

// Quantifier Elimination for Formulas Constrained by Quadratic Equations via Slope Resultants
// Hoon Hong, The computer J., 1993

// ///////////////////////////////////////////////
// quadEQ_t... 等式制約が複数ある場合に，
//             どの等式制約を選択するかを
//             決定するための情報
// ///////////////////////////////////////////////

type quadEQ_t struct {
	lv    Level // 消去する変数レベル
	a     *Atom // specialQE で利用するを等式制約
	p     *Poly // a.Poly()
	z     RObj  // p の主係数
	idx   int   // ?
	deg   int   // p の lv に関する次数/
	lc    bool  // lc(a.p) is constant?
	uni   bool  // p is univariate
	coefs []RObj
	atoms map[OP]*Atom
	diff  map[int]*vs_sample_point
}

func (minatom *quadEQ_t) String() string {
	return fmt.Sprintf("quadEQ_t{lv=%d, a=%v, p=%v, z=%v, idx=%d, deg=%d, lc=%v, uni=%v}", minatom.lv, minatom.a, minatom.p, minatom.z, minatom.idx, minatom.deg, minatom.lc, minatom.uni)
}

/**
 * 現在設定されているものよりも簡易だったら，更新する
 * @param d int = poly.Deg(q)
 * @param q Level Quantifier
 */
func (minatom *quadEQ_t) SetAtomIfEasy(fof FofQ, atom *Atom, dmin, dmax, ii int) {
	poly := atom.getPoly()
	for _, q := range fof.Qs() {
		d := poly.Deg(q)
		if d < dmin || d > dmax {
			continue
		}

		z := poly.Coef(q, uint(d))
		_, lc := z.(NObj)
		univ := poly.IsUnivariate()

		// fmt.Printf("EQ!  poly=%v, d=%d, lc=%v\n", poly, d, lc)

		// 次数が低いか，主係数が定数なものを選択する
		if minatom.a == nil || minatom.deg > d ||
			(minatom.deg == d && univ) ||
			(minatom.deg == d && !minatom.uni && lc) ||
			(minatom.deg == d && !minatom.lc) {
			minatom.lv = q
			minatom.a = atom
			minatom.p = poly
			minatom.z = z
			minatom.idx = ii
			minatom.deg = d
			minatom.lc = lc
			minatom.uni = univ
		}
	}
}

type fof_quad_eq struct {
	p       *Poly
	sgn_lcp int // sign of lc(p)
	sgn_s   int // for quadratic case: 2つの根のうち大きいほうなら +1, 小さいなら -1
	lv      Level
	g       *Ganrac
}

func NewFofQuadEq(g *Ganrac, p *Poly, lv Level) *fof_quad_eq {
	tbl := new(fof_quad_eq)
	tbl.g = g
	tbl.p = p
	tbl.lv = lv
	return tbl
}

// ///////////////////////////////////////////////
// 主係数の符号によって不等号の向きが変わらないような論理式か
// ///////////////////////////////////////////////
func quadeq_isEven(f Fof, lv Level) bool {
	stack := make([]Fof, 1)
	stack[0] = f
	for len(stack) > 0 {
		f = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch ff := f.(type) {
		case *Atom:
			if ff.op == EQ || ff.op == NE {
				continue
			}
			if ff.Deg(lv)%2 != 0 {
				return false
			}
		case FofAO:
			stack = append(stack, ff.Fmls()...)
		case FofQ:
			stack = append(stack, ff.Fml())
		}
	}
	return true
}

// ///////////////////////////////////////////////
// 核
// ///////////////////////////////////////////////
func qe_lineq(a *Atom, param any) (Fof, bool) {
	t := param.(*fof_quad_eq)
	if !a.hasVar(t.lv) {
		return a, false
	}
	res := make([]RObj, len(a.p))
	for i, p := range a.p {
		res[i] = t.g.ox.Resultant(t.p, p, t.lv)
	}
	op := a.op
	if t.sgn_lcp < 0 && a.Deg(t.lv)%2 != 0 {
		op = op.neg()
	}
	return NewAtoms(res, op), true
}

func quadeq_getrst(gan *Ganrac, f, g *Poly, lv Level) (RObj, RObj, RObj) {
	r := gan.ox.Resultant(f, g, lv)
	t := gan.ox.Slope(f, g, lv, 0)
	s := gan.ox.Slope(f, g, lv, 1)
	return r, t, s
}

func qe_quadeq(a *Atom, param interface{}) (Fof, bool) {
	u := param.(*fof_quad_eq)
	if !a.hasVar(u.lv) {
		return a, false
	}
	f := u.p
	g := a.getPoly()
	aop := a.op
	switch aop {
	case EQ, NE:
		opt := LE
		ops := LE
		if u.sgn_s < 0 {
			ops = GE
		}
		r, t, s := quadeq_getrst(u.g, f, g, u.lv)
		ret := NewFmlAnd(NewAtom(r, EQ),
			NewFmlOr(
				NewFmlAnd(NewAtom(s, ops), NewAtom(t, opt.neg())),
				NewFmlAnd(NewAtom(t, opt), NewAtom(s, ops.neg()))))
		if a.op == NE {
			ret = ret.Not()
		}
		return ret, true
	case GT, LE, GE, LT:
		if aop == GE || aop == LT {
			aop = aop.neg()
			g = g.Neg().(*Poly)
		}
		if aop != GT && aop != LE {
			panic(fmt.Sprintf("invalid aop=%v", aop))
		}
		if a.op&EQ != 0 {
			aop = aop.not()
		}
		if aop != GT {
			panic(fmt.Sprintf("invalid op=%v", aop))
		}
		op := aop
		//fmt.Printf("u.sgn_lcp=%+d, u.sgn_s=%+d, a.Deg(u.lv)=%d, aop=%v: [%v %v 0]\n", u.sgn_lcp, u.sgn_s, a.Deg(u.lv), a, g, aop)
		if u.sgn_lcp < 0 && a.Deg(u.lv)%2 != 0 {
			// a2 の符号の影響を受けるパターン
			op = op.neg()
		}
		ops := op
		if u.sgn_s < 0 {
			ops = ops.neg()
		}
		r, t, s := quadeq_getrst(u.g, f, g, u.lv)
		sa := NewAtom(s, ops)
		ta := NewAtom(t, op)
		ret := NewFmlOrs(
			NewFmlAnd(NewAtom(r, op), ta),
			NewFmlAnd(NewAtom(r, op.neg()), sa),
			NewFmlAnd(ta, sa))
		if a.op&EQ != 0 {
			ret = ret.Not()
		}
		return ret, true
	default:
		panic(fmt.Sprintf("op=%d", a.op))
	}
}

/////////////////////////////////////////////////
// 共通部分
/////////////////////////////////////////////////

func (fof *Atom) qe_quadeq(fm func(a *Atom, p interface{}) Fof, p interface{}) Fof {
	return fm(fof, p)
}

func (qeopt QEopt) qe_quadeq(fof FofQ, cond qeCond) Fof {
	// fmt.Printf("qe_quadeq(%v)\n", fof)
	var op OP
	if _, ok := fof.(*Exists); ok {
		op = EQ
	} else {
		op = NE
	}

	dmin := 1
	dmax := 2
	if qeopt.Algo&QEALGO_EQLIN == 0 {
		dmin = 2
	}
	if qeopt.Algo&QEALGO_EQQUAD == 0 {
		dmax = 1
	}

	fff, ok := fof.Fml().(FofAO)
	if !ok {
		return nil
	}
	minatom := &quadEQ_t{}
	minatom.idx = -3

	for ii, fffi := range fff.Fmls() {
		if atom, ok := fffi.(*Atom); ok && atom.op == op {
			minatom.SetAtomIfEasy(fof, atom, dmin, dmax, ii)
		}
	}

	if minatom.a == nil {
		return nil
	}

	if op == NE {
		fff = fff.Not().(FofAO)
	}

	tbl := NewFofQuadEq(qeopt.g, minatom.p, minatom.lv)

	if minatom.deg == 2 {
		// minatom.deg == 2
		even := quadeq_isEven(fff, minatom.lv)
		discrim := NewAtom(qeopt.g.ox.Discrim(minatom.p, minatom.lv), GE)
		qeopt.log(cond, 2, "eq2", "%v [%v] discrim=%v, even=%v\n", fof, minatom.p, discrim, even)
		var o Fof = falseObj
		for _, sgns := range []struct {
			sgn_s int // 2つの根のうち，大きい方なら正.
			op    OP  // 主係数の符号
			skip  bool
		}{
			{+1, GT, false},
			{+1, LT, even},
			{-1, GT, false},
			{-1, LT, even},
		} {
			if sgns.skip {
				continue
			}
			tbl.sgn_s = sgns.sgn_s
			if sgns.op == GT {
				tbl.sgn_lcp = 1
			} else {
				tbl.sgn_lcp = -1
			}
			var aop OP
			if even {
				aop = NE
			} else {
				aop = sgns.op
			}
			ffx, _ := fff.Apply(qe_quadeq, tbl, true)
			opp := NewFmlAnds(ffx, NewAtom(minatom.z, aop), discrim)
			if false {
				cad, _ := funcCAD(qeopt.g, "CAD", []interface{}{opp})
				fmt.Printf("  opp=%v ==> %v\n", opp, cad)
			}
			o = NewFmlOr(o, opp)
		}

		if _, ok := minatom.z.(NObj); ok {
			if op == NE {
				o = o.Not()
			}
			return o
		}

		eq := NewAtom(minatom.z, EQ)
		fml := qeopt.g.simplFof(fff, eq, falseObj) // 等式制約で簡単化
		fml = NewFmlAnd(fml, eq)
		o = NewFmlOr(o, fml)
		if op == NE {
			o = o.Not()
		}
		return o
	}

	// 1次の等式制約の場合
	qeopt.log(cond, 2, "eq1", "%v [%v]\n", fff, minatom.p)
	tbl.sgn_lcp = 1
	ffp, _ := fff.Apply(qe_lineq, tbl, true)
	opos := NewFmlAnd(ffp, NewAtom(minatom.z, GT))

	tbl.sgn_lcp = -1
	ffn, _ := fff.Apply(qe_lineq, tbl, true)
	oneg := NewFmlAnd(ffn, NewAtom(minatom.z, LT))

	fs := make([]Fof, len(fff.Fmls())+1)
	copy(fs, fff.Fmls())
	c := minatom.p.Coef(tbl.lv, 0)
	fs[minatom.idx] = NewAtom(c, EQ)
	fs[len(fs)-1] = NewFmlAnd(fs[minatom.idx], NewAtom(minatom.z, EQ))

	ret := NewFmlOrs(opos, oneg, fff.gen(fs))
	if op == NE {
		ret = ret.Not()
	}
	return ret
}
