package fixedheap

import (
	"container/heap"
	"fmt"
	"reflect"
	"sort"
)

// Numeric is an interface that your Item.Priority has to satisfy. As long as
// it is a numeric value, you should specify a simple LessThan function using
// type assertions of the Item.Priority type you are using.
//
// I had to do this because I had Item.Priority as ints for some use cases, and
// as floats for some instances.
type Numeric interface {
	LessThan(b Numeric) bool
}

// Item is the type that our Heap or Fix holds. It holds arbitrary data and
// a Priority field that should be implement Numeric
type Item struct {
	Value    interface{} // arbitrary data
	Priority Numeric
}

//////////////////////////////////////////////////////////////////////

// Heap implements heap.Interface methods
type Heap []*Item

func (h Heap) Len() int { return len(h) }

func (h Heap) Less(i, j int) bool {
	return (h[i].Priority).LessThan(h[j].Priority)
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Push ...
func (h *Heap) Push(x interface{}) {
	item := x.(Item)
	*h = append(*h, &item)
}

// Pop ...
func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// New returns a initialized instance of a Heap pointer
func New() *Heap {
	h := &Heap{}
	heap.Init(h)
	return h
}

//////////////////////////////////////////////////////////////////////

// Fix is a specialized Heap that has an int field Cap that is the
// maximum number of items that the Heap can hold
type Fix struct {
	Heap
	Cap int
}

// ManagedPush pushes the top N highest priority items to the Fix,
// where N = fh.Cap
func ManagedPush(fh *Fix, it *Item) {
	i := *it
	pq := &fh.Heap
	if pq.Len() < fh.Cap {
		heap.Push(pq, i)
		return
	} else if ((*pq)[0].Priority).LessThan(i.Priority) {
		(*pq)[0] = &i
		heap.Fix(pq, 0)
	}
}

// NewFix returns an initialized Fix pointer with capacity = cap
func NewFix(cap int) *Fix {
	pq := &Heap{}
	heap.Init(pq)
	fh := &Fix{*pq, cap}
	return fh
}

//////////////////////////////////////////////////////////////////////

func (f Fix) String() string {
	h := f.Heap

	var outputLength int
	if f.Cap < f.Heap.Len() {
		outputLength = f.Cap
	} else {
		outputLength = f.Heap.Len()
	}
	strings := make([]string, outputLength)
	// TODO not sure if sorting in Stringer is OK
	sort.Sort(h)
	for i := range strings {
		strings[i] = h[i].String()
	}
	return fmt.Sprintf("%v --> %d", strings, f.Cap)
}

func (h Heap) String() string {
	strings := make([]string, h.Len())
	// TODO not sure if sorting in Stringer is OK
	sort.Sort(h)
	for i, v := range h {
		strings[i] = v.String()
	}
	return fmt.Sprintf("%v", strings)
}

func (i Item) String() string {
	v := reflect.ValueOf(i.Priority)
	if v.Kind() == reflect.Float64 {
		return fmt.Sprintf("(%v, %.3f)", i.Value, v.Interface())
	}
	return fmt.Sprintf("(%v, %v)", i.Value, i.Priority)
}
