package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aymond/hive.discordbot/pkg/cmd"
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

func ready(session *discordgo.Session, event *discordgo.Ready) {
	// Set the bots status.
	session.UpdateGameStatus(0, gamestatus)
}

func messageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {
	cmd.MessageCreate(session, m)
}
