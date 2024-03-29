package ganrac

// ganrac で実装していない機能を
// 外部 CAS を利用
// - Q[x] 上の因数分解，など

import (
	"log"
)

type CAS interface {
	Gcd(p, q *Poly) RObj

	/* asir format
	 * Factor(p) = [[c, 1], [q2, n2], [q3, n3], ..., [qm, nm]]
	 *  where
	 *     c in Q and
	 * 	   q2, ..., qm in Q[X] and
	 * 	   n2, ..., nm in N and
	 *     p = c* q2^n2 * q3^n3 * ... * qm^nm
	 */
	Factor(p *Poly) *List

	Discrim(p *Poly, lv Level) RObj
	Resultant(p *Poly, q *Poly, lv Level) RObj

	/* principal subresultant coefficient */
	Psc(p *Poly, q *Poly, lv Level, j int32) RObj

	/*
	 * slope resultant
	 * H. Hong. Quantifier elimination for formulas constrained by quadratic equations
	 */
	Slope(p *Poly, q *Poly, lv Level, k int32) RObj

	// Subresultant seq.
	// cc の値により, psc 等に変える
	//   0: subresultant seq.
	//   1: principal subresultant seq.
	//   2: constant term
	//   3: 1+2
	Sres(p *Poly, q *Poly, lv Level, cc int32) *List

	// groebner basis
	// vars 変数順序を表す
	// 束縛変数と自由変数を分けるために block order を用いる. n はそのための区切りの位置を表す.
	// 区切りが必要のない場合は，n=0 を用いる
	GB(p *List, vars *List, n int) *List

	// returns reduce(p, gb), sgn < 0
	// (vars, n) については，GB を参照
	Reduce(p *Poly, gb *List, vars *List, n int) (RObj, bool)

	Eval(p string) (GObj, error)

	Close() error

	SetLogger(logger *log.Logger)
}
