package ganrac

// A Generalized Framework for Virtual Substitution
// M. Kosta, T. Sturm 2015

// 3次等式制約
// f(x) = a3 x^3 + a2 x^2 + a1 x + a0
//
// P = f(x) == 0 && Q
//
//    f' | f'' | f'''| index
// +-----+-----+-----+------
//   >=  |  +  |  +  |  3
//    -  |  *  |  +  |  2
//   >=  |  0  |  +  |  2
//   >=  |  -  |  +  |  1
//
// これが実用的でないなら，VS-cube が使えるときはない気がする.
// free-var が多いときに限定するか
import (
	"fmt"
)

type vs3_target struct {
	lv      Level
	p       *Poly
	a, b, c RObj
	op      OP
}

func (minatom *quadEQ_t) log(g *Ganrac, deg, idx int, op OP, msg string, f Fof) {
	dic := NewDict()
	dic.Set("var", NewInt(1))
	c, _ := funcCAD(g, "cad", []interface{}{f, dic})
	fmt.Printf("       vs3-%d-%d-%2v: %s: %v   <==   %v\n", deg, idx, op, msg, c, f)
}

func (minatom *quadEQ_t) appendvs3(g *Ganrac, deg, idx int, op OP, msg string, f Fof, ret []Fof) []Fof {
	fmt.Printf("appendvs3: %v\n", f)
	ret = append(ret, f)
	if true {
		minatom.log(g, deg, idx, op, msg, f)
	}
	if true {
		minatom.testCompare(g, f, "appendvs3")
	}
	return ret
}

type fn_vs3 func(v vs3_target, neccon Fof, ret []Fof, qeopt QEopt, cond qeCond) []Fof

/////////////////////////////////////////////////////////////////////////////////
// atom に X を仮想代入する
/////////////////////////////////////////////////////////////////////////////////

// ==============================================================
// 3次等式制約が虚根をもち，実根を唯一持つ場合
//
//	discrim < 0
//
// ==============================================================
func (minatom *quadEQ_t) vscube_x_quad(v vs3_target, neccon Fof, ret []Fof, qeopt QEopt, cond qeCond) []Fof {
	lv := v.lv
	a, b, c := v.a, v.b, v.c
	op := v.op

	spp := newVsQuadSamplePoint(3, a, b, c)
	spn := spp.newVsQuadSamplePointNeg()

	d2 := spp.sqr
	dge := NewAtom(d2, GE) // p が実根をもつ
	dlt := NewAtom(d2, LT) // p が実根をもたない

	switch op {
	case EQ:
		ret = append(ret, NewFmlAnds(neccon, dge, spp.virtual_subst_sqr(minatom.a, lv)))
		ret = append(ret, NewFmlAnds(neccon, dge, spn.virtual_subst_sqr(minatom.a, lv)))
	case NE:
		ret = append(ret, NewFmlAnds(neccon, dlt))
		ane := NewAtom(minatom.p, NE).(*Atom)
		sppane := spp.virtual_subst_sqr(ane, lv)
		spnane := spn.virtual_subst_sqr(ane, lv)
		ret = append(ret, NewFmlAnds(neccon, sppane, spnane))
	case LT, LE, GT, GE:
		// 変数名は LT の場合を基本としてつけているので，
		// GT, GE の場合は，逆 (lt -> gt, le -> ge) になる
		plt := minatom.atoms[LT|(op&EQ)]
		pgt := minatom.atoms[GT|(op&EQ)]

		var spn_plt, spp_pgt Fof

		var alt, agt Fof
		if op == LT || op == LE {
			alt = NewAtom(a, LT)
			agt = NewAtom(a, GT)
			spn_plt = spn.virtual_subst_sqr(plt, lv)
			spp_pgt = spp.virtual_subst_sqr(pgt, lv)
		} else {
			alt = NewAtom(a, GT)
			agt = NewAtom(a, LT)
			spn_plt = spp.virtual_subst_sqr(plt, lv)
			spp_pgt = spn.virtual_subst_sqr(pgt, lv)
		}

		// a > 0 (spn のほうが小さい根)
		ret = append(ret, NewFmlAnds(agt, dge, spn_plt, spp_pgt))

		// a < 0 (spp のほうが小さい根)
		ret = append(ret, NewFmlAnds(alt, dlt))
		ret = append(ret, NewFmlAnds(alt, spn_plt))
		ret = append(ret, NewFmlAnds(alt, spp_pgt))

	default:
		panic(fmt.Sprintf("not implemented: %v", op))
	}

	return ret
}

