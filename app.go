package main

import (
	"context"
	"encoding/base64"
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

// Результат операции
type OperationResult struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	BinaryKey     string `json:"binaryKey"`
	KeyStream     string `json:"keyStream"`
	OriginalSize  int    `json:"originalSize"`
	ProcessedSize int    `json:"processedSize"`
	OutputPath    string `json:"outputPath"`
}

// Результат валидации ключа
type KeyValidationResult struct {
	Valid     bool   `json:"valid"`
	BinaryKey string `json:"binaryKey"`
	Message   string `json:"message"`
	KeyLength int    `json:"keyLength"`
}

// Конвертация любого ключа в бинарный
func (a *App) ConvertToBinaryKey(input string) KeyValidationResult {
	if input == "" {
		return KeyValidationResult{
			Valid:   false,
			Message: "Ключ не может быть пустым",
		}
	}

	binaryKey := convertInputToBinary(input)

	return KeyValidationResult{
		Valid:     true,
		BinaryKey: binaryKey,
		KeyLength: len(binaryKey),
		Message:   fmt.Sprintf("Ключ успешно преобразован. Длина: %d бит", len(binaryKey)),
	}
}

// Конвертация входных данных в бинарную строку
func convertInputToBinary(input string) string {
	// Проверяем, является ли ввод уже бинарным
	isBinary := true
	for _, c := range input {
		if c != '0' && c != '1' {
			isBinary = false
			break
		}
	}

	if isBinary && len(input) > 0 {
		return input
	}

	// Конвертируем каждый символ в бинарное представление
	var binary strings.Builder
	for _, c := range input {
		binary.WriteString(fmt.Sprintf("%08b", c))
	}

	return binary.String()
}

// Выбор файла для обработки
func (a *App) SelectInputFile() (string, error) {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Выберите файл для обработки",
		Filters: []runtime.FileFilter{
			{DisplayName: "Все файлы", Pattern: "*.*"},
			{DisplayName: "Изображения", Pattern: "*.png;*.jpg;*.jpeg;*.gif;*.bmp"},
			{DisplayName: "Видео", Pattern: "*.mp4;*.avi;*.mkv;*.mov"},
			{DisplayName: "Аудио", Pattern: "*.mp3;*.wav;*.flac;*.ogg"},
			{DisplayName: "Документы", Pattern: "*.txt;*.pdf;*.doc;*.docx"},
		},
	})
	return file, err
}

// Выбор места сохранения
func (a *App) SelectOutputFile(defaultName string) (string, error) {
	file, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Сохранить зашифрованный файл",
		DefaultFilename: defaultName,
	})
	return file, err
}

// Шифрование файла
func (a *App) EncryptFile(inputPath, key string) OperationResult {
	if inputPath == "" {
		return OperationResult{
			Success: false,
			Message: "Не выбран входной файл",
		}
	}

	if key == "" {
		return OperationResult{
			Success: false,
			Message: "Ключ не может быть пустым",
		}
	}

	// Читаем входной файл
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Ошибка чтения файла: %v", err),
		}
	}

	// Конвертируем ключ в бинарный
	binaryKey := convertInputToBinary(key)

	// Проверяем минимальную длину ключа (37 бит для LFSR)
	if len(binaryKey) < 37 {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Длина ключа должна быть минимум 37 бит. Текущая длина: %d бит", len(binaryKey)),
		}
	}

	// Создаем LFSR и шифруем
	lfsr := NewLFSR(binaryKey)
	keyStream := lfsr.GenerateKeyStream(len(data))
	encrypted := xorBytes(data, keyStream)

	// Генерируем имя выходного файла
	ext := filepath.Ext(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), ext)
	defaultOutput := baseName + "_encrypted" + ext

	// Выбираем место сохранения
	outputPath, err := a.SelectOutputFile(defaultOutput)
	if err != nil || outputPath == "" {
		return OperationResult{
			Success: false,
			Message: "Не выбрано место сохранения",
		}
	}

	// Сохраняем зашифрованный файл
	err = os.WriteFile(outputPath, encrypted, 0644)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Ошибка сохранения файла: %v", err),
		}
	}

	// Формируем отображение keystream (первые 256 бит)
	keyStreamDisplay := formatKeyStream(keyStream, 256)

	return OperationResult{
		Success:       true,
		Message:       "Файл успешно зашифрован",
		BinaryKey:     truncateString(binaryKey, 512),
		KeyStream:     keyStreamDisplay,
		OriginalSize:  len(data),
		ProcessedSize: len(encrypted),
		OutputPath:    outputPath,
	}
}

