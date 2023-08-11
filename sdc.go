package ganrac

//
// QE for atom.
//  ex([x], f(x) <= 0)
//
// QE for sign definite condition.
//   ex([x], x >= 0 && f(x) <= 0)

import (
	"fmt"
)

func atomQEconstruct(ret []Fof, sres *List, atom_formulas [][]OP, neccon Fof) []Fof {
	if _, ok := neccon.(*AtomF); ok {
		return ret
	}
	// fmt.Printf("atomQEconstruct(%v) %v\n", sres, neccon)

	ps := make([]RObj, sres.Len())
	for i := 0; i < len(ps); i++ {
		ps[i] = sres.geti(i).(RObj)
	}

	for _, afml := range atom_formulas {
		r := make([]Fof, 1, len(afml)+1)
		r[0] = neccon
		for j, v := range afml {
			if v != OP_TRUE {
				r = append(r, NewAtom(ps[len(afml)-j-1], v))
			}
		}
		ret = append(ret, NewFmlAnds(r...))
	}
	// fmt.Printf("atomQEconstruct() returns %v\n", ret)
	return ret
}

// atomQE. 変数を消去できなかったときは nil を返す
func atomQE(atom *Atom, lv Level, neccon Fof, qeopt QEopt) Fof {
	if _, ok := neccon.(*AtomF); ok {
		return falseObj
	}

	p := atom.getPoly()
	deg := p.Deg(lv)

	// fmt.Printf("atomQE(%v, %v, nec=%v), deg=%v\n", atom, VarStr(lv), neccon, deg)
	if deg <= 0 {
		return NewFmlAnd(atom, neccon)
	}

	lc := p.Coef(lv, uint(deg))
	cont := p.Coef(lv, 0)

	switch atom.op {
	case GE, GT:
		if lc.IsNumeric() && lc.Sign() > 0 {
			return neccon
		}
		if cont.IsNumeric() && cont.Sign() > 0 {
			return neccon
		}
	case LE, LT:
		if lc.IsNumeric() && lc.Sign() < 0 {
			return neccon
		}
		if cont.IsNumeric() && cont.Sign() < 0 {
			return neccon
		}
	case NE: // neqQE におまかせしたいところですが.
		return NewFmlAnd(neccon, apply_neqQEatom(atom, lv))
	}

	if deg%2 != 0 {
		// 奇数次なら，すべての符号をもつので，必ず真である．
		if lc.IsNumeric() {
			return neccon
		}

		qq := p.setZero(lv, deg)
		// 主係数 0 の場合を考える.
		// @TODO ここで simplify する.
		if false {
			// neccon := NewFmlAnd(cond.neccon, NewAtom(lc, EQ))
			// a := qeopt.g.simplFof(atom, neccon, cond.sufcon)
			// if a == atom { // debug用後で削除. テスト中の問題ではここは通らない
			// 	panic(fmt.Sprintf("why!\natom=%v\nlc=%v\nneccon=%v\n", atom, lc, neccon))
			// }
		}

		switch qa := NewAtom(qq, atom.op).(type) {
		case *Atom:
			if qa.Deg(lv) >= deg {
				panic(fmt.Sprintf("atomQE: qa.Deg(lv)=%v, lv=%v, deg=%d; atom=%v", qa.Deg(lv), lv, deg, atom))
			}
			// lc != 0 で真のため， lc=0 は neccon への追加が不要
			// NOTE: lc=0 を追加すること atomQE 内での simplify に使えるか?
			rr := atomQE(qa, lv, neccon, qeopt)
			if rr == nil {
				return nil
			}
			return NewFmlOr(NewFmlAnd(NewAtom(lc, NE), neccon), rr)
		default:
			return NewFmlOr(NewFmlAnd(NewAtom(lc, NE), neccon), NewFmlAnd(qa, neccon))
		}
	} else if deg > ATOMQE_MAX_DEG || atom.op&EQ == 0 {
		return nil
	} else {
		////////////////////////////////
		// 偶数次
		////////////////////////////////
		aop := atom.op
		if atom.op == GE {
			p = p.Neg().(*Poly)
			aop = LE
		}

		if p.isEven(lv) { // 偶関数なら次数を落として SDCにできる.
			// ここの neccon は simplify 用ではなく，結果復帰用なので，
			// 限量変数である lv が含まれないことは保証されているはず.
			if neccon.hasVar(lv) {
				panic(fmt.Sprintf("atomQE: neccon=%v, lv=%v", neccon, lv))
			}

			q := p.redEven(lv)
			// f(x) = f(-x) = F(x^2)
			// ex(x, F(x^2) <= 0) => ex(y, F(y) <= 0 and y >= 0)
			return sdcQEcont2(q, aop, lv, neccon, qeopt)
		}

		// 公式を用いる
		afml := atomqe_formula[deg/2]
		ret := make([]Fof, 0, 4+len(afml)*2)

		nec2 := NewFmlAnd(NewAtom(lc, EQ), neccon)
		// 次数落ち
		switch atom2 := NewAtom(p.setZero(lv, deg), aop).(type) {
		case *Atom:
			ret = append(ret, atomQE(atom2, lv, nec2, qeopt))
		default:
			ret = append(ret, NewFmlAnds(nec2, atom2))
		}

		pd := p.Diff(lv).(*Poly)
		sres := qeopt.g.ox.Sres(p, pd, lv, 1)
		// fmt.Printf("sres=%v\n", sres)

		if atom.op == EQ {
			// 主係数と定数項の符号が異なれば ok
			ret = append(ret, NewFmlAnds(NewAtom(lc, LT), NewAtom(cont, GE), neccon))
			ret = append(ret, NewFmlAnds(NewAtom(lc, GT), NewAtom(cont, LE), neccon))
			nec2 = NewFmlAnd(NewAtom(lc, GT), neccon)
			ret = atomQEconstruct(ret, sres, afml, nec2)

			nec2 = NewFmlAnd(NewAtom(lc, LT), neccon)
			sresneg := NewListN(sres.Len())
			for i := 0; i < sres.Len(); i++ {
				u := sres.getiPoly(i)
				sresneg.Append(u.Neg())
			}
			ret = atomQEconstruct(ret, sresneg, afml, nec2)

		} else {
			// 主係数で成立確定
			ret = append(ret, NewFmlAnds(neccon, NewAtom(lc, atom.op & ^EQ)))

			nec2 = NewFmlAnd(neccon, NewAtom(lc, atom.op.not()))
			ret = atomQEconstruct(ret, sres, afml, nec2)
		}
		return NewFmlOrs(ret...)
	}
}