func (minatom *quadEQ_t) vscube_x_lin(v vs3_target, neccon Fof, ret []Fof, qeopt QEopt, cond qeCond) []Fof {
	lv := v.lv
	_, b, c := v.a, v.b, v.c
	op := v.op

	dens := makeDenAry(b, minatom.deg)
	fbeta := minatom.p.subst_frac(c.Neg(), dens, lv)
	// fmt.Printf("%v [%d // -(%v)/%v] => %v\n", minatom.p, lv, c, b, fbeta)
	ret = append(ret, NewFmlAnds(neccon, NewAtom(b, NE), NewAtom(fbeta, op.neg())))
	return ret
}

//==============================================================
// 等式制約が虚根をもたず，実根を重複を込めて３つもつ場合
//   discrim >= 0
//==============================================================

// 重複を含めて 3 つの実根を持つ場合,
// i.e., 3次等式制約の判別式が非負の場合
// かつ，１つ目の実根

/**
 * 3次等式制約 minatom.p は主係数が正で，虚根を持たないと仮定
 * 2次のatom (a*x^2+b*x+c op 0) に対して，
 * minatom.p の最小の実根を VS する
 */
func (minatom *quadEQ_t) vscube_1_quad(v vs3_target, neccon Fof, ret []Fof, qeopt QEopt, cond qeCond) []Fof {
	fmt.Printf("vscube_1_quad: (%v)*x^2+ (%v)*x+ (%v) %v 0\n", v.a, v.b, v.c, v.op)
	return ret
}

/**
 *
 * minatom: 3次等式制約
 * @param v 代入先 1次の多項式
 */