// Дешифрование файла
func (a *App) DecryptFile(inputPath, key string) OperationResult {
	if inputPath == "" {
		return OperationResult{
			Success: false,
			Message: "Не выбран входной файл",
		}
	}

	if key == "" {
		return OperationResult{
			Success: false,
			Message: "Ключ не может быть пустым",
		}
	}

	// Читаем входной файл
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Ошибка чтения файла: %v", err),
		}
	}

	// Конвертируем ключ в бинарный
	binaryKey := convertInputToBinary(key)

	// Проверяем минимальную длину ключа
	if len(binaryKey) < 37 {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Длина ключа должна быть минимум 37 бит. Текущая длина: %d бит", len(binaryKey)),
		}
	}

	// Создаем LFSR и дешифруем (XOR симметричен)
	lfsr := NewLFSR(binaryKey)
	keyStream := lfsr.GenerateKeyStream(len(data))
	decrypted := xorBytes(data, keyStream)

	// Генерируем имя выходного файла
	ext := filepath.Ext(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), ext)
	baseName = strings.TrimSuffix(baseName, "_encrypted")
	defaultOutput := baseName + "_decrypted" + ext

	// Выбираем место сохранения
	outputPath, err := a.SelectOutputFile(defaultOutput)
	if err != nil || outputPath == "" {
		return OperationResult{
			Success: false,
			Message: "Не выбрано место сохранения",
		}
	}

	// Сохраняем расшифрованный файл
	err = os.WriteFile(outputPath, decrypted, 0644)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Ошибка сохранения файла: %v", err),
		}
	}

	// Формируем отображение keystream
	keyStreamDisplay := formatKeyStream(keyStream, 256)

	return OperationResult{
		Success:       true,
		Message:       "Файл успешно расшифрован",
		BinaryKey:     truncateString(binaryKey, 512),
		KeyStream:     keyStreamDisplay,
		OriginalSize:  len(data),
		ProcessedSize: len(decrypted),
		OutputPath:    outputPath,
	}
}

// Шифрование текста
func (a *App) EncryptText(text, key string) OperationResult {
	if text == "" {
		return OperationResult{
			Success: false,
			Message: "Текст не может быть пустым",
		}
	}

	if key == "" {
		return OperationResult{
			Success: false,
			Message: "Ключ не может быть пустым",
		}
	}

	binaryKey := convertInputToBinary(key)

	if len(binaryKey) < 37 {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Длина ключа должна быть минимум 37 бит. Текущая длина: %d бит", len(binaryKey)),
		}
	}

	data := []byte(text)
	lfsr := NewLFSR(binaryKey)
	keyStream := lfsr.GenerateKeyStream(len(data))
	encrypted := xorBytes(data, keyStream)

	// Кодируем в base64 для безопасного отображения
	encodedResult := base64.StdEncoding.EncodeToString(encrypted)
	keyStreamDisplay := formatKeyStream(keyStream, 256)

	return OperationResult{
		Success:       true,
		Message:       encodedResult,
		BinaryKey:     truncateString(binaryKey, 512),
		KeyStream:     keyStreamDisplay,
		OriginalSize:  len(data),
		ProcessedSize: len(encrypted),
	}
}

// Дешифрование текста
func (a *App) DecryptText(encodedText, key string) OperationResult {
	if encodedText == "" {
		return OperationResult{
			Success: false,
			Message: "Текст не может быть пустым",
		}
	}

	if key == "" {
		return OperationResult{
			Success: false,
			Message: "Ключ не может быть пустым",
		}
	}

	// Декодируем из base64
	data, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: "Ошибка декодирования: неверный формат зашифрованного текста",
		}
	}

	binaryKey := convertInputToBinary(key)

	if len(binaryKey) < 37 {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Длина ключа должна быть минимум 37 бит. Текущая длина: %d бит", len(binaryKey)),
		}
	}

	lfsr := NewLFSR(binaryKey)
	keyStream := lfsr.GenerateKeyStream(len(data))
	decrypted := xorBytes(data, keyStream)

	keyStreamDisplay := formatKeyStream(keyStream, 256)

	return OperationResult{
		Success:       true,
		Message:       string(decrypted),
		BinaryKey:     truncateString(binaryKey, 512),
		KeyStream:     keyStreamDisplay,
		OriginalSize:  len(data),
		ProcessedSize: len(decrypted),
	}
}

// XOR двух массивов байт
func xorBytes(data, keyStream []byte) []byte {
	result := make([]byte, len(data))
	for i := range data {
		result[i] = data[i] ^ keyStream[i]
	}
	return result
}

// Форматирование keystream для отображения
func formatKeyStream(keyStream []byte, maxBits int) string {
	var sb strings.Builder
	bitCount := 0

	for _, b := range keyStream {
		for i := 7; i >= 0 && bitCount < maxBits; i-- {
			if (b>>i)&1 == 1 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
			bitCount++
			if bitCount%8 == 0 && bitCount < maxBits {
				sb.WriteByte(' ')
			}
		}
		if bitCount >= maxBits {
			break
		}
	}

	if len(keyStream)*8 > maxBits {
		sb.WriteString("...")
	}

	return sb.String()
}

// Обрезка строки
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
