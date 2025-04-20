package ganrac

import (
	"testing"
	// "fmt"
)

func TestMatDeterminant(t *testing.T) {
	funcname := "TestMatDeterminant"

	cache := &NoCache[RObj]{}

	for ii, tbl := range []struct {
		matrix [][]int64
		det    int64
	}{
		{
			matrix: [][]int64{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			det: 0,
		},
		{
			matrix: [][]int64{
				{6, 1, 1},
				{4, -2, 5},
				{2, 8, 7},
			},
			det: -306,
		},
		{
			matrix: [][]int64{
				{3, 0, 2},
				{2, 0, -2},
				{0, 1, 1},
			},
			det: 10,
		},
		{
			matrix: [][]int64{
				{1, 2, 3},
				{0, 1, 4},
				{5, 6, 0},
			},
			det: 1,
		},
		{
			matrix: [][]int64{
				{2, -3, 1},
				{2, 0, -1},
				{1, 4, 5},
			},
			det: 49,
		},
		// ここから４字行列
		{
			matrix: [][]int64{
				{1, 2, 3, 4},
				{5, 6, 7, 8},
				{9, 1, 2, 3},
				{4, 5, 6, 7},
			},
			det: 0,
		},
		{
			matrix: [][]int64{
				{1, 0, 0, 0},
				{0, 1, 0, 0},
				{0, 0, 1, 0},
				{0, 0, 0, 2},
			},
			det: 2,
		},
		{
			matrix: [][]int64{
				{2, 3, 1, 5},
				{4, 1, 3, 2},
				{3, 2, 4, 1},
				{1, 5, 2, 3},
			},
			det: 47,
		},
		{
			matrix: [][]int64{
				{0, 2, 1, 3},
				{4, 0, 2, 1},
				{3, 1, 0, 2},
				{1, 4, 3, 0},
			},
			det: -101,
		},
		{
			matrix: [][]int64{
				{1, 2, 3, 0},
				{0, 1, 4, 0},
				{0, 0, 1, 5},
				{0, 0, 0, 1},
			},
			det: 1,
		},
		{
			matrix: [][]int64{
				{1, 2, 3, 4},
				{0, 6, 7, 8},
				{9, 0, 2, 3},
				{4, 5, 0, 7},
			},
			det: -339,
		},
		{
			matrix: [][]int64{
				{2, 0, 1, 5},
				{4, 1, 0, 2},
				{3, 2, 4, 0},
				{0, 5, 2, 3},
			},
			det: 355,
		},
	} {
		m := NewMatrix(len(tbl.matrix), len(tbl.matrix[0]))
		for i := 0; i < len(tbl.matrix); i++ {
			for j := 0; j < len(tbl.matrix[i]); j++ {
				m.Set(i, j, NewInt(tbl.matrix[i][j]))
			}
		}

		// fmt.Printf("%s(%d): mat=%v\n", funcname, ii, m)
		det := NewInt(tbl.det)

		d, err := m.Determinant()
		if err != nil {
			t.Errorf("%s: error calculating determinant: %v", funcname, err)
			continue
		}
		if !d.Equals(det) {
			t.Errorf("%s(%d): expected determinant %v, got %v, mat=%v", funcname, ii, det, d, m)
			continue
		}

		for i := 0; i < m.row; i++ {
			d = m.detRow(0, i, cache)
			if !d.Equals(det) {
				t.Errorf("%s(%d,detRow%d): expected determinant %v, got %v, mat=%v", funcname, ii, i, det, d, m)
				continue
			}

			d = m.detCol(0, i, cache)
			if !d.Equals(det) {
				t.Errorf("%s(detRow%d): expected determinant %v, got %v", funcname, i, det, d)
				continue
			}
		}
	}
}

func TestMatrixSetUsedMask(t *testing.T) {
	m := NewMatrix(3, 3)

	for rowcol := 0; rowcol < 1; rowcol++ {
		mask1 := Hash(0)

		mask1 = m.setUsedMask(mask1, 2, rowcol)
		if x := mask1.NumSetBits(0, 1); x != 1 {
			t.Errorf("A1: Expected mask1 (%d)  x=%d", mask1, x)
			continue
		}
		if x := mask1.NumSetBits(rowcol, 1); x != 1 {
			t.Errorf("A2: Expected mask1 (%d)  x=%d", mask1, x)
			continue
		}

		mask1 = m.setUsedMask(mask1, 0, rowcol)
		if x := mask1.NumSetBits(0, 1); x != 2 {
			t.Errorf("A3: Expected mask1 (%d)  x=%d", mask1, x)
			continue
		}
		if x := mask1.NumSetBits(rowcol, 1); x != 2 {
			t.Errorf("A4: Expected mask1 (%d)  x=%d", mask1, x)
			continue
		}

		mask2 := Hash(0)
		mask2 = m.setUsedMask(mask2, 0, rowcol)
		if x := mask2.NumSetBits(0, 1); x != 1 {
			t.Errorf("A1: Expected mask2 (%d)  x=%d", mask2, x)
			continue
		}
		if x := mask2.NumSetBits(rowcol, 1); x != 1 {
			t.Errorf("A2: Expected mask2 (%d)  x=%d", mask2, x)
			continue
		}

		mask2 = m.setUsedMask(mask2, 1, rowcol)
		if x := mask2.NumSetBits(0, 1); x != 2 {
			t.Errorf("A3: Expected mask2 (%d)  x=%d", mask2, x)
			continue
		}
		if x := mask2.NumSetBits(rowcol, 1); x != 2 {
			t.Errorf("A4: Expected mask2 (%d)  x=%d", mask2, x)
			continue
		}

		if mask1 != mask2 {
			t.Errorf("A: Expected mask1 (%d) to equal mask2 (%d)", mask1, mask2)
			continue
		}

		mask1 = Hash(0)
		mask1 = m.setUsedMask(mask1, 1, rowcol)
		mask1 = m.setUsedMask(mask1, 0, rowcol)

		mask2 = Hash(0)
		mask2 = m.setUsedMask(mask2, 0, rowcol)
		mask2 = m.setUsedMask(mask2, 0, rowcol)

		if mask1 != mask2 {
			t.Errorf("B: Expected mask1 (%d) to equal mask2 (%d)", mask1, mask2)
			continue
		}
	}

}
