package ganrac

// Quantifier elimination for inequational constraints.
// Hidenao IWANE. 2015 (Japanese)

// ex([x], f(x) != 0 && phi(x))

import (
	"fmt"
)

func is_neq_only(fof Fof, lv Level) bool {
	switch pp := fof.(type) {
	case FofQ:
		return is_neq_only(pp.Fml(), lv)
	case FofAO:
		for _, f := range pp.Fmls() {
			if !is_neq_only(f, lv) {
				return false
			}
		}
		return true
	case *Atom:
		if pp.op == NE {
			return true
		}
		return !pp.hasVar(lv)
	}
	return false
}

func is_strict_only(fof Fof, lv Level) bool {
	switch pp := fof.(type) {
	case FofQ:
		return is_strict_only(pp.Fml(), lv)
	case FofAO:
		for _, f := range pp.Fmls() {
			if !is_strict_only(f, lv) {
				return false
			}
		}
		return true
	case *Atom:
		if pp.op == NE || pp.op == LT || pp.op == GT {
			return true
		}
		return !pp.hasVar(lv)
	}
	return false
}

/*
 * strict でないものが少ししかない, かつ, 2次以下
 * strict でないものが許されるのは depth=0 の atom な場合のみ
 */
func is_strict_or_quad(fof Fof, lv Level, depth int) bool {
	switch pp := fof.(type) {
	case FofQ:
		return is_strict_or_quad(pp.Fml(), lv, depth)
	case FofAO:
		for _, f := range pp.Fmls() {
			if !is_strict_or_quad(f, lv, depth+1) {
				if a, ok := f.(*Atom); ok && depth == 0 && a.Deg(lv) <= 2 && a.op != EQ {
					// 通常は hong93 を適用済みのはずなので，
					// EQ はなく，GE か LE は保証されているはず
					continue
				}
				return false
			}
		}
		return true
	case *Atom:
		if (pp.op & EQ) == 0 { // LT or GT or NE
			return true
		}
		return !pp.hasVar(lv)
	}
	return false
}

/*
 * 非等式制約部分とそれ以外で分割する
 */
func divide_neq(finput Fof, lv Level, qeopt QEopt) (Fof, Fof) {
	switch fof := finput.(type) {
	case *Atom:
		if qeopt.assert && fof.Deg(lv) == 0 {
			panic("lv not found")
		}
		if fof.op == NE {
			return fof, trueObj
		} else {
			return trueObj, fof
		}
	case *FmlAnd:
		fne := make([]Fof, 0, len(fof.Fmls()))
		fot := make([]Fof, 0, len(fof.Fmls()))
		for _, f := range fof.Fmls() {
			if is_neq_only(f, lv) {
				fne = append(fne, f)
			} else {
				fot = append(fot, f)
			}
		}
		if len(fne) == 0 {
			return trueObj, fof
		} else if len(fot) == 0 {
			return fof, trueObj
		} else {
			return NewFmlAnds(fne...), NewFmlAnds(fot...)
		}
	}
	return nil, nil
}

// atom に対する QE を行う; ex([lv], pp)
//
// 類似関数名 apply_neqQE_atom() と間違えないよう注意
func apply_neqQEatom(pp *Atom, lv Level) Fof {
	if pp.op != NE || !pp.hasVar(lv) {
		return pp
	}
	ands := make([]Fof, 0, len(pp.p))
	for _, p := range pp.p {
		deg := p.Deg(lv)
		ors := make([]Fof, deg+1)
		for d := 0; d <= deg; d++ {
			ors[d] = NewAtom(p.Coef(lv, uint(d)), NE)
		}
		ands = append(ands, NewFmlOrs(ors...))
	}
	return NewFmlAnds(ands...)
}

func apply_neqQE(fof Fof, lv Level) Fof {
	switch pp := fof.(type) {
	case FofAO:
		fmls := pp.Fmls()
		ret := make([]Fof, len(fmls))
		for i, f := range fmls {
			ret[i] = apply_neqQE(f, lv)
		}
		return pp.gen(ret)
	case *Atom:
		return apply_neqQEatom(pp, lv)
	case FofQ:
		return pp.gen(pp.Qs(), apply_neqQE(pp.Fml(), lv))
	}
	return nil
}

