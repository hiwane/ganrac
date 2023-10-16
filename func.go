package ganrac

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func (g *Ganrac) FunctionNames() []string {
	ret := make([]string, 0, len(g.builtin_func_table))
	for _, s := range g.builtin_func_table {
		if g.ox != nil || !s.ox {
			ret = append(ret, s.name)
		}
	}
	return ret
}

func (g *Ganrac) setBuiltinFuncTable() {
	// 関数テーブル
	g.builtin_func_table = []func_table{
		// sorted by name
		{"all", 2, 2, funcForAll, false, "([x], FOF)", "universal quantifier", nil, "", "", ""},
		//		{"and", 2, 2, funcAnd, false, "(FOF, ...):\t\tconjunction (&&)", ""},
		{"cad", 1, 2, funcCAD, true, "(fof [, option])", "", []func_help{
			{"fof", "first-order formula or example name", true, "", nil, ""},
			{"option", "dictionary", false, "", []func_help{
				{"proj", "(m|h) projection operator", false, "m", nil, ""},
				{"var", "(0|1) variable order.         1 if auto", false, "0", nil, ""},
				{"stat", "(0|1) print statistics data.  1 if print", false, "0", nil, ""},
				{"ls", "(0|1) lifting strategy.", false, "1", []func_help{
					{"0", "basic", false, "", nil, ""},
					{"1", "improved", false, "", nil, ""},
				}, ""},
			}, ""},
		}, "", "", ""},
		{"cadinit", 1, 2, funcCADinit, true, "(fof [, option])", "", nil, "", "", ""},
		{"cadlift", 1, 10, funcCADlift, true, "(CAD [, index, ...])", "", nil, "", "", ""},
		{"cadproj", 1, 2, funcCADproj, true, "(CAD [, proj])", "", nil, "", "", ""},
		{"cadsfc", 1, 1, funcCADsfc, true, "(CAD)", "", nil, "", "", ""},
		{"coef", 3, 3, funcCoef, false, "(poly, var, deg)", "", nil, "", "", ""}, // coef(F, x, 2)
		{"counter", 0, 0, funcCounter, false, "()", "", nil, "", "", ""},
		{"deg", 2, 2, funcDeg, false, "(poly|fof, var)", "degree of a polynomial with respect to var", []func_help{
			{"poly", "polynomial", true, "", nil, ""},
			{"fof", "first-order formula", true, "", nil, ""},
			{"var", "variable", true, "", nil, ""},
		}, "degree of a polynomial w.r.t a variable", `
> deg(x^2+3*y+1, x);
2
> deg(x^2+3*y+1, y);
1
> deg(x^2+3*y+1>0 && x^3+y^3==0, y);
3
> deg(0, y);
-1
`, ""}, // deg(F, x)
		{"diff", 2, 101, funcDiff, false, "(poly, var, ...)", "differential", nil, "", "", ""},
		{"discrim", 2, 2, funcOXDiscrim, true, "(poly, var)", "discriminant", []func_help{
			{"poly", "polynomial", true, "", nil, ""},
			{"var", "variable", true, "", nil, ""},
		}, "", `
> discrim(2*x^2-3*x-3, x);
33
> discrim(a*x^2+b*x+c, x);
-4*c*a+b^2
> discrim(a*x^2+b*x+c, y);
0
`, ""},
		{"equiv", 2, 2, funcEquiv, false, "(fof1, fof2)", "fof1 is equivalent to fof2", nil, "", "", ""},
		{"evalcas", 1, 1, funcOXStr, true, "(str)", "evaluate str by CAS", []func_help{
			{"str", "string", true, "", nil, ""},
		}, "", `
> evalcas("fctr(x^2-4);");
[[1,1],[x-2,1],[x+2,1]]
`, ""},
		{"ex", 2, 2, funcExists, false, "(vars, FOF)", "existential quantifier", []func_help{
			{"vars", "list of variables", true, "", nil, ""},
			{"FOF", "first-order formula", true, "", nil, ""},
		}, "", `
> ex([x], a*x^2+b*x+c == 0):
`, ""},
		{"example", 0, 1, funcExample, false, "([name])", "example", nil, "", "", ""},
		{"fctr", 1, 1, funcOXFctr, true, "(poly)", "factorize polynomial over the rationals", nil, "", "", ""},
		{"fctrp", 1, 1, funcOXFctrp, true, "(poly)", "factorize polynomial over the rationals. para", nil, "", "", ""},
		{"gb", 2, 3, funcOXGB, true, "(polys, vars)", "Groebner basis", nil, "", "", ""},
		{"help", 0, 1, nil, false, "([str])", "show help", nil, "", "", ""},
		//		{"igcd", 2, 2, funcIGCD, false, "(int1, int2)", "The integer greatest common divisor", ""},
		{"impl", 2, 2, funcImpl, false, "(fof1, fof2)", "fof1 impies fof2", nil, "", "", ""},
		{"indets", 1, 1, funcIndets, false, "(mobj)", "find indeterminates of an expression", nil, "", "", ""},
		{"intv", 1, 3, funcIntv, false, "(lb, ub [, prec])", "make an interval", nil, "", "", ""},
		{"len", 1, 1, funcLen, false, "(mobj)", "length of an object", nil, "", "", ""},
		{"load", 2, 2, funcLoad, false, "(fname)@", "load file", nil, "", "", ""},
		{"not", 1, 1, funcNot, false, "(fof)", "", []func_help{
			{"fof", "first-order formula", true, "", nil, ""},
		}, "", `
> not(x > 0);
x <= 0
> not(ex([x], a*x^2+b*x+c==0));
all([x], a*x^2+b*x+c != 0)
`, ""},
		/*
		   		{"oxfunc", 2, 100, funcOXFunc, true, "(fname, args...)*\tcall ox-function by ox-asir", `
		   Args
		   ========
		   fname : string, function name of ox-server
		   args  : arguments of the function

		   Examples
		   ========
		     > oxfunc("deg", x^2-1, x);
		     2
		     > oxfunc("igcd", 8, 12);
		     4
		   `},
		*/
		{"print", 1, 10, funcPrint, false, "(obj [, kind, ...])", "print object", nil, "", `
> print(x^10+y-3);
y+x^10-3
> print(x^10+y-3, "tex");
y+x^{10}-3
> print(ex([x], x^2>1 && y +x == 0), "tex");
\exists x x^2-1 > 0 \land y+x = 0

> ` + init_var_funcname + `(a,b,c,x);
> C = cadinit(ex([x], a*x^2+b*x+c==0));
> cadproj(C);
> print(C, "proj");
> print(C, "proji");
> cadlift(C);
> print(C, "signatures");
> print(C, "sig", 1);  # "sig" is an abbreviation for "signatures"
> print(C, "cell", 1, 1);
> print(C, "stat");
`, ""},
		{"psc", 4, 4, funcOXPsc, true, "(poly, poly, var, int)", "principal subresultant coefficient", nil, "", "", ""},
		{"qe", 1, 2, funcQE, true, "(fof [, opt])", "real quantifier elimination", []func_help{
			{"fof", "first-order formula", true, "", nil, ""},
			{"opt", "dictionary", false, "",
				[]func_help{
					{getQEoptStr(QEALGO_VSLIN), "(0|1) linear       virtual substitution", false, "1", nil, ""},
					{getQEoptStr(QEALGO_VSQUAD), "(0|1) quadratic    virtual substitution", false, "1", nil, ""},
					{getQEoptStr(QEALGO_EQLIN), "(0|1) linear       equational constraint [Hong93]", false, "1", nil, ""},
					{getQEoptStr(QEALGO_EQQUAD), "(0|1) quadratic    equational constraint [Hong93]", false, "1", nil, ""},
					{getQEoptStr(QEALGO_NEQ), "(0|1) inequational constraints [Iwane15]", false, "1", nil, ""},
					{getQEoptStr(QEALGO_SMPL_EVEN), "(0|1) simplify     even formula", false, "1", nil, ""},
					{getQEoptStr(QEALGO_SMPL_HOMO), "(0|1) simplify     homogeneous formula", false, "1", nil, ""},
					{getQEoptStr(QEALGO_SMPL_TRAN), "(0|1) simplify     translatation formula", false, "1", nil, ""},
				}, ""},
		}, "", fmt.Sprintf(`
> vars(x, b, c):
> F = ex([x], x^2+b*x+c == 0, {%s: 0, %s: 1});
`, getQEoptStr(QEALGO_VSQUAD), getQEoptStr(QEALGO_SMPL_HOMO)), ""},
		{"quit", 0, 1, funcQuit, false, "([code])", "bye", nil, "", "", ""},
		{"realroot", 2, 2, funcRealRoot, false, "(uni-poly)", "real root isolation", nil, "", "", ""},
		{"res", 3, 3, funcOXRes, true, "(poly, poly, var)", "resultant", nil, "", "", ""},
		{"rootbound", 1, 1, funcRootBound, false, "(unipoly)", "root bound",
			[]func_help{
				{"unipoly", "univariate polynomial", true, "", nil, ""},
			}, "", `
> rootbound(x^2-2);
3
`, ""},
		{"save", 2, 3, funcSave, false, "(obj, fname)@", "save object...", nil, "", "", ""},
		{"simpl", 1, 3, funcSimplify, true, "(Fof [, neccon [, sufcon]])", "simplify formula FoF", nil, "", "", ""},
		{"simplnum", 1, 1, funcSimplNum, true, "(Fof [, neccon [, sufcon]])", "simplify formula FoF for DEBUG", nil, "", "", ""},
		{"sleep", 1, 1, funcSleep, false, "(milisecond)", "zzz", nil, "", "", ""},
		// {"sqfr", 1, 1, funcSqfr, false, "(poly)* square-free factorization", nil, "", "", ""},
		{"slope", 4, 4, funcOXSlope, true, "(poly, poly, var, int)", "slope resultant", nil, "", "", ""},
		{"sres", 3, 4, funcOXSres, true, "(poly, poly, var[, int])", "subresultant sequence", nil, "", "", ""},
		{"subst", 1, 101, funcSubst, false, "(poly|FOF|List,x,vx,y,vy,...)", "", nil, "", "", ""},
		{"time", 1, 1, funcTime, false, "(expr)", "run command and system resource usage", nil, "", "", ""},
		{init_var_funcname, 0, 0, nil, false, "(var, ...)", "init variable order",
			[]func_help{
				{"var", "variable", true, "", nil, ""},
			}, "", `
> ` + init_var_funcname + `(x,y,z);
> F = x^2+2;
> F;
> ` + init_var_funcname + `(a,b,c);
> F;
a^2+2
> x;
error: undefined variable ` + "`x`\n", ""},
		{"verbose", 1, 2, funcVerbose, false, "(int [, int])", "set verbose level", nil, "", "", ""},
		{"vs", 1, 2, funcVS, true, "(FOF, int) ", "virtual substitution", nil, "", "", ""},
	}
}

