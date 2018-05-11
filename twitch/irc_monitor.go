package twitch

import (
	"fmt"
	"os"
	"sync"
)

func (c *IRCConnection) Listen(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for {
		m, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		if m.Nick() == os.Getenv("TWITCH_BOT_NAME") {
			continue // Ignore us
		}

		if m.Command == "PING" {
			c.Send("PONG %s", m.Trailing)
		}

		if m.Data.Message == "modestYO" {
			c.Send("PRIVMSG #" + m.Data.StreamerName + " :modestYO")
		}

		// fmt.Println(m.Data.DisplayName, m.Data.Message)
		fmt.Println(m.Data)
	}
}
