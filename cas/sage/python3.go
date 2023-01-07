package sage

// #include <Python.h>
// #cgo CFLAGS: -I/usr/include/python3.10/
// #cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu/ -lpython3.10
import "C"

import (
	"fmt"
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

func toGoString(s *C.PyObject) string {
	cutf8 := C.PyUnicode_AsUTF8(s)
	return C.GoString(cutf8)
}