// func (p *pNode) callFunction(args []interface{}) (interface{}, error) {
func (g *Ganrac) callFunction(funcname string, args []interface{}) (interface{}, error) {
	// とりあえず素朴に
	for _, f := range g.builtin_func_table {
		if f.name == funcname {
			if len(args) < f.min {
				return nil, fmt.Errorf("too few argument: function %s()", funcname)
			}
			if len(args) > f.max {
				return nil, fmt.Errorf("too many argument: function %s()", funcname)
			}
			if f.ox && g.ox == nil {
				return nil, fmt.Errorf("required OX server: function %s()", funcname)
			}
			if f.name == "help" {
				return funcHelp(g.builtin_func_table, f.name, args)
			} else {
				return f.f(g, f.name, args)
			}
		}
	}

	return nil, fmt.Errorf("unknown function: %s", funcname)
}

// //////////////////////////////////////////////////////////
// 論理式
// //////////////////////////////////////////////////////////
func funcNot(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	f, ok := args[0].(Fof)
	if !ok {
		return nil, fmt.Errorf("not(): unsupported for %v", args[0])
	}
	return f.Not(), nil
}

func funcImpl(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	f0, ok := args[0].(Fof)
	if !ok {
		return nil, fmt.Errorf("%s(1st-arg): expected a first-order formula", name)
	}
	f1, ok := args[1].(Fof)
	if !ok {
		return nil, fmt.Errorf("%s(2nd-arg): expected a first-order formula", name)
	}

	return NewFmlImpl(f0, f1), nil
}

