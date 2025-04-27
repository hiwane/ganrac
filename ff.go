package ganrac

import (
	"fmt"
)

// Field でいいんだけど, add とか別で定義しちゃったし，
// mod 目的でしか使わないから
type FiniteField interface {
	// FF 上での 0 を返す.
	zero() NObj

	// Z[x] とかが与えられて, FF[x] に変換するときに利用する
	toFF(f NObj) NObj

	// 以下の f, g は FF 上の要素と仮定
	negMod(f NObj) NObj	// -f
	invMode(f NObj) NObj	// f^-1
	addMod(f, g NObj) NObj	// f + g mod FF
	subMod(f, g NObj) NObj	// f - g mod FF
	mulMod(f, g NObj) NObj	// f * g mod FF

	// f が FF 上で valid か判定する
	validMod(f NObj) error
}

func negMod(aa RObj, p FiniteField) RObj {
	switch a := aa.(type) {
	case *Poly:
		z := a.NewPoly()
		for i := 1; i < len(a.c); i++ {
			z.c[i] = negMod(a.c[i], p)
		}
		return z
	case NObj:
		return p.negMod(a)
	}
	return nil
}

func addMod(aa, bb RObj, p FiniteField) RObj {
	switch a := aa.(type) {
	case NObj:
		switch b := bb.(type) {
		case NObj:
			return p.addMod(a, b)
		case *Poly:
			z := b.Clone()
			z.c[0] = addMod(a, b.c[0], p)
			return z
		}
	case *Poly:
		switch b := bb.(type) {
		case NObj:
			z := a.Clone()
			z.c[0] = addMod(a.c[0], b, p)
			return z
		case *Poly:
			return addModPoly(a, b, p)
		}
	}
	return nil
}

// 入力 は Mod p 上な要素と仮定
func addModPoly(f, g *Poly, p FiniteField) RObj {
	if f.lv < g.lv {
		z := g.Clone()
		z.c[0] = addMod(z.c[0], f, p)
		return z
	} else if f.lv > g.lv {
		z := f.Clone()
		z.c[0] = addMod(z.c[0], g, p)
		return z
	} else if len(f.c) == len(g.c) {
		var i int
		var c RObj
		for i = len(f.c) - 1; i >= 0; i-- {
			c = addMod(f.c[i], g.c[i], p)
			if !c.IsZero() {
				break
			}
		}
		if i <= 0 {
			return c
		}
		q := NewPoly(f.lv, i+1)
		q.c[i] = c
		for i--; i >= 0; i-- {
			q.c[i] = addMod(f.c[i], g.c[i], p)
		}
		return q
	} else {
		var dmin int
		var q *Poly
		if len(f.c) < len(g.c) {
			dmin = len(f.c)
			q = g
		} else {
			dmin = len(g.c)
			q = f
		}
		z := NewPoly(q.lv, len(q.c))
		copy(z.c[dmin:], q.c[dmin:])
		for i := 0; i < dmin; i++ {
			z.c[i] = addMod(f.c[i], g.c[i], p)
		}
		return z
	}
}

func subMod(aa, bb RObj, p FiniteField) RObj {
	switch a := aa.(type) {
	case NObj:
		switch b := bb.(type) {
		case NObj:
			return p.subMod(a, b)
		case *Poly:
			z := b.NewPoly()
			for i := 1; i < len(b.c); i++ {
				z.c[i] = negMod(b.c[i], p)
			}
			z.c[0] = subMod(a, z.c[0], p)
			return z
		}
	case *Poly:
		switch b := bb.(type) {
		case NObj:
			z := a.Clone()
			z.c[0] = subMod(a.c[0], b, p)
			return z
		case *Poly:
			return subModPoly(a, b, p)
		}
	}
	return nil
}

// 入力 は Mod p 上な要素と仮定
func subModPoly(f, g *Poly, p FiniteField) RObj {
	if f.lv < g.lv {
		z := g.NewPoly()
		for i := 1; i < len(g.c); i++ {
			z.c[i] = negMod(g.c[i], p)
		}
		z.c[0] = subMod(f, g.c[0], p)
		return z
	} else if f.lv > g.lv {
		z := f.Clone()
		z.c[0] = subMod(z.c[0], g, p)
		return z
	} else if len(f.c) == len(g.c) {
		var i int
		var c RObj
		for i = len(f.c) - 1; i >= 0; i-- {
			c = subMod(f.c[i], g.c[i], p)
			if !c.IsZero() {
				break
			}
		}
		if i <= 0 {
			return c
		}
		q := NewPoly(f.lv, i+1)
		q.c[i] = c
		for i--; i >= 0; i-- {
			q.c[i] = subMod(f.c[i], g.c[i], p)
		}
		return q
	} else if len(f.c) <= len(g.c) {
		z := NewPoly(g.lv, len(g.c))
		for i := 0; i < len(f.c); i++ {
			z.c[i] = subMod(f.c[i], g.c[i], p)
		}
		for i := len(f.c); i < len(g.c); i++ {
			z.c[i] = negMod(g.c[i], p)
		}
		return z
	} else {
		z := NewPoly(f.lv, len(f.c))
		copy(z.c[len(g.c):], g.c[len(g.c):])
		for i := 0; i < len(g.c); i++ {
			z.c[i] = subMod(f.c[i], g.c[i], p)
		}
		return z
	}
}

