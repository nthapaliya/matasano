package matasano

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	// TODO: IMPORTANT, fixed seed will produce same "random" numbers
	// each time. Good for testing, bad for everything else
	// RND = rand.New(rand.NewSource(100.0))
	t := time.Now().UTC().UnixNano()
	RND = rand.New(rand.NewSource(t))
	KEY = GenRandBuffer(aes.BlockSize)
}

var (
	RND *rand.Rand
	KEY []byte
)

// AES-rindjael encryption. This files has two versions, the simpler
// ECB mode, and the slightly more complicated CBC mode

// DetectAES tries to determine whether a 160 byte buffer is encrypted
// using ECB block encoding
// Useful only in context of this challenge, not very general
func AESDetectECB(buff []byte) bool {
	blocksize := aes.BlockSize
	inputsize := len(buff)
	// TODO: cleanup, don't use panic
	if inputsize%blocksize != 0 {
		// fmt.Println("PANIC")
		// return false
		panic("Wrong")
	}
	r := inputsize / blocksize
	for i := 0; i < r; i++ {
		b1 := buff[blocksize*i : blocksize*(i+1)]
		for j := i + 1; j < r; j++ {
			b2 := buff[blocksize*j : blocksize*(j+1)]
			if bytes.Compare(b1, b2) == 0 {
				return true
			}
		}
	}
	return false
}

func initAESBlock(buff, key, iv []byte) (cipher.Block, error) {
	if len(key) != 16 {
		return nil, errors.New("key is not proper length")
	}
	if iv != nil && len(iv) != 16 {
		return nil, errors.New("initialization vector is of improper length")
	}

	blocksize := len(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if b := block.BlockSize(); b != blocksize {
		return nil, fmt.Errorf("something wrong, block size %d", b)
	}
	if len(buff)%blocksize != 0 {
		return nil, errors.New("buffer length is not a multiple of blocksize," +
			" calling function should pad")
	}

	return block, nil
}

// AESDecrypt decrypts the buffer with the given key using AES ECB block mode
func AESDecryptECB(buff, key []byte) error {
	block, err := initAESBlock(buff, key, nil)
	if err != nil {
		return err
	}
	blocksize := len(key)
	for i := 0; i < len(buff); i += blocksize {
		block.Decrypt(buff[i:i+blocksize], buff[i:i+blocksize])
	}
	return nil
}

func AESEncryptECB(buff, key []byte) error {
	block, err := initAESBlock(buff, key, nil)
	if err != nil {
		return err
	}
	blocksize := len(key)
	for i := 0; i < len(buff); i += blocksize {
		block.Encrypt(buff[i:i+blocksize], buff[i:i+blocksize])
	}

	return nil
}

func AESDecryptCBC(buff, key, iv []byte) error {
	block, err := initAESBlock(buff, key, iv)
	if err != nil {
		return err
	}
	blocksize := len(key)

	scratch := make([]byte, blocksize)     // these two arrays so that we don't
	cipherblock := make([]byte, blocksize) // change key and iv as side effect
	copy(cipherblock, iv)

	for i := 0; i < len(buff); i += blocksize {
		block.Decrypt(scratch, buff[i:i+blocksize])
		scratch, _ = Add(scratch, cipherblock)
		copy(cipherblock, buff[i:i+blocksize])
		copy(buff[i:i+blocksize], scratch)
	}
	return nil
}

func AESEncryptCBC(buff, key, iv []byte) error {
	block, err := initAESBlock(buff, key, iv)
	if err != nil {
		return err
	}
	blocksize := len(key)

	scratch := make([]byte, blocksize)     // these two arrays so that we don't
	cipherblock := make([]byte, blocksize) // change key and iv as side effect
	copy(cipherblock, iv)

	for i := 0; i < len(buff); i += blocksize {
		scratch, _ = Add(buff[i:i+blocksize], cipherblock)
		block.Encrypt(scratch, scratch)
		copy(cipherblock, scratch)
		copy(buff[i:i+blocksize], scratch)
	}
	return nil
}

// Consider making this non-exported later
func GenRandBuffer(keylen int) []byte {
	out := make([]byte, keylen)
	for i := range out {
		out[i] = byte(RND.Intn(256))
	}
	return out
}

// stupid function
func Oracle(in []byte) ([]byte, error) {
	blocksize := aes.BlockSize

	appendcount := 5 + RND.Intn(6)
	outlen := (1 + (len(in)+2*appendcount)/blocksize) * blocksize

	out := GenRandBuffer(outlen)
	in, err := PadToNext(in, blocksize)
	if err != nil {
		return nil, err
	}
	copy(out[appendcount:], in)

	if flip := RND.Intn(2); flip == 0 {
		// fmt.Println("ecb") // for debug
		randkey := GenRandBuffer(16)
		err := AESEncryptECB(out, randkey)
		if err != nil {
			return nil, err
		}
	} else {
		// fmt.Println("cbc") // for debug
		randkey, iv := GenRandBuffer(16), GenRandBuffer(16)
		err := AESEncryptCBC(out, randkey, iv)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func ECBOracle(in, unknown []byte) ([]byte, error) {
	blocksize := aes.BlockSize

	in, err := PadToNext(append(in, unknown...), blocksize)
	if err != nil {
		return nil, err
	}

	out := make([]byte, len(in))
	copy(out, in)

	if err := AESEncryptECB(out, KEY); err != nil {
		return nil, err
	}
	return out, nil
}

func detectBlockSize(upto int) (int, error) {
	in := make([]byte, upto)
	for i := range in {
		in[i] = 'A'
	}

	for i := range in {
		buff, err := ECBOracle(in[:i], nil)
		if err != nil {
			return -1, err
		}
		if AESDetectECB(buff) {
			return i / 2, nil
		}
	}
	return -1, errors.New("unable to detect block size")
}

func AESBreakECBEncryption(unknown []byte) []byte {
	blocksize, err := detectBlockSize(100)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	in := make([]byte, blocksize)
	for i := range in {
		in[i] = 'A'
	}

	m := make(map[string]byte)
	for i := 0; i < 256; i++ {
		in[len(in)-1] = byte(i)
		temp, _ := ECBOracle(in, nil)
		temp = temp[:blocksize]
		m[string(temp)] = byte(i)
	}

	in = in[:blocksize-1]

	outbuffer := make([]byte, len(unknown))
	for k := range outbuffer {
		out, _ := ECBOracle(in, unknown[k:])
		out = out[:blocksize]
		outbuffer[k] = m[string(out)]
	}
	return outbuffer
}

func ProfileOracle(email string) ([]byte, error) {
	cleartext := ProfileFor(email)
	return ECBOracle([]byte(cleartext), nil)
}
