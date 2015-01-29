package matasano

import (
	"container/heap"
	"encoding/hex"
	"errors"
	"math"
	"strconv"

	"github.com/nthapaliya/matasano/heap"
	"github.com/nthapaliya/matasano/utils"
)

// TODO: MAJOR MAJOR CLEANUP
// TODO: HEAP MAYBE NOT NEEDED, REMOVE DEPENDENCY
// TODO: finish this. you already made it
// work, but it's too half assed to be any good

// FixedXOR is a simple function that takes two byte buffers and
// returns a byte buffer with the XOR'ed product
//
func Add(buffA, buffB []byte) ([]byte, error) {
	return FixedXOR(buffA, buffB)
}

func FixedXOR(buffA, buffB []byte) ([]byte, error) {
	if len(buffA) != len(buffB) {
		return nil, errors.New("lengths of input buffers do not match")
	}
	return EncodeXORCipher(buffA, buffB), nil
}

// EncodeXORCipher takes two non equal byte buffers and outputs an encoded
// buffer using the second parameter as a key
// TODO: Make sure that this doesn't explode when keylen > bufferlen
//
func EncodeXORCipher(input, key []byte) []byte {
	keylen := len(key)
	output := make([]byte, len(input))
	for i := range output {
		output[i] = input[i] ^ key[i%keylen]

	}
	return output
}

// DecodeXORCipher decodes buffer with the given key
//
func DecodeXORCipher(input, key []byte) []byte {
	return EncodeXORCipher(input, key)
}

// XORwithChar xor's buffer with single char taken as second parameter
// TODO: Check if this can be made unexported
//
func XORwithChar(buffer []byte, b byte) []byte {
	out := make([]byte, len(buffer))
	for i := 0; i < len(buffer); i++ {
		out[i] = buffer[i] ^ b
	}
	return out
}

// DetectSingleCharXOR is not very useful, it reads a file (given as first param)
// and tries to detect which line is XOR-ed with a single character, and if it
// finds one, outputs it. If not, returns an error
//
func DetectSingleCharXOR(name string) (string, error) {
	s, err := utils.ReadLines(name)
	if err != nil {
		return "", err
	}
	for _, value := range s {
		buff, err := hex.DecodeString(value)
		if err != nil || buff == nil {
			continue
		} else if buff != nil {
			return string(buff), nil
		}
	}
	return "", errors.New("empty string returned")
}

/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////

// TODO: BETTER NAMING CONVENTIONS

// Solve takes a single byte xor-encoded byte slice, then tries to solve it by:
// 	a) iterating through a list of possible bytes (0 to 254)
// 	b) xor-ing the buffer with each byte
// 	c) constructing a histogram of the xor-ed buffer
// 	d) comparing it with a pre-computed histogram
// If found, returns the BYTE character buffer is xor'ed with, else an error
//
func Solve(buffer []byte) (byte, error) {
	b, err := _solve(buffer)
	if err != nil {
		return 0, err
	}
	return b, nil
}

func _solve(buffer []byte) (byte, error) {
	score, maxv := -1000, byte(0)
	for v := byte(0); v < 255; v++ {
		outputBuffer := XORwithChar(buffer, v)
		// wrap this functionality vv
		if byteSlicePrintable(outputBuffer) {
			m := makeHistogram(outputBuffer)
			if s := compareHistograms(utils.FreqLowerCaseAlpha, m); s > score {
				score = s
				maxv = v
			}
		}
		// wrap this functionality ^^
	}
	if score == -1000 || maxv == byte(0) {
		return 0, errors.New("empty string returned")
	}
	return maxv, nil
}

// Scores the two inputs on relative similarity.
// First parameter is the "good" histogram, second is the
// one we're scoring
//
func compareHistograms(a, b map[byte]float64) (score int) {
	score = 0
	for k, v := range b {
		score -= int(math.Abs(v - a[k]))
	}
	return score
}

// Histograms are normalized to 1, then the numbers are math.Log10()'ed
// for easier comparisions
//
func makeHistogram(buffer []byte) map[byte]float64 {
	m := make(map[byte]float64)
	for _, v := range buffer {
		m[v] += 1.0
	}
	for k := range m {
		m[k] /= float64(len(buffer))
		m[k] = math.Log10(m[k])
	}
	return m
}

// returns if contains printable runes including '\n'
func byteSlicePrintable(buff []byte) bool {
	for _, v := range buff {
		if vrune := rune(v); !strconv.IsPrint(vrune) && vrune != rune('\n') {
			return false
		}
	}
	return true
}

/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////

// TODO: cleanup and complete this

