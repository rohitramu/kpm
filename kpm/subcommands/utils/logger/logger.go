package logger

import (
	"log"
	"os"
)

// TODO: abstract away the "log" library so users of this library won't need to import it

// Default is the default logger
var Default = NewLogger()

// Logger represents a logger
type Logger struct {
	Verbose *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

// NewLogger creates a new logger
func NewLogger() Logger {
	var defaultLogger = Logger{
		log.New(os.Stdout, "[VERB] ", log.LstdFlags),
		log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		log.New(os.Stdout, "[WARN] ", log.LstdFlags),
		log.New(os.Stderr, "[_ERR] ", log.LstdFlags),
	}

	return defaultLogger
}
