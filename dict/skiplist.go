package dict

import (
	"github.com/joexzh/dsa"
	"github.com/joexzh/dsa/list"
	"math/rand"
)

type SkipList struct {
	list.LinkedList // holding layers of quadList
}

func NewSkipList() SkipList {
	sl := SkipList{LinkedList: list.NewLinkedList()}
	return sl
}

func (sl *SkipList) Size() int {
	if sl.Empty() {
		return 0
	}
	return sl.Last().Data.(*quadList).Size()
}

// GetRange by [startK, endK) of key range, returns ordered entries
func (sl *SkipList) GetRange(startK dsa.Comparable, endK dsa.Comparable) list.LinkedList {
	ents := list.NewLinkedList()
	if sl.Empty() {
		return ents
	}
	qlist := sl.First()
	nd := qlist.Data.(*quadList).First()

	ok := sl.search(&qlist, &nd, startK)
	for nd.below != nil {
		nd = nd.below
	}
	if ok {
		for nd.pred.valid() && nd.pred.entry.K.Compare(startK) >= 0 { // try to lookup before
			nd = nd.pred
		}
	} else {
		nd = nd.succ
	}
	for ; nd.valid() && nd.entry.K.Compare(endK) < 0; nd = nd.succ {
		ents.InsertEnd(nd.entry)
	}
	return ents
}

// Get value of k. If k is duplicated, pick the first one with the highest tower
func (sl *SkipList) Get(k dsa.Comparable) interface{} {
	qlist := sl.First()
	nd := qlist.Data.(*quadList).First()
	if ok := sl.search(&qlist, &nd, k); ok {
		return nd.entry.V
	}
	return nil
}

// Put will always successes, same key already sitting there, new key will insert after it.
func (sl *SkipList) Put(k dsa.Comparable, v interface{}) bool {
	e := dsa.Entry{K: k, V: v}
	if sl.Empty() {
		sl.addAsFirstLayer()
	}

	qlist := sl.First()
	qnd := qlist.Data.(*quadList).First()
	if ok := sl.search(&qlist, &qnd, k); ok {
		for qnd.below != nil {
			qnd = qnd.below
		}
	}

	for qnd.succ.valid() && qnd.succ.entry.K.Compare(k) <= 0 { // lookup to right side
		qnd = qnd.succ
	}
	qlist = sl.Last()
	b := qlist.Data.(*quadList).InsertAfterAbove(e, qnd, nil) // insert first node

	for rand.Intn(2)&1 == 1 { // 50% chance to add addition node on top of the tower
		for qnd.valid() && qnd.above == nil {
			qnd = qnd.pred
		}
		if !qnd.valid() {
			if qlist == sl.First() {
				sl.addAsFirstLayer() // if run out of layer, add one
			}
			qnd = qlist.Pred().Data.(*quadList).header // move qnd to above layer's header
		} else {
			qnd = qnd.above
		}
		qlist = qlist.Pred()
		b = qlist.Data.(*quadList).InsertAfterAbove(e, qnd, b)
	}
	return true
}

// Remove one key, if there are duplicate keys, remove the first one with the highest tower
func (sl *SkipList) Remove(k dsa.Comparable) bool {
	if sl.Empty() {
		return false
	}

	qlist := sl.First()
	qnd := sl.First().Data.(*quadList).First()
	if !sl.search(&qlist, &qnd, k) {
		return false
	}
	for qlist.Succ() != nil {
		lower := qnd.below
		qlist.Data.(*quadList).Remove(qnd)
		qnd = lower
		qlist = qlist.Succ()
	}
	for !sl.Empty() && sl.First().Data.(*quadList).Empty() {
		sl.LinkedList.Remove(sl.First())
	}
	return true
}

func (sl *SkipList) Empty() bool {
	return !sl.First().Valid()
}

func (sl *SkipList) Traverse(f func(e dsa.Entry)) {
	if sl.Empty() {
		return
	}

	sl.Last().Data.(*quadList).Traverse(f)
}

// Walk each tower from bottom to top, from left to right, tower num start from 0
func (sl *SkipList) Walk(f func(e dsa.Entry, tower int)) {
	if sl.Empty() {
		return
	}
	qlist := sl.Last()
	towerSeq := 0
	for bottom := qlist.Data.(*quadList).First(); bottom.valid(); bottom = bottom.succ {
		f(bottom.entry, towerSeq)
		for qnd := bottom.above; qnd != nil; qnd = qnd.above {
			f(qnd.entry, towerSeq)
		}
		towerSeq++
	}
}

func (sl *SkipList) addAsFirstLayer() {
	sl.InsertBefore(sl.First(), NewQuadList())
}

// for structure print and test
func (sl *SkipList) level() int {
	level := 0
	for qlist := sl.First(); qlist.Valid(); qlist = qlist.Succ() {
		level++
	}
	return level
}

// search for the first match node in the top of the tower, qnd will be that node.
// If not found, qnd will be in the lowest level, with the biggest key whose smaller than k.
// If qnd's key is duplicated, get the last node
// This function accepts pointer of pointer, so it can replace it.
func (sl *SkipList) search(qlistP **list.LinkedNode, qndP **quadNode, k dsa.Comparable) bool {
	for true {
		for (*qndP).succ != nil && (*qndP).entry.K.Compare(k) <= 0 {
			*qndP = (*qndP).succ
		}
		*qndP = (*qndP).pred
		if (*qndP).pred != nil && (*qndP).entry.K == k {
			return true
		}
		*qlistP = (*qlistP).Succ()
		if (*qlistP).Succ() == nil {
			return false
		}
		if (*qndP).pred != nil {
			*qndP = (*qndP).below
		} else {
			*qndP = (*qlistP).Data.(*quadList).First()
		}
	}
	return false
}

type quadNode struct {
	pred  *quadNode
	succ  *quadNode
	above *quadNode
	below *quadNode
	entry dsa.Entry
}

func (nd *quadNode) valid() bool {
	return nd != nil && nd.pred != nil && nd.succ != nil
}

type quadList struct {
	header  *quadNode
	trailer *quadNode
	size    int
}

func NewQuadList() *quadList {
	ql := &quadList{}
	ql.header = new(quadNode)
	ql.trailer = new(quadNode)
	ql.header.succ = ql.trailer
	ql.trailer.pred = ql.header
	return ql
}

func (ql *quadList) First() *quadNode {
	return ql.header.succ
}

func (ql *quadList) Last() *quadNode {
	return ql.trailer.pred
}

func (ql *quadList) Empty() bool {
	return ql.size <= 0
}

func (ql *quadList) Size() int {
	return ql.size
}

// Remove *quadNode, returns entry of it
func (ql *quadList) Remove(x *quadNode) (v interface{}) {
	x.pred.succ = x.succ
	x.succ.pred = x.pred
	ql.size--
	return x.entry
}

func (ql *quadList) InsertAfterAbove(e dsa.Entry, p *quadNode, b *quadNode) *quadNode {
	succ := p.succ
	newQnd := &quadNode{pred: p, succ: p.succ, below: b, entry: e}
	if b != nil {
		b.above = newQnd
	}
	p.succ = newQnd
	succ.pred = newQnd
	ql.size++
	return newQnd
}

func (ql *quadList) Traverse(f func(e dsa.Entry)) {
	for qnd := ql.First(); qnd != nil && qnd.succ != nil; qnd = qnd.succ {
		f(qnd.entry)
	}
}

func (ql *quadList) clear() {
	ql.header.succ = ql.trailer
	ql.trailer.pred = ql.header
	ql.size = 0
}
