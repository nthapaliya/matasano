package utils_test

import (
	"math"
	"testing"

	"bytes"

	"github.com/nthapaliya/matasano/utils"
)

const PREFIX = "github.com/nthapaliya/matasano/"

func TestFrequencyDist(t *testing.T) {
	s := 0.0
	for _, v := range utils.FreqLowerCaseAlpha {
		s += math.Pow(10.0, v)
	}
	if math.Abs(s-1.0) > 0.005 {
		t.Errorf("Frequency percentage is way off: %f", s)
	}
}

func TestHammingDistance(t *testing.T) {
	b1 := []byte("this is a test")
	b2 := []byte("wokka wokka!!!")
	d := utils.HammingDistance(b1, b2)
	if d != 37 {
		t.Errorf("hamming distance incorrect, got %d", d)
	}
}

func TestHamDist2(t *testing.T) {
	LIM := byte(255)
	for a := byte(0); a < LIM; a++ {
		for b := byte(0); b < LIM; b++ {
			if utils.PopCountTable(a, b) != utils.PopCountSimple(a, b) {
				t.Errorf("breaks at %d and %d\n", a, b)
			}
		}
	}
}

func TestReadLines(t *testing.T) {
	s, err := utils.ReadLines("simpleTest.txt")
	if err != nil {
		t.Error(err)
	}
	if len(s) != 10 {
		t.Errorf("file array len=%d, should be %d", len(s), 10)
	}
}

func TestReadAll(t *testing.T) {
	s, err := utils.ReadAll("simpleTest.txt", "")
	if err != nil {
		t.Error(err)
	}
	if len(s) != 100 {
		t.Errorf("readall failed, got strlen %d", len(s))
	}
}

func TestMatrix(t *testing.T) {
	st, err := utils.ReadAll("simpleTest.txt", "")
	if err != nil {
		t.Error(err)
	}
	buff := []byte(st)
	m := utils.Split(buff, 20)

	outBytes := m.Join()
	if bytes.Compare(buff, outBytes) != 0 || &buff == &outBytes { // making sure they're not the same thing
		t.Errorf("two buffers do not match, should have len %d, got len %d",
			len(buff), len(outBytes))
	}

	m = utils.NewMatrix(10, 10)
	if l1, l2 := len(m), len(m[0]); l1 != 10 || l2 != 10 {
		t.Errorf("something wrong with new matrix, has len %d by %d", l1, l2)
	}
}