func (minatom *quadEQ_t) vscube_1_lin(v vs3_target, neccon Fof, ret []Fof, qeopt QEopt, cond qeCond) []Fof {
	lv := v.lv
	_, b, c := v.a, v.b, v.c
	op := v.op

	fmt.Printf("vscube_1_lin: (%v)*x+ (%v) %v 0\n", b, c, op)
	z := newVsLinSamplePoint(3, c, b)

	switch op { // P(1, 1, op) on vs3.md
	case EQ: // P(1, 1, =)
		g := NewAtoms([]RObj{v.p, b}, GE).(*Atom)
		gy1 := minatom.diff[1].virtual_subst(g, lv)
		hzeq := z.virtual_subst(NewAtom(minatom.p, EQ).(*Atom), lv)
		ret = append(ret, NewFmlAnds(neccon, hzeq, gy1))
		return ret
	case NE: // P(1, 1, !=)
		g := NewAtoms([]RObj{v.p, b}, LT).(*Atom)
		gy1 := minatom.diff[1].virtual_subst(g, lv)
		hzne := z.virtual_subst(NewAtom(minatom.p, NE).(*Atom), lv)
		ret = append(ret, NewFmlAnd(neccon, NewFmlOr(hzne, gy1)))
		return ret
	}

	densgn := z.DenSign()
	var hzgt Fof = falseObj
	var hzlt Fof = falseObj

	bgt := NewAtom(b, GT)
	blt := NewAtom(b, LT)

	switch op { // P(1, 1, op)
	case GT, GE:
		g := NewAtom(v.p, GE).(*Atom)
		// gy1 := g(x // y1)
		gy1 := virtual_subst(g, minatom.diff[1], lv)

		fmt.Printf("vs3-1: %v, densgn=%v\n", op, densgn)
		fmt.Printf("diff(y)-= %v\n", minatom.diff[1])
		fmt.Printf("g       = %v\n", g)
		fmt.Printf("gy1     = %v\n", gy1)

		if densgn >= 0 {
			z.SetDenSign(1)
			hzlt = z.virtual_subst(NewAtom(minatom.p, LT|(op&EQ)).(*Atom), lv)
			ret = minatom.appendvs3(qeopt.g, 1, 1, op, "b>0 & x in R1", NewFmlAnds(neccon, bgt, hzlt, gy1), ret)

		}
		if densgn <= 0 {
			z.SetDenSign(-1)

			fmt.Printf("b=%v\n", b)

			ret = minatom.appendvs3(qeopt.g, 1, 1, op, "b<0 & x >= y1", NewFmlAnds(neccon, blt, gy1), ret)

			fmt.Printf("z=%v\n", z)
			hzgt = z.virtual_subst(NewAtom(minatom.p, GT|(op&EQ)).(*Atom), lv)
			ret = minatom.appendvs3(qeopt.g, 1, 1, op, "b<0 & hz > 0 ", NewFmlAnds(neccon, blt, hzgt), ret)
		}
		fmt.Printf("h =%v\n", minatom.p)
		fmt.Printf("y1=%v\n", minatom.diff[1])
		fmt.Printf("hzgt=%v\n", hzgt)
	case LT, LE:
		g := NewAtom(v.p, LE).(*Atom)
		gy1 := minatom.diff[1].virtual_subst(g, lv)
		if densgn >= 0 {
			z.SetDenSign(1)
			hzgt = z.virtual_subst(NewAtom(minatom.p, GT|(op&EQ)).(*Atom), lv)
			ret = append(ret, NewFmlAnds(neccon, bgt, hzgt))
		}
		if densgn <= 0 {
			z.SetDenSign(-1)
			hzlt = z.virtual_subst(NewAtom(minatom.p, LT|(op&EQ)).(*Atom), lv)
			ret = append(ret, NewFmlAnds(neccon, blt, hzlt, gy1))
		}
		ret = append(ret, NewFmlAnds(neccon, bgt, gy1))
	}

	return ret
}
func (minatom *quadEQ_t) vscube_2_quad(v vs3_target, neccon Fof, ret []Fof, qeopt QEopt, cond qeCond) []Fof {
	return ret
}
func (minatom *quadEQ_t) vscube_2_lin(v vs3_target, neccon Fof, ret []Fof, qeopt QEopt, cond qeCond) []Fof {
	deg := 1
	idx := 2
	lv := v.lv
	_, b, c := v.a, v.b, v.c
	op := v.op

	z := newVsLinSamplePoint(3, c, b)

	switch op { // P(1, 1, op) on vs3.md
	case EQ: // P(1, 1, =)
		g1 := NewAtoms([]RObj{v.p, b}, LE).(*Atom)
		gy1 := minatom.diff[1].virtual_subst(g1, lv)
		g2 := NewAtoms([]RObj{v.p, b}, GE).(*Atom)
		gy2 := minatom.diff[2].virtual_subst(g2, lv)
		hz := z.virtual_subst(NewAtom(minatom.p, EQ).(*Atom), lv)
		ret = minatom.appendvs3(qeopt.g, deg, idx, op, "hz=0 +a", NewFmlAnds(neccon, hz, gy1, gy2), ret)
		return ret
	case NE:
		g1 := NewAtoms([]RObj{v.p, b}, GT).(*Atom)
		gy1 := minatom.diff[1].virtual_subst(g1, lv)
		g2 := NewAtoms([]RObj{v.p, b}, LT).(*Atom)
		gy2 := minatom.diff[2].virtual_subst(g2, lv)
		hz := z.virtual_subst(NewAtom(minatom.p, NE).(*Atom), lv)
		ret = minatom.appendvs3(qeopt.g, deg, idx, op, "hz!=0 +a", NewFmlAnds(neccon, NewFmlOrs(hz, gy1, gy2)), ret)
		return ret
	}

	densgn := z.DenSign()
	fmt.Printf("z::: b=%v, c=%v, densgn=%d\n", b, c, densgn)

	bgt := NewAtom(b, GT)
	blt := NewAtom(b, LT)

	fmt.Printf("g=%v\n", v.p)
	fmt.Printf("diff1=%v\n", minatom.diff[1])
	g1 := NewAtom(v.p, op|EQ).(*Atom)
	gy1 := minatom.diff[1].virtual_subst(g1, lv)
	g2 := NewAtom(v.p, op|EQ).(*Atom)
	gy2 := minatom.diff[2].virtual_subst(g2, lv)

	fmt.Printf("@@ g1=%v\n", gy1)
	fmt.Printf("@@ g2=%v\n", gy2)

	switch op { // P(1, 1, op)
	case GT, GE:

		if densgn >= 0 {
			z.SetDenSign(+1)
			hzgt := z.virtual_subst(NewAtom(minatom.p, LT|(op&EQ)).(*Atom), lv)
			fmt.Printf("h=%v\n", minatom.p)
			fmt.Printf("z=%v\n", z)
			fmt.Printf("h(z)>=0: %v\n", qeopt.g.simplFof(hzgt, cond.neccon, cond.sufcon))
			ret = minatom.appendvs3(qeopt.g, deg, idx, op, "b>0 & hz > 0+", NewFmlAnds(neccon, bgt, hzgt, gy2), ret)

			ret = minatom.appendvs3(qeopt.g, deg, idx, op, "b>0 & x <= y1", NewFmlAnds(neccon, bgt, gy1), ret)
		}

		if densgn <= 0 {
			z.SetDenSign(-1)
			hzlt := z.virtual_subst(NewAtom(minatom.p, LT|(op&EQ)).(*Atom), lv)
			ret = minatom.appendvs3(qeopt.g, deg, idx, op, "b<0 & hz < 0+", NewFmlAnds(neccon, blt, hzlt, gy1), ret)

			ret = minatom.appendvs3(qeopt.g, deg, idx, op, "b<0 & x >= y2", NewFmlAnds(neccon, blt, gy2), ret)
		}

	case LT, LE:
		if densgn >= 0 {
			z.SetDenSign(+1)
			hzlt := z.virtual_subst(NewAtom(minatom.p, LT|(op&EQ)).(*Atom), lv)
			ret = append(ret, NewFmlAnds(neccon, bgt, hzlt, gy1))
		}
		if densgn <= 0 {
			z.SetDenSign(-1)
			hzgt := z.virtual_subst(NewAtom(minatom.p, GT|(op&EQ)).(*Atom), lv)
			ret = append(ret, NewFmlAnds(neccon, blt, hzgt, gy2))
		}

		ret = append(ret, NewFmlAnds(neccon, bgt, gy2))
		ret = append(ret, NewFmlAnds(neccon, blt, gy1))
	}

	return ret
}
func (minatom *quadEQ_t) vscube_3_quad(v vs3_target, neccon Fof, ret []Fof, qeopt QEopt, cond qeCond) []Fof {
	return ret
}
func (minatom *quadEQ_t) vscube_3_lin(v vs3_target, neccon Fof, ret []Fof, qeopt QEopt, cond qeCond) []Fof {
	deg := 1
	idx := 3
	lv := v.lv
	_, b, c := v.a, v.b, v.c
	op := v.op

	z := newVsLinSamplePoint(3, c, b)

	switch op { // P(1, 1, op) on vs3.md
	case EQ: // P(1, 1, =)
		g := NewAtoms([]RObj{v.p, b}, LE).(*Atom)
		gy2 := minatom.diff[2].virtual_subst(g, lv)
		hz := z.virtual_subst(NewAtom(minatom.p, EQ).(*Atom), lv)
		ret = append(ret, NewFmlAnds(neccon, hz, gy2))
		return ret
	case NE:
		g := NewAtoms([]RObj{v.p, b}, GT).(*Atom)
		gy1 := minatom.diff[1].virtual_subst(g, lv)
		hz := z.virtual_subst(NewAtom(minatom.p, NE).(*Atom), lv)
		ret = append(ret, NewFmlAnd(neccon, NewFmlOr(hz, gy1)))
		return ret
	}

	densgn := z.DenSign()

	bgt := NewAtom(b, GT)
	blt := NewAtom(b, LT)

	switch op { // P(1, 1, op)
	case GT, GE:
		g := NewAtom(v.p, GE).(*Atom)
		gy2 := minatom.diff[2].virtual_subst(g, lv)
		if densgn >= 0 {
			z.SetDenSign(1)
			hzlt := z.virtual_subst(NewAtom(minatom.p, LT|(op&EQ)).(*Atom), lv)
			ret = minatom.appendvs3(qeopt.g, deg, idx, op, "b>0 & hz < 0 ",
				NewFmlAnds(neccon, bgt, hzlt), ret)
		}
		if densgn <= 0 {
			z.SetDenSign(-1)
			hzgt := z.virtual_subst(NewAtom(minatom.p, GT|(op&EQ)).(*Atom), lv)
			ret = minatom.appendvs3(qeopt.g, deg, idx, op, "b<0 & hz > 0+",
				NewFmlAnds(neccon, blt, hzgt, gy2), ret)
		}
		ret = minatom.appendvs3(qeopt.g, deg, idx, op, "b<0 & x >= y2",
			NewFmlAnds(neccon, bgt, gy2), ret)
		fmt.Printf("y2=%v\n", minatom.diff[2])
		fmt.Printf("b=%v\n", b)
	case LT, LE:
		g := NewAtom(v.p, LE).(*Atom)
		gy2 := minatom.diff[2].virtual_subst(g, lv)
		if densgn >= 0 {
			z.SetDenSign(1)
			hzgt := z.virtual_subst(NewAtom(minatom.p, GT|(op&EQ)).(*Atom), lv)
			ret = append(ret, NewFmlAnds(neccon, bgt, hzgt, gy2))
		}
		if densgn <= 0 {
			z.SetDenSign(-1)
			hzlt := z.virtual_subst(NewAtom(minatom.p, LT|(op&EQ)).(*Atom), lv)
			ret = append(ret, NewFmlAnds(neccon, blt, hzlt))
		}
		ret = append(ret, NewFmlAnds(neccon, blt, gy2))
	}
	return ret
}

