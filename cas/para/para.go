package para

/*
 * CAS が並列化できないので，
 *
 *
 *
 */

import (
	"fmt"
	. "github.com/hiwane/ganrac"
	"log"
)

type ParaCAS struct {
	cas CAS
	ch  chan paraCASarg
}

type paraCASarg struct {
	kind string
	gob  []GObj
	lv   Level
	j    int
	ch   chan GObj
}

func (para *ParaCAS) loop() {
	for {
		// fmt.Printf("<<paraCAS>>: zzz\n")
		arg := <-para.ch
		// fmt.Printf("<<paraCAS>>: wake up! %v\n", arg.kind)
		switch arg.kind {
		case "gcd":
			r := para.cas.Gcd(arg.gob[0].(*Poly), arg.gob[1].(*Poly))
			arg.ch <- r
		case "fctr":
			r := para.cas.Factor(arg.gob[0].(*Poly))
			arg.ch <- r
		case "discrim":
			r := para.cas.Discrim(arg.gob[0].(*Poly), arg.lv)
			arg.ch <- r
		case "res":
			r := para.cas.Resultant(arg.gob[0].(*Poly), arg.gob[1].(*Poly), arg.lv)
			arg.ch <- r
		case "psc":
			r := para.cas.Psc(arg.gob[0].(*Poly), arg.gob[1].(*Poly), arg.lv, int32(arg.j))
			arg.ch <- r
		case "slope":
			r := para.cas.Slope(arg.gob[0].(*Poly), arg.gob[1].(*Poly), arg.lv, int32(arg.j))
			arg.ch <- r
		case "sres":
			r := para.cas.Sres(arg.gob[0].(*Poly), arg.gob[1].(*Poly), arg.lv, int32(arg.j))
			arg.ch <- r
		case "gb":
			r := para.cas.GB(arg.gob[0].(*List), arg.gob[1].(*List), arg.j)
			arg.ch <- r
		case "reduce":
			r, b := para.cas.Reduce(arg.gob[0].(*Poly), arg.gob[1].(*List), arg.gob[2].(*List), arg.j)
			var bs RObj
			if b {
				bs = NewInt(1)
			} else {
				bs = NewInt(0)
			}
			arg.ch <- NewList(r, bs)
		case "eval":
			r, err := para.cas.Eval(arg.gob[0].(*String).String())
			if err != nil {
				arg.ch <- nil
			} else {
				arg.ch <- r
			}
		}
	}
}

func NewParaCAS(cas CAS) *ParaCAS {
	p := &ParaCAS{cas: cas, ch: make(chan paraCASarg)}
	go p.loop()
	return p
}

func (para *ParaCAS) Close() error {
	return para.cas.Close()
}

func (para *ParaCAS) Gcd(p, q *Poly) RObj {
	var arg paraCASarg
	arg.ch = make(chan GObj)
	arg.kind = "gcd"
	arg.gob = []GObj{p, q}
	para.ch <- arg
	gob := <-arg.ch
	return gob.(RObj)
}

func (para *ParaCAS) Factor(p *Poly) *List {
	var arg paraCASarg
	arg.ch = make(chan GObj)
	arg.kind = "fctr"
	arg.gob = []GObj{p}
	para.ch <- arg
	gob := <-arg.ch
	if gob == nil {
		fmt.Printf("<<main>>: error %v\n", p)
		panic("Factor")
	}
	return gob.(*List)
}

func (para *ParaCAS) Discrim(p *Poly, lv Level) RObj {
	var arg paraCASarg
	arg.ch = make(chan GObj)
	arg.kind = "discrim"
	arg.gob = []GObj{p}
	arg.lv = lv
	para.ch <- arg
	gob := <-arg.ch
	return gob.(RObj)
}

func (para *ParaCAS) Resultant(p *Poly, q *Poly, lv Level) RObj {
	var arg paraCASarg
	arg.ch = make(chan GObj)
	arg.kind = "res"
	arg.gob = []GObj{p, q}
	arg.lv = lv
	para.ch <- arg
	gob := <-arg.ch
	return gob.(RObj)
}

func (para *ParaCAS) Psc(p *Poly, q *Poly, lv Level, j int32) RObj {
	var arg paraCASarg
	arg.ch = make(chan GObj)
	arg.kind = "psc"
	arg.gob = []GObj{p, q}
	arg.lv = lv
	arg.j = int(j)
	para.ch <- arg
	gob := <-arg.ch
	return gob.(RObj)
}

func (para *ParaCAS) Slope(p *Poly, q *Poly, lv Level, k int32) RObj {
	var arg paraCASarg
	arg.ch = make(chan GObj)
	arg.kind = "slope"
	arg.gob = []GObj{p, q}
	arg.lv = lv
	arg.j = int(k)
	para.ch <- arg
	gob := <-arg.ch
	return gob.(RObj)
}

func (para *ParaCAS) Sres(p *Poly, q *Poly, lv Level, cc int32) *List {
	var arg paraCASarg
	arg.ch = make(chan GObj)
	arg.kind = "sres"
	arg.gob = []GObj{p, q}
	arg.lv = lv
	arg.j = int(cc)
	para.ch <- arg
	gob := <-arg.ch
	return gob.(*List)
}

func (para *ParaCAS) GB(p *List, vars *List, n int) *List {
	var arg paraCASarg
	arg.ch = make(chan GObj)
	arg.kind = "gb"
	arg.gob = []GObj{p, vars}
	arg.j = n
	para.ch <- arg
	gob := <-arg.ch
	return gob.(*List)
}

func (para *ParaCAS) Reduce(p *Poly, gb *List, vars *List, n int) (RObj, bool) {
	var arg paraCASarg
	arg.ch = make(chan GObj)
	arg.kind = "reduce"
	arg.gob = []GObj{p, gb, vars}
	arg.j = n
	para.ch <- arg
	gob := <-arg.ch
	glist := gob.(*List)
	r0, _ := glist.Geti(0)
	r1, _ := glist.Geti(1)
	return r0.(RObj), !r1.(RObj).IsZero()
}

func (para *ParaCAS) Eval(p string) (GObj, error) {
	var arg paraCASarg
	arg.ch = make(chan GObj)
	arg.kind = "eval"
	arg.gob = []GObj{NewString(p)}
	para.ch <- arg
	gob := <-arg.ch
	if gob == nil {
		return nil, fmt.Errorf("error in eval")
	}
	return gob, nil
}

func (para *ParaCAS) SetLogger(logger *log.Logger) {
	para.cas.SetLogger(logger)
}
