package ganrac

/***
 * LRUCache: Least Recently Used Cache
 */

import (
	ll "container/list"
	"fmt"
)

// LRUCache.(queue|map) に格納する要素の型
type LRUCacheValue[T any] struct {
	key    Hashable
	value  T
	queptr *ll.Element
	mapptr *ll.Element
}

type LRUCache[T any] struct {
	capacity int
	m        map[Hash]*ll.List
	queue    *ll.List // front: recently used ==> back: least recently used

	cntHit  int
	cntMiss int
	cntDel  int
	cntPut  int
}

func NewLRUCache[T any](capacity int) *LRUCache[T] {
	lru := &LRUCache[T]{capacity: capacity}
	lru.m = make(map[Hash]*ll.List, capacity*2)
	lru.queue = ll.New()
	return lru
}

func (lru *LRUCache[T]) Len() int {
	return lru.queue.Len()
}

func (lru *LRUCache[T]) lenMapList() int {
	n := 0
	for _, v := range lru.m {
		n += v.Len()
	}
	return n
}

func (lru *LRUCache[T]) Capacity() int {
	return lru.capacity
}

func (lru *LRUCache[T]) Put(key Hashable, value T) {
	if l := lru.queue.Len(); l >= lru.capacity {
		lru.removeN(l - lru.capacity + 1)
	}
	lru.add(key, value)
	lru.cntPut++
}

func (lru *LRUCache[T]) Get(key Hashable) (T, bool) {
	h := key.Hash()
	if v, ok := lru.m[h]; ok {
		// v を走査
		for e := v.Front(); e != nil; e = e.Next() {
			u := e.Value.(*LRUCacheValue[T])
			if key.Equals(u.key) {
				lru.queue.MoveToFront(u.queptr)
				lru.cntHit++
				return u.value, true
			}
		}
	}
	lru.cntMiss++
	var v T
	return v, false
}

func (lru *LRUCache[T]) removeN(n int) {
	for i := 0; i < n; i++ {
		lru.remove()
	}
}

func (lru *LRUCache[T]) remove() {
	elem := lru.queue.Back() // lru element
	lru.queue.Remove(elem)

	node := elem.Value.(*LRUCacheValue[T])
	h := node.key.Hash()
	if mp, ok := lru.m[h]; ok {
		if mp.Len() <= 1 {
			delete(lru.m, h)
		} else {
			mp.Remove(node.mapptr)
		}
	} else {
		panic(fmt.Sprintf("remove map failed: h=%v; node=%v", h, node))
	}

	lru.cntDel++
}

func (lru *LRUCache[T]) add(key Hashable, value T) {
	h := key.Hash()
	vv := &LRUCacheValue[T]{key: key, value: value}
	vv.queptr = lru.queue.PushFront(vv)
	if _, ok := lru.m[h]; !ok {
		lru.m[h] = ll.New()
	}
	vv.mapptr = lru.m[h].PushBack(vv)
}

func (v *LRUCacheValue[T]) String() string {
	return fmt.Sprintf("LRUCacheValue{k=%v, v=%v}", v.key.Hash(), v.value)
}

func (lru *LRUCache[T]) String() string {
	return fmt.Sprintf("LRUCache{#=%v/%v, #map=%v, hit=%d/%d=%.2f, del=%d, put=%d}",
		lru.Len(), lru.capacity,
		len(lru.m),
		lru.cntHit, lru.cntHit+lru.cntMiss,
		float64(lru.cntHit)/float64(lru.cntHit+lru.cntMiss),
		lru.cntDel, lru.cntPut)
}

////////////////////////////////////////
// for DEBUG
////////////////////////////////////////

// テスト用: キューの状態を文字列で返す
func (lru *LRUCache[T]) queueString() string {
	str := ""
	for e := lru.queue.Front(); e != nil; e = e.Next() {
		v := e.Value.(*LRUCacheValue[T])
		str += fmt.Sprintf("%x ", v.key.Hash())
	}
	return str
}
