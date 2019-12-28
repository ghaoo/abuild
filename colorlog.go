package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	DEBUG   = "DEBU"
	TRAC    = "TRAC"
	INFO    = "INFO"
	WARNING = "WARN"
	ERROR   = "ERRO"
	SUCCESS = "SUCC"
	SKIP    = "SKIP"
)

func ColorLog(format string, a ...interface{}) {
	color.White(colorLogS(format, a...))
}

func colorLogS(format string, a ...interface{}) string {
	log := fmt.Sprintf(format, a...)

	var clog string

	i := strings.Index(log, "]")
	if log[0] == '[' && i > -1 {
		clog += "[" + getColorLevel(log[1:i]) + "]"
	}

	log = log[i+1:]

	// 将限定字符内的文字渲染成特定颜色
	// [...] ...
	if strings.Contains(log, "[") && strings.Contains(log, "]") {
		i1 := strings.Index(log, "[")
		i2 := strings.Index(log, "]")
		log = strings.Replace(log, log[i1:i2+1], color.Set(color.FgRed).Sprintf(log[i1+1:i2]), -1)
	}

	// (...) ...
	if strings.Contains(log, "(") && strings.Contains(log, ")") {
		i1 := strings.Index(log, "(")
		i2 := strings.Index(log, ")")
		log = strings.Replace(log, log[i1:i2+1], color.Set(color.FgGreen).Sprintf(log[i1+1:i2]), -1)
	}

	// <...> ...
	if strings.Contains(log, "<") && strings.Contains(log, ">") {
		i1 := strings.Index(log, "<")
		i2 := strings.Index(log, ">")
		log = strings.Replace(log, log[i1:i2+1], color.Set(color.FgCyan).Sprintf(log[i1+1:i2]), -1)
	}

	// {...} ...
	if strings.Contains(log, "{") && strings.Contains(log, "}") {
		i1 := strings.Index(log, "{")
		i2 := strings.Index(log, "}")
		log = strings.Replace(log, log[i1:i2+1], color.Set(color.FgBlue).Sprintf(log[i1+1:i2]), -1)
	}

	log = clog + log

	return time.Now().Format("2006/01/02 15:04:05 ") + log

}

func getColorLevel(level string) string {
	level = strings.ToUpper(level)
	switch level {
	case DEBUG:
		return color.Set(color.FgWhite).Sprintf(level)
	case TRAC:
		return color.Set(color.FgCyan).Sprintf(level)
	case INFO:
		return color.Set(color.FgBlue).Sprintf(level)
	case WARNING:
		return color.Set(color.FgYellow).Sprintf(level)
	case ERROR:
		return color.Set(color.FgRed).Sprintf(level)
	case SUCCESS:
		return color.Set(color.FgGreen).Sprintf(level)
	case SKIP:
		return color.Set(color.FgMagenta).Sprintf(level)
	default:
		return level
	}
	return level
}
