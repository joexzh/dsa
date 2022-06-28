package dsa

import "strings"

//
// type Comparable interface {
// 	Compare(e Comparable) int // -1 less, 0 equal, 1 greater
// }

// Item is used to compare two element.
// If both direction Less return false, the two element is considered equal
type Item interface {
	Less(than Item) bool
}

type Int64 int64

func (i Int64) Less(than Item) bool {
	if i < than.(Int64) {
		return true
	}
	return false
}

type Uint64 uint64

func (u Uint64) Less(than Item) bool {
	if u < than.(Uint64) {
		return true
	}
	return false
}

type String string

func (s String) Less(than Item) bool {
	if strings.Compare(string(s), string(than.(String))) < 0 {
		return true
	}
	return false
}
