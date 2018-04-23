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
	c, err := twitch.Connect()
	if err != nil {
		log.Fatalf("NO CONNECTION WutFace %s", err)
	}

	c.Send("USER " + os.Getenv("TWITCH_BOT_NAME"))
	c.Send("PASS oauth:" + os.Getenv("TWITCH_BOT_AUTH"))
	c.Send("NICK " + os.Getenv("TWITCH_BOT_NAME"))
	c.Send("CAP REQ :twitch.tv/commands twitch.tv/tags")

	c.Send("JOIN #fedmyster") // Don't stay in FEDMYSTERs chat lol just need messages
	c.Send("JOIN #modesttim")

	wg.Add(2)
	go monitor(c)
	wg.Wait()
}

func monitor(c *twitch.Conn) {
	for {
		m, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(m.Data.DisplayName, m.Data.Message)
	}
}
