package ganrac

// George E. Collins
// Quantifier elimination for real closed fields by cylindrical algebraic decomposition

import (
	"fmt"
	"io"
	"os"
	"time"
)

type ProjectionAlgo int

type sign_t int8
type mult_t int8

const (
	t_undef  = -1 // まだ評価していない
	t_false  = 0
	t_true   = 1
	t_other  = 2 // 兄弟の情報で親の真偽値が確定したのでもう評価しない
	q_free   = t_undef
	q_forall = t_false
	q_exists = t_true

	PF_EVAL_UNKNOWN = -1
	PF_EVAL_NO      = 0
	PF_EVAL_YES     = 1

	PROJ_McCallum = 0
	PROJ_HONG     = 1

	CAD_STAGE_INITED int8 = 0
	CAD_STAGE_PROJED int8 = 1
	CAD_STAGE_LIFTED int8 = 2
)

type qIntrval struct {
	// int, rat, binint が入ると想定
	inf, sup NObj
}

type AtomProj struct {
	Atom
	pl *ProjLink
}

type CADStat struct {
	fctr         int
	qrealroot    int
	irealroot    int
	irealroot_ok int
	sqrt         int
	sqrt_ok      int
	discriminant int
	resultant    int
	psc          int
	cell         []int
	true_cell    []int
	false_cell   []int
	precision    int
	lift         []int
	rlift        []int
	tm           []time.Duration
	sector       int
}

type CAD struct {
	qfml     Fof           // quantified formula: input
	fml      Fof           // qff
	output   Fof           // qff
	q        []int8        // quantifier
	proj     []ProjFactors // [level]
	u        []*Interval   // [level]
	pl4const []*ProjLink   // 定数用 0, +, -
	apppoly  []*Poly       // makepdf 用. 入力以外の多項式
	stack    *cellStack
	root     *Cell
	rootp    *Cellmod
	g        *Ganrac
	stat     CADStat
	nwo      bool // well-oriented
	stage    int8
	palgo    ProjectionAlgo
}

type CAD_error_wo struct {
	cell *Cell
}

func (e *CAD_error_wo) Error() string {
	return fmt.Sprintf("NOT well-oriented: %v", e.cell.Index())
}

func NewCadErrorWO(cell *Cell) error {
	e := new(CAD_error_wo)
	e.cell = cell
	return e
}

//var CAD_NO_WO = errors.New("NOT well-oriented")

func qeCAD(fml Fof) Fof {
	return fml
}

