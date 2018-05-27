package main

import (
	"os"

	"github.com/TimothyCole/WeetBot-V2/helper"
	"github.com/TimothyCole/WeetBot-V2/twitch"
	_ "github.com/joho/godotenv/autoload"
)

var (
	nick = os.Getenv("TWITCH_BOT_NAME")
	pass = os.Getenv("TWITCH_BOT_AUTH")
)

func main() {
	bot := twitch.NewClient(nick, pass)
	if err := bot.Connect(); err != nil {
		panic(err)
	}

	db, err := helper.InitDB()
	if err != nil {
		panic(err)
	}

	results, err := db.Query("SELECT id, login, display_name, followers, views, type, `join` FROM users WHERE `join`=1")
	if err != nil {
		panic(err)
	}

	var users []User
	for results.Next() {
		var user User
		results.Scan(&user.ID, &user.Login, &user.DisplayName, &user.Followers, &user.Views, &user.Type, &user.Join)
		go joinChannel(bot, user.ID, user.Login)
		users = append(users, user)
	}

	<-bot.Done
}
