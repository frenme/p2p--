package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	levelNames = map[Level]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
	}
	currentLevel = INFO
)

func SetLevel(level Level) {
	currentLevel = level
}

func Debug(format string, args ...interface{}) {
	if currentLevel <= DEBUG {
		logMessage(DEBUG, format, args...)
	}
}

func Info(format string, args ...interface{}) {
	if currentLevel <= INFO {
		logMessage(INFO, format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	if currentLevel <= WARN {
		logMessage(WARN, format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if currentLevel <= ERROR {
		logMessage(ERROR, format, args...)
	}
}

func logMessage(level Level, format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	levelName := levelNames[level]
	message := fmt.Sprintf(format, args...)
	
	output := fmt.Sprintf("[%s] %s: %s", timestamp, levelName, message)
	
	if level >= ERROR {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(os.Stdout)
	}
	
	log.SetFlags(0)
	log.Print(output)
}