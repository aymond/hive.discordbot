package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	token      string
	gamestatus string
)

func init() {
	flag.StringVar(&token, "t", "", "Discord Bot Token")
	flag.StringVar(&gamestatus, "s", "Hacking!", "Game Status")
	flag.Parse()
}

func main() {
	if token == "" {
		fmt.Println("No token provided. Please run: with option -t <bot token>")
		return
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register callbacks for the discord events.
	discord.AddHandler(ready)
	discord.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Listener is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close the websocket and stop listening.
	discord.Close()
}

// ready will be called when the bot receives the "ready" event from Discord.
func ready(session *discordgo.Session, event *discordgo.Ready) {
	// Set the bots status.
	session.UpdateGameStatus(0, gamestatus)
}

// messageCreated will be called every time a new message is created on any channel that the authenticated bot has access to.
func messageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages that the bot creates
	if m.Author.ID == session.State.User.ID {
		return
	}

	// check if the message is "!Hello"
	if strings.HasPrefix(m.Content, "!Hello") {
		answerHello(session, m)
	}

	if strings.HasPrefix(m.Content, "!hello") {
		answerHello(session, m)
	}

	if strings.HasPrefix(m.Content, "!bgg") {
		answerBgg(session, m)
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		session.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		session.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func answerHello(session *discordgo.Session, m *discordgo.MessageCreate) {
	// Find the channel that the message came from.
	c, err := session.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	// Find the guild for that channel.
	_, err = session.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return
	}

	session.ChannelMessageSend(m.ChannelID, "Hello "+m.Author.Username)
}

func answerBgg(session *discordgo.Session, m *discordgo.MessageCreate) {

	// https://boardgamegeek.com/wiki/page/BGG_XML_API2

	// Find the channel that the message came from.
	c, err := session.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	// Find the guild for that channel.
	_, err = session.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return
	}

	parts := strings.Split(m.Content, " ")
	if len(parts) < 2 {
		session.ChannelMessageSend(m.ChannelID, "Should have had three parts "+m.Author.Username)
		return
	}

	switch parts[1] {
	case "find":
		session.ChannelMessageSend(m.ChannelID, "Finding "+m.Author.Username)
	case "search":
		results := SearchItems(parts[2], "boardgame", false)
		session.ChannelMessageSend(m.ChannelID, "Found "+results.Total+" results.")
	default:
		session.ChannelMessageSend(m.ChannelID, "Hello "+m.Author.Username)
	}
}
