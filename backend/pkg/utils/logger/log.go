package logger

import (
	"log"
	"os"
	"sync"
)

const (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
	BoldBlue   = "\033[1;34m"
	BoldPurple = "\033[1;35m"
	BoldCyan   = "\033[1;36m"
	BoldWhite  = "\033[1;37m"
)

type Logger struct {
	*log.Logger
}

var instance *Logger
var once sync.Once

func GetInstance() *Logger {
	once.Do(func() {
		instance = &Logger{
			Logger: log.New(os.Stdout, "", log.LstdFlags),
		}
	})
	return instance
}

func (l *Logger) Info(message string) {
	l.Printf("%s%s[*** INFO ***]%s %s%s%s", Bold, BoldBlue, Reset, BoldBlue, message, Reset)
}

func (l *Logger) Warn(message string) {
	l.Printf("%s%s[*** WARN ***]%s %s%s%s", Bold, BoldYellow, Reset, BoldYellow, message, Reset)
}

func (l *Logger) Error(message string) {
	l.Printf("%s%s[*** ERROR ***]%s %s%s%s", Bold, BoldRed, Reset, BoldRed, message, Reset)
}

func (l *Logger) Debug(message string) {
	l.Printf("%s%s[*** DEBUG ***]%s %s%s%s", Bold, BoldCyan, Reset, BoldCyan, message, Reset)
}
