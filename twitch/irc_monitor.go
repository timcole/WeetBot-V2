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

		go bot.callEvent(m)
	}
}

func (bot *Bot) callEvent(m *Message) {
	var event interface{}
	switch m.Command {
	case "PRIVMSG":
		event = bot.irc.events.onNewMessage
		break
	case "WHISPER":
		event = bot.irc.events.onNewWhisper
		break
	case "USERNOTICE":
		if m.Data.NoticeType == "sub" || m.Data.NoticeType == "resub" || m.Data.NoticeType == "subgift" {
			event = bot.irc.events.onNewSub

			dbug, _ := json.Marshal(m)
			fmt.Println(string(dbug))

			js, _ := json.Marshal(m.Data)
			bot.Say(m.Data.StreamerName, string(js))
			break
		} else if m.Data.NoticeType == "raid" {
			event = bot.irc.events.onNewRaid
			break
		} else {
			dbug, _ := json.Marshal(m)
			fmt.Println(string(dbug))
		}
		break
	}

	type handlerType func(m *Message)
	if f, ok := event.(func(*Message)); ok {
		go handlerType(f)(m)
	} else if bot.verbose == true {
		fmt.Println(m.String())
	}
}

// OnNewMessage fires on a new message
func (bot *Bot) OnNewMessage(cb interface{}) {
	bot.irc.events.onNewMessage = cb
}

// OnNewWhisper fires on a new whisper
func (bot *Bot) OnNewWhisper(cb interface{}) {
	bot.irc.events.onNewWhisper = cb
}

// OnNewSub fires on a new sub
func (bot *Bot) OnNewSub(cb interface{}) {
	bot.irc.events.onNewSub = cb
}

// OnNewRaid fires on a new raid
func (bot *Bot) OnNewRaid(cb interface{}) {
	bot.irc.events.onNewRaid = cb
}
