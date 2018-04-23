package twitch

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	ws "github.com/gorilla/websocket"
)

func WSConnect() (*Connections, error) {
	c, _, err := ws.DefaultDialer.Dial("wss://pubsub-edge.twitch.tv/", nil)
	if err != nil {
		return nil, err
	}

	connection := &Connections{
		WSConnection: c,
	}

	go wsListen(connection)

	// Send pings every 2.5 minutes to keep the connection alive (5 minutes is the required time but 2.5 just to be safe)
	go func(c *Connections) {
		ticker := time.NewTicker(150 * time.Second)
		for _ = range ticker.C {
			c.Send(`{ "type": "PING" }`)
		}
	}(connection)

	return connection, nil
}

type PubSubResponse struct {
	Type  string `json:"type"`
	Nonce string `json:"nonce,omitempty"`
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

func wsListen(c *Connections) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			m := PubSubResponse{}
			err := c.WSConnection.ReadJSON(&m)
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
				HandlerType(f)(m)
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

func (c *Connections) AddHandler(topic string, q interface{}) {
	if EventHandlers[topic] == nil {
		EventHandlers[topic] = q
	}
}

func (c *Connections) Send(msg string) {
	c.WSConnection.WriteMessage(ws.TextMessage, []byte(msg))
}
