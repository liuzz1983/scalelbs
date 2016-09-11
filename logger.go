package main

import (
	"github.com/uber-common/bark"
	"github.com/uber/tchannel-go"
)

type ProxyLogger struct {
	Logger bark.Logger
}

func (l ProxyLogger) Enabled(_ tchannel.LogLevel) bool { return false }

func (l ProxyLogger) Fatal(msg string) {
	l.Logger.Fatal(msg)
}
func (l ProxyLogger) Error(msg string) {
	l.Logger.Error(msg)
}
func (l ProxyLogger) Warn(msg string) {
	l.Logger.Warn(msg)
}
func (l ProxyLogger) Infof(msg string, args ...interface{}) {
	l.Logger.Infof(msg, args...)
}
func (l ProxyLogger) Info(msg string) {
	l.Logger.Info(msg)
}
func (l ProxyLogger) Debugf(msg string, args ...interface{}) {
	l.Logger.Debugf(msg, args...)
}
func (l ProxyLogger) Debug(msg string) {
	l.Logger.Debug(msg)
}
func (l ProxyLogger) Fields() tchannel.LogFields {
	return nil
}

func (l ProxyLogger) WithFields(fields ...tchannel.LogField) tchannel.Logger {
	logFields := make(bark.Fields)
	for _, field := range fields {
		logFields[field.Key] = field.Value
	}
	l.Logger.WithFields(logFields)
	return l
}
