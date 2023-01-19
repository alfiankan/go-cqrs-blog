package common

import (
	"log"
)

const (
	LOG_LEVEL_ERROR = "error"
	LOG_LEVEL_INFO  = "info"
	LOG_LEVEL_WARN  = "warn"
)

func Log(level, message string) {
	log.Printf("level=%s message=%s", level, message)
}

func LogExit(err error, level string) {
	if err != nil {
		log.Fatalf("level=%s message=%s", level, err.Error())
	}
}
