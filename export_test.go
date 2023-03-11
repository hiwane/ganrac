package ganrac

import (
	"math/big"
)

const (
	T_undef = -1 // まだ評価していない
	T_false = 0
	T_true  = 1
	T_other = 2 // 兄弟の情報で親の真偽値が確定したのでもう評価しない
)

var TrueObj = trueObj
var FalseObj = falseObj

type Mult_t mult_t
type Algo_t algo_t

func (iv *Interval) Inf() *big.Float {
	return iv.inf
}
func (iv *Interval) Sup() *big.Float {
	return iv.sup
}

func (qeopt *QEopt) Qe_init(g *Ganrac, fof Fof) {
	qeopt.qe_init(g, fof)
	if g.ox == nil {
		panic("g.ox is null")
	}
}

func (qeopt *QEopt) QE(fof Fof, cond qeCond) Fof {
	return qeopt.qe(fof, cond)
}

func (qeopt QEopt) Qe_neq(fof FofQ, cond qeCond) Fof {
	return qeopt.qe_neq(fof, cond)
}
func (qeopt QEopt) Qe_tran(fof Fof, cond qeCond) Fof {
	return qeopt.qe_tran(fof, cond)
}

func (cad *CAD) Sym_sqfr2(porg *Poly, cell *Cell) []*cadSqfr {
	return cad.sym_sqfr2(porg, cell)
}

func (cad *CAD) Sym_zero_chk(p *Poly, c *Cell) bool {
	return cad.sym_zero_chk(p, c)
}

func (cad *CAD) Symde_gcd2(forg, gorg *Poly, cell *Cell, pi int) (*Poly, *Poly) {
	return cad.symde_gcd2(forg, gorg, cell, pi)
}

func (cad *CAD) Symde_gcd_mod(forg, gorg *Poly, cell *Cellmod, p Uint, need_t bool) (*Poly, Moder, Moder) {
	return cad.symde_gcd_mod(forg, gorg, cell, p, need_t)
}

func (c *Cell) SetParent(d *Cell) {
	c.parent = d
}

func (c *Cell) SetLevel(lv Level) {
	c.lv = lv
}

func (c *Cell) SetDefPoly(p *Poly) {
	c.defpoly = p
}

func (c *Cell) SetIntv(inf, sup *BinInt) {
	c.intv.inf = inf
	c.intv.sup = sup
}

func (c *Cell) SetNintv(intv *Interval) {
	c.nintv = intv
}

func (cell *Cell) Mod(cad *CAD, p Uint) (*Cellmod, bool) {
	return cell.mod(cad, p)
}

func (z *Poly) LC() RObj {
	return z.lc()
}

func (f *Poly) Mod(p Uint) Moder {
	return f.mod(p)
}

func FuncCAD(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	return funcCAD(g, name, args)
}

func (tbl *fof_quad_eq) SetSgnLcp(v int) {
	tbl.sgn_lcp = v
}

func (tbl *fof_quad_eq) SetSgnS(v int) {
	tbl.sgn_s = v
}

func (qeopt *QEopt) Qe_evenq(prenex_fof Fof, cond qeCond) Fof {
	if qeopt.g.ox == nil {
		panic("muri-yan")
	}
	return qeopt.qe_evenq(prenex_fof, cond)
}

func QeQuadEq(a Fof, tbl *fof_quad_eq) Fof {
	return a.qe_quadeq(qe_quadeq, tbl)
}

func QeLinEq(a Fof, tbl *fof_quad_eq) Fof {
	return a.qe_quadeq(qe_lineq, tbl)
}

func SimplNum(
	fof simpler, g *Ganrac, true_region, false_region *NumRegion) (Fof, *NumRegion, *NumRegion) {
	return fof.simplNum(g, true_region, false_region)
}

func (poly *Poly) SimplNumUniPoly(t, f *NumRegion) (OP, *NumRegion, *NumRegion) {
	return poly.simplNumUniPoly(t.getU(f, poly.lv))
}

func SimplFctr(fof simpler, g *Ganrac) Fof {
	return fof.simplFctr(g)
}

func SimplReduce(p Fof, g *Ganrac, inf *reduce_info) Fof {
	return p.simplReduce(g, inf)
}

func (t *NumRegion) Append(lv Level, inf, sup NObj) {
	t.r[lv] = append(t.r[lv], &ninterval{inf, sup})
}

func FofTag(fof Fof) uint {
	return fof.fofTag()
}

func (p *Poly) Powi(y int64) RObj {
	return p.powi(y)
}

func DegModer(m Moder) int {
	return m.deg()
}

func (cell *Cellmod) Factor1() *Poly {
	return cell.factor1
}
func (cell *Cellmod) Factor2() *Poly {
	return cell.factor2
}

func Mul_mod(a, b Moder, p Uint) Moder {
	return a.mul_mod(b, p)
}

func Add_mod(a, b Moder, p Uint) Moder {
	return a.add_mod(b, p)
}
