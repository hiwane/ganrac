package ganrac

// Applying Linear Quantifier Elimination
// R. Loos, V. Weispfenning 1993

// [Weispfenning94]
// Quantifier Elimination for Real Algebra -- the Quadratic Case and Beyond
// V. Weispfenning 1994

// A Generalized Framework for Virtual Substitution
// M. Kosta, T. Sturm 2015

import (
	"fmt"
)

const OPVS_SHIFT = 8

type fof_vser interface {
	apply_vs(fm func(atom *Atom, p *vs_sample_point, lv Level) Fof, p *vs_sample_point, lv Level) Fof
}

type vs_sample_point struct {
	// (num + sqrt(v)) / den
	num    RObj
	den    []RObj // [den^0, den, den^2, den^3, ...]
	sqr    RObj
	lc     RObj
	deg    int
	sqrc   int
	densgn int
	// 線形なら +1, ２次なら +1 or -1.
	// 一番大きな根なら +1 = 主係数の符号と x+e の符号が一致
	// 一番小さな根なら -1 = 主係数の符号と x+e の符号が一致しない

	idx    int
	neccon Fof
}

func (sp *vs_sample_point) Format(s fmt.State, format rune) {
	fmt.Fprintf(s, "[sp, (")
	if sp.sqr != nil {
		fmt.Fprintf(s, "(")
		sp.num.Format(s, format)
		op := "+"
		if sp.sqrc < 0 {
			op = "-"
		}
		fmt.Fprintf(s, ")%ssqrt(", op)
		sp.sqr.Format(s, format)
		fmt.Fprintf(s, ")")
	} else {
		sp.num.Format(s, format)
	}
	fmt.Fprintf(s, ")/(")
	sp.den[1].Format(s, format)
	fmt.Fprintf(s, "), @densgn=%+d,idx=%+d,lc=", sp.densgn, sp.idx)
	sp.lc.Format(s, format)
	fmt.Fprintf(s, ",nec=")
	sp.neccon.Format(s, format)
	fmt.Fprintf(s, ",den=[")
	for i := 0; i < len(sp.den); i++ {
		sp.den[i].Format(s, format)
		fmt.Fprintf(s, ",")
	}
	fmt.Fprintf(s, "]")
}

func makeDenAry(v RObj, deg int) []RObj {
	ret := make([]RObj, deg+1)
	ret[0] = one
	ret[1] = v
	for i := 2; i <= deg; i++ {
		ret[i] = v.Mul(ret[i-1])
	}
	return ret
}

func newVslinSamplePoint(deg, maxd int, num, den RObj) *vs_sample_point {
	sp := new(vs_sample_point)
	sp.num = num
	sp.deg = deg
	if den.IsNumeric() {
		sp.densgn = den.Sign()
	} else {
		sp.densgn = 0
	}
	sp.lc = den
	sp.num = sp.num.Neg()

	sp.idx = 1
	sp.den = makeDenAry(den, maxd)
	sp.neccon = trueObj
	return sp
}

func (sp *vs_sample_point) DenSign() int {
	return sp.densgn
}

func (sp *vs_sample_point) SetDenSign(sgn int) {
	sp.densgn = sgn
}

func (sp *vs_sample_point) LC() RObj {
	return sp.lc
}

func (sp *vs_sample_point) necessaryCondition() Fof {
	return sp.neccon
}

func (sp *vs_sample_point) isRational() bool {
	return sp.sqr == nil || sp.sqr.IsZero()
}

type vs_elimination_set struct {
	ps []*vs_elimination_elem
}

type vs_elimination_elem struct {
	p *Poly
	// EQ if =, <=, >=
	// GT if >, !=
	// LT if <, !=
	kind int
}

func newVsEliminationSet() *vs_elimination_set {
	p := new(vs_elimination_set)
	p.ps = make([]*vs_elimination_elem, 0)
	return p
}

func (es *vs_elimination_set) append(p *Poly, oop OP) {
	op := int(oop)
	for _, qq := range es.ps {
		if qq.p.Equals(p) {
			qq.kind |= (1 << (OPVS_SHIFT + op)) | op
			return
		}
	}
	q := new(vs_elimination_elem)
	q.p = p
	q.kind = (1 << (OPVS_SHIFT + op)) | op
	es.ps = append(es.ps, q)
}

