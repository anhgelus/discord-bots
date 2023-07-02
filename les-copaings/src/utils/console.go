package utils

import "fmt"

var reset = "\033[0m"

/*
SendSuccess
Send input strings as ANSI color code for console
As success the color is green.
*/
func SendSuccess(message string) {
	color := "\033[32m" // green color in ANSI
	fmt.Println(color, message, reset)
}

/*
SendWarn
Send input strings as ANSI color code for console
As warning the color is yellow.
*/
func SendWarn(message string) {
	color := "\033[33m" // yellow color in ANSI
	fmt.Println(color, message, reset)
}

/*
SendDebug
Send input strings as ANSI color code for console
As debug the color is ??.
*/
func SendDebug(message any) {
	color := "\033[34m"
	fmt.Println(color, message, reset)
}

/*
SendAlert
Send input strings as ANSI color code for console
As warning the color is red.
*/
func SendAlert(message string) {
	color := "\033[31m" // red color in ANSI
	fmt.Println(color, message, reset)
}

func SendError(err error) {
	panic(err)
}
