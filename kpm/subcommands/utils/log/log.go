package log

import (
	"fmt"
	stdLog "log"
	"os"
	"strings"

	"github.com/emirpasic/gods/maps/hashbidimap"
)

// Level enum
const (
	// LevelNone indicates that logs should not be sent to output.
	LevelNone Level = iota

	// LevelError indicates that only error logs should be sent to output.
	LevelError

	// LevelWarning indicates that only warning and error logs should be sent to output.
	LevelWarning

	// LevelInfo indicates that only info, warning and error logs should be sent to output.
	LevelInfo

	// LevelVerbose indicates that all logs should be sent to output.
	LevelVerbose
)

// loggers is the list of defined loggers.
var loggers = func() []*stdLog.Logger {
	var result = make([]*stdLog.Logger, MaxLevel+1)

	result[LevelNone] = nil
	result[LevelError] = stdLog.New(os.Stderr, "[_ERR] ", stdLog.LstdFlags)
	result[LevelWarning] = stdLog.New(os.Stdout, "[WARN] ", stdLog.LstdFlags)
	result[LevelInfo] = stdLog.New(os.Stdout, "[INFO] ", stdLog.LstdFlags)
	result[LevelVerbose] = stdLog.New(os.Stdout, "[VERB] ", stdLog.LstdFlags)

	return result
}()

var logLevelNames = func() *hashbidimap.Map {
	var m = hashbidimap.New()

	m.Put(LevelNone, "none")
	m.Put(LevelError, "error")
	m.Put(LevelWarning, "warning")
	m.Put(LevelInfo, "info")
	m.Put(LevelVerbose, "verbose")

	return m
}()

// DefaultLevel is the default log level.
const DefaultLevel = LevelInfo

// MaxLevel is the highest value allowed for a log level.
const MaxLevel = LevelVerbose

// MinLevel is the lowest value allowed for a log level.
const MinLevel = LevelNone

// Level identifies the severity of a log line.
type Level int

// currentLogLevel is the currently selected log level.
var currentLogLevel = DefaultLevel

// Parse converts a string representation of a log level to a Level object.
func Parse(logLevelString string) (Level, error) {
	var result Level

	// Get log level value
	var levelObj, found = logLevelNames.GetKey(strings.ToLower(logLevelString))
	if !found {
		return result, fmt.Errorf("Invalid log level string: %s", logLevelString)
	}

	// Cast it to a Level object
	var ok bool
	result, ok = levelObj.(Level)
	if !ok {
		Panic(fmt.Sprintf("Unexpected log level key in log level names map: %s", levelObj))
	}

	return result, nil
}

// String converts a log leve into a string representation.
func (logLevel Level) String() (string, error) {
	var result string

	// Get log level string
	var levelStringObj, found = logLevelNames.Get(logLevel)
	if !found {
		return result, fmt.Errorf("Invalid log level: %s", string(logLevel))
	}

	// Cast it to a string
	var ok bool
	result, ok = levelStringObj.(string)
	if !ok {
		Panic(fmt.Sprintf("Unexpected log level string in log level names map: %s", levelStringObj))
	}

	return result, nil
}

// GetLevel returns the currently selected log level.
func GetLevel() Level {
	return currentLogLevel
}

// SetLevel sets the log level to the provided value.
func SetLevel(level Level) {
	if level > MaxLevel {
		currentLogLevel = MaxLevel
	} else if level < MinLevel {
		currentLogLevel = MinLevel
	} else {
		currentLogLevel = level
	}
}

// Panic logs the formatted message as an error, and then panics.
func Panic(toLog ...interface{}) {
	internalLog(LevelError, toLog...)
}

// Fatal logs the formatted message as an error, and then exits.
func Fatal(toLog ...interface{}) {
	internalLog(LevelError, toLog...)
}

// Error logs the formatted message as an error.
func Error(toLog ...interface{}) {
	internalLog(LevelError, toLog...)
}

// Warning logs the formatted message as a warning.
func Warning(toLog ...interface{}) {
	internalLog(LevelWarning, toLog...)
}

// Info logs the formatted message as an informational message.
func Info(toLog ...interface{}) {
	internalLog(LevelInfo, toLog...)
}

// Verbose logs the formatted message as a verbose message.
func Verbose(toLog ...interface{}) {
	internalLog(LevelVerbose, toLog...)
}

func internalLog(level Level, toLog ...interface{}) {
	if currentLogLevel >= level {
		loggers[level].Println(toLog...)
	}
}
