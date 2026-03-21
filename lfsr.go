package main

// LFSR структура для Linear Feedback Shift Register
// Полином: x^37 + x^12 + x^10 + x^2 + 1
type LFSR struct {
	state     uint64
	degree    int
	stateMask uint64
	isZeroKey bool
}

// NewLFSR создает новый LFSR с заданным начальным состоянием из бинарной строки (37 бит)
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

// NextBit генерирует следующий бит и сдвигает регистр
func (l *LFSR) NextBit() byte {
	if l.isZeroKey {
		return 0
	}

	outputBit := byte(l.state & 1)

	bit36 := (l.state >> 36) & 1
	bit11 := (l.state >> 11) & 1
	bit9 := (l.state >> 9) & 1
	bit1 := (l.state >> 1) & 1

	newBit := bit36 ^ bit11 ^ bit9 ^ bit1

	l.state >>= 1
	l.state |= newBit << 36

	return outputBit
}

// GenerateKeyStreamBits генерирует поток ключей как строку битов
func (l *LFSR) GenerateKeyStreamBits(length int) string {
	bits := make([]byte, length)

	if l.isZeroKey {
		for i := 0; i < length; i++ {
			bits[i] = '0'
		}
		return string(bits)
	}

	for i := 0; i < length; i++ {
		if l.NextBit() == 1 {
			bits[i] = '1'
		} else {
			bits[i] = '0'
		}
	}
	return string(bits)
}

// NextByte генерирует следующий байт (8 бит)
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

// GenerateKeyStream генерирует поток ключей заданной длины в байтах
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

// GetState возвращает текущее состояние регистра как бинарную строку
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

// IsZeroKey возвращает true, если ключ состоит из всех нулей
func (l *LFSR) IsZeroKey() bool {
	return l.isZeroKey
}
