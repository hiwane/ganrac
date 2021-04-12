package ganrac

import (
	"fmt"
	"io"
	"strings"
)

type Level uint

// poly in K[x_lv,...,x_n]
type Poly struct { // recursive expression
	lv Level
	indeter
	c []RObj
}

func NewPolyVar(lv Level) *Poly {
	return varlist[lv].p
}

func newPolyVarn(lv Level, deg int) *Poly {
	// return x[lv]^deg
	p := NewPoly(lv, deg+1)
	for i := 0; i < deg; i++ {
		p.c[i] = zero
	}
	p.c[deg] = one
	return p
}

func NewPoly(lv Level, deg_1 int) *Poly {
	p := new(Poly)
	p.c = make([]RObj, deg_1)
	p.lv = lv
	return p
}

func NewPolyInts(lv Level, coeffs ...int64) *Poly {
	p := NewPoly(lv, len(coeffs))
	for i, c := range coeffs {
		p.c[i] = NewInt(c)
	}
	return p
}

func NewPolyCoef(lv Level, coeffs ...RObj) *Poly {
	p := NewPoly(lv, len(coeffs))
	for i, c := range coeffs {
		p.c[i] = c
	}
	return p
}

func (z *Poly) valid() error {
	if z.c == nil {
		return fmt.Errorf("coefs is null")
	}
	if len(z.c) < 2 {
		return fmt.Errorf("numeric... %v", z)
	}
	if z.c[len(z.c)-1].IsZero() {
		return fmt.Errorf("lc should not be zero... %v", z)
	}
	for i, c := range z.c {
		if c == nil {
			return fmt.Errorf("coef[%d] is null", i)
		}
		err := c.valid()
		if err != nil {
			return err
		}
		if cp, ok := c.(*Poly); ok {
			if cp.lv >= z.lv {
				return fmt.Errorf("invalid level z=%d, coef[%d][%v]", z.lv, i, cp)
			}
		}

	}
	return nil
}

func (z *Poly) Equals(x interface{}) bool {
	p, ok := x.(*Poly)
	if !ok {
		return false
	}
	if p.lv != z.lv || len(p.c) != len(z.c) {
		return false
	}
	for i := 0; i < len(p.c); i++ {
		if !z.c[i].Equals(p.c[i]) {
			return false
		}
	}
	return true
}

func (z *Poly) Deg(lv Level) int {
	if lv == z.lv {
		return len(z.c) - 1
	} else if lv > z.lv {
		return 0
	}
	m := 0
	for _, c := range z.c {
		p, ok := c.(*Poly)
		if !ok {
			continue
		}
		d := p.Deg(lv)
		if d > m {
			m = d
		}
	}
	return m
}

func (z *Poly) Coef(lv Level, deg uint) RObj {
	if lv == z.lv {
		if deg >= uint(len(z.c)) {
			return zero
		} else {
			return z.c[deg]
		}
	} else if lv > z.lv {
		if deg == 0 {
			return z
		} else {
			return zero
		}
	}
	r := NewPoly(z.lv, len(z.c))
	for i, c := range z.c {
		p, ok := c.(*Poly)
		if ok {
			r.c[i] = p.Coef(lv, deg)
		} else {
			if deg == 0 {
				r.c[i] = c
			} else {
				r.c[i] = zero
			}
		}
	}
	for i := len(z.c) - 1; i > 0; i-- {
		if !r.c[i].IsZero() {
			r.c = r.c[:i+1]
			return r
		}
	}
	return r.c[0]
}

func (z *Poly) Tag() uint {
	return TAG_POLY
}

func (z *Poly) hasVar(lv Level) bool {
	if z.lv < lv {
		return false
	} else if z.lv == lv {
		return true
	}
	for _, c := range z.c {
		cc, ok := c.(*Poly)
		if ok && cc.hasVar(lv) {
			return true
		}
	}
	return false
}

func (z *Poly) Sign() int {
	// sign of leading coefficient
	return z.c[len(z.c)-1].Sign()
}

func (z *Poly) String() string {
	var b strings.Builder
	z.write(&b, false, "*")
	return b.String()
}

func (z *Poly) dump(b io.Writer) {
	fmt.Fprintf(b, "(poly %d %d (", z.lv, len(z.c))
	for _, c := range z.c {
		if c.IsNumeric() {
			fmt.Fprintf(b, "%v", c)
		} else {
			cp := c.(*Poly)
			cp.dump(b)
		}
		fmt.Fprintf(b, " ")
	}
	fmt.Fprintf(b, "))")
}

