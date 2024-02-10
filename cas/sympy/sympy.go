package sympy

// https://www.sympy.org/

// #include <Python.h>
// PyAPI_FUNC(PyObject *) Py_CompileStringGanrac(const char *str, const char *p, int s) {
//   return Py_CompileString(str, p, s);
// }
import "C"

import (
	"fmt"
	. "github.com/hiwane/ganrac"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"
)

type SymPy struct {
	pModule    *C.PyObject
	pstop      *C.PyObject
	pGC        *C.PyObject
	pGCD       *C.PyObject
	pResultant *C.PyObject
	pDiscrim   *C.PyObject
	pPsc       *C.PyObject
	pSlope     *C.PyObject
	pSres      *C.PyObject
	pGB        *C.PyObject
	pReduce    *C.PyObject
	pEval      *C.PyObject
	pValid     *C.PyObject
	pFactor    *C.PyObject
	pVarInit   *C.PyObject

	ganrac *Ganrac
	logger *log.Logger

	varnum int

	cnt     int
	is_gc   bool
	verbose bool
}

func (sympy *SymPy) log(format string, a ...interface{}) {
	if sympy.verbose {
		_, file, line, _ := runtime.Caller(1)
		// runtime.FuncForPC(pc).Name()  (pc := 1st element of Caller())
		_, fname := filepath.Split(file)
		fname = fname[:len(fname)-3] // suffix .go は不要
		fmt.Printf("[%d] %14s:%4d: ", sympy.cnt, fname, line)
		fmt.Printf(format, a...)
	}
}

func (sympy *SymPy) Inc(n int) int {
	sympy.cnt += n
	return sympy.cnt
}

func (sympy *SymPy) Close() error {
	sympy.log("sympy.Close() start\n")

	if sympy.is_gc {
		sympy.log("sympy.Close(): tracemalloc.stop()\n")
		r := callFunction(sympy.pstop)
		C.Py_DecRef(r)
	}

	C.Py_DecRef(sympy.pModule)
	C.Py_DecRef(sympy.pstop)
	C.Py_DecRef(sympy.pGCD)
	C.Py_DecRef(sympy.pResultant)
	C.Py_DecRef(sympy.pDiscrim)
	C.Py_DecRef(sympy.pPsc)
	C.Py_DecRef(sympy.pSlope)
	C.Py_DecRef(sympy.pSres)
	C.Py_DecRef(sympy.pGB)
	C.Py_DecRef(sympy.pReduce)
	C.Py_DecRef(sympy.pEval)
	C.Py_DecRef(sympy.pValid)
	C.Py_DecRef(sympy.pFactor)
	C.Py_DecRef(sympy.pVarInit)

	if false {
		ns := C.PyLong_FromLong(C.long(1))
		ret := callFunction(sympy.pGC, ns)
		C.Py_DecRef(ret)
		C.Py_DecRef(sympy.pGC)
	}

	// C.Py_Finalize()
	sympy.log("sympy.Close() end\n")
	return nil
}

