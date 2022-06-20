package dsa

type Comparable interface {
	Compare(e Comparable) int // -1 less, 0 equal, 1 greater
}