// BreakXOR takes a byte buffer, guesses the length of the key, and tries
// to decode the XOR-ed buffer. Returns the key (as a []byte) if found, else
// an error
//
func BreakXOR(buff []byte) ([]byte, error) {
	// keysize := MostLikelyKeysizes(buff)
	//MostLikelyKeysizes(buff)
	keysize := 29
	key, err := bx(buff, keysize)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func bx(buff []byte, keysize int) (key []byte, err error) {
	m := utils.Split(buff, keysize)
	key = make([]byte, keysize)

	for i, row := range m {
		solvedByte, err := Solve(row)
		if err != nil {
			return nil, err
		}
		key[i] = solvedByte
	}
	return key, nil
}

////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

// guess key size, GuessKeySize, guessKeySize, guesskeysize portion of file

// This package implements functions that will mostly try to guess the
// key length of an xor-ed cipher.
// Three methods that I will try here:
//
//   a) First, the method given in the matasano challenge text. Take a block of
//      length = n which you are guessing. Take the next block of length = n.
//      Take the hamming distance of these blocks. Try with a range of n. The
//      value of n that gives the smallest hamming distance is probably the
//      key length.
//
//   b) A simple implementation of the Kasiski attack described at
//      http://en.wikipedia.org/wiki/Kasiski_examination#Superposition
//      What we do here is as follows:
//        Take two copies of the text, shift one text one letter to the right, and
//        superpose them. Count the number of 'coincidences', ie where the top
//        and bottom are the same. Continue with shifting and counting.
//        Apparently, the number of coincidences rises sharply when the shifted
//        number is a multiple of keylen.
//   c) Using the Index of Coincidence as outlined here:
//      http://en.wikipedia.org/wiki/Index_of_coincidence#Calculation
//      This method is a little more involved, but can give us a statistically
//      stronger method of guessing the keylength. I hope I don't really have to
//      implement this. I'll try b) first

const inf = math.MaxFloat64

const FIX_CAP = 5

// *******************************************************************
// *******************************************************************

// Type declarations that implement fixedheap.Numeric
type Int int

func (a Int) LessThan(b fixedheap.Numeric) bool {
	c := b.(Int)
	return a < c
}

type Float float64

func (a Float) LessThan(b fixedheap.Numeric) bool {
	c := b.(Float)
	return a < c
}

// TODO: return a union of factors, most likely
// func MostLikelyKeysizes(buff []byte) ([]int, error) {
// 	r1, _ := CountCoincidences(buff)
// 	r2, _ := GuessKeySize(buff)
// 	fmt.Println(r1)
// 	for _, v := range r1 {
// 		fmt.Println(factor(v))
// 	}
// 	fmt.Println(r2)
// 	for _, v := range r2 {
// 		fmt.Println(factor(v))
// 	}
//
// 	return nil, nil
// }
//
// func factor(n int) []int {
// 	factors := []int{}
// 	d := 2
// 	for n > 1 {
// 		if n%d == 0 {
// 			n /= d
// 			factors = append(factors, d)
// 		} else {
// 			d++
// 		}
//
// 	}
// 	return factors
// }

// *******************************************************************
// *******************************************************************

// MOST Coincidences are more likely
// TODO rewrite func name.

// Using "Coincidence Counting". Wikipedia: Breaking Vigniere cipher
func CountCoincidences(buff []byte) ([]int, error) {
	f := fixedheap.NewFix(FIX_CAP) // POINTER TO f

	for value := 2; value < len(buff); value++ {
		priority := coincidences(buff, value)
		item := &fixedheap.Item{Value: value, Priority: Int(priority)}
		fixedheap.ManagedPush(f, item)
	}
	r := make([]int, f.Heap.Len())
	for i, v := range f.Heap {
		k, _ := v.Value.(int), v.Priority
		r[i] = k
	}
	return r, nil
}

func coincidences(buff []byte, shiftlen int) (count int) {
	lbuff := len(buff)
	for i := range buff {
		if buff[i] == buff[(i+shiftlen)%lbuff] {
			count++
		}
	}
	return count
}

///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

// Least HammingDistances are more likely
// Based on taking equal block sizes and returning the least "normalized" one
func GuessKeySize(buff []byte) ([]int, error) {
	f := fixedheap.NewFix(FIX_CAP) // pointer to fixedheap

	for i := 2; i < 100; i++ {
		score := keysizeBlockHamming(buff, i)
		item := fixedheap.Item{Value: i, Priority: Float(score)}
		heap.Push(&f.Heap, item)
	}
	var outputLength int
	if f.Cap < f.Heap.Len() {
		outputLength = f.Cap
	} else {
		outputLength = f.Heap.Len()
	}
	r := make([]int, outputLength)
	for i := range r {
		item := (heap.Pop(&f.Heap)).(*fixedheap.Item)
		r[i] = item.Value.(int)
	}
	return r, nil
}

func keysizeBlockHamming(buff []byte, k int) (fscore float64) {
	// k == keysize

	buffLen := len(buff)
	ratio := buffLen / k
	if k < 2 || ratio < 3 {
		return inf
	}

	var score int
	loops := ratio - 2

	for i := 0; i <= loops; i++ {
		a := buff[(i+0)*k : (i+1)*k]
		b := buff[(i+1)*k : (i+2)*k]
		t := utils.HammingDistance(a, b)
		if t == -1 {
			panic("Hamming is -ve, check algorithm")
		}
		score += t
	}
	// normalized score
	fscore = float64(score) / float64(loops*k)
	return fscore
}