// get_vs_polys から呼び出される
func (es *vs_elimination_set) addAtom(atom *Atom, lv Level) {
	for _, pol := range atom.p {
		if pol.hasVar(lv) {
			if len(atom.p) == 1 {
				es.append(pol, atom.op)
			} else if atom.op&EQ != 0 {
				es.append(pol, GE)
				es.append(pol, LE)
			} else {
				es.append(pol, NE)
			}
		}
	}
}

func (fof *AtomT) apply_vs(fm func(atom *Atom, p *vs_sample_point, lv Level) Fof, p *vs_sample_point, lv Level) Fof {
	return fof
}
func (fof *AtomF) apply_vs(fm func(atom *Atom, p *vs_sample_point, lv Level) Fof, p *vs_sample_point, lv Level) Fof {
	return fof
}

func (fof *Atom) apply_vs(fm func(atom *Atom, p *vs_sample_point, lv Level) Fof, ptt *vs_sample_point, lv Level) Fof {
	if fof.hasVar(lv) {
		return fm(fof, ptt, lv)
	} else {
		return fof

	}
}

func (fof *FmlAnd) apply_vs(fm func(atom *Atom, p *vs_sample_point, lv Level) Fof, p *vs_sample_point, lv Level) Fof {
	var r Fof = trueObj

	for _, f := range fof.fml {
		r = NewFmlAnd(r, f.apply_vs(fm, p, lv))
	}

	return r
}

func (fof *FmlOr) apply_vs(fm func(atom *Atom, p *vs_sample_point, lv Level) Fof, p *vs_sample_point, lv Level) Fof {
	var r Fof = falseObj

	for _, f := range fof.fml {
		r = NewFmlOr(r, f.apply_vs(fm, p, lv))
	}

	return r
}

func (fof *ForAll) apply_vs(fm func(atom *Atom, p *vs_sample_point, lv Level) Fof, p *vs_sample_point, lv Level) Fof {
	fmt.Printf("forall %v\n", fof)
	panic("invalid......... apply_vs(forall)")
}

func (fof *Exists) apply_vs(fm func(atom *Atom, p *vs_sample_point, lv Level) Fof, p *vs_sample_point, lv Level) Fof {
	fmt.Printf("exists %v\n", fof)
	panic("invalid......... apply_vs(exists)")
}

// サンプル点を生成する.
// 2 次の場合に 3 サンプルになるので slice で返す
func (vee *vs_elimination_elem) vs_make_samples(lv Level, maxd int) []*vs_sample_point {
	p := vee.p
	if p.Deg(lv) == 1 {
		sp := gen_sample_vslin(p, lv, maxd)
		return []*vs_sample_point{sp}
	}
	if p.Deg(lv) == 2 {
		return gen_sample_vsquad(p, lv, maxd)
	}
	return nil
}

func gen_sample_vsquad(p *Poly, lv Level, maxd int) []*vs_sample_point {
	sps := make([]*vs_sample_point, 0, 3)

	c0 := p.Coef(lv, 0)
	c1 := p.Coef(lv, 1)
	c2 := p.Coef(lv, 2)

	// c2=0 の場合は，1次のものに条件付き
	if !c2.IsNumeric() {
		sp := newVslinSamplePoint(1, maxd, c0, c1)
		sp.neccon = NewAtom(c2, EQ)
		sps = append(sps, sp)
	}

	// 2次の場合は，解の公式を用いて...
	sp1 := newVslinSamplePoint(2, maxd, c1, Mul(c2, two))
	sp1.sqr = Sub(c1.Mul(c1), Mul(Mul(c2, c0), four)) // 混合部分は判別式
	sp1.neccon = NewAtom(sp1.sqr, GE)                 // 判別式が非負
	sp1.sqrc = 1
	sp1.deg = 2
	sp1.lc = c2 // 2倍しているのが嫌い
	sp1.idx = 1

	sps = append(sps, sp1)

	sp2 := new(vs_sample_point)
	*sp2 = *sp1
	sp2.sqrc = -1
	sp2.idx = -sp1.idx
	sps = append(sps, sp2)

	return sps
}

func gen_sample_vslin(p *Poly, lv Level, maxd int) *vs_sample_point {
	num := p.Coef(lv, 0)
	den := p.Coef(lv, 1)
	return newVslinSamplePoint(1, maxd, num, den)
}

func virtual_subst(atom *Atom, ptt *vs_sample_point, lv Level) Fof {
	// 線形なサンプル点の代入
	// 分母は非ゼロと仮定し，
	// その符号は pt.densgn
	if ptt.isRational() {
		return ptt.virtual_subst(atom, lv)
	} else {
		return ptt.virtual_subst_sqr(atom, lv)
	}
}

