package matasano_test

import (
	"testing"

	"github.com/nthapaliya/matasano"
)

func TestParser(t *testing.T) {
	s := "foo=bar&baz=qux&zap=zazzle"
	_, err := matasano.Parse(s)
	if err != nil {
		t.Error(err)
	}
}

func TestEncoder(t *testing.T) {
	s := "foo=bar&baz=qux&zap=zazzle"
	m, _ := matasano.Parse(s)
	if tmp := matasano.Encode(m); s != tmp {
		t.Errorf("error")
	}
}

func TestProfileFor(t *testing.T) {
	emails := []string{
		"foo@bar.com",
		"foo@bar.com&role=admin",
	}
	for _, email := range emails {
		encoded := matasano.ProfileFor(email)
		decoded, _ := matasano.Parse(encoded)
		if matasano.Encode(decoded) != encoded {
			t.Errorf("out does not match in, test cases fail")
		}
	}
}