/*
* fof: inequational constraints
* atom: f <= 0 or f >= 0: f is univariate.
*
* Returns: qff which is equivalent to ex([x], f <= 0 && fof_neq)
:        : nil if fails
*/
func apply_neqQE_atom_univ(fof, qffneq Fof, atom *Atom, lv Level, qeopt QEopt, cond qeCond) Fof {
	// fmt.Printf("univ: %s AND %s\n", fof, atom)
	// atom.p は univariate
	// qffneq := apply_neqQE(fof, lv)
	p := atom.getPoly()

	// ex([x], sgn * f >= 0) がわかった
	if p.Sign() > 0 && atom.op == GE || p.Sign() < 0 && atom.op == LE {
		return qffneq
	}

	bak_deg := p.Deg(lv)
	if bak_deg%2 != 0 {
		return qffneq
	}

	ps := qeopt.g.ox.Factor(p)
	evens := make([]*Poly, 0, ps.Len())

	for i := 1; i < ps.Len(); i++ {
		_fr, _ := ps.Geti(i) // f^r
		fr := _fr.(*List)
		f := fr.getiPoly(0)
		r := fr.getiInt(1)

		rr, _ := f.RealRootIsolation(1)
		if rr.Len() == 0 {
			// ゼロ点がない => 符号一定
			continue
		}

		if r.Bit(0) == 0 { // r % 2 == 0
			// 有限個のゼロ点以外は符号を変えない
			evens = append(evens, f)
		} else {
			// 符号が変化する区間があることが確定
			return qffneq
		}
	}

	// 有限個のゼロ点でのみ条件を満たすことがわかった => 等式制約 QE へ.
	var ret Fof = falseObj
	for _, z := range evens {
		f := z
		ret = NewFmlOr(ret, qeopt.qe(NewQuantifier(false, []Level{lv},
			NewFmlAnd(NewAtom(f, EQ), fof)), cond))
	}

	return ret
}

func apply_neqQE_pstrict(fof, fne Fof, fot_fmls []Fof, lv Level, qeopt QEopt, cond qeCond) Fof {
	// GE, LE が一部含まれるが，それは２次以下なので，
	// hong93 により高速に解ける.
	// fne が大きければ，全体を解くよりも strict 部分のみの CAD になるのでお得

	// fot_fmls := fot.Fmls() // 壊してはだめ
	n := 1
	if p, ok := fne.(*FmlAnd); ok {
		n = p.Len()
	}
	fmls := make([]Fof, len(fot_fmls)+n)
	copy(fmls, fot_fmls)
	if n == 1 {
		fmls[len(fot_fmls)] = fne
	} else {
		p := fne.(*FmlAnd)
		copy(fmls[len(fot_fmls):], p.Fmls())
	}

	ors := make([]Fof, 0, len(fmls)+1)
	for i, v := range fot_fmls {
		if a, ok := v.(*Atom); ok && a.op&EQ != 0 {
			// AND の要素のうちひとつの GE/LE を 等式制約に変えた論理式にする
			fmls[i] = newAtoms(a.p, EQ)
			f := NewExists([]Level{lv}, NewFmlAnds(fmls...))
			ors = append(ors, f)
			fmls[i] = v
		}
	}

	// 全部 strict な場合
	for i, v := range fot_fmls {
		if a, ok := v.(*Atom); ok && a.op&EQ != 0 {
			fmls[i] = newAtoms(a.p, a.op.strict())
		}
	}
	fmls = fmls[:len(fot_fmls)]
	fstrict := NewExists([]Level{lv}, NewFmlAnds(fmls...))
	ff := NewFmlAnd(apply_neqQE(fne, lv), fstrict)
	ors = append(ors, ff)
	newfml := NewFmlOrs(ors...)
	qeopt.log(cond, 3, "neq", "<%s> pstricto %v\n", VarStr(lv), newfml)
	return qeopt.qe(newfml, cond)
}

