package main

import (
	"fmt"

	"github.com/TimothyCole/WeetBot-V2/twitch"
)

func joinChannel(bot *twitch.Bot, channelID int64, channelName string) {
	bot.Join(channelName)

	testTopic := fmt.Sprintf("video-playback-by-id.%v", channelID)
	bot.Listen(testTopic)
	bot.AddTopicHandler(testTopic, func(msg twitch.PubSubResponse) {
		fmt.Println(msg)
	})
}
