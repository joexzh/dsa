package dict

import (
	"github.com/joexzh/dsa"
	"github.com/joexzh/dsa/list"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
)

type myInt int

func (m myInt) Compare(other dsa.Comparable) int {
	o, ok := other.(myInt)
	if !ok {
		panic("myInt.Compare method's parameter type must be myInt as well")
	}
	if m > o {
		return 1
	}
	if m < o {
		return -1
	}
	return 0
}

func SkipListSizeTest(t *testing.T, sl SkipList, expected int) {
	size := sl.Size()
	if size != expected {
		t.Fatalf("expected %v, got %v", expected, size)
	}
}

func TestSkipList_Empty(t *testing.T) {
	sl := NewSkipList()
	sl.Put(myInt(2), 4.3)
	sl.Put(myInt(3), 4.3)
	sl.Put(myInt(2), 4.3)
	sl.Put(myInt(2), 4.3)
	sl.Put(myInt(4), 4.3)
	sl.Put(myInt(6), 4.3)
	sl.Put(myInt(5), 4.3)
	SkipListSizeTest(t, sl, 7)

	if sl.Empty() {
		t.Fatalf("excepte not empty, got empty")
	}
	if !sl.Remove(myInt(3)) {
		t.Fatalf("key %v expected remove success, got fail", myInt(3))
	}
	if sl.Remove(myInt(1)) {
		t.Fatalf("key %v expected remove fail, got success", myInt(1))
	}
	if !sl.Remove(myInt(4)) {
		t.Fatalf("key %v expected remove success, got fail", myInt(4))
	}
	if !sl.Remove(myInt(5)) {
		t.Fatalf("key %v expected remove success, got fail", myInt(5))
	}
	if !sl.Remove(myInt(6)) {
		t.Fatalf("key %v expected remove success, got fail", myInt(6))
	}
	if !sl.Remove(myInt(2)) {
		t.Fatalf("key %v expected remove success, got fail", myInt(2))
	}
	if !sl.Remove(myInt(2)) {
		t.Fatalf("key %v expected remove success, got fail", myInt(2))
	}
	if !sl.Remove(myInt(2)) {
		t.Fatalf("key %v expected remove success, got fail", myInt(2))
	}

	if !sl.Empty() {
		t.Fatalf("expect empty, got not empty")
	}
	SkipListSizeTest(t, sl, 0)
}

func TestSkipList_GetRange(t *testing.T) {
	sl := NewSkipList()
	sl.Put(myInt(1), 1.2)
	sl.Put(myInt(11), 11.2)
	sl.Put(myInt(4), 4.2)
	sl.Put(myInt(1), 1.22)
	sl.Put(myInt(1), 1.23)
	sl.Put(myInt(7), 7.2)
	sl.Put(myInt(9), 9.2)
	sl.Put(myInt(4), 4.22)
	sl.Put(myInt(4), 4.222)
	sl.Put(myInt(4), 4.24)

	rl := sl.GetRange(myInt(1), myInt(5))

	expectList := list.NewLinkedList()
	expectList.InsertEnd(dsa.Entry{K: myInt(1), V: 1.2})
	expectList.InsertEnd(dsa.Entry{K: myInt(1), V: 1.22})
	expectList.InsertEnd(dsa.Entry{K: myInt(1), V: 1.23})
	expectList.InsertEnd(dsa.Entry{K: myInt(4), V: 4.2})
	expectList.InsertEnd(dsa.Entry{K: myInt(4), V: 4.22})
	expectList.InsertEnd(dsa.Entry{K: myInt(4), V: 4.222})
	expectList.InsertEnd(dsa.Entry{K: myInt(4), V: 4.24})

	if !rl.Equal(expectList) {
		t.Fatalf("expected %s, got %s", expectList.String(), rl.String())
	}
	SkipListSizeTest(t, sl, 10)
}

func TestSkipList_Get(t *testing.T) {
	sl := NewSkipList()
	sl.Put(myInt(1), 1.2)
	sl.Put(myInt(11), 11.2)
	sl.Put(myInt(4), 4.2)
	sl.Put(myInt(1), 1.22)

	got := sl.Get(myInt(1))
	expected1 := 1.2
	expected2 := 1.22
	if got == nil || (got.(float64) != expected1 && got.(float64) != expected2) {
		t.Fatalf("expected %v or %v, got %v", expected1, expected2, got)
	}
	wrongGot := sl.Get(myInt(2))
	if wrongGot != nil {
		t.Fatalf("expected nil, got %v", wrongGot)
	}
	SkipListSizeTest(t, sl, 4)
}

func TestSkipList_Put(t *testing.T) {
	sl := NewSkipList()
	sl.Put(myInt(2), 2.2)
	sl.Put(myInt(5), 5.2)
	sl.Put(myInt(3), 3.2)

	expected1 := 2.2
	got1 := sl.Get(myInt(2))
	if got1 == nil || got1.(float64) != expected1 {
		t.Fatalf("expected %v, got %v", expected1, got1)
	}
	expected2 := 3.2
	got2 := sl.Get(myInt(3))
	if got2 == nil || got2.(float64) != expected2 {
		t.Fatalf("expected %v, got %v", expected2, got2)
	}
	expected3 := 5.2
	got3 := sl.Get(myInt(5))
	if got3 == nil || got3.(float64) != expected3 {
		t.Fatalf("expected %v, got %v", expected3, got3)
	}

	SkipListSizeTest(t, sl, 3)
}

type SortableEntry []dsa.Entry

func (s SortableEntry) Len() int {
	return len(s)
}

func (s SortableEntry) Less(i, j int) bool {
	if s[i].K.Compare(s[j].K) < 0 {
		return true
	}
	return false
}

func (s SortableEntry) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func TestSkipList_Traverse(t *testing.T) {
	sl := NewSkipList()
	size := 101

	entries := make([]dsa.Entry, size, size)
	for i := range entries {
		key := myInt(rand.Intn(math.MaxInt))
		val := myInt(rand.Intn(math.MaxInt))

		e := dsa.Entry{K: key, V: val}
		entries[i] = e
		sl.Put(e.K, e.V)
	}
	sort.Sort(SortableEntry(entries))

	expected := entries
	got := make([]dsa.Entry, 0, size)
	sl.Traverse(func(e dsa.Entry) {
		got = append(got, e)
	})
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

// if you want to see the result, add -v flag
func TestSkipList_Walk(t *testing.T) {
	sl := NewSkipList()
	rand.Seed(time.Now().UnixNano())

	for range [83]struct{}{} {
		key := myInt(rand.Intn(math.MaxInt))
		val := myInt(rand.Intn(math.MaxInt))
		sl.Put(key, val)
	}
	towers := make([][]dsa.Entry, 0, 1)
	sl.Walk(func(e dsa.Entry, tower int) {
		if len(towers) < tower+1 {
			towers = append(towers, []dsa.Entry{})
		}
		towers[tower] = append(towers[tower], e)
	})

	for i := range towers {
		t.Logf("tower %d: %v\n", i, towers[i])
	}
}
