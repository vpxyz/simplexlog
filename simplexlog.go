// Package simplexlog simple wrapper for the standard log package, that adds log level.
// simplexlog is concurrent safe.
package simplexlog

import (
	"io"
	"log"
	"os"
	"sync"
)

const (
	// DefaultLogFlags default flagsfor log
	DefaultLogFlags = log.Ldate | log.Ltime | log.Lmicroseconds

	// Some usefull predefined label

	// LabelCritical critical label
	LabelCritical = "CRITICAL "
	// LabelError error label
	LabelError = "ERROR "
	// LabelWarning warning label
	LabelWarning = "WARNING "
	// LabelNotice notice label
	LabelNotice = "NOTICE "
	// LabelInfo info label
	LabelInfo = "INFO "
	// LabelDebug debug label
	LabelDebug = "DEBUG "
	// LabelTrace trace label
	LabelTrace = "TRACE "

	// Log levels

	// Critical log level
	Critical LogLevel = iota
	// Error log level
	Error
	// Warning log level
	Warning
	// Notice log level
	Notice
	// Info log level
	Info
	// Debug log level
	Debug
	// Trace log level
	Trace
	// All log level
	All
)

// LogLevel level of log
type LogLevel uint

// Log simple log wrapper
type Logger struct {
	logCritical,
	logError,
	logWarning,
	logNotice,
	logInfo,
	logDebug,
	logTrace *log.Logger
	mutex sync.Mutex // guard the log level
	level LogLevel
}

// Config log option
type Config struct {
	// Out is the output writer
	Out io.Writer
	// Label the prefix of a log line
	Label string
	// Flags are the same combination of flag of standard log package
	Flags int
}

// SetDefault set the options of default logger used by all log level except Error, Critical level and Fatal and Panic
func SetDefault(c Config) func(*Logger) {
	return func(l *Logger) {
		l.logWarning = log.New(c.Out, c.Label, c.Flags)
		l.logInfo = l.logWarning
		l.logNotice = l.logWarning
		l.logDebug = l.logWarning
		l.logTrace = l.logWarning
	}
}

// SetErrorDefault set the options of default logger for error (used by Error, Critical level and by Fatal and Panic)
func SetErrorDefault(c Config) func(*Logger) {
	return func(l *Logger) {
		l.logError = log.New(c.Out, c.Label, c.Flags)
		l.logCritical = l.logError
	}
}

// SetAllDefault set the options of default logger used by all the log level
func SetAllDefault(c Config) func(*Logger) {
	return func(l *Logger) {
		l.logWarning = log.New(c.Out, c.Label, c.Flags)
		l.logInfo = l.logWarning
		l.logNotice = l.logWarning
		l.logDebug = l.logWarning
		l.logTrace = l.logWarning
		l.logError = l.logWarning
		l.logCritical = l.logWarning
	}
}

// SetDebug set the options of debug logger
func SetDebug(c Config) func(*Logger) {
	return func(l *Logger) {
		l.logDebug = log.New(c.Out, c.Label, c.Flags)
	}
}

// SetTrace set the options of trace logger
func SetTrace(c Config) func(*Logger) {
	return func(l *Logger) {
		l.logTrace = log.New(c.Out, c.Label, c.Flags)
	}
}

// SetCritical set the options of critical logger
func SetCritical(c Config) func(*Logger) {
	return func(l *Logger) {
		l.logCritical = log.New(c.Out, c.Label, c.Flags)
	}
}

// SetCritical set the options of error logger
func SetError(c Config) func(*Logger) {
	return func(l *Logger) {
		l.logError = log.New(c.Out, c.Label, c.Flags)
	}
}

// SetWarning set the option of warning logger
func SetWarning(o Config) func(*Logger) {
	return func(l *Logger) {
		l.logWarning = log.New(o.Out, o.Label, o.Flags)
	}
}

// SetNotice set the option of notice logger
func SetNotice(c Config) func(*Logger) {
	return func(l *Logger) {
		l.logNotice = log.New(c.Out, c.Label, c.Flags)
	}
}

// SetInfo set the option of info logger
func SetInfo(c Config) func(*Logger) {
	return func(l *Logger) {
		l.logInfo = log.New(c.Out, c.Label, c.Flags)
	}
}