/////////////////////////////////////////////////////////////////////////////////
// FOF に X を仮想代入する (汎用)
/////////////////////////////////////////////////////////////////////////////////

/**
 * atom の minatom.lv に対して VS する
 *
 * atom.Deg(lv) < 3 を仮定 (@see prem_qff())
 */
func (minatom *quadEQ_t) vscube_z_atom(atom *Atom, qeopt QEopt, cond qeCond,
	fn_quad, fn_lin fn_vs3) Fof {
	p := atom.getPoly()

	vv := vs3_target{
		minatom.lv,
		p,
		p.Coef(minatom.lv, 2),
		p.Coef(minatom.lv, 1),
		p.Coef(minatom.lv, 0),
		atom.op}

	d := p.Deg(vv.lv)
	N := 7
	ret := make([]Fof, 0, N)
	var neccon Fof = trueObj
	if d == 2 {
		vv.a = p.Coef(vv.lv, 2)
		ret = fn_quad(vv, NewAtom(vv.a, NE), ret, qeopt, cond)
		neccon = NewAtom(vv.a, EQ)
	}

	if _, ok := neccon.(*AtomF); !ok {
		// 1 次の場合
		// a == 0
		// -c/b を minatom.z に代入する
		neccon1 := NewFmlAnds(neccon, NewAtom(vv.b, NE))
		neccon1 = qeopt.g.simplFof(neccon1, cond.neccon, cond.sufcon)
		if _, ok := neccon1.(*AtomF); !ok {
			ret = fn_lin(vv, neccon1, ret, qeopt, cond)
		}
		ret = append(ret, NewFmlAnds(neccon, NewAtom(vv.b, EQ), NewAtom(vv.c, atom.op)))
	}

	for i, r := range ret {
		fmt.Printf("vscube_z_atom ret[%d]=%v\n", i, r)
	}

	if len(ret) >= N {
		panic(fmt.Sprintf("too small N! N=%d, len(ret)=%d", N, len(ret)))
	}

	return NewFmlOrs(ret...)
}

