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

// Send sends message to the channel with timestamp
func store(severity int, message string, a ...interface{}) {
	go func() {
		logs = append(logs, strconv.Itoa(severity)+time.Now().Format("02/01/2006 15:04:05")+": "+fmt.Sprintf(message, a...))
	}()
}

// Infof logs to dashboard (sends message through log channel) plus echoes to standard output
func Infof(message string, a ...interface{}) {
	if len(a) > 0 {
		log.Printf(message+"\n", a...)
	} else {
		log.Println(message)
	}
	store(infoLog, message, a...)
}

// Errorf logs to dashboard (sends message through log channel) plus echoes to standard output
func Errorf(message string, a ...interface{}) {
	if len(a) > 0 {
		log.Printf(message+"\n", a...)
	} else {
		log.Println(message)
	}
	store(errorLog, message, a...)
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
