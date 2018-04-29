package twitch

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	ws "github.com/gorilla/websocket"
)

func (bot *Connections) WSConnect() (*Connections, error) {
	c, _, err := ws.DefaultDialer.Dial("wss://pubsub-edge.twitch.tv/", nil)
	if err != nil {
		return nil, err
	}

	// connection := &Connections{
	// 	WSConnection: c,
	// }

	bot.WS = &WSConnection{c}

	go wsListen(bot.WS)

	// Send pings every 2.5 minutes to keep the connection alive (5 minutes is the required time but 2.5 just to be safe)
	go func(c *WSConnection) {
		ticker := time.NewTicker(150 * time.Second)
		for _ = range ticker.C {
			c.Send(`{ "type": "PING" }`)
		}
	}(bot.WS)

	return bot, nil
}

type PubSubResponse struct {
	Type  string `json:"type"`
	Nonce string `json:"nonce,omitempty"`
	Error string `json:"error,omitempty"`
	Data  struct {
		Topics     []string `json:"topics,omitempty"`
		Topic      string   `json:"topic,omitempty"`
		AuthToken  string   `json:"auth_token,omitempty"`
		RawMessage string   `json:"message,omitempty"`
		Message    *PubSubMessage
	}
}

type PubSubMessage struct {
	DisplayName string `json:"display_name,omitempty"`
	Username    string `json:"username,omitempty"`
	UserID      string `json:"user_id,omitempty"`

	Type       string  `json:"type,omitempty"`
	ServerTime float64 `json:"server_time,omitempty"`
	Viewers    int     `json:"viewers,omitempty"`
}

var EventHandlers = make(map[string]interface{})

func wsListen(c *WSConnection) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			m := PubSubResponse{}
			err := c.conn.ReadJSON(&m)
			if err != nil {
				log.Println("JSON Read Error: ", err)
				return
			}

			if m.Data.RawMessage != "" {
				msg := &PubSubMessage{}
				err = json.Unmarshal([]byte(m.Data.RawMessage), msg)
				if err != nil {
					fmt.Println("NOOOO: ", err)
				}
				m.Data.Message = msg
			}

			type HandlerType func(m PubSubResponse)
			if f, ok := EventHandlers[m.Data.Topic].(func(PubSubResponse)); ok {
				go HandlerType(f)(m)
			} else {
				if m.Data.Topic != "" {
					fmt.Printf("not running func topic=%s interface=%v\n", m.Data.Topic, EventHandlers[m.Data.Topic])
				} else {
					fmt.Println(m)
				}
			}
		}
	}()
}

func (c *WSConnection) AddHandler(topic string, q interface{}) {
	if EventHandlers[topic] == nil {
		EventHandlers[topic] = q
	}
}

func (c *WSConnection) Send(msg string) {
	c.conn.WriteMessage(ws.TextMessage, []byte(msg))
}
