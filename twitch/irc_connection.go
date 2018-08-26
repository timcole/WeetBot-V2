package twitch

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func (bot *Bot) ircConnect() error {
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		return err
	}

	bot.irc.conn = conn
	bot.irc.rd = bufio.NewReader(conn)

	return nil
}

func (bot *Bot) readMessage() (*Message, error) {
	line, err := bot.irc.rd.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return parseLine(line)
}

// SendRawIRC a raw message over IRC
func (bot *Bot) SendRawIRC(format string, args ...interface{}) error {
	if _, err := fmt.Fprintf(bot.irc.conn, format, args...); err != nil {
		return err
	}
	if !strings.HasSuffix(format, "\r\n") {
		if _, err := fmt.Fprint(bot.irc.conn, "\r\n"); err != nil {
			return nil
		}
	}
	return nil
}

// Close IRC connection
func (bot *Bot) Close() error {
	return bot.irc.conn.Close()
}

// Say Sends a Private Channel Message
func (bot *Bot) Say(channel string, message ...string) {
	var msg = strings.Join(message, " ")
	fmt.Println("PRIVMSG #" + channel + " :" + msg)
	bot.SendRawIRC("PRIVMSG #" + channel + " :" + msg)
}

// Join a channel
func (bot *Bot) Join(channel string) {
	bot.SendRawIRC("JOIN #" + strings.ToLower(channel))
}

// Part a channel
func (bot *Bot) Part(channel string) {
	bot.SendRawIRC("PART #" + strings.ToLower(channel))
}

// Whisper a user
func (bot *Bot) Whisper(user, message string) {
	bot.Say("jtv", "/w", user, message)
}
