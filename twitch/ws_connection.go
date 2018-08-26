package twitch

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	ws "github.com/gorilla/websocket"
)

func (bot *Bot) wsConnect() error {
	c, _, err := ws.DefaultDialer.Dial("wss://pubsub-edge.twitch.tv/", nil)
	if err != nil {
		return err
	}

	// Send pings every 2.5 minutes to keep the connection alive (5 minutes is the required time but 2.5 just to be safe)
	go func(bot *Bot) {
		var ticker = time.NewTicker(150 * time.Second)
		for range ticker.C {
			bot.wsSend(`{ "type": "PING" }`)
		}
	}(bot)

	bot.ws = c

	return nil
}

// PubSubResponse is the response from PubSub
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

// PubSubMessage is a message part of the PubSubResponse
type PubSubMessage struct {
	DisplayName string `json:"display_name,omitempty"`
	Username    string `json:"username,omitempty"`
	UserID      string `json:"user_id,omitempty"`

	Type       string  `json:"type,omitempty"`
	ServerTime float64 `json:"server_time,omitempty"`
	Viewers    int     `json:"viewers,omitempty"`
}

var eventHandlers = make(map[string]interface{})

func (bot *Bot) subscribe() {
	for {
		m := PubSubResponse{}
		err := bot.ws.ReadJSON(&m)
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

		type handlerType func(m PubSubResponse)
		if f, ok := eventHandlers[m.Data.Topic].(func(PubSubResponse)); ok {
			go handlerType(f)(m)
		} else {
			if m.Data.Topic != "" {
				fmt.Printf("not running func topic=%s interface=%v\n", m.Data.Topic, eventHandlers[m.Data.Topic])
			} else {
				fmt.Println(m)
			}
		}
	}
}

// AddTopicHandler adds a function handle for PubSub topics
func (*Bot) AddTopicHandler(topic string, q interface{}) {
	if eventHandlers[topic] == nil {
		eventHandlers[topic] = q
	}
}

func (bot *Bot) wsSend(msg string) {
	bot.wsMutex.Lock()
	bot.ws.WriteMessage(ws.TextMessage, []byte(msg))
	bot.wsMutex.Unlock()
}

// Listen to a new pubsub topic(s)
func (bot *Bot) Listen(topics ...string) {
	sTopics := strings.Join(topics, `", "`)
	json := `{"type": "LISTEN","data": {"topics": ["` + sTopics + `"]}}`
	bot.wsSend(json)
}

// UnListen to a pubsub topic(s)
func (bot *Bot) UnListen(topics ...string) {
	sTopics := strings.Join(topics, `", "`)
	json := `{"type": "UNLISTEN","data": {"topics": ["` + sTopics + `"]}}`
	bot.wsSend(json)
}