func (z *Poly) write(b io.Writer, out_sgn bool, mul string) {
	for i := len(z.c) - 1; i >= 0; i-- {
		if s := z.c[i].Sign(); s == 0 {
			continue
		} else {
			if z.c[i].IsNumeric() {
				if s > 0 {
					if i != len(z.c)-1 || out_sgn {
						fmt.Fprintf(b, "+")
					}
					if i == 0 || !z.c[i].IsOne() {
						fmt.Fprintf(b, "%v", z.c[i])
						if i != 0 {
							fmt.Fprintf(b, "*")
						}
					}
				} else {
					if i != 0 && z.c[i].IsMinusOne() {
						fmt.Fprintf(b, "-")
					} else {
						fmt.Fprintf(b, "%v", z.c[i])
						if i != 0 {
							fmt.Fprintf(b, "*")
						}
					}
				}
			} else if p, ok := z.c[i].(*Poly); ok {
				if p.isMono() {
					p.write(b, i != len(z.c)-1, mul)
					if i > 0 {
						fmt.Fprintf(b, "*")
					}
				} else {
					if i != len(z.c)-1 {
						fmt.Fprintf(b, "+")
					}
					if i > 0 {
						fmt.Fprintf(b, "(")
						p.write(b, false, mul)
						fmt.Fprintf(b, ")*")
					} else {
						p.write(b, false, mul)
					}
				}
			}
			if i > 0 {
				fmt.Fprintf(b, "%s", varlist[z.lv].v)
				if i > 1 {
					fmt.Fprintf(b, "^%d", i)
				}
			}
		}
	}
}

func (p *Poly) write_tex(b io.Writer) {
	p.write(b, false, " ")
}

func (p *Poly) write_src(b io.Writer) {
	if p.isUnivariate() {
		bints := true
		for _, cc := range p.c {
			if c, ok := cc.(*Int); !ok || !c.IsInt64() {
				bints = false
				break
			}
		}
		if bints {
			fmt.Fprintf(b, "NewPolyInts(%d", p.lv)
			for _, cc := range p.c {
				c := cc.(*Int).Int64()
				fmt.Fprintf(b, ",%d", c)
			}
			fmt.Fprintf(b, ")")
			return
		}

	}

	fmt.Fprintf(b, "NewPolyCoef(%d", p.lv)
	for _, cc := range p.c {
		fmt.Fprintf(b, ",")
		switch c := cc.(type) {
		case *Poly:
			c.write_src(b)
		case *Int:
			if c.n.IsInt64() {
				fmt.Fprintf(b, "NewInt(%d)", c.Int64())
			} else {
				fmt.Fprintf(b, "ParseInt(\"%v\")", c)
			}
		default:
			panic("nooooo")
		}
	}
	fmt.Fprintf(b, ")")
}

func (z *Poly) isVar() bool {
	return len(z.c) == 2 && z.c[0].IsZero() && z.c[1].IsOne()
}

func (z *Poly) IsZero() bool {
	return false
}

func (z *Poly) IsOne() bool {
	return false
}

func (z *Poly) IsMinusOne() bool {
	return false
}

func (z *Poly) IsNumeric() bool {
	return false
}

func (z *Poly) Set(x RObj) RObj {
	return z
}

func (z *Poly) Copy() RObj {
	return z.copy()
}

func (z *Poly) copy() *Poly {
	u := NewPoly(z.lv, len(z.c))
	for i, c := range z.c {
		u.c[i] = c
	}
	return u
}

func (z *Poly) Neg() RObj {
	x := z.copy()
	for i := 0; i < len(x.c); i++ {
		x.c[i] = x.c[i].Neg()
	}
	return x
}

