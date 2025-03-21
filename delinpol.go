package ganrac

// The McCallum projection, lifting, and order-invariance
// Christopher W. Brown. USNA-TM2005

import (
	"fmt"
)

func (cell *Cell) dim() int {
	n := 0
	for cell.lv >= 0 {
		if cell.isSector() {
			n++
		}
		cell = cell.parent
	}
	return n
}

func (cell *Cell) ancestor(lv Level) *Cell {
	for cell.lv > lv {
		cell = cell.parent
	}
	return cell
}

func (cad *CAD) constcoord_test(cell *Cell, pf ProjFactor) bool {
	// @@1 の前処理. proj. factor にその変数が含まれていたらダメ
	for c := cell.parent; c.lv >= 0; c = c.parent {
		if c.isSector() {
			if pf.P().Deg(c.lv) > 0 {
				return false
			}
		}
	}

	b := false
	ps := NewList()
	for i := Level(0); i <= cell.lv; i++ {
		c := cell.ancestor(i)
		if c.isSector() { // @@1
			// add xi=xi to set $a
			// 前処理により，その変数が含まれていないことは確定
			b = true
			continue
		}

		// step 1: L is a list of definint proj factor
		pfL := make([]ProjFactor, 0)
		n := ps.Len()
		for j := 0; j < len(c.multiplicity); j++ {
			if c.multiplicity[j] > 0 {
				pf := cad.proj[i].get(uint(j))
				pfL = append(pfL, pf)
				ps.Append(pf.P())
				n++
			}
		}
		if !b {
			continue
		}

		// 飽きた
		// step 2:
		// //		gb := cad.g.ox.GB(ps, uint(i+1))
		// 		fmt.Printf(" L=%v\n", pfL)
		// 		fmt.Printf("gb=%v\n", gb)
		panic("stop")
	}

	return true
}

// 射影因子追加不要なら true を返す
func (cad *CAD) need_delineating_poly(cell *Cell, pf ProjFactor) bool {
	// t-order partials の GCD を計算して，それが定数かすでに射影因子に含まれているなら ok
	if err := cell.valid(cad); err != nil {
		fmt.Printf("err: %v\n", err)
		panic("stop")
	}
	plv := pf.P().lv
	a := []*Poly{pf.P()}
	for t := Level(0); ; t++ { // t-order
		b := make([]*Poly, 0, len(a))
		diff := make([]*Poly, 0, len(a))
		// fmt.Printf("need_delineating_poly() a[%d]=%v\n", t, a)
		for _, p := range a {
			for j := Level(0); j <= plv; j++ { // 微分対象
				switch q := p.Diff(j).(type) {
				case *Poly:
					diff = append(diff, q)
					switch qc := cell.reduce(q).(type) {
					case *Poly:
						// 0次元なので，主変数以外は消えるはず.
						if plv != qc.lv {
							// fmt.Printf("[%d,%d/%d] pf=%v, q=%v, qc=%v\n", t, j, pf.P().lv, pf.P(), q, qc)
							return false
						}
						if !qc.IsUnivariate() {
							// 代入できんかったし
							return false
						}
						if qc.Sign() < 0 {
							qc = qc.neg() // projection に保存した形式で
						}
						b = append(b, qc)
					default:
						if qc.IsZero() {
							continue
						} else {
							return true
						}
					}
				default:
					if q.IsZero() {
						continue
					} else {
						return true
					}
				}
				// fmt.Printf("ndp.p=%v\n", cell.reduce(p.Diff(j)))
			}
		}

		if len(b) == 0 { // 代入後全部 0 になった
			// 次のレベルへ
			a = diff
			// fmt.Printf("need_delineating_poly() diff=%v\n", diff)
			if len(diff) == 0 {
				return true
			}
			continue
		}

		g := b[0]
		for k := 1; k < len(b); k++ {
			gg := cad.g.ox.Gcd(g, b[k])
			if gp, ok := gg.(*Poly); ok {
				g = gp
			} else {
				return true
			}
		}

		// g がすでに含まれているか.
		gx := cad.g.ox.Factor(g)
		for k := gx.Len() - 1; k >= 1; k-- {
			fctr, _ := gx.Geti(k)
			g = fctr.(*List).getiPoly(0)
			if g.Sign() < 0 {
				g = g.neg()
			}
			found := false
			for _, p := range cad.proj[g.lv].gets() {
				if g.Equals(p.P()) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
}

// 処理を継続できない場合 false を返す
func (pf *ProjFactorMC) vanishChk(cad *CAD, cell *Cell) bool {

	if int(pf.P().lv) == len(cad.q)-1 {
		return true
	}

	if cell.dim() > 0 {
		// @TODO. constcoord_test() 実装中
		return false
		// if cad.constcoord_test(cell, pf) {
		// }
		// panic("constcoord_test=true")
	} else {
		if cad.need_delineating_poly(cell, pf) {
			return true
		}
	}

	return false
}
