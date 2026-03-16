package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Logger struct {
	mu      sync.Mutex
	logFile string
}

func NewLogger() *Logger {
	return &Logger{
		logFile: filepath.Join(getAppDataDir(), "app.log"),
	}
}

func (l *Logger) write(level, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("[%s] [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), level, msg)

	f, err := os.OpenFile(l.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(line)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.write("INFO", format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.write("ERROR", format, args...)
}

func (l *Logger) ReadAll() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	data, err := os.ReadFile(l.logFile)
	if err != nil {
		return ""
	}
	return string(data)
}

func (l *Logger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	os.WriteFile(l.logFile, []byte{}, 0644)
}