func (x *Poly) Add(y RObj) RObj {
	if y.IsNumeric() {
		z := x.copy()
		z.c[0] = z.c[0].Add(y)
		return z
	}
	p, _ := y.(*Poly)
	if p.lv < x.lv {
		z := x.copy()
		z.c[0] = p.Add(z.c[0])
		return z
	} else if p.lv > x.lv {
		z := p.copy()
		z.c[0] = x.Add(z.c[0])
		return z
	} else {
		var dmin int
		var q *Poly
		if len(p.c) < len(x.c) {
			dmin = len(p.c)
			q = x
		} else {
			dmin = len(x.c)
			q = p
		}
		z := NewPoly(p.lv, len(q.c))
		for i := 0; i < dmin; i++ {
			z.c[i] = Add(x.c[i], p.c[i])
		}
		for i := dmin; i < len(q.c); i++ {
			z.c[i] = q.c[i]
		}
		return z.normalize()
	}
}

func (z *Poly) normalize() RObj {
	for i := len(z.c) - 1; i > 0; i-- {
		if !z.c[i].IsZero() {
			z.c = z.c[:i+1]
			return z
		}
	}
	return z.c[0]
}

func (z *Poly) Sub(y RObj) RObj {
	// @TODO とりまサボり.
	yn := y.Neg()
	return z.Add(yn)
}

func (x *Poly) Mul(yy RObj) RObj {
	// @TODO とりあえず素朴版 -> Karatsuba へ
	if yy.IsNumeric() {
		if yy.IsZero() {
			return yy
		}
		z := NewPoly(x.lv, len(x.c))
		for i := 0; i < len(x.c); i++ {
			z.c[i] = x.c[i].Mul(yy)
		}
		return z
	}
	y, _ := yy.(*Poly)
	if y.lv < x.lv {
		z := NewPoly(x.lv, len(x.c))
		for i := 0; i < len(x.c); i++ {
			z.c[i] = y.Mul(x.c[i])
		}
		return z
	} else if y.lv > x.lv {
		z := NewPoly(y.lv, len(y.c))
		for i := 0; i < len(y.c); i++ {
			z.c[i] = x.Mul(y.c[i])
		}
		return z
	}
	z := NewPoly(x.lv, len(y.c)+len(x.c)-1)
	for i := 0; i < len(z.c); i++ {
		z.c[i] = zero
	}
	for i := 0; i < len(x.c); i++ {
		if x.c[i].IsZero() {
			continue
		}
		xiyy := y.Mul(x.c[i])
		xiy, _ := xiyy.(*Poly)
		for j := len(xiy.c) - 1; j >= 0; j-- {
			z.c[i+j] = Add(z.c[i+j], xiy.c[j])
		}
	}

	return z
}

func (x *Poly) Div(y NObj) RObj {
	z := NewPoly(x.lv, len(x.c))
	for i, c := range x.c {
		z.c[i] = c.Div(y)
	}
	return z
}

func (p *Poly) leadingCoef() NObj {
	for {
		switch q := p.c[len(p.c)-1].(type) {
		case *Poly:
			p = q
		case NObj:
			return q
		default:
			panic("invalid")
		}
	}
}

func (x *Poly) leadingTerm() *Poly {
	p := NewPoly(x.lv, len(x.c))
	for i := 0; i < len(x.c)-1; i++ {
		p.c[i] = zero
	}
	p.c[len(p.c)-1] = x.c[len(x.c)-1]

	q := p
	for {
		switch c := q.c[len(q.c)-1].(type) {
		case NObj:
			return p
		case *Poly:
			cq := NewPoly(c.lv, len(c.c))
			for i := 0; i < len(cq.c)-1; i++ {
				cq.c[i] = zero
			}
			cq.c[len(cq.c)-1] = c.c[len(c.c)-1]
			q.c[len(q.c)-1] = cq
			q = cq
		default:
			panic("unknown")
		}
	}
}

