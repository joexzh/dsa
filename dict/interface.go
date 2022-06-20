package dict

import "github.com/joexzh/dsa"

type Dictionary interface {
	Put(k dsa.Comparable, v interface{}) bool
	Remove(k dsa.Comparable) bool
	Get(k dsa.Comparable) interface{}
}
