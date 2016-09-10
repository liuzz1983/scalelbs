package main

import (
	"github.com/uber-common/bark"
	"github.com/uber/tchannel-go"
)

type LbsLogger struct {
	Logger bark.Logger
}

func (l LbsLogger) Enabled(_ tchannel.LogLevel) bool { return false }

func (l LbsLogger) Fatal(msg string) {
	l.Logger.Fatal(msg)
}
func (l LbsLogger) Error(msg string) {
	l.Logger.Error(msg)
}
func (l LbsLogger) Warn(msg string) {
	l.Logger.Warn(msg)
}
func (l LbsLogger) Infof(msg string, args ...interface{}) {
	l.Logger.Infof(msg, args...)
}
func (l LbsLogger) Info(msg string) {
	l.Logger.Info(msg)
}
func (l LbsLogger) Debugf(msg string, args ...interface{}) {
	l.Logger.Debugf(msg, args...)
}
func (l LbsLogger) Debug(msg string) {
	l.Logger.Debug(msg)
}
func (l LbsLogger) Fields() tchannel.LogFields {
	return nil
}

func (l LbsLogger) WithFields(fields ...tchannel.LogField) tchannel.Logger {
	logFields := make(bark.Fields)
	for _, field := range fields {
		logFields[field.Key] = field.Value
	}
	l.Logger.WithFields(logFields)
	return l
}