// 範囲条件が多いので reduction する
func sdcQEredRange(f *Atom, rminmax [][]*Atom, idx int, lv Level, neccon Fof, qeopt QEopt) Fof {

	// コメントは idx=0 のときを想定して記述している.
	rmin := rminmax[idx]
	ret := make([]Fof, 0, len(rmin)*5)

	a := make([]RObj, len(rmin))
	b := make([]RObj, len(rmin))
	for i, rr := range rmin {
		pi := rr.getPoly()
		a[i] = pi.Coef(lv, 1)
		if !a[i].IsNumeric() || a[i].Sign() <= 0 {
			panic(fmt.Sprintf("rmin[%d]=%v, a=%v\n", i, rr, a[i]))
		}
		b[i] = pi.Coef(lv, 0)
	}

	// 下限を表すものが複数あったから，一つになるまで絞り込む
	for i, rr := range rmin {
		op := rr.op
		ands := make([]Fof, 0, len(rmin)) // i 番目が一番大きいための条件

		for j := 0; j < len(rmin); j++ {
			if i != j {
				// a*x+b >= 0 && c*x + d >= 0  ==> a*x+b >= 0
				// x >= -b/a && x >= -d/c      ==> -b/a >= -d/c  ==> a*d-c*b >= 0
				switch bb := NewAtom(Sub(Mul(a[i], b[j]), Mul(a[j], b[i])), op).(type) {
				case *Atom:
					ands = append(ands, bb) // rmin[i] >= rmin[j]
				case *AtomF:
					// simplify されていたらありえないはずだが.
					ands = nil
					break
				}
			}
		}
		if ands == nil {
			continue
		}
		var ff Fof
		neccon := NewFmlAnds(ands...)
		if idx == 0 {
			ff = sdcQEcont(f, []*Atom{rr}, rminmax[1], lv, neccon, qeopt)
		} else {
			ff = sdcQEcont(f, rminmax[0], []*Atom{rr}, lv, neccon, qeopt)
		}
		if ff == nil {
			return nil
		}

		ret = append(ret, ff)
	}
	return NewFmlOrs(ret...)
}

