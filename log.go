package slog

import (
	"fmt"
	"io"
	"os"
	"sync"
)

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

// Hook is a callback function
type Hook func(message string) error

// var logs []string
var hooks map[int][]Hook
var minlevel int
var m sync.Mutex

var out io.Writer = os.Stderr

// RegisterHook will execute given hook function on every message
func RegisterHook(h Hook, level int) {
	hooks[level] = append(hooks[level], h)
}

// SetMinLevel sets the log level below that messages will be dropped
func SetMinLevel(level int) {
	minlevel = level
}

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) {
	m.Lock()
	defer m.Unlock()
	out = w
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
func Infof(message string, a ...interface{}) {
	f(InfoLog, message, a...)
}

// Errorf logs to dashboard (sends message through log channel) plus echoes to standard output
func Errorf(message string, a ...interface{}) {
	f(ErrorLog, message, a...)
}

// Fatalf logs to dashboard (sends message through log channel) plus echoes to standard output
func Fatalf(message string, a ...interface{}) {
	f(FatalLog, message, a...)
	os.Exit(1)
}

func f(level int, message string, a ...interface{}) {
	if level < minlevel {
		return
	}
	if len(hooks[level]) > 0 {
		for _, h := range hooks[level] {
			if h == nil {
				continue
			}
			err := h(fmt.Sprintf(message, a...))
			if err != nil {
				fmt.Fprintf(out, "Log: hook returned error %#v %v\n", h, err)
				// store(level, fmt.Sprintf("Log: hook returned error %#v %v\n", h, err))
			}
		}
	}
	if len(a) > 0 {
		message = truncateString(fmt.Sprintf(message, a...), 5000)
	}
	fmt.Fprintln(out, message)
	//store(level, message)
}

// // Filter returns logs for minimum level
// func Filter(minlevel int) (ret []string) {
// 	for _, s := range logs {
// 		z, _ := strconv.Atoi(s[0:1])
// 		if z&^minlevel > 0 {
// 			ret = append(ret, s[1:])
// 		}
// 	}
// 	return
// }

func init() {
	hooks = make(map[int][]Hook, 1)
}