// 複素根を 2 つ持つ場合 (実根を 1 つのみもつ）
// i.e., 三次等式制約の判別式が負の場合
func (minatom *quadEQ_t) vscube_z(fml Fof, qeopt QEopt, cond qeCond,
	fn_quad, fn_lin fn_vs3) Fof {

	fmt.Printf("vscube_z(%v)\n", fml)
	switch fff := fml.(type) {
	case FofAO:
		same := true
		fmls := fff.Fmls()
		ret := make([]Fof, len(fmls))
		for i := 0; i < len(fmls); i++ {
			ret[i] = minatom.vscube_z(fmls[i], qeopt, cond, fn_quad, fn_lin)
			if ret[i] != fmls[i] {
				same = false
			}
		}
		if same {
			return fml
		}
		return fff.gen(ret)
	case *Atom:
		p := fff.getPoly()
		d := p.Deg(minatom.lv)
		if d <= 0 {
			return fml
		} else if d > 2 {
			panic(fmt.Sprintf("why?: %v", p))
		}
		return minatom.vscube_z_atom(fff, qeopt, cond, fn_quad, fn_lin)
	case *AtomF:
		return falseObj
	case *AtomT:
		return trueObj
	}
	panic(fmt.Sprintf("why! %v", fml))
}

/////////////////////////////////////////////////////////////////////////////////
// 準備
/////////////////////////////////////////////////////////////////////////////////

