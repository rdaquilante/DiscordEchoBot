package bot

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/writ"
)

//Config sets up writ library to handle commands for me and add a helpful print Help
var Config struct {
	Script bool `flag:"script" description:"Upload the script to this chat."`
	Help   bool `flag:"help, h" description:"Print the help and commands for this bot."`
}

//Main sets up the bot and gets it running
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

	//this opens the websocket. From tutorial
	dproc.Open()
	defer dproc.Close()

	var channel *discordgo.Channel
	var chanErr error
	//get into the selected voice chat
	if channel, chanErr = dproc.Channel("163811859596705792"); chanErr != nil {
		fmt.Println(chanErr)
		return
	}
	voice, _ := dproc.ChannelVoiceJoin(channel.GuildID, "163811898419052544", false, false)

	for !voice.Ready {
		time.After(1000)
	}
	defer voice.Close()
	//echo()
	//keeps the program running until I tell it to stop
	var input string
	//read until read error or "end" token
	for true {
		if _, err := fmt.Scanln(&input); err == nil {
			if input == "end" {
				//close up shop we're done here
				dproc.Logout()
				break
			}
		}
	}
	return
}

// //alright here's the meat
// func echo(sess *discordgo.Session, vserv *discordgo.VoiceStateUpdate) {
// 	recv := make(chan *discordgo.Packet, 2)
// 	go dgvoice.ReceivePCM(v, recv)
//
// 	send := make(chan []int16, 2)
// 	go dgvoice.SendPCM(v, send)
//
// 	v.Speaking(true)
// 	defer v.Speaking(false)
//
// 	for {
//
// 		p, ok := <-recv
// 		if !ok {
// 			return
// 		}
//
// 		send <- p.PCM
// 	}
// }

//Copying this from the tutorial, it'll be called every time a new message is
//created on a channel this user has access to
func messageCreate(sess *discordgo.Session, mess *discordgo.MessageCreate) {
	//ignore the bot's own messages
	if mess.Author.Username == "FEAEchoBot" {
		return
	}
	//record the command for bot's use
	fmt.Printf("%20s %20s %20s > %s\n", mess.ChannelID, time.Now().Format(time.Stamp),
		mess.Author.Username, mess.Content)

	cmd := writ.New("bot", &Config) //create the config object for writ
	cmd.Help.Usage = "Call the bot with --[command]. Your own mechanical Jeeves."

	cmd.Decode(strings.Split(mess.Content, " ")) //decode message
	fmt.Printf("%#v", Config)

	if Config.Script == true { //this will handle ! commands like most bots
		//print who requested the script
		r, _ := os.Open("/home/roby/Dropbox/Abridged Emblem.txt")
		sess.ChannelFileSend(mess.ChannelID, "Abridged Emblem.txt", r)
	}
	if Config.Help == true {
		//print help and exit gracefully
		var buf []byte
		w := bytes.NewBuffer(buf)
		cmd.WriteHelp(w) //write the help to a []byte, convert to string, send
		sess.ChannelMessageSend(mess.ChannelID, w.String())
	}
	Config.Script = false
	Config.Help = false
}
