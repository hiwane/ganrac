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
// これが実用的にならないなら，この手法が使えるときはない気がする.
import ()

func (qeopt QEopt) qe_cubeeq(fof FofQ, cond qeCond) Fof {
	fff, ok := fof.Fml().(FofAO)
	if !ok {
		return nil
	}
	if fff.Len() < 4 {
		return nil
	}

	var op OP
	if _, ok := fof.(*Exists); ok {
		op = EQ
	} else {
		op = NE
	}

	minatom := &quadeq_t{}

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
	}

	lv := minatom.lv
	p1 := minatom.p.Diff(lv)
	p2 := p1.(*Poly).Diff(lv)
	p3 := minatom.p.Coef(lv, 3)
	atom := minatom.a

	ret := make([]Fof, 0, 8) // 8 は次の for 文の配列サイズ
	for _, tbl := range []struct {
		f Fof
	}{
		{
			NewFmlAnds(atom, NewAtom(p3, GT), NewAtom(p1, GE), NewAtom(p2, GT)),
		}, {
			NewFmlAnds(atom, NewAtom(p3, GT), NewAtom(p1, GE), NewAtom(p2, EQ)),
		}, {
			NewFmlAnds(atom, NewAtom(p3, GT), NewAtom(p1, GE), NewAtom(p2, LT)),
		}, {
			NewFmlAnds(atom, NewAtom(p3, GT), NewAtom(p1, LT)),
		}, {
			NewFmlAnds(atom, NewAtom(p3, LT), NewAtom(p1, LE), NewAtom(p2, GT)),
		}, {
			NewFmlAnds(atom, NewAtom(p3, LT), NewAtom(p1, LE), NewAtom(p2, EQ)),
		}, {
			NewFmlAnds(atom, NewAtom(p3, LT), NewAtom(p1, LE), NewAtom(p2, LT)),
		}, {
			NewFmlAnds(atom, NewAtom(p3, LT), NewAtom(p1, GT)),
		},
	} {

		o := make([]Fof, 0, fff.Len())
		for ii, fffi := range fff.Fmls() {
			if ii == minatom.idx {
				continue
			}
			newfml := NewExists([]Level{lv}, NewFmlAnd(tbl.f, fffi))
			qeopt.DelAlgo(QEALGO_EQCUBE) // 無限ループ回避
			qff := qeopt.qe(newfml.(FofQ), cond)
			o = append(o, qff)
		}

		ret = append(ret, NewFmlAnds(o...))
	}

	r := NewFmlOrs(ret...)
	if op == NE {
		r = r.Not()
	}

	return r
}
