package log

import (
	"fmt"
	"io"
	"log"
	stdLog "log"
	"os"
	"runtime"
	"strings"

	"github.com/emirpasic/gods/maps/hashbidimap"
)

// Level enum
const (
	// LevelNone indicates that logs should not be sent to output.
	LevelNone Level = iota

	// LevelError indicates that error logs should be sent to output.
	LevelError

	// LevelWarning indicates that warning and error logs should be sent to output.
	LevelWarning

	// LevelInfo indicates that info, warning and error logs should be sent to output.
	LevelInfo

	// LevelVerbose indicates that verbose, info, warning and error logs should be sent to output.
	LevelVerbose

	// LevelDebug indicates that debug, verbose, info, warning and error logs should be sent to output.
	LevelDebug

	//TODO: Replace most verbose logs with debug logs
)

// WriterOut is the stream to use when writing program output.
var WriterOut = os.Stdout

// WriterInfo is the stream to use when writing info or error logs.
var WriterInfo = os.Stdout

// WriterErr is the stream to use when writing error or warning logs.
var WriterErr = os.Stderr

// DefaultLevel is the default log level.
const DefaultLevel = LevelInfo

// MinLevel is the lowest value allowed for a log level.
const MinLevel = LevelNone

// MaxLevel is the highest value allowed for a log level.
const MaxLevel = LevelDebug

// loggers is the list of defined loggers.
var loggers = func() []*stdLog.Logger {
	var result = make([]*stdLog.Logger, MaxLevel+1)

	result[LevelNone] = nil
	result[LevelError] = stdLog.New(WriterErr, "[ERR] ", stdLog.LstdFlags)
	result[LevelWarning] = stdLog.New(WriterErr, "[WRN] ", stdLog.LstdFlags)
	result[LevelInfo] = stdLog.New(WriterInfo, "[INF] ", stdLog.LstdFlags)
	result[LevelVerbose] = stdLog.New(WriterInfo, "[VRB] ", stdLog.LstdFlags)
	result[LevelDebug] = stdLog.New(WriterInfo, "[DBG] ", stdLog.LstdFlags)

	return result
}()

var logLevelNames = func() *hashbidimap.Map {
	var m = hashbidimap.New()

	m.Put(LevelNone, "none")
	m.Put(LevelError, "error")
	m.Put(LevelWarning, "warning")
	m.Put(LevelInfo, "info")
	m.Put(LevelVerbose, "verbose")
	m.Put(LevelDebug, "debug")

	return m
}()

// Level identifies the severity of a log line.
type Level int

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
		Panic("Unexpected log level key in log level names map: %s", levelObj)
	}

	return result, nil
}

// String converts a log level into a string representation.
func (logLevel Level) String() (string, error) {
	var result string

	// Get log level string
	var levelStringObj, found = logLevelNames.Get(logLevel)
	if !found {
		return result, fmt.Errorf("Invalid log level: %d", logLevel)
	}

	// Cast it to a string
	var ok bool
	result, ok = levelStringObj.(string)
	if !ok {
		Panic("Unexpected log level string in log level names map: %s", levelStringObj)
	}

	return result, nil
}

// currentLogLevel is the currently selected log level.
var currentLogLevel = DefaultLevel

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

// Output logs the formatted message as output, without any prefixes or logging flags turned on.
func Output(format string, toLog ...interface{}) {
	var output = fmt.Sprintf(format, toLog...)
	var _, err = fmt.Fprintln(WriterOut, output)
	if err != nil {
		Error("Failed to write output: %s", output)
	}
}

// OutputStream writes all bytes from the given reader as output.
func OutputStream(reader io.Reader) {
	var _, err = io.Copy(WriterOut, reader)
	if err != nil {
		Error("Failed to write stream to output")
	}
}

// Panic logs the formatted message as an error, and then panics.
func Panic(format string, toLog ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	checkAndLog(LevelDebug, func(logger *log.Logger) {
		userMessage := fmt.Sprintf(format, toLog...)
		logLocationInfo := fmt.Sprintf("%s:%d", filename, line)
		logger.Println(fmt.Sprintf("[PANIC] %s [%s]", userMessage, logLocationInfo))
	})
	checkAndLog(LevelError, func(logger *log.Logger) {
		logger.Panicln(fmt.Sprintf(format, toLog...))
	})
}

// Fatal logs the formatted message as an error, and then exits.
func Fatal(format string, toLog ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	checkAndLog(LevelDebug, func(logger *log.Logger) {
		userMessage := fmt.Sprintf(format, toLog...)
		logLocationInfo := fmt.Sprintf("%s:%d", filename, line)
		logger.Println(fmt.Sprintf("[FATAL] %s [%s]", userMessage, logLocationInfo))
	})
	checkAndLog(LevelError, func(logger *log.Logger) {
		logger.Fatalln(fmt.Sprintf(format, toLog...))
	})
}

// Error logs the formatted message as an error.
func Error(format string, toLog ...interface{}) {
	checkAndLog(LevelError, func(logger *log.Logger) {
		logger.Println(fmt.Sprintf(format, toLog...))
	})
}

// Warning logs the formatted message as a warning.
func Warning(format string, toLog ...interface{}) {
	checkAndLog(LevelWarning, func(logger *log.Logger) {
		logger.Println(fmt.Sprintf(format, toLog...))
	})
}

// Info logs the formatted message as an informational message.
func Info(format string, toLog ...interface{}) {
	checkAndLog(LevelInfo, func(logger *log.Logger) {
		logger.Println(fmt.Sprintf(format, toLog...))
	})
}

// Verbose logs the formatted message as a verbose message.
func Verbose(format string, toLog ...interface{}) {
	checkAndLog(LevelVerbose, func(logger *log.Logger) {
		logger.Println(fmt.Sprintf(format, toLog...))
	})
}

// Debug logs the formatted message as a debug message.
func Debug(format string, toLog ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	checkAndLog(LevelDebug, func(logger *log.Logger) {
		userMessage := fmt.Sprintf(format, toLog...)
		logLocationInfo := fmt.Sprintf("%s:%d", filename, line)
		logger.Println(fmt.Sprintf("%s [%s]", userMessage, logLocationInfo))
	})
}

func checkAndLog(level Level, doLog func(logger *log.Logger)) {
	if currentLogLevel >= level {
		doLog(loggers[level])
	}
}
