package twitch

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Conn struct {
	conn net.Conn
	rd   *bufio.Reader
}

func Connect() (*Conn, error) {
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		return nil, err
	}
	c := &Conn{
		conn: conn,
		rd:   bufio.NewReader(conn),
	}
	return c, nil
}

func (c *Conn) Send(format string, args ...interface{}) error {
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

func (c *Conn) ReadMessage() (*Message, error) {
	line, err := c.rd.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return ParseLine(line)
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
