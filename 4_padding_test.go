package matasano_test

import (
	"bytes"
	"testing"

	"github.com/nthapaliya/matasano"
)

// Challenge 9
func TestPadding(t *testing.T) {
	in := []byte("YELLOW SUBMARINE")
	blocksize := 20
	out, err := matasano.PadBuffer(in, blocksize)
	if err != nil {
		t.Error(err)
	}
	if len(out) != blocksize {
		t.Errorf("blocksize not 20, got %d", len(out))
	}
	if out[blocksize-1] != 4 {
		t.Errorf("last bit not 4, got %d", out[blocksize-1])
	}
	if _, err := matasano.PadBuffer(in, 4); err == nil {
		t.Error("something wrong with algorithm, should fail here")
	}
	logpass(t, 8)
}

// Challenge 14
func TestPaddingValidation(t *testing.T) {
	s1 := []byte("ICE ICE BABY\x04\x04\x04\x04")
	s2 := []byte("ICE ICE BABY\x05\x05\x05\x05")
	s3 := []byte("ICE ICE BABY\x01\x02\x03\x04")
	s4 := []byte("ICE ICE BABY\xff\xff\xff\xff")

	s, err := matasano.ValidatePadding(s1)
	if err != nil {
		t.Error(err)
	} else if bytes.Compare(s, []byte("ICE ICE BABY")) != 0 {
		t.Errorf("problem in padding validation")
	}

	s, err = matasano.ValidatePadding(s2)
	if err == nil || s != nil {
		t.Errorf("padding validation failed")
	}

	s, err = matasano.ValidatePadding(s3)
	if err == nil || s != nil {
		t.Errorf("padding validation failed")
	}

	s, err = matasano.ValidatePadding(s4)
	if err == nil || s != nil {
		t.Errorf("padding validation failed")
	}
	logpass(t, 14)
}

func TestPadToNext(t *testing.T) {
	blocksize := 16
	in := []byte("this is a string i just typed in with no regard to anything")
	b, err := matasano.PadToNext(in, blocksize)
	if err != nil {
		t.Error(err)
	}
	if n := len(b); n%blocksize != 0 {
		t.Errorf("pad to next failed")
	}
}
