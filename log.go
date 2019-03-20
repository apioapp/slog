package slog

import (
	"io"
	"os"
	"time"

	"runtime"

	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
)

const LogFilePath = "logs/misc.log"

const (
	// InfoLog is Info level (lowest) for SetMinLevel
	InfoLog = 1 << iota
	// WarningLog is Warning level for SetMinLevel, between InfoLog and ErrorLog
	WarningLog
	// ErrorLog is Error level for SetMinLevel, between WarningLog and FatalLog
	ErrorLog
	// FatalLog is highest log level for SetMinLevel (logging into Fatalf will also throw panic())
	FatalLog
)

// HookFunc is a callback function type, getting message as argument
type HookFunc func(message string) error

var lumberjackLogrotate *lumberjack.Logger
var maxlen = 5000

// RegisterHook will execute given hook function on every message
func RegisterHook(h HookFunc) {
	log.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogrotate, HookWriter{Hook: h}))
}

func truncateString(str string, num int) string {
	s := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		s = str[0:num] + "..."
	}
	return s
}

// Send sends message to the channel with ISO timestamp
// func store(severity int, message string, a ...interface{}) {
// 	go func() {
// 		m.Lock()
// 		logs = append(logs, strconv.Itoa(severity)+time.Now().Format("2006-01-02 15:04:05")+": "+fmt.Sprintf(message, a...))
// 		m.Unlock()
// 	}()
// }

// Infof logs to dashboard (sends message through log channel) plus echoes to standard output
func Infof(message string, args ...interface{}) {
	log.Infof(message, args...)
}

// Errorf logs to dashboard (sends message through log channel) plus echoes to standard output
func Errorf(message string, args ...interface{}) {
	log.Errorf(message, args...)
}

// Fatalf logs to dashboard (sends message through log channel) plus echoes to standard output
func Fatalf(message string, args ...interface{}) {
	log.Fatalf(message, args...)
}

type HookWriter struct {
	Hook HookFunc
}

func (h HookWriter) Write(p []byte) (n int, err error) {
	h.Hook(string(p))
	return len(p), nil
}

func init() {
	// Setup logger
	lumberjackLogrotate = &lumberjack.Logger{
		Filename:   LogFilePath,
		MaxSize:    5, // Max megabytes before log is rotated
		MaxBackups: 0, // Max number of old log files to keep
		MaxAge:     0, // Max number of days to retain log files
		Compress:   true,
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, TimestampFormat: time.RFC1123Z})

	logMultiWriter := io.MultiWriter(os.Stdout, lumberjackLogrotate)
	log.SetOutput(logMultiWriter)

	log.WithFields(log.Fields{
		"Runtime Version": runtime.Version(),
		"Number of CPUs":  runtime.NumCPU(),
		"Arch":            runtime.GOARCH,
	}).Info("Application Initializing")
}
