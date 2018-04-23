package twitch

import (
	"github.com/gorilla/websocket"
)

type Connections struct {
	WSConnection *websocket.Conn
}
