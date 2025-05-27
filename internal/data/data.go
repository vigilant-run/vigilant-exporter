package data

import (
	"time"
)

type LogLevel = string

const (
	LogLevelTrace       LogLevel = "TRACE"
	LogLevelDebugName   LogLevel = "DEBUG"
	LogLevelInfoName    LogLevel = "INFO"
	LogLevelWarningName LogLevel = "WARNING"
	LogLevelErrorName   LogLevel = "ERROR"
	LogLevelFatalName   LogLevel = "FATAL"
)

type Log struct {
	Timestamp  time.Time         `json:"timestamp"`
	Level      LogLevel          `json:"level"`
	Body       string            `json:"body"`
	Attributes map[string]string `json:"attributes"`
}

func NewLog(
	timestamp time.Time,
	level LogLevel,
	body string,
	attributes map[string]string,
) *Log {
	return &Log{
		Timestamp:  timestamp,
		Level:      level,
		Body:       body,
		Attributes: attributes,
	}
}

type LogBatch struct {
	Logs []*Log `json:"logs"`
}

func NewLogBatch(
	logs []*Log,
) *LogBatch {
	return &LogBatch{
		Logs: logs,
	}
}

type MessageBatch struct {
	Token string `json:"token"`
	Logs  []*Log `json:"logs"`
}

func NewMessageBatch(
	token string,
	logs []*Log,
) *MessageBatch {
	return &MessageBatch{
		Token: token,
		Logs:  logs,
	}
}
