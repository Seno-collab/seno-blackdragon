package utils

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type DailyFileWriter struct {
	basePath string
	curDate  string
	file     *os.File
	mu       sync.Mutex
}

func NewDailyFileWriter(basePath string) *DailyFileWriter {
	return &DailyFileWriter{basePath: basePath}
}

func (w *DailyFileWriter) rotateIfNeededLocked() error {
	today := time.Now().Format("2006-01-02")
	if w.file != nil && today == w.curDate {
		return nil
	}
	if w.file != nil {
		_ = w.file.Sync()
		_ = w.file.Close()
	}
	filePath := fmt.Sprintf("%s-%s.log", w.basePath, today)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.file = f
	w.curDate = today
	return nil
}

func (w *DailyFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.rotateIfNeededLocked(); err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

func (w *DailyFileWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

func (w *DailyFileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		if err := w.file.Sync(); err != nil {
			_ = w.file.Close()
			w.file = nil
			return err
		}
		err := w.file.Close()
		w.file = nil
		return err
	}
	return nil
}
