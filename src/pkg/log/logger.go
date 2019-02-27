// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

// L is an alias for the the standard logger.
var L = logrus.NewEntry(logrus.StandardLogger())

type loggerKey struct{}

// WithContext returns a new context with the provided logger. Use in
// combination with logger.WithField(s) for great effect.
func WithContext(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// FromContext retrieves the current logger from the context. If no
// logger is available, the default logger is returned.
func FromContext(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(loggerKey{})
	if logger == nil {
		return L
	}
	return logger.(*logrus.Entry)
}

// FromRequest retrieves the current logger from the request. If no
// logger is available, the default logger is returned.
func FromRequest(r *http.Request) *logrus.Entry {
	return FromContext(r.Context())
}

func WithError(err error) *logrus.Entry {
	return logrus.WithError(err)
}

// Debug logs debug level messages with default logger.
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Debugf logs debug level messages with default logger in printf-style.
func Debugf(msg string, args ...interface{}) {
	logrus.Debugf(msg, args...)
}

// Info logs Info level messages with default logger in structured-style.
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Infof logs Info level messages with default logger in printf-style.
func Infof(msg string, args ...interface{}) {
	logrus.Infof(msg, args...)
}

// Warn logs Warn level messages with default logger in structured-style.
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

// Warnf logs Warn level messages with default logger in printf-style.
func Warnf(msg string, args ...interface{}) {
	logrus.Warnf(msg, args...)
}

// Error logs Error level messages with default logger in structured-style.
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Errorf logs Error level messages with default logger in printf-style.
func Errorf(msg string, args ...interface{}) {
	logrus.Errorf(msg, args...)
}

// Panic logs Panic level messages with default logger in structured-style.
func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

// Panicf logs Panicf level messages with default logger in printf-style.
func Panicf(msg string, args ...interface{}) {
	logrus.Panicf(msg, args...)
}

// Fatal logs Fatal level messages with default logger in structured-style.
func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

// Fatalf logs Fatalf level messages with default logger in printf-style.
func Fatalf(msg string, args ...interface{}) {
	logrus.Fatalf(msg, args...)
}
