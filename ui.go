package ganrac

import (
// "math/big"
)

func (p Uint) zero() NObj {
	return Uint(0)
}

func (p Uint) toFF(a NObj) NObj {
	switch a := a.(type) {
	case Uint:
		return Uint(a % p)
	case *Int:
		return a
	}
	return a
}

func (p Uint) negMod(aa NObj) NObj {
	switch a := aa.(type) {
	case Uint:
		return Uint((p - a) % p)
	}
	return nil
}

// 入力 は Mod p 上な要素と仮定
func (p Uint) addMod(aa, bb NObj) NObj {
	switch a := aa.(type) {
	case Uint:
		switch b := bb.(type) {
		case Uint:
			return Uint((a + b) % p)
		}
	}
	return nil
}

func (p Uint) subMod(aa, bb NObj) NObj {
	switch a := aa.(type) {
	case Uint:
		switch b := bb.(type) {
		case Uint:
			if a > b {
				return Uint(a - b)
			} else {
				return Uint(p - b + a)
			}
		}
	}
	return nil
}

func (p Uint) mulMod(aa, bb NObj) NObj {
	switch a := aa.(type) {
	case Uint:
		switch b := bb.(type) {
		case Uint:
			return Uint((uint64(a) * uint64(b)) % uint64(p))
		}
	}
	return nil
}

func (p Uint) validMod(f NObj) error {
	return nil
}