func NewSymPy(g *Ganrac, fname string) (*SymPy, error) {
	C.Py_Initialize()

	var str string

	sympy := new(SymPy)
	sympy.ganrac = g
	sympy.is_gc = false
	sympy.verbose = false

	/*
			   [<Signals.SIGINT: 2>, <cyfunction python_check_interrupt at 0x7f41066a0ee0>]
			   [8, <built-in function default_int_handler>]
			   [8, None]

		def debuglog(s):
			with open("/tmp/boo.txt", "a") as fp:
				print(s, file=fp)

		# __call__ の引数 monomial (たぶん，次数の列）が固定長ではない.
		# つまり，n 変数の GB 計算途中で, len(monomial) < n な場合がある.
		# 多項式の実装が変わらない限り block order には対応できないと思われる
		from sympy.polys.orderings import MonomialOrder
		class BlockOrder(MonomialOrder):

		  def __init__(self, n, order):
		    self.block_n = n
		    self.block_order = order
		    self.alias = f'block({order},{n})'

		  def __call__(self, monomial):
		    s = monomial[:self.block_n]
		    t = monomial[self.block_n:]
		    return (sum(s), tuple(reversed([-m for m in s])),
		      sum(t), tuple(reversed([-m for m in t])))
			#return self.block_order(s) + self.block_order(t)


		F = [x*x-z, x*y-1, x**3-x**2*y-x*x-1]
		print(sympy.groebner(F, x, y, z, order=BlockOrder(2, grevlex)))

	*/

	// TypeError: signal handler must be signal.SIG_IGN, signal.SIG_DFL, or a callable object
	str = `
import signal
try:
	signal.signal(8, signal.SIG_DFL)
except Exception:
	pass
import sympy
from sympy.polys import polyoptions
from sympy.polys.polytools import GroebnerBasis

def debuglog(s):
	with open("/tmp/ganrac_sympy.debug", "a") as fp:
		print(s, file=fp, flush=True)
	pass

def gan_factor(polystr: str):
	F = sympy.factor_list(polystr)
	G = [[F[0], 1]] + [[u[0], u[1]] for u in F[1]]
	return str(G)

def gan_gcd(p, q):
	debuglog(["gcd", p, q])
	F = sympy.parse_expr(p)
	G = sympy.parse_expr(q)
	debuglog(["gcd", F, G])
	return str(sympy.gcd(F, G))

def gan_res(p, q, x):
	debuglog(["res", p, q, x])
	F = sympy.parse_expr(p)
	G = sympy.parse_expr(q)
	X = sympy.Symbol(x)
	debuglog(["res", F, G, X])
	return str(sympy.resultant(F, G, X))

def gan_discrim(p, x):
	debuglog(["discrim", p, x])
	F = sympy.parse_expr(p)
	X = sympy.Symbol(x)
	debuglog(["discrim", F, X])
	return str(sympy.discriminant(F, X))

def gan_coef(f, x, i):
	if hasattr(f, 'coeff'):
		return f.coeff(x, i)
	elif i == 0:
		return f
	else:
		return 0

def gan_sres(p, q, x, CC):
	debuglog(["sres", p, q, x, CC])
	F = sympy.parse_expr(p)
	G = sympy.parse_expr(q)
	X = sympy.Symbol(x)
	debuglog(["sres", F, G, X])
	degF = sympy.Poly(F).degree(X)
	degG = sympy.Poly(G).degree(X)
	if degF > degG:
		N = degF - 1
	elif degF < degG:
		N = degG - 1
	else:
		N = degF - 0
	V = [0] * (N + 2)
	V[N] = G
	V[N+1] = F
	N -= 1
	for v in sympy.subresultants(F, G, X)[2:]:
		w = sympy.Poly(v)
		if w.has(X):
			D = w.degree(X)
		else:
			D = 0
		V[N] = v
		if N != D:
			# N[D] を算出する. subresultant theorem (2)
			numer = V[N] * (gan_coef(V[N], X, D) ** (N - D))
			denom = gan_coef(V[N + 1], X, N + 1) ** (N - D)
			q, r = sympy.div(numer, denom)
			V[D] = q
		N = D - 1
	V = V[:-2]

	if CC == 0:
		return str(V)
	if CC == 1:
		return str([gan_coef(v, X, i) for i, v in enumerate(V)])
	if CC == 2:
		return str([gan_coef(v, X, 0) for v in V])
	if CC == 3:
		return str([V[0]] + [gan_coef(V[i], X, 0) + X**i * gan_coef(V[i], X, i) for i in range(1, len(V))])
	return str(V)

def gan_psc(p, q, x, J):
	debuglog(["psc", p, q, x, J])
	F = sympy.parse_expr(p)
	G = sympy.parse_expr(q)
	X = sympy.Symbol(x)
	debuglog(["psc", F, G, X])
	M = sympy.Poly(F).degree(X)
	N = sympy.Poly(G).degree(X)
	L = M+N-2*J;
	S = sympy.zeros(L, L)
	for D in range(M, -1, -1):
		AI = F.coeff(X, D)
		for I in range(min([N - J, L + D - M - 1])):
			S[I, M-D+I] = AI
	for I in range(N - J - 1, max([0, (N-J-1)-J]) -1, -1):
		S[I, L-1] = F.coeff(X, I-(N-J-1)+J);
	for D in range(N, -1, -1):
		BI = G.coeff(X, D)
		for I in range(min([M - J, L + D - N - 1])):
			S[I+N-J, N-D+I] = BI
	for I in range(M - J - 1, max([0, M-J-1-J])-1, -1):
		S[I+N-J, L-1] = G.coeff(X, I-(M-J-1)+J);

	return str(S.det())

def gan_comb(A, B):
	C = 1
	for I in range(1, B+1):
		C *= (A - I + 1)
		C //= I
	return C

def gan_slope(p, q, x, K):
	debuglog(["slope-i", p, x, K])
	F = sympy.parse_expr(p)
	G = sympy.parse_expr(q)
	X = sympy.Symbol(x)
	debuglog(["slope-o", F, G, X])
	M = sympy.Poly(F).degree(X)
	N = sympy.Poly(G).degree(X)
	L = N - K
	S = sympy.zeros(L+1, L+1)
	CMK = gan_comb(M, K + 1)
	for J in range(L):
		S[0, J] = F.coeff(X, M-J)
		for I in range(1, L-J):
			S[I, I+J] = S[0, J]
		if 0 <= L - J <= L:
			S[L-J, L] = (CMK - gan_comb(M-J, K+1)) * S[0, J]
		S[L, J] = G.coeff(X, N - J)
	J = L
	if M-J >= 0:
		S[0, L] = (CMK - gan_comb(M - J, K + 1)) * F.coeff(X, M-J)
	S[L, L] = CMK * G.coeff(X, K)
	return str(S.det())

def gan_gb(polys, vars, n):
	# https://docs.sympy.org/latest/modules/polys/reference.html#sympy.polys.polytools.groebner
	order = 'grevlex'
	debuglog(["gb", polys, vars, n])
	F = sympy.parse_expr(polys)
	V = []
	for v in vars:
		V.append(sympy.Symbol(v))
	debuglog(["gbo", F, V])
	G = sympy.groebner(F, *V, order=order)
	# G = [g / g.content() for g in G]
	return str(G.exprs)

def gan_reduce(p, gb, vars, n):
	order = 'grevlex'
	debuglog(["reducei", p, gb, vars, n])
	P = sympy.parse_expr(p)
	FF = sympy.parse_expr(gb)
	VARS = [sympy.Symbol(v) for v in vars]
	F, opt = sympy.parallel_poly_from_expr(FF, *VARS, order=order)
	debuglog(["reduceo", P, F, VARS])

	try:
		debuglog(["reduce 1"])
		debuglog(["reduce 2", opt])
		GB = GroebnerBasis._new(F, opt)
		debuglog(["reduce 3", GB])
		debuglog(["reduce 3.5", [GB.reduce(P)]])
		UU = GB.reduce(P)
		debuglog(["reduce 4", UU])
		U = UU[-1]
		debuglog(["reduce 4"])
		debuglog(U)
		debuglog(["reduce 5"])
		U = str(U)
		debuglog(["reduce 7"])
		return U
	except Exception as e:
		print("reduce() failed")
		print(GB)
		print(P)
		debuglog(GB)
		debuglog(P)
		debuglog("reduce failed: " + str(e))
		raise


def gan_eval(s):
	F = sympy.parse_expr(s)
	return str(F)
`

	/*
	   po0=&{1 0x7f4c97776840}
	   po1=&{2 0x7f4c97770920}
	   pname=0x30fb9c0: ganrac
	   pModule=&{3 0x7f4c97770920}
	   not callable! ganrac
	   pFunc(varinit)=&{2 0x7f4c97774080}
	   callable 1
	   args=&{1 0x7f4c9776f480}, vlist=&{1 0x7f4c97773480}
	   retval=&{1 0x7f4c97772da0}
	   retva2=
	   pFunc(varinit)=&{2 0x7f4c97774080}
	   callable 1
	   retval=&{1 0x7f4c9776f0a0}
	   retva2=2*(2*x + 3)*(2*x - 3)
	*/
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	po := C.Py_CompileStringGanrac(cstr, C.CString(fname), C.Py_file_input)

	mname := "ganrac" // module名
	pName := C.CString(mname)
	defer C.free(unsafe.Pointer(pName))
	po = C.PyImport_ExecCodeModule(pName, po) // ganrac module として str を実行
	if po == nil {
		if err := C.PyErr_Occurred(); err != nil {
			C.PyErr_Print()
			return nil, fmt.Errorf("PyImport_ExecCodeModule failed.")
		}
	}

	sympy.pModule = C.PyImport_ImportModule(pName)
	if sympy.pModule == nil {
		return nil, fmt.Errorf("import module %s failed", mname)
	}

	var psnapstr *C.PyObject

	for _, tbl := range []struct {
		fname string
		pos   **C.PyObject
	}{
		//		{"gan_snapstart", &psnapstr},
		//		{"gan_snapstop", &sympy.pstop},
		//		{"gan_gc", &sympy.pGC},
		{"gan_gcd", &sympy.pGCD},
		{"gan_res", &sympy.pResultant},
		{"gan_discrim", &sympy.pDiscrim},
		{"gan_psc", &sympy.pPsc},
		{"gan_slope", &sympy.pSlope},
		{"gan_sres", &sympy.pSres},
		{"gan_gb", &sympy.pGB},
		{"gan_reduce", &sympy.pReduce},
		{"gan_eval", &sympy.pEval},
		{"gan_factor", &sympy.pFactor},
		//		{"varinit", &sympy.pVarInit},
	} {
		p, err := loadFunction(sympy.pModule, tbl.fname)
		if err != nil {
			return nil, err
		}
		*tbl.pos = p
	}

	if sympy.is_gc {
		fmt.Printf("NewSymPy() call tracemalloc.start()\n")
		r := callFunction(psnapstr)
		C.Py_DecRef(r)
	}

	return sympy, nil
}

