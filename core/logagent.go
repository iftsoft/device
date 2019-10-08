package core

import (
	"fmt"
	"runtime"
	"time"
)

type EnumLogLevel uint16

// log level
const (
	LogLevelEmpty EnumLogLevel = iota
	LogLevelPanic
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
	LogLevelDump
	LogLevelTrace
//	LogLevelMax
)

func (e EnumLogLevel) String() string {
	switch e {
	case LogLevelEmpty:			return "EMPTY"
	case LogLevelPanic:			return "PANIC"
	case LogLevelError:			return "ERROR"
	case LogLevelWarn:			return "WARN"
	case LogLevelInfo:			return "INFO"
	case LogLevelDebug:			return "DEBUG"
	case LogLevelDump:			return "DUMP"
	case LogLevelTrace:			return "TRACE"
	default:					return "UNDEF"
	}
}

//var logLevelNames = [LogLevelMax]string{
//	"EMPTY", "PANIC", "ERROR", "WARN ", "INFO ", "DEBUG", "DUMP ", "TRACE",
//}

//func GetLogLevelText(level EnumLogLevel) string {
//	if level >= 0 && level < LogLevelMax {
//		return logLevelNames[level]
//	}
//	return "UNDEF"
//}

func GetLogAgent(level EnumLogLevel, title string) *LogAgent {
	if level > LogLevelTrace {
		level = LogLevelEmpty
	}
	return &LogAgent{level, title, getLogSource()}
}

type LogAgent struct {
	logLevel EnumLogLevel
	modTitle string
	source   bool
}

func (log *LogAgent) Init(level EnumLogLevel, title string, src bool) {
	log.logLevel = level
	log.modTitle = title
	log.source = src
}

func (log *LogAgent) SetLevel(level EnumLogLevel) {
	log.logLevel = level
}

func (log LogAgent) IsTrace() bool {
	return log.logLevel >= LogLevelTrace
}
func (log LogAgent) IsDump() bool {
	return log.logLevel >= LogLevelDump
}
func (log LogAgent) IsDebug() bool {
	return log.logLevel >= LogLevelDebug
}
func (log LogAgent) IsInfo() bool {
	return log.logLevel >= LogLevelInfo
}
func (log LogAgent) IsWarn() bool {
	return log.logLevel >= LogLevelWarn
}
func (log LogAgent) IsError() bool {
	return log.logLevel >= LogLevelError
}
func (log LogAgent) IsPanic() bool {
	return log.logLevel >= LogLevelPanic
}
func (log LogAgent) IsEmpty() bool {
	return log.logLevel == LogLevelEmpty
}

func (log *LogAgent) Trace(format string, args ...interface{}) {
	if log != nil && log.IsTrace() {
		log.formatLine(LogLevelTrace, fmt.Sprintf(format, args...))
	}
}
func (log *LogAgent) Dump(format string, args ...interface{}) {
	if log != nil && log.IsDump() {
		log.formatLine(LogLevelDump, fmt.Sprintf(format, args...))
	}
}
func (log *LogAgent) Debug(format string, args ...interface{}) {
	if log != nil && log.IsDebug() {
		log.formatLine(LogLevelDebug, fmt.Sprintf(format, args...))
	}
}
func (log *LogAgent) Info(format string, args ...interface{}) {
	if log != nil && log.IsInfo() {
		log.formatLine(LogLevelInfo, fmt.Sprintf(format, args...))
	}
}
func (log *LogAgent) Warn(format string, args ...interface{}) {
	if log != nil && log.IsWarn() {
		log.formatLine(LogLevelWarn, fmt.Sprintf(format, args...))
	}
}
func (log *LogAgent) Error(format string, args ...interface{}) {
	if log != nil && log.IsError() {
		log.formatLine(LogLevelError, fmt.Sprintf(format, args...))
	}
}
func (log *LogAgent) Panic(format string, args ...interface{}) {
	if log != nil && log.IsPanic() {
		text := fmt.Sprintf(format, args...)
		log.formatLine(LogLevelPanic, TraceCallStack(text, 2))
	}
}

func (log *LogAgent) formatLine(level EnumLogLevel, text string) {
	t := time.Now()
	moment := t.Format("2006-01-02 15:04:05.999999")
	for size := len(moment); size < 26; size++ {
		moment += "0"
	}
	gid := GetGID()
	var mesg string
	if log.source {
		pc, file, line, ok := runtime.Caller(2)
		if ok {
			mesg = fmt.Sprintf("%s [%s %s %s] %s:%d %s() %s\n",
				moment, level.String(), gid, log.modTitle, file, line, runtime.FuncForPC(pc).Name(), text)
		} else {
			mesg = fmt.Sprintf("%s [%s %s %s] %s\n", moment, level.String(), gid, log.modTitle, text)
		}
	} else {
		mesg = fmt.Sprintf("%s [%s %s %s] %s\n", moment, level.String(), gid, log.modTitle, text)
	}
	LogToFile(level, mesg)
}

func (log *LogAgent) PanicRecover() {
	if r := recover(); r != nil {
		if log != nil {
			log.Warn("Panic Recovered: %+v", r)
		}
	}
}
