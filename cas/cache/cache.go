package cache

import (
	"fmt"
	"github.com/hiwane/ganrac"
	"log"
)

type Cache struct {
	cas ganrac.CAS
	qm  *QueMap
}

func New(cas ganrac.CAS, capacity int) *Cache {
	qm := NewQueMap(capacity)
	return &Cache{cas: cas, qm: qm}
}

func (c *Cache) Factor(p *ganrac.Poly) *ganrac.List {
	key := QueHash(fmt.Sprintf("Fctr(%V)", p))
	if v := c.qm.Get(key); v != nil {
		return v.l
	}
	ret := c.cas.Factor(p)
	c.qm.Put(&QueValue{p: p, l: ret, key: key})
	return ret
}

func (c *Cache) Gcd(p, q *ganrac.Poly) ganrac.RObj {
	key := QueHash(fmt.Sprintf("GCD(%V,%V)", p, q))
	if v := c.qm.Get(key); v != nil {
		return v.r
	}
	ret := c.cas.Gcd(p, q)
	c.qm.Put(&QueValue{p: p, q: q, r: ret, key: key})
	return ret
}

func (c *Cache) Discrim(p *ganrac.Poly, lv ganrac.Level) ganrac.RObj {
	key := QueHash(fmt.Sprintf("Dis(%V,%x)", p, lv))
	if v := c.qm.Get(key); v != nil {
		return v.r
	}
	ret := c.cas.Discrim(p, lv)
	c.qm.Put(&QueValue{p: p, lv: lv, r: ret, key: key})
	return ret
}

func (c *Cache) Resultant(p, q *ganrac.Poly, lv ganrac.Level) ganrac.RObj {
	key := QueHash(fmt.Sprintf("Res(%V,%V,%x)", p, q, lv))
	if v := c.qm.Get(key); v != nil {
		return v.r
	}
	ret := c.cas.Resultant(p, q, lv)
	c.qm.Put(&QueValue{p: p, q: q, lv: lv, r: ret, key: key})
	return ret
}

func (c *Cache) Psc(p *ganrac.Poly, q *ganrac.Poly, lv ganrac.Level, j int32) ganrac.RObj {
	key := QueHash(fmt.Sprintf("Psc(%V,%V,%x,%x)", p, q, lv, j))
	if v := c.qm.Get(key); v != nil {
		return v.r
	}
	ret := c.cas.Psc(p, q, lv, j)
	c.qm.Put(&QueValue{p: p, q: q, lv: lv, cc: j, r: ret, key: key})
	return ret
}

func (c *Cache) Slope(p *ganrac.Poly, q *ganrac.Poly, lv ganrac.Level, j int32) ganrac.RObj {
	key := QueHash(fmt.Sprintf("Slp(%V,%V,%x,%x)", p, q, lv, j))
	if v := c.qm.Get(key); v != nil {
		return v.r
	}
	ret := c.cas.Slope(p, q, lv, j)
	c.qm.Put(&QueValue{p: p, q: q, lv: lv, cc: j, r: ret, key: key})
	return ret
}

func (c *Cache) Sres(p *ganrac.Poly, q *ganrac.Poly, lv ganrac.Level, j int32) *ganrac.List {
	key := QueHash(fmt.Sprintf("Srs(%V,%V,%x,%x)", p, q, lv, j))
	if v := c.qm.Get(key); v != nil {
		return v.l
	}
	ret := c.cas.Sres(p, q, lv, j)
	c.qm.Put(&QueValue{p: p, q: q, lv: lv, cc: j, l: ret, key: key})
	return ret
}

func (c *Cache) GB(p *ganrac.List, vars *ganrac.List, n int) *ganrac.List {
	return c.cas.GB(p, vars, n)
}

func (c *Cache) Reduce(p *ganrac.Poly, gb *ganrac.List, vars *ganrac.List, n int) (ganrac.RObj, bool) {
	return c.cas.Reduce(p, gb, vars, n)
}

func (c *Cache) Eval(p string) (ganrac.GObj, error) {
	return c.cas.Eval(p)
}

func (c *Cache) Close() error {
	return c.cas.Close()
}

func (c *Cache) SetLogger(logger *log.Logger) {
	c.cas.SetLogger(logger)
}

func (c *Cache) Hit() int {
	return c.qm.hit
}

func (c *Cache) Miss() int {
	return c.qm.miss
}

func (c *Cache) RemoveCount() int {
	return c.qm.rm
}

func (c *Cache) Capacity() int {
	return c.qm.capacity
}

func (c *Cache) Len() int {
	return len(c.qm.m)
}