/**
 * 符号を考慮して擬剰余
 * minatom.z (minatom.p の lv に関する主係数 $lc）は非ゼロを仮定
 *
 * returns (fpos, fneg)
 *    fpos .. 主係数 $lc の符号が正のとき
 *    fneg .. 主係数 $lc の符号が負のとき
 */
func (minatom *quadEQ_t) prem_atom(atom *Atom) (Fof, Fof) {

	lv := minatom.lv
	f := atom.getPoly()
	var ff RObj = f
	org_deg := atom.Deg(lv) + 1
	m := 0
	for m < org_deg+5 {
		d := f.Deg(lv)
		if d < minatom.deg {
			break
		}

		m++
		fm := f.Mul(minatom.z)
		lc := f.Coef(lv, uint(d))
		gm := minatom.p.Mul(lc)
		if d > minatom.deg {
			xm := NewPolyVarn(lv, d-minatom.deg)
			gm = gm.Mul(xm)
		}

		ff = fm.Sub(gm)
		if ff_poly, ok := ff.(*Poly); ok {
			f = ff_poly
		} else {
			break
		}
	}
	if m >= org_deg {
		panic(fmt.Sprintf("stop!: atom=%v, minatom=%v", atom, minatom))
	}

	if m == 0 {
		return atom, atom
	} else if m%2 == 0 {
		a := NewAtom(ff, atom.op)
		return a, a
	} else {
		return NewAtom(ff, atom.op), NewAtom(ff, atom.op.neg())
	}
}

/**
 * fff を簡単化して，fff.Deg(lv) < 3 にする
 * minatom.z (minatom.p の lv に関する主係数 $lc）は非ゼロを仮定
 *
 * returns (fpos, fneg)
 *    fpos .. 主係数 $lc の符号が正のとき
 *    fneg .. 主係数 $lc の符号が負のとき
 */
func (minatom *quadEQ_t) prem_qff(fff Fof) (Fof, Fof) {
	switch f := fff.(type) {
	case FofAO:
		fmls := f.Fmls()
		pos := make([]Fof, len(fmls))
		neg := make([]Fof, len(fmls))
		same := true
		for i := 0; i < len(fmls); i++ {
			pos[i], neg[i] = minatom.prem_qff(fmls[i])
			if pos[i] != neg[i] {
				same = false
			}
		}
		fp := f.gen(pos)
		if same {
			return fp, fp
		}
		return fp, f.gen(neg)
	case *Atom:
		return minatom.prem_atom(f)
	default:
		panic("why!")
	}
}

var _cubeeq_expect Fof = nil

func CubeeqSetExpect(fml Fof) {
	_cubeeq_expect = fml
}

