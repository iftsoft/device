package core

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
)

// Get go routine ID as a string
func GetGID() string {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	return string(b)
}

// Extract error text from Message object
func GetErrorText(err error) string {
	if err != nil {
		return err.Error()
	}
	return "Success"
}

// Get Call Stack Trace as a string
func TraceCallStack(text string, i int) string {
	//	i := 2
	for {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		text += fmt.Sprintf("\n%s:%d %s", file, line, runtime.FuncForPC(pc).Name())
		i++
	}
	return text
}

func PanicRecover(err *error, log *LogAgent) {
	if r := recover(); r != nil {
		if log != nil {
			log.Panic("Panic happens: %+v", r)
		}
		if err != nil {
			str := fmt.Sprintf("panic is recovered: %+v", r)
			*err = errors.New(str)
		}
	}
}
