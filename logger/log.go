package logger

import (
	"fmt"
	"os"
	"strings"
)

type Logger struct {
	UseColor bool
	UseDebug bool
	Quiet    bool
	writer   *os.File
}

func New() *Logger {
	return &Logger{writer: os.Stdout}
}

func (l *Logger) SetLevel(level int) *Logger {
	switch level {
	case 0:
		l.writer = os.Stdin
	case 1:
		l.writer = os.Stdout
	case 2:
		l.writer = os.Stderr
	default:
		panic("invalid log level")
	}

	return l
}

var colorMap = strings.NewReplacer("$R", "\x1b[31m", "$Y", "\x1b[33m", "$B", "\x1b[34m", "$Z", "\x1b[0m")
var blankMap = strings.NewReplacer("$R", "\x1b[0m", "$Y", "\x1b[0m", "$B", "\x1b[0m", "$Z", "\x1b[0m")

func (l *Logger) color(data string, args ...interface{}) string {
	str := fmt.Sprintf(data, args...)
	if l.UseColor {
		return colorMap.Replace(str)
	} else {
		return blankMap.Replace(str)
	}
}

func (l *Logger) Debug(data string, args ...interface{}) {
	if l.UseDebug {
		l.writer.WriteString(fmt.Sprintf(data, args...))
	}
}

func (l *Logger) Line(data string, args ...interface{}) *Logger {
	l.writer.WriteString(fmt.Sprintf(data, args...))
	return l
}

func (l *Logger) WithCmd(cmd string) *Logger {
	l.writer.WriteString(l.color("run '%s' for more information", cmd))
	return l
}

func (l *Logger) Info(data string, args ...interface{}) {
	l.writer.WriteString(l.color("$Binfo$Z: "))
	l.writer.WriteString(fmt.Sprintf(data, args...) + "\n")
}

func (l *Logger) Warn(data string, args ...interface{}) {
	l.writer.WriteString(l.color("$Ywarn$Z: "))
	l.writer.WriteString(fmt.Sprintf(data, args...) + "\n")
}

func (l *Logger) Error(data string, args ...interface{}) *Logger {
	os.Stderr.WriteString(l.color("$Rerror$Z: "))
	os.Stderr.WriteString(fmt.Sprintf(data, args...) + "\n")
	return l
}

func (l *Logger) WithError(err error) *Logger {
	l.Error(err.Error())
	return l
}
