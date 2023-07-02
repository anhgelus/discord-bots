package main

import (
	start "github.com/anhgelus/discord-bots/les-copaings/src/init"
	"os"
)

func main() {
	start.Bot(os.Args[1])
}
