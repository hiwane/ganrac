package ganrac

import (
	"fmt"
)

/////////////////////////////////////
// 論理式の簡単化.
// 等式制約を利用して，atom を簡単化する
/////////////////////////////////////

// 新しい等式制約が見つかったので, GB (inf.eqns) を更新する
// Returns: GB が更新されたかどうか.
func (inf *reduce_info) updateGB(g *Ganrac, maxvar Level) bool {
	// maxvar := p.maxVar()
	for _, eq := range inf.eqns.Iter() {
		mv := eq.(*Poly).maxVar()
		if mv > maxvar {
			maxvar = mv
		}
	}

	eqns := inf.eqns
	inf.eqns = inf.GB(g, maxvar)
	num := inf.eqns.Len()
	if eqns.Len() != num {
		return true
	}
	for i := 0; i < num; i++ {
		p := inf.eqns.getiPoly(i)
		has_same_poly := false
		for j := i + 1; j < num; j++ {
			q := eqns.getiPoly(j)
			if p.Equals(q) {
				has_same_poly = true
				break
			}
		}
		if !has_same_poly {
			return true
		}
	}
	return false
}

func (ri *reduce_info) Append(p *Poly) {
	ri.eqns.Append(p)
}

// neccon: 必要条件
// sufcon: 十分条件
func NewReduceInfo(g *Ganrac, c, neccon, sufcon Fof) *reduce_info {
	inf := new(reduce_info)
	inf.eqns = NewList()
	if neccon == trueObj && sufcon == falseObj {
		return inf
	}

	eqns := make([]*Poly, 0)
	switch n := neccon.(type) {
	case *FmlOr:
		for _, f := range n.Fmls() {
			if eq, ok := f.(*Atom); ok && eq.op == EQ {
				eqns = append(eqns, eq.getPoly())
			}
		}
	case *Atom:
		if n.op == EQ {
			eqns = append(eqns, n.getPoly())
		}
	}

	switch n := sufcon.(type) {
	case *Atom:
		if n.op == NE {
			eqns = append(eqns, n.getPoly())
		}
	}
	if len(eqns) == 0 {
		return inf
	}

	// 不要な変数が入っていると困る
	maxv := c.maxVar()
	m := len(eqns)
	for lv := Level(0); lv <= maxv; lv++ {
		if !c.hasVar(lv) {
			for i, eq := range eqns {
				if eq != nil && eq.hasVar(lv) {
					eqns[i] = nil
					m--
				}
			}
		}
	}
	if m == 0 {
		return inf
	}
	for _, eq := range eqns {
		if eq != nil {
			inf.Append(eq)
		}
	}
	inf.updateGB(g, c.maxVar())

	return inf
}

func (src *reduce_info) Clone() *reduce_info {
	dest := new(reduce_info)
	*dest = *src
	dest.q = make([]Level, len(src.q))
	copy(dest.q, src.q)
	dest.varb = make([]bool, len(src.varb))
	copy(dest.varb, src.varb)
	dest.depth = src.depth + 1
	dest.eqns = NewListN(src.eqns.Len() + 5)
	for _, p := range src.eqns.Iter() {
		dest.eqns.Append(p)
	}
	return dest
}

func (inf *reduce_info) isQ(lv Level) bool {
	for _, q := range inf.q {
		if q == lv {
			return true
		}
	}
	return false
}

func (inf *reduce_info) Reduce(g *Ganrac, p *Poly) (RObj, bool) {
	if int(p.lv) >= len(inf.varb) {
		b := make([]bool, p.lv+1)
		copy(b, inf.varb)
		inf.varb = b
	}

	n := 0 // 一時的に変数リストを増やす
	for lv := p.lv; lv >= Level(0); lv-- {
		if !inf.varb[lv] && p.hasVar(lv) {
			inf.vars.Append(NewPolyVar(lv))
			n++
		}
	}

	r, neg := g.ox.Reduce(p, inf.eqns, inf.vars, inf.qn+n)
	inf.vars.v = inf.vars.v[:inf.vars.Len()-n] // 元に戻す

	return r, neg
}

func (inf *reduce_info) GB(g *Ganrac, lvmax Level) *List {
	quan := make([]Level, 0, lvmax)
	free := make([]Level, 0, lvmax)
	varb := make([]bool, lvmax+1)
	for lv := Level(0); lv <= lvmax; lv++ {
		for i := inf.eqns.Len() - 1; i >= 0; i-- {
			e := inf.eqns.getiPoly(i)
			if e.hasVar(lv) {
				if inf.isQ(lv) {
					quan = append(quan, lv)
				} else {
					free = append(free, lv)
				}
				varb[lv] = true
				break
			}
		}
	}

	vars := NewList()
	for _, lv := range free {
		vars.Append(NewPolyVar(lv))
	}
	for _, lv := range quan {
		vars.Append(NewPolyVar(lv))
	}
	inf.vars = vars
	inf.varb = varb
	inf.qn = len(quan)

	g.log(8, 1, "GBi=%v\n", inf.eqns)
	gb := g.ox.GB(inf.eqns, vars, inf.qn)
	g.log(8, 1, "GBo=%v\n", gb)
	return gb
}

