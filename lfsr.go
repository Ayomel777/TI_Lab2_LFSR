package main

// Полином: x^37 + x^12 + x^10 + x^2 + 1
type LFSR struct {
	state     uint64
	degree    int
	stateMask uint64
	isZeroKey bool
}

func NewLFSR(binaryKey string) *LFSR {
	const degree = 37

	stateMask := (uint64(1) << degree) - 1

	var state uint64 = 0

	for i := 0; i < len(binaryKey) && i < degree; i++ {
		if binaryKey[i] == '1' {
			state |= uint64(1) << (degree - 1 - i)
		}
	}

	isZeroKey := state == 0

	return &LFSR{
		state:     state,
		degree:    degree,
		stateMask: stateMask,
		isZeroKey: isZeroKey,
	}
}

func (l *LFSR) NextBit() byte {
	if l.isZeroKey {
		return 0
	}

	outputBit := byte(l.state & 1)

	// Полином: x^37 + x^12 + x^10 + x^2 + 1
	bit36 := (l.state >> 36) & 1
	bit11 := (l.state >> 11) & 1
	bit9 := (l.state >> 9) & 1
	bit1 := (l.state >> 1) & 1

	newBit := bit36 ^ bit11 ^ bit9 ^ bit1

	l.state >>= 1

	l.state |= newBit << 36

	return outputBit
}

func (l *LFSR) NextByte() byte {
	if l.isZeroKey {
		return 0
	}

	var result byte = 0
	for i := 0; i < 8; i++ {
		bit := l.NextBit()
		result |= bit << i
	}
	return result
}

func (l *LFSR) GenerateKeyStream(length int) []byte {
	keyStream := make([]byte, length)

	if l.isZeroKey {
		return keyStream
	}

	for i := 0; i < length; i++ {
		keyStream[i] = l.NextByte()
	}
	return keyStream
}

func (l *LFSR) GetState() string {
	result := make([]byte, l.degree)
	for i := 0; i < l.degree; i++ {
		if (l.state>>(l.degree-1-i))&1 == 1 {
			result[i] = '1'
		} else {
			result[i] = '0'
		}
	}
	return string(result)
}

func (l *LFSR) IsZeroKey() bool {
	return l.isZeroKey
}
