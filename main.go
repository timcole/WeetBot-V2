package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/TimothyCole/WeetBot-v2/twitch"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	wg := new(sync.WaitGroup)

	// IRC CONNNECTION
	c, err := twitch.IRCConnect()
	if err != nil {
		log.Fatalf("NO IRC CONNECTION WutFace %s", err)
	}

	c.Send("USER " + os.Getenv("TWITCH_BOT_NAME"))
	c.Send("PASS oauth:" + os.Getenv("TWITCH_BOT_AUTH"))
	c.Send("NICK " + os.Getenv("TWITCH_BOT_NAME"))
	c.Send("CAP REQ :twitch.tv/commands twitch.tv/tags")

	c.Send("JOIN #modesttim")

	wg.Add(1)
	go monitor(c)

	// PUBSUB CONNECTION
	ws, err := twitch.WSConnect()
	if err != nil {
		log.Fatalf("NO WS CONNECTION WutFace %s", err)
	}
	ws.Send(`{ "type": "LISTEN", "data": { "topics": ["video-playback-by-id.51684790"] } }`)
	ws.AddHandler("video-playback-by-id.51684790", func(m twitch.PubSubResponse) {
		fmt.Println(m.Data.RawMessage)
		fmt.Println("There is", m.Data.Message.Viewers, "viewers")
	})

	ws.Send(`{ "type": "LISTEN", "data": { "topics": ["following.51684790"] } }`)
	ws.AddHandler("following.29829912", func(m twitch.PubSubResponse) {
		f := m.Data.Message
		fmt.Println("New Followers!", f.DisplayName)
	})

	wg.Add(1)

	wg.Wait()
}

func monitor(c *twitch.Conn) {
	for {
		_, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
		}

		// fmt.Println(m.Data.DisplayName, m.Data.Message)
	}
}
