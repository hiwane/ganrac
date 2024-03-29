package ganrac

import (
	"fmt"
)

type List struct {
	GObj
	v []GObj
}

func (z *List) Iter() []GObj {
	return z.v
}

func (z *List) Tag() uint {
	return TAG_LIST
}

func (z *List) String() string {
	s := "["
	for i := 0; i < len(z.v); i++ {
		if i != 0 {
			s += ","
		}
		s += z.v[i].String()
	}
	return s + "]"
}

func (z *List) Format(s fmt.State, format rune) {
	if z == nil {
		fmt.Fprintf(s, "<<nil list>>")
		return
	}
	var left, right string
	switch format {
	case FORMAT_DUMP:
		left = fmt.Sprintf("(list %d", len(z.v))
		right = ")"
	case FORMAT_TEX:
		left = "\\left["
		right = "\\right]"
	case FORMAT_SRC:
		left = "NewList("
		right = ")"
	default:
		left = "["
		right = "]"
	}

	fmt.Fprintf(s, left)
	for i, v := range z.v {
		if i != 0 {
			fmt.Fprintf(s, ", ")
		}
		v.Format(s, format)
	}
	fmt.Fprintf(s, right)
}

func (z *List) Set(ii *Int, v GObj) error {
	ilen := NewInt(int64(len(z.v)))
	if ii.Sign() < 0 || ii.Cmp(ilen) >= 0 {
		return fmt.Errorf("list index out of range")
	}
	m := int(ii.n.Int64())
	z.v[m] = v
	return nil
}

func (z *List) Seti(i int, v GObj) error {
	if i < 0 || i >= len(z.v) {
		return fmt.Errorf("list index out of range")
	}
	z.v[i] = v
	return nil
}

func (z *List) Get(ii *Int) (GObj, error) {
	ilen := NewInt(int64(len(z.v)))
	if ii.Sign() < 0 || ii.Cmp(ilen) >= 0 {
		return nil, fmt.Errorf("list index out of range")
	}
	m := int(ii.n.Int64())
	return z.v[m], nil
}

func (z *List) Geti(idx ...int) (GObj, error) {

	var r GObj = z
	var ok bool
	for j, i := range idx {
		z, ok = r.(*List)
		if !ok {
			return nil, fmt.Errorf("not a list (%d,%d)", j, i)
		}
		if i < 0 || i >= len(z.v) {
			return nil, fmt.Errorf("list index out of range")
		}
		r = z.v[i]
	}
	return r, nil
}

func (z *List) Len() int {
	return len(z.v)
}

func (z *List) Equals(vv interface{}) bool {
	v, ok := vv.(*List)
	if !ok || z.Len() != v.Len() {
		return false
	}
	for i := z.Len() - 1; i >= 0; i-- {
		if c, ok := z.v[i].(equaler); !ok || !c.Equals(v.v[i]) {
			return false
		}
	}
	return true
}

func (z *List) Indets(b []bool) {
	for _, p := range z.v {
		q, ok := p.(indeter)
		if ok {
			q.Indets(b)
		}
	}
}

func (z *List) Swap(i, j int) error {
	if i < 0 || i >= len(z.v) || j < 0 || j >= len(z.v) {
		return fmt.Errorf("list index out of range")
	}
	v := z.v[i]
	z.v[i] = z.v[j]
	z.v[j] = v
	return nil
}

func (z *List) Append(a GObj) {
	z.v = append(z.v, a)
}

func NewList(args ...interface{}) *List {
	lst := new(List)
	lst.v = make([]GObj, len(args))
	for i := 0; i < len(args); i++ {
		lst.v[i] = args[i].(GObj)
	}
	return lst
}

func NewListN(n int, args ...interface{}) *List {
	lst := new(List)
	lst.v = make([]GObj, len(args), n)
	for i := 0; i < len(args); i++ {
		lst.v[i] = args[i].(GObj)
	}
	return lst
}

func (z *List) getiPoly(i int) *Poly {
	// i は正しいと仮定
	return z.geti(i).(*Poly)
}

func (z *List) getiList(i int) *List {
	// i は正しいと仮定
	return z.geti(i).(*List)
}

func (z *List) getiInt(i int) *Int {
	// i は正しいと仮定
	return z.geti(i).(*Int)
}

func (z *List) geti(i int) GObj {
	return z.v[i]
}

func (z *List) Subst(xs RObj, lvs Level) *List {
	p := NewListN(z.Len())
	for i := 0; i < len(z.v); i++ {
		p.Append(gobjSubst(z.v[i], xs, lvs))
	}
	return p
}

func (z *List) toIntv(prec uint) *List {
	p := NewListN(z.Len())
	for i := 0; i < len(z.v); i++ {
		p.Append(gobjToIntv(z.v[i], prec))
	}
	return p
}
