package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/TimothyCole/WeetBot-V2/twitch"
	_ "github.com/joho/godotenv/autoload"
)

var (
	nick = os.Getenv("TWITCH_BOT_NAME")
	pass = os.Getenv("TWITCH_BOT_AUTH")
	api  CommandsAPI
)

type CommandsAPI struct {
	Data []Command `json:"data"`
}

type Command struct {
	Channel   int    `json:"channel"`
	Command   string `json:"command"`
	Response  string `json:"response"`
	Cooldown  int    `json:"cooldown"`
	Userlevel int    `json:"userlevel"`
	Points    int    `json:"points"`
	Hidden    bool   `json:"hidden"`
}

func GetCommandsForTim() {
	commands := "https://tcole.me/api/stream/51684790/commands"
	req, _ := http.NewRequest("GET", commands, nil)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	api = CommandsAPI{}
	err = json.Unmarshal(body, &api)

	api.Data = append(api.Data, Command{
		Channel:   51684790,
		Command:   "!ping",
		Response:  "PONG KappaPride",
		Cooldown:  0,
		Userlevel: 0,
		Points:    0,
	})

	fmt.Println("Commands Loaded for Tim")
}

func main() {
	bot := twitch.NewClient(nick, pass)
	if err := bot.Connect(); err != nil {
		panic(err)
	}

	go GetCommandsForTim()

	bot.OnNewMessage(func(msg *twitch.Message) {
		fmt.Println("> New Message: ", msg.Data.DisplayName, msg.Data.Message)

		if "!reload" == msg.Data.Message && msg.Data.DisplayName == "ModestTim" {
			go GetCommandsForTim()
		}

		for _, command := range api.Data {
			if command.Command != msg.Data.Message || command.Channel != msg.Data.StreamerID {
				continue
			}

			bot.Say(msg.Data.StreamerName, command.Response)
		}
	})
	bot.OnNewWhisper(func(msg *twitch.Message) {
		fmt.Println("> New Whisper: ", msg.Data.DisplayName, msg.Data.Message)
	})
	bot.OnNewSub(func(msg *twitch.Message) {
		fmt.Println("> New Sub: ", msg.Data.DisplayName, msg.Data.Sub.Plan, msg.Data.GiftSub.Login)

		// dbug, _ := json.Marshal(msg)
		// bot.Say(msg.Data.StreamerName, string(dbug))
	})
	bot.OnNewRaid(func(msg *twitch.Message) {
		fmt.Println("> New Raid: ", msg.Data.Raid.DisplayName, msg.Data.Raid.Viewers)
	})

	go joinChannel(bot, 51684790, "ModestTim")
	go joinChannel(bot, 88560344, "WeetBot")

	<-bot.Done
}
