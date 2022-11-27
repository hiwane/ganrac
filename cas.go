package ganrac

// ganrac で実装していない機能を
// 外部 CAS を利用
// - Q[x] 上の因数分解，など

type CAS interface {
	Gcd(p, q *Poly) RObj
	Factor(p *Poly) *List

	Discrim(p *Poly, lv Level) RObj
	Resultant(p *Poly, q *Poly, lv Level) RObj
	Psc(p *Poly, q *Poly, lv Level, j int32) RObj
	Sres(p *Poly, q *Poly, lv Level, k int32) RObj
	GB(p *List, vars *List, n int) *List
	Reduce(p *Poly, gb *List, vars *List, n int) (RObj, bool)

	Eval(p string) (GObj, error)

	Close() error
}
