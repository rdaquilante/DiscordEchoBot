package bot

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

//Main gon do shit in here yo
func Main(user, pass string) {
	dproc, err := discordgo.New(user, pass)
	if err != nil {
		fmt.Println("Bot failed to start:", err)
		return
	}
	if dproc == nil {
		return
	}
	dproc.AddHandler(messageCreate)

	//this opens the websocket. Hell if I know what that means.
	dproc.Open()
	defer dproc.Close()
	//get into the selected voice chat
	if channel, err := dproc.Channel("163811859596705792"); err != nil {
		fmt.Println(err)
		return
	} else if connErr := dproc.ChannelVoiceJoin(channel.GuildID, "163811898419052544", false, false); connErr != nil {
		fmt.Println("connErr: ", connErr)
		return
	}
	defer dproc.Voice.Close()
	for !dproc.Voice.Ready {
		time.After(1000)
	}
	if err := dproc.Voice.Open(); err != nil {
		fmt.Println(err)
		return
	}
	//echo()
	//keeps the program running until I tell it to stop
	var input string
	//read until read error or "end" token
	for true {
		if _, err := fmt.Scanln(&input); err == nil {
			if input == "end" {
				//close up shop we're done here
				break
			}
		}
	}
	return
}

//alright here's the meat
func echo(sess *discordgo.Session, vserv *discordgo.VoiceStateUpdate) {
	recv := make(chan *discordgo.Packet, 2)
	go dgvoice.ReceivePCM(v, recv)

	send := make(chan []int16, 2)
	go dgvoice.SendPCM(v, send)

	v.Speaking(true)
	defer v.Speaking(false)

	for {

		p, ok := <-recv
		if !ok {
			return
		}

		send <- p.PCM
	}
}

//Copying this from the tutorial, it'll be called every time a new message is
//created on a channel this user has access to
func messageCreate(sess *discordgo.Session, mess *discordgo.MessageCreate) {
	switch {
	case mess.Content == "!script":
		//print who requested the script
		fmt.Printf("%20s %20s %20s > %s\n", mess.ChannelID, time.Now().Format(time.Stamp), mess.Author.Username, mess.Content)
		r, _ := os.Open("H:/Dropbox/Dropbox/Abridged Emblem.txt")
		sess.ChannelFileSend(mess.ChannelID, "Abridged Emblem.txt", r)
	}
}
