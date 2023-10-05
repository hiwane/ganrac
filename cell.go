package ganrac

import (
	"fmt"
	"io"
	"math/big"
	"os"
)

type Cell struct {
	de           bool
	vanish       bool
	truth        int8
	sgn_of_left  sign_t
	lv           Level
	parent       *Cell
	children     []*Cell
	index        uint
	defpoly      *Poly
	intv         qIntrval  // 有理数=defpoly=nil か，bin-interval
	nintv        *Interval // 数値計算. defpoly=multivariate, de=true
	ex_deg       int       // 拡大次数.... のつもりだったけど使っていない
	signature    []sign_t
	multiplicity []mult_t
}

// for liftig phase
type cellStack struct {
	stack []*Cell
}

func newCellStack(cap int) *cellStack {
	cs := new(cellStack)
	cs.stack = make([]*Cell, 0, cap)
	return cs
}

func (cs *cellStack) empty() bool {
	return len(cs.stack) == 0
}

func (cs *cellStack) push(c *Cell) {
	cs.stack = append(cs.stack, c)
}

func (cs *cellStack) pop() *Cell {
	cell := cs.stack[len(cs.stack)-1]
	cs.stack = cs.stack[:len(cs.stack)-1]
	return cell
}

func (cs *cellStack) size() int {
	return len(cs.stack)
}

// signature が重複するセルがあるか.
func (cs *cellStack) hasSig(c *Cell) bool {
	for _, g := range cs.stack {
		if g.isSameSignature(c) {
			return true
		}
	}
	return false
}

// signature が重複しなければ cell を push する
func (cs *cellStack) pushSig(c *Cell) {
	if !cs.hasSig(c) {
		cs.push(c)
	}
}

// signature (projection factor の符号列)を返す.
// +      if positive
// -      if negative
// number otherwise. returns multiplicity
func (cell *Cell) printSignature(b io.Writer) {
	ch := '('
	str := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHJIKLMNOPQRSTUVWXYZ")
	strn := mult_t(len(str))
	for j := 0; j < len(cell.signature); j++ {
		sgns := '+'
		if cell.signature[j] < 0 {
			sgns = '-'
		} else if cell.signature[j] == 0 {
			// そもそも10次程度までしか動かないでしょう.
			if cell.multiplicity[j] < strn {
				sgns = str[cell.multiplicity[j]]
			} else {
				sgns = '?'
			}
		}
		fmt.Fprintf(b, "%c%c", ch, sgns)
		ch = ' '
	}
	fmt.Fprintf(b, ")")
}

func (cell *Cell) stringTruth() string {
	if cell.truth < 0 {
		return "?"
	}
	return []string{"f", "t", "."}[cell.truth]
}

func (cell *Cell) Truth() int8 {
	return cell.truth
}

func (cell *Cell) Print(args ...interface{}) error {
	return cell.Fprint(os.Stdout, args...)
}

