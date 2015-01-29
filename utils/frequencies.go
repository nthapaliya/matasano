package utils

import "math"

var _freqLowerCaseAlpha = map[byte]float64{
	'e': 13.0001,
	't': 9.056,
	'a': 8.167,
	'o': 7.507,
	'i': 6.966,
	'n': 6.749,
	's': 6.327,
	'h': 6.094,
	'r': 5.987,
	'd': 4.253,
	'l': 4.025,
	'c': 2.782,
	'u': 2.758,
	'm': 2.406,
	'w': 2.360,
	'f': 2.228,
	'g': 2.015,
	'y': 1.974,
	'p': 1.929,
	'b': 1.492,
	'v': 0.978,
	'k': 0.772,
	'j': 0.153,
	'x': 0.150,
	'q': 0.095,
	'z': 0.074,
}
var _freq = map[byte]float64{
	'a': 0.0651738,
	'b': 0.0124248,
	'c': 0.0217339,
	'd': 0.0349835,
	'e': 0.1041442,
	'f': 0.0197881,
	'g': 0.0158610,
	'h': 0.0492888,
	'i': 0.0558094,
	'j': 0.0009033,
	'k': 0.0050529,
	'l': 0.0331490,
	'm': 0.0202124,
	'n': 0.0564513,
	'o': 0.0596302,
	'p': 0.0137645,
	'q': 0.0008606,
	'r': 0.0497563,
	's': 0.0515760,
	't': 0.0729357,
	'u': 0.0225134,
	'v': 0.0082903,
	'w': 0.0171272,
	'x': 0.0013692,
	'y': 0.0145984,
	'z': 0.0007836,
	' ': 0.1918182,
}
var FreqLowerCaseAlpha map[byte]float64

func init() {
	FreqLowerCaseAlpha = make(map[byte]float64)
	for k, v := range _freq {
		FreqLowerCaseAlpha[k] = math.Log10(v)
	}

}