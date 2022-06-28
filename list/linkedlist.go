package list

import (
	"fmt"
	"github.com/joexzh/dsa"
	"strings"
)

type LinkedNode struct {
	pred *LinkedNode
	succ *LinkedNode
	Data interface{}
}

func (nd *LinkedNode) Pred() *LinkedNode {
	return nd.pred
}

func (nd *LinkedNode) Succ() *LinkedNode {
	return nd.succ
}

func (nd *LinkedNode) Valid() bool {
	return nd != nil && nd.pred != nil && nd.succ != nil
}

type LinkedList struct {
	header  *LinkedNode
	trailer *LinkedNode
	size    int
}

func NewLinkedList() LinkedList {
	header := &LinkedNode{}
	trailer := &LinkedNode{}
	header.succ = trailer
	trailer.pred = header
	return LinkedList{header: header, trailer: trailer}
}

func (l *LinkedList) Size() int {
	return l.size
}

// Find in an unordered list, before node p, traverse from right to left at most n elements,
// looking for a node.Data equals to Data
func (l *LinkedList) Find(e interface{}, n int, p *LinkedNode) *LinkedNode {
	for n > 0 {
		n--
		p := p.pred
		if p == nil || p == l.header {
			break
		}
		if e == p.Data {
			return p
		}

	}
	return nil
}

// Search only make sense for a sorted list, o(n)
func (l *LinkedList) Search(item dsa.Item) (*LinkedNode, bool) {
	curr := l.header.succ
	for curr != l.trailer {
		currItem, ok := curr.Data.(dsa.Item)
		if !ok {
			panic("LinkedList.Search method only support type implements dsa.Comparable interface!")
		}
		if !currItem.Less(item) && !item.Less(currItem) {
			return curr, true
		}
		if !currItem.Less(item) {
			return curr.pred, false
		}
		curr = curr.succ
	}
	return curr, false
}

func (l *LinkedList) InsertEnd(e interface{}) {
	l.InsertBefore(l.trailer, e)
}

func (l *LinkedList) InsertStart(e interface{}) {
	l.InsertAfter(l.header, e)
}

// InsertAfter a node, if the node is trailer or nil, insert will fail
func (l *LinkedList) InsertAfter(x *LinkedNode, e interface{}) bool {
	if x == l.trailer || x == nil {
		return false
	}
	succ := x.succ
	newnd := &LinkedNode{pred: x, succ: succ, Data: e}
	x.succ = newnd
	succ.pred = newnd
	l.size++
	return true
}

// InsertBefore a node, if the node is header or nil, insert will fail
func (l *LinkedList) InsertBefore(x *LinkedNode, e interface{}) bool {
	if x == l.header || x == nil {
		return false
	}
	pred := x.pred
	newnd := &LinkedNode{pred: pred, succ: x, Data: e}
	pred.succ = newnd
	x.pred = newnd
	l.size++
	return true
}

// Equal each element must implement dsa.Comparable interface
func (l *LinkedList) Equal(other LinkedList) bool {
	if l.Size() != other.Size() {
		return false
	}

	for me, oe := l.First(), other.First(); me.Valid(); me, oe = me.Succ(), oe.Succ() {
		if me.Data != oe.Data {
			return false
		}
	}
	return true
}

func (l *LinkedList) String() string {
	sb := strings.Builder{}
	sb.WriteString("header->")
	for nd := l.First(); nd.Valid(); nd = nd.Succ() {
		sb.WriteString(fmt.Sprint(nd.Data))
		sb.WriteString("->")
	}
	sb.WriteString("trailer\n")
	return sb.String()
}

func (l *LinkedList) Sort() {
	// todo linked list
	panic("not implement")
}

func (l *LinkedList) Remove(q *LinkedNode) {
	if !q.Valid() {
		return
	}
	q.pred.succ = q.succ
	q.succ.pred = q.pred
	q.pred, q.succ = nil, nil
	l.size--
}

func (l *LinkedList) First() *LinkedNode {
	return l.header.succ
}

func (l *LinkedList) Last() *LinkedNode {
	return l.trailer.pred
}

func (l *LinkedList) Traverse(f func(node *LinkedNode)) {
	for nd := l.First(); nd.Valid(); nd = nd.succ {
		f(nd)
	}
}

func (l LinkedList) mergeSort() {
	// todo linked list
	panic("not implement")
}
