package cas

import (
	"fmt"
	"github.com/hiwane/ganrac"
	"strings"
)

func evalstr(g *ganrac.Ganrac, s string) (interface{}, error) {
	if !strings.HasSuffix(s, ";") {
		s += ";"
	}
	return g.Eval(strings.NewReader(s))
}

func str2list(g *ganrac.Ganrac, s string) (*ganrac.List, error) {
	x, err := evalstr(g, s)
	if err != nil {
		return nil, err
	}

	y, ok := x.(*ganrac.List)
	if !ok {
		return nil, fmt.Errorf("not a List")
	}
	return y, nil
}

func str2poly(g *ganrac.Ganrac, s string) (*ganrac.Poly, error) {
	x, err := evalstr(g, s)
	if err != nil {
		return nil, err
	}

	y, ok := x.(*ganrac.Poly)
	if !ok {
		return nil, fmt.Errorf("not a poly")
	}
	return y, nil
}

func str2fof(g *ganrac.Ganrac, s string) (ganrac.Fof, error) {
	x, err := evalstr(g, s)
	if err != nil {
		return nil, err
	}

	y, ok := x.(ganrac.Fof)
	if !ok {
		return nil, fmt.Errorf("not a formula")
	}
	return y, nil
}

func str2robj(g *ganrac.Ganrac, s string) (ganrac.RObj, error) {
	x, err := evalstr(g, s)
	if err != nil {
		return nil, err
	}

	y, ok := x.(ganrac.RObj)
	if !ok {
		return nil, fmt.Errorf("not an RObj")
	}
	return y, nil
}