func sdivlt(x, y *Poly) RObj {
	// return lt(y)/lt(x) if lt(y) is a factor of lt(x)
	// return nil otherwise
	if x.lv != y.lv || len(x.c) < len(y.c) {
		return nil
	}
	var zret *Poly

	if len(x.c) != len(y.c) {
		zret = NewPoly(x.lv, len(x.c)-len(y.c)+1)
		for i := 0; i < len(zret.c)-1; i++ {
			zret.c[i] = zero
		}
		zret.c[len(zret.c)-1] = one

	}
	z := zret

	for j := len(x.c); j >= 0; j-- {

		switch yp := y.c[len(y.c)-1].(type) {
		case NObj:
			c := x.c[len(x.c)-1].Div(yp)
			if z == nil {
				return c
			}
			z.c[len(z.c)-1] = c
			return zret
		case *Poly:
			switch xp := x.c[len(x.c)-1].(type) {
			case *Poly:
				ybak := y
				x = xp
				y = yp
				var c *Poly
				if x.lv < y.lv || len(x.c) < len(y.c) {
					return nil
				} else if x.lv > y.lv {
					c = NewPoly(x.lv, len(x.c))
					y = ybak
				} else if len(x.c) == len(y.c) {
					continue
				} else {
					c = NewPoly(x.lv, len(x.c)-len(y.c)+1)
				}
				for i := 0; i < len(c.c)-1; i++ {
					c.c[i] = zero
				}
				c.c[len(c.c)-1] = one
				if z == nil {
					zret = c
				} else {
					z.c[len(z.c)-1] = c
				}
				z = c
			default:
				fmt.Printf("unexpected: xp=%v, yp=%v\n", xp, yp)
				return nil
			}
		}
	}
	panic("toooooo")
}

func (x *Poly) sdiv(y *Poly) RObj {
	// assume: y is a factor of x
	var ret RObj = zero
	for i := len(x.c); i >= 0; i-- {
		m := sdivlt(x, y)
		if m == nil {
			return nil
		}
		ret = Add(ret, m)
		xx := x.Sub(y.Mul(m))
		if xx.IsNumeric() {
			if xx.IsZero() {
				return ret
			} else {
				return nil
			}
		}
		x = xx.(*Poly)
	}
	panic("toooooo")
}

func (x *Poly) powi(y int64) RObj {
	return x.Pow(NewInt(y))
}

func (x *Poly) Pow(y *Int) RObj {
	// return x^y
	// int版と同じ手法. 通常 x^m 以外では使わないから放置
	// @TODO 2項定理使ったほうが効率的?
	if y.Sign() < 0 {
		return nil // unsupported.
	}
	m := y.n.BitLen() - 1
	if m < 0 {
		return NewInt(1)
	}

	t := x
	var z *Poly
	for i := 0; i < m; i++ {
		if y.n.Bit(i) != 0 {
			if z == nil {
				z = t
			} else {
				zz := z.Mul(t)
				z, _ = zz.(*Poly)
			}
		}

		tt := t.Mul(t)
		t, _ = tt.(*Poly)
	}
	if z == nil {
		return t
	}

	return z.Mul(t)
}

func (z *Poly) Subst(xs []RObj, lvs []Level, idx int) RObj {
	// lvs: sorted

	for ; idx < len(xs) && z.lv < lvs[idx]; idx++ {
	}
	if idx == len(xs) {
		return z
	}
	if z.lv > lvs[idx] {
		p := NewPoly(z.lv, len(z.c))
		for i := 0; i < len(z.c); i++ {
			p.c[i] = z.c[i].Subst(xs, lvs, idx)
		}
		for i := len(z.c) - 1; i > 0; i-- {
			if !p.c[i].IsZero() {
				p.c = p.c[:i+1]
				return p
			}
		}
		return p.c[0]
	}
	x := xs[idx]
	p := z.c[len(z.c)-1].Subst(xs, lvs, idx+1)
	for i := len(z.c) - 2; i >= 0; i-- {
		p = Add(Mul(p, x), z.c[i].Subst(xs, lvs, idx+1))
	}
	return p
}

func (z *Poly) subst1(x RObj, lv Level) RObj {
	return z.Subst([]RObj{x}, []Level{lv}, 0)
}

func (z *Poly) subst_frac(num RObj, dens []RObj, lv Level) RObj {
	// VS からの呼び出しを仮定.
	// dens = [1, den, ..., den^d]
	// d = len(dens) - 1
	// z(x[lv]=num/den) * den^d
	if z.lv > lv {
		p := make([]RObj, len(z.c))
		for i := 0; i < len(z.c); i++ {
			switch zc := z.c[i].(type) {
			case *Poly:
				p[i] = zc.subst_frac(num, dens, lv)
			case NObj:
				p[i] = dens[len(dens)-1].Mul(zc)
			default:
				fmt.Printf("panic! %v\n", zc)
			}
		}
		x := NewPolyVar(z.lv)
		xn := x
		ret := p[0]
		for i := 1; i < len(p); i++ {
			ret = Add(ret, xn.Mul(p[i]))
			xn = xn.Mul(x).(*Poly)
		}
		if err := ret.valid(); err != nil {
			panic("!")
		}

		return ret
	} else if z.lv < lv {
		vv := z.Mul(dens[len(dens)-1])
		if err := vv.valid(); err != nil {
			panic("!")
		}
		return vv
	}

	dd := len(dens) - len(z.c)
	p := Mul(dens[dd], z.c[len(z.c)-1])
	for i := len(z.c) - 2; i >= 0; i-- {
		p = Add(Mul(z.c[i], dens[len(z.c)-i-1+dd]), Mul(p, num))
	}
	if err := p.valid(); err != nil {
		panic("!")
	}
	return p
}