func (sympy *SymPy) GC() error {
	if !sympy.is_gc {
		return nil
	}
	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	n, _ := runtime.MemProfile(nil, true)

	rets := "boo"
	if true {
		ns := C.PyLong_FromLong(C.long(0))
		ret := callFunction(sympy.pGC, ns)
		defer C.Py_DecRef(ret)
		rets = toGoString(ret)
	}
	fmt.Printf("GC<%7d> py=%9s, go=%9d, %9d, %9d\n", sympy.cnt, rets, ms.Alloc, ms.HeapAlloc, n)
	return nil
}

/* 変数の最大レベルを返す */
func (sympy *SymPy) varn(polys ...*Poly) int {
	n := 0
	for _, p := range polys {
		if n < int(p.Level())+1 {
			n = int(p.Level() + 1)
		}
	}
	return n + 1
}

func (sympy *SymPy) toGaNRAC(s string) interface{} {
	s = strings.ReplaceAll(strings.ReplaceAll(s, "x", "$"), "**", "^") + ";" // GaNRAC 用に変換
	u, err := sympy.ganrac.Eval(strings.NewReader(s))
	if err != nil {
		fmt.Fprintf(os.Stderr, "sympy.toGaNRAC(%s) failed: %s.", s, err.Error())
		return nil
	}
	return u
}