func (pt *vs_sample_point) virtual_subst_sqr(atom *Atom, lv Level) Fof {
	d := atom.Deg(lv)
	if d <= 0 {
		return atom
	}
	op := atom.op
	if pt.densgn < 0 && atom.Deg(lv)%2 != 0 {
		op = op.neg()
	}
	// fmt.Printf("      virtual_subst_sqr densgn=%d, Deg=%d, atom=%v : %v: \n", pt.densgn, atom.Deg(lv), atom, op)

	p := atom.getPoly()
	q, r := p.subst_frac_sqr(pt.num, pt.sqr, pt.sqrc, pt.den[:d+1], lv)
	if err := q.valid(); err != nil {
		panic(err)
	}
	if op == GE || op == GT {
		op = op.neg()
		q = q.Neg()
		r = r.Neg()
	}

	q2_r2x := Sub(Mul(q, q), Mul(Mul(r, r), pt.sqr))
	switch op {
	case EQ:
		return NewFmlAnds(
			NewAtoms([]RObj{q, r}, LE),
			NewAtom(q2_r2x, EQ))
	case NE:
		return NewFmlOrs(
			NewAtoms([]RObj{q, r}, GT),
			NewAtom(q2_r2x, NE),
		)
	case LE:
		return NewFmlOrs(
			NewFmlAnd(NewAtom(q, LE), NewAtom(q2_r2x, GE)),
			NewFmlAnd(NewAtom(r, LE), NewAtom(q2_r2x, LE)),
		)
	case LT:
		qlt := NewAtom(q, LT)
		rle := NewAtom(r, LE)
		return NewFmlOrs(
			NewFmlAnd(qlt, NewAtom(q2_r2x, GT)),
			NewFmlAnd(qlt, rle),
			NewFmlAnd(rle, NewAtom(q2_r2x, LT)),
		)
	default:
		panic(fmt.Sprintf("why? op=%v, lv=%d, pt=%v, atom=%v", op, lv, pt, atom))
	}
}

func (pt *vs_sample_point) virtual_subst(atom *Atom, lv Level) Fof {
	op := atom.op
	if pt.densgn < 0 && atom.Deg(lv)%2 != 0 {
		op = op.neg()
	}

	pp := make([]RObj, len(atom.p))
	for i, p := range atom.p {
		d := p.Deg(lv)
		pp[i] = p.subst_frac(pt.num, pt.den[:d+1], lv)
		if err := pp[i].valid(); err != nil {
			panic(err)
		}
		//			fmt.Printf("pp[%d]=%v\n", i, pp[i])
	}
	//		fmt.Printf("virtual_subst_lin() %v -> %v<%d>::%v\n", atom, pp, atom.op, NewAtoms(pp, atom.op))
	return NewAtoms(pp, op)
}

// [Weispfenning94] p.89
// (f < 0)[e + epsilon / x] = nu(f) [e / x]
func vs_nu(polys []*Poly, op OP, pt *vs_sample_point, lv Level) Fof {
	// pt + epsilon を代入する
	d := 0
	if pt.densgn < 0 {
		for _, p := range polys {
			d += p.Deg(lv)
		}
	}

	var opx OP
	if d%2 == 0 || true {
		// 偶数次数か，分母の符号が正ならそのまま.
		opx = op
	} else {
		opx = op.neg()
	}
	f1 := virtual_subst(newAtoms(polys, opx), pt, lv)

	// fmt.Printf("    vsnu(): f1=%v, \tatom=%v, lv=%v\n", f1, newAtoms(polys, op), lv)
	if err := f1.valid(); err != nil { // debug
		fmt.Printf("%V\ninvalid f1 %v\natom=%v\npt=%v\nlv=%d\n", f1, f1, newAtoms(polys, op), pt, lv)
		panic(err.Error())
	}

	rs := make([]RObj, 0, len(polys))
	var pmul *Poly = nil
	for i, p := range polys {
		if !p.hasVar(lv) {
			rs = append(rs, p)
		} else if pmul == nil {
			pmul = polys[i]
		} else {
			pmul = pmul.Mul(polys[i]).(*Poly)
		}
	}
	if pmul == nil {
		return f1
	}

	var f2 Fof = virtual_subst(newAtoms(polys, EQ), pt, lv)
	// fmt.Printf("    vsnu(): f2=%v\n", f2)
	if err := f2.valid(); err != nil { // debug
		fmt.Printf("%V\ninvalid f2 %v\n", f1, f1)
		panic(err.Error())
	}

	var v RObj
	if pmul == nil {
		v = zero
	} else {
		v = pmul.Diff(lv)
	}
	// fmt.Printf("    vsnu(): pmul=%v, v=%v: %v\n", pmul, v, v.IsNumeric())

	if v.IsNumeric() {
		// fmt.Printf("vsnu(): rs=%v, %d\n", rs, v.Sign())
		if v.Sign() < 0 {
			op = op.neg()
		}
		return NewFmlOr(f1, NewFmlAnd(f2, NewAtoms(rs, op)))
	} else {
		rs = append(rs, v)
		atm := NewAtoms(rs, op).(*Atom)

		// fmt.Printf("    vsnu(): ps=%v\n", rs)
		sfml := vs_nu(atm.p, atm.op, pt, lv)
		// fmt.Printf("    vsnu(): sfml=%v\n", sfml)
		return NewFmlOr(f1, NewFmlAnd(f2, sfml))
	}
}

