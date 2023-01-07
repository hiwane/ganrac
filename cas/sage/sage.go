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
	"os"
	"runtime"
	"strings"
	"unsafe"
)

type Sage struct {
	pModule  *C.PyObject
	pstop    *C.PyObject
	pGC      *C.PyObject
	pGCD     *C.PyObject
	pValid   *C.PyObject
	pFactor  *C.PyObject
	pVarInit *C.PyObject
	ganrac   *Ganrac

	varnum  int

	cnt   int
	is_gc bool
}

func (sage *Sage) Inc(n int) int {
	sage.cnt += n
	return sage.cnt
}

func (sage *Sage) Close() error {
	fmt.Printf("sage.Close(): start\n")

	if sage.is_gc {
		fmt.Printf("sage.Close(): tracemalloc.stop()\n")
		r := callFunction(sage.pstop)
		C.Py_DecRef(r)
	}

	// C.Py_Finalize()
	// fmt.Printf("sage.Close(): end finalize()\n")
	fmt.Printf("sage.Close(): end pstop()\n")
	return nil
}

func NewSage(g *Ganrac, fname string) (*Sage, error) {
	C.Py_Initialize()

	var str string

	sage := new(Sage)
	sage.ganrac = g
	sage.is_gc = true

	/*
	   [<Signals.SIGINT: 2>, <cyfunction python_check_interrupt at 0x7f41066a0ee0>]
	   [8, <built-in function default_int_handler>]
	   [8, None]
	*/

	// TypeError: signal handler must be signal.SIG_IGN, signal.SIG_DFL, or a callable object
	str = `
import signal
signal.signal(8, signal.SIG_DFL)
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

def gan_gc():
	gc.collect()
	return str(sum([stat.size for stat in tracemalloc.take_snapshot().statistics('filename')]))

def varinit(num):
	d = {}
	for n in range(num):
		v = 'x' + str(n)
		d[v] = sage.all.var(v)
	return d

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
			fmt.Printf("errrrro on PyImport_ExecCodeModule()\n")
			C.PyErr_Print()
			return nil, fmt.Errorf("PyImport_ExecCodeModule failed")
		}
	}

	fmt.Printf("pname=%v: %s, po=%p\n", pName, mname, po)
	sage.pModule = C.PyImport_ImportModule(pName)
	if sage.pModule == nil {
		return nil, fmt.Errorf("import module %s failed", mname)
	}

	var psnapstr *C.PyObject

	fmt.Printf("loadFunction\n")
	for _, tbl := range []struct {
		fname string
		pos   **C.PyObject
	}{
		{"gan_snapstart", &psnapstr},
		{"gan_snapstop", &sage.pstop},
		{"gan_gc", &sage.pGC},
		{"gan_gcd", &sage.pGCD},
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

	fmt.Printf("NewSage() end\n")
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
		ret := callFunction(sage.pGC)
		defer C.Py_DecRef(ret)
		rets = toGoString(ret)
	}
	fmt.Printf("GC<%5d> py=%9s, go=%9d, %9d, %9d\n", sage.cnt, rets, ms.Alloc, ms.HeapAlloc, n)
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

	// fmt.Printf("Gcd(%s,%s) start!\n", p, q)
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
	return NewInt(1)
}
func (sage *Sage) Resultant(p *Poly, q *Poly, lv Level) RObj {
	return NewInt(1)
}
func (sage *Sage) Psc(p *Poly, q *Poly, lv Level, j int32) RObj {
	return NewInt(1)
}
func (sage *Sage) Sres(p *Poly, q *Poly, lv Level, k int32) RObj {
	return NewInt(1)
}
func (sage *Sage) GB(p *List, vars *List, n int) *List {
	return NewList()
}
func (sage *Sage) Reduce(p *Poly, gb *List, vars *List, n int) (RObj, bool) {
	return NewInt(1), true
}
func (sage *Sage) Eval(p string) (GObj, error) {
	return nil, fmt.Errorf("unsupported")
}
