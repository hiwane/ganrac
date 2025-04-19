package ganrac

import (
	"fmt"
	"testing"
)

type testHash struct {
	v uint
}

func (h *testHash) Hash() Hash {
	return Hash(h.v % 10)
}

func (h *testHash) Equals(v any) bool {
	if h2, ok := v.(*testHash); ok {
		return h.v == h2.v
	}
	return false
}

func (h *testHash) String() string {
	return fmt.Sprintf("testHash(%d)", h.v)
}

func TestLRUCache(t *testing.T) {

	lru := NewLRUCache[string](3)

	////////////////////////////////////////
	// 1つ目の要素を追加
	////////////////////////////////////////
	key := &testHash{v: 1}
	value := "a"

	if expectLen := 0; lru.Len() != lru.lenMapList() && lru.Len() != expectLen {
		t.Errorf("Expected size %d, got Len=%d, LenMapList=%d", expectLen, lru.Len(), lru.lenMapList())
		return
	}
	if _, ok := lru.Get(key); ok {
		t.Errorf("Expected key %v to not exist", key)
		return
	}

	lru.Put(key, value)

	if expectLen := 1; lru.Len() != lru.lenMapList() && lru.Len() != expectLen {
		t.Errorf("Expected size %d, got Len=%d, LenMapList=%d", expectLen, lru.Len(), lru.lenMapList())
		return
	}

	if v, ok := lru.Get(key); !ok || v != value {
		t.Errorf("Expected key %v; ok=%v value=%v, v=%v", key, ok, value, v)
		return
	}

	if questr := "1 "; lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}

	////////////////////////////////////////
	// 2つ目の要素を追加
	////////////////////////////////////////
	key2 := &testHash{v: 2}
	value2 := "b"

	if expectLen := 1; lru.Len() != lru.lenMapList() && lru.Len() != expectLen {
		t.Errorf("Expected size %d, got Len=%d, LenMapList=%d", expectLen, lru.Len(), lru.lenMapList())
		return
	}
	if _, ok := lru.Get(key2); ok {
		t.Errorf("Expected key %v to not exist", key2)
		return
	}

	lru.Put(key2, value2)

	if questr := "2 1 "; lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}

	if expectLen := 2; lru.Len() != lru.lenMapList() && lru.Len() != expectLen {
		t.Errorf("Expected size %d, got Len=%d, LenMapList=%d", expectLen, lru.Len(), lru.lenMapList())
		return
	}

	if v, ok := lru.Get(key); !ok || v != value {
		t.Errorf("1 Expected key %v; ok=%v value=%v, v=%v", key, ok, value, v)
		return
	}
	if questr := "1 2 "; lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}
	if v, ok := lru.Get(key2); !ok || v != value2 {
		t.Errorf("1 Expected key %v; ok=%v value=%v, v=%v", key2, ok, value2, v)
		return
	}
	if questr := "2 1 "; lru.queueString() != questr { // 2 を GET したので，2 が先頭に来る
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}

	////////////////////////////////////////
	// 3つ目の要素を追加
	////////////////////////////////////////
	key3 := &testHash{v: 3}
	value3 := "c"

	if expectLen := 2; lru.Len() != lru.lenMapList() && lru.Len() != expectLen {
		t.Errorf("Expected size %d, got Len=%d, LenMapList=%d", expectLen, lru.Len(), lru.lenMapList())
		return
	}
	if _, ok := lru.Get(key3); ok {
		t.Errorf("Expected key %v to not exist", key3)
		return
	}

	lru.Put(key3, value3)

	if questr := "3 2 1 "; lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}

	if expectLen := 3; lru.Len() != lru.lenMapList() && lru.Len() != expectLen {
		t.Errorf("Expected size %d, got Len=%d, LenMapList=%d", expectLen, lru.Len(), lru.lenMapList())
		return
	}

	if v, ok := lru.Get(key3); !ok || v != value3 {
		t.Errorf("1 Expected key %v; ok=%v value=%v, v=%v", key3, ok, value3, v)
		return
	}
	if questr := "3 2 1 "; lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}
	if v, ok := lru.Get(key2); !ok || v != value2 {
		t.Errorf("1 Expected key %v; ok=%v value=%v, v=%v", key2, ok, value2, v)
		return
	}
	if questr := "2 3 1 "; lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}
	if v, ok := lru.Get(key); !ok || v != value {
		t.Errorf("1 Expected key %v; ok=%v value=%v, v=%v", key, ok, value, v)
		return
	}
	if questr := "1 2 3 "; lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}

	////////////////////////////////////////
	// 4つ目の要素を追加：キャパシティを超える
	// 直前に，key3, key2, key の順で
	// Get() しているので， key3 が消える
	////////////////////////////////////////
	key4 := &testHash{v: 4}
	value4 := "d"

	if expectLen := 3; lru.Len() != lru.lenMapList() && lru.Len() != expectLen {
		t.Errorf("Expected size %d, got Len=%d, LenMapList=%d", expectLen, lru.Len(), lru.lenMapList())
		return
	}

	lru.Put(key4, value4)
	if expectLen := 3; lru.Len() != lru.lenMapList() && lru.Len() != expectLen {
		t.Errorf("Expected size %d, got Len=%d, LenMapList=%d", expectLen, lru.Len(), lru.lenMapList())
		return
	}

	if v, ok := lru.Get(key4); !ok || v != value4 {
		t.Errorf("1 Expected key %v; ok=%v value=%v, v=%v", key4, ok, value4, v)
		return
	}
	if questr := "4 1 2 "; lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}
	if v, ok := lru.Get(key3); ok {
		fmt.Printf("Cache=%v\nQueue=", lru)
		t.Errorf("1 Expected key %v; ok=%v value=%v, v=%v", key3, ok, value3, v)
		return
	}
	if v, ok := lru.Get(key2); !ok || v != value2 {
		fmt.Printf("Cache=%v\n", lru)
		t.Errorf("1 Expected key %v; ok=%v value=%v, v=%v", key2, ok, value2, v)
		return
	}
	if questr := "2 4 1 "; lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}
	if v, ok := lru.Get(key); !ok || v != value {
		t.Errorf("1 Expected key %v; ok=%v value=%v, v=%v", key, ok, value, v)
		return
	}
	if questr := "1 2 4 "; lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}

	value4 += "newnew"
	lru.Put(key4, value4)
	if expectLen := 3; lru.Len() != lru.lenMapList() && lru.Len() != expectLen {
		t.Errorf("Expected size %d, got Len=%d, LenMapList=%d", expectLen, lru.Len(), lru.lenMapList())
		return
	}

	if v, ok := lru.Get(key4); !ok || v != value4 {
		t.Errorf("1 Expected key %v; ok=%v value=%v, v=%v", key4, ok, value4, v)
		return
	}
	if questr := "4 1 2 "; lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}

}

