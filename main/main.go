package main

import (
	"DiscordEchoBot/bot"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Need a user name and password to log in with!")
		return
	}
	bot.Main(os.Args[1], os.Args[2])
}
