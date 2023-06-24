package mylog

import (
	"fmt"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	HIGHLIGHT
)

var levelMap map[string]LogLevel = map[string]LogLevel{
	"debug":     DEBUG,
	"info":      INFO,
	"warn":      WARN,
	"error":     ERROR,
	"highlight": HIGHLIGHT,
}

var levelTag map[LogLevel]string = map[LogLevel]string{
	DEBUG:     MagentaString("[D]"),
	INFO:      CyanString("[I]"),
	WARN:      YellowString("[W]"),
	ERROR:     RedString("[E]"),
	HIGHLIGHT: GreenString("[H]"),
}

var Level LogLevel = DEBUG
var FormatTime = time.TimeOnly

func Warn(v ...interface{}) (int, error) {
	if Level <= WARN {
		v = append([]interface{}{time.Now().Format(FormatTime), levelTag[WARN]}, v...)
		return fmt.Println(v...)
	} else {
		return 0, nil
	}
}

func Info(v ...interface{}) (int, error) {
	if Level <= INFO {
		v = append([]interface{}{time.Now().Format(FormatTime), levelTag[INFO]}, v...)
		return fmt.Println(v...)
	} else {
		return 0, nil
	}
}

func Debug(v ...interface{}) (int, error) {
	if Level <= DEBUG {
		v = append([]interface{}{time.Now().Format(FormatTime), levelTag[DEBUG]}, v...)
		return fmt.Println(v...)
	} else {
		return 0, nil
	}
}

func Error(v ...interface{}) (int, error) {
	if Level <= ERROR {
		v = append([]interface{}{time.Now().Format(FormatTime), levelTag[ERROR]}, v...)
		return fmt.Println(v...)
	} else {
		return 0, nil
	}
}

func Highlight(v ...interface{}) (int, error) {
	if Level <= HIGHLIGHT {
		v = append([]interface{}{time.Now().Format(FormatTime), levelTag[HIGHLIGHT]}, v...)
		return fmt.Println(v...)
	} else {
		return 0, nil
	}
}

func GetLevelBy(s string) LogLevel {
	return levelMap[strings.ToLower(s)]
}