func (sympy *SymPy) EvalRObj(s string) RObj {
	u := sympy.toGaNRAC(s)
	if u == nil {
		return nil
	}
	if uu, ok := u.(RObj); ok {
		return uu
	} else {
		fmt.Fprintf(os.Stderr, "sympy.EvalRObj(%s) invalid.", s)
		return nil
	}
}

func (sympy *SymPy) EvalList(s string) *List {
	u := sympy.toGaNRAC(s)
	if u == nil {
		return nil
	}
	if uu, ok := u.(*List); ok {
		return uu
	}
	fmt.Fprintf(os.Stderr, "sympy.EvalList(%s) invalid.", s)
	return nil
}

func (sympy *SymPy) Gcd(p, q *Poly) RObj {
	sympy.log("Gcd(%s,%s) start!\n", p, q)
	varn := sympy.varn(p, q)

	// 変数はすべて xi 形式にする
	ps := polyToPyString(p)
	qs := polyToPyString(q)
	// fmt.Printf("ps=%s qs=%p\n", fmt.Sprintf("%I", p), qs)

	ret := callFunctionv(sympy.pGCD, varn, ps, qs)

	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pGCD\n", sympy.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	if strings.HasPrefix(retstr, "none") {
		panic(fmt.Sprintf("<%d> stop..... varlist is broken. %s.", sympy.cnt, retstr))
	}
	return sympy.EvalRObj(retstr)
}

func (sympy *SymPy) Factor(q *Poly) *List {
	sympy.Inc(1)
	sympy.log("Factor(%s) start!\n", q)

	ps := polyToPyString(q) // 変数はすべて xi 形式にする

	ret := callFunction(sympy.pFactor, ps)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "call object failed pFactor\n")
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	rets := toGoString(ret)
	uu := sympy.EvalList(rets)
	if uu == nil {
		fmt.Fprintf(os.Stderr, "sympy.Factor eval() failed: %s.\n", rets)
		return nil
	}
	return uu
}

func (sympy *SymPy) Discrim(p *Poly, lv Level) RObj {
	sympy.log("Discrim(%s) start!\n", p)

	// 変数はすべて xi 形式にする
	ps := polyToPyString(p)
	xs := lvToPyString(lv)

	ret := callFunction(sympy.pDiscrim, ps, xs)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pDiscrim\n", sympy.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sympy.EvalRObj(retstr)
}
func (sympy *SymPy) Resultant(p *Poly, q *Poly, lv Level) RObj {
	sympy.log("Resultant(%s) start!\n", p)

	// 変数はすべて xi 形式にする
	ps := polyToPyString(p)
	qs := polyToPyString(q)
	xs := lvToPyString(lv)

	ret := callFunction(sympy.pResultant, ps, qs, xs)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pResultant\n", sympy.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sympy.EvalRObj(retstr)
}

