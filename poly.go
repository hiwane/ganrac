package ganrac

// 多項式は，再帰表現を採用している.
// CAD における projection を想定し，
// その多項式のレベルは，係数の多項式のレベルよりも大きくしている．
// また，QE への適用を想定しているため，高い次数のものは扱わないので，
// 係数はリストではなく配列で保持している．

import (
	"encoding/binary"
	"fmt"
	"github.com/hiwane/ganrac/cache"
	"hash/fnv"
	"io"
	"math/big"
	"math/rand"
)

type Level int8

type Levels []Level

// poly in K[x_lv,...,x_n]
type Poly struct { // recursive expression
	lv Level
	indeter
	c []RObj
}

func NewPolyVar(lv Level) *Poly {
	if 0 <= int(lv) && int(lv) < len(varlist) {
		return varlist[lv].p
	} else {
		return NewPolyVarn(lv, 1)
	}
}

func NewPolyVarn(lv Level, deg int) *Poly {
	// return x[lv]^deg
	p := NewPoly(lv, deg+1)
	for i := 0; i < deg; i++ {
		p.c[i] = zero
	}
	p.c[deg] = one
	return p
}

func NewPoly(lv Level, deg_plus_1 int) *Poly {
	p := new(Poly)
	p.c = make([]RObj, deg_plus_1)
	p.lv = lv
	return p
}

func NewPolyInts(lv Level, coeffs ...int64) *Poly {
	p := NewPoly(lv, len(coeffs))
	for i, c := range coeffs {
		p.c[i] = NewInt(c)
	}
	if err := p.valid(); err != nil {
		panic(err.Error())
	}
	return p
}

func NewPolyCoef(lv Level, coeffs ...interface{}) *Poly {
	p := NewPoly(lv, len(coeffs))
	for i, cc := range coeffs {
		switch c := cc.(type) {
		case RObj:
			p.c[i] = c
		case int:
			p.c[i] = NewInt(int64(c))
		case int64:
			p.c[i] = NewInt(c)
		default:
			panic("!")
		}
	}
	if err := p.valid(); err != nil {
		panic(err.Error())
	}
	return p
}

func (p *Poly) NewPoly() *Poly {
	return NewPoly(p.lv, len(p.c))
}

func (z *Poly) Clone() *Poly {
	u := z.NewPoly()
	copy(u.c, z.c)
	return u
}