func funcEquiv(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	f0, ok := args[0].(Fof)
	if !ok {
		return nil, fmt.Errorf("%s(1st-arg): expected a first-order formula", name)
	}
	f1, ok := args[1].(Fof)
	if !ok {
		return nil, fmt.Errorf("%s(2nd-arg): expected a first-order formula", name)
	}

	return NewFmlEquiv(f0, f1), nil
}

func funcExists(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	return funcForEx(false, name, args)
}

func funcForAll(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	return funcForEx(true, name, args)
}

func funcForEx(forex bool, name string, args []interface{}) (interface{}, error) {
	f0, ok := args[0].(*List)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected list: %d:%v", name, args[0].(GObj).Tag(), args[0])
	}

	lv := make([]Level, len(f0.v))
	for i, qq := range f0.v {
		q, ok := qq.(*Poly)
		if !ok || !q.isVar() {
			return nil, fmt.Errorf("%s(1st arg:%d): expected var-list", name, i)
		}
		lv[i] = q.lv
	}

	f1, ok := args[1].(Fof)
	if !ok {
		return nil, fmt.Errorf("%s(2nd arg): expected formula", name)
	}
	return NewQuantifier(forex, lv, f1), nil
}

// //////////////////////////////////////////////////////////
// OpenXM
// //////////////////////////////////////////////////////////
func funcOXStr(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	f0, ok := args[0].(*String)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected string: %d:%v", name, args[0].(GObj).Tag(), args[0])
	}

	return g.ox.Eval(f0.s)
}

/*
func funcOXFunc(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	f0, ok := args[0].(*String)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected string: %d:%v", name, args[0].(GObj).Tag(), args[0])
	}

	err := g.ox.ExecFunction(f0.s, args[1:]...)
	if err != nil {
		return nil, fmt.Errorf("%s(): required OX server", name)
	}
	s, err := g.ox.PopCMO()
	if err != nil {
		return nil, fmt.Errorf("%s(): popCMO failed %w", name, err)
	}
	gob := g.ox.toGObj(s)

	return gob, nil
}
*/

func funcDiff(g *Ganrac, name string, args []interface{}) (interface{}, error) {

	for i := 1; i < len(args); i++ {
		c, ok := args[i].(*Poly)
		if !ok || !c.isVar() {
			return nil, fmt.Errorf("%s(%dth arg): expected var: %v", name, i, args[i])
		}
	}

	_, ok := args[0].(NObj)
	if ok {
		return zero, nil
	}
	p, ok := args[0].(*Poly)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected polynomial: %v", name, args[0])
	}

	var ret RObj = p
	for i := 1; i < len(args); i++ {
		c := args[i].(*Poly)
		switch r := ret.(type) {
		case *Poly:
			ret = r.Diff(c.lv)
		default:
			return zero, nil
		}
	}
	return ret, nil
}

func funcOXDiscrim(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	c, ok := args[1].(*Poly)
	if !ok || !c.isVar() {
		return nil, fmt.Errorf("%s(2nd arg): expected var: %v", name, args[1])
	}

	switch p := args[0].(type) {
	case *Poly:
		return g.ox.Discrim(p, c.lv), nil
	case NObj:
		return zero, nil
	default:
		return nil, fmt.Errorf("%s(1st arg): expected poly: %d:%v", name, args[0].(GObj).Tag(), args[0])
	}
}

func funcOXFctr(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	fmt.Printf("go funcOXFctr\n")
	f0, ok := args[0].(*Poly)
	fmt.Printf("f0 =%v\n", f0)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected poly: %d:%v", name, args[0].(GObj).Tag(), args[0])
	}

	fmt.Printf("go oxFctor\n")
	return g.ox.Factor(f0), nil
}

func funcOXFctrp(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	fmt.Printf("go funcOXFctrpara\n")
	f0, ok := args[0].(*Poly)
	fmt.Printf("f0 =%v\n", f0)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected poly: %d:%v", name, args[0].(GObj).Tag(), args[0])
	}

	ch := make(chan any)
	go func(p *Poly, ch chan any) {
		ch <- g.ox.Factor(p)
	}(f0, ch)

	q := <-ch
	fmt.Printf("go oxFctor\n")
	return q, nil
}

func funcOXGB(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	f0, ok := args[0].(*List)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected poly-list: %d:%v", name, args[0].(GObj).Tag(), args[0])
	}

	f1, ok := args[1].(*List)
	if !ok {
		return nil, fmt.Errorf("%s(2nd arg): expected var-list: %d:%v", name, args[1].(GObj).Tag(), args[1])
	}

	n := 0
	if len(args) == 3 {
		f2, ok := args[2].(*Int)
		if !ok || f2.Sign() < 0 || !f2.IsInt64() {
			return nil, fmt.Errorf("%s(3rd arg): expected nonnegint: %d:%v", name, args[2].(GObj).Tag(), args[2])
		}
		n = int(f2.Int64())
	}

	return g.ox.GB(f0, f1, n), nil
}

