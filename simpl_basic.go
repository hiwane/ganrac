package ganrac

import (
	"fmt"
)

/////////////////////////////////////
// 論理式の簡単化.
//
// simplification of quantifier-free formulas over ordered fields
// A. Dolzmann, T. Sturm
/////////////////////////////////////

func (p *AtomT) simplBasic(neccon, sufcon Fof) Fof {
	return p
}

func (p *AtomF) simplBasic(neccon, sufcon Fof) Fof {
	return p
}

func simplAtomAnd(p *Atom, neccon *Atom) Fof {
	// fctr されているはず
	if len(neccon.p) > len(p.p) {
		return p
	}
	if len(neccon.p) > 1 {

		found := false
		for i := 0; i < len(neccon.p); i++ {
			found := false
			for j := 0; j < len(p.p); j++ {
				if p.p[j].Equals(neccon.p[i]) {
					found = true
				}
			}
			if !found {
				break
			}
		}
		if found {
			if len(p.p) == len(neccon.p) {
				if (p.op & neccon.op) == 0 {
					return falseObj
				}
				if (p.op | neccon.op) == p.op {
					return trueObj
				}
				return newAtoms(p.p, p.op&neccon.op)
			} else if neccon.op == GT {
				// かぶるやつは不要です...
			}
		}
	}

	flags := make([]bool, len(p.p)) // 更新されたか
	s := 1
	nec := neccon.getPoly()
	for i, pp := range p.p {
		c, b := pp.diffConst(nec)
		// fmt.Printf("c=%d[%1v] nec=%v, target=%v\n", c, b, neccon, pp)
		if !b {
			continue
		}
		// p + c op1 0 : atom
		// p     op2 0 : neccon
		if c < 0 {
			if (neccon.op & (LT | GT)) == LT {
				// if p < 0 or p <= 0, then p-|c| is negative
				flags[i] = true
				s *= -1
			} else if neccon.op == EQ {
				flags[i] = true
				s *= -1
			}
		} else if c > 0 {
			if (neccon.op & (LT | GT)) == GT {
				// if p > 0 or p >= 0, then p+|c| < 0 is positive
				flags[i] = true
			} else if neccon.op == EQ {
				flags[i] = true
			}
		} else if len(p.p) == 1 {
			if (p.op & neccon.op) == 0 {
				return falseObj
			}
			if (p.op | neccon.op) == p.op {
				return trueObj
			}
			return NewAtom(p.p[0], p.op&neccon.op)
		} else {
			switch neccon.op {
			case EQ:
				// 符号確定. 積なので全体で 0
				return NewBool((p.op & EQ) != 0)
			case LT:
				// 符号確定
				flags[i] = true
				s *= -1
			case GT:
				// 符号確定
				flags[i] = true
			case NE:
				// 非ゼロ確定... 符号が影響しない場合は除去できる
				if p.op == EQ || p.op == NE {
					flags[i] = true
				}
			}
		}
	}
	up := false
	for _, f := range flags {
		if f {
			up = true
		}
	}
	if !up {
		return p
	}

	fmls := make([]*Poly, 0, len(p.p))
	for i, fg := range flags {
		if !fg {
			fmls = append(fmls, p.p[i])
		}
	}
	opp := p.op
	if s < 0 {
		opp = opp.neg()
	}
	var ret Fof
	if len(fmls) == 0 {
		ret = NewAtom(one, opp)
	} else {
		ret = newAtoms(fmls, opp)
	}
	// fmt.Printf("p%d=`%v`, nec=`%v` => `%v`\n", len(p.p),p, neccon, ret)

	return ret
}
func simplAtomOr(p *Atom, q *Atom) Fof {
	return simplAtomAnd(p.Not().(*Atom), q.Not().(*Atom)).Not()
}

func (p *Atom) simplBasic(neccon, sufcon Fof) Fof {
	if len(p.p) > 1 && (p.op == EQ || p.op == NE) {
		if p.op == EQ {
			var ret Fof = falseObj
			for _, poly := range p.p {
				ret = NewFmlOr(ret, NewAtom(poly, p.op).simplBasic(neccon, sufcon))
			}
			return ret
		} else {
			var ret Fof = trueObj
			for _, poly := range p.p {
				ret = NewFmlOr(ret, NewAtom(poly, p.op).simplBasic(neccon, sufcon))
			}
			return ret
		}
	}

	switch nn := neccon.(type) {
	case *Atom:
		pp := simplAtomAnd(p, nn)
		if ppp, ok := pp.(*Atom); !ok {
			return pp
		} else {
			p = ppp
		}
	case *FmlAnd:
		for _, f := range nn.fml {
			ff, ok := f.(*Atom)
			if ok {
				pp := simplAtomAnd(p, ff)
				if ppp, ok := pp.(*Atom); !ok {
					return pp
				} else {
					p = ppp
				}
			}
		}
	}
	switch nn := sufcon.(type) {
	case *Atom:
		pp := simplAtomOr(p, nn)
		if ppp, ok := pp.(*Atom); !ok {
			return pp
		} else {
			p = ppp
		}
	case *FmlOr:
		for _, f := range nn.fml {
			ff, ok := f.(*Atom)
			if ok {
				switch q := simplAtomOr(p, ff).(type) {
				case *Atom:
					p = q
				default:
					return q
				}
			}
		}
	}

	return p
}

