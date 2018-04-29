package twitch

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func (bot *Bot) IRCConnect() (*Bot, error) {
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		return nil, err
	}

	bot.IRC = &IRCConnection{
		conn: conn,
		rd:   bufio.NewReader(conn),
	}

	bot.IRC.Send("USER " + bot.Name)
	bot.IRC.Send("PASS oauth:" + strings.Replace(bot.OAuth, "oauth:", "", 0))
	bot.IRC.Send("NICK " + bot.Name)
	bot.IRC.Send("CAP REQ :twitch.tv/commands twitch.tv/tags")

	return bot, nil
}

func (c *IRCConnection) Send(format string, args ...interface{}) error {
	if _, err := fmt.Fprintf(c.conn, format, args...); err != nil {
		return err
	}
	if !strings.HasSuffix(format, "\r\n") {
		if _, err := fmt.Fprint(c.conn, "\r\n"); err != nil {
			return nil
		}
	}
	return nil
}

func (c *IRCConnection) ReadMessage() (*Message, error) {
	line, err := c.rd.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return ParseLine(line)
}

func (c *IRCConnection) Close() error {
	return c.conn.Close()
}
