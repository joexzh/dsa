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

func (m myInt) Less(than dsa.Item) bool {
	th, ok := than.(myInt)
	if !ok {
		panic("myInt.Compare method's parameter type must be myInt as well")
	}
	if m < th {
		return true
	}
	return false
}

func SkipListSizeTest(t *testing.T, sl SkipList, expected int) {
	size := sl.Size()
	if size != expected {
		t.Fatalf("expected %v, got %v", expected, size)
	}
}

func TestSkipList_Empty(t *testing.T) {
	sl := NewSkipList()
	sl.Put(myInt(2), 2.1)
	sl.Put(myInt(3), 3.1)
	sl.Put(myInt(2), 2.2)
	sl.Put(myInt(2), 2.3)
	sl.Put(myInt(4), 4.1)
	sl.Put(myInt(6), 6.1)
	sl.Put(myInt(5), 5.1)
	SkipListSizeTest(t, sl, 7)

	if sl.Empty() {
		t.Fatalf("excepte not empty, got empty")
	}
	n := sl.Remove(myInt(3))
	if n != 1 {
		t.Fatalf("key %v expected remove 1, got %d", myInt(3), n)
	}
	n = sl.Remove(myInt(1))
	if n != 0 {
		t.Fatalf("key %v expected remove 0, got %d", myInt(1), n)
	}
	n = sl.Remove(myInt(4))
	if n != 1 {
		t.Fatalf("key %v expected remove 1, got %d", myInt(4), n)
	}
	n = sl.Remove(myInt(5))
	if n != 1 {
		t.Fatalf("key %v expected remove 5, got %d", myInt(5), n)
	}
	n = sl.Remove(myInt(6))
	if n != 1 {
		t.Fatalf("key %v expected remove success, got %d", myInt(6), n)
	}
	n = sl.Remove(myInt(2))
	if n != 3 {
		t.Fatalf("key %v expected remove 3, got %d", myInt(2), n)
	}
	if !sl.Empty() {
		t.Fatalf("expect empty, got not empty")
	}
	SkipListSizeTest(t, sl, 0)
}

