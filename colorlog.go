package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	DEBUG   = "DEBU"
	INFO    = "INFO"
	TRAC    = "TRAC"
	ERROR   = "ERRO"
	WARNING = "WARN"
	SUCCESS = "SUCC"
	SKIP    = "SKIP"
)

func ColorLog(format string, a ...interface{}) {
	fmt.Print(colorLogS(format, a...))
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
		log = strings.Replace(log, log[i1:i2+1], color.RedString(log[i1+1:i2]), -1)
	}


	// (...) ...
	if strings.Contains(log, "(") && strings.Contains(log, ")") {
		i1 := strings.Index(log, "(")
		i2 := strings.Index(log, ")")
		log = strings.Replace(log, log[i1:i2+1], color.MagentaString(log[i1+1:i2]), -1)
	}

	// <...> ...
	if strings.Contains(log, "<") && strings.Contains(log, ">") {
		i1 := strings.Index(log, "<")
		i2 := strings.Index(log, ">")
		log = strings.Replace(log, log[i1:i2+1], color.CyanString(log[i1+1:i2]), -1)
	}

	// {...} ...
	if strings.Contains(log, "{") && strings.Contains(log, "}") {
		i1 := strings.Index(log, "{")
		i2 := strings.Index(log, "}")
		log = strings.Replace(log, log[i1:i2+1], color.BlueString(log[i1+1:i2]), -1)
	}

	log = clog + log

	return time.Now().Format("2006/01/02 15:04:05 ") + log

}

func getColorLevel(level string) string {
	level = strings.ToUpper(level)
	switch level {
	case DEBUG:
		return color.WhiteString(level)
	case TRAC:
		return color.CyanString(level)
	case INFO:
		return color.BlueString(level)
	case WARNING:
		return color.YellowString(level)
	case ERROR:
		return color.RedString(level)
	case SUCCESS:
		return color.GreenString(level)
	case SKIP:
		return color.MagentaString(level)
	default:
		return level
	}
	return level
}