// sdcQE. 変数を消去できなかったときは nil を返す
//
// ex([lv], f && rng)
// where rng.Deg(lv) = 1 and f.Deg(lv) > 1
func sdcQE(f *Atom, rngs []*Atom, lv Level, qeopt QEopt, cond qeCond) Fof {
	var rmin, rmax []*Atom

	var rng_lcpol *Atom

	if deg := f.Deg(lv); deg <= 0 || deg > SDCQE_MAX_DEG || f.op&EQ == 0 {
		return nil
	}

	for _, rng := range rngs {
		if rng.op == EQ || rng.op&EQ == 0 {
			return nil
		}
		rr := rng.getPoly()
		lc := rr.Coef(lv, 1)
		if !lc.IsNumeric() {
			if rng_lcpol != nil {
				return nil
			}
			rng_lcpol = rng
			continue
		}

		if rng.op == GE {
			rmin = append(rmin, rng)
		} else { // LE
			rmax = append(rmax, rng)
		}
	}
	if rng_lcpol != nil {
		if len(rmin)+len(rmax) > 0 {
			return nil // unsupported
		}

		// ex([x], f(x) <= 0 && a*x+b >= 0)
		return sdcQEpoly(f, rng_lcpol, lv, trueObj, qeopt)
	}

	// ex([x], f(x) <= 0 && land_i x >= ai && land_j x <= bj
	return sdcQEcont(f, rmin, rmax, lv, trueObj, qeopt)
}

// lv を主変数として，多項式を表示する
func sdcPolyFormat(f *Poly, lv Level) string {
	if f == nil {
		return "nil"
	}
	ret := ""
	sss := ""
	deg := f.Deg(lv)
	for i := deg; i >= 0; i-- {
		c := f.Coef(lv, uint(i))
		if c.IsZero() {
			continue
		}
		ret += fmt.Sprintf("%s(%v)", sss, c)
		sss = "+"
		if i >= 1 {
			ret += fmt.Sprintf("*%s", VarStr(lv))
			if i >= 2 {
				ret += fmt.Sprintf("^%d", i)
			}
		}
	}
	return ret
}