func (z *Poly) valid() error {
	if z == nil {
		return fmt.Errorf("poly is null")
	}
	if z.c == nil {
		return fmt.Errorf("coefs is null")
	}
	if len(z.c) < 2 {
		return fmt.Errorf("[lv=%d,deg=%d] poly deg()<1 ... `%v`", z.lv, len(z.c), z)
	}
	if z.c[len(z.c)-1].IsZero() {
		st := ""
		switch tt := z.c[z.deg()].(type) {
		case *Poly:
			st = "poly"
		case *Int:
			st = "int"
		case Uint:
			st = fmt.Sprintf("uint:%v", uint32(tt))
		case *Rat:
			st = "rat"
		case NObj:
			st = fmt.Sprintf("nobj<%d>", tt.numTag())
		default:
			st = "?"
		}

		return fmt.Errorf("[lv=%d,deg=%d]lc<%s:%v> should not be zero... `%v`", z.lv, z.deg(), st, z.c[z.deg()], z)
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

func (z *Poly) Hash() cache.Hash {
	var buf [8]byte
	h := fnv.New64a()
	h.Write([]byte("$"))
	binary.LittleEndian.PutUint64(buf[:], uint64(z.lv))
	h.Write(buf[:])
	binary.LittleEndian.PutUint64(buf[:], uint64(len(z.c)))
	h.Write(buf[:])
	for _, c := range z.c {
		h.Write([]byte("^"))
		binary.LittleEndian.PutUint64(buf[:], uint64(c.Hash()))
		h.Write(buf[:])
	}
	return cache.Hash(h.Sum64())
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
	// Deg(0) = 0
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
	r := z.NewPoly()
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

func (z *Poly) Level() Level {
	return z.lv
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
	return fmt.Sprintf("%v", z)
}

func (z *Poly) Format(s fmt.State, format rune) {
	switch format {
	case FORMAT_DUMP: // DEBUG (dump)
		fmt.Fprintf(s, "(poly %d %d (", z.lv, len(z.c))
		for _, c := range z.c {
			if c == nil {
				fmt.Fprintf(s, "nil")
			} else if cp, ok := c.(fmt.Formatter); ok {
				cp.Format(s, format)
			} else {
				fmt.Fprintf(s, "%"+string(FORMAT_DUMP), c)
			}
			fmt.Fprintf(s, " ")
		}
		fmt.Fprintf(s, "))")
	case FORMAT_SRC: // source
		z.write_src(s)
	case FORMAT_TEX, FORMAT_QEPCAD:
		z.write(s, format, false, " ", "^", true)
	default:
		if p, ok := s.Precision(); ok {
			ss := fmt.Sprintf("%v", z)
			fmt.Fprintf(s, "%.*s", p, ss)
		} else {
			var pow string
			if format == FORMAT_SYMPY {
				pow = "**"
			} else {
				pow = "^"
			}
			z.write(s, format, false, "*", pow,
				format != FORMAT_INDEX && format != FORMAT_SYMPY)
		}
	}
}

func (z *Poly) write(b fmt.State, format rune, out_sgn bool, mul, pow string, vars bool) {
	// out_sgn 主係数で + の出力が必要ですよ
	for i := len(z.c) - 1; i >= 0; i-- {
		if z.c[i].IsZero() {
			continue
		} else {
			if z.c[i].IsNumeric() {
				s := z.c[i].Sign()
				if s >= 0 {
					if i != len(z.c)-1 || out_sgn {
						fmt.Fprintf(b, "+")
					}
					if i == 0 || !z.c[i].IsOne() {
						z.c[i].Format(b, format)
						if i != 0 {
							fmt.Fprintf(b, "%s", mul)
						}
					}
				} else {
					if i != 0 && z.c[i].IsMinusOne() {
						fmt.Fprintf(b, "-")
					} else {
						z.c[i].Format(b, format)
						if i != 0 {
							fmt.Fprintf(b, "%s", mul)
						}
					}
				}
			} else if p, ok := z.c[i].(*Poly); ok {
				if p.isMono() { // 括弧不要
					p.write(b, format, i != len(z.c)-1 || out_sgn, mul, pow, vars)
					if i > 0 {
						fmt.Fprintf(b, "%s", mul)
					}
				} else {
					if i > 0 {
						if out_sgn || i != len(z.c)-1 {
							fmt.Fprintf(b, "+")
						}
						fmt.Fprintf(b, "(")
						p.write(b, format, false, mul, pow, vars)
						fmt.Fprintf(b, ")%s", mul)
					} else {
						p.write(b, format, true, mul, pow, vars)
					}
				}
			}
			if i > 0 {
				if vars {
					fmt.Fprintf(b, "%s", VarStr(z.lv))
				} else {
					fmt.Fprintf(b, "x%d", z.lv)
				}
				if i >= 10 && mul == " " { // TeX
					fmt.Fprintf(b, "%s{%d}", pow, i)
				} else if i > 1 {
					fmt.Fprintf(b, "%s%d", pow, i)
				}
			}
		}
	}
}

func (p *Poly) write_src(b io.Writer) {
	fmt.Fprintf(b, "NewPolyCoef(%d", p.lv)
	for _, cc := range p.c {
		if c, ok := cc.(*Int); ok && c.IsInt64() {
			fmt.Fprintf(b, ", %v", c)
		} else {
			fmt.Fprintf(b, ", %S", cc)
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

func (z *Poly) Neg() RObj {
	x := z.Clone()
	for i := 0; i < len(x.c); i++ {
		x.c[i] = x.c[i].Neg()
	}
	return x
}

func (z *Poly) neg() *Poly {
	return z.Neg().(*Poly)
}

func (x *Poly) Add(y RObj) RObj {
	p, ok := y.(*Poly)
	if !ok {
		z := x.Clone()
		z.c[0] = z.c[0].Add(y)
		return z
	} else if p.lv < x.lv {
		z := x.Clone()
		z.c[0] = p.Add(z.c[0])
		return z
	} else if p.lv > x.lv {
		z := p.Clone()
		z.c[0] = x.Add(z.c[0])
		return z
	} else if len(p.c) != len(x.c) {
		var q *Poly
		if len(p.c) <= len(x.c) {
			q = x
		} else {
			q = p
			p = x
		}
		// deg(p) < deg(q)
		z := NewPoly(p.lv, len(q.c))
		n := len(p.c)
		for i := 0; i < n; i++ {
			z.c[i] = Add(p.c[i], q.c[i])
		}
		copy(z.c[n:], q.c[n:])
		return z
	} else {
		var i int
		var c RObj
		for i = len(p.c) - 1; i >= 0; i-- {
			c = Add(p.c[i], x.c[i])
			if !c.IsZero() {
				break
			}
		}
		if i <= 0 {
			return c
		}
		q := NewPoly(p.lv, i+1)
		q.c[i] = c
		for i--; i >= 0; i-- {
			q.c[i] = Add(p.c[i], x.c[i])
		}
		return q
	}
}

func (z *Poly) normalize() RObj {
	for i := len(z.c) - 1; i > 0; i-- {
		if !z.c[i].IsZero() {
			z.c = z.c[:i+1]
			return z
		}
	}
	if p, ok := z.c[0].(*Poly); ok {
		return p.normalize()
	} else {
		return z.c[0]
	}
}

func (p *Poly) Sub(y RObj) RObj {
	q, ok := y.(*Poly)
	if !ok || p.lv > q.lv {
		// y は数かレベルが低い
		z := p.Clone()
		z.c[0] = Sub(z.c[0], y)
		return z
	} else if p.lv < q.lv {
		z := q.NewPoly()
		z.c[0] = p.Sub(q.c[0])
		for i := 1; i < len(q.c); i++ {
			z.c[i] = q.c[i].Neg()
		}
		return z
	} else if len(p.c) > len(q.c) {
		z := p.NewPoly()
		m := len(q.c)
		for i := 0; i < m; i++ {
			z.c[i] = Sub(p.c[i], q.c[i])
		}
		copy(z.c[m:], p.c[m:])
		return z
	} else if len(p.c) < len(q.c) {
		z := q.NewPoly()
		for i := 0; i < len(p.c); i++ {
			z.c[i] = Sub(p.c[i], q.c[i])
		}
		for i := len(p.c); i < len(q.c); i++ {
			z.c[i] = q.c[i].Neg()
		}
		return z
	} else {
		var i int
		var c RObj
		for i = len(p.c) - 1; i >= 0; i-- {
			c = Sub(p.c[i], q.c[i])
			if !c.IsZero() {
				break
			}
		}
		if i <= 0 {
			return c
		}
		z := NewPoly(p.lv, i+1)
		z.c[i] = c
		for i--; i >= 0; i-- {
			z.c[i] = Sub(p.c[i], q.c[i])
		}
		return z
	}
}

func (x *Poly) Mul(yy RObj) RObj {
	// karatsuba {"Subst":1497,"Mul":169551}
	// IncDebugCounter("Poly.Mul", 1)
	y, ok := yy.(*Poly)
	if !ok {
		if yy.IsZero() {
			// NOTE: zero を返すと y.(*Interval) のときに動作しなくなる
			return yy
		}
		if yy.IsOne() {
			return x
		}
		z := NewPoly(x.lv, len(x.c))
		for i := 0; i < len(x.c); i++ {
			z.c[i] = x.c[i].Mul(yy)
		}
		return z
	}
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
	} else if len(x.c) > 2 && len(y.c) > 2 {
		return x.mulKaratsuba(y)
	} else {
		return x.mulBasic(y)
	}
}

func mulPolyChan(f, g RObj, ch chan RObj) {
	ch <- Mul(f, g)
}

func (f *Poly) mulKaratsuba(g *Poly) RObj {
	// assert deg(f)>0,deg(g)>0
	// assert max(deg(f),deg(g))>1
	// assert f.lv == g.lv
	m := len(f.c)
	n := len(g.c)
	var d int
	if m >= n {
		d = m / 2
	} else {
		d = n / 2
	}
	f1, f0 := f.karatsuba_divide(d)
	g1, g0 := g.karatsuba_divide(d)

	var f1g1, f0g0, f10g01 RObj
	if true {
		f1g1 = Mul(f1, g1)
		f0g0 = Mul(f0, g0)
		f10 := Sub(f1, f0)
		g01 := Sub(g0, g1)
		f10g01 = Mul(f10, g01)
	} else {
		var cf1g1, cf0g0, cf10g01 chan RObj
		cf1g1 = make(chan RObj)
		go mulPolyChan(f1, g1, cf1g1)
		cf0g0 = make(chan RObj)
		go mulPolyChan(f0, g0, cf0g0)
		f10 := Sub(f1, f0)
		g01 := Sub(g0, g1)

		cf10g01 = make(chan RObj)
		go mulPolyChan(f10, g01, cf10g01)
		f1g1 = <-cf1g1
		f0g0 = <-cf0g0
		f10g01 = <-cf10g01
	}

	h1 := Add(Add(f10g01, f1g1), f0g0)

	// h = f1g1 x^2d + h1 x^d + f0g0
	h := NewPoly(f.lv, m+n-1)
	if p, ok := f1g1.(*Poly); ok && p.lv == f.lv {
		// assert f1g1.lv == f.lv
		if p, ok := f1g1.(*Poly); !ok || p.lv != f.lv {
			panic(fmt.Sprintf("mnd=(%d,%d,%d) f=%v, g=%v", m, n, d, f, g))
		}
		copy(h.c[2*d:], p.c)
	} else {
		h.c[2*d] = f1g1
	}

	var i int
	if p, ok := f0g0.(*Poly); ok && p.lv == f.lv {
		copy(h.c, p.c)
		i = len(p.c)
	} else {
		h.c[0] = f0g0
		i = 1
	}
	for ; i < d; i++ {
		h.c[i] = zero
	}

	// fmt.Printf("d=%d, <%v|%v>, <%v|%v> h=%V, h1=%v\n", d, f0, f1, g0, g1, h, h1)
	if p, ok := h1.(*Poly); ok && p.lv == f.lv {
		for i = 0; i < len(p.c); i++ {
			if h.c[i+d] == nil {
				h.c[i+d] = p.c[i]
			} else {
				h.c[i+d] = Add(h.c[i+d], p.c[i])
			}
		}
	} else {
		if h.c[d] == nil {
			h.c[d] = h1
		} else {
			h.c[d] = Add(h1, h.c[d])
		}
		i = 1
	}
	for ; i < d; i++ {
		if h.c[i+d] == nil {
			h.c[i+d] = zero
		}
	}
	return h
}

/* 素朴な多項式の積 */
func (x *Poly) mulBasic(y *Poly) RObj {
	z := NewPoly(x.lv, len(y.c)+len(x.c)-1)
	for i := range z.c {
		z.c[i] = zero
	}
	for i := range x.c {
		if x.c[i].IsZero() {
			continue
		}
		xiyy := y.Mul(x.c[i])
		xiy, _ := xiyy.(*Poly)
		for j := len(xiy.c) - 1; j >= 0; j-- {
			if z.c[i+j] == zero {
				z.c[i+j] = xiy.c[j]
			} else {
				z.c[i+j] = Add(z.c[i+j], xiy.c[j])
			}
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

// 整数係数多項式か
func (x *Poly) IsZZ() bool {
	for _, c := range x.c {
		switch cc := c.(type) {
		case *Poly:
			if !cc.IsZZ() {
				return false
			}
		case *Int:
			continue
		default:
			return false
		}
	}
	return true
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
	var zret *Poly
	if x.lv < y.lv {
		return nil
	} else if x.lv > y.lv {
		zret = NewPoly(x.lv, len(x.c))
		for i, cc := range x.c {
			switch c := cc.(type) {
			case *Poly:
				zret.c[i] = sdivlt(c, y)
				if zret.c[i] == nil {
					return nil
				}
			default:
				if c.IsZero() {
					zret.c[i] = zero
				} else {
					return nil
				}
			}
		}
		return zret
	}

	if len(x.c) != len(y.c) {
		zret = NewPoly(x.lv, len(x.c)-len(y.c)+1)
		for i := 0; i < len(zret.c)-1; i++ {
			zret.c[i] = zero
		}
		zret.c[len(zret.c)-1] = one

	}
	z := zret

	for j := len(x.c); j >= 0; j-- {
		switch yp := y.lc().(type) {
		case NObj: // y の主係数が定数である.
			c := x.lc().Div(yp)
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
				if x.lv < y.lv || x.lv == y.lv && len(x.c) < len(y.c) {
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

func (xorg *Poly) Sdiv(y *Poly) RObj {
	// assume: y is a factor of x
	x := xorg
	var ret RObj = zero
	for {
		m := sdivlt(x, y)
		if m == nil {
			fmt.Printf("Poly.Sdiv() X=%v\n", xorg)
			fmt.Printf("Poly.Sdiv() x=%v\n", x)
			fmt.Printf("Poly.Sdiv() y=%v\n", y)
			fmt.Printf("Poly.Sdiv() m=%v, \n", m)
			panic("1")
		}
		ret = Add(ret, m)
		xx := x.Sub(y.Mul(m))
		if xx.IsNumeric() {
			if xx.IsZero() {
				return ret
			} else {
				panic("2")
			}
		}
		x = xx.(*Poly)
	}
	// unreachable.
}

func (x *Poly) powi(y int64) RObj {
	if y <= 1 {
		if y == 0 {
			return one
		} else {
			return x
		}
	}
	return x.Pow(NewInt(y))
}

func (x *Poly) Pow(y *Int) RObj {
	// return x^y
	// int版と同じ手法. 通常 x^m 以外では使わないから放置
	// @TODO 2項定理使ったほうが効率的?
	if y.Sign() < 0 {
		return nil // unsupported. pow(-n)
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

// 区間の代入なので 2 乗を利用したい．
func (z *Poly) SubstIntv(x *Interval, lv Level, prec uint) RObj {
	d := z.Deg(lv)

	if d <= 0 {
		return z
	}

	var x2 *Interval
	if d > 1 {
		if x.inf.Sign() > 0 || x.sup.Sign() < 0 {
			x2 = x.Mul(x).(*Interval)
		} else {
			x2 = newInterval(prec)
			x2.inf.SetInt64(0)
			if x.inf.IsInf() || x.sup.IsInf() {
				x2.sup.SetInf(false)
			} else {
				x2.sup.Neg(x.inf)
				if c := x2.sup.Cmp(x.sup); c < 0 {
					x2.sup.Mul(x.sup, x.sup)
				} else {
					x2.sup.Mul(x.inf, x.inf)
				}
			}
		}
	}

	return z.subst_intv(x, x2, lv, prec)
}

// 区間の代入なので 2 乗を利用したい．
func (z *Poly) subst_intv(x, x2 *Interval, lv Level, prec uint) RObj {
	if z.lv > lv {
		p := NewPoly(z.lv, len(z.c))
		for i := 0; i < len(z.c); i++ {
			switch zc := z.c[i].(type) {
			case *Poly:
				p.c[i] = zc.subst_intv(x, x2, lv, prec)
			case NObj:
				p.c[i] = zc
			default:
				fmt.Printf("panic! %v\n", zc)
			}
		}
		return p
	} else if z.lv < lv {
		return z
	}

	var modd int
	var meven int
	if len(z.c)%2 == 0 {
		// 奇数次
		modd = len(z.c) - 1
		meven = len(z.c) - 2

	} else {
		// 偶数次
		modd = len(z.c) - 2
		meven = len(z.c) - 1
	}
	var a RObj
	var b RObj
	a = z.c[modd]
	b = z.c[meven]

	for i := modd - 2; i >= 0; i -= 2 {
		a = Add(a.Mul(x2), z.c[i])
	}
	a = a.Mul(x)
	for i := meven - 2; i >= 0; i -= 2 {
		b = Add(b.Mul(x2), z.c[i])
	}
	return Add(a, b)
}

func (z *Poly) Subst(xs RObj, lv Level) RObj {
	// lvs: sorted
	// IncDebugCounter("Poly.Subst", 1)
	var p RObj
	if z.lv == lv {
		x := xs
		p = z.c[len(z.c)-1]
		for i := len(z.c) - 2; i >= 0; i-- {
			p = Add(Mul(p, x), z.c[i])
		}
	} else {
		x := NewPolyVar(z.lv)
		p = z.c[len(z.c)-1].Subst(xs, lv)
		for i := len(z.c) - 2; i >= 0; i-- {
			p = Add(Mul(p, x), z.c[i].Subst(xs, lv))
		}
	}

	if err := p.valid(); err != nil {
		panic(err.Error())
	}

	return p
}

// returns to, tc s.t.
// z(num + sgn *sqrt($sqr)) = to + tc * sqrt($sqr)
func (z *Poly) subst_frac_sqr(num, sqr RObj, sgn int, dens []RObj, lv Level) (RObj, RObj) {
	// VS からの呼び出しを仮定.
	// dens = [1, den, ..., den^d]
	// d = len(dens) - 1
	// z(x[lv]=(num+sgn*root(sqr)/den) * den^d = to + tc * root(sqr)

	if z.lv > lv {
		p := make([]RObj, len(z.c))
		q := make([]RObj, len(z.c))
		for i := 0; i < len(z.c); i++ {
			switch zc := z.c[i].(type) {
			case *Poly:
				p[i], q[i] = zc.subst_frac_sqr(num, sqr, sgn, dens, lv)
			case NObj:
				p[i] = dens[len(dens)-1].Mul(zc)
				q[i] = zero
			default:
				fmt.Printf("panic! %v\n", zc)
			}
		}
		// 代入している num/den 等がパラメータ付きである場合，
		// level によっては係数に持てないので，別途和を計算する必要がある
		x := NewPolyVar(z.lv)
		xn := x
		to := p[0]
		tc := q[0]
		for i := 1; i < len(p); i++ {
			to = Add(to, xn.Mul(p[i]))
			tc = Add(tc, xn.Mul(q[i]))
			xn = xn.Mul(x).(*Poly)
		}
		return to, tc
	} else if z.lv < lv {
		vv := z.Mul(dens[len(dens)-1])
		if err := vv.valid(); err != nil {
			panic("!")
		}
		return vv, zero
	}

	dd := len(dens) - 1

	// Horner 法は使えない.
	// [z2*x^2 + z1*x + z0](x // (a +- sqrt(b))/d)
	// z2(a^2+b +- 2*a*sqrt(b)) + z1*d*(a +- sqrt(b)) + d^2*z0
	// to + tc*sqrt(b)

	var to RObj
	var tc RObj // 根号の係数

	// 0次
	to = Mul(z.c[0], dens[dd-0]) // 根号の外

	// 1次
	tc = Mul(dens[dd-1], z.c[1])
	to = Add(to, Mul(tc, num))
	if sgn < 0 {
		tc = tc.Neg()
	}

	if z.deg() >= 2 {
		// 2次
		// (a+sqrt(b))^2 = a^2 + b +- 2a*sqrt(b)
		vv := Mul(z.c[2], dens[dd-2])
		to = Add(to, Mul(vv, Add(Mul(num, num), sqr)))
		ss := Mul(vv, Mul(num, two))
		tc = AddOrSub(tc, ss, sgn)
	}
	if z.deg() >= 3 {
		// 3次
		// (a+-sqrt(b))^3 = a^3 +- 3a^2*sqrt(b) + 3a*b +- sqrt(b)^3
		//               =(a^3+3*a*b) +- (3a^2+b)*sqrt(b)
		a3 := Mul(num, NewInt(3))
		oo := Add(Mul(Mul(num, num), num), Mul(sqr, a3))
		cc := Add(Mul(a3, num), sqr)

		vv := Mul(z.c[3], dens[dd-3])
		to = Add(to, Mul(vv, oo))
		tc = AddOrSub(tc, Mul(vv, cc), sgn)
	}
	if z.deg() >= 4 {
		panic(fmt.Sprintf("unsupported VS (deg); %v", z))
	}

	return to, tc
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

func (z *Poly) IsUnivariate() bool {
	for _, c := range z.c {
		if _, ok := c.(NObj); !ok {
			return false
		}
	}
	return true
}

func (z *Poly) isIntPoly() bool {
	for _, cc := range z.c {
		switch c := cc.(type) {
		case *Poly:
			if !c.isIntPoly() {
				return false
			}
		case *Int:
			continue
		default:
			return false
		}
	}
	return true
}

func (z *Poly) isIntvPoly() bool {
	for _, cc := range z.c {
		switch c := cc.(type) {
		case *Poly:
			if !c.isIntPoly() {
				return false
			}
		case *Interval:
			continue
		default:
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
	switch c := z.c[len(z.c)-1].(type) {
	case *Poly:
		return c.isMono()
	default:
		return true
	}
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

func (z *Poly) lc() RObj {
	return z.c[len(z.c)-1]
}

func (z *Poly) deg() int {
	return len(z.c) - 1
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

func (z *Poly) Diff(lv Level) RObj {
	// 微分
	if z.lv > lv {
		p := NewPoly(z.lv, len(z.c))
		for i, c := range z.c {
			if cp, ok := c.(*Poly); ok {
				p.c[i] = cp.Diff(lv)
			} else {
				p.c[i] = zero
			}
		}
		return p.normalize()
	} else if z.lv < lv {
		return zero
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

	u := f.Mul(b).Sub(g.Mul(a))
	if _, ok := u.(*Poly); ok {
		return 0, false
	}
	return u.Sign(), true
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
	for j := 0; p.lv == q.lv && len(p.c) >= len(q.c); j++ {
		cc := p.c[len(p.c)-1].Div(lc)
		qq := NewPoly(p.lv, len(p.c))
		df := len(p.c) - len(q.c)
		for i := 0; i < df; i++ {
			qq.c[i] = zero
		}
		for i := 0; i < len(q.c); i++ {
			qq.c[i+df] = Mul(q.c[i], cc)
		}
		pp := p.Sub(qq)
		if ppp, ok := pp.(*Poly); ok {
			p = ppp
		} else {
			return pp
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

// lc(g) がかけられており, f は擬剰余の準備ができている状態
// return (q, r) where r = forg - g*q
func (forg *Poly) _quorem(g *Poly) (RObj, RObj) {
	var q RObj
	q = zero
	f := forg
	gc := g.lc()
	for {
		var c RObj
		switch gcc := gc.(type) {
		case *Poly:
			// 擬除算やっているので gcc が poly なら f.lc() も poly のはず
			if _, ok := f.lc().(*Poly); !ok {
				fmt.Printf("forg=[%d, %d]\n", forg.lv, len(forg.c)-1)
				fmt.Printf("f   =[%d, %d]\n", f.lv, len(f.c)-1)
			}
			c = sdivlt(f.lc().(*Poly), gcc)
		case *Int:
			c = f.lc().Div(gcc)
		default:
			panic("a")
		}
		var qc RObj
		if len(f.c) > len(g.c) {
			qc = Mul(NewPolyVarn(g.lv, len(f.c)-len(g.c)), c)
		} else {
			qc = c
		}
		q = Add(qc, q)
		if err := q.valid(); err != nil {
			panic(err)
		}
		switch ff := f.Sub(g.Mul(qc)).(type) {
		case NObj:
			return q, ff
		case *Poly:
			if ff.lv == g.lv && len(ff.c) >= len(g.c) {
				f = ff
			} else {
				return q, ff
			}
		}
	}
}

func (fx *Poly) prem(g *Poly) RObj {
	// 擬剰余
	// pquorem は f.lv == g.lv を仮定するが本関数は仮定しない
	// lv := g.lv について, deg(f, lv) < deg(g, lv) になるように剰余を求める

	f := fx
	d := f.Deg(g.lv)
	e := g.deg()
	if d < e {
		return f
	}

	gc := g.lc()
	a := gc.Pow(NewInt(int64(d - e + 1)))
	f = f.Mul(a).(*Poly)
	if err := f.valid(); err != nil {
		panic(err.Error())
	}
	if !f.isZZ() {
		fmt.Printf("fin=%v\na  =%v\nfou=%v\n", fx, a, f)
		panic("not ZZ")
	}

	for d >= e {
		du := uint(d)
		var c RObj
		switch gcc := gc.(type) {
		case *Poly:
			// f が gcc で割り切れることを保証しているため
			// gc が Poly なら lc(f) も poly
			// c は単項で，f-c*g が g よりも次数が落ちるとは限らない
			c = sdivlt(f.Coef(g.lv, du).(*Poly), gcc)
		case *Int:
			c = f.Coef(g.lv, du).Div(gcc)
		default:
			panic("a")
		}
		var qc RObj
		if d > g.deg() {
			qc = Mul(NewPolyVarn(g.lv, d-g.deg()), c)
		} else {
			qc = c
		}
		switch ff := f.Sub(g.Mul(qc)).(type) {
		case NObj:
			return ff
		case *Poly:
			f = ff
			d = f.Deg(g.lv)
		}
		if err := f.valid(); err != nil {
			panic(err.Error())
		}
		if d >= int(du) && false {
			fmt.Printf("f[%d,%s]=%v\n", fx.Deg(g.lv), VarStr(g.lv), fx)
			fmt.Printf("g[%d,%s]=%v\n", g.Deg(g.lv), VarStr(g.lv), g)
			fmt.Printf("gc   =%v\n", gc)
			fmt.Printf("fc   =%v\n", fx.Coef(g.lv, uint(fx.Deg(g.lv))))
			fmt.Printf("F    =%v\n", f)
			fmt.Printf("Fc   =%v\n", f.Coef(g.lv, uint(f.Deg(g.lv))))
			fmt.Printf("d=%d, du=%d, deg(fx)=%d\n", d, du, fx.Deg(g.lv))
			panic("nande?")
		}
	}
	f = f.primpart()
	return f
}

func (f *Poly) pquorem(g *Poly) (RObj, RObj, RObj) {
	// assume: f.lv == g.lv
	// return: (a, q, r) where a * f = q * g + r and deg(r) < deg(g)
	if f.lv != g.lv {
		panic("invalid")
	}
	if len(f.c) < len(g.c) {
		return one, zero, f
	}

	gc := g.lc()
	m := len(f.c) - len(g.c) + 1
	a := gc.Pow(NewInt(int64(m)))
	pp := f.Mul(a).(*Poly)
	q, r := pp._quorem(g)
	if true {
		gq := Mul(g, q)
		sub := pp.Sub(gq)
		if !sub.Equals(r) {
			fmt.Printf("f =%v\n", f)
			fmt.Printf("pp=%v\n", pp)
			fmt.Printf("g =%v\n", g)
			fmt.Printf("q =%v\n", q)
			fmt.Printf("r =%v\n", r)
			panic("boo")
		}
		if rp, ok := r.(*Poly); ok && rp.lv == f.lv && rp.deg() >= g.deg() {
			fmt.Printf("f =%v\n", f)
			fmt.Printf("pp=%v\n", pp)
			fmt.Printf("g =%v\n", g)
			fmt.Printf("q =%v\n", q)
			fmt.Printf("r =%v\n", r)
			panic("baa")
		}
	}

	if err := a.valid(); err != nil {
		panic(err)
	}
	if err := q.valid(); err != nil {
		panic(err)
	}
	if err := r.valid(); err != nil {
		panic(err)
	}
	if gg, ok := gc.(*Poly); ok {
		var m1, m2 int
		q, m1 = pquorem_simplify_max(q, gg, m)
		r, m2 = pquorem_simplify_max(r, gg, m1)
		fmt.Printf(" m1=%d, m2=%d\n", m1, m2)
		if m2 < m1 { // q を割りすぎたので...
			m = m2
			for i := m2; i < m1; i++ {
				q = q.Mul(gc)
			}
		} else {
			m = m1
		}
		a, _ = pquorem_simplify_max(a, gg, m)
	}
	return a, q, r
}

func (f *Poly) div_lt(a RObj) RObj {
	switch aa := a.(type) {
	case *Poly:
		if f.lv == aa.lv {
			if f.deg() < aa.deg() {
				return nil
			}
			c := div_lt(f.lc(), aa.lc())
			if c == nil {
				return nil
			}
			if df := f.deg() - aa.deg(); df > 0 {
				x := NewPolyVarn(f.lv, f.deg()-aa.deg())
				return x.Mul(c)
			} else {
				return c
			}
		} else if f.lv > aa.lv {
			c := div_lt(f.lc(), aa)
			if c == nil {
				return nil
			}
			x := NewPolyVarn(f.lv, f.deg())
			return x.Mul(c)
		}
		return nil
	case *Int:
		c := div_lt(f.lc(), aa)
		if c == nil {
			return nil
		}
		x := NewPolyVarn(f.lv, f.deg())
		return x.Mul(c)
	default:
		return nil
	}
}

func (f *Int) div_lt(a RObj) RObj {
	if aa, ok := a.(*Int); ok {
		q, r := f.QuoRem(aa)
		if !r.IsZero() {
			return nil
		}
		return q
	} else {
		return nil
	}

}

// 主係数を割る
func div_lt(f, a RObj) RObj {
	switch ff := f.(type) {
	case *Poly:
		return ff.div_lt(a)
	case *Int:
		return ff.div_lt(a)
	default:
		return nil
	}
}

func (forg *Poly) pquorem_simplify(a RObj) RObj {
	f := forg
	var ret RObj = zero
	for {
		q := f.div_lt(a)
		if q == nil {
			return nil
		}
		fsub := f.Sub(Mul(q, a))
		switch aa := a.(type) {
		case *Poly:
			ret = Add(ret, q)
			switch fp := fsub.(type) {
			case *Poly:
				if fp.lv < aa.lv || fp.lv == aa.lv && fp.deg() < aa.deg() {
					return nil
				}
				f = fp
			case NObj:
				if fp.IsZero() {
					return ret
				} else {
					return nil
				}
			default:
				return nil
			}
		case *Int:
			return q
		default:
			return nil
		}
	}
}

// q を a で割り切れるだけ割る
// max 割って良い最大回数
func pquorem_simplify_max(qorg RObj, a *Poly, maxn int) (RObj, int) {
	qold := qorg
	maxi := maxn
	for ; maxi > 0; maxi-- {
		if qq, ok2 := qold.(*Poly); ok2 {
			q := qq.pquorem_simplify(a)
			if q == nil {
				break
			}
			qold = q
		} else {
			break
		}
	}
	retn := maxn - maxi
	if true {
		qx := Mul(qold, a.Pow(NewInt(int64(retn))))
		if !qx.Equals(qorg) {
			fmt.Printf("pquorem_simplify_max() invalid\n")
			fmt.Printf("q=%v\n", qorg)
			fmt.Printf("a=%v\n", a)
			fmt.Printf("m=%v, max=%d, maxi=%d\n", retn, maxn, maxi)
			fmt.Printf("q*a^m=%v\n", qx)
			fmt.Printf("ret=%v\n", qold)
			panic("stop")
		}
	}

	return qold, retn
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
				k = c.Abs().(*Int)
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

func (p *Poly) PPC() (*Poly, *Int) {
	c := p.content(nil)
	return p.Div(c).(*Poly), c
}

func (p *Poly) Cmp(q *Poly) int {
	if p.lv != q.lv {
		return int(p.lv - q.lv)
	}
	if p.IsUnivariate() {
		if !q.IsUnivariate() {
			return -1
		}
	} else {
		if q.IsUnivariate() {
			return +1
		}
	}
	if p.deg() != q.deg() {
		return p.deg() - q.deg()
	}

	switch pc := p.lc().(type) {
	case *Int:
		if qc, ok := q.lc().(*Int); ok {
			if pc.IsOne() && !qc.IsOne() {
				return -1
			} else if !pc.IsOne() && qc.IsOne() {
				return +1
			}
		} else {
			return -1
		}
	case *Poly:
		switch qc := q.lc().(type) {
		case *Int:
			return 1
		case *Poly:
			return pc.Cmp(qc)
		}
	}

	return 0
}

func randPoly(r rand.Source, varn, deg, ccoef, num int) *Poly {
	var ret RObj = zero

	coef := int64(ccoef)
	for i := 1; ; i++ {
		var vv RObj

		c := r.Int63()%(2*coef) - coef
		vv = NewInt(c)
		for j := 0; j < varn; j++ {
			d := r.Int63() % int64(deg+1)
			if d != 0 {
				v := NewPolyVar(Level(j)).powi(d)
				vv = v.Mul(vv)
			}
		}
		ret = Add(ret, vv)
		if p, ok := ret.(*Poly); i >= num && ok {
			if err := p.valid(); err != nil {
				panic(fmt.Sprintf("randPoly invalid: %v", p))
			}
			return p
		}
	}
}

func randPolyMod(r rand.Source, varn, deg, ccoef, num int, p Uint) *Poly {
	for i := 0; i < 1000; i++ {
		r := randPoly(r, varn, deg, ccoef, num)
		q, ok := r.mod(p).(*Poly)
		if ok {
			return q
		}
	}
	panic(fmt.Sprintf("randPolyMod van=%v, deg=%v, deg=%v, ccoef=%v, p=%v", varn, deg, ccoef, num, p))
}

func (porg *Poly) pp() (*Poly, RObj) {
	ps := make([]*Poly, 1)
	ps[0] = porg

	var lcm *big.Int
	var gcd *big.Int
	wk := new(big.Int)
	for len(ps) > 0 {
		p := ps[len(ps)-1]
		ps = ps[:len(ps)-1]

		for i := range p.c {
			if p.c[i].IsZero() {
				continue
			}
			switch c := p.c[i].(type) {
			case *Poly:
				ps = append(ps, c)
			case *Int:
				if gcd == nil {
					gcd = new(big.Int)
					gcd.Abs(c.n)
				} else {
					wk.Abs(c.n)
					gcd.GCD(nil, nil, gcd, wk)
				}
			case *Rat:
				if lcm == nil {
					lcm = new(big.Int)
					lcm.Set(c.n.Denom())
				} else {
					wk.GCD(nil, nil, lcm, c.n.Denom())
					wk.Quo(c.n.Denom(), wk)
					lcm.Mul(lcm, wk)
				}
				if gcd == nil {
					gcd = new(big.Int)
					gcd.Abs(c.n.Num())
				} else {
					wk.Abs(c.n.Num())
					gcd.GCD(nil, nil, gcd, wk)
				}
			default:
				panic("")
			}
		}
	}

	if lcm == nil {
		q := newInt()
		q.n = gcd
		p := porg.Div(q)
		return p.(*Poly), q
	} else {
		q := newRat()
		q.n.SetFrac(lcm, gcd)
		p := porg.Mul(q)
		return p.(*Poly), q
	}

}

func (p *Poly) isEven(lv Level) bool {
	if p.lv == lv {
		for i := 1; i < len(p.c); i += 2 {
			if !p.c[i].IsZero() {
				return false
			}
		}
		return true
	} else if p.lv < lv {
		return true
	} else {
		for _, cc := range p.c {
			if c, ok := cc.(*Poly); ok {
				if !c.isEven(lv) {
					return false
				}
			}
		}
		return true
	}
}

func (p *Poly) Less(q *Poly) bool {
	if p.lv != q.lv {
		return p.lv < q.lv
	}
	if len(p.c) != len(q.c) {
		return len(p.c) < len(q.c)
	}
	for i := p.lv - 1; i >= 0; i-- {
		if p.hasVar(i) {
			if !q.hasVar(i) {
				return false
			}
		} else if q.hasVar(i) {
			return true
		}
	}
	for i := p.lv - 1; i >= 0; i-- {
		pd := p.Deg(i)
		qd := q.Deg(i)
		if pd != qd {
			return pd < qd
		}
	}

	return true
}

func (p *Poly) constantTerm() NObj {
	switch q := p.c[0].(type) {
	case *Poly:
		return q.constantTerm()
	case NObj:
		return q
	default:
		panic("ho")
	}
}

func (p *Poly) lm(deg []int) {
	deg[p.lv] = p.deg()
	if q, ok := p.lc().(*Poly); ok {
		q.lm(deg)
	}
}

func (lvs Levels) contains(lv Level) bool {
	// lvs is sorted in ascending order
	for _, v := range lvs {
		if v == lv {
			return true
		}
	}
	return false
}

func (p *Poly) tdeg(lvs Levels) int {
	// total degree w.r.t lvs
	// lvs is sorted in ascending order
	u := 0
	if lvs.contains(p.lv) {
		u = 1
	}

	deg := -1
	for i, cc := range p.c {
		var d int
		if c, ok := cc.(*Poly); ok {
			d = c.tdeg(lvs)
		} else {
			d = 0
		}
		d += u * i
		if deg < d {
			deg = d
		}
	}

	return deg
}

func (p *Poly) discrim2(lv Level) RObj {
	// 次数 2 の場合
	a := p.Coef(lv, 2)
	b := p.Coef(lv, 1)
	c := p.Coef(lv, 0)
	return Sub(Mul(b, b), Mul(four, Mul(a, c)))
}

func (f *Poly) karatsuba_divide(d int) (RObj, RObj) {
	// returns f = f1 * x^d + f0 where deg(f0) < d
	// assert len(f.c) > d
	// assert d > 1
	if len(f.c) <= d {
		return zero, f
	}
	f0 := NewPoly(f.lv, d)
	f1 := NewPoly(f.lv, len(f.c)-d)
	copy(f0.c, f.c[:d])
	copy(f1.c, f.c[d:])
	//	fmt.Printf("KARA.divide %d\nf [%d]=%v\nf1[%d]=%v\nf0[%d]=%v\n", d, len(f.c), f, len(f1.c), f1, len(f0.c), f0)
	//	fmt.Printf("kara.divide %d\nf=%v\nf1=%v\nf0=%v\n", d, f, f1.normalize(), f0.normalize())
	return f1.normalize(), f0.normalize()
}

func (f *Poly) isZZ() bool {
	for _, cc := range f.c {
		switch c := cc.(type) {
		case *Poly:
			if !c.isZZ() {
				return false
			}
		case *Int:
			continue
		default:
			return false
		}
	}
	return true
}

// returns f(-x)
func (f *Poly) NegX(lv Level) *Poly {
	if f.lv == lv {
		z := f.NewPoly()
		for i, cc := range f.c {
			if i%2 == 0 {
				z.c[i] = cc
			} else {
				z.c[i] = cc.Neg()
			}
		}
		return z
	} else if f.lv < lv {
		return f
	} else {
		z := f.NewPoly()
		for i, cc := range f.c {
			if cp, ok := cc.(*Poly); ok {
				z.c[i] = cp.NegX(lv)
			} else {
				z.c[i] = cc
			}
		}
		return z
	}
}

// f.coeff(lv, deg) = 0 になるように f を変形する
func (f *Poly) setZero(lv Level, deg int) RObj {
	if f.lv == lv {
		if deg >= len(f.c) {
			return f
		} else if deg == len(f.c)-1 {
			z := NewPoly(f.lv, deg)
			copy(z.c, f.c)
			return z.normalize()
		} else {
			z := f.Clone()
			z.c[deg] = zero
			return z
		}
	} else if f.lv < lv {
		return f
	} else {
		z := f.NewPoly()
		for i, cc := range f.c {
			if cp, ok := cc.(*Poly); ok {
				if cp.lv < lv && deg == 0 {
					z.c[i] = zero
				} else if cp.lv < lv {
					z.c[i] = cp
				} else {
					z.c[i] = cp.setZero(lv, deg)
				}
			} else {
				if deg == 0 {
					z.c[i] = zero
				} else {
					z.c[i] = cc
				}
			}
		}
		return z.normalize()
	}
}

// return x^n f(1/x)
// assume f(0) != 0, deg(f) <= n
func (f *Poly) SubstXinvLv(lv Level, n int) RObj {
	if f.lv == lv {
		z := NewPoly(lv, n+1)
		for i, cc := range f.c {
			z.c[n-i] = cc
		}
		for i := len(f.c); i <= n; i++ {
			z.c[n-i] = zero
		}
		return z.normalize()
	} else if f.lv < lv {
		return f
	} else {
		z := f.NewPoly()
		for i, cc := range f.c {
			if cp, ok := cc.(*Poly); ok && cp.lv >= lv {
				z.c[i] = cp.SubstXinvLv(lv, n)
			} else {
				z.c[i] = Mul(cc, NewPolyVarn(lv, n))
			}
		}
		return z
	}
}

// ax^2+bx+c の判別式
func Discrim2(a, b, c RObj) RObj {
	return Sub(Mul(b, b), Mul(four, Mul(a, c)))
}
