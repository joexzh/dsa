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

// GetRange by [startK, endK] of key range, returns ordered entries
func (sl *SkipList) GetRange(startK dsa.Item, endK dsa.Item) list.LinkedList {
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
		for predNd := nd.pred; predNd.valid() && !predNd.entry.K.Less(startK); predNd = predNd.pred { // try to lookup before
			ents.InsertStart(predNd.entry)
		}
	} else {
		nd = nd.succ
	}
	for ; nd.valid() && !endK.Less(nd.entry.K); nd = nd.succ {
		ents.InsertEnd(nd.entry)
	}
	return ents
}

// Get value of key.
//
// If key is duplicated, and Entry.V implements dsa.Item interface, get the first match in ascending order,
// otherwise simply take the highest tower which results like random.
func (sl *SkipList) Get(key dsa.Item) interface{} {
	qlist := sl.First()
	nd := qlist.Data.(*quadList).First()
	if ok := sl.search(&qlist, &nd, key); ok {
		_, ok := nd.entry.V.(dsa.Item)
		if ok {
			for nd.below.valid() {
				nd = nd.below
			}
			for ; nd.pred.valid() && !nd.pred.entry.K.Less(key); nd = nd.pred {
				_, ok := nd.pred.entry.V.(dsa.Item)
				if !ok { // pred entry.V not support dsa.Item interface
					return nd.entry.V
				}
			}
		}
		return nd.entry.V
	}
	return nil
}

// Put inserts a key-value pair, will always succeed.
// If key is duplicated, and v implements dsa.Item interface, put among them in ascending order,
// otherwise, insert after the founded node (acts like random).
func (sl *SkipList) Put(key dsa.Item, value interface{}) {
	e := dsa.Entry{K: key, V: value}
	if sl.Empty() {
		sl.addAsFirstLayer()
	}

	qlist := sl.First()
	nd := qlist.Data.(*quadList).First()
	if ok := sl.search(&qlist, &nd, key); ok {
		for nd.below != nil {
			nd = nd.below
		}

		itemV, vOk := value.(dsa.Item)
		itemNdV, ok := nd.entry.V.(dsa.Item)
		if vOk && ok { // compare value and nd.entry.V to determine forward or backward
			if itemNdV.Less(itemV) { // backward
				for ; nd.succ.valid() && !key.Less(nd.succ.entry.K); nd = nd.succ {
					itemNdV, ok := nd.succ.entry.V.(dsa.Item)
					if !ok || itemV.Less(itemNdV) {
						break
					}
				}

			} else { // forward
				for nd = nd.pred; nd.valid() && !nd.entry.K.Less(key); nd = nd.pred {
					itemNdV, ok := nd.entry.V.(dsa.Item)
					if !ok || !itemV.Less(itemNdV) {
						break
					}
				}
			}
		}
	}

	sl.insert(&e, nd)
}

// Replace key-value pair, if key not exist, insert it.
// Make sure keys are not duplicated, and don't mix with Put method.
func (sl *SkipList) Replace(key dsa.Item, value interface{}) {
	e := dsa.Entry{K: key, V: value}
	if sl.Empty() {
		sl.addAsFirstLayer()
	}

	qlist := sl.First()
	nd := qlist.Data.(*quadList).First()
	if ok := sl.search(&qlist, &nd, key); ok {
		for ; nd != nil; nd = nd.below { // replace entry from top to bottom
			nd.entry = e
		}
	} else { // not exist
		sl.insert(&e, nd)
	}
}

// Remove all nodes of the same key.
func (sl *SkipList) Remove(key dsa.Item) int {
	if sl.Empty() {
		return 0
	}

	qlist := sl.First()
	nd := sl.First().Data.(*quadList).First()
	if !sl.search(&qlist, &nd, key) {
		return 0
	}

	predBottom, succBottom := sl.remove(qlist, nd)
	n := 1
	for nd = predBottom; nd.valid() && !nd.entry.K.Less(key); nd = predBottom { // remove left
		qlist = sl.Last()
		for nd.above != nil {
			nd = nd.above
			qlist = qlist.Pred()
		}
		predBottom, _ = sl.remove(qlist, nd)
		n++
	}
	for nd = succBottom; nd.valid() && !key.Less(nd.entry.K); nd = succBottom { // remove right
		qlist = sl.Last()
		for nd.above != nil {
			nd = nd.above
			qlist = qlist.Pred()
		}
		_, succBottom = sl.remove(qlist, nd)
		n++
	}

	return n
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

// insert e after nd, build a tower. nd must be a bottom node.
// You should use search function getting nd at top of tower first, then travel to bottom.
func (sl *SkipList) insert(e *dsa.Entry, nd *quadNode) {
	qlist := sl.Last()
	b := qlist.Data.(*quadList).InsertAfterAbove(*e, nd, nil) // insert first node

	for rand.Intn(2)&1 == 1 { // 50% chance to add addition node on top of the tower
		for nd.valid() && nd.above == nil {
			nd = nd.pred
		}
		if !nd.valid() {
			if qlist == sl.First() {
				sl.addAsFirstLayer() // if run out of layer, add one
			}
			nd = qlist.Pred().Data.(*quadList).header // move nd to above layer's header
		} else {
			nd = nd.above
		}
		qlist = qlist.Pred()
		b = qlist.Data.(*quadList).InsertAfterAbove(*e, nd, b)
	}
}

// remove from top to bottom.
// You should use search method first to get qlist and nd
func (sl *SkipList) remove(qlist *list.LinkedNode, nd *quadNode) (predBottom *quadNode, succBottom *quadNode) {
	for qlist.Succ() != nil { // remove from top to bottom
		lower := nd.below
		if lower == nil {
			predBottom = nd.pred
			succBottom = nd.succ
		}
		qlist.Data.(*quadList).Remove(nd)
		nd = lower
		qlist = qlist.Succ()
	}
	for !sl.Empty() && sl.First().Data.(*quadList).Empty() {
		sl.LinkedList.Remove(sl.First())
	}
	return
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
func (sl *SkipList) search(qlistP **list.LinkedNode, qndP **quadNode, k dsa.Item) bool {
	for true {
		for (*qndP).succ != nil && !k.Less((*qndP).entry.K) {
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
func (ql *quadList) Remove(x *quadNode) dsa.Entry {
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
