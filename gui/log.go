package gui

import "github.com/sirupsen/logrus"

type Logger struct{}

func (l *Logger) Write(p []byte) (n int, err error) {
	UI.AddLog(string(p))

	return len(p), nil
}

type LogFormatter struct{}

func (l *LogFormatter) Format(logLine *logrus.Entry) ([]byte, error) {
	return []byte(logLine.Message + "\n"), nil
}