func funcOXPsc(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	f, ok := args[0].(*Poly)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected poly: %d:%v", name, args[0].(GObj).Tag(), args[0])
	}
	h, ok := args[1].(*Poly)
	if !ok {
		return nil, fmt.Errorf("%s(2nd arg): expected poly: %d:%v", name, args[1].(GObj).Tag(), args[1])
	}

	x, ok := args[2].(*Poly)
	if !ok || !x.isVar() {
		return nil, fmt.Errorf("%s(3rd arg): expected var: %d:%v", name, args[2].(GObj).Tag(), args[2])
	}

	j, ok := args[3].(*Int)
	if !ok || !j.IsInt64() || j.Sign() < 0 {
		return nil, fmt.Errorf("%s(4th arg): expected nonnegint: %v", name, args[3])
	}

	return g.ox.Psc(f, h, x.lv, int32(j.Int64())), nil
}

func funcOXSlope(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	f, ok := args[0].(*Poly)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected poly: %d:%v", name, args[0].(GObj).Tag(), args[0])
	}
	h, ok := args[1].(*Poly)
	if !ok {
		return nil, fmt.Errorf("%s(2nd arg): expected poly: %d:%v", name, args[1].(GObj).Tag(), args[1])
	}

	x, ok := args[2].(*Poly)
	if !ok || !x.isVar() {
		return nil, fmt.Errorf("%s(3rd arg): expected var: %d:%v", name, args[2].(GObj).Tag(), args[2])
	}

	j, ok := args[3].(*Int)
	if !ok || !j.IsInt64() || j.Sign() < 0 {
		return nil, fmt.Errorf("%s(4th arg): expected nonnegint: %v", name, args[3])
	}

	return g.ox.Slope(f, h, x.lv, int32(j.Int64())), nil
}

func funcOXSres(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	f, ok := args[0].(*Poly)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected poly: %d:%v", name, args[0].(GObj).Tag(), args[0])
	}
	h, ok := args[1].(*Poly)
	if !ok {
		return nil, fmt.Errorf("%s(2nd arg): expected poly: %d:%v", name, args[1].(GObj).Tag(), args[1])
	}

	x, ok := args[2].(*Poly)
	if !ok || !x.isVar() {
		return nil, fmt.Errorf("%s(3rd arg): expected var: %d:%v", name, args[2].(GObj).Tag(), args[2])
	}

	coef := int32(0)
	if len(args) > 3 {
		cc, ok := args[3].(*Int)
		if !ok || !cc.IsInt64() || cc.Sign() < 0 {
			return nil, fmt.Errorf("%s(4th arg): expected nonnegint: %v", name, args[3])
		}
		coef = int32(cc.Int64())
	}

	return g.ox.Sres(f, h, x.lv, coef), nil
}

// //////////////////////////////////////////////////////////
// CAD
// //////////////////////////////////////////////////////////
func funcExample(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	str := ""
	if len(args) == 1 {
		c, ok := args[0].(*String)
		if !ok {
			return nil, fmt.Errorf("%s() expected string", name)
		}
		str = c.s
	}

	ex := GetExampleFof(str)
	if ex == nil {
		if str == "" {
			return nil, nil
		}
		return nil, fmt.Errorf("%s() invalid 1st arg %s", name, str)
	}

	ll := NewList()
	ll.Append(ex.Input)
	ll.Append(ex.Output)
	ll.Append(NewString(ex.Ref))
	ll.Append(NewString(ex.DOI))

	return ll, nil
}

func funcSimplify(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	c, ok := args[0].(Fof)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg) expected FOF", name)
	}
	var neccon Fof = trueObj
	var sufcon Fof = falseObj

	if len(args) > 1 {
		neccon, ok = args[1].(Fof)
		if !ok {
			return nil, fmt.Errorf("%s(2nd arg) expected FOF", name)
		}
	}
	if len(args) > 2 {
		sufcon, ok = args[2].(Fof)
		if !ok {
			return nil, fmt.Errorf("%s(3rd arg) expected FOF", name)
		}
	}

	return g.simplFof(c, neccon, sufcon), nil
}

func funcSimplNum(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	c, ok := args[0].(Fof)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg) expected FOF", name)
	}

	r, _, _ := c.simplNum(g, nil, nil)
	return r, nil
}

func funcGetFormula(name string, arg interface{}) (Fof, error) {
	fof, ok := arg.(Fof)
	if !ok {
		fstr, ok := arg.(*String)
		if ok {
			example := GetExampleFof(fstr.s)
			if example == nil {
				return nil, fmt.Errorf("%s(): invalid name: %v", name, arg)
			}
			fof = example.Input
		} else {
			return nil, fmt.Errorf("%s(1st arg): expected Fof: %v", name, arg)
		}
	}
	return fof, nil
}

func cadArgProj(name string, v any) (ProjectionAlgo, error) {
	var algo ProjectionAlgo = PROJ_McCallum
	if vstr, ok := v.(*String); ok && vstr.s == "m" {
		algo = PROJ_McCallum
	} else if ok && vstr.s == "h" {
		algo = PROJ_HONG
	} else {
		return algo, fmt.Errorf("%s(2nd arg): invalid proj value: %v", name, v)
	}
	return algo, nil
}