func (cell *Cell) Fprint(b io.Writer, args ...interface{}) error {
	s := "cell"
	idx := 0
	var cad *CAD
	if len(args) > idx { // *CAD の取得
		if s2, ok := args[idx].(*CAD); ok {
			cad = s2
			idx++
		}
	}
	if len(args) > idx { // 表示方式の取得("sig", "signatures", "cell", ...)
		switch s2 := args[idx].(type) {
		case *String:
			s = s2.s
			idx++
		case string:
			s = s2
			idx++
		}
	}

	for i := idx; i < len(args); i++ { // index の取得
		ii, ok := args[i].(*Int)
		if !ok {
			return fmt.Errorf("invalid argument [expect integer]")
		}
		if ii.Sign() < 0 || !ii.IsInt64() || cell.children == nil || ii.Int64() >= int64(len(cell.children)) {
			return fmt.Errorf("invalid argument [invalid index]")
		}
		cell = cell.children[ii.Int64()]
	}

	switch s {
	case "sig", "signatures", "tcells", "fcells":
		if cell.children == nil {
			return fmt.Errorf("invalid argument [no child]")
		}
		truth := int8(-1)
		switch s {
		case "tcells":
			truth = t_true
		case "fcells":
			truth = t_false
		}

		fmt.Fprintf(b, "%s(%v) :: index=%v, truth=%d\n", s, args[1:], cell.Index(), cell.truth)
		if cad != nil {
			fmt.Fprintf(b, "         (")
			for i, pf := range cad.proj[cell.lv+1].gets() {
				if i != 0 {
					fmt.Fprintf(b, " ")
				}
				if pf.Input() {
					fmt.Fprintf(b, "i")
				} else {
					fmt.Fprintf(b, " ")
				}
			}
			fmt.Fprintf(b, ")\n")
		}

		for i, c := range cell.children {
			if truth >= 0 && c.truth != truth {
				continue
			}
			fmt.Fprintf(b, "%3d,%s,", i, c.stringTruth())
			if c.children == nil {
				fmt.Fprintf(b, "  ")
			} else {
				fmt.Fprintf(b, "%2d", len(c.children))
			}
			fmt.Fprintf(b, " ")
			c.printSignature(b)
			if c.intv.inf != nil {
				fmt.Fprintf(b, " [% e,% e]", c.intv.inf.Float(), c.intv.sup.Float())
			} else if c.nintv != nil {
				fmt.Fprintf(b, " [% e,% e]", c.nintv.inf, c.nintv.sup)
			}
			if c.defpoly != nil {
				fmt.Fprintf(b, " %.50v", c.defpoly)
			} else if c.isSection() {
				if c.intv.inf != c.intv.sup {
					panic("invlaid")
				}
				fmt.Fprintf(b, " %v", c.intv.inf)
			}
			fmt.Fprintf(b, "\n")
		}
	case "cell":
		fmt.Fprintf(b, "--- information about the cell %v %p ---\n", cell.Index(), cell)
		fmt.Fprintf(b, "lv=%d:%s, de=%v, exdeg=%d, truth=%d sgn=%d\n",
			cell.lv, VarStr(cell.lv), cell.de, cell.ex_deg, cell.truth, cell.sgn_of_left)
		var num int
		if cell.children == nil {
			num = -1
		} else {
			num = len(cell.children)
		}
		fmt.Fprintf(b, "# of children=%d\n", num)
		if cell.defpoly != nil {
			fmt.Fprintf(b, "def.poly     =%v\n", cell.defpoly)
		} else if cell.isSection() {
			if cell.intv.inf != cell.intv.sup {
				panic("invlaid")
			}
			fmt.Fprintf(b, "def.value    =%v\n", cell.intv.inf)
		}
		if cell.signature != nil {
			fmt.Fprintf(b, "signature    =")
			cell.printSignature(b)
			fmt.Fprintf(b, "\n")
		}
		if cell.intv.inf != nil {
			sup := cell.intv.sup.Float()
			inf := cell.intv.inf.Float()
			dist := sup - inf
			fmt.Fprintf(b, "iso.intv     =[%v,%v]\n", cell.intv.inf, cell.intv.sup)
			fmt.Fprintf(b, "             =[%e,%e].  dist=%e\n", inf, sup, dist)
		}
		if cell.nintv != nil {
			bb := new(big.Float)
			bb.Sub(cell.nintv.sup, cell.nintv.inf)
			fmt.Fprintf(b, "iso.nintv    =%f\n", cell.nintv)
			fmt.Fprintf(b, "             =%e.  dist=%e\n", cell.nintv, bb)
		}
	case "cellp":
		for cell.lv >= 0 {
			if err := cell.Fprint(b); err != nil {
				return err
			}
			cell = cell.parent
			fmt.Fprintf(b, "cell %d: %v\n", cell.lv, cell.Index())
		}
		return cell.Fprint(b)
	default:
		return fmt.Errorf("invalid argument [kind=%s]", s)
	}
	return nil
}

func (cell *Cell) Neg() *Cell {
	c := new(Cell)
	*c = *cell
	c.nintv = newInterval(cell.nintv.Prec())
	c.nintv.inf.Neg(cell.nintv.sup)
	c.nintv.sup.Neg(cell.nintv.inf)
	return c
}

// returns the sign of the sample point
func (cell *Cell) Sign() int {
	if cell.intv.inf != nil {
		if cell.intv.inf.Sign() > 0 {
			return 1
		} else if cell.intv.sup.Sign() < 0 {
			return -1
		}
	} else {
		if cell.nintv.inf.Sign() > 0 {
			return 1
		} else if cell.nintv.sup.Sign() < 0 {
			return -1
		}
	}
	return 0
}

