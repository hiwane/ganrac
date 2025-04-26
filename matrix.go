package ganrac

import (
	"fmt"
	"github.com/hiwane/ganrac/cache"
)

type Matrix struct {
	row, col int
	m        []RObj
}

func NewMatrix(rows, cols int) *Matrix {
	if rows <= 0 || cols <= 0 {
		return nil
	}
	mat := Matrix{row: rows, col: cols, m: make([]RObj, rows*cols)}
	return &mat
}

func (m *Matrix) Get(i, j int) (RObj, bool) {
	if i < 0 || i >= m.row || j < 0 || j >= m.col {
		return nil, false
	}
	return m.get(i, j), true
}

func (m *Matrix) get(i, j int) RObj {
	// fmt.Printf("get(%d/%d, %d/%d)\n", i, m.row, j, m.col)
	return m.m[i*m.col+j]
}

func (m *Matrix) Set(i, j int, v RObj) bool {
	if i < 0 || i >= m.row || j < 0 || j >= m.col {
		return false
	}
	m.set(i, j, v)
	return true
}

func (m *Matrix) set(i, j int, v RObj) {
	// fmt.Printf("set(%d/%d, %d/%d, %v)\n", i, m.row, j, m.col, v)
	m.m[i*m.col+j] = v
}

func (m *Matrix) Determinant() (RObj, error) {
	// Check if the matrix is square
	if m.row != m.col {
		return nil, fmt.Errorf("Matrix must be square to calculate determinant; got %dx%d", m.row, m.col)
	}

	var lru cache.Cacher[RObj]
	if m.row <= 3 {
		lru = &cache.NoCache[RObj]{}
	} else {
		lru = cache.NewLRUCache[RObj](100)
	}

	d := m.det(0, lru)
	// fmt.Printf("Det(%dx%d) = %s\n", m.row, m.col, lru)
	return d, nil
}

func (m *Matrix) det(usedMask cache.Hash, cache cache.Cacher[RObj]) RObj {
	// fmt.Printf("        det(%#x=%d,%v)\n", uint64(usedMask), uint64(usedMask), m)

	// Base case for 2x2 matrix
	if m.row == 2 {
		return Sub(Mul(m.get(0, 0), m.get(1, 1)), Mul(m.get(0, 1), m.get(1, 0)))
	}

	//////////////////////////////////////
	// Recursive case for larger matrices
	//////////////////////////////////////
	rows := make([]int, m.row)
	cols := make([]int, m.col)

	// 展開する行・列を決定する. なるべくゼロが多い行・列を選ぶ
	maxr := 0
	rowi := 0
	maxc := 0
	coli := 0
	for i := 0; i < m.row; i++ {
		for j := 0; j < m.col; j++ {
			if m.get(i, j).IsZero() {
				rows[i]++
				cols[j]++
				if maxr < rows[i] {
					maxr = rows[i]
					rowi = i
				}
				if maxc < cols[j] {
					maxc = cols[j]
					coli = j
				}
			}
		}
	}

	if maxr > 0 {
		return m.detRow(usedMask, rowi, cache)
	} else {
		return m.detCol(usedMask, coli, cache)
	}
}

// 32次正方行列まで対応. それ以上はそもそも計算できないでしょう
// rowcol: 1: row, 0: col
// pos: 小行列の取り除く場所
// usedMask: 元の行列での取り除いた場所を表すビットマスク
//
//	pos は 1 行目を取り除いたあとに 3 行目を取り除く場合は，元の行列の４行目を表している点に注意
func (m *Matrix) setUsedMask(usedMask cache.Hash, _pos, rowcol int) cache.Hash {
	c := uint(rowcol)
	pos := _pos
	mask := uint64(usedMask)
	for c = uint(rowcol); c < 64; c += 2 {
		if (mask & (1 << uint(c))) == 0 {
			if pos == 0 {
				break
			}
			pos--
		}
	}

	ret := cache.Hash(mask | (1 << uint(c)))
	// fmt.Printf("           setUsedMask: ipt=%#x, pos=%d, rowcol=%d => c=%d: ret=%#x\n", uint64(usedMask), _pos, rowcol, c, uint64(ret))
	if ret.NumSetBits(0, 1) != 1+usedMask.NumSetBits(0, 1) {
		m.printMinor(-1, -1, usedMask)
		fmt.Printf("A:input=%#x:%d, ret=%#x:%d\n", uint64(usedMask), usedMask.NumSetBits(0, 1), uint64(ret), ret.NumSetBits(0, 1))
		panic("stopA")
	}
	if ret.NumSetBits(rowcol, 2) != 1+usedMask.NumSetBits(rowcol, 2) {
		fmt.Printf("B:input=%#x:%d, ret=%#x:%d\n", uint64(usedMask), usedMask.NumSetBits(0, 1), uint64(ret), ret.NumSetBits(0, 1))
		panic("stopB")
	}

	return ret
}

