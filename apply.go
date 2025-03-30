package ganrac

import (
	"fmt"
)

func (p *AtomT) Apply(fn ApplyFunc, arg any, qff bool) (Fof, bool) {
	return p, false
}

func (p *AtomF) Apply(fn ApplyFunc, arg any, qff bool) (Fof, bool) {
	return p, false
}

func (p *ForAll) Apply(fn ApplyFunc, arg any, qff bool) (Fof, bool) {
	if qff {
		fmt.Printf("ForAll: %v\n", p)
		panic("!")
	}
	fml, update := p.fml.Apply(fn, arg, qff)
	if !update {
		return p, false
	}

	return NewQuantifier(true, p.q, fml), true
}

func (p *Exists) Apply(fn ApplyFunc, arg any, qff bool) (Fof, bool) {
	if qff {
		fmt.Printf("Exists: %v\n", p)
		panic("!")
	}
	fml, update := p.fml.Apply(fn, arg, qff)
	if !update {
		return p, false
	}

	return NewQuantifier(false, p.q, fml), true
}

func (p *FmlAnd) Apply(fn ApplyFunc, arg any, qff bool) (Fof, bool) {
	update := false
	fs := make([]Fof, p.Len())
	for i, f := range p.Fmls() {
		var u bool
		fs[i], u = f.Apply(fn, arg, qff)
		update = update || u
	}
	if update {
		return p.gen(fs), update
	} else {
		return p, update
	}
}

func (p *FmlOr) Apply(fn ApplyFunc, arg any, qff bool) (Fof, bool) {
	update := false
	fs := make([]Fof, p.Len())
	for i, f := range p.Fmls() {
		var u bool
		fs[i], u = f.Apply(fn, arg, qff)
		update = update || u
	}
	if update {
		return p.gen(fs), update
	} else {
		return p, update
	}
}

func (atom *Atom) Apply(fn ApplyFunc, arg any, qff bool) (Fof, bool) {
	return fn(atom, arg)
}
