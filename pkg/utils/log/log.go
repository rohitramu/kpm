package log

import (
	"fmt"
	"io"
	stdLog "log"
	"os"
	"runtime"
	"strings"

	"github.com/vishalkuo/bimap"
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

var logLevelNamesBidiMap = bimap.NewBiMapFromMap[Level, string](map[Level]string{
	LevelNone:    "none",
	LevelError:   "error",
	LevelWarning: "warning",
	LevelInfo:    "info",
	LevelVerbose: "verbose",
	LevelDebug:   "debug",
})

var LevelNames = logLevelNamesBidiMap.GetForwardMap()

// Level identifies the severity of a log line.
type Level int

// Parse converts a string representation of a log level to a Level object.
func Parse(logLevelString string) (Level, error) {
	// Get log level value
	var result, found = logLevelNamesBidiMap.GetInverse(strings.ToLower(logLevelString))
	if !found {
		return result, fmt.Errorf("invalid log level string: %s", logLevelString)
	}

	return result, nil
}

// String converts a log level into a string representation.
func (logLevel Level) String() (string, error) {
	// Get log level string
	var result, found = logLevelNamesBidiMap.Get(logLevel)
	if !found {
		return result, fmt.Errorf("invalid log level: %d", logLevel)
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

// Outputf logs the formatted message as output, without any prefixes or logging flags turned on.
func Outputf(format string, toLog ...interface{}) {
	var output = fmt.Sprintf(format, toLog...)
	var _, err = fmt.Fprintln(WriterOut, output)
	if err != nil {
		Errorf("Failed to write output: %s", output)
	}
}

// OutputStream writes all bytes from the given reader as output.
func OutputStream(reader io.Reader) {
	var _, err = io.Copy(WriterOut, reader)
	if err != nil {
		Errorf("Failed to write stream to output")
	}
}

// Panicf logs the formatted message as an error, and then panics.
func Panicf(format string, toLog ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	checkAndLog(LevelDebug, func(logger *stdLog.Logger) {
		userMessage := fmt.Sprintf(format, toLog...)
		logLocationInfo := fmt.Sprintf("%s:%d", filename, line)
		logger.Printf("[PANIC] %s [%s]", userMessage, logLocationInfo)
	})
	checkAndLog(LevelError, func(logger *stdLog.Logger) {
		logger.Panicf(format, toLog...)
	})
}

// Fatalf logs the formatted message as an error, and then exits.
func Fatalf(format string, toLog ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	checkAndLog(LevelDebug, func(logger *stdLog.Logger) {
		userMessage := fmt.Sprintf(format, toLog...)
		logLocationInfo := fmt.Sprintf("%s:%d", filename, line)
		logger.Printf("[FATAL] %s [%s]", userMessage, logLocationInfo)
	})
	checkAndLog(LevelError, func(logger *stdLog.Logger) {
		logger.Fatalf(format, toLog...)
	})
}

// Errorf logs the formatted message as an error.
func Errorf(format string, toLog ...interface{}) {
	checkAndLog(LevelError, func(logger *stdLog.Logger) {
		logger.Printf(format, toLog...)
	})
}

// Warningf logs the formatted message as a warning.
func Warningf(format string, toLog ...interface{}) {
	checkAndLog(LevelWarning, func(logger *stdLog.Logger) {
		logger.Printf(format, toLog...)
	})
}

// Infof logs the formatted message as an informational message.
func Infof(format string, toLog ...interface{}) {
	checkAndLog(LevelInfo, func(logger *stdLog.Logger) {
		logger.Printf(format, toLog...)
	})
}

// Verbosef logs the formatted message as a verbose message.
func Verbosef(format string, toLog ...interface{}) {
	checkAndLog(LevelVerbose, func(logger *stdLog.Logger) {
		logger.Printf(format, toLog...)
	})
}

// Debugf logs the formatted message as a debug message.
func Debugf(format string, toLog ...interface{}) {
	_, filename, line, _ := runtime.Caller(1)
	checkAndLog(LevelDebug, func(logger *stdLog.Logger) {
		userMessage := fmt.Sprintf(format, toLog...)
		logLocationInfo := fmt.Sprintf("%s:%d", filename, line)
		logger.Printf("%s [%s]", userMessage, logLocationInfo)
	})
}

func checkAndLog(level Level, doLog func(logger *stdLog.Logger)) {
	if currentLogLevel >= level {
		doLog(loggers[level])
	}
}
