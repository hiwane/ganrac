package cache

import (
	"container/list"
	"github.com/hiwane/ganrac"
	// "fmt"
)

type QueHash string

type QueValue struct {
	p   *ganrac.Poly
	q   *ganrac.Poly
	r   ganrac.RObj
	lv  ganrac.Level
	l   *ganrac.List
	cc  int32
	n   int
	key QueHash
}

type QueMap struct {
	capacity int
	m        map[QueHash]*list.Element
	list     *list.List

	hit  int
	miss int
	rm   int
}

func NewQueMap(capacity int) *QueMap {
	qm := &QueMap{capacity: capacity}
	qm.m = make(map[QueHash]*list.Element, capacity)
	qm.list = list.New()
	return qm
}

func (qm *QueMap) Size() int {
	return len(qm.m)
}

func (qm *QueMap) Capacity() int {
	return qm.capacity
}

func (qv *QueValue) Hash() QueHash {
	return QueHash(qv.p.String())
}

func (qm *QueMap) Put(value *QueValue) {
	if len(qm.m) >= qm.capacity {
		qm.remove(qm.list.Back())
	}
	qm.add(value)
}

func (qm *QueMap) Get(key QueHash) *QueValue {
	if v, ok := qm.m[key]; ok {
		qm.list.MoveToFront(v)
		qm.hit++
		// fmt.Printf("Cache.QueMap.Get(%d:%d:%s) hit\n", len(qm.m), qm.list.Len(), key)
		return v.Value.(*QueValue)
	}
	qm.miss++
	return nil
}

func (qm *QueMap) remove(elem *list.Element) {
	node := elem.Value.(*QueValue)
	qm.list.Remove(elem)
	delete(qm.m, node.key)
	qm.rm++
}

func (qm *QueMap) add(value *QueValue) {
	elem := qm.list.PushFront(value)
	qm.m[value.key] = elem
}