func (z *Poly) mul_2exp(m uint) RObj {
	// assume: z in Z[X]
	p := NewPoly(z.lv, len(z.c))
	for i, cc := range z.c {
		p.c[i] = cc.mul_2exp(m)
	}
	return p
}

func (z *Poly) subst_num_2exp(num RObj, den uint, lv Level, deg int) RObj {
	// z(x = num / 2^den)
	if z.lv > lv {
		p := make([]RObj, len(z.c))
		for i := 0; i < len(z.c); i++ {
			switch zc := z.c[i].(type) {
			case *Poly:
				p[i] = zc.subst_num_2exp(num, den, lv, deg)
			case *Int:
				p[i] = zc.mul_2exp(den * uint(deg))
			default:
				fmt.Printf("panic! %v\n", zc)
			}
		}
		x := NewPolyVar(z.lv)
		xn := x
		ret := p[0]
		for i := 1; i < len(p); i++ {
			ret = Add(ret, xn.Mul(p[i]))
			xn = xn.Mul(x).(*Poly)
		}
		if err := ret.valid(); err != nil {
			panic("!")
		}

		return ret
	} else if z.lv < lv {
		return z.mul_2exp(den * uint(deg))
	}

	dd := deg - len(z.c) // かさ上げ
	p := z.c[len(z.c)-1].mul_2exp(den * uint(dd+1))
	for i := len(z.c) - 2; i >= 0; i-- {
		p = Add(z.c[i].mul_2exp(den*uint(dd+len(z.c)-i)), Mul(p, num))
	}
	if err := p.valid(); err != nil {
		panic("! @TODO")
	}
	return p
}

func (z *Poly) subst_binint_1var(numer *Int, denom uint) RObj {
	// called from realroot()
	// 2^(denom*deg)*p(x=x + numer/2^denom)
	// assume: z is level- lv univariate polynomial in Z[x]
	cc := newInt()
	cc.n.Lsh(one.n, denom)
	q := NewPolyCoef(z.lv, numer, cc)
	p := z.c[len(z.c)-1]
	for i := len(z.c) - 2; i >= 0; i-- {
		p = Add(Mul(p, q), z.c[i].(NObj).mul_2exp(denom*uint(len(z.c)-i-1)))
	}
	return p
}

func (z *Poly) isUnivariate() bool {
	for _, c := range z.c {
		if _, ok := c.(NObj); !ok {
			return false
		}
	}
	return true
}

func (z *Poly) Indets(b []bool) {
	b[z.lv] = true
	for _, c := range z.c {
		if _, ok := c.(indeter); ok {
			c.(indeter).Indets(b)
		}
	}
}

func (z *Poly) maxVar() Level {
	return z.lv + 1
}

func (z *Poly) isMono() bool {
	for i := len(z.c) - 2; i >= 0; i-- {
		if !z.c[i].IsZero() {
			return false
		}
	}
	return true
}

func (z *Poly) LeadinfCoef() NObj {
	switch c := z.c[len(z.c)-1].(type) {
	case NObj:
		return c
	case *Poly:
		return c.LeadinfCoef()
	}
	return nil
}

func (z *Poly) hasSameTerm(pp RObj, lowest bool) bool {
	// 定数以外同じ項をもつか.
	p, ok := pp.(*Poly)
	if !ok {
		return false
	}

	if z.lv != p.lv || len(p.c) != len(z.c) {
		return false
	}
	for i := len(p.c) - 1; i >= 0; i-- {
		switch zz := z.c[i].(type) {
		case *Poly:
			if !zz.hasSameTerm(p.c[i], lowest && i == 0) {
				return false
			}
		case NObj:
			pc, ok := p.c[i].(NObj)
			if !ok {
				return false
			}
			if pc.IsZero() != zz.IsZero() && !lowest {
				return false
			}
		default:
			// しらない
			panic("unknown")
		}
	}
	return true
}

