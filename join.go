package main

import (
	"encoding/json"
	"fmt"

	"github.com/TimothyCole/WeetBot-V2/twitch"
)

func joinChannel(bot *twitch.Bot, channelID int64, channelName string) {
	bot.Join(channelName)

	testTopic := fmt.Sprintf("video-playback-by-id.%v", channelID)
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
}