func (m *Matrix) validUsedMask(usedMask cache.Hash) bool {
	n := usedMask.NumSetBits(0, 1)
	if n%2 != 0 {
		panic(fmt.Sprintf("stop1; %x", usedMask))
	}
	n0 := usedMask.NumSetBits(0, 2)
	if n0*2 != n {
		panic(fmt.Sprintf("stop2; %x", usedMask))
	}
	return true
}

func (m *Matrix) printMinor(i, j int, usedMask cache.Hash) {
	fmt.Printf("%#02x  ", uint64(usedMask))
	for c := m.row; c < 4; c++ {
		fmt.Printf("  ")
	}
	fmt.Printf("(%d,%d)\t\tprintMinor()\n", i, j)
}

// 行 row で小行列展開
func (m *Matrix) detRow(_usedMask cache.Hash, row int, lcache cache.Cacher[RObj]) RObj {
	sign := 1
	if row%2 == 0 {
		sign = -1
	}
	usedMask := m.setUsedMask(_usedMask, row, 1)

	var det RObj = zero
	for j := 0; j < m.col; j++ {
		sign *= -1
		s := m.get(row, j)
		if s.IsZero() {
			continue
		}
		subMat := m.submatrix(row, j)
		usedMaskC := m.setUsedMask(usedMask, j, 0)
		// m.printMinor(row, j, usedMaskC)
		m.validUsedMask(usedMaskC)
		var subDet RObj
		if v, ok := lcache.Get(usedMaskC); ok {
			subDet = v
		} else {
			subDet = subMat.det(usedMaskC, lcache)
			lcache.Put(subMat, subDet)
		}

		s = Mul(s, subDet)
		if sign > 0 {
			det = Add(det, s)
		} else {
			det = Sub(det, s)
		}
	}
	return det
}

// 列 col で小行列展開
func (m *Matrix) detCol(usedMask cache.Hash, col int, lcache cache.Cacher[RObj]) RObj {
	sign := 1
	if col%2 == 0 {
		sign = -1
	}
	usedMask = m.setUsedMask(usedMask, col, 0)
	var det RObj = zero
	for i := 0; i < m.row; i++ {
		sign *= -1
		s := m.get(i, col)
		if s.IsZero() {
			continue
		}
		subMat := m.submatrix(i, col)
		usedMaskR := m.setUsedMask(usedMask, i, 1)
		// m.printMinor(i, col, usedMaskR)
		m.validUsedMask(usedMaskR)
		var subDet RObj
		if v, ok := lcache.Get(usedMaskR); ok {
			subDet = v
		} else {
			subDet = subMat.det(usedMaskR, lcache)
			lcache.Put(usedMaskR, subDet)
		}

		s = Mul(s, subDet)
		if sign > 0 {
			det = Add(det, s)
		} else {
			det = Sub(det, s)
		}
	}
	return det
}

func (m *Matrix) submatrix(row, col int) *Matrix {
	// Create a new matrix with one less row and column
	// suppose 0 <= row < m.row
	// suppose 0 <= col < m.col
	sub := NewMatrix(m.row-1, m.col-1)
	ni := 0
	for i := 0; i < m.row; i++ {
		if i == row {
			ni = 1
			continue
		}
		newRow := i - ni
		nj := 0
		for j := 0; j < m.col; j++ {
			if j == col {
				nj = 1
				continue
			}
			newCol := j - nj
			sub.set(newRow, newCol, m.get(i, j))
		}
	}
	return sub
}

func (m *Matrix) String() string {
	str := ""
	sep0 := "["
	for i := 0; i < m.row; i++ {
		sep := "["
		str += sep0
		sep0 = ","
		for j := 0; j < m.col; j++ {
			v := m.get(i, j)
			if v == nil {
				str += sep + "nil"
			} else {
				str += fmt.Sprintf("%s%s", sep, m.get(i, j))
			}
			sep = ","
		}
		str += "]"
	}
	return str + "]"

}

func (m *Matrix) Hash() cache.Hash {
	h := cache.Hash(0)
	for i := 0; i < m.row; i++ {
		for j := 0; j < m.col; j++ {
			h = h*31 + m.get(i, j).Hash()
		}
	}
	return h
}

func (m *Matrix) Equals(other any) bool {
	if o, ok := other.(*Matrix); ok {
		if m.row != o.row || m.col != o.col {
			return false
		}
		for i := 0; i < m.row; i++ {
			for j := 0; j < m.col; j++ {
				if !m.get(i, j).Equals(o.get(i, j)) {
					return false
				}
			}
		}
		return true
	}
	return false
}
