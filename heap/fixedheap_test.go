package fixedheap_test

import (
	"container/heap"
	"fmt"
	"testing"

	"github.com/nthapaliya/matasano/heap"
)

var (
	m1 = map[string]int{
		"0":  0,
		"1":  1,
		"2":  2,
		"3":  3,
		"4":  4,
		"5":  5,
		"6":  6,
		"7":  7,
		"8":  8,
		"9":  9,
		"10": 10,
		"11": 11,
	}

	m2 = map[string]float64{
		"0.01":  0.0123,
		"1.01":  1.0123,
		"2.01":  2.0123,
		"3.01":  3.0123,
		"4.01":  4.0123,
		"5.01":  5.0123,
		"6.01":  6.0123,
		"7.01":  7.0123,
		"8.01":  8.0123,
		"9.01":  9.0123,
		"10.01": 10.0123,
		"11.01": 11.0123,
	}
)

type (
	Int   int
	Float float64
)

func (i Int) LessThan(b fixedheap.Numeric) bool {
	return i < b.(Int)
}
func (f Float) LessThan(b fixedheap.Numeric) bool {
	return f < b.(Float)
}

var verbose = testing.Verbose()

func TestNew(t *testing.T) {
	f := fixedheap.New()
	if f == nil {
		t.Errorf("error creating Heap")
	} else if verbose {
		fmt.Println("Test New() passed")
	}
}

func TestStringerInts(t *testing.T) {
	f := fixedheap.New()
	m := m1
	for k, v := range m {
		i := fixedheap.Item{Value: k, Priority: Int(v)}
		heap.Push(f, i)
	}
	ans := "[(0, 0) (1, 1) (2, 2) (3, 3) (4, 4) (5, 5) (6, 6) (7, 7) (8, 8) (9, 9) (10, 10) (11, 11)]"
	if tmp := fmt.Sprint(f); ans != tmp {
		t.Errorf("stringer function failed, got \n%s\nneed\n%s", tmp, ans)
	} else if verbose {
		fmt.Println("Test stringer for Ints passed")
	}
}

func TestStringerFloats(t *testing.T) {
	f := fixedheap.New()
	m := m2
	for k, v := range m {
		i := fixedheap.Item{Value: k, Priority: Float(v)}
		heap.Push(f, i)
	}
	ans := "[(0.01, 0.012) (1.01, 1.012) (2.01, 2.012) (3.01, 3.012) (4.01, 4.012) (5.01, 5.012) (6.01, 6.012) (7.01, 7.012) (8.01, 8.012) (9.01, 9.012) (10.01, 10.012) (11.01, 11.012)]"
	if tmp := fmt.Sprint(f); ans != tmp {
		t.Errorf("stringer function failed, got \n%s\nneed\n%s", tmp, ans)
	} else if verbose {
		fmt.Println("Test stringer for Floats passed")
	}
}

func TestHeapPush(t *testing.T) {
	f := fixedheap.New()
	m := m2
	for k, v := range m {
		i := fixedheap.Item{Value: k, Priority: Float(v)}
		heap.Push(f, i)
	}
	if f.Len() != len(m) {
		t.Errorf("variable f empty, Push() not working")
	} else if verbose {
		fmt.Println("Test basic Push() to Heap instance passed")
	}
}

func TestHeapPop(t *testing.T) {
	f := fixedheap.New()
	m := m2
	for k, v := range m {
		i := fixedheap.Item{Value: k, Priority: Float(v)}
		heap.Push(f, i)
	}
	for f.Len() > 0 {
		f.Pop()
	}
	if verbose {
		fmt.Println("Test basic Pop() for Heap instance passed")
	}
}

func TestFHPush(t *testing.T) {
	f := fixedheap.NewFix(5)
	m := m2
	for k, v := range m {
		i := fixedheap.Item{Value: k, Priority: Float(v)}
		heap.Push(f, i)
	}
	if f.Cap != 5 || f.Heap.Len() != len(m) {
		t.Errorf("variable f empty, Push() not working")
	} else if verbose {
		fmt.Println("Test basic Push() to Fix instance passed")
	}
}

func testPushMax(t *testing.T, cap int) {
	f := fixedheap.NewFix(cap)
	m := m2

	for k, v := range m {
		i := &fixedheap.Item{Value: k, Priority: Float(v)}
		fixedheap.ManagedPush(f, i)
	}

	if f.Cap != cap || f.Heap.Len() > f.Cap {
		t.Errorf("variable f empty, Push() failed. cap: %d\n", cap)
	} else if verbose {
		fmt.Printf("Test ManagedPush() to Fix.Cap = %d instance passed\n", cap)
	}
}

func TestPushMax(t *testing.T) {
	capacity := []int{5, 3, 6, 12}
	for _, v := range capacity {
		testPushMax(t, v)
	}
}