func mulMod(aa, bb RObj, p FiniteField) RObj {
	switch a := aa.(type) {
	case NObj:
		switch b := bb.(type) {
		case NObj:
			return p.mulMod(a, b)
		case *Poly:
			return mulModPolyNum(b, a, p)
		}
	case *Poly:
		switch b := bb.(type) {
		case NObj:
			return mulModPolyNum(a, b, p)
		case *Poly:
			return mulModPolyPoly(a, b, p)
		}
	}
	return nil
}

func mulModPolyNum(a *Poly, b NObj, p FiniteField) RObj {
	z := a.NewPoly()
	for i := 0; i < len(a.c); i++ {
		z.c[i] = mulMod(a.c[i], b, p)
	}
	return z.normalize()
}

func mulModPolyPoly(f, g *Poly, p FiniteField) RObj {
	if f.lv != g.lv {
		if f.lv > g.lv {
			f, g = g, f
		}
		z := NewPoly(g.lv, len(g.c))
		for i := range z.c {
			z.c[i] = mulMod(f, g.c[i], p)
		}
		if err := z.valid(); err != nil {
			panic(fmt.Sprintf("invalid 1 %v\nz=%v\nf=%v\ng=%v\n", err, z, f, g))
		}
		return z
	}
	if len(f.c) > KARATSUBA_DEG_MOD && len(g.c) > KARATSUBA_DEG_MOD {
		return mulModPolyKaratsuba(f, g, p)
	} else {
		return mulModPolyBasic(f, g, p)
	}
}

func mulModPolyBasic(f, g *Poly, p FiniteField) RObj {
	// assert f.lv = g.lv
	z := NewPoly(f.lv, len(f.c)+len(g.c)-1)
	zero := p.zero()
	for i := range z.c {
		z.c[i] = zero
	}
	for i, c := range f.c {
		if c.IsZero() {
			continue
		}
		fig := mulMod(g, c, p).(*Poly)
		for j := range fig.c {
			c := fig.c[j].(Moder)
			z.c[i+j] = addMod(c, z.c[i+j], p)
		}
	}
	return f
}

func mulModPolyKaratsuba(f, g *Poly, p FiniteField) RObj {
	// returns f*g mod p
	// assert f.lv = g.lv
	// assert len(f.c) > KARATSUBA_DEG_MOD
	// assert len(g.c) > KARATSUBA_DEG_MOD
	var d int
	if len(f.c) > len(g.c) {
		d = len(f.c) / 2
	} else {
		d = len(g.c) / 2
	}

	zero := p.zero()
	f1, f0 := f.karatsuba_divide(d, zero)
	g1, g0 := g.karatsuba_divide(d, zero)

	f1g1 := mulMod(f1, g1, p)
	f0g0 := mulMod(f0, g0, p)
	f10 := subMod(f1, f0, p)
	g01 := subMod(g0, g1, p)
	fg := mulMod(f10, g01, p)
	fg = addMod(fg, f1g1, p)
	fg = addMod(fg, f0g0, p)

	d2 := 2 * d

	// 係数配列に変換
	var cf1g1 []RObj
	if ptmp, ok := f1g1.(*Poly); ok && ptmp.lv == f.lv {
		d2 += len(ptmp.c)
		cf1g1 = ptmp.c
	} else {
		cf1g1 = []RObj{f1g1}
	}
	var cf0g0 []RObj
	if ptmp, ok := f0g0.(*Poly); ok && ptmp.lv == f.lv {
		cf0g0 = ptmp.c
	} else {
		cf0g0 = []RObj{f0g0}
	}
	dx := -1
	if q, ok := fg.(*Poly); ok && q.lv == f.lv {
		dx = len(q.c)
	}

	dd := maxint(2*d+len(cf1g1), d+dx)
	ret := NewPoly(f.lv, dd)
	for i := 0; i < len(ret.c); i++ {
		ret.c[i] = zero
	}
	copy(ret.c, cf0g0)
	copy(ret.c[2*d:], cf1g1)

	if q, ok := fg.(*Poly); ok && q.lv == f.lv {
		for i := 0; i < len(q.c); i++ {
			ret.c[i+d] = addMod(q.c[i], ret.c[i+d], p)
		}
	} else {
		ret.c[d] = addMod(fg, ret.c[d], p)
	}

	return ret
}