func funcCAD(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	fof, err := funcGetFormula(name, args[0])
	if err != nil {
		return nil, err
	}
	if err = fof.valid(); err != nil {
		return nil, err
	}
	if !fof.isPrenex() {
		return nil, fmt.Errorf("%s(1st arg): prenext formula is expected: %v", name, args[0])
	}
	var algo ProjectionAlgo = PROJ_McCallum
	var var_order bool
	var print_stat bool
	var lifting_strategy bool = true
	if len(args) > 1 {
		dic, ok := args[1].(*Dict)
		if !ok {
			return nil, fmt.Errorf("%s(2nd arg): expected Dict: %v", name, args[1])
		}
		for k, v := range dic.v {
			switch k {
			case "proj":
				algo, err = cadArgProj(name, v)
				if err != nil {
					return nil, err
				}
			case "var":
				var_order = funcArgBoolVal(v)
			case "ls":
				lifting_strategy = funcArgBoolVal(v)
			case "stat":
				print_stat = funcArgBoolVal(v)
			default:
				return nil, fmt.Errorf("%s(2nd arg): unknown option: %s", name, k)
			}
		}
	}

	switch fof.(type) {
	case *AtomT, *AtomF:
		return fof, nil
	}

	var qeopt QEopt
	var cond qeCond
	var o2 []Level
	if var_order {
		cond.neccon = trueObj
		cond.sufcon = falseObj

		fof, o2 = qeopt.qe_cad_varorder_pre(fof, cond, fof.maxVar()+1)
	}

	cad, err := NewCAD(fof, g)
	if err != nil {
		return nil, err
	}
	cad.lift_strategy = lifting_strategy
	_, err = cad.Projection(algo)
	if err != nil {
		return nil, err
	}
	g.log(1, 1, "go lift\n")
	err = cad.Lift()
	if err != nil {
		return nil, err
	}
	g.log(1, 1, "go sfc\n")
	fml, err := cad.Sfc()
	if err != nil {
		return nil, err
	}
	if var_order {
		fml = qeopt.qe_cad_varorder_post(fml, cond, fof.maxVar()+1, o2)
	}
	if print_stat {
		err = cad.Print(NewString("stat"))
		if err != nil {
			return nil, err
		}
	}

	return fml, nil
}

func funcCADinit(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	c, err := funcGetFormula(name, args[0])
	if err != nil {
		return nil, err
	}

	var lifting_strategy bool = true
	if len(args) > 1 {
		dic, ok := args[1].(*Dict)
		if !ok {
			return nil, fmt.Errorf("%s(2nd arg): expected Dict: %v", name, args[1])
		}
		for k, v := range dic.v {
			switch k {
			case "ls":
				lifting_strategy = funcArgBoolVal(v)
			default:
				return nil, fmt.Errorf("%s(2nd arg): unknown option: %s", name, k)
			}
		}
	}

	cad, err := NewCAD(c, g)
	if err != nil {
		return nil, err
	}
	cad.lift_strategy = lifting_strategy
	return cad, nil
}

func funcCADproj(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	c, ok := args[0].(*CAD)
	if !ok {
		return nil, fmt.Errorf("%s() expected CAD generated by cadinit()", name)
	}

	var algo ProjectionAlgo = PROJ_McCallum
	if len(args) > 1 {
		var err error
		algo, err = cadArgProj(name, args[1])
		if err != nil {
			return nil, err
		}
	}

	p, err := c.Projection(algo)
	return p, err
}

func funcCADlift(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	c, ok := args[0].(*CAD)
	if !ok {
		return nil, fmt.Errorf("%s(1st-arg) expected CAD generated by cadinit()", name)
	}

	index := make([]int, len(args)-1)
	for i := 1; i < len(args); i++ {
		v, ok := args[i].(*Int)
		if !ok || !v.IsInt64() || v.Int64() > 10000000 || v.Int64() < -1 {
			return nil, fmt.Errorf("%s(%dth arg) expected index", name, i)
		}
		index[i-1] = int(v.Int64())
	}

	err := c.Lift(index...)
	return c, err
}

func funcCADsfc(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	cad, ok := args[0].(*CAD)
	if !ok {
		return nil, fmt.Errorf("%s(1st-arg) expected CAD generated by cadinit()", name)
	}

	return cad.Sfc()
}

func funcVS(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	fof, ok := args[0].(FofQ)
	if !ok || !fof.Fml().IsQff() {
		return nil, fmt.Errorf("%s() expected prenex-FOF", name)
	}
	maxdeg := 1
	if len(args) > 1 {
		deg, ok := args[1].(*Int)
		if !ok || (!deg.Equals(two) && !deg.IsOne()) {
			return nil, fmt.Errorf("%s(2nd-arg) expected 1 or 2", name)
		}
		maxdeg = int(deg.Int64())
	}

	var fml Fof
	fml = fof
	for _, q := range fof.Qs() {
		fml = vs_main(fml, q, maxdeg, g)
	}
	return fml, nil
}

// //////////////////////////////////////////////////////////
// util
// //////////////////////////////////////////////////////////
func printQEPCADheader(fml Fof) error {
	if !fml.isPrenex() {
		return fmt.Errorf("not a prenex formula")
	}

	maxv := fml.maxVar()
	freev := make([]Level, 0, maxv+1)
	num := 0
	for i := Level(0); i <= maxv; i++ {
		if fml.hasVar(i) {
			num++
			if fml.hasFreeVar(i) {
				freev = append(freev, i)
			}
		}
	}
	sep := '('
	for _, v := range freev {
		fmt.Printf("%c%s", sep, VarStr(v))
		sep = ','
	}
	if fq, ok := fml.(FofQ); ok {
		for fq != nil {
			for _, v := range fq.Qs() {
				fmt.Printf("%c%s", sep, VarStr(v))
				sep = ','
			}
			fq, ok = fq.Fml().(FofQ)
			if !ok {
				break
			}
		}
	}
	fmt.Printf(")\n%d\n", len(freev))
	return nil
}