func TestSkipList_GetRange(t *testing.T) {
	sl := NewSkipList()
	sl.Put(myInt(1), 2^0)
	sl.Put(myInt(11), 2^1)
	sl.Put(myInt(4), 2^2)
	sl.Put(myInt(1), 2^3)
	sl.Put(myInt(1), 2^4)
	sl.Put(myInt(7), 2^5)
	sl.Put(myInt(9), 2^6)
	sl.Put(myInt(4), 2^7)
	sl.Put(myInt(4), 2^8)
	sl.Put(myInt(4), 2^9)

	rl := sl.GetRange(myInt(1), myInt(5))
	gotEntries := make([]dsa.Entry, 0)
	rl.Traverse(func(nd *list.LinkedNode) {
		e := nd.Data.(dsa.Entry)
		gotEntries = append(gotEntries, e)
	})

	expectedEntries := []dsa.Entry{{myInt(1), 2 ^ 0}, {myInt(1), 2 ^ 3}, {myInt(1), 2 ^ 4},
		{myInt(4), 2 ^ 2}, {myInt(4), 2 ^ 7}, {myInt(4), 2 ^ 8}, {myInt(4), 2 ^ 9}}

	expectedLen := len(expectedEntries)
	gotLen := len(gotEntries)
	if len(gotEntries) != len(expectedEntries) {
		t.Logf("expectedEnt %v, gotEnt %v", expectedEntries, gotEntries)
		t.Fatalf("expected gotEntries len %v, got %v", expectedLen, gotLen)
	}

	for i := range gotEntries {
		if gotEntries[i].K != expectedEntries[i].K {
			t.Fatalf("expecte %d, got %d, unmatch index %d", expectedEntries, gotEntries, i)
		}
	}

	var gotSum1 int
	for _, e := range gotEntries {
		if e.K == myInt(1) {
			gotSum1 |= e.V.(int)
		}
	}
	var expectedSum1 int
	for _, e := range expectedEntries {
		if e.K == myInt(1) {
			expectedSum1 |= e.V.(int)
		}
	}
	if expectedSum1 != gotSum1 {
		t.Fatalf("expectedSum1 %d, gotSum1 %d", expectedSum1, gotSum1)
	}

	var gotSum4 int
	for _, e := range gotEntries {
		if e.K == myInt(4) {
			gotSum4 |= e.V.(int)
		}
	}
	var expectedSum4 int
	for _, e := range expectedEntries {
		if e.K == myInt(4) {
			expectedSum4 |= e.V.(int)
		}
	}
	if expectedSum4 != gotSum4 {
		t.Fatalf("expectedSum4 %d, gotSum4 %d", expectedSum4, gotSum4)
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

func TestSkipList_Replace(t *testing.T) {
	sl := NewSkipList()
	sl.Put(myInt(1), 2^10)
	sl.Replace(myInt(1), 2^11)

	expected := 2 ^ 11
	got := sl.Get(myInt(1)).(int)
	if expected != got {
		t.Fatalf("expected %d, got %d", expected, got)
	}

	SkipListSizeTest(t, sl, 1)

	sl.Replace(myInt(100), 2^12)
	sl.Replace(myInt(101), 2^13)
	sl.Replace(myInt(102), 2^14)
	sl.Replace(myInt(103), 2^15)

	sl.Replace(myInt(101), 0)

	expected = 0
	got = sl.Get(myInt(101)).(int)
	if expected != got {
		t.Fatalf("expected %d, got %d", expected, got)
	}

	expected = 2 ^ 14
	got = sl.Get(myInt(102)).(int)
	if expected != got {
		t.Fatalf("expected %d, got %d", expected, got)
	}

	SkipListSizeTest(t, sl, 5)

}

func TestSkipList_DuplicatedKey_Order(t *testing.T) {
	sl := NewSkipList()
	sl.Put(myInt(1), myInt(5))
	sl.Put(myInt(1), myInt(3))
	sl.Put(myInt(1), myInt(1))

	sl.Walk(func(e dsa.Entry, tower int) {
		t.Logf("walk: e %v, tower %d", e, tower)
	})
	expected := myInt(1)
	got := sl.Get(myInt(1)).(myInt)
	if expected != got {
		t.Fatalf("expected %d, got %d", expected, got)
	}

	SkipListSizeTest(t, sl, 3)

	sl.Put(myInt(11), myInt(2))
	sl.Put(myInt(11), myInt(5))
	sl.Put(myInt(11), myInt(2))
	sl.Put(myInt(11), myInt(4))
	sl.Put(myInt(11), myInt(4))
	sl.Put(myInt(11), myInt(6))
	sl.Put(myInt(11), myInt(6))
	expected = myInt(2)
	got = sl.Get(myInt(11)).(myInt)
	if expected != got {
		t.Fatalf("expected %d, got %d", expected, got)
	}

	expectedll := list.NewLinkedList()
	expectedll.InsertEnd(dsa.Entry{myInt(11), myInt(2)})
	expectedll.InsertEnd(dsa.Entry{myInt(11), myInt(2)})
	expectedll.InsertEnd(dsa.Entry{myInt(11), myInt(4)})
	expectedll.InsertEnd(dsa.Entry{myInt(11), myInt(4)})
	expectedll.InsertEnd(dsa.Entry{myInt(11), myInt(5)})
	expectedll.InsertEnd(dsa.Entry{myInt(11), myInt(6)})
	expectedll.InsertEnd(dsa.Entry{myInt(11), myInt(6)})
	gotll := sl.GetRange(myInt(11), myInt(11))

	sl.Walk(func(e dsa.Entry, tower int) {
		t.Logf("walk: e %v, tower %d", e, tower)
	})

	if !expectedll.Equal(gotll) {
		t.Fatalf("expected linkedlist %s, got linkedlist %s", expectedll.String(), gotll.String())
	}

	SkipListSizeTest(t, sl, 10)
}

type SortableEntry []dsa.Entry

func (s SortableEntry) Len() int {
	return len(s)
}

func (s SortableEntry) Less(i, j int) bool {
	return s[i].K.Less(s[j].K)
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