func (stat CADStat) Fprint(b io.Writer, cad *CAD) {
	fmt.Fprintf(b, "time....\n")
	fmt.Fprintf(b, "=========================\n")
	fmt.Fprintf(b, "proj: %9.3f sec\n", stat.tm[0].Seconds())
	fmt.Fprintf(b, "lift: %9.3f sec\n", stat.tm[1].Seconds())
	fmt.Fprintf(b, "sfc : %9.3f sec\n", stat.tm[2].Seconds())
	fmt.Fprintf(b, "\n")

	if cad.stage >= CAD_STAGE_PROJED {
		fmt.Fprintf(b, "CAD proj. stat....\n")
		fmt.Fprintf(b, "===============================\n")
		fmt.Fprintf(b, "LV | # proj | deg | tdeg | npr\n")
		for lv := len(cad.q) - 1; lv >= 0; lv-- {
			deg := 0
			tdeg := 0
			npr := 0
			for _, p := range cad.proj[lv].gets() {
				d := p.P().deg()
				if d > deg {
					deg = d
				}
				if p.Sign() != 0 {
					npr++
				}
			}

			fmt.Fprintf(b, "%2d |%7d |%4d |%5d |%4d\n", lv, cad.proj[lv].Len(), deg, tdeg, npr)
		}
		fmt.Fprintf(b, "\n")
	}
	if cad.stage >= CAD_STAGE_LIFTED {
		fmt.Fprintf(b, "CAD cell stat....\n")
		fmt.Fprintf(b, "=====================================================\n")
		fmt.Fprintf(b, "LV |    cell |    true |   false |    lift |   rlift\n")
		fmt.Fprintf(b, "---+---------+---------+---------+---------+---------\n")
		sn := make([]int, 5)
		for i := 0; i < len(cad.q); i++ {
			fmt.Fprintf(b, "%2d |%8d |%8d |%8d |%8d |%8d\n", i, stat.cell[i], stat.true_cell[i], stat.false_cell[i], stat.lift[i], stat.rlift[i])
			sn[0] += stat.cell[i]
			sn[1] += stat.true_cell[i]
			sn[2] += stat.false_cell[i]
			sn[3] += stat.lift[i]
			sn[4] += stat.rlift[i]
		}
		fmt.Fprintf(b, "---+---------+---------+---------+---------+---------\n")
		fmt.Fprintf(b, "%2d |%8d |%8d |%8d |%8d |%8d\n", -1, sn[0], sn[1], sn[2], sn[3], sn[4])
		fmt.Fprintf(b, "\n")
	}
	fmt.Fprintf(b, "CA stat....\n")
	fmt.Fprintf(b, "==================================================\n")
	if cad.stage >= CAD_STAGE_PROJED {
		fmt.Fprintf(b, " - proj | discrim              : %8d\n", stat.discriminant)
		fmt.Fprintf(b, " - proj | resultant            : %8d\n", stat.resultant)
		fmt.Fprintf(b, " - proj | psc                  : %8d\n", stat.psc)
		fmt.Fprintf(b, " - proj | factorization over Z : %8d\n", stat.fctr)
	}
	if cad.stage >= CAD_STAGE_LIFTED {
		fmt.Fprintf(b, " - lift | real root in Z[x]    : %8d\n", stat.qrealroot)
		fmt.Fprintf(b, " - lift | real root in intv[x] : %8d / %d\n", stat.irealroot_ok, stat.irealroot)
		fmt.Fprintf(b, " - lift | sector sample point  : %8d\n", stat.sector)
	}
}

func (c *CAD) Tag() uint {
	return TAG_CAD
}

func (c *CAD) String() string {
	return fmt.Sprintf("CAD[%v]", c.fml)
}

func (c *CAD) log(lv int, format string, a ...interface{}) {
	if lv <= c.g.verbose_cad {
		c.g.logger.Printf(format, a...)
	}
}

func (c *CAD) InitRoot() {
	c.root = NewCell(c, nil, 0)
	c.rootp = NewCellmod(c.root)
}

func NewCAD(prenex_formula Fof, g *Ganrac) (*CAD, error) {
	if err := prenex_formula.valid(); err != nil {
		return nil, err
	}
	if g.ox == nil {
		return nil, fmt.Errorf("CAS is required")
	}
	switch prenex_formula.(type) {
	case *AtomT, *AtomF:
		return nil, fmt.Errorf("prenex formula is expected")
	}

	c := new(CAD)
	c.g = g
	c.stage = CAD_STAGE_INITED

	///////////////////////////////////
	// 変数順序の妥当性チェック
	///////////////////////////////////
	v := prenex_formula.maxVar()
	c.q = make([]int8, v)
	for i := 0; i < len(c.q); i++ {
		c.q[i] = -1
	}
	c.qfml = prenex_formula
	c.fml = prenex_formula
	vmax := Level(0)
	for cnt := 0; ; cnt++ {
		var qq []Level
		var qval int8
		switch f := c.fml.(type) {
		case *ForAll:
			qq = f.q
			qval = q_forall
			c.fml = f.fml
		case *Exists:
			qq = f.q
			qval = q_exists
			c.fml = f.fml
		default:
			goto _NEXT
		}

		max := vmax
		min := v
		for _, qi := range qq {
			c.q[qi] = qval
			if min > qi {
				min = qi
			}
			if max < qi {
				max = qi
			}
		}
		if int(max-min) != len(qq)-1 || (cnt > 0 && min != vmax+1) {
			return nil, fmt.Errorf("CAD: invalid variable order [%d,%d,%d]", min, max, vmax)
		}

		vmax = max
	}
_NEXT:

	if !c.fml.IsQff() {
		return nil, fmt.Errorf("prenex formula is expected")
	}
	// 隙間があると面倒なのでエラーにする
	for i := Level(0); int(i) < len(c.q); i++ {
		if !c.fml.hasVar(i) {
			return nil, fmt.Errorf("CAD: invalid variable order [%d,%d]", i, vmax)
		}
	}

	qdx := false
	for i := Level(0); int(i) < len(c.q); i++ {
		if c.q[i] >= 0 {
			qdx = true
		} else if qdx {
			return nil, fmt.Errorf("CAD: invalid variable order [%d,%d]", i, vmax)
		}
	}

	c.InitRoot()
	c.stack = newCellStack()
	c.stack.push(c.root)
	c.stat.cell = make([]int, len(c.q))
	c.stat.true_cell = make([]int, len(c.q))
	c.stat.false_cell = make([]int, len(c.q))
	c.stat.lift = make([]int, len(c.q))
	c.stat.rlift = make([]int, len(c.q))
	c.stat.tm = make([]time.Duration, 3)

	return c, nil
}