func (z *Poly) diff(lv Level) RObj {
	// 微分
	if z.lv > lv {
		p := NewPoly(z.lv, len(z.c))
		for i, c := range z.c {
			if cp, ok := c.(*Poly); ok {
				p.c[i] = cp.diff(lv)
			} else {
				p.c[i] = zero
			}
		}
		for i := len(p.c) - 1; i > 0; i-- {
			if !p.c[i].IsZero() {
				p.c = p.c[:i+1]
				return p
			}
		}
		return p.c[0]
	} else if z.lv < lv {
		return z
	}

	if len(z.c) == 2 {
		return z.c[1]
	}

	p := NewPoly(z.lv, len(z.c)-1)
	for i := 0; i < len(z.c)-1; i++ {
		p.c[i] = Mul(z.c[i+1], NewInt(int64(i+1)))
	}
	return p
}

func (f *Poly) diffConst(g *Poly) (int, bool) {
	// if exisis c, d in Q s.t. f(x) = c*p(x) + d, returns sign(d), true
	// otherwise returns 0, false
	if !f.hasSameTerm(g, true) {
		return 0, false
	}
	a := f.LeadinfCoef().Abs()
	b := g.LeadinfCoef().Abs()

	return f.Mul(b).Sub(g.Mul(a)).Sign(), true
}

func (f *Poly) toIntv(prec uint) RObj {
	p := NewPoly(f.lv, len(f.c))
	for i, c := range f.c {
		p.c[i] = c.toIntv(prec)
	}
	return p
}

func (p *Poly) _reduce(q *Poly) RObj {
	if p.lv < q.lv {
		return p
	} else if p.lv > q.lv {
		pp := NewPoly(p.lv, len(p.c))
		for i := 0; i < len(p.c); i++ {
			switch c := p.c[i].(type) {
			case *Poly:
				pp.c[i] = c._reduce(q)
			default:
				pp.c[i] = p.c[i]
			}
		}
		return pp.normalize()
	}

	lc := q.c[len(q.c)-1].(NObj)
	for i := 0; p.lv == q.lv && len(p.c) >= len(q.c); i++ {
		cc := p.c[len(p.c)-1].Div(lc)
		qq := NewPoly(p.lv, len(p.c))
		df := len(p.c) - len(q.c)
		for i := 0; i < df; i++ {
			qq.c[i] = zero
		}
		for i := 0; i < len(q.c); i++ {
			qq.c[i+df] = q.c[i].Mul(cc)
		}
		pp := p.Sub(qq)
		if ppp, ok := pp.(*Poly); ok {
			p = ppp
		} else {
			return pp
		}
		if i > 100 {
			panic("err") // DEBUG
		}
	}
	return p
}

func (p *Poly) reduce(q *Poly) RObj {
	// q を使って可能な限り簡単化する.
	// assume: lc(q) in Z, lv(p) >= lv(q)
	d := p.Deg(q.lv)
	if d == 0 {
		return p
	}
	d = d - q.Deg(q.lv)
	if d < 0 {
		return p
	}
	var c NObj
	switch cc := q.c[len(q.c)-1].(type) {
	case NObj:
		c = cc.Abs()
	default:
		return p
	}
	c = c.Pow(NewInt(int64(d + 1))).(NObj)
	p = p.Mul(c).(*Poly)

	switch pp := p._reduce(q).(type) {
	case *Poly:
		return pp.primpart()
	default:
		return pp
	}
}

func (p *Poly) content(k *Int) *Int {
	for i, cc := range p.c {
		switch c := cc.(type) {
		case *Poly:
			k = c.content(k)
		case *Int:
			if c.Equals(zero) {
				p.c[i] = zero
			} else if k == nil {
				k = c
			} else {
				k = k.Gcd(c)
			}
		case *BinInt:
			panic(fmt.Sprintf("unexpected binint %v", c))
		case *Rat:
			panic(fmt.Sprintf("unexpected rat %v", c))
		default:
			panic(fmt.Sprintf("unexpected %v", c))
		}
	}
	return k
}

func (p *Poly) primpart() *Poly {
	// assume: p in Z[X]
	c := p.content(nil)
	return p.Div(c).(*Poly)
}