func (minatom *quadEQ_t) testCompare(ganrac *Ganrac, fml Fof, eyec string) {
	if _cubeeq_expect == nil {
		return
	}

	fmt.Printf("testComparet(%s): %v\n", eyec, fml)
	expect := _cubeeq_expect

	dic := NewDict()
	dic.Set("var", NewInt(1))
	u, _ := funcCAD(ganrac, "cad", []any{NewFmlImpl(fml, expect), dic})

	if _, ok := u.(*AtomT); !ok {
		// 型を表示する
		fmt.Printf("WIP: [eyec=%s] u should be True: type(u)=%T, u=%v\n", eyec, u, u)

		simpl, err := funcCAD(ganrac, "cad", []any{fml, dic})
		if err != nil {
			fmt.Printf("cad(err)=%v\n", err)
			panic("stop")
		}
		fmt.Printf("eactual   : %v\n", simpl)
		fmt.Printf("expect    : %v\n", expect)
		fmt.Printf("exp <=  act: %v\n", u)
		u, _ = funcCAD(ganrac, "cad", []any{NewFmlImpl(expect, fml), dic})
		fmt.Printf("exp  => act: %v\n", u)

		fml2 := NewFmlAnd(fml, expect.Not())
		u, _ = funcCAD(ganrac, "cad", []any{fml2, dic})
		fmt.Printf("Correctly, a false region: actual & not(expect):\n%v\n", u)
		panic("stop")
	}

}

func (minatom *quadEQ_t) vscube(fff FofAO, qeopt QEopt, cond qeCond) Fof {

	// minatom.a を 3 次の等式制約条件として，
	// その根で VS する

	// fpos, fneg は minatom.p の主係数が正の場合と負の場合の
	// 3次等式制約の擬剰余
	fpos, fneg := minatom.prem_qff(fff)
	fmt.Printf("fpos=%v, fneg=%v\n", fpos, fneg)

	disc := qeopt.g.ox.Discrim(minatom.p, minatom.lv)

	ret := make([]Fof, 9)

	minatom.coefs = make([]RObj, 4)

	cgt := NewAtom(minatom.z, GT)
	clt := NewAtom(minatom.z, LT)
	dlt := NewAtom(disc, LT)
	dge := NewAtom(disc, GE)

	///////////////////////////////////////////////
	// minatom.p の主係数が正の場合
	///////////////////////////////////////////////
	for k := uint(0); k < 4; k++ {
		minatom.coefs[k] = minatom.p.Coef(minatom.lv, k)
	}
	minatom.atoms = make(map[OP]*Atom, 4)
	minatom.atoms[LT] = NewAtom(minatom.p, LT).(*Atom)
	minatom.atoms[GT] = NewAtom(minatom.p, GT).(*Atom)
	minatom.atoms[LE] = NewAtom(minatom.p, LE).(*Atom)
	minatom.atoms[GE] = NewAtom(minatom.p, GE).(*Atom)

	minatom.diff = make(map[int]*vs_sample_point, 2)
	minatom.diff[2] = newVsQuadSamplePoint(2, Mul(NewInt(3), minatom.coefs[3]), Mul(two, minatom.coefs[2]), minatom.coefs[1]) // y1
	minatom.diff[1] = minatom.diff[2].newVsQuadSamplePointNeg()                                                               // y2
	for _, v := range minatom.diff {
		v.SetDenSign(1)
	}
	fmt.Printf("minatom=%v\n", minatom)
	fmt.Printf("minatom.p=%v\n", minatom.p)
	fmt.Printf("cgt=%v, clt=%v, minatom.(3: %v, 2:%v)\n", cgt, clt, minatom.coefs[3], minatom.coefs[2])

	for i, tbl := range []struct {
		eyec    string
		disc    Fof // 判別式の符号
		c       Fof // 主係数の符号
		prem    Fof
		fn_quad fn_vs3
		fn_lin  fn_vs3
	}{
		///////////////////////////////////////////////
		// minatom.p の主係数が正の場合
		///////////////////////////////////////////////
		{"+1-x", dlt, cgt, fpos, minatom.vscube_x_quad, minatom.vscube_x_lin}, // 虚根を持つ場合
		{"+3-1", dge, cgt, fpos, minatom.vscube_1_quad, minatom.vscube_1_lin}, // 虚根をもたない＋一番小さい実根
		{"+3-2", dge, cgt, fpos, minatom.vscube_2_quad, minatom.vscube_2_lin}, // 虚根をもたない＋真ん中の実根
		{"+3-3", dge, cgt, fpos, minatom.vscube_3_quad, minatom.vscube_3_lin}, // 虚根をもたない＋一番大きい実根

		{"sep", dlt, clt, fpos, nil, nil}, // separator
		///////////////////////////////////////////////
		// minatom.p の主係数が負の場合
		///////////////////////////////////////////////
		{"-1-x", dlt, clt, fneg, minatom.vscube_x_quad, minatom.vscube_x_lin}, // 虚根を持つ場合
		{"-3-1", dge, clt, fneg, minatom.vscube_1_quad, minatom.vscube_1_lin},
		{"-3-2", dge, clt, fneg, minatom.vscube_2_quad, minatom.vscube_2_lin},
		{"-3-3", dge, clt, fneg, minatom.vscube_3_quad, minatom.vscube_3_lin},
	} {
		ret[i] = falseObj
		if tbl.fn_quad == nil { // separator
			for j := 0; j < 4; j++ {
				minatom.coefs[j] = minatom.coefs[j].Neg()
			}
			minatom.p = minatom.p.neg()
			minatom.atoms[LT], minatom.atoms[GT] = minatom.atoms[GT], minatom.atoms[LT]
			minatom.atoms[LE], minatom.atoms[GE] = minatom.atoms[GE], minatom.atoms[LE]
			minatom.diff[1], minatom.diff[2] = minatom.diff[2], minatom.diff[1]

			continue
		}
		neccon := NewFmlAnds(tbl.disc, tbl.c)
		if _, ok := neccon.(*AtomF); ok {
			continue
		}
		fmt.Printf("   @@@ in [%s,%d] disc=%v, c=%v, prem=%v\n", tbl.eyec, i, tbl.disc, tbl.c, tbl.prem)
		ret[i] = NewFmlAnds(neccon,
			minatom.vscube_z(tbl.prem, qeopt, cond, tbl.fn_quad, tbl.fn_lin))
		fmt.Printf("   @@@ out[%s,%d] = %v\n", tbl.eyec, i, ret[i])

		if true { // @TEST
			// @TEST
			fmt.Printf("minatom =%v\n", minatom)
			fmt.Printf(" diff[1]=%v\n", minatom.diff[1])
			fmt.Printf(" diff[2]=%v\n", minatom.diff[2])
			minatom.testCompare(qeopt.g, ret[i], tbl.eyec)
		}

	}

	return NewFmlOrs(ret...)
}