// ex([x], f(x) <= 0 && x >= 0) && neccon
// f.Coef(lv, 0) <= 0 の場合は考慮しない
func sdcQEmain(f *Poly, lv Level, neccon Fof, qeopt QEopt) Fof {
	print_log := false
	deg := f.Deg(lv)
	if deg > SDCQE_MAX_DEG {
		return nil
	} else if deg <= 1 {
		if deg == 0 {
			return falseObj
		}
		lc := f.Coef(lv, uint(deg))
		return NewFmlAnd(neccon, NewAtom(lc, LT))
	} else if deg%2 == 0 {
		for f.isEven(lv) {
			deg /= 2
			f = f.redEven(lv)
		}
	}

	if false {
		fmt.Printf("sdcQEmain(neccon = %v)\n", neccon)
		fmt.Printf("FF = %s\n", sdcPolyFormat(f, lv))
	}

	fd := f.Diff(lv).(*Poly)
	sres := qeopt.g.ox.Sres(f, fd, lv, 3)
	if print_log {
		fmt.Printf("sres=%v\n", sres)
	}
	hc := make([]RObj, 2*deg-2)

	// hc... h[n-2], h[n-3], ..., h[0], c[n-1], c[n-2], ..., c[1]
	for i := 0; i < deg-1; i++ {
		d := deg - i - 2
		var h, c RObj
		switch v := sres.geti(d).(type) {
		case *Poly:
			h = v.Coef(lv, uint(d))
			c = v.Coef(lv, uint(0))
		case RObj:
			if d > 0 {
				h = zero
			} else {
				h = v
			}
			c = v
		}

		hc[i] = h
		if d > 0 {
			hc[i+deg] = c
		} else {
			hc[deg-1] = fd.Coef(lv, 0)
		}
	}
	if print_log {
		fmt.Printf("hc=%v\n", hc)
	}

	N := len(sdcqe_formula[deg])
	ret := make([]Fof, N+1, N+2)
	lc := f.Coef(lv, uint(deg))
	lcpos := NewAtom(lc, GT)
	for i, sdcfml := range sdcqe_formula[deg] {
		// fmt.Printf("fml[%d]=%v\n", i, sdcfml)
		vv := make([]Fof, 2, len(hc)+2)
		vv[0] = neccon
		vv[1] = lcpos
		for j, v := range sdcfml {
			if v != OP_TRUE {
				vv = append(vv, NewAtom(hc[j], v))
			}
		}

		ret[i] = NewFmlAnds(vv...)
		if print_log {
			fmt.Printf("vv[%d/%d]=%v => %v\n", i, N, vv, ret[i])
		}
	}
	ret[N] = NewFmlAnds(neccon, NewAtom(lc, LT))
	if print_log {
		fmt.Printf("vV[%d/%d]=%v\n", N, N, ret[N])
	}

	if f2, ok := f.setZero(lv, deg).(*Poly); ok {
		ret = append(ret, sdcQEmain(f2, lv, NewFmlAnd(NewAtom(lc, EQ), neccon), qeopt))
		if print_log {
			fmt.Printf("VV[%d/%d]=%v\n", N+1, N, ret[N+1])
		}
	} else {
		// 定数は正であることが必要条件なので，ここは不要
	}

	return NewFmlOrs(ret...)
}

// ex([lv], f && a*x+b >= 0)
// rng の次数は 1, 主係数は変数も可能
func sdcQEpoly(f, rng *Atom, lv Level, neccon Fof, qeopt QEopt) Fof {
	if f.op&EQ == 0 {
		return nil
	}
	fp := f.getPoly()

	p := rng.getPoly()
	a := p.Coef(lv, 1)
	b := p.Coef(lv, 0)
	q := Mul(NewPolyVar(lv), a).Sub(b)

	// a*x+b >= 0 && fp(x) >= 0 or
	// a*x+b >= 0 && fp(x) == 0
	deg := fp.Deg(lv)
	den := makeDenAry(a, deg)
	fp = fp.subst_frac(q, den, lv).(*Poly)
	// fmt.Printf("  fp=%v %s 0\n", sdcPolyFormat(fp, lv), f.op)

	ret := make([]Fof, 1, 6)

	// rng の主係数 $a が 0 で範囲条件がなくなるケース
	{
		nec2 := NewFmlAnds(neccon, NewAtom(a, EQ), NewAtom(b, rng.op))
		ret[0] = atomQE(f, lv, nec2, qeopt)
		// fmt.Printf("  ret[0]=%v; %v\n", ret[0], nec2)
	}

	c0 := fp.Coef(lv, 0)

	// 定数項が正(GT)の場合. LE で真が確定なので，必要条件に加える必要はない
	//       c0 <= 0 || c0 > 0 && SDC
	// <==>  c0 <= 0 || SDC

	// a > 0 の場合
	{
		// fmt.Printf("@@@ case a+ := %v>0; %v\n", a, rng)
		nec2 := neccon
		nec2 = NewFmlAnds(neccon, NewAtom(a, GT))
		cop := NewAtom(c0, f.op)
		ret = append(ret, NewFmlAnd(cop, nec2))
		// fmt.Printf("  retA1[%d]=%v; %v %s 0\n", len(ret), ret[len(ret)-1], c0, f.op)

		_, ok2 := cop.(*AtomT)
		if _, ok := nec2.(*AtomF); !ok && !ok2 {
			fq := fp
			opf := f.op
			if opf&GT != 0 {
				fq = fq.Neg().(*Poly)
				opf = opf.neg()
			}
			if rng.op == LE {
				fq = fq.NegX(lv)
			}

			rr := sdcQEmain(fq, lv, nec2, qeopt)
			ret = append(ret, rr)
			// fmt.Printf("  retA2[%d]=%v\n", len(ret), ret[len(ret)-1])
		}
	}

	{ // a < 0 の場合
		// fmt.Printf("@@@ case a- := %v<0; %v\n", a, rng)
		nec2 := neccon
		nec2 = NewFmlAnds(neccon, NewAtom(a, LT))
		opf := f.op
		opr := rng.op.neg() // 上下限が逆になる
		if deg%2 != 0 {
			opf = opf.neg()
		}
		cop := NewAtom(c0, opf)
		ret = append(ret, NewFmlAnd(cop, nec2))
		// fmt.Printf("  retB1[%d]=%v\n", len(ret), ret[len(ret)-1])
		_, ok2 := cop.(*AtomT)
		if _, ok := nec2.(*AtomF); !ok && !ok2 {
			fq := fp
			if opf&GT != 0 {
				fq = fq.Neg().(*Poly)
				opf = opf.neg()
			}
			if opr == LE {
				fq = fq.NegX(lv)
			}

			rr := sdcQEmain(fq, lv, nec2, qeopt)
			ret = append(ret, rr)
			// fmt.Printf("  retB2[%d]=%v\n", len(ret), ret[len(ret)-1])
		}
	}

	return NewFmlOrs(ret...)
}

