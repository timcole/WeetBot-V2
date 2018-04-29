package main

import (
	"log"
	"os"
	"sync"

	"github.com/TimothyCole/WeetBot-v2/twitch"
	_ "github.com/joho/godotenv/autoload"
)

var (
	wg       = new(sync.WaitGroup)
	autoJoin = []string{
		"weetbot",
		"modesttim",
	}
)

func main() {
	bot := &twitch.Bot{
		Name:  os.Getenv("TWITCH_BOT_NAME"),
		OAuth: os.Getenv("TWITCH_BOT_AUTH"),
	}
	bot, err := bot.IRCConnect()
	if err != nil {
		log.Fatalf("NO IRC CONNECTION WutFace %s", err)
	}
	bot, err = bot.WSConnect()
	if err != nil {
		log.Fatalf("NO WS CONNECTION WutFace %s", err)
	}
	go bot.IRC.Listen(wg)
	go bot.WS.Listen(wg)

	go func(irc *twitch.IRCConnection) {
		for _, name := range autoJoin {
			go irc.Send("JOIN #" + name)
		}
	}(bot.IRC)

	// WS
	bot.WS.Send(`{ "type": "LISTEN", "data": { "topics": ["video-playback-by-id.51684790"] } }`)
	bot.WS.AddHandler("video-playback-by-id.51684790", func(m twitch.PubSubResponse) {
		// fmt.Println("There is", m.Data.Message.Viewers, "viewers")
	})

	bot.WS.Send(`{ "type": "LISTEN", "data": { "topics": ["following.51684790"] } }`)
	bot.WS.AddHandler("following.51684790", func(m twitch.PubSubResponse) {
		// f := m.Data.Message
		// fmt.Println("New Followers!", f.DisplayName)
	})

	wg.Wait()
}
