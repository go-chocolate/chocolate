package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-chocolate/chocolate/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

type Formatter func(e *logrus.Entry) ([]byte, error)

func (f Formatter) Format(e *logrus.Entry) ([]byte, error) {
	return f(e)
}

var (
	JSONFormatter = Formatter(jsonFormatter)
	TextFormatter = Formatter(textFormatter)

	globalFields = map[string]any{}
)

func SetField(key string, val any) {
	globalFields[key] = val
}

func SetFields(fields map[string]any) {
	for k, v := range fields {
		SetField(k, v)
	}
}

type logEntry struct {
	traceId  string
	fields   map[string]any
	time     time.Time
	level    logrus.Level
	file     string
	line     int
	function string
	message  string
}

func (e *logEntry) build() map[string]any {
	m := make(map[string]any)
	for k, v := range globalFields {
		m[k] = v
	}
	if e.traceId != "" {
		m["trace_id"] = e.traceId
	}
	for k, v := range e.fields {
		m[k] = v
	}

	if e.level != 0 {
		m["level"] = e.level.String()
	}
	if e.file != "" {
		m["file"] = fmt.Sprintf("%s:%d", e.file, e.line)
	}
	if e.function != "" {
		m["function"] = e.function
	}

	if e.message != "" {
		m["message"] = e.message
	}
	return m
}

func (e *logEntry) text() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	for k, v := range e.build() {
		fmt.Fprintf(b, "%s=%v ", k, v)
	}
	b.Write([]byte{'\n'})
	return b.Bytes(), nil
}

func (e *logEntry) json() ([]byte, error) {
	m := e.build()
	return json.Marshal(m)
}

func build(e *logrus.Entry) *logEntry {
	entry := &logEntry{
		fields:  e.Data,
		time:    e.Time,
		level:   e.Level,
		message: e.Message,
	}
	if e.Context != nil {
		entry.traceId = telemetry.TraceIDFromContext(e.Context)
	}
	if e.Caller != nil {
		entry.file = e.Caller.File
		entry.line = e.Caller.Line
		entry.function = e.Caller.Function
	}
	return entry
}

func jsonFormatter(e *logrus.Entry) ([]byte, error) {
	return build(e).json()

}

func textFormatter(e *logrus.Entry) ([]byte, error) {
	return build(e).text()
}