// ハッシュキーが一致する場合のテスト
func TestLRUCacheMap(t *testing.T) {

	type Pair struct {
		value string
		key   *testHash
	}

	lru := NewLRUCache[string](10)

	vx := []Pair{
		{"a0", &testHash{v: 1}},
		{"a1", &testHash{v: 11}},
		{"a2", &testHash{v: 21}},
		{"a3", &testHash{v: 31}},
		{"a4", &testHash{v: 41}},
		{"a5", &testHash{v: 51}},
	}

	// 1 を追加
	questr := ""
	for _, pair := range vx {
		lru.Put(pair.key, pair.value)
		questr = "1 " + questr
	}

	if lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}

	v2 := []Pair{
		{"b0", &testHash{v: 2}},
		{"b1", &testHash{v: 12}},
		{"b2", &testHash{v: 22}},
		{"b3", &testHash{v: 32}},
		{"b8", &testHash{v: 82}},
		{"b9", &testHash{v: 92}},
	}

	// 2 を追加
	for _, pair := range v2 {
		lru.Put(pair.key, pair.value)
		questr = "2 " + questr
	}
	questr = questr[:2*lru.capacity]
	if lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}

	v3 := []Pair{
		{"c0", &testHash{v: 3}},
		{"c1", &testHash{v: 13}},
		{"c2", &testHash{v: 23}},
		{"c3", &testHash{v: 33}},
		{"c8", &testHash{v: 83}},
	}

	// 3 を追加: 最初に追加した 1 は全部消える
	for _, pair := range v3 {
		lru.Put(pair.key, pair.value)
		questr = "3 " + questr
	}
	questr = questr[:2*lru.capacity]

	if lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}

	v1 := []Pair{
		{"a8", &testHash{v: 81}},
		{"a7", &testHash{v: 71}},
	}

	v2 = v2[len(v2)-10+len(v1)+len(v3):]
	if len(v1)+len(v2)+len(v3) != lru.Capacity() {
		t.Errorf("v1=%d, v2=%d, v3=%d, lru.Capacity()=%d", len(v1), len(v2), len(v3), lru.Capacity())
		return
	}

	// 1 を追加
	for _, pair := range v1 {
		lru.Put(pair.key, pair.value)
		questr = "1 " + questr
	}
	questr = questr[:2*lru.capacity]
	if lru.queueString() != questr {
		t.Errorf("Expected queue to be '%s', got %s", questr, lru.queueString())
		return
	}

	if lru.Len() != lru.Capacity() {
		t.Errorf("Expected size %d, got Len=%d", lru.Capacity(), lru.Len())
		return
	}

	for _, vv := range [][]Pair{v1, v2, v3} {
		for _, pair := range vv {
			if v, ok := lru.Get(pair.key); !ok || v != pair.value {
				t.Errorf("Expected key %v; ok=%v value=%v, v=%v", pair.key, ok, pair.value, v)
				return
			}
		}
	}
}
