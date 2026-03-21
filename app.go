package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type OperationResult struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	CipherText    string `json:"cipherText"`
	KeyStream     string `json:"keyStream"`
	BitsCount     int    `json:"bitsCount"`
	ExtractedBits string `json:"extractedBits"`
}

type KeyValidationResult struct {
	Valid     bool   `json:"valid"`
	BinaryKey string `json:"binaryKey"`
	Message   string `json:"message"`
	KeyLength int    `json:"keyLength"`
}

type InputValidationResult struct {
	Valid         bool   `json:"valid"`
	ExtractedBits string `json:"extractedBits"`
	BitsCount     int    `json:"bitsCount"`
	Message       string `json:"message"`
}

type FileReadResult struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Bits     string `json:"bits"`
	FilePath string `json:"filePath"`
	FileSize int    `json:"fileSize"`
}

func extractBits(input string) string {
	var bits strings.Builder
	for _, c := range input {
		if c == '0' || c == '1' {
			bits.WriteRune(c)
		}
	}
	return bits.String()
}

func (a *App) ValidateInput(input string) InputValidationResult {
	if input == "" {
		return InputValidationResult{
			Valid:   false,
			Message: "Поле ввода пустое",
		}
	}

	extracted := extractBits(input)

	if len(extracted) == 0 {
		return InputValidationResult{
			Valid:         false,
			ExtractedBits: "",
			BitsCount:     0,
			Message:       "Ошибка: в тексте нет ни одного символа 0 или 1",
		}
	}

	return InputValidationResult{
		Valid:         true,
		ExtractedBits: extracted,
		BitsCount:     len(extracted),
		Message:       fmt.Sprintf("Извлечено %d бит", len(extracted)),
	}
}

func (a *App) ValidateKey(input string) KeyValidationResult {
	if input == "" {
		return KeyValidationResult{
			Valid:   false,
			Message: "Ключ не может быть пустым",
		}
	}

	for i, c := range input {
		if c != '0' && c != '1' {
			return KeyValidationResult{
				Valid:     false,
				BinaryKey: input,
				KeyLength: len(input),
				Message:   fmt.Sprintf("Ошибка: недопустимый символ '%c' на позиции %d. Разрешены только 0 и 1", c, i+1),
			}
		}
	}

	if len(input) != 37 {
		return KeyValidationResult{
			Valid:     false,
			BinaryKey: input,
			KeyLength: len(input),
			Message:   fmt.Sprintf("Ошибка: длина ключа должна быть ровно 37 бит. Текущая длина: %d", len(input)),
		}
	}

	return KeyValidationResult{
		Valid:     true,
		BinaryKey: input,
		KeyLength: len(input),
		Message:   "Ключ валиден. Длина: 37 бит",
	}
}

func (a *App) SelectAndReadFile() FileReadResult {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Выберите файл для обработки",
		Filters: []runtime.FileFilter{
			{DisplayName: "Все файлы", Pattern: "*.*"},
		},
	})

	if file == "" {
		return FileReadResult{
			Success: false,
			Message: "Файл не выбран",
		}
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return FileReadResult{
			Success: false,
			Message: fmt.Sprintf("Ошибка чтения файла: %v", err),
		}
	}

	if len(data) == 0 {
		return FileReadResult{
			Success: false,
			Message: "Файл пустой",
		}
	}

	bits := bytesToBits(data)

	return FileReadResult{
		Success:  true,
		Message:  fmt.Sprintf("Файл загружен: %d байт = %d бит", len(data), len(bits)),
		Bits:     bits,
		FilePath: file,
		FileSize: len(data),
	}
}

func bytesToBits(data []byte) string {
	var bits strings.Builder
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			if (b>>i)&1 == 1 {
				bits.WriteByte('1')
			} else {
				bits.WriteByte('0')
			}
		}
	}
	return bits.String()
}

func bitsToBytes(bits string) []byte {
	for len(bits)%8 != 0 {
		bits = "0" + bits
	}

	bytes := make([]byte, len(bits)/8)
	for i := 0; i < len(bytes); i++ {
		var b byte = 0
		for j := 0; j < 8; j++ {
			if bits[i*8+j] == '1' {
				b |= 1 << (7 - j)
			}
		}
		bytes[i] = b
	}
	return bytes
}

func (a *App) Encrypt(input string, key string) OperationResult {
	inputValidation := a.ValidateInput(input)
	if !inputValidation.Valid {
		return OperationResult{
			Success: false,
			Message: inputValidation.Message,
		}
	}

	keyValidation := a.ValidateKey(key)
	if !keyValidation.Valid {
		return OperationResult{
			Success: false,
			Message: keyValidation.Message,
		}
	}

	bits := inputValidation.ExtractedBits

	lfsr := NewLFSR(key)
	keyStreamBits := lfsr.GenerateKeyStreamBits(len(bits))

	cipherBits := xorBits(bits, keyStreamBits)

	return OperationResult{
		Success:       true,
		Message:       fmt.Sprintf("Зашифровано %d бит", len(bits)),
		CipherText:    cipherBits,
		KeyStream:     truncateString(keyStreamBits, 1000000000),
		BitsCount:     len(bits),
		ExtractedBits: bits,
	}
}

func (a *App) Decrypt(cipherText string, key string) OperationResult {
	inputValidation := a.ValidateInput(cipherText)
	if !inputValidation.Valid {
		return OperationResult{
			Success: false,
			Message: "Шифротекст: " + inputValidation.Message,
		}
	}

	keyValidation := a.ValidateKey(key)
	if !keyValidation.Valid {
		return OperationResult{
			Success: false,
			Message: keyValidation.Message,
		}
	}

	bits := inputValidation.ExtractedBits

	lfsr := NewLFSR(key)
	keyStreamBits := lfsr.GenerateKeyStreamBits(len(bits))

	plainBits := xorBits(bits, keyStreamBits)

	return OperationResult{
		Success:       true,
		Message:       fmt.Sprintf("Расшифровано %d бит", len(bits)),
		CipherText:    plainBits,
		KeyStream:     truncateString(keyStreamBits, 1000000000),
		BitsCount:     len(bits),
		ExtractedBits: bits,
	}
}

func xorBits(bits1, bits2 string) string {
	result := make([]byte, len(bits1))
	for i := 0; i < len(bits1); i++ {
		if bits1[i] == bits2[i] {
			result[i] = '0'
		} else {
			result[i] = '1'
		}
	}
	return string(result)
}

func (a *App) SaveToFile(bits string, defaultName string) OperationResult {
	inputValidation := a.ValidateInput(bits)
	if !inputValidation.Valid {
		return OperationResult{
			Success: false,
			Message: inputValidation.Message,
		}
	}

	cleanBits := inputValidation.ExtractedBits

	outputPath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Сохранить файл",
		DefaultFilename: defaultName,
	})

	if err != nil || outputPath == "" {
		return OperationResult{
			Success: false,
			Message: "Не выбрано место сохранения",
		}
	}

	data := bitsToBytes(cleanBits)

	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Ошибка сохранения: %v", err),
		}
	}

	return OperationResult{
		Success:   true,
		Message:   fmt.Sprintf("Сохранено: %s (%d байт)", filepath.Base(outputPath), len(data)),
		BitsCount: len(cleanBits),
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
