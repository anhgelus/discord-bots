package utils

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
)

var reset = "\033[0m"

/*
SendSuccess
Send input strings as ANSI color code for console
As success the color is green.
*/
func SendSuccess(message string) {
	log.Default().Println(logger.Green, message, reset)
}

/*
SendWarn
Send input strings as ANSI color code for console
As warning the color is yellow.
*/
func SendWarn(message string) {
	log.Default().Println(logger.Yellow, message, reset)
}

/*
SendDebug
Send input strings as ANSI color code for console
As debug the color is ??.
*/
func SendDebug(message ...any) {
	log.Default().Println(logger.Cyan, message, reset)
}

/*
SendAlert
Send input strings as ANSI color code for console
As warning the color is red.
*/
func SendAlert(pos string, message string) {
	color := "\033[31m" // red color in ANSI
	fmt.Printf("[%s] %s%s%s\n", pos, color, message, reset)
	log.Default().Println(logger.Red, fmt.Sprintf("[%s] %s%s%s", pos, color, message, reset))
}

func SendError(err error) {
	log.Default().Panic(err)
}
