package twitch

import (
	"bufio"
	"net"

	"github.com/gorilla/websocket"
)

type IRCConnection struct {
	conn net.Conn
	rd   *bufio.Reader
}

type WSConnection struct {
	conn *websocket.Conn
}

type Connections struct {
	WS  *WSConnection
	IRC *IRCConnection
}
