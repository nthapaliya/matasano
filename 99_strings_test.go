package matasano_test

import (
	"fmt"
	"testing"
)

var teststrings = []string{
	"1.  Convert hex to base64 and back",
	"2.  Fixed XOR",
	"3.  Single-character XOR Cipher",
	"4.  Detect single-character XOR",
	"5.  Repeating-key XOR Cipher",
	"6.  Break repeating-key XOR",
	"7.  AES in ECB Mode",
	"8.  Detecting ECB",
	"9.  Implement PKCS#7 padding",
	"10. Implement CBC Mode",
	"11. Oracle function to detect ECB",
	"12. Byte-at-a time ECB decryption, Full control version",
	"13. ECB cut-and-paste",
	"14. Byte-at-a-time ECB decryption, Partial control version",
	"15. PKCS#7 padding validation",
	"16. CBC bit flipping",
}

var testpass = make([]bool, len(teststrings))

func logpass(t *testing.T, n int) {
	if !t.Failed() {
		testpass[n] = true
	}
}

func _TestMain(t *testing.T) {
	for n, v := range testpass {
		if v {
			fmt.Printf("PASS:\t%s\n", teststrings[n])

		} else {
			fmt.Printf("    :\t%s\n", teststrings[n])
		}
	}
}