// cell の妥当性評価 for debug
func (cell *Cell) valid(cad *CAD) error {
	if cell.lv >= 0 {
		if cell.defpoly == nil {
			if cell.intv.inf != cell.intv.sup {
				return fmt.Errorf("defpoly=nil but inf != sup %v", cell.Index())
			}
		}
	}

	idx := cell.Index()
	if len(idx) > 0 && (cad.stage < 2 || cad.q[cell.lv] == q_free) {
		c := cad.root
		for _, x := range idx {
			c = c.children[x]
		}
		if c != cell {
			return fmt.Errorf("cell=%v... but ...", idx)
		}
	}

	if cell.children != nil {
		nt := 0
		nf := 0
		ntc := 0
		nfc := 0
		children := make([]*Cell, 0, len(cell.children))
		for i, c := range children {
			if c.index != uint(i) {
				return fmt.Errorf("index invalid. expect=%d, actual=%v", i, c.Index())
			}
		}
		for i := 0; i < len(cell.children); i += 2 {
			children = append(children, cell.children[i])
		}
		for i := 1; i < len(cell.children); i += 2 {
			children = append(children, cell.children[i])
		}
		for _, c := range children {
			if err := c.valid(cad); err != nil {
				return err
			}
			if c.truth == t_true {
				nt++
				if c.children != nil {
					ntc++
				}

			} else if c.truth == t_false {
				nf++
				if c.children != nil {
					nfc++
				}
			}
			if c.parent != cell && (cad.stage < 2 || cell.truth < 0) {
				return fmt.Errorf("%v.parent=%v != %v invalid [%p:%p]", c.Index(), c.parent.Index(), cell.Index(), c.parent, cell)
			}
		}

		errmes := ""
		if cad.q[cell.lv+1] == q_forall {
			if nfc > 1 {
				errmes = fmt.Sprintf("forall + many false cells")
			} else if nf > 0 && cell.truth != t_false {
				errmes = fmt.Sprintf("forall + false false false")
			} else if nf == 0 && nt > 0 && cell.truth != t_true {
				errmes = fmt.Sprintf("forall + true true true")
			}
		} else if cad.q[cell.lv+1] == q_exists {
			if ntc > 1 {
				errmes = fmt.Sprintf("exists + many false cells")
			} else if nt > 0 && cell.truth != t_true {
				errmes = fmt.Sprintf("exists + true true true")
			} else if nt == 0 && nf > 0 && cell.truth != t_false {
				errmes = fmt.Sprintf("exists + false false false")
			}
		}
		if errmes != "" {
			cell.Print("cellp")
			cell.Print("signatures")
			return fmt.Errorf("%s: index=%v, truth=%d, true=(%d/%d), false=(%d/%d)", errmes, cell.Index(), cell.truth, ntc, nt, nfc, nf)
		}
	}

	return nil
}

func (cell *Cell) Precs() []uint {
	prec := make([]uint, cell.lv+1)
	c := cell
	for c.lv >= 0 {
		prec[c.lv] = c.Prec()
		c = c.parent
	}
	return prec
}

func (cell *Cell) Prec() uint {
	if cell.defpoly == nil {
		return 0
	}
	switch l := cell.intv.inf.(type) {
	case *BinInt:
		if l.m > 0 {
			return 0
		} else {
			return uint(-l.m)
		}
	}
	return cell.nintv.Prec()
}

// インデックスを返す. 論文等は 1 始まりだが，0 始まりであることに注意
// つまり， section は奇数である
func (cell *Cell) Index() []uint {
	idx := make([]uint, cell.lv+1)
	c := cell
	for c.lv >= 0 {
		idx[c.lv] = c.index
		c = c.parent
	}
	return idx
}

// 自分〜先祖に section がいるか.
func (cell *Cell) hasSection() bool {
	for c := cell; c.lv >= 0; c = c.parent {
		if c.isSection() {
			return true
		}
	}
	return false
}

func (cell *Cell) isSector() bool {
	return cell.index%2 == 0
}

func (cell *Cell) isSection() bool {
	return cell.index%2 != 0
}

func (cell *Cell) setDefPoly(p *Poly) {
	cell.defpoly = p
	cell.ex_deg = len(p.c) - 1
}

func (c *Cell) isSameSignature(d *Cell) bool {
	if c.lv != d.lv {
		panic(fmt.Sprintf("not same level c=%v, d=%v", c.Index(), d.Index()))
	}
	if len(c.signature) != len(d.signature) {
		panic(fmt.Sprintf("not same len c[%d]=%v, d[%d]=%v", len(c.signature), c.Index(), len(d.signature), d.Index()))
	}
	if c.isSection() != d.isSection() {
		return false
	}

	for i := 0; i < len(c.signature); i++ {
		if c.signature[i] != d.signature[i] {
			return false
		}
	}
	return true
}
