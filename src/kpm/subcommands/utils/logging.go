package utils

import (
	"log"
	"os"
)

type _logger struct {
	Verbose *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

func NewLogger() _logger {
	var defaultLogger = _logger{
		log.New(os.Stdout, "[VERB] ", log.LstdFlags),
		log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		log.New(os.Stdout, "[WARN] ", log.LstdFlags),
		log.New(os.Stderr, "[ERR!] ", log.LstdFlags),
	}

	return defaultLogger
}