func (p *FmlAnd) simplBasic(neccon, sufcon Fof) Fof {
	fmls := make([]Fof, len(p.fml))
	for i, f := range p.fml {
		if f.IsQff() {
			fmls[i] = f
		} else {
			fmls[i] = trueObj
		}
	}
	ret := make([]Fof, len(p.fml))
	update := false
	for i := len(fmls) - 1; i >= 0; i-- {
		fmls[i] = neccon
		nc := NewFmlAnds(fmls...)
		ret[i] = p.fml[i].simplBasic(nc, sufcon)
		if ret[i] != p.fml[i] {
			update = true
		}
		if ret[i].IsQff() {
			fmls[i] = ret[i]
		} else {
			fmls[i] = trueObj
		}
	}
	if !update {
		return p
	}
	return NewFmlAnds(ret...)
}

func (p *FmlOr) simplBasic(neccon, sufcon Fof) Fof {
	fmls := make([]Fof, len(p.fml))
	for i, f := range p.fml {
		if f.IsQff() {
			fmls[i] = f
		} else {
			fmls[i] = falseObj
		}
	}
	ret := make([]Fof, len(p.fml))
	update := false
	for i := len(fmls) - 1; i >= 0; i-- {
		fmls[i] = sufcon
		sf := NewFmlOrs(fmls...)
		ret[i] = p.fml[i].simplBasic(neccon, sf)
		if ret[i] != p.fml[i] {
			update = true
		}
		if ret[i].IsQff() {
			fmls[i] = ret[i]
		} else {
			fmls[i] = falseObj
		}
	}
	if !update {
		return p
	}
	return NewFmlOrs(ret...)
}

func simplRreduceHasQvar(f Fof, qs []Level) bool {
	for _, q := range qs {
		if f.hasVar(q) {
			return true
		}
	}
	return false
}

// f から不要な条件を取り除く.
// 不要とは，束縛変数により，実際には違う変数であり，条件として扱ってはいけないもの
//
// 例:   a == 0 && ex([a], a != 0)
// <==>  a == 0 && ex([b], b != 0)
// であり，a != 0 部分は a == 0 では簡単化できない.
//
// Args:
//
//	q: quantified variables
//	isnec: 必要条件なら, true.
func simplBasicRemoveQvar(f Fof, qs []Level, isnec bool) Fof {
	if !f.IsQff() {
		panic(fmt.Sprintf("simplBasicRemoveQvar: not qff: %v", f))
	}
	if fao, ok := f.(FofAO); ok {
		if _, ok := f.(*FmlAnd); ok == isnec {
			// And/Or の要素を個別に見る
			ret := make([]Fof, 0, len(fao.Fmls()))
			up := false
			for _, h := range fao.Fmls() {
				if simplRreduceHasQvar(h, qs) {
					up = true
				} else {
					ret = append(ret, h)
				}
			}
			if up {
				return fao.gen(ret)
			}
			return f
		}
	}

	if simplRreduceHasQvar(f, qs) {
		return NewBool(isnec)
	} else {
		return f
	}
}

func (p *ForAll) simplBasic(neccon, sufcon Fof) Fof {

	neccon = simplBasicRemoveQvar(neccon, p.q, true)
	sufcon = simplBasicRemoveQvar(sufcon, p.q, false)

	fml := p.fml.simplBasic(neccon, sufcon)
	if fml == p.fml {
		return p
	}
	return NewQuantifier(true, p.q, fml)
}

func (p *Exists) simplBasic(neccon, sufcon Fof) Fof {

	neccon = simplBasicRemoveQvar(neccon, p.q, true)
	sufcon = simplBasicRemoveQvar(sufcon, p.q, false)

	fml := p.fml.simplBasic(neccon, sufcon)
	if fml == p.fml {
		return p
	}
	return NewQuantifier(false, p.q, fml)
}