func (c *CAD) InitProj(algo ProjectionAlgo) {
	vnum := Level(len(c.q))
	c.proj = make([]ProjFactors, vnum)

	for i := Level(0); i < vnum; i++ {
		if algo == PROJ_McCallum {
			c.proj[i] = newProjFactorsMC()
		} else if algo == PROJ_HONG {
			c.proj[i] = newProjFactorsHH()
		} else {
			panic(fmt.Sprintf("unknown %v", algo))
		}
	}

	c.fml = clone4CAD(c.fml, c)

	// 定数（符号確定）用の ProjLink を構築
	c.pl4const = make([]*ProjLink, 3)
	for i, s := range []OP{EQ, GT, LT} {
		c.pl4const[i] = newProjLink()
		c.pl4const[i].op = s
	}
}

func clone4CAD(formula Fof, c *CAD) Fof {
	switch fml := formula.(type) {
	case *FmlAnd:
		var t Fof = trueObj
		for _, f := range fml.fml {
			t = NewFmlAnd(t, clone4CAD(f, c))
		}
		return t
	case *FmlOr:
		var t Fof = falseObj
		for _, f := range fml.fml {
			t = NewFmlOr(t, clone4CAD(f, c))
		}
		return t
	case *Atom:
		t := new(AtomProj)
		t.op = fml.op
		t.p = make([]*Poly, len(fml.p))
		t.pl = new(ProjLink)
		t.pl.op = GT
		for i, poly := range fml.p {
			pl2 := c.addPoly(poly, true)
			t.pl.merge(pl2)
			t.p[i] = poly
		}
		return t
	case *AtomT, *AtomF:
		return fml
	}
	fmt.Printf("fml=%v\n", formula)
	panic("stop")
}

func NewCell(cad *CAD, parent *Cell, idx uint) *Cell {
	cell := new(Cell)
	cell.parent = parent
	cell.index = idx
	cell.ex_deg = 1
	cell.truth = t_undef
	if parent == nil {
		cell.lv = -1
	} else {
		cell.lv = parent.lv + 1
		cell.signature = make([]sign_t, cad.proj[cell.lv].Len())
		cell.multiplicity = make([]mult_t, cad.proj[cell.lv].Len())
	}
	return cell
}

func (c *CAD) Root() *Cell {
	return c.root
}

func (cad *CAD) Print(args ...interface{}) error {
	return cad.Fprint(os.Stdout, args...)
}

func (cad *CAD) Fprint(b io.Writer, args ...interface{}) error {
	if len(args) == 0 {
		fmt.Fprintf(b, "input=%v\n", cad.fml)
		return nil
	}
	s, ok := args[0].(*String)
	if !ok {
		return fmt.Errorf("invalid argument [expect string]")
	}

	switch s.s {
	case "stat":
		cad.stat.Fprint(b, cad)
	case "proj":
		return cad.FprintProj(b, args[1:]...)
	case "proji":
		cad.FprintInput(b, args[1:]...)
	case "sig", "signatures", "cell", "cellp", "fcells", "tcells":
		aa := make([]interface{}, len(args)+1)
		copy(aa[1:], args)
		aa[0] = cad
		return cad.root.Fprint(b, aa...)
	default:
		return fmt.Errorf("invalid argument")
	}

	return nil
}
