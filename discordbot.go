package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
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

	if strings.HasPrefix(m.Content, "!random") {
		randomPlayer(session, m)
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

	session.ChannelMessageSend(m.ChannelID, "Hello "+m.Author.Username)
}

func answerBgg(session *discordgo.Session, m *discordgo.MessageCreate) {

	// https://boardgamegeek.com/wiki/page/BGG_XML_API2

	parts := strings.Split(m.Content, " ")
	if len(parts) < 2 {
		session.ChannelMessageSend(m.ChannelID, "Should have had three parts "+m.Author.Username)
		return
	}

	switch parts[1] {
	case "find":
		session.ChannelMessageSend(m.ChannelID, "Finding "+m.Author.Username)
	case "search":
		searchstring := strings.Join(parts[2:], "+")
		exact := true
		results, searchurl := BGGSearchItems(searchstring, "boardgame", exact)
		switch results.Total {
		case "0":
			session.ChannelMessageSend(m.ChannelID, "No results found.")
		case "1":
			gameID := results.Items[0].ID
			meta := BGGGetItemPage("https://boardgamegeek.com/boardgame/" + gameID)
			//thumbnail := BGGGetThumbnail(meta.Image)
			complexMessage := discordgo.MessageEmbed{

				Title:       results.Items[0].Names[0].Value,
				Description: meta.Description,
				URL:         "https://boardgamegeek.com/boardgame/" + gameID,
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: meta.Image},
			}

			session.ChannelMessageSendEmbed(m.ChannelID, &complexMessage)

		case "2":
			session.ChannelMessageSend(m.ChannelID, "Found "+results.Total+" results.")
		default:
			//session.ChannelMessageSend(m.ChannelID, "Found "+results.Total+" results. Please refine the search.")
			//gameID := results.Items[0].ID
			complexMessage := discordgo.MessageEmbed{

				Title: "Found " + results.Total + " results.",
				URL:   searchurl,
			}
			session.ChannelMessageSendEmbed(m.ChannelID, &complexMessage)
		}

	case "exact":
		searchstring := strings.Join(parts[2:], "+")
		exact := true
		results, _ := BGGSearchItems(searchstring, "boardgame", exact)
		gameID := results.Items[0].ID

		meta := BGGGetItemPage("https://boardgamegeek.com/boardgame/" + gameID)
		complexMessage := discordgo.MessageEmbed{

			Title:       results.Items[0].Names[0].Value,
			Description: meta.Description,
			URL:         "https://boardgamegeek.com/boardgame/" + gameID,
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: meta.Image},
		}
		session.ChannelMessageSendEmbed(m.ChannelID, &complexMessage)
	default:
		session.ChannelMessageSend(m.ChannelID, "Hello "+m.Author.Username)
	}
}

func randomPlayer(session *discordgo.Session, m *discordgo.MessageCreate) {

	parts := strings.Split(m.Content, " ")
	if len(parts) != 2 {
		session.ChannelMessageSend(m.ChannelID, "Should have a single number e.g. !random 3 "+m.Author.Username)
		return
	}
	max, _ := strconv.Atoi(parts[1])
	// Choose random number between 1 and max
	randomNumber := rand.Intn(max-1) + 1
	log.Println("The secret number is", randomNumber)
	session.ChannelMessageSend(m.ChannelID, "Random Number: "+strconv.Itoa(randomNumber))
}

func sendMessageEmbed() {

}
