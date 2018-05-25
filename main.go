package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TimothyCole/WeetBot-V2/twitch"
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

	bot.Join("modesttim")

	testTopic := "video-playback-by-id.51684790"
	bot.Listen(testTopic)
	bot.AddTopicHandler(testTopic, func(msg twitch.PubSubResponse) {
		// fmt.Println(msg)
	})
	bot.OnNewMessage(func(msg *twitch.Message) {
		fmt.Println("> New Message: ", msg.Data.DisplayName, msg.Data.Message)
	})
	bot.OnNewWhisper(func(msg *twitch.Message) {
		fmt.Println("> New Whisper: ", msg.Data.DisplayName, msg.Data.Message)
	})
	bot.OnNewSub(func(msg *twitch.Message) {
		fmt.Println("> New Sub: ", msg.Data.DisplayName, msg.Data.Sub.Plan, msg.Data.GiftSub.Login)

		dbug, _ := json.Marshal(msg)
		bot.Say(msg.Data.StreamerName, string(dbug))
	})
	bot.OnNewRaid(func(msg *twitch.Message) {
		fmt.Println("> New Raid: ", msg.Data.Raid.DisplayName, msg.Data.Raid.Viewers)
	})

	<-bot.Done
}
