package slog

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	InfoLog = 1 << iota
	WarningLog
	ErrorLog
	FatalLog
)

// Hook is a callback function
type Hook func(message string) error

var logs []string
var hooks map[int][]Hook

// RegisterHook will execute given hook function on every message
func RegisterHook(h Hook, level int) {
	hooks[level] = append(hooks[level], h)
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
func store(severity int, message string, a ...interface{}) {
	go func() {
		logs = append(logs, strconv.Itoa(severity)+time.Now().Format("2006-01-02 15:04:05")+": "+fmt.Sprintf(message, a...))
	}()
}

// Infof logs to dashboard (sends message through log channel) plus echoes to standard output
func Infof(message string, a ...interface{}) {
	f(InfoLog, message, a)
}

// Warningf logs to dashboard (sends message through log channel) plus echoes to standard output
func Warningf(message string, a ...interface{}) {
	f(WarningLog, message, a...)
}

// Errorf logs to dashboard (sends message through log channel) plus echoes to standard output
func Errorf(message string, a ...interface{}) {
	f(ErrorLog, message, a...)
}

// Fatalf logs to dashboard (sends message through log channel) plus echoes to standard output
func Fatalf(message string, a ...interface{}) {
	f(FatalLog, message, a...)
}

func f(level int, message string, a ...interface{}) {
	if len(hooks[level]) > 0 {
		for _, h := range hooks[level] {
			err := h(fmt.Sprintf(message, a...))
			if err != nil {
				log.Printf("Log: hook returned error %#v %v\n", h, err)
				store(level, fmt.Sprintf("Log: hook returned error %#v %v\n", h, err))
			}
		}
	}
	if len(a) > 0 {
		message = truncateString(fmt.Sprintf(message, a...), 5000)
	}
	log.Println(message)
	store(level, message)
}

// Filter returns logs for minimum level
func Filter(minlevel int) (ret []string) {
	for _, s := range logs {
		z, _ := strconv.Atoi(s[0:1])
		if z&^minlevel > 0 {
			ret = append(ret, s[1:])
		}
	}
	return
}

func init() {
	hooks = make(map[int][]Hook, 1)
}
