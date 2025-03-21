package ganrac

import (
	"fmt"
	"strings"
)

// グローバル変数にしたくないのだけど，
// どうやって保持すべきなのか....
var varlist []varInfo

var varstr2lv map[string]Level

var coloredFml bool

type varInfo struct {
	v string
	p *Poly
}

func buildFormatString(f fmt.State, verb rune) string {
	var format strings.Builder

	format.WriteByte('%')
	for _, flag := range []byte{'-', '+', ' ', '#', '0'} {
		if f.Flag(int(flag)) {
			format.WriteByte(flag)
		}
	}

	if w, ok := f.Width(); ok {
		format.WriteString(fmt.Sprintf("%d", w))
	}

	if p, ok := f.Precision(); ok {
		format.WriteString(fmt.Sprintf(".%d", p))
	}

	format.WriteRune(verb)
	return format.String()
}

func VarStr(lv Level) string {
	if 0 <= lv && int(lv) < len(varlist) {
		return varlist[lv].v
	} else {
		return fmt.Sprintf("$%d", lv)
	}
}

func VarPoly(lv Level) *Poly {
	return varlist[lv].p
}

func VarNum() int {
	return len(varlist)
}

func VarStr2Lv(s string) (Level, bool) {
	lv, ok := varstr2lv[s]
	return lv, ok
}

func var2lv(v string) (Level, error) {
	for i, x := range varlist {
		if x.v == v {
			return Level(i), nil
		}
	}
	return 0, fmt.Errorf("undefined variable `%s`.", v)
}

func SetColordFml(b bool) {
	coloredFml = b
}

func GetColordFml() bool {
	return coloredFml
}

func esc_sgr(m int) string {
	if coloredFml {
		return fmt.Sprintf("\x1b[%dm", m)
	} else {
		return ""
	}
}

func init() {
	g := new(Ganrac)
	g.InitVarList([]string{
		"x", "y", "z", "w", "a", "b", "c", "d", "e", "f", "g", "h",
	})
}

func (g *Ganrac) InitVarList(vlist []string) error {
	for i, v := range vlist {
		if v == init_var_funcname {
			return fmt.Errorf("%s is reserved", v)
		}
		for _, bft := range g.builtin_func_table {
			if v == bft.name {
				return fmt.Errorf("%s is reserved", v)
			}
		}
		for j := 0; j < i; j++ {
			if vlist[j] == v {
				return fmt.Errorf("%s is duplicated", v)
			}
		}
	}

	varlist = make([]varInfo, len(vlist))
	varstr2lv = make(map[string]Level, len(vlist))
	for i := 0; i < len(vlist); i++ {
		varlist[i] = varInfo{vlist[i], NewPolyCoef(Level(i), 0, 1)}
		varstr2lv[vlist[i]] = Level(i)
	}

	return nil
}

func Add(x, y RObj) RObj {
	if x.Tag() >= y.Tag() {
		return x.Add(y)
	} else {
		return y.Add(x)
	}
}

func Sub(x, y RObj) RObj {
	if y.IsNumeric() || !x.IsNumeric() {
		return x.Sub(y)
	} else {
		// num - poly
		return y.Neg().Add(x)
	}
}

func AddOrSub(x, y RObj, sgn int) RObj {
	if sgn > 0 {
		return Add(x, y)
	} else {
		return Sub(x, y)
	}
}

func Mul(x, y RObj) RObj {
	if x.Tag() >= y.Tag() {
		return x.Mul(y)
	} else {
		return y.Mul(x)
	}
}