func (p *AtomT) simplReduce(g *Ganrac, inf *reduce_info) Fof {
	return p
}

func (p *AtomF) simplReduce(g *Ganrac, inf *reduce_info) Fof {
	return p
}

func (p *Atom) simplReduce(g *Ganrac, inf *reduce_info) Fof {
	if inf.eqns.Len() == 0 {
		return p
	}
	q := p.getPoly()
	r, neg := inf.Reduce(g, q)
	if !q.Equals(r) {
		g.log(3, 1, "simplReduce(Atom) %v => %v [eq=%v]\n", q, r, inf.eqns)
		var a Fof
		if neg {
			a = NewAtom(r, p.op.neg())
		} else {
			a = NewAtom(r, p.op)
		}
		return a.simplFctr(g)
	}
	return p
}

// p の非等式制約を用いて簡単化する.
// p = f != 0 && G として，
//
//	f == 0 && G <==> false ならば，
//
// p <==> G と簡単化できる
func simplReduceAO2(g *Ganrac, infbase *reduce_info, p FofAO, fmls []Fof, op OP, eqn_const []*Poly) Fof {
	if true { // とても遅かった. テスト時間が 7分→11分になった
		return p
	}
	if op != EQ && op != NE {
		panic(fmt.Sprintf("simplReduceAO2() invalid op=%v", op))
	}
	nop := op.not()
	var update bool

	for i, fml := range p.Fmls() {
		if eqn_const[i] != nil {
			continue
		}
		qs := make([]*Atom, 0, len(eqn_const))
		switch f := fml.(type) {
		case *Atom:
			if f.op == nop {
				qs = append(qs, f)
			}
		case FofAO:
			for _, g := range f.Fmls() {
				if a, ok := g.(*Atom); ok && a.op == nop {
					qs = append(qs, a)
				}
			}
		}

		if len(qs) == 0 {
			continue
		}
		for _, qatom := range qs {
			inf := infbase.Clone()
			q := qatom.getPoly()
			inf.Append(q)

			upgb := inf.updateGB(g, p.maxVar())
			if v, ok := inf.eqns.geti(0).(NObj); ok {
				// （コメントの簡略化のため, op=EQ, つまり，p が AND な場合）
				// not(fml[i]) を加えて GB 計算すると false になることがわかった
				// つまり，fml[i] は条件として不要であることがわかった
				if v.Sign() == 0 {
					panic("????")
				}
				// この条件は不要です
				update = true
				fmls[i] = NewBool(nop == NE) // この条件は不要です
				break
			}
			if upgb {
				// （コメントの簡略化のため, op=EQ, つまり，p が AND な場合）
				// not(fmls[i]) を加えて GB 計算し，他の AND 要素を簡単化すると false になった.
				// つまり，fmls[i] は条件として不要であることがわかった
				for j := 0; j < len(fmls); j++ {
					if eqn_const[j] == nil && j != i {
						f := fmls[j].simplReduce(g, inf)
						_, okt := f.(*AtomT)
						_, okf := f.(*AtomF)
						if okt && op == NE || okf && op == EQ {
							update = true
							fmls[i] = f.Not() // この条件は不要です
							upgb = false
							break
						}
					}
				}
				if !upgb {
					break
				}
			}
		}
	}

	if update {
		return p.gen(fmls)
	} else {
		return p
	}
}

