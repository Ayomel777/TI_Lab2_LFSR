package main

// LFSR структура для Linear Feedback Shift Register
// Полином: x^37 + x^12 + x^10 + x^2 + 1
// Тапы: позиции 37, 12, 10, 2 (1 не считается, это выход)
type LFSR struct {
	state     uint64
	degree    int
	tapMask   uint64
	stateMask uint64
}

// NewLFSR создает новый LFSR с заданным начальным состоянием из бинарной строки
func NewLFSR(binaryKey string) *LFSR {
	const degree = 37

	// Создаем маску для тапов: позиции 37, 12, 10, 2
	// В нашей нумерации (0-индексация): 36, 11, 9, 1
	// x^37 + x^12 + x^10 + x^2 + 1
	// Тапы: биты 36 (x^37), 11 (x^12), 9 (x^10), 1 (x^2)
	tapMask := uint64(1)<<36 | uint64(1)<<11 | uint64(1)<<9 | uint64(1)<<1

	// Маска состояния для 37 бит
	stateMask := (uint64(1) << degree) - 1

	// Инициализируем состояние из ключа
	var state uint64 = 0

	// Берем первые 37 бит из ключа
	keyLen := len(binaryKey)
	bitsToUse := degree
	if keyLen < bitsToUse {
		bitsToUse = keyLen
	}

	for i := 0; i < bitsToUse; i++ {
		if binaryKey[i] == '1' {
			state |= uint64(1) << (degree - 1 - i)
		}
	}

	// Если ключ длиннее 37 бит, XOR-им остальные биты
	if keyLen > degree {
		for i := degree; i < keyLen; i++ {
			if binaryKey[i] == '1' {
				pos := (i - degree) % degree
				state ^= uint64(1) << (degree - 1 - pos)
			}
		}
	}

	// Убедимся, что состояние не нулевое (LFSR застрянет на 0)
	if state == 0 {
		state = 1
	}

	return &LFSR{
		state:     state,
		degree:    degree,
		tapMask:   tapMask,
		stateMask: stateMask,
	}
}

// NextBit генерирует следующий бит
func (l *LFSR) NextBit() byte {
	// Выходной бит - младший бит состояния
	outputBit := byte(l.state & 1)

	// Вычисляем бит обратной связи через XOR тапов
	feedback := l.state & l.tapMask

	// Считаем четность (XOR всех бит в feedback)
	feedbackBit := uint64(0)
	temp := feedback
	for temp != 0 {
		feedbackBit ^= temp & 1
		temp >>= 1
	}

	// Сдвигаем регистр вправо
	l.state >>= 1

	// Вставляем бит обратной связи в старший бит
	l.state |= feedbackBit << (l.degree - 1)

	// Применяем маску
	l.state &= l.stateMask

	return outputBit
}

// NextByte генерирует следующий байт (8 бит)
func (l *LFSR) NextByte() byte {
	var result byte = 0
	for i := 0; i < 8; i++ {
		result |= l.NextBit() << i
	}
	return result
}

// GenerateKeyStream генерирует поток ключей заданной длины в байтах
func (l *LFSR) GenerateKeyStream(length int) []byte {
	keyStream := make([]byte, length)
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
