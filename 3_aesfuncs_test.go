package matasano_test

// TESTAES

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"bytes"

	"github.com/nthapaliya/matasano"
	"github.com/nthapaliya/matasano/utils"
)

const UNKNOWN = "Um9sbGluJyBpbiBteSA1LjA" +
	"KV2l0aCBteSByYWctdG9wIG" +
	"Rvd24gc28gbXkgaGFpciBjY" +
	"W4gYmxvdwpUaGUgZ2lybGll" +
	"cyBvbiBzdGFuZGJ5IHdhdml" +
	"uZyBqdXN0IHRvIHNheSBoaQ" +
	"pEaWQgeW91IHN0b3A/IE5vL" +
	"CBJIGp1c3QgZHJvdmUgYnkK"

// Challenge 7
func TestAESDecryptECB(t *testing.T) {
	filename := "txt/aesData.txt"
	s, err := utils.ReadAll(filename, "")
	if err != nil {
		t.Error(err)
	}
	buff, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Error(err)
	}

	err = matasano.AESDecryptECB(buff, []byte("YELLOW SUBMARINE"))
	if err != nil || buff == nil {
		t.Error(err)
	}
	logpass(t, 6)
}

// Challenge 8
func TestAESDetectECB(t *testing.T) {
	filename := "txt/detectAes.txt"
	strings, err := utils.ReadLines(filename)
	if err != nil {
		t.Error(err)
	}
	detected := false
	for _, s := range strings {
		buff, _ := hex.DecodeString(s) // ignore error for now
		if matasano.AESDetectECB(buff) {
			detected = true
		}
	}
	if !detected {
		t.Errorf("Did not detect ecb string")
	}
	logpass(t, 7)

}

func TestAESEncryptECB(t *testing.T) {
	filename := "txt/aesData.txt"
	s, err := utils.ReadAll(filename, "")
	if err != nil {
		t.Error(err)
	}
	buff, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Error(err)
	}
	originalbuff := make([]byte, len(buff))
	copy(originalbuff, buff)

	err = matasano.AESDecryptECB(buff, []byte("YELLOW SUBMARINE"))
	if err != nil || buff == nil {
		t.Error(err)
	}
	err = matasano.AESEncryptECB(buff, []byte("YELLOW SUBMARINE"))
	if err != nil || buff == nil {
		t.Error(err)
	}
	if bytes.Compare(originalbuff, buff) != 0 {
		fmt.Printf("original buff \n%s\ngot\n%s\n", string(originalbuff), string(buff))
		t.Errorf("decryption did not work")
	}
}

// Challenge 10
func TestAESDecryptCBC(t *testing.T) {
	s, err := utils.ReadAll("txt/10.txt", "")
	if err != nil {
		t.Error(err)
	}

	buff, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Error(err)
	}
	key := []byte("YELLOW SUBMARINE")
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	if err := matasano.AESDecryptCBC(buff, key, iv); err != nil {
		t.Error(err)
	}
	logpass(t, 9)
}

func TestAESEncryptCBC(t *testing.T) {
	s, err := utils.ReadAll("txt/10.txt", "")
	if err != nil {
		t.Error(err)
	}

	buff, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Error(err)
	}
	originalbuff := make([]byte, len(buff))
	copy(originalbuff, buff)

	key := []byte("YELLOW SUBMARINE")
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	if err := matasano.AESDecryptCBC(buff, key, iv); err != nil {
		t.Error(err)
	}
	if err := matasano.AESEncryptCBC(buff, key, iv); err != nil {
		t.Error(err)
	}

	if bytes.Compare(buff, originalbuff) != 0 {
		t.Errorf("decoding CBC did not work")
	}
}

func _TestGenRandomKey(t *testing.T) {
	out := [][]byte{
		[]byte{7, 56, 216, 204, 46, 140, 253, 81, 248, 93, 25, 198, 65, 60, 123, 44},
		[]byte{92, 50, 159, 83, 147, 170, 26, 120, 29, 113, 179, 20, 59, 87, 99, 122},
		[]byte{137, 17, 161, 118, 24, 190, 35, 126, 76, 134, 68, 26, 3, 13, 191, 208},
		[]byte{105, 110, 32, 183, 95, 199, 113, 21, 46, 109, 185, 66, 216, 156, 133, 230},
		[]byte{199, 10, 226, 31, 6, 31, 76, 122, 26, 88, 240, 154, 243, 78, 167, 247},
		[]byte{27, 84, 128, 83, 199, 219, 23, 250, 50, 110, 97, 178, 53, 10, 240, 8},
	}

	// NOTE: The above values only work for rnd = rnd.New(rand.NewSource(100.0))
	// ie the seed has to be 100.0 otherwise it won't work.
	// Consider making this function non-exported

	for i := 0; i < 6; i++ {
		b := matasano.GenRandBuffer(16)
		if bytes.Compare(b, out[i]) != 0 {
			t.Errorf("genrandkey produces wrong output")
		}
	}
}

func TestOracle(t *testing.T) {
	s := "this is a stupid string that I just typed out"
	buff := []byte(s)
	out, err := matasano.Oracle(buff)
	if err != nil || out == nil {
		t.Errorf("oracle did not work")
	}
}

// Challenge 11
func TestDetection(t *testing.T) {
	s := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	buff := []byte(s)
	for i := 0; i < 100; i++ {
		out, err := matasano.Oracle(buff)
		if err != nil || out == nil {
			t.Errorf("oracle did not work")
		}
		if matasano.AESDetectECB(out) {
			logpass(t, 10)
			return
		}
	}
	// if they get here, it's a problem
	t.Fail()
}

func TestECBOracle(t *testing.T) {
	unknown, err := base64.StdEncoding.DecodeString(UNKNOWN)
	if err != nil {
		t.Error(err)
	}
	in := []byte("AAAAAAAAAAAAAAAAAAAAAAAAA")

	if out, err := matasano.ECBOracle(in, unknown); err != nil {
		t.Error(err)
	} else if len(out)%aes.BlockSize != 0 {
		t.Errorf("something went wrong with oracle")
	}
}

// Challenge 12
func TestAESBreakECBEncryption(t *testing.T) {
	unknown, err := base64.StdEncoding.DecodeString(UNKNOWN)
	if err != nil {
		t.Error(err)
	}
	out := matasano.AESBreakECBEncryption(unknown)
	if s := string(out); bytes.Compare(out, unknown) != 0 {
		t.Errorf("breaking ECB failed, got %s", s)
	}
	logpass(t, 11)
}

func TestProfileOracle(t *testing.T) {
	_, err := matasano.ProfileOracle("foo.bar@gmail.com")
	if err != nil {
		t.Error(err)
	}
}

func TestEncryptingProfiles(t *testing.T) {
	email := "foo.bar@gmail.com"

	key := matasano.GenRandBuffer(aes.BlockSize)
	p := matasano.ProfileFor(email) // string

	profile := []byte(p)
	buff := append([]byte{}, profile...) // copy profile to new slice buff

	matasano.AESEncryptECB(buff, key)
	matasano.AESDecryptECB(buff, key)

	if bytes.Compare(buff, profile) != 0 {
		t.Errorf("error in decrypting profile")
	}
}

func TestBreakProfileEncryption(t *testing.T) {

}
