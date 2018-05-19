package twitch

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func (bot *Bot) monitor() {
	for {
		m, err := bot.readMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		var verbose = true
		if verbose == true {
			fmt.Println(m.String())
		}

		if m.Nick() == bot.name {
			continue // Ignore us
		}

		if m.Command == "PING" {
			bot.SendRawIRC("PONG %s", m.Trailing)
		}

		if len(m.Data.Arguments) <= 0 {
			continue
		}

		if m.Data.Arguments[0] == "!debug" {
			var dbug []byte
			var f reflect.Value
			if len(m.Data.Arguments) > 1 {
				r := reflect.ValueOf(m.Data)
				f = reflect.Indirect(r)
				for i := 1; i < len(m.Data.Arguments); i++ {
					f = f.FieldByName(m.Data.Arguments[i])
				}

				if f.IsValid() {
					dbug, _ = json.Marshal(f.Interface())
				} else {
					dbug = []byte("Invalid debug")
				}
			} else {
				dbug, _ = json.Marshal(m.Data)
			}
			bot.Say(m.Data.StreamerName, string(dbug))
		}

		if m.Data.Arguments[0] == "!kill" && strings.ToLower(m.Data.DisplayName) == "modesttim" {
			bot.Say(m.Data.StreamerName, "BibleThump")
			bot.Done <- true
		}
	}
}
