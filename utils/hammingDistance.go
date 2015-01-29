package utils

var lookup [256]int

func init() {
	initLookup()
}

// Stolen from H. Warren's Hacker's Delight
func initLookup() {
	lookup[0] = 0
	for i := 1; i < 256; i++ {
		lookup[i] = lookup[i/2] + (i & 1)
	}
}

func HammingDistance(b1, b2 []byte) int {
	if len(b1) != len(b2) {
		return -1
	}
	distance := 0
	for i := 0; i < len(b1); i++ {
		distance += lookup[b1[i]^b2[i]]
	}
	return distance
}

// TODO: Benchmark whether function calls are more expensive.
// For now, I don't use this, but just 'inline' this functionality
// in the inner loop in HammingDistance, at the expense of readability later
func PopCountTable(x, y byte) int {
	a := x ^ y
	return lookup[a]
}

// Keeping it in here just in case, for debugging etc.
// Procedure will work even for ints.
func PopCountSimple(x, y byte) int {
	distance, val := 0, x^y

	for val > 0 {
		val &= val - 1
		distance++
	}
	return distance
}
