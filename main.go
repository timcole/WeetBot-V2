package main

import (
	"fmt"
	"os"

	"github.com/TimothyCole/WeetBot-v2/twitch"
	_ "github.com/joho/godotenv/autoload"
)

var (
	settings = struct {
		nick string
		pass string
	}{
		nick: os.Getenv("TWITCH_BOT_NAME"),
		pass: os.Getenv("TWITCH_BOT_AUTH"),
	}
)

func main() {
	bot := twitch.NewClient(settings.nick, settings.pass)
	if err := bot.Connect(); err != nil {
		panic(err)
	}

	bot.Join("jamie254")
	bot.Join("weetbot")

	fmt.Println("Running :D")

	<-bot.Done
}
