package main

import (
	"github.com/piotrostr/discord-greeter/bot"
)

// initialize once per bot
// be careful not to initialize same bot with another ip
// its best to create token - proxy pairs and make sure ip is sticky

// proxy shall be in the form of
// username:password@host:port
// sample token
// OTQwNzg3MDY2MTU1OTYyNDA4.YgMetA.ILdepf9Gi1ehwCtzwEysLObpqbo

func main() {
	bot := bot.Bot{}
	bot.Initialize()
}