// p から等式制約を見つけて，and/or の他要素を簡単化する
func simplReduceAO(g *Ganrac, inf *reduce_info, p FofAO, op OP) Fof {
	update := false
	fmls := make([]Fof, len(p.Fmls()))
	// 外側にある等式制約を用いて簡単化
	for i, fml := range p.Fmls() {
		fmls[i] = fml.simplReduce(g, inf)
		if fmls[i] != fml {
			update = true
		}
	}

	n := 0
	// fmls のうち等式制約部分を保持する
	eqn_const := make([]*Poly, len(fmls))
	for i, fml := range fmls {
		switch f := fml.(type) {
		case *Atom:
			if f.op != op {
				continue
			}
			eqn_const[i] = f.getPoly()
			n++
		case FofAO:
			// f_i := poly, q_j := fof として, 以下も等式制約
			//       (f_1 = 0 || f_2 == 0 || f_3 == 0) && q_2 && q_3 && ...
			//  <==>  f_1 * f_2 * f_3 = 0              && q_2 && q_3 && ...
			noeq := false
			for _, g := range f.Fmls() {
				if a, ok := g.(*Atom); !ok || a.op != op {
					noeq = true
					break
				}
			}
			if noeq {
				continue
			}
			// 積が等式制約
			var mul RObj
			for j, f := range fml.(FofAO).Fmls() {
				if j == 0 {
					mul = f.(*Atom).getPoly()
				} else {
					mul = Mul(mul, f.(*Atom).getPoly())
				}
			}
			eqn_const[i] = mul.(*Poly)
			n++
		}
	}

	// eqcon の更新
	if n > 0 {
		// 今の等式制約で簡約... 定数になったら?

		// 新たに等式制約を追加して GB 計算
		inf = inf.Clone()
		for _, eqpoly := range eqn_const {
			if eqpoly != nil {
				inf.Append(eqpoly)
			}
		}

		upgb := inf.updateGB(g, p.maxVar())
		if v, ok := inf.eqns.geti(0).(NObj); ok {
			// GB が [1] になったということは，
			// 等式制約を満たすものがいない
			if v.Sign() == 0 {
				panic("????")
			}
			return NewAtom(v, op)
		}

		if upgb {
			// 更新された GB で再度簡単化
			for i, fml := range fmls {
				if eqn_const[i] != nil {
					continue
				}
				fmls[i] = fml.simplReduce(g, inf)
				if fmls[i] != fml {
					update = true
				}
			}
		}
		if upgb && n >= inf.eqns.Len() {
			// GB 計算したことによって，等式制約の数が減った
			fmls2 := make([]Fof, 0, len(fmls))
			for i, fml := range fmls {
				if eqn_const[i] == nil {
					fmls2 = append(fmls2, fml)
				}
			}

			for _, eq := range inf.eqns.Iter() {
				fmls2 = append(fmls2, NewAtom(eq.(*Poly), op))
			}

			fmls = fmls2
			update = true
		}
	}
	if update {
		return p.gen(fmls)
	} else {
		return simplReduceAO2(g, inf, p, fmls, op, eqn_const)
	}
}

func (p *FmlAnd) simplReduce(g *Ganrac, inf *reduce_info) Fof {
	return simplReduceAO(g, inf, p, EQ)
}

func (p *FmlOr) simplReduce(g *Ganrac, inf *reduce_info) Fof {
	return simplReduceAO(g, inf, p, NE)
}

func simplReduceQ(g *Ganrac, inf *reduce_info, p FofQ) Fof {
	// inf に p.Qs() な変数が含まれていたら，
	// それは別の変数扱いなので，除去が必要
	qs := make([]Level, 0, len(p.Qs()))
	for _, q := range p.Qs() {
		if int(q) < len(inf.varb) && inf.varb[q] {
			qs = append(qs, q)
		}
	}

	if len(qs) > 0 {
		// qs に含まれる変数を, 等式制約から消去する
		infbackup := inf
		inf = infbackup.Clone()
		inf.vars = NewList()
		inf.varb = make([]bool, len(inf.varb))

		for _, q := range qs {
			inf.vars.Append(NewPolyVar(q))
			inf.varb[q] = true
		}
		for _, _v := range infbackup.vars.Iter() {
			v := _v.(*Poly)
			if !inf.varb[v.lv] {
				inf.vars.Append(v)
				inf.varb[v.lv] = true
			}
		}
		if infbackup.vars.Len() != inf.vars.Len() {
			fmt.Printf("old=%v\n", infbackup.vars)
			fmt.Printf("new=%v\n", inf.vars)
			fmt.Printf("qs =%v\n", qs)
			panic("?")
		}
		inf.eqns = g.ox.GB(inf.eqns, inf.vars, inf.vars.Len()-len(qs))

		// GB の要素のなかから，qs が含まれる部分を除去する.
		gb := NewList()
		for _, p := range inf.eqns.Iter() {
			b := true
			for _, q := range qs {
				if p.(*Poly).hasVar(q) {
					b = false
					break
				}
			}
			if b {
				gb.Append(p)
			}
		}

		// inf.varb, inf.qn を更新
		inf.eqns = gb
		g.log(9, 1, "[x] p=%v, varb=%v, %v, gb=%v\n", p, inf.varb, inf.vars, gb)
	}

	fml := p.Fml().simplReduce(g, inf)
	if fml == p.Fml() {
		return p
	}
	return p.gen(p.Qs(), fml)
}

func (p *ForAll) simplReduce(g *Ganrac, inf *reduce_info) Fof {
	return simplReduceQ(g, inf, p)
}

func (p *Exists) simplReduce(g *Ganrac, inf *reduce_info) Fof {
	return simplReduceQ(g, inf, p)
}
