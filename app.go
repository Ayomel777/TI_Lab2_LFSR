package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx  context.Context
	lfsr *LFSR
	key  string // двоичное представление ключа (37 символов)
}

// LFSR структура регистра сдвига с линейной обратной связью
type LFSR struct {
	state uint64
}

// константы для многочлена x^37 + x^12 + x^10 + x^2 + 1
const (
	lfsrTaps = (1 << 37) | (1 << 12) | (1 << 10) | (1 << 2) | 1
	lfsrMask = (1 << 37) - 1
)

// NewLFSR создает новый LFSR с начальным состоянием
func NewLFSR(seed uint64) *LFSR {
	return &LFSR{state: seed & lfsrMask}
}

// NextByte возвращает следующий байт гаммы и обновляет состояние
func (l *LFSR) NextByte() byte {
	var out byte
	for i := 0; i < 8; i++ {
		// Вычисляем обратную связь: XOR битов на позициях отводов
		// отводы: 37, 6, 4, 1, 0
		feedback := ((l.state >> 37) & 1) ^
			((l.state >> 6) & 1) ^
			((l.state >> 4) & 1) ^
			((l.state >> 1) & 1) ^
			(l.state & 1)
		// Сдвигаем влево, младший бит = feedback
		l.state = ((l.state << 1) | uint64(feedback)) & lfsrMask
		// Формируем байт: первый бит становится старшим
		out = (out << 1) | byte(feedback)
	}
	return out
}

// Startup вызывается при старте приложения
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// SetKey устанавливает ключ из двоичной строки (ровно 37 символов '0'/'1')
func (a *App) SetKey(key string) error {
	if len(key) != 37 {
		return fmt.Errorf("key must be exactly 37 bits, got %d", len(key))
	}
	var seed uint64
	for i, ch := range key {
		if ch != '0' && ch != '1' {
			return fmt.Errorf("key must contain only 0 and 1, invalid character at position %d", i)
		}
		seed = (seed << 1) | uint64(ch-'0')
	}
	a.lfsr = NewLFSR(seed)
	a.key = key
	return nil
}

// GetKeyInfo возвращает текущий ключ (для отображения)
func (a *App) GetKeyInfo() string {
	return a.key
}

// Encrypt шифрует файл (input → output) с использованием текущего ключа
func (a *App) Encrypt(inputPath, outputPath string) error {
	return a.processFile(inputPath, outputPath)
}

// Decrypt дешифрует (аналогично Encrypt, так как XOR симметричен)
func (a *App) Decrypt(inputPath, outputPath string) error {
	return a.processFile(inputPath, outputPath)
}

// processFile содержит общую логику шифрования/дешифрования
func (a *App) processFile(inPath, outPath string) error {
	if a.lfsr == nil {
		return errors.New("key not set")
	}

	inFile, err := os.Open(inPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer outFile.Close()

	const bufSize = 4096
	buf := make([]byte, bufSize)

	for {
		n, err := inFile.Read(buf)
		if n > 0 {
			// Генерируем гамму для прочитанных байт
			gamma := make([]byte, n)
			for i := 0; i < n; i++ {
				gamma[i] = a.lfsr.NextByte()
			}
			// XOR
			for i := 0; i < n; i++ {
				buf[i] ^= gamma[i]
			}
			if _, err := outFile.Write(buf[:n]); err != nil {
				return fmt.Errorf("error writing output file: %w", err)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading input file: %w", err)
		}
	}
	return nil
}

// OpenFileDialog открывает диалог выбора файла и возвращает выбранный путь
func (a *App) OpenFileDialog() (string, error) {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Выберите файл для шифрования/дешифрования",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Все файлы",
				Pattern:     "*",
			},
		},
	})
	if err != nil {
		return "", err
	}
	return file, nil
}

// SaveFileDialog открывает диалог сохранения файла и возвращает выбранный путь
func (a *App) SaveFileDialog() (string, error) {
	file, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "Сохранить результат как",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Все файлы",
				Pattern:     "*",
			},
		},
	})
	if err != nil {
		return "", err
	}
	return file, nil
}