func virtual_subst_e(atom *Atom, ptt *vs_sample_point, lv Level) Fof {
	// ptt+ infinitesimal を代入する
	return ptt.virtual_subst_e(atom, lv)
}

func (pt *vs_sample_point) virtual_subst_e(atom *Atom, lv Level) Fof {
	if atom.op == EQ {
		var ret Fof = falseObj
		for _, p := range atom.p {
			var pi Fof = trueObj
			d := p.Deg(lv)
			for i := 0; i <= d; i++ {
				c := p.Coef(lv, uint(i))
				pi = NewFmlAnd(pi, NewAtom(c, EQ))
			}
			ret = NewFmlOr(ret, pi)
		}
		return ret
	} else if atom.op == NE {
		return pt.virtual_subst_e(newAtoms(atom.p, EQ), lv).Not()
	} else if atom.op == LT || atom.op == GT {
		if err := atom.valid(); err != nil {
			fmt.Printf("%V\ninvalid atom %v\n", atom, atom)
			panic(err.Error())
		}
		// fmt.Printf("    go vs_nu(%v, %v)\n", atom, pt)
		ret := vs_nu(atom.p, atom.op, pt, lv)
		// fmt.Printf("    vs_nu(LT|GT) %v\n", ret)
		return ret
	} else if atom.op == LE || atom.op == GE {
		ret := pt.virtual_subst_e(newAtoms(atom.p, atom.op.not()), lv).Not()
		// fmt.Printf("    vs_nu(LE|GE) %v\n", ret)
		return ret
	} else {
		panic("invalid op")
	}
}

// [Weispfenning94] p.90
// (f < 0)[-inf / x] = mu(f)
func vs_mu(atom *Atom, lv Level) Fof {
	// polys < 0 の atom に対して -infty を代入する

	pmul := atom.getPoly()

	d := pmul.Deg(lv)

	coeffs := make([]RObj, d+1)
	for i := 0; i <= d; i++ {
		coeffs[i] = pmul.Coef(lv, uint(i))
	}

	var ret Fof = falseObj
	for i := 0; i <= d; i++ {
		var u Fof
		if i%2 == 0 {
			u = NewAtom(coeffs[i], LT)
		} else {
			u = NewAtom(coeffs[i], GT)
		}
		for j := i + 1; j <= d; j++ {
			u = NewFmlAnd(u, NewAtom(coeffs[j], EQ))
		}
		ret = NewFmlOr(ret, u)
	}
	return ret
}

func virtual_subst_i(atom *Atom, ptt *vs_sample_point, lv Level) Fof {
	return virtual_subst_i_lv(atom, lv)
}

func virtual_subst_i_lv(atom *Atom, lv Level) Fof {
	// -inf を代入する.
	switch atom.op {
	case EQ:
		var ret Fof = falseObj
		for _, p := range atom.p {
			var pi Fof = trueObj
			d := p.Deg(lv)
			for i := 0; i <= d; i++ {
				c := p.Coef(lv, uint(i))
				pi = NewFmlAnd(pi, NewAtom(c, EQ))
			}
			ret = NewFmlOr(ret, pi)
		}
		return ret
	case NE:
		return virtual_subst_i_lv(newAtoms(atom.p, EQ), lv).Not()
	case LT:
		return vs_mu(atom, lv)
	case LE:
		return NewFmlOr(
			virtual_subst_i_lv(newAtoms(atom.p, EQ), lv),
			virtual_subst_i_lv(newAtoms(atom.p, LT), lv))
	case GT:
		return NewFmlOr(
			virtual_subst_i_lv(newAtoms(atom.p, EQ), lv),
			virtual_subst_i_lv(newAtoms(atom.p, LT), lv)).Not()
	case GE:
		return virtual_subst_i_lv(newAtoms(atom.p, LT), lv).Not()
	default:
		panic("invalid op")
	}
}

