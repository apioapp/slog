// Slog uses logrus and lumberjack to provide a simple logging interface with two main funcs:
// Infof(string, ...interface{}) and Errorf(string, ...interface{})
// for logging to stdout and logfile at the same time

package slog

import (
	"fmt"
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
)

const LogFilePath = "logs/misc.log"

const (
	// InfoLog is Info level (lowest) for SetMinLevel
	InfoLog = iota + 1
	// ErrorLog is Error level for SetMinLevel, between WarningLog and FatalLog
	ErrorLog
	// FatalLog is highest log level for SetMinLevel (logging into Fatalf will also throw panic())
	FatalLog
)

// HookFunc is a callback function type, getting message as argument
type HookFunc func(message string) error

// LevelHookFunc is hook func with log level input attribute
type LevelHookFunc func(message string, level int) error

var lumberjackLogrotate *lumberjack.Logger
var maxlen = 5000
var hooks map[int][]HookFunc
var lhooks map[int][]LevelHookFunc

// SetOutput sets the standard logger output to a writer
func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

// SetFilename modifies the output file path
func SetFilename(filename string) {
	lumberjackLogrotate.Filename = filename
}

// RegisterHook will execute given hook function on every message
func RegisterHook(h HookFunc, minlevel int) {
	if hooks == nil {
		hooks = map[int][]HookFunc{}
	}
	for l := minlevel; l <= FatalLog; l++ {
		hooks[l] = append(hooks[l], h)
	}
	//	log.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogrotate))
}

// RegisterLevelHook will execute given hook function on every message (message and level)
func RegisterLevelHook(h LevelHookFunc, minlevel int) {
	if lhooks == nil {
		lhooks = map[int][]LevelHookFunc{}
	}
	for l := minlevel; l <= FatalLog; l++ {
		lhooks[l] = append(lhooks[l], h)
	}
	//	log.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogrotate))
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
	if hooks != nil && hooks[InfoLog] != nil {
		for _, h := range hooks[InfoLog] {
			h(fmt.Sprintf(message, args...))
		}
	}
	if lhooks != nil && lhooks[InfoLog] != nil {
		for _, h := range lhooks[InfoLog] {
			h(fmt.Sprintf(message, args...), InfoLog)
		}
	}
	log.Infof(message, args...)
}

// Errorf logs to dashboard (sends message through log channel) plus echoes to standard output
func Errorf(message string, args ...interface{}) {
	if hooks != nil && hooks[ErrorLog] != nil {
		for _, h := range hooks[ErrorLog] {
			h(fmt.Sprintf(message, args...))
		}
	}
	if lhooks != nil && lhooks[ErrorLog] != nil {
		for _, h := range lhooks[ErrorLog] {
			h(fmt.Sprintf(message, args...), ErrorLog)
		}
	}
	log.Errorf(message, args...)
}

// Fatalf logs to dashboard (sends message through log channel) plus echoes to standard output
func Fatalf(message string, args ...interface{}) {
	if hooks != nil && hooks[FatalLog] != nil {
		for _, h := range hooks[FatalLog] {
			h(fmt.Sprintf(message, args...))
		}
	}
	if lhooks != nil && lhooks[FatalLog] != nil {
		for _, h := range lhooks[FatalLog] {
			h(fmt.Sprintf(message, args...), FatalLog)
		}
	}
	log.Fatalf(message, args...)
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

	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	},
	)

	logMultiWriter := io.MultiWriter(os.Stdout, lumberjackLogrotate)
	log.SetOutput(logMultiWriter)
}
