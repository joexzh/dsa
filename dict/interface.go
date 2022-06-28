package dict

import "github.com/joexzh/dsa"

type Dictionary interface {
	Put(k dsa.Item, v interface{}) bool
	Remove(k dsa.Item) bool
	Get(k dsa.Item) interface{}
}