func funcPrint(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	switch cc := args[0].(type) {
	case *CAD:
		return nil, cc.Print(args[1:]...)
	case fmt.Formatter:
		if len(args) > 2 {
			return nil, fmt.Errorf("invalid # of arg")
		}
		t := "org"
		if len(args) == 2 {
			s, ok := args[1].(*String)
			if !ok {
				return nil, fmt.Errorf("invalid 2nd arg")
			}
			t = s.s
		}
		switch t {
		case "org":
			fmt.Printf("%v\n", cc)
		case "tex":
			fmt.Printf("%P\n", cc)
		case "src":
			fmt.Printf("%S\n", cc)
		case "dump":
			fmt.Printf("%V\n", cc)
		case "qepcad":
			if fml, ok := cc.(Fof); ok {
				err := printQEPCADheader(fml)
				if err != nil {
					return nil, err
				}
			}
			fmt.Printf("%Q\n", cc)
		default:
			if len(t) > 0 && t[0] == '%' {
				fmt.Printf(t+"\n", cc)
			} else {
				return nil, fmt.Errorf("invalid 2nd arg")
			}
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("print(): unsupported object is specified")
	}
}

func funcQuit(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	code := 0
	if len(args) > 0 {
		c, ok := args[0].(*Int)
		if !ok {
			return nil, fmt.Errorf("%s() expected int", name)
		}
		if c.Sign() < 0 || !c.IsInt64() || c.Int64() > 125 {
			return nil, fmt.Errorf("%s() expected integer in the range [0, 125]", name)
		}
		code = int(c.Int64())
	}
	os.Exit(code)
	return nil, nil
}

func funcSleep(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	c, ok := args[0].(*Int)
	if !ok {
		return nil, fmt.Errorf("%s() expected int", name)
	}
	if c.Sign() <= 0 {
		return nil, nil
	}

	v := c.Int64()
	time.Sleep(time.Millisecond * time.Duration(v))
	return nil, nil
}

func funcVerbose(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	c, ok := args[0].(*Int)
	if !ok || !c.IsInt64() {
		return nil, fmt.Errorf("%s(1st arg) expected int", name)
	}

	if len(args) > 1 {
		d, ok := args[1].(*Int)
		if !ok || !c.IsInt64() {
			return nil, fmt.Errorf("%s(2nd arg) expected int", name)
		}
		if d.Sign() <= 0 {
			g.verbose_cad = 0
		} else {
			g.verbose_cad = int(d.Int64())
		}
	}

	if c.Sign() <= 0 {
		g.verbose = 0
	} else {
		g.verbose = int(c.Int64())
	}
	return nil, nil
}

func funcTime(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	return nil, nil
}

// //////////////////////////////////////////////////////////
// poly
// //////////////////////////////////////////////////////////
func funcSubst(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	if len(args)%2 != 1 {
		return nil, fmt.Errorf("%s() invalid args", name)
	} else if len(args) == 1 {
		return args[0], nil
	}

	rlv := make([]struct {
		r  RObj
		lv Level
	}, (len(args)-1)/2)

	j := 0
	for i := 1; i < len(args); i += 2 {
		p, ok := args[i].(*Poly)
		if !ok || !p.isVar() {
			return nil, fmt.Errorf("%s() invalid %d'th arg: %v", name, i+1, args[i])
		}

		rlv[j].lv = p.lv

		v, ok := args[i+1].(RObj)
		if !ok {
			return nil, fmt.Errorf("%s() invalid %d'th arg", name, i+2)
		}
		rlv[j].r = v
		j++
	}
	rlv = rlv[:j]

	o := args[0].(GObj)
	for _, r := range rlv {
		o = gobjSubst(o, r.r, r.lv)
	}

	return o, nil
}

func funcDeg(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	d, ok := args[1].(*Poly)
	if !ok || !d.isVar() {
		return nil, fmt.Errorf("%s(2nd arg): expected var: %v", name, args[1])
	}

	var deg int
	switch p := args[0].(type) {
	case Fof:
		deg = p.Deg(d.lv)
	case *Poly:
		deg = p.Deg(d.lv)
	case RObj:
		if p.IsZero() {
			deg = -1
		} else {
			deg = 0
		}
	default:
		return nil, fmt.Errorf("%s(1st arg): expected poly or FOF: %v", name, args[0])
	}

	return NewInt(int64(deg)), nil
}

func funcCoef(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	_, ok := args[0].(RObj)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected RObj: %v", name, args[0])
	}

	c, ok := args[1].(*Poly)
	if !ok || !c.isVar() {
		return nil, fmt.Errorf("%s(2nd arg): expected var: %v", name, args[1])
	}

	d, ok := args[2].(*Int)
	if !ok {
		return nil, fmt.Errorf("%s(3rd arg): expected int: %v", name, args[2])
	}

	if d.Sign() < 0 {
		return zero, nil
	}
	rr, ok := args[0].(*Poly)
	if !ok {
		if d.Sign() == 0 {
			return args[0], nil
		} else {
			return zero, nil
		}
	}
	if !d.n.IsUint64() {
		return zero, nil
	}

	return rr.Coef(c.lv, uint(d.n.Uint64())), nil
}

func funcCounter(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	d := NewDict()
	for k, v := range DebugCounter {
		d.Set(k, NewInt(int64(v)))
	}
	return d, nil
}

func funcArgBoolVal(val GObj) bool {
	switch v := val.(type) {
	case RObj:
		return !v.IsZero()
	case *String:
		return v.s == "yes" || v.s == "Y" || v.s == "y"
	default:
		return false
	}
}

