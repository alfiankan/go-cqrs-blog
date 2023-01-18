package common

import (
	"log"
	"os"
)

const (
	LOG_LEVEL_ERROR = "error"
	LOG_LEVEL_INFO  = "info"
	LOG_LEVEL_WARN  = "warn"
)

func Log(level, message string) {
	log.Printf("level=%s message=%s", level, message)
}

func LogExit(level, message string) {
	log.Printf("level=%s message=%s", level, message)
	os.Exit(1)
}
