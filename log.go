package slog

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	infoLog    = 1
	warningLog = 3
	errorLog   = 7
	fatalLog   = 15
)

var logs []string

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

// Send sends message to the channel with ISO timestamp
func store(severity int, message string, a ...interface{}) {
	go func() {
		logs = append(logs, strconv.Itoa(severity)+time.Now().Format("2006-01-02 15:04:05")+": "+fmt.Sprintf(message, a...))
	}()
}

// Infof logs to dashboard (sends message through log channel) plus echoes to standard output
func Infof(message string, a ...interface{}) {
	if len(a) > 0 {
		message = truncateString(fmt.Sprintf(message, a...), 5000)
	}
	log.Println(message)
	store(infoLog, message)
}

// Errorf logs to dashboard (sends message through log channel) plus echoes to standard output
func Errorf(message string, a ...interface{}) {
	if len(a) > 0 {
		message = truncateString(fmt.Sprintf(message, a...), 5000)
	}
	log.Println(message)
	store(errorLog, message)
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