func (sympy *SymPy) Psc(p *Poly, q *Poly, lv Level, j int32) RObj {
	// sympy.all.Matrix([[x,y],[z,w]]).det()
	ps := polyToPyString(p)
	qs := polyToPyString(q)
	xs := lvToPyString(lv)
	js := C.PyLong_FromLong(C.long(j))

	ret := callFunction(sympy.pPsc, ps, qs, xs, js)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pPsc\n", sympy.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sympy.EvalRObj(retstr)
}

func (sympy *SymPy) Sres(p *Poly, q *Poly, lv Level, k int32) *List {
	sympy.log("Sres(%s,%s,%d,%d) start!\n", p, q, lv, k)
	ps := polyToPyString(p)
	qs := polyToPyString(q)
	xs := lvToPyString(lv)
	ks := C.PyLong_FromLong(C.long(k))

	ret := callFunction(sympy.pSres, ps, qs, xs, ks)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pSres\n", sympy.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sympy.EvalList(retstr)
}

func (sympy *SymPy) Slope(p *Poly, q *Poly, lv Level, k int32) RObj {
	sympy.log("Slope(%s) start!\n", p)
	ps := polyToPyString(p)
	qs := polyToPyString(q)
	xs := lvToPyString(lv)
	ks := C.PyLong_FromLong(C.long(k))

	ret := callFunction(sympy.pSlope, ps, qs, xs, ks)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pSlope\n", sympy.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sympy.EvalRObj(retstr)
}

func (sympy *SymPy) toPyVars(vars *List) *C.PyObject {
	vs := C.PyTuple_New(C.long(vars.Len()))
	for i := 0; i < vars.Len(); i++ {
		v, _ := vars.Geti(i)
		C.PyTuple_SetItem(vs, C.long(i), toPyString(fmt.Sprintf("x%d", v.(*Poly).Level())))
	}
	return vs
}

/*
[1843] F = [x^2-z, x*y-1, x^3-x^2*y-x^2-1]$
[1844] gr(F, [x,y,z], [[0, 2], [0, 1]]);
[z^3-3*z^2-z-1,2*y+z^2-4*z+1,2*x-z^2+2*z+1]
[1845] gr(F, [x,y,z], [[0, 1], [0, 2]]);
[2*y+z^2-4*z+1,(-z-1)*y+z-1,-y^2-2*y+z-2,-x-y+z-1]
*/
func (sympy *SymPy) GB(p *List, vars *List, n int) *List {
	sympy.cnt++
	sympy.log("[%d] GB(%s,%s,%d) start!\n", sympy.cnt, p, vars, n)
	// 変数はすべて xi 形式にする
	ps := listToPyString(p)
	ns := C.PyLong_FromLong(C.long(n))

	vs := sympy.toPyVars(vars)
	defer C.Py_DecRef(vs)

	ret := callFunction(sympy.pGB, ps, vs, ns)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pGB\n", sympy.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sympy.EvalList(retstr)
}

func (sympy *SymPy) Reduce(p *Poly, gb *List, vars *List, n int) (RObj, bool) {
	sympy.cnt++
	sympy.log("[%d] Reduce(%s,%s,%s,%d) start!\n", sympy.cnt, p, gb, vars, n)
	ps := polyToPyString(p)
	gbs := listToPyString(gb)
	vs := sympy.toPyVars(vars)
	for i := 0; i < vars.Len(); i++ {
		_vi, _ := vars.Geti(i)
		vi := _vi.(RObj)
		for j := i + 1; j < vars.Len(); j++ {
			_vj, _ := vars.Geti(j)
			vj := _vj.(RObj)
			if vi.Equals(vj) {
				panic(fmt.Sprintf("stop reduce [%d,%d] %s invalid\n", i, j, vars))
			}
		}
	}
	defer C.Py_DecRef(vs)
	ns := C.PyLong_FromLong(C.long(n))

	ret := callFunction(sympy.pReduce, ps, gbs, vs, ns)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pReduce\n", sympy.cnt)
		C.PyErr_Print()
		return p, false
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sympy.EvalRObj(retstr), false
}

func (sympy *SymPy) Eval(p string) (GObj, error) {
	sympy.log("Eval(%s) start!\n", p)
	ps := toPyString(p)
	ret := callFunction(sympy.pEval, ps) // 変数設定は??
	if ret == nil {
		C.PyErr_Print()
		return nil, fmt.Errorf("<%d> call object failed pEval\n", sympy.cnt)
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	s := sympy.toGaNRAC(retstr)
	if s == nil {
		return nil, fmt.Errorf("SymPy: returns null")
	}
	if sg, ok := s.(GObj); ok {
		return sg, nil
	} else {
		return nil, fmt.Errorf("SymPy: unsupported " + retstr)
	}
}

func (sympy *SymPy) SetLogger(_ *log.Logger) {
}