func (qeopt QEopt) qe_cubeeq(fof FofQ, cond qeCond) Fof {
	fmt.Printf("cubeeq %v\n", fof)
	fff, ok := fof.Fml().(FofAO)
	if !ok {
		fmt.Printf("not a FofAO\n")
		return nil
	}

	var op OP
	if _, ok := fof.(*Exists); ok {
		op = EQ
	} else {
		op = NE
	}

	minatom := &quadEQ_t{}

	// 3 次の等式制約を探す
	for ii, fffi := range fff.Fmls() {
		if atom, ok := fffi.(*Atom); ok && atom.op == op {
			minatom.SetAtomIfEasy(fof, atom, 3, 3, ii)
		}
	}

	// 全体を CAD で解くよりも以下の式を CAD で解くほうが簡単にならないといけない
	// f = a3 x^3 + a2 x^2 + a1 x + a0 として，
	// ex([x], f == 0 && f' rho1 0 && f'' rho2 0 && f''' rho3 0 && fff[i])

	if minatom.a == nil {
		return nil
	}

	if op == NE {
		fff = fff.Not().(FofAO)
		minatom.a = minatom.a.Not().(*Atom)
	}

	fmt.Printf("go! vscube %v\n", fff)
	r := minatom.vscube(fff, qeopt, cond)
	fmt.Printf("gone! vscube %v\n", r)
	if op == NE {
		r = r.Not()
	}
	fmt.Printf("gonn! vscube %v\n", r)

	return r
}