// 区間範囲をずらす
// f(x) <= 0 && p(x) = x-5 >= 0
// =>
// f(x+5) <= 0 && x >= 0
//
// deg(p) = 1
func (f *Poly) sdcSubst(p *Poly, lv Level) RObj {
	if p.Deg(lv) != 1 {
		panic(fmt.Sprintf("sdcSubst: p.Deg(%s) != 1: %v", VarStr(lv), p))
	}
	if c := p.Coef(lv, 1); !c.IsNumeric() || c.Sign() <= 0 {
		panic(fmt.Sprintf("sdcSubst: p.Coef(%s, 1) is not a positive number: %v", VarStr(lv), p))
	}

	c1 := p.Coef(lv, 1)
	c0 := p.Coef(lv, 0)
	q := NewPolyVar(lv).Mul(c1).Sub(c0)
	if c1.IsOne() {
		g := f.Subst(q, lv).(*Poly)
		// fmt.Printf("sdcSubst(1) f=%v\np=%v\nq=%v\n", f, p, q)
		// fmt.Printf("g=%v\n", sdcPolyFormat(g, lv))
		return g
	}

	deg := f.Deg(lv)
	den := makeDenAry(c1, deg)
	g := f.subst_frac(q, den, lv)
	// fmt.Printf("sdcSubst(2) f=%v\np=%v\ng=%v\n", f, p, g)
	return g
}

