package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	Green   = "\033[97;42m"
	White   = "\033[90;47m"
	Yellow  = "\033[97;43m"
	Red     = "\033[97;41m"
	Blue    = "\033[97;44m"
	Magenta = "\033[97;45m"
	Cyan    = "\033[97;46m"
	Reset   = "\033[0m"
)

func LogErr(format string, a ...any) {
	logTypeImpl("error", format, a...)
}

func LogWarn(format string, a ...any) {
	logTypeImpl("warn", format, a...)
}

func LogInfo(format string, a ...any) {
	logTypeImpl("info", format, a...)
}

func LogSuccess(format string, a ...any) {
	logTypeImpl("success", format, a...)
}

func getNowTimeString() string {
	return time.Now().String()[0:19]
}

func logTypeImpl(typeStr string, format string, a ...any) {
	color := Blue
	tip := "INFO"
	switch typeStr {
	case "error":
		color = Red
		tip = "ERROR"
	case "warn":
		color = Yellow
		tip = "WARN"
	case "success":
		color = Green
		tip = "DONE"
	}
	logImpl(Reset+"["+getNowTimeString()+"] |"+color+tip+Reset+"|"+": "+format+"\n"+Reset, a...)
}

func logImpl(format string, a ...any) {
	_, err := fmt.Fprintf(os.Stdout, format, a...)
	if err != nil {
		return
	}
}