func funcQE(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	fof, err := funcGetFormula(name, args[0])
	if err != nil {
		return nil, err
	}
	if err = fof.valid(); err != nil {
		return nil, err
	}
	opt := NewQEopt()
	set_verbose := false
	verbose := g.verbose
	if len(args) > 1 {
		dic, ok := args[1].(*Dict)
		if !ok {
			return nil, fmt.Errorf("%s(2nd arg): expected Dict: %v", name, args[1])
		}
		for k, v := range dic.v {
			switch k {
			case getQEoptStr(QEALGO_EQQUAD):
				opt.SetAlgo(QEALGO_EQQUAD, funcArgBoolVal(v))
			case getQEoptStr(QEALGO_EQLIN):
				opt.SetAlgo(QEALGO_EQLIN, funcArgBoolVal(v))
			case getQEoptStr(QEALGO_VSQUAD):
				opt.SetAlgo(QEALGO_VSQUAD, funcArgBoolVal(v))
			case getQEoptStr(QEALGO_VSLIN):
				opt.SetAlgo(QEALGO_VSLIN, funcArgBoolVal(v))
			case getQEoptStr(QEALGO_NEQ):
				opt.SetAlgo(QEALGO_NEQ, funcArgBoolVal(v))
			case getQEoptStr(QEALGO_SMPL_EVEN):
				opt.SetAlgo(QEALGO_SMPL_EVEN, funcArgBoolVal(v))
			case getQEoptStr(QEALGO_SMPL_HOMO):
				opt.SetAlgo(QEALGO_SMPL_HOMO, funcArgBoolVal(v))
			case getQEoptStr(QEALGO_SMPL_TRAN):
				opt.SetAlgo(QEALGO_SMPL_TRAN, funcArgBoolVal(v))
			case getQEoptStr(QEALGO_SMPL_ROTA):
				opt.SetAlgo(QEALGO_SMPL_ROTA, funcArgBoolVal(v))
			case getQEoptStr(QEALGO_ATOM):
				opt.SetAlgo(QEALGO_ATOM, funcArgBoolVal(v))
			case getQEoptStr(QEALGO_SDC):
				opt.SetAlgo(QEALGO_SDC, funcArgBoolVal(v))
			case "verbose":
				if val, ok := v.(*Int); ok && val.IsInt64() {
					g.verbose = int(val.Int64())
					set_verbose = true
				} else {
					return nil, fmt.Errorf("%s(2nd arg): invalid option value: %s: %v.", name, k, v)
				}
			default:
				return nil, fmt.Errorf("%s(2nd arg): unknown option: %s", name, k)
			}
		}
	}

	qff := g.QE(fof, opt)
	if set_verbose {
		g.verbose = verbose
	}
	return qff, nil
}

func funcRealRoot(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	p, ok := args[0].(*Poly)
	if !ok {
		return nil, fmt.Errorf("%s(): expected poly: %v", name, args[0])
	}

	q, ok := args[1].(*Int)
	if !ok {
		return nil, fmt.Errorf("%s(): expected int: %v", name, args[1])
	}

	return p.RealRootIsolation(int(q.Int64()))
}

func funcOXRes(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	f, ok := args[0].(*Poly)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected poly: %d:%v", name, args[0].(GObj).Tag(), args[0])
	}
	h, ok := args[1].(*Poly)
	if !ok {
		return nil, fmt.Errorf("%s(2nd arg): expected poly: %d:%v", name, args[1].(GObj).Tag(), args[1])
	}
	x, ok := args[2].(*Poly)
	if !ok || !x.isVar() {
		return nil, fmt.Errorf("%s(3rd arg): expected var: %d:%v", name, args[2].(GObj).Tag(), args[2])
	}

	return g.ox.Resultant(f, h, x.lv), nil
}

func funcRootBound(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	p, ok := args[0].(*Poly)
	if !ok {
		return nil, fmt.Errorf("%s(): expected poly: %v", name, args[0])
	}
	return p.RootBound()
}

func funcIndets(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	b := make([]bool, len(varlist))
	p, ok := args[0].(indeter)
	if !ok {
		return NewList(), nil
	}
	p.Indets(b)

	ret := make([]interface{}, 0, len(b))
	for i := 0; i < len(b); i++ {
		if b[i] {
			ret = append(ret, NewPolyVar(Level(i)))
		}
	}
	return NewList(ret...), nil
}

func funcIntv(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	prec := uint(53)
	if g, ok := args[0].(GObj); ok && len(args) == 1 {
		return gobjToIntv(g, prec), nil
	}

	a, ok := args[0].(NObj)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected number: %v", name, args[0])
	}
	b := a
	if len(args) > 1 {
		bb, ok := args[1].(NObj)
		if !ok {
			return nil, fmt.Errorf("%s(2nd arg): expected number: %v", name, args[1])
		}
		if _, ok = a.(NObj); !ok {
			return nil, fmt.Errorf("%s(1st arg): expected number: %v", name, args[0])
		}
		b = bb
	}
	if len(args) > 2 {
		pp, ok := args[2].(*Int)
		if !ok || !pp.IsInt64() || pp.Sign() <= 0 {
			return nil, fmt.Errorf("%s(3rd arg): expected int: %v", name, args[2])
		}

		prec = uint(pp.Int64())
	}
	aa := a.toIntv(prec)
	if aintv, ok := aa.(*Interval); ok {
		bintv, ok := b.toIntv(prec).(*Interval)
		if !ok {
			return nil, fmt.Errorf("%s(2nd arg): expected number: %v", name, args[1])
		}
		u := newInterval(prec)
		u.inf = aintv.inf
		u.sup = bintv.sup
		return u, nil
	} else {
		return aa, nil
	}
}