func get_vs_polys(p Fof, lv Level) *vs_elimination_set {
	// algorithm 2 is not implemented
	// p に含まれる lv変数を含む多項式のリスト.
	peqlt := newVsEliminationSet()

	pstack := make([]Fof, 1)
	pstack[0] = p
	for len(pstack) > 0 {
		qq := pstack[len(pstack)-1]
		switch q := qq.(type) {
		case *ForAll:
			pstack[len(pstack)-1] = q.fml
		case *Exists:
			pstack[len(pstack)-1] = q.fml
		case *Atom:
			pstack = pstack[:len(pstack)-1]
			peqlt.addAtom(q, lv)
		case *FmlAnd:
			pstack[len(pstack)-1] = q.fml[0]
			for i := 1; i < len(q.fml); i++ {
				pstack = append(pstack, q.fml[i])
			}
		case *FmlOr:
			pstack[len(pstack)-1] = q.fml[0]
			for i := 1; i < len(q.fml); i++ {
				pstack = append(pstack, q.fml[i])
			}
		}
	}
	return peqlt
}

// テストから呼び出すための関数
func vsLinear(fof Fof, lv Level) Fof {
	return vs_main(fof, lv, 1, NewGANRAC())
}

// target_deg in [1, 2]: vs 対象とする最大次数
// 2 を指定したときは，線形しかないケースは適用対象外
// *Ganrac はログ出力用
func vs_main(fof Fof, lv Level, target_deg int, gan *Ganrac) Fof {
	loglv := 9
	maxd := fof.Deg(lv)
	if maxd != target_deg {
		return fof
	}

	var fml Fof
	switch pp := fof.(type) {
	case *ForAll:
		fml = pp.fml.Not()
	case *Exists:
		fml = pp.fml
	default:
		return fof
	}

	if !fml.IsQff() { // 一番内側であること
		return fof
	}

	var ret Fof = falseObj
	elset := get_vs_polys(fml, lv)
	// fmt.Printf("vsLin=%v\n", fml)
	// fmt.Printf("peq[%d]=%v\n", len(elset.equ), elset.equ)
	// fmt.Printf("plt[%d]=%v\n", len(elset.ine), elset.ine)
	required_zero := true
	required_minf := false

	const OPVS_LT = 1 << (LT + OPVS_SHIFT)
	const OPVS_GT = 1 << (GT + OPVS_SHIFT)
	const OPVS_LE = 1 << (LE + OPVS_SHIFT)
	const OPVS_GE = 1 << (GE + OPVS_SHIFT)
	const OPVS_EQ = 1 << (EQ + OPVS_SHIFT)
	const OPVS_NE = 1 << (NE + OPVS_SHIFT)

	for i, pp := range elset.ps {
		for j, pt := range pp.vs_make_samples(lv, maxd) {
			sgn := pt.DenSign()
			if sgn != 0 && pt.neccon.Equals(trueObj) { // 主係数が定数, かつ, 判別式が必ず非負/線形
				required_zero = false // もう 0 評価は不要
			}

			gan.log(loglv, 1, " @@  i,j=%d,%d, %v,%#x, zero=%v, sgn=%d, %v\n", i, j, pt, pp.kind, required_zero, sgn, pp.p)

			for _, stbl := range []struct {
				do  bool // 実行するか
				sgn int  // lc の符号
				op  OP   // lc の符号に対応する OP
			}{
				{sgn >= 0, +1, GT},
				{sgn <= 0, -1, LT},
			} {
				if !stbl.do {
					continue
				}
				pt.SetDenSign(stbl.sgn)

				lc := NewAtom(pt.LC(), stbl.op)
				if pp.kind&int(EQ) != 0 {
					// >=, <=, == のどれかまたは複数だが，>=, <= のみの場合は主係数の符号を制限できる
					// しかし，主係数の符号を制限する場合には，-inf の追加が必要になる
					if false && pp.kind&(OPVS_LE|OPVS_EQ) == 0 {
						// GE のみ ... sgnl < 0
						if stbl.sgn*pt.idx > 0 {
							goto _NEXT_GT
						}
					} else if false && pp.kind&(OPVS_GE|OPVS_EQ) == 0 {
						// LE のみ ... sgnl > 0
						if stbl.sgn*pt.idx < 0 {
							goto _NEXT_GT
						}
					}
					sfml := fml.apply_vs(virtual_subst, pt, lv)
					gan.log(loglv, 1, "add2:=:%d:%d: {%v}  [lc %v][nec %v]\n", stbl.sgn, pt.idx, sfml, lc, pt.necessaryCondition())
					if err := sfml.valid(); err != nil {
						panic(err)
					}
					ret = NewFmlOr(ret, NewFmlAnds(sfml, lc, pt.necessaryCondition()))
				}

			_NEXT_GT:
				if pp.kind&(OPVS_GT|OPVS_NE|OPVS_LT) != 0 {
					required_minf = true
				}

				// 微分が非正である.
				//   1次の場合.... 主係数 < 0
				//   2次の場合.... idx < 0
				if pp.kind&(OPVS_LT|OPVS_NE) != 0 && (pt.deg == 1 && stbl.sgn < 0 || pt.deg == 2 && pt.idx < 0) {
					// p < 0 ==> signr < 0
					sfml := fml.apply_vs(virtual_subst_e, pt, lv)
					gan.log(loglv, 1, "add6:<e:%d:%d {%v}  [lc %v][nec %v]\n", stbl.sgn, pt.idx, sfml, lc, pt.necessaryCondition())
					if err := sfml.valid(); err != nil {
						panic(err)
					}
					ret = NewFmlOr(ret, NewFmlAnds(sfml, lc, pt.necessaryCondition()))
					if err := ret.valid(); err != nil {
						panic(err)
					}
				}
				if pp.kind&(OPVS_GT|OPVS_NE) != 0 && (pt.deg == 1 && stbl.sgn > 0 || pt.deg == 2 && pt.idx > 0) {
					// p > 0 ==> sgnr > 0
					sfml := fml.apply_vs(virtual_subst_e, pt, lv)
					gan.log(loglv, 1, "add7:>e:%d:%d {%v}  [lc %v][nec %v]\n", stbl.sgn, pt.idx, sfml, lc, pt.necessaryCondition())
					if err := sfml.valid(); err != nil {
						panic(err)
					}
					ret = NewFmlOr(ret, NewFmlAnds(sfml, lc, pt.necessaryCondition()))
					if err := ret.valid(); err != nil {
						panic(err)
					}
				}
			}
		}
	}

	if err := ret.valid(); err != nil {
		panic(err)
	}
	if required_minf { // -inf
		pt := new(vs_sample_point) // ダミー
		sfml := fml.apply_vs(virtual_subst_i, pt, lv)
		gan.log(loglv, 1, "j=x, [-inf] {%v}\n", sfml)
		if err := sfml.valid(); err != nil {
			panic(err)
		}
		// fmt.Printf("before\nret%x= %v\nsfm%x= %v\n", ret.fofTag(), ret, sfml.fofTag(), sfml)
		ret = NewFmlOr(ret, sfml)
		if err := ret.valid(); err != nil {
			ppp, ok := ret.(*FmlOr)
			if ok {
				fmt.Printf("len=%d\n", ppp.Len())
			}

			panic(err)
		}
	}
	if required_zero { // サンプル点の分母がパラメータのみの場合に必要
		sfml := fml.Subst(zero, lv)
		gan.log(loglv, 1, "j=x, [zero] {%v}\n", sfml)
		ret = NewFmlOr(ret, sfml)
		if err := ret.valid(); err != nil {
			panic(err)
		}
	}

	if q, ok := fof.(*ForAll); ok {
		ret = ret.Not()
		ret = NewQuantifier(true, q.q, ret)
	} else if q, ok := fof.(*Exists); ok {
		ret = NewQuantifier(false, q.q, ret)
	}

	// fmt.Printf("LinearVS ret=%v\n", ret)
	return ret
}

func (qeopt QEopt) qe_vs(fof FofQ, cond qeCond, d int) Fof {
	for _, q := range fof.Qs() {
		ff := vs_main(fof, q, d, qeopt.g)
		if ff != fof {
			return ff
		}
	}
	return nil
}
