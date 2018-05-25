package twitch

import (
	"bufio"
	"net"
	"strings"

	"github.com/gorilla/websocket"
)

// Bot is the structure for a bot client instance
type Bot struct {
	name  string
	oauth string
	Done  chan bool

	verbose bool
	ws      *websocket.Conn
	irc     struct {
		conn   net.Conn
		rd     *bufio.Reader
		events struct {
			onNewMessage interface{}
			onNewWhisper interface{}
			onNewSub     interface{}
			onNewRaid    interface{}
		}
	}
}

// NewClient Initalizes the Bot struct with needed data
func NewClient(name, oauth string) *Bot {
	return &Bot{
		name:    name,
		oauth:   oauth,
		Done:    make(chan bool),
		verbose: true,
	}
}

// Connect to IRC and PubSub
func (bot *Bot) Connect() error {
	// IRC Connection
	if err := bot.ircConnect(); err != nil {
		return err
	}

	// Send the needed data to Twitch to connect
	go func(bot *Bot) {
		bot.SendRawIRC("USER " + bot.name)
		bot.SendRawIRC("PASS oauth:" + strings.Replace(bot.oauth, "oauth:", "", 0))
		bot.SendRawIRC("NICK " + bot.name)
		bot.SendRawIRC("CAP REQ :twitch.tv/commands twitch.tv/tags")
	}(bot)
	go bot.monitor()

	// PubSub Connection
	if err := bot.wsConnect(); err != nil {
		return err
	}
	go bot.subscribe()

	return nil
}
