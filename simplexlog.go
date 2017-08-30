// Package simplexlog simple wrapper for the standard log package, that adds log level.
// simplexlog is concurrent safe.
package simplexlog

import (
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	// DefaultLogFlags default flagsfor log
	DefaultLogFlags = log.Ldate | log.Ltime | log.Lmicroseconds

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

	// Level Name

	// LevelCritical critical level
	LevelCritical = "CRITICAL"
	// LevelError error level
	LevelError = "ERROR"
	// LevelWarning warning level
	LevelWarning = "WARNING"
	// LevelNotice notice level
	LevelNotice = "NOTICE"
	// LevelInfo info level
	LevelInfo = "INFO"
	// LevelDebug debug level
	LevelDebug = "DEBUG"
	// LevelTrace trace level
	LevelTrace = "TRACE"
	// LevelAll all level
	LevelAll = "ALL"
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
	// default log config
	logger := Logger{
		logCritical: log.New(os.Stderr, LevelCritical+" ", DefaultLogFlags),
		logError:    log.New(os.Stderr, LevelError+" ", DefaultLogFlags),
		logWarning:  log.New(os.Stdout, LevelWarning+" ", DefaultLogFlags),
		logNotice:   log.New(os.Stdout, LevelNotice+" ", DefaultLogFlags),
		logInfo:     log.New(os.Stdout, LevelInfo+" ", DefaultLogFlags),
		logDebug:    log.New(os.Stdout, LevelDebug+" ", DefaultLogFlags),
		logTrace:    log.New(os.Stdout, LevelTrace+" ", DefaultLogFlags),
		level:       Info,
	}

	// now customize logger
	for _, config := range configurations {
		config(&logger)
	}

	return &logger
}

// SwitchTo change the log level
func (l *Logger) SwitchTo(level LogLevel) {
	if level < Critical || level > All {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.level = level
}

// SwitchToLevel change log level, must match (case insensitive) level name (like LevelTrace, LevelCritical etc)
func (l *Logger) SwitchToLevel(level string) {
	level = strings.TrimSpace(strings.ToUpper(level))

	l.mutex.Lock()
	defer l.mutex.Unlock()
	switch level {
	case LevelCritical:
		l.level = Critical
	case LevelError:
		l.level = Error
	case LevelWarning:
		l.level = Warning
	case LevelNotice:
		l.level = Notice
	case LevelInfo:
		l.level = Info
	case LevelDebug:
		l.level = Debug
	case LevelTrace:
		l.level = Trace
	case LevelAll:
		l.level = All
	}
}

// Level return the current log level
func (l *Logger) Level() LogLevel {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.level
}

// Infof print, accordind to format, to the Info logger
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.Level() >= Info {
		l.logInfo.Printf(format, v...)
	}
}

// Noticef print, accordind to format, to the Notice logger
func (l *Logger) Noticef(format string, v ...interface{}) {
	if l.Level() >= Notice {
		l.logNotice.Printf(format, v...)
	}
}

// Warningf print, accordind to format, to the Warning logger
func (l *Logger) Warningf(format string, v ...interface{}) {
	if l.Level() >= Warning {
		l.logWarning.Printf(format, v...)
	}
}

// Errorf print, accordind to format, to the Error logger
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.Level() >= Error {
		l.logError.Printf(format, v...)
	}
}

// Criticalf print, accordind to format, to the Critical logger
func (l *Logger) Criticalf(format string, v ...interface{}) {
	if l.Level() >= Critical {
		l.logCritical.Printf(format, v...)
	}
}

// Debugf print, accordind to format, to the Debug logger
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.Level() >= Debug {
		l.logDebug.Printf(format, v...)
	}
}

// Tracef print, accordind to format, to the Debug logger
func (l *Logger) Tracef(format string, v ...interface{}) {
	if l.Level() >= Trace {
		l.logTrace.Printf(format, v...)
	}
}

// Fatalf print fatal message, accordind to format, to critical logger, followed by call to os.Exit(1)
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logCritical.Fatalf(format, v...)

}

// Panicf print panic message to the critical logger, followed by call to panic()
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.logCritical.Panicf(format, v...)
}

// Info print to the Info logger
func (l *Logger) Info(a ...interface{}) {
	if l.Level() >= Info {
		l.logInfo.Print(a...)
	}
}

// Notice print to the Notice logger
func (l *Logger) Notice(a ...interface{}) {
	if l.Level() >= Notice {
		l.logNotice.Print(a...)
	}
}

// Warning print to the Warning logger
func (l *Logger) Warning(a ...interface{}) {
	if l.Level() >= Warning {
		l.logWarning.Print(a...)
	}
}

// Error print to the Error logger
func (l *Logger) Error(a ...interface{}) {
	if l.Level() >= Error {
		l.logError.Print(a...)
	}
}

// Critical print to the Critical logger
func (l *Logger) Critical(a ...interface{}) {
	if l.Level() >= Critical {
		l.logCritical.Print(a...)
	}
}

// Debug print to the Debug logger
func (l *Logger) Debug(a ...interface{}) {
	if l.Level() >= Debug {
		l.logDebug.Print(a...)
	}
}

// Trace print to the Debug logger
func (l *Logger) Trace(a ...interface{}) {
	if l.Level() >= Trace {
		l.logTrace.Print(a...)
	}
}

// Fatal print fatal message to critical logger, followed by call to os.Exit(1)
func (l *Logger) Fatal(a ...interface{}) {
	l.logCritical.Fatal(a...)

}

// Panic print panic message to the critical logger, followed by call to panic()
func (l *Logger) Panic(a ...interface{}) {
	l.logCritical.Panic(a...)
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
