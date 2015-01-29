package matasano_test

// TESTBASECONV

import (
	"testing"

	"github.com/nthapaliya/matasano"
)

const (
	hexstring = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	b64string = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
)

func testHexToBase64(t *testing.T) {
	input := hexstring
	output := b64string

	s, err := matasano.HexToBase64(input)
	if err != nil {
		t.Error(err)
	}
	if s != output {
		t.Errorf("incorrect output: %s", s)
	}
}

func testBase64ToHex(t *testing.T) {
	input := b64string
	output := hexstring
	s, err := matasano.Base64ToHex(input)
	if err != nil {
		t.Error(err)
	}
	if s != output {
		t.Errorf("incorrect output: %s", s)
	}
}

func testBackAndForth(t *testing.T) {
	s1 := b64string
	tmp, err := matasano.Base64ToHex(s1)
	if err != nil {
		t.Error(err)
	}
	tmp, err = matasano.HexToBase64(tmp)
	if err != nil {
		t.Error(err)
	}
	if tmp != s1 {
		t.Errorf("conversion from B64 to Hex and back failed. Got: %s", tmp)
	}

	s2 := hexstring
	tmp, err = matasano.HexToBase64(s2)
	if err != nil {
		t.Error(err)
	}
	tmp, err = matasano.Base64ToHex(tmp)
	if err != nil {
		t.Error(err)
	}
	if tmp != s2 {
		t.Errorf("conversion from Hex to B64 and back failed")
	}
}

// Challenge 1
func TestBaseConv(t *testing.T) {
	testHexToBase64(t)
	testBase64ToHex(t)
	testBackAndForth(t)
	logpass(t, 0)
}
