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

// 主係数の符号を設定
func (tbl *fof_quad_eq) SetSgnLcp(v int) {
	tbl.sgn_lcp = v
}

// ２つの根のうち，+ なら 1, - なら -1
func (tbl *fof_quad_eq) SetSgnS(v int) {
	tbl.sgn_s = v
}

func (qeopt *QEopt) Qe_evenq(prenex_fof Fof, cond qeCond, maxdeg int) Fof {
	if qeopt.g.ox == nil {
		panic("muri-yan")
	}
	return qeopt.qe_evenq(prenex_fof, cond, maxdeg)
}

func QeQuadEq(a Fof, tbl *fof_quad_eq) Fof {
	return a.qe_quadeq(qe_quadeq, tbl)
}

func QeLinEq(a Fof, tbl *fof_quad_eq) Fof {
	return a.qe_quadeq(qe_lineq, tbl)
}

func QeVS(a Fof, lv Level, m int, g *Ganrac) Fof {
	return vs_main(a, lv, m, g)
}

func SimplNum(
	fof simpler, g *Ganrac, true_region, false_region *NumRegion) (Fof, *NumRegion, *NumRegion) {
	return fof.simplNum(g, true_region, false_region)
}

func (poly *Poly) SimplNumUniPoly(t, f *NumRegion) (OP, *NumRegion, *NumRegion) {
	gray := newGrayResion(poly.lv+1, 32)
	gray.setTF(t, f)
	gray.upGray(poly.lv)
	return poly.simplNumUniPoly(gray)
}

func SimplBasic(fof, nec, suf Fof) Fof {
	return fof.simplBasic(nec, suf)
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

func (g *Ganrac) SimplFof(f Fof) Fof {
	return g.simplFof(f, trueObj, falseObj)
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

func (p *Atom) Normalize() Fof {
	return p.normalize()
}
func (p *FmlAnd) Normalize() Fof {
	return p.normalize()
}
func (p *FmlOr) Normalize() Fof {
	return p.normalize()
}

func HasVar(q Fof, lv Level) bool {
	return q.hasVar(lv)
}
func ValidFof(q Fof) error {
	return q.valid()
}
func ValidRObj(q RObj) error {
	return q.valid()
}

func (p *Poly) SetZero(lv Level, deg int) RObj {
	return p.setZero(lv, deg)
}

func SdcQEmain(f *Poly, lv Level, neccon Fof, qeopt QEopt) Fof {
	return sdcQEmain(f, lv, neccon, qeopt)
}

func SdcQEcont(f *Atom, rmin, rmax []*Atom, lv Level, neccon Fof, qeopt QEopt) Fof {
	return sdcQEcont(f, rmin, rmax, lv, neccon, qeopt)
}

func SdcQEpoly(f, rng *Atom, lv Level, neccon Fof, qeopt QEopt) Fof {
	return sdcQEpoly(f, rng, lv, neccon, qeopt)
}

func AtomQE(atom *Atom, lv Level, neccon Fof, qeopt QEopt) Fof {
	return atomQE(atom, lv, neccon, qeopt)
}

func (qeopt *QEopt) SetG(g *Ganrac) {
	qeopt.g = g
}
