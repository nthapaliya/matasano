package matasano

import "errors"

// PKCS padding
func PadToNext(in []byte, blocksize int) ([]byte, error) {
	outlen := blocksize * (1 + len(in)/blocksize)
	return PadBuffer(in, outlen)
}
func PadBuffer(in []byte, blocksize int) ([]byte, error) {
	n := len(in)
	if n > blocksize {
		return nil, errors.New("blocksize smaller than input block")
	}
	out := make([]byte, blocksize)
	copy(out, in)
	pad := byte(blocksize - n)
	for i := n; i < len(out); i++ {
		out[i] = pad
	}
	return out, nil
}

// If padding is valid, returns the stripped buffer, else an error
func ValidatePadding(in []byte) ([]byte, error) {
	blocksize := len(in)
	lastbyte := int(in[blocksize-1])
	if lastbyte > blocksize {
		return nil, errors.New("padding byte is larger than blocksize")
	}

	for i := blocksize - lastbyte; i < blocksize; i++ {
		if in[i] != byte(lastbyte) {
			return nil, errors.New("invalid padding")
		}
	}
	out := make([]byte, blocksize-lastbyte)

	copy(out, in)
	return out, nil
}
