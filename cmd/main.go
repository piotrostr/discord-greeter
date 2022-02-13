package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/piotrostr/discord-greeter/pkg/bot"
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

	color.Blue("bot initialized with:")
	fmt.Printf("\n\t token: %s ", bot.Token)
	fmt.Printf("\n\t captcha key: %s", bot.CaptchaKey)
	fmt.Printf("\n\t invite code: %s", bot.InviteCode)
	fmt.Printf("\n\t proxy: %s\n", bot.Proxy)

	bot.JoinServer()
}
