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
	BinaryKey     string `json:"binaryKey"`
	KeyStream     string `json:"keyStream"`
	OriginalSize  int    `json:"originalSize"`
	ProcessedSize int    `json:"processedSize"`
	OutputPath    string `json:"outputPath"`
}

type KeyValidationResult struct {
	Valid     bool   `json:"valid"`
	BinaryKey string `json:"binaryKey"`
	Message   string `json:"message"`
	KeyLength int    `json:"keyLength"`
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

func (a *App) SelectInputFile() (string, error) {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Выберите файл для обработки",
		Filters: []runtime.FileFilter{
			{DisplayName: "Все файлы", Pattern: "*.*"},
			{DisplayName: "Изображения", Pattern: "*.png;*.jpg;*.jpeg;*.gif;*.bmp;*.webp"},
			{DisplayName: "Видео", Pattern: "*.mp4;*.avi;*.mkv;*.mov;*.wmv"},
			{DisplayName: "Аудио", Pattern: "*.mp3;*.wav;*.flac;*.ogg;*.aac"},
			{DisplayName: "Документы", Pattern: "*.txt;*.pdf;*.doc;*.docx;*.xls;*.xlsx"},
			{DisplayName: "Архивы", Pattern: "*.zip;*.rar;*.7z;*.tar;*.gz"},
		},
	})
	return file, err
}

func (a *App) SelectOutputFile(defaultName string) (string, error) {
	file, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Сохранить файл",
		DefaultFilename: defaultName,
	})
	return file, err
}

func (a *App) EncryptFile(inputPath, key string) OperationResult {
	if inputPath == "" {
		return OperationResult{
			Success: false,
			Message: "Не выбран входной файл",
		}
	}

	keyValidation := a.ValidateKey(key)
	if !keyValidation.Valid {
		return OperationResult{
			Success: false,
			Message: keyValidation.Message,
		}
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Ошибка чтения файла: %v", err),
		}
	}

	if len(data) == 0 {
		return OperationResult{
			Success: false,
			Message: "Файл пустой",
		}
	}

	lfsr := NewLFSR(key)
	keyStream := lfsr.GenerateKeyStream(len(data))
	encrypted := xorBytes(data, keyStream)

	ext := filepath.Ext(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), ext)
	defaultOutput := baseName + "_encrypted" + ext

	outputPath, err := a.SelectOutputFile(defaultOutput)
	if err != nil || outputPath == "" {
		return OperationResult{
			Success: false,
			Message: "Не выбрано место сохранения",
		}
	}

	err = os.WriteFile(outputPath, encrypted, 0644)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Ошибка сохранения файла: %v", err),
		}
	}

	keyStreamDisplay := formatKeyStream(keyStream, 256)

	return OperationResult{
		Success:       true,
		Message:       "Файл успешно зашифрован",
		BinaryKey:     key,
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

	// Валидация ключа
	keyValidation := a.ValidateKey(key)
	if !keyValidation.Valid {
		return OperationResult{
			Success: false,
			Message: keyValidation.Message,
		}
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Ошибка чтения файла: %v", err),
		}
	}

	if len(data) == 0 {
		return OperationResult{
			Success: false,
			Message: "Файл пустой",
		}
	}

	lfsr := NewLFSR(key)
	keyStream := lfsr.GenerateKeyStream(len(data))
	decrypted := xorBytes(data, keyStream)

	ext := filepath.Ext(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), ext)
	baseName = strings.TrimSuffix(baseName, "_encrypted")
	defaultOutput := baseName + "_decrypted" + ext

	outputPath, err := a.SelectOutputFile(defaultOutput)
	if err != nil || outputPath == "" {
		return OperationResult{
			Success: false,
			Message: "Не выбрано место сохранения",
		}
	}

	err = os.WriteFile(outputPath, decrypted, 0644)
	if err != nil {
		return OperationResult{
			Success: false,
			Message: fmt.Sprintf("Ошибка сохранения файла: %v", err),
		}
	}

	keyStreamDisplay := formatKeyStream(keyStream, 256)

	return OperationResult{
		Success:       true,
		Message:       "Файл успешно расшифрован",
		BinaryKey:     key,
		KeyStream:     keyStreamDisplay,
		OriginalSize:  len(data),
		ProcessedSize: len(decrypted),
		OutputPath:    outputPath,
	}
}

func xorBytes(data, keyStream []byte) []byte {
	result := make([]byte, len(data))
	for i := range data {
		result[i] = data[i] ^ keyStream[i]
	}
	return result
}

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
