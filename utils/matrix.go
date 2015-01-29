package utils

type Matrix [][]byte

func Split(buff []byte, keysize int) Matrix {
	m := Matrix{}
	for i := 0; i < keysize; i++ {
		m = append(m, []byte{})
	}
	for i, v := range buff {
		m[i%keysize] = append(m[i%keysize], v)
	}
	return m
}

func (m Matrix) Join() []byte {
	// below here it's basic bookkeeping
	keysize := len(m)
	if keysize < 2 {
		return nil
	}
	var byteLen int
	for _, v := range m {
		byteLen += len(v)
	}
	outBuffer := make([]byte, byteLen)
	// up till here it's basic bookkeeping

	for i := range outBuffer {
		outBuffer[i] = m[i%keysize][i/keysize]
	}
	return outBuffer

}

func NewMatrix(rows, columns int) Matrix {
	m := Matrix{}
	for i := 0; i < rows; i++ {
		m = append(m, []byte{})
	}
	if columns > 0 {
		for i := 0; i < rows; i++ {
			m[i] = make([]byte, columns)
		}
	}
	return m
}
