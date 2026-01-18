package log

import (
	"log"
)

type Logger struct {
	level string
}

func New(level string) *Logger {
	return &Logger{level: level}
}

func (l *Logger) Info(msg string, kv ...any) {
	log.Println(" [INFO]: ", msg, kv)
}

func (l *Logger) Error(msg string, kv ...any) {
	log.Println("[ERROR]: ", msg, kv)
}

func (l *Logger) Fatal(msg string, kv ...any) {
	log.Fatal("[FATAL]: ", msg, kv)
}
