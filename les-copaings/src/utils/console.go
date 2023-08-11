package utils

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
)

var reset = "\033[0m"

// SendSuccess sends a success message
func SendSuccess(message string) {
	log.Default().Println(logger.Green, message, reset)
}

// SendWarn sends a warning message
func SendWarn(message string) {
	log.Default().Println(logger.Yellow, message, reset)
}

// SendDebug sends a debug message
func SendDebug(message ...any) {
	log.Default().Println(logger.Cyan, message, reset)
}

// SendAlert sends an alert
func SendAlert(pos string, message string) {
	log.Default().Println(fmt.Sprintf("[%s] %s%s%s", pos, logger.Red, message, reset))
}

// SendError sends an error (like a panic(any...))
func SendError(err error) {
	log.Default().Panic(err)
}