//	apply_neqQE_atom(*Atom, Level): 非等式でない部分が atom だった
//
// fof: inequational constraints
// atom: f <= 0 or f >= 0
//
// Returns: qff which is equivalent to ex([x], f <= 0 && fof_neq)
//
//	: nil if fails
//
// see apply_neqQEatom(); similar function
func apply_neqQE_atom(fof Fof, atom *Atom, lv Level, qeopt QEopt, cond qeCond) Fof {
	// fmt.Printf("atom: %s AND %s\n", fof, atom)
	if qeopt.assert && (atom.op&EQ) == 0 {
		panic(fmt.Sprintf("unexpected op %d, expected [%d,%d]", atom.op, GE, LE))
	}

	var ret Fof = falseObj
	poly := atom.getPoly()
	for {
		qffneq := apply_neqQE(fof, lv)

		deg := poly.Deg(lv)
		// fmt.Printf("atom.poly[%d]=%v\n", deg, poly)
		lc := poly.Coef(lv, uint(deg))

		lccond := NewAtom(lc, atom.op)
		lccond = qeopt.simplify(lccond, cond)
		if _, ok := lccond.(*AtomT); ok {
			lccond := NewAtom(lc, atom.op.strict())
			lccond = qeopt.simplify(lccond, cond)
			if _, ok := lccond.(*AtomT); ok {
				ret = NewFmlOr(ret, qffneq)
				return ret
			}
			ret = NewFmlOr(ret, NewFmlAnd(NewAtom(lc, atom.op.strict()), qffneq))
		} else if deg%2 != 0 {
			ret = NewFmlOr(ret, NewFmlAnd(NewAtom(lc, NE), qffneq))
		} else if deg == 0 {
			return NewFmlOr(ret, NewFmlAnd(NewAtom(lc, atom.op), qffneq))
		} else if deg == 2 {
			discrim := poly.discrim2(lv)
			op := LT
			if atom.op == GE {
				op = GT
			}
			// ex([x], ax^2+bx+c >= 0)
			// <==>
			// infinite: a > 0 || b^2-4ac > 0 || (a=0 && b=0 && c >= 0)
			// __finite: b^2-4ac=0 /\ a !=0
			c1 := poly.Coef(lv, 1)
			c0 := poly.Coef(lv, 0)
			ret = NewFmlOr(ret, NewFmlAnd(NewAtom(lc, op), qffneq))
			ret = NewFmlOr(ret, NewFmlAnd(NewAtom(discrim, GT), qffneq))
			ret = NewFmlOr(ret, NewFmlAnds(NewAtom(lc, EQ), NewAtom(c1, EQ), NewAtom(c0, atom.op), qffneq))
			qq := qeopt.qe(NewExists([]Level{lv},
				NewFmlAnd(fof,
					NewAtom(Add(Mul(Mul(two, lc), NewPolyVar(lv)), c1), EQ))), cond)
			// fmt.Printf("qq=%v\n", qq)
			ret = NewFmlOr(ret, NewFmlAnds(
				// ex([x], 2ax+b=0 && b^2-4ac = 0 && a != 0 && NEQ)
				NewAtom(lc, NE), // atom.op とどちらが良いか.
				NewAtom(discrim, EQ),
				qq))
			return ret

		} else if poly.IsUnivariate() {
			r := apply_neqQE_atom_univ(fof, qffneq, NewAtom(poly, atom.op).(*Atom), lv, qeopt, cond)
			if r == nil {
				return nil
			}
			ret = NewFmlOr(ret, r)
			return ret
		} else {
			return nil
		}

		fof = NewFmlAnd(fof, NewAtom(lc, EQ))
		fof = qeopt.simplify(fof, cond)
		if _, ok := fof.(*AtomF); ok {
			return ret
		}
		switch pp := poly.Sub(Mul(lc, NewPolyVarn(lv, deg))).(type) {
		case *Poly:
			poly = pp
		default:
			// 定数になった...
			ret = NewFmlOr(ret, NewFmlAnd(NewAtom(pp, atom.op), qffneq))
			return ret
		}
	}
}

/*
 * fof: prenex first-order formula ex([x, y, z], p1 & p2 & ... & p2)
 */
func neqQE(fof Fof, lv Level, qeopt QEopt, cond qeCond) Fof {
	fne, fot := divide_neq(fof, lv, qeopt)
	if !fne.hasVar(lv) {
		return fof
	}

	if fot == trueObj {
		qeopt.log(cond, 3, "neq", "<%s> all %v\n", VarStr(lv), fof)
		return apply_neqQE(fof, lv)
	}
	if fne == trueObj {
		return fof
	}

	// @TODO fot が OR になったら分解可能

	if is_strict_only(fot, lv) {
		qeopt.log(cond, 3, "neq", "<%s> strict [%v] %v\n", VarStr(lv), fne, fof)
		fstrict := NewQuantifier(false, []Level{lv}, fot)
		return NewFmlAnd(apply_neqQE(fne, lv), qeopt.qe(fstrict, cond))
	}
	if atom, ok := fot.(*Atom); ok {
		if atom.op == EQ { // Hong93 とかで対応可能
			return fof
		}
		qeopt.log(cond, 3, "neq", "<%s> atom %v\n", VarStr(lv), fof)
		return apply_neqQE_atom(fne, atom, lv, qeopt, cond)
	}
	if qeopt.Algo&(QEALGO_EQLIN|QEALGO_EQQUAD) != 0 && is_strict_or_quad(fot, lv, 0) {
		// @TODO fne がそれなりに複雑である場合に限定したほうが良いかも
		if fotand, ok := fot.(*FmlAnd); ok {
			qeopt.log(cond, 3, "neq", "<%s> pstrict and[%d] %v\n", VarStr(lv), len(fotand.fml), fof)
			return apply_neqQE_pstrict(fof, fne, fotand.fml, lv, qeopt, cond)
		}
		if fotatom, ok := fot.(*Atom); ok {
			qeopt.log(cond, 3, "neq", "<%s> pstrict atom %v\n", VarStr(lv), fof)
			return apply_neqQE_pstrict(fof, fne, []Fof{fotatom}, lv, qeopt, cond)
		}
	}

	return fof
}

func (qeopt QEopt) qe_neq(fof FofQ, cond qeCond) Fof {
	var fml Fof
	not := false
	switch pp := fof.(type) {
	case *ForAll:
		fml = pp.fml.Not()
		not = true
	case *Exists:
		fml = pp.fml
	default:
		return fof
	}

	for _, q := range fof.Qs() {
		ff := neqQE(fml, q, qeopt, cond)
		if ff != fml {
			if not {
				return ff.Not()
			} else {
				return ff
			}
		}
	}
	return nil
}