// //////////////////////////////////////////////////////////
// integer
// //////////////////////////////////////////////////////////
func funcIGCD(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	a, ok := args[0].(*Int)
	if !ok {
		return nil, fmt.Errorf("%s(1st arg): expected int: %v", name, args[0])
	}
	b, ok := args[1].(*Int)
	if !ok {
		return nil, fmt.Errorf("%s(2nd arg): expected int: %v", name, args[1])
	}

	return a.Gcd(b), nil
}

////////////////////////////////////////////////////////////
// system
////////////////////////////////////////////////////////////

func funcLoad(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	return nil, fmt.Errorf("%s not implemented", name) // @TODO
}

func funcSave(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	return nil, fmt.Errorf("%s not implemented", name) // @TODO
}

func funcLen(g *Ganrac, name string, args []interface{}) (interface{}, error) {
	p, ok := args[0].(lener)
	if !ok {
		return nil, fmt.Errorf("%s(): not supported: %v", name, args[0])
	}
	return NewInt(int64(p.Len())), nil
}

func funcHelp(builtin_func_table []func_table, name string, args []interface{}) (interface{}, error) {
	if len(args) == 0 {
		return funcHelps(builtin_func_table, "@")
	}

	p, ok := args[0].(*String)
	if !ok {
		return nil, fmt.Errorf("%s(): required help(\"string\"):", name)
	}

	return funcHelps(builtin_func_table, p.s)
}

func funcHelps(builtin_func_table []func_table, name string) (interface{}, error) {
	if name == "@" { // 引数なし
		fmt.Printf("GANRAC\n")
		fmt.Printf("==========\n")
		fmt.Printf("TOKENS:\n")
		fmt.Printf("  integer      : `[0-9]+`\n")
		fmt.Printf("  string       : `\"[^\"]*\"`\n")
		fmt.Printf("  indeterminate: `[a-z][a-zA-Z0-9_]*`\n")
		fmt.Printf("  variable     : `[A-Z][a-zA-Z0-9_]*`\n")
		fmt.Printf("  true/false   : `true`/`false`\n")
		fmt.Printf("\n")
		fmt.Printf("OPERATORS:\n")
		fmt.Printf("  + - * / ^\n")
		fmt.Printf("  < <= > >= == !=\n")
		fmt.Printf("  && ||\n")
		fmt.Printf("\n")
		fmt.Printf("FUNCTIONS:\n")

		maxlen := 15
		for _, fv := range builtin_func_table {
			if s := len(fv.name) + len(fv.args); maxlen < s {
				maxlen = s
			}
		}
		for _, fv := range builtin_func_table {
			mark := " "
			label := fv.name + fv.args
			if fv.ox {
				mark = "*"
			} else if fv.args[len(fv.args)-1] == '@' {
				mark = "@"
				label = label[:len(label)-1]
			}
			fmt.Printf("%s %-*s\t%s\n", mark, maxlen, label, fv.descript)
		}
		fmt.Printf("\n")
		fmt.Printf(" * ... required CAS\n")
		fmt.Printf(" @ ... not implemented\n")
		fmt.Printf("\n")
		fmt.Printf("\n")
		fmt.Printf("EXAMPLES:\n")
		fmt.Printf("  > %s(x, y, z);  # init variable order.\n", init_var_funcname)
		fmt.Printf("  > F = x^2 + 2;\n")
		fmt.Printf("  > deg(F, x);\n")
		fmt.Printf("  2\n")
		fmt.Printf("  > t;\n")
		fmt.Printf("  error: undefined variable `t`\n")
		fmt.Printf("  > qe(ex([x], x+1 = 0 && x < 0));\n")
		fmt.Printf("  true\n")
		fmt.Printf("  > help(\"deg\");\n")
		return nil, nil
	}
	// 引数ありのばあい
	const sep = "================\n"
	const indent = "  "
	for _, fv := range builtin_func_table {
		if fv.name == name {
			fmt.Printf("%s%s: %s\n", fv.name, fv.args, fv.descript)
			if fv.arguments != nil {
				fmt.Printf("\nArguments:\n%s", sep)
				maxn := 0
				for _, a := range fv.arguments {
					if maxn < len(a.name) {
						maxn = len(a.name)
					}
				}
				for _, a := range fv.arguments {
					fmt.Printf("%s%-*s: %s", indent, maxn, a.name, a.atype)
					if a.defval != "" {
						fmt.Printf("\t(default: %s)", a.defval)
					}
					fmt.Printf("\n")
					maxm := 0
					maxd := 0
					for _, b := range a.dic {
						if ln := len(b.name); maxm < ln {
							maxm = ln
						}
						if ln := len(b.atype); maxd < ln {
							maxd = ln
						}
					}
					if maxm+len(indent) == maxn {
						// コロンの位置が揃うと見づらいため，ずらす
						maxm++
					}
					for _, b := range a.dic {
						fmt.Printf("%s%s%-*s: %-*s", indent, indent, maxm, b.name, maxd, b.atype)
						if b.defval != "" {
							fmt.Printf("  (default: %s)", b.defval)
						}
						fmt.Printf("\n")
					}
				}
			}
			if fv.returns != "" {
				fmt.Printf("\nReturns:\n%s", sep)
				fmt.Printf("%s%s\n", indent, fv.returns)
			}
			if fv.examples != "" {
				fmt.Printf("\nExamples:\n%s", sep)
				fmt.Printf("%s", add_indent(fv.examples, indent))
			}

			return nil, nil
		}
	}

	return nil, fmt.Errorf("unknown function `%s()`\n", name)
}

func add_indent(input, ws string) string {
	vs := strings.Split(input, "\n")
	ret := ""
	for i, v := range vs {
		if i == 0 && v == "" {
			continue
		} else if v == "" {
			ret += "\n"
		} else {
			ret += ws + v + "\n"
		}
	}
	return ret
}
