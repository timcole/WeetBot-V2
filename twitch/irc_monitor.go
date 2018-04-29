package twitch

import (
	"fmt"
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

		if m.Command == "PING" {
			c.Send("PONG %s", m.Trailing)
		}

		fmt.Println(m.Data.DisplayName, m.Data.Message)
	}
}
