package slog

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	JSON("test")
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	SetOutput(w)
	Infof("test?")
	w.Flush()
	var logEntry struct {
		Level   string `json:"level"`
		Message string `json:"msg"`
		Service string `json:"service"`
		Time    time.Time
	}
	err := json.Unmarshal(b.Bytes(), &logEntry)
	if err != nil {
		t.Errorf("can't decode logentry '%s'", b.String())
	}
	if logEntry.Level != "info" {
		t.Errorf("Log level is not info %s", b.String())
	}
	if logEntry.Service != "test" {
		t.Errorf("Log service is not test %s", b.String())
	}
	if logEntry.Message != "test?" {
		t.Errorf("Log message is not 'test?' %s", b.String())
	}
	if logEntry.Time.After(time.Now()) { //well, testing time...
		t.Errorf("Log time is in the future %s", b.String())
	}
}