// ex([lv], f && land_i x >= a_i && land_j x <= b_j)
func sdcQEcont(f *Atom, rmin, rmax []*Atom, lv Level, neccon Fof, qeopt QEopt) Fof {
	// fmt.Printf("\n\n@@@@@@@@\nsdcQEcont(f=%v, rmin=%v, rmax=%v, lv=%s, neccon=%v, qeopt=%v)\n", f, rmin, rmax, VarStr(lv), neccon, qeopt)
	if len(rmin) > 1 {
		// 下限を表すものが複数あったから，一つになるまで絞り込む
		return sdcQEredRange(f, [][]*Atom{rmin, rmax}, 0, lv, neccon, qeopt)
	}
	if len(rmax) > 1 {
		// 上限を表すものが複数あったから，一つになるまで絞り込む
		return sdcQEredRange(f, [][]*Atom{rmin, rmax}, 1, lv, neccon, qeopt)
	}

	fp := f.getPoly()
	opf := f.op
	if f.op == GE {
		fp = fp.Neg().(*Poly)
		opf = LE
	}
	if len(rmax) == 0 || len(rmin) == 0 {
		// 上限か下限どちらか一方しかない場合
		var r *Poly
		if len(rmin) != 0 {
			r = rmin[0].getPoly()
			if rmin[0].op == LE {
				r = r.Neg().(*Poly)
			}
		} else if len(rmax) != 0 {
			r = rmax[0].getPoly()
			if rmax[0].op == GE {
				r = r.Neg().(*Poly)
			}
		} else {
			// atom だった
			return atomQE(f, lv, neccon, qeopt)
		}
		fp = fp.sdcSubst(r, lv).(*Poly)
		if len(rmax) > 0 {
			// x<=0 を x>=0 に変換
			fp = fp.NegX(lv)
		}

		return sdcQEcont2(fp, opf, lv, neccon, qeopt)
	}

	// 上下限ある場合
	// a <= x <= b  :: -a0/a1 <= x <= -b0/b1
	deg := fp.Deg(lv)
	ap := rmin[0].getPoly()
	bp := rmax[0].getPoly()
	if rmin[0].op == LE {
		ap = ap.Neg().(*Poly)
	}
	if rmax[0].op == GE {
		bp = bp.Neg().(*Poly)
	}

	a1 := ap.Coef(lv, 1)
	a0 := ap.Coef(lv, 0)
	b1 := bp.Coef(lv, 1)
	b0 := bp.Coef(lv, 0)
	if !a1.IsNumeric() || a1.Sign() <= 0 || !b1.IsNumeric() || b1.Sign() <= 0 {
		// 1次の係数は定数であり，正と仮定
		panic(fmt.Sprintf("invalid A=%v, B=%v", rmin[0], rmax[0]))
	}
	dist_num := Sub(Mul(b1, a0), Mul(b0, a1)) // bug..... 分数のまま処理しないと，distとして正しくない

	//////////////////////////////////////
	// 端点で成立するところを取り除き
	//////////////////////////////////////
	cap_ret := 7
	ret := make([]Fof, 2, cap_ret)
	nec2 := NewFmlAnd(neccon, NewAtom(dist_num, GE))
	// fmt.Printf("nec=%v, dist_num=%v\n", neccon, dist_num)

	dena := makeDenAry(a1, deg)
	fpa := fp.subst_frac(a0.Neg(), dena, lv)
	ret[0] = NewFmlAnd(nec2, NewAtom(fpa, opf)) // 左端点
	// fmt.Printf("ret[0]=%v\n", ret[0])

	if dist_num.IsNumeric() && dist_num.Sign() <= 0 {
		if !dist_num.IsZero() {
			return falseObj
		}
		return ret[0]
	}

	denb := makeDenAry(b1, deg)
	fpb := fp.subst_frac(b0.Neg(), denb, lv)
	ret[1] = NewFmlAnd(nec2, NewAtom(fpb, opf)) // 右端点
	// fmt.Printf("ret[1]=%v\n", ret[1])

	//////////////////////////////////////
	// 区間を x >= 0 に変換
	//////////////////////////////////////
	// fmt.Printf("fp[a..b  ] = %v\n", fp)
	fp = fp.sdcSubst(ap, lv).(*Poly) // 0 <= x <= b-a
	// fmt.Printf("fp[0..b-a] = %v @ %v\n", fp, ap)

	dist_den := makeDenAry(Mul(a1, b1), deg)
	fp = fp.subst_frac(Mul(NewPolyVar(lv), dist_num), dist_den, lv).(*Poly) // 0 <= x <= 1
	// fmt.Printf("fp[0..1  ] = %v @ %v\n", fp, dist_num)
	fp = fp.SubstXinvLv(lv, deg).(*Poly) // 1 <= x <= +inf
	// fmt.Printf("fp[1..inf] = %v\n", fp)
	x1 := NewPoly(lv, 2)
	x1.c[0] = one
	x1.c[1] = one
	fp = fp.Subst(x1, lv).(*Poly) // 0 <= x <= +inf
	// fmt.Printf("fp[0..inf] = %v\n", fp)

	neccon = NewFmlAnds(neccon, NewAtom(dist_num, GT))

	//////////////////////////////////////
	// 締め
	//////////////////////////////////////

	if f.op != EQ {
		rr := sdcQEmain(fp, lv, neccon, qeopt)
		ret = append(ret, rr)
	} else {
		fpagt := NewAtom(fpa, GT)
		fpalt := NewAtom(fpa, LT)
		fpbgt := NewAtom(fpb, GT)
		fpblt := NewAtom(fpb, LT)
		nec2 := NewFmlAnds(neccon, fpagt, fpbgt)
		ret = append(ret, sdcQEmain(fp, lv, nec2, qeopt))
		nec3 := NewFmlAnds(neccon, fpalt, fpblt)
		ret = append(ret, sdcQEmain(fp.Neg().(*Poly), lv, nec3, qeopt))
		ret = append(ret, NewFmlAnds(neccon, fpagt, fpblt))
		ret = append(ret, NewFmlAnds(neccon, fpalt, fpbgt))
	}
	if len(ret) > cap_ret {
		panic(fmt.Sprintf("#ret=%d, cap=%d\n", len(ret), cap_ret))
	}
	return NewFmlOrs(ret...)
}

