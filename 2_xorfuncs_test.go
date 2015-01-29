package matasano_test

// TESTXORFUNCS

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/nthapaliya/matasano"
	"github.com/nthapaliya/matasano/utils"
)

// Challenge 2
func TestFixedXOR(t *testing.T) {
	inputString := "1c0111001f010100061a024b53535009181c"
	xorString := "686974207468652062756c6c277320657965"
	supposedOut := "746865206b696420646f6e277420706c6179"

	inputBuffer, _ := hex.DecodeString(inputString)
	xorBuffer, _ := hex.DecodeString(xorString)

	outputByte, err := matasano.FixedXOR(inputBuffer, xorBuffer)
	if err != nil {
		t.Error(err)
	}
	output := make([]byte, len(inputString))
	hex.Encode(output, outputByte)
	if tmp := string(output); tmp != supposedOut {
		t.Errorf("fixed XOR Test Failed, got: %s", tmp)
	}
	logpass(t, 1)
}

// Challenge 3
func TestSingleCharXOR(t *testing.T) {
	input := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	buff, err := hex.DecodeString(input)
	if err != nil {
		t.Error(err)
	}
	_, err = matasano.Solve(buff)
	if err != nil {
		t.Error(err)
	}
	logpass(t, 2)
}

// Challenge 4
func TestXORDetection(t *testing.T) {
	filename := "txt/singleCharXORStrings.txt"

	s, err := matasano.DetectSingleCharXOR(filename)
	if err != nil {
		t.Error(err)
	}
	if s == "" {
		t.Errorf("detection failed, string empty")
	}
	logpass(t, 3)
}

func testEncodeKeyXOR(t *testing.T) {
	input := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	key := "ICE"

	outputString := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
	correctOutput, err := hex.DecodeString(outputString)
	if err != nil {
		t.Error(err)
	}

	testOutput := matasano.EncodeXORCipher([]byte(input), []byte(key))

	if len(correctOutput) != len(testOutput) {
		t.Errorf("output length is incorrect, need %d, got %d",
			len(correctOutput), len(testOutput))
	}
	if bytes.Compare(correctOutput, testOutput) != 0 {
		t.Errorf("repeating key xor did not work, got %s",
			hex.EncodeToString(testOutput))
	}
}

func testEncodeDecode(t *testing.T) {
	input := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	key := "ICE"

	a := matasano.EncodeXORCipher([]byte(input), []byte(key))
	b := matasano.DecodeXORCipher(a, []byte(key))

	if bytes.Compare([]byte(input), b) != 0 {
		t.Errorf("repeating key xor did not work, got %s",
			string(b))
	}
}

// Challenge 5
func TestRepeatingKeyXOR(t *testing.T) {
	testEncodeKeyXOR(t)
	testEncodeDecode(t)
	logpass(t, 4)
}

func fileHelper(t *testing.T) []byte {
	filename := "txt/repeatingXORData.txt"
	s, err := utils.ReadAll(filename, "")
	if err != nil {
		t.Error(err)
	}
	buff, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Error(err)
	}
	return buff
}

func TestCC(t *testing.T) {
	buff := fileHelper(t)
	if _, err := matasano.CountCoincidences(buff); err != nil {
		t.Error(err)
	}
	if _, err := matasano.GuessKeySize(buff); err != nil {
		t.Error(err)
	}
}

// Challenge 6
func TestBreakXOR(t *testing.T) {
	buff := fileHelper(t)
	_, err := matasano.BreakXOR(buff)
	if err != nil {
		t.Error(err)
	}
	logpass(t, 5)
}
