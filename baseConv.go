package matasano

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
)

// HexToBase64 converts a string encoded in Hex to Base64
func HexToBase64(inputstring string) (string, error) {
	if len(inputstring)%3 != 0 {
		return "", errors.New("incorrect length of input")
	}
	bytes, err := hex.DecodeString(inputstring)
	if err != nil {
		return "", err
	}
	output := base64.StdEncoding.EncodeToString(bytes)
	return output, nil
}

// Base64ToHex converts a string encoded in Base64 To Hex
func Base64ToHex(inputstring string) (string, error) {
	if len(inputstring)%2 != 0 {
		return "", errors.New("incorrect length of input")
	}
	bytes, err := base64.StdEncoding.DecodeString(inputstring)
	if err != nil {
		return "", err
	}
	output := hex.EncodeToString(bytes)
	return output, nil
}
