package sympy

// #include <Python.h>
import "C"

// export CFLAGS := -I/usr/include/python3.10/
// export LDFLAGS := -L/usr/lib/x86_64-linux-gnu/ -lpython3.10

import (
	"fmt"
	"github.com/hiwane/ganrac"
	"unsafe"
)

/*
 * pModule 内の funname へのポインタを返す.
 */
func loadFunction(pModule *C.PyObject, funname string) (*C.PyObject, error) {
	cstr := C.CString(funname)
	defer C.free(unsafe.Pointer(cstr))
	pFunc := C.PyObject_GetAttrString(pModule, cstr)
	if pFunc == nil {
		return nil, fmt.Errorf("GetAttrString(%s) faield", funname)
	}
	if C.PyCallable_Check(pFunc) <= 0 {
		return nil, fmt.Errorf("Callable_Check(%s) failed", funname)
	}

	return pFunc, nil
}

/*
 * 関数 pFunc を実行する
 */
func callFunction(pFunc *C.PyObject, argv ...*C.PyObject) *C.PyObject {
	atuple := C.PyTuple_New(C.long(len(argv)))
	for i, v := range argv {
		C.PyTuple_SetItem(atuple, C.long(i), v)
	}
	defer C.Py_DecRef(atuple)
	return C.PyObject_CallObject(pFunc, atuple)
}

func callFunctionv(pFunc *C.PyObject, v int, argv ...*C.PyObject) *C.PyObject {
	atuple := C.PyTuple_New(C.long(len(argv) + 1))

	C.PyTuple_SetItem(atuple, C.long(0), C.PyLong_FromLong(C.long(v)))
	for i, v := range argv {
		C.PyTuple_SetItem(atuple, C.long(i+1), v)
	}
	defer C.Py_DecRef(atuple)
	return C.PyObject_CallObject(pFunc, atuple)
}

/*
 * python 文字列へ変換する
 */
func toPyString(s string) *C.PyObject {
	p := C.CString(s)
	defer C.free(unsafe.Pointer(p))
	return C.PyUnicode_FromString(p)
}

// 変数はすべて xi 形式にする
func polyToPyString(p *ganrac.Poly) *C.PyObject {
	s := fmt.Sprintf("%"+string(ganrac.FORMAT_SYMPY), p)
	return toPyString(s)
}

func listToPyString(p *ganrac.List) *C.PyObject {
	s := fmt.Sprintf("%"+string(ganrac.FORMAT_SYMPY), p)
	return toPyString(s)
}

func lvToPyString(lv ganrac.Level) *C.PyObject {
	s := fmt.Sprintf("x%d", lv)
	return toPyString(s)
}

func toGoString(s *C.PyObject) string {
	cutf8 := C.PyUnicode_AsUTF8(s)
	return C.GoString(cutf8)
}