func sdcQEcont2(fp *Poly, opf OP, lv Level, neccon Fof, qeopt QEopt) Fof {
	c0 := fp.Coef(lv, 0)
	if opf != EQ {
		if opf != LE {
			panic(fmt.Sprintf("why? op=%s[%x]", opf, opf))
		}
		ret := sdcQEmain(fp, lv, neccon, qeopt)
		return NewFmlOr(ret, NewFmlAnd(NewAtom(c0, LE), neccon))
	} else {
		deg := fp.Deg(lv)
		ret := make([]Fof, 3, 3+2*deg)
		c0gt := NewAtom(c0, GT)
		c0lt := NewAtom(c0, LT)
		ret[0] = sdcQEmain(fp, lv, c0gt, qeopt)
		ret[1] = sdcQEmain(fp.Neg().(*Poly), lv, c0lt, qeopt)
		ret[2] = NewFmlAnd(neccon, NewAtom(c0, EQ))

		// 主係数と定数項が異符号
		for i := uint(deg); i > 0; i-- {
			// @TODO simplify
			ci := fp.Coef(lv, i)
			ret = append(ret, NewFmlAnds(neccon, NewAtom(ci, GT), c0lt))
			ret = append(ret, NewFmlAnds(neccon, NewAtom(ci, LT), c0gt))
			neccon = NewFmlAnd(neccon, NewAtom(ci, EQ))
		}

		return NewFmlOrs(ret...)
	}
}

// fof: prenex first-order formula ex([x, y, z], p1 & p2 & ... & p2)
func sdcAtomQE(fof Fof, lv Level, qeopt QEopt, cond qeCond) Fof {
	// fmt.Printf("sdcAtomQE[%v] %v\n", VarStr(lv), fof)
	switch fofx := fof.(type) {
	case *Atom:
		rr := atomQE(fofx, lv, trueObj, qeopt)
		if rr == nil {
			return fof
		} else {
			return rr
		}
	case *FmlAnd:
		CAPRNG := 5
		var f *Atom                     // 2次以上の
		rng := make([]*Atom, 0, CAPRNG) // 1次の．範囲部分
		var other []Fof

		other = make([]Fof, 0, len(fofx.fml)+1)
		for _, p := range fofx.fml {
			if !p.hasVar(lv) {
				other = append(other, p)
				continue
			}
			if atom, ok := p.(*Atom); !ok {
				// atom でないものがあったら対応できない
				return fof
			} else if atom.Deg(lv) > 1 {
				if f != nil {
					// 2次以上のものが複数あったら対応できない
					return fof
				}
				f = atom
			} else {
				if len(rng) > CAPRNG { // 範囲条件がたくさんすぎた
					return fof
				}
				rng = append(rng, atom)
			}
		}
		if f == nil {
			return fof
		} else if len(rng) == 0 {
			rr := atomQE(f, lv, trueObj, qeopt)
			if rr == nil {
				return fof
			}
			other = append(other, rr)
			return NewFmlAnds(other...)
		} else {
			rr := sdcQE(f, rng, lv, qeopt, cond)
			if rr == nil {
				return fof
			}
			other = append(other, rr)
			return NewFmlAnds(other...)
		}
	default:
		panic(fmt.Sprintf("unexpected %v", fof))
	}
}

func (qeopt QEopt) qe_sdcatom(fof FofQ, cond qeCond) Fof {
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
		ff := sdcAtomQE(fml, q, qeopt, cond)
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
