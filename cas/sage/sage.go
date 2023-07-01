package sage

/*
404350,
*/

// #include <Python.h>
// #cgo CFLAGS: -I/usr/include/python3.10/
// #cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu/ -lpython3.10
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

type Sage struct {
	pModule    *C.PyObject
	pstop      *C.PyObject
	pGC        *C.PyObject
	pGCD       *C.PyObject
	pResultant *C.PyObject
	pDiscrim   *C.PyObject
	pPsc       *C.PyObject
	pSlope     *C.PyObject
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

func (sage *Sage) log(format string, a ...interface{}) {
	if sage.verbose {
		_, file, line, _ := runtime.Caller(1)
		// runtime.FuncForPC(pc).Name()  (pc := 1st element of Caller())
		_, fname := filepath.Split(file)
		fname = fname[:len(fname)-3] // suffix .go は不要
		fmt.Printf("[%d] %14s:%4d: ", sage.cnt, fname, line)
		fmt.Printf(format, a...)
	}
}

func (sage *Sage) Inc(n int) int {
	sage.cnt += n
	return sage.cnt
}

func (sage *Sage) Close() error {
	sage.log("sage.Close() start\n")

	if sage.is_gc {
		sage.log("sage.Close(): tracemalloc.stop()\n")
		r := callFunction(sage.pstop)
		C.Py_DecRef(r)
	}

	C.Py_DecRef(sage.pModule)
	C.Py_DecRef(sage.pstop)
	C.Py_DecRef(sage.pGCD)
	C.Py_DecRef(sage.pResultant)
	C.Py_DecRef(sage.pDiscrim)
	C.Py_DecRef(sage.pPsc)
	C.Py_DecRef(sage.pSlope)
	C.Py_DecRef(sage.pGB)
	C.Py_DecRef(sage.pReduce)
	C.Py_DecRef(sage.pEval)
	C.Py_DecRef(sage.pValid)
	C.Py_DecRef(sage.pFactor)
	C.Py_DecRef(sage.pVarInit)

	ns := C.PyLong_FromLong(C.long(1))
	ret := callFunction(sage.pGC, ns)
	C.Py_DecRef(ret)
	C.Py_DecRef(sage.pGC)

	// C.Py_Finalize()
	sage.log("sage.Close() end\n")
	return nil
}

func NewSage(g *Ganrac, fname string) (*Sage, error) {
	C.Py_Initialize()

	var str string

	sage := new(Sage)
	sage.ganrac = g
	sage.is_gc = false
	sage.verbose = false

	/*
	   [<Signals.SIGINT: 2>, <cyfunction python_check_interrupt at 0x7f41066a0ee0>]
	   [8, <built-in function default_int_handler>]
	   [8, None]
	*/

	// TypeError: signal handler must be signal.SIG_IGN, signal.SIG_DFL, or a callable object
	str = `
import signal
try:
	signal.signal(8, signal.SIG_DFL)
except Exception:
	pass
import sage.all
import gc
import os
import tracemalloc

def gan_snapstart():
	print("tracemalloc.start(py)")
	tracemalloc.start()
	return

def gan_snapstop():
	tracemalloc.stop()
	return

def gan_gc(B):
	gc.collect()
	if B != 0:
		return "0"
	return str(sum([stat.size for stat in tracemalloc.take_snapshot().statistics('filename')]))

def varinit(num):
	d = {}
	for n in range(num):
		v = 'x' + str(n)
		d[v] = sage.all.var(v)
	return d

def polyringinit(varn):
	V = ['x' + str(i) for i in range(varn)]
	P = sage.all.PolynomialRing(sage.all.QQ, V)
	vardic = {V[i]: v for i, v in enumerate(P.gens())}
	return P, vardic

def gan_factor(varn: int, polystr: str):
	vardic = varinit(varn)
	F = sage.all.sage_eval(polystr, locals=vardic)
	G = F.factor_list()
	for i in range(len(G)):
		G[i] = list(G[i])
	return str(G)

def gan_gcd(varn: int, p, q):
	vardic = varinit(varn)
	F = sage.all.sage_eval(p, locals=vardic)
	G = sage.all.sage_eval(q, locals=vardic)
	return str(sage.all.gcd(F, G))

def gan_res(varn: int, p, q, x):
	vardic = varinit(varn)
	F = sage.all.sage_eval(p, locals=vardic)
	G = sage.all.sage_eval(q, locals=vardic)
	X = vardic[x]
	return str(F.resultant(G, X))

def gan_discrim(varn: int, p, x):
	P, vardic = polyringinit(varn)
	F = sage.all.sage_eval(p, locals=vardic)
	X = vardic[x]
	return str(F.discriminant(X))

def gan_psc(varn: int, p, q, x, J):
	P, vardic = polyringinit(varn)
	F = sage.all.sage_eval(p, locals=vardic)
	G = sage.all.sage_eval(q, locals=vardic)
	X = vardic['x' + str(x)]
	M = F.degree(X)
	N = G.degree(X)
	L = M+N-2*J;
	S = sage.all.Matrix(P, nrows=L, ncols=L)
	for D in range(M, -1, -1):
		AI = F.coefficient({X: D})
		for I in range(min([N - J, L + D - M - 1])):
			S[I, M-D+I] = AI
	for I in range(N - J - 1, max([0, (N-J-1)-J]) -1, -1):
		S[I, L-1] = F.coefficient({X: I-(N-J-1)+J});
	for D in range(N, -1, -1):
		BI = G.coefficient({X: D})
		for I in range(min([M - J, L + D - N - 1])):
			S[I+N-J, N-D+I] = BI
	for I in range(M - J - 1, max([0, M-J-1-J])-1, -1):
		S[I+N-J, L-1] = G.coefficient({X: I-(M-J-1)+J});

	return str(S.det())

def gan_comb(A, B):
	C = 1
	for I in range(1, B+1):
		C *= (A - I + 1)
		C //= I
	return C

def gan_slope(varn: int, p, q, x, K):
	P, vardic = polyringinit(varn)
	F = sage.all.sage_eval(p, locals=vardic)
	G = sage.all.sage_eval(q, locals=vardic)
	X = vardic['x' + str(x)]
	M = F.degree(X)
	N = G.degree(X)
	L = N - K
	S = sage.all.Matrix(P, nrows=L+1, ncols=L+1)
	CMK = gan_comb(M, K + 1)
	for J in range(L):
		S[0, J] = F.coefficient({X: M-J})
		for I in range(1, L-J):
			S[I, I+J] = S[0, J]
		if 0 <= L - J <= L:
			S[L-J, L] = (CMK - gan_comb(M-J, K+1)) * S[0, J]
		S[L, J] = G.coefficient({X: N - J})
	J = L
	if M-J >= 0:
		S[0, L] = (CMK - gan_comb(M - J, K + 1)) * F.coefficient({X: M-J})
	S[L, L] = CMK * G.coefficient({X: K})
	return str(S.det())

def gan_gb(polys, vars, n):
	"""
	https://doc.sagemath.org/html/en/reference/polynomial_rings/sage/rings/polynomial/multi_polynomial_ideal.html
	>>> P = sage.all.PolynomialRing(sage.all.QQ, ('x', 'y', 'z'))
	>>> x, y, z = P.gens()
	>>> I = sage.all.ideal(x**5 + y**4 + z**3 - 1,  x**3 + y**3 + z**2 - 1)
	>>> G = I.groebner_basis()
	>>> (x**2).reduce(G)
	x^2
	>>> (x**2).reduce(I)
	x^2
	"""
	V = vars
	if 0 < n < len(vars):
		order = f'degrevlex({n}),degrevlex({len(vars) - n})'
	else:
		order = 'degrevlex'
	if len(vars) == 1:
		V += ("y",)
	P = sage.all.PolynomialRing(sage.all.QQ, V, order=order)
	vardic = {V[i]: v for i, v in enumerate(P.gens())}
	F = sage.all.sage_eval(polys, locals=vardic)
	I = sage.all.ideal(F)
	G = I.groebner_basis()
	G = [g / g.content() for g in G]
	return str(G)

def gan_reduce(p, gb, vars, n):
	V = vars
	if 0 < n < len(vars):
		order = f'degrevlex({n}),degrevlex({len(vars) - n})'
	else:
		order = 'degrevlex'
	if len(vars) == 1:
		V += ("y",)
	try:
		P = sage.all.PolynomialRing(sage.all.QQ, V, order=order)
	except Exception as e:
		print(vars)
		print(V)
		raise

	vardic = {V[i]: v for i, v in enumerate(P.gens())}
	F = sage.all.sage_eval(p, locals=vardic)
	G = sage.all.sage_eval(gb, locals=vardic)
	I = sage.all.ideal(G)
	U = F.reduce(G)
	if U != 0:
		U /= U.content()
	return str(U)

def gan_eval(varn: int, s):
	return ""
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
	po := C.Py_CompileString(cstr, C.CString(fname), C.Py_file_input)

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

	sage.pModule = C.PyImport_ImportModule(pName)
	if sage.pModule == nil {
		return nil, fmt.Errorf("import module %s failed", mname)
	}

	var psnapstr *C.PyObject

	for _, tbl := range []struct {
		fname string
		pos   **C.PyObject
	}{
		{"gan_snapstart", &psnapstr},
		{"gan_snapstop", &sage.pstop},
		{"gan_gc", &sage.pGC},
		{"gan_gcd", &sage.pGCD},
		{"gan_res", &sage.pResultant},
		{"gan_discrim", &sage.pDiscrim},
		{"gan_psc", &sage.pPsc},
		{"gan_slope", &sage.pSlope},
		{"gan_gb", &sage.pGB},
		{"gan_reduce", &sage.pReduce},
		{"gan_eval", &sage.pEval},
		{"gan_factor", &sage.pFactor},
		{"varinit", &sage.pVarInit},
	} {
		if p, err := loadFunction(sage.pModule, tbl.fname); err != nil {
			return nil, err
		} else {
			*tbl.pos = p
		}
	}

	if sage.is_gc {
		fmt.Printf("NewSage() call tracemalloc.start()\n")
		r := callFunction(psnapstr)
		C.Py_DecRef(r)
	}

	return sage, nil
}

func (sage *Sage) GC() error {
	if !sage.is_gc {
		return nil
	}
	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	n, _ := runtime.MemProfile(nil, true)

	rets := "boo"
	if true {
		ns := C.PyLong_FromLong(C.long(0))
		ret := callFunction(sage.pGC, ns)
		defer C.Py_DecRef(ret)
		rets = toGoString(ret)
	}
	fmt.Printf("GC<%7d> py=%9s, go=%9d, %9d, %9d\n", sage.cnt, rets, ms.Alloc, ms.HeapAlloc, n)
	return nil
}

/* 変数の最大レベルを返す */
func (sage *Sage) varn(polys ...*Poly) int {
	n := 0
	for _, p := range polys {
		if n < int(p.Level())+1 {
			n = int(p.Level() + 1)
		}
	}
	return n + 1
}

func (sage *Sage) toGaNRAC(s string) interface{} {
	s = strings.ReplaceAll(s, "x", "$") + ";" // GaNRAC 用に変換
	u, err := sage.ganrac.Eval(strings.NewReader(s))
	if err != nil {
		fmt.Fprintf(os.Stderr, "sage.toGaNRAC(%s) failed: %s.", s, err.Error())
		return nil
	}
	return u
}

func (sage *Sage) EvalRObj(s string) RObj {
	u := sage.toGaNRAC(s)
	if u == nil {
		return nil
	}
	if uu, ok := u.(RObj); ok {
		return uu
	} else {
		fmt.Fprintf(os.Stderr, "sage.EvalRObj(%s) invalid.", s)
		return nil
	}
}

func (sage *Sage) EvalList(s string) *List {
	u := sage.toGaNRAC(s)
	if u == nil {
		return nil
	}
	if uu, ok := u.(*List); ok {
		return uu
	} else {
		fmt.Fprintf(os.Stderr, "sage.EvalList(%s) invalid.", s)
		return nil
	}
}

func (sage *Sage) Gcd(p, q *Poly) RObj {
	sage.log("Gcd(%s,%s) start!\n", p, q)
	varn := sage.varn(p, q)

	// 変数はすべて xi 形式にする
	ps := toPyString(fmt.Sprintf("%I", p))
	qs := toPyString(fmt.Sprintf("%I", q))
	// fmt.Printf("ps=%s qs=%p\n", fmt.Sprintf("%I", p), qs)

	ret := callFunctionv(sage.pGCD, varn, ps, qs)

	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pGCD\n", sage.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	if strings.HasPrefix(retstr, "none") {
		panic(fmt.Sprintf("<%d> stop..... varlist is broken. %s.", sage.cnt, retstr))
	}
	return sage.EvalRObj(retstr)
}

func (sage *Sage) Factor(q *Poly) *List {
	sage.cnt++
	sage.log("[%d] Factor(%s) start!\n", sage.cnt, q)
	varn := sage.varn(q)

	p, cont := q.PPC()
	ps := toPyString(fmt.Sprintf("%I", p)) // 変数はすべて xi 形式にする

	ret := callFunctionv(sage.pFactor, varn, ps)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "call object failed pFactor\n")
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	rets := toGoString(ret)
	uu := sage.EvalList(rets)
	if uu == nil {
		fmt.Fprintf(os.Stderr, "sage.Factor eval() failed: %s.\n", rets)
		return nil
	}
	// 第１項を数要素にする
	for i := 0; i < uu.Len(); i++ {
		vi, _ := uu.Geti(i)
		vi0, _ := vi.(*List).Geti(0)
		if _, ok := vi0.(*Poly); !ok {
			// 非多項式がみつかった
			if i != 0 {
				uu.Swap(0, i)
			}
			if !cont.IsOne() {
				v0, _ := uu.Geti(0)
				v00, _ := v0.(*List).Geti(0)
				v0.(*List).Seti(0, Mul(cont, v00.(RObj)))
			}
			return uu
		}
	}
	// 定数要素がなかったから追加する
	v := NewList()
	v.Append(cont)
	v.Append(NewInt(1))
	uu.Append(v)
	uu.Swap(0, uu.Len()-1)
	return uu
}

func (sage *Sage) Discrim(p *Poly, lv Level) RObj {
	sage.log("Discrim(%s) start!\n", p)
	varn := sage.varn(p)

	// 変数はすべて xi 形式にする
	ps := toPyString(fmt.Sprintf("%I", p))
	xs := toPyString(fmt.Sprintf("x%d", lv))

	ret := callFunctionv(sage.pDiscrim, varn, ps, xs)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pDiscrim\n", sage.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sage.EvalRObj(retstr)
}
func (sage *Sage) Resultant(p *Poly, q *Poly, lv Level) RObj {
	sage.log("Resultant(%s) start!\n", p)
	varn := sage.varn(p, q)

	// 変数はすべて xi 形式にする
	ps := toPyString(fmt.Sprintf("%I", p))
	qs := toPyString(fmt.Sprintf("%I", q))
	xs := toPyString(fmt.Sprintf("x%d", lv))

	ret := callFunctionv(sage.pResultant, varn, ps, qs, xs)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pResultant\n", sage.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sage.EvalRObj(retstr)
}

func (sage *Sage) Psc(p *Poly, q *Poly, lv Level, j int32) RObj {
	// sage.all.Matrix([[x,y],[z,w]]).det()
	varn := sage.varn(p, q)
	ps := toPyString(fmt.Sprintf("%I", p))
	qs := toPyString(fmt.Sprintf("%I", q))
	lvs := C.PyLong_FromLong(C.long(lv))
	js := C.PyLong_FromLong(C.long(j))

	ret := callFunctionv(sage.pPsc, varn, ps, qs, lvs, js)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pPsc\n", sage.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sage.EvalRObj(retstr)
}

func (sage *Sage) Slope(p *Poly, q *Poly, lv Level, k int32) RObj {
	sage.log("Slope(%s) start!\n", p)
	varn := sage.varn(p, q)
	ps := toPyString(fmt.Sprintf("%I", p))
	qs := toPyString(fmt.Sprintf("%I", q))
	lvs := C.PyLong_FromLong(C.long(lv))
	ks := C.PyLong_FromLong(C.long(k))

	ret := callFunctionv(sage.pSlope, varn, ps, qs, lvs, ks)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pSlope\n", sage.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sage.EvalRObj(retstr)
}

func (sage *Sage) toPyVars(vars *List) *C.PyObject {
	vs := C.PyTuple_New(C.long(vars.Len()))
	for i := 0; i < vars.Len(); i++ {
		v, _ := vars.Geti(i)
		C.PyTuple_SetItem(vs, C.long(i), toPyString(fmt.Sprintf("x%d", v.(*Poly).Level())))
	}
	return vs
}

func (sage *Sage) GB(p *List, vars *List, n int) *List {
	sage.cnt++
	sage.log("[%d] GB(%s,%s,%d) start!\n", sage.cnt, p, vars, n)
	// 変数はすべて xi 形式にする
	ps := toPyString(fmt.Sprintf("%I", p))
	ns := C.PyLong_FromLong(C.long(n))

	vs := sage.toPyVars(vars)
	defer C.Py_DecRef(vs)

	ret := callFunction(sage.pGB, ps, vs, ns)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pGB\n", sage.cnt)
		C.PyErr_Print()
		return nil
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sage.EvalList(retstr)
}

func (sage *Sage) Reduce(p *Poly, gb *List, vars *List, n int) (RObj, bool) {
	sage.cnt++
	sage.log("[%d] Reduce(%s,%s,%s,%d) start!\n", sage.cnt, p, gb, vars, n)
	ps := toPyString(fmt.Sprintf("%I", p))
	gbs := toPyString(fmt.Sprintf("%I", gb))
	vs := sage.toPyVars(vars)
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

	ret := callFunction(sage.pReduce, ps, gbs, vs, ns)
	if ret == nil {
		fmt.Fprintf(os.Stderr, "<%d> call object failed pReduce\n", sage.cnt)
		C.PyErr_Print()
		return p, false
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	return sage.EvalRObj(retstr), false
}

func (sage *Sage) Eval(p string) (GObj, error) {
	sage.log("Eval(%s) start!\n", p)
	ps := toPyString(p)
	ret := callFunction(sage.pEval, ps) // 変数設定は??
	if ret == nil {
		C.PyErr_Print()
		return nil, fmt.Errorf("<%d> call object failed pEval\n", sage.cnt)
	}
	defer C.Py_DecRef(ret)

	retstr := toGoString(ret)
	s := sage.toGaNRAC(retstr)
	if s == nil {
		return nil, fmt.Errorf("Sage: returns null")
	}
	if sg, ok := s.(GObj); ok {
		return sg, nil
	} else {
		return nil, fmt.Errorf("Sage: unsupported " + retstr)
	}
}

func (sage *Sage) SetLogger(logger *log.Logger) {
}