// New return a new logger. By default, all logs message are output to os.Stdout, except "error" and "critical" message that are logged to os.Stderr.
func New(configurations ...func(*Logger)) *Logger {
	// default logger
	dl := log.New(os.Stdout, "", DefaultLogFlags)
	// default error logger
	el := log.New(os.Stderr, "", DefaultLogFlags)

	logger := Logger{logCritical: el, logError: el, logWarning: dl, logInfo: dl, logDebug: dl, logTrace: dl, level: Info}

	// now customize logger
	for _, config := range configurations {
		config(&logger)
	}

	return &logger
}

// SetLevel change the log level
func (l *Logger) SetLevel(level LogLevel) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.level = level
}

// Level return the current log level
func (l *Logger) Level() LogLevel {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.level
}

// InfoPrintf print to the Info logger
func (l *Logger) InfoPrintf(format string, v ...interface{}) {
	if l.Level() >= Info {
		l.logInfo.Printf(format, v...)
	}
}

// NoticePrintf print to the Notice logger
func (l *Logger) NoticePrintf(format string, v ...interface{}) {
	if l.Level() >= Notice {
		l.logNotice.Printf(format, v...)
	}
}

// WarningPrintf print to the Warning logger
func (l *Logger) WarningPrintf(format string, v ...interface{}) {
	if l.Level() >= Warning {
		l.logWarning.Printf(format, v...)
	}
}

// ErrorPrintf print to the Error logger
func (l *Logger) ErrorPrintf(format string, v ...interface{}) {
	if l.Level() >= Error {
		l.logError.Printf(format, v...)
	}
}

// CriticalPrintf print to the Critical logger
func (l *Logger) CriticalPrintf(format string, v ...interface{}) {
	if l.Level() >= Critical {
		l.logCritical.Printf(format, v...)
	}
}

// DebugPrintf print to the Debug logger
func (l *Logger) DebugPrintf(format string, v ...interface{}) {
	if l.Level() >= Debug {
		l.logDebug.Printf(format, v...)
	}
}

// TracePrintf print to the Debug logger
func (l *Logger) TracePrintf(format string, v ...interface{}) {
	if l.Level() >= Trace {
		l.logTrace.Printf(format, v...)
	}
}

// Fatalf print fatal message to critical logger, followed by call to os.Exit(1)
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logCritical.Fatalf(format, v...)

}

// Panicf print panic message to the critical logger, followed by call to panic()
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.logCritical.Panicf(format, v...)
}

// InfoPrint print to the Info logger
func (l *Logger) InfoPrint(format string) {
	if l.Level() >= Info {
		l.logInfo.Print(format)
	}
}

// NoticePrint print to the Notice logger
func (l *Logger) NoticePrint(format string) {
	if l.Level() >= Notice {
		l.logNotice.Print(format)
	}
}

// WarningPrint print to the Warning logger
func (l *Logger) WarningPrint(format string) {
	if l.Level() >= Warning {
		l.logWarning.Print(format)
	}
}

// ErrorPrint print to the Error logger
func (l *Logger) ErrorPrint(format string) {
	if l.Level() >= Error {
		l.logError.Print(format)
	}
}

// CriticalPrint print to the Critical logger
func (l *Logger) CriticalPrint(format string) {
	if l.Level() >= Critical {
		l.logCritical.Print(format)
	}
}

// DebugPrint print to the Debug logger
func (l *Logger) DebugPrint(format string) {
	if l.Level() >= Debug {
		l.logDebug.Print(format)
	}
}

// TracePrint print to the Debug logger
func (l *Logger) TracePrint(format string) {
	if l.Level() >= Trace {
		l.logTrace.Print(format)
	}
}

// Fatal print fatal message to critical logger, followed by call to os.Exit(1)
func (l *Logger) Fatal(format string) {
	l.logCritical.Fatal(format)

}

// Panic print panic message to the critical logger, followed by call to panic()
func (l *Logger) Panic(format string) {
	l.logCritical.Panic(format)
}

// InfoLogger return the info logger
func (l *Logger) InfoLogger() *log.Logger {
	return l.logInfo
}

// NoticeLogger return the error logger
func (l *Logger) NoticeLogger() *log.Logger {
	return l.logCritical
}

// WarningLogger return the warning logger
func (l *Logger) WarningLogger() *log.Logger {
	return l.logWarning
}

// Errorlogger Return the error logger
func (l *Logger) ErrorLogger() *log.Logger {
	return l.logError
}

// CriticalLogger return the error logger
func (l *Logger) CriticalLogger() *log.Logger {
	return l.logCritical
}

// DebugLogger return the debug logger
func (l *Logger) DebugLogger() *log.Logger {
	return l.logDebug
}

// TraceLogger return the trace logger
func (l *Logger) TraceLogger() *log.Logger {
	return l.logTrace
}
