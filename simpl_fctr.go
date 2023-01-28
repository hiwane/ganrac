package ganrac

/////////////////////////////////////
// 因数分解するよ
// CAS が必要
//
// simplification of quantifier-free formulas over ordered firlds
// A. Dolzmann, T. Sturm
/////////////////////////////////////

import "fmt"

func (p *AtomT) simplFctr(g *Ganrac) Fof {
	return p
}

func (p *AtomF) simplFctr(g *Ganrac) Fof {
	return p
}

func (p *Atom) simplFctr(g *Ganrac) Fof {
	if p.irreducible {
		return p
	}
	if g.ox == nil {
		fmt.Printf("g.ox=%p:%p\n", g, g.ox)
		fmt.Printf("p=%v\n", p)
		panic("stop")
	}

	pp := [][]*Poly{
		make([]*Poly, 0), // 偶数冪な因数
		make([]*Poly, 0)} // 奇数冪な因数
	sgn := 1
	up := false
	for _, p := range p.p {
		fctr := g.ox.Factor(p)
		// Factor の復帰値の 0 番目は容量
		fctrn, _ := fctr.Geti(0)
		cont, _ := fctrn.(*List).Geti(0)
		sgn *= cont.(RObj).Sign()
		if !cont.(RObj).IsOne() && !cont.(RObj).IsMinusOne() {
			up = true
		}
		for i := fctr.Len() - 1; i > 0; i-- {
			fctrn, _ = fctr.Geti(i)
			ei, _ := fctrn.(*List).Geti(1)
			e := ei.(*Int).Int64()
			pi, _ := fctrn.(*List).Geti(0)
			pp[e%2] = append(pp[e%2], pi.(*Poly))
			if e > 1 {
				up = true
			}
		}
	}
	if !up && len(pp[0]) == 0 && len(pp[1]) == len(p.p) {
		p.irreducible = true
		return p
	}
	var ret Fof
	switch p.op {
	case EQ:
		ret = falseObj
		for _, qq := range pp {
			for _, q := range qq {
				a := NewAtom(q, p.op).(*Atom)
				a.irreducible = true
				a.pmul = q
				ret = NewFmlOr(ret, a)
			}
		}
		return ret
	case NE:
		ret = trueObj
		for _, qq := range pp {
			for _, q := range qq {
				a := NewAtom(q, p.op).(*Atom)
				a.irreducible = true
				a.pmul = q
				ret = NewFmlAnd(ret, a)
			}
		}
		return ret
	}
	op := p.op
	if sgn < 0 {
		op = op.neg()
	}

	if (op & EQ) != 0 { // LE || GE
		if len(pp[1]) == 0 && op == GE {
			return trueObj
		}
		ret = falseObj
		for _, q := range pp[0] {
			a := NewAtom(q, EQ).(*Atom)
			a.irreducible = true
			a.pmul = q
			ret = NewFmlOr(ret, a)
		}
		if len(pp[1]) > 0 {
			qq := make([]RObj, len(pp[1]))
			for i := 0; i < len(qq); i++ {
				qq[i] = pp[1][i]
			}
			a := NewAtoms(qq, op).(*Atom)
			a.irreducible = true
			ret = NewFmlOr(ret, a)
		}
	} else if len(pp[1]) == 0 && op == LT {
		return falseObj
	} else if len(pp[1]) == 0 && op == GE {
		return trueObj
	} else { // LT || GT
		ret = trueObj
		for _, q := range pp[0] {
			a := NewAtom(q, NE).(*Atom)
			a.irreducible = true
			a.pmul = q
			ret = NewFmlAnd(ret, a)
		}
		if len(pp[1]) > 0 {
			qq := make([]RObj, len(pp[1]))
			for i := 0; i < len(qq); i++ {
				qq[i] = pp[1][i]
			}
			a := NewAtoms(qq, op).(*Atom)
			a.irreducible = true
			ret = NewFmlAnd(ret, a)
		}
	}
	return ret
}

func (p *FmlAnd) simplFctr(g *Ganrac) Fof {
	fml := make([]Fof, len(p.fml))
	for i := 0; i < len(fml); i++ {
		fml[i] = p.fml[i].simplFctr(g)
	}

	return NewFmlAnds(fml...)
}

func (p *FmlOr) simplFctr(g *Ganrac) Fof {
	fml := make([]Fof, len(p.fml))
	for i := 0; i < len(fml); i++ {
		fml[i] = p.fml[i].simplFctr(g)
	}

	return NewFmlOrs(fml...)
}

func (p *ForAll) simplFctr(g *Ganrac) Fof {
	fml := p.fml.simplFctr(g)
	if fml == p.fml {
		return p
	}
	return NewQuantifier(true, p.q, fml)
}

func (p *Exists) simplFctr(g *Ganrac) Fof {
	fml := p.fml.simplFctr(g)
	if fml == p.fml {
		return p
	}
	return NewQuantifier(false, p.q, fml)
}
