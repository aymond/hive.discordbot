package cmd

import (
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/aymond/hive.discordbot/internal/pkg/bgg"
	"github.com/bwmarrin/discordgo"
)

var gamestatus string

// Ready will be called when the bot receives the "ready" event from Discord.
func Ready(session *discordgo.Session, event *discordgo.Ready) {
	// Set the bots status.
	session.UpdateGameStatus(0, gamestatus)
}

// MessageCreate will be called every time a new message is created on any channel that the authenticated bot has access to.
func MessageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages that the bot creates
	if m.Author.ID == session.State.User.ID {
		return
	}

	// check if the message is "!Hello"
	if strings.HasPrefix(m.Content, "!Hello") {
		AnswerHello(session, m)
	}

	if strings.HasPrefix(m.Content, "!hello") {
		AnswerHello(session, m)
	}

	if strings.HasPrefix(m.Content, "!bgg") {
		AnswerBgg(session, m)
	}

	if strings.HasPrefix(m.Content, "!random") {
		RandomPlayer(session, m)
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

// AnswerHello answers the Hello command
func AnswerHello(session *discordgo.Session, m *discordgo.MessageCreate) {

	session.ChannelMessageSend(m.ChannelID, "Hello "+m.Author.Username)
}

// AnswerBgg answers the BGG command
func AnswerBgg(session *discordgo.Session, m *discordgo.MessageCreate) {
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
		results, searchurl := bgg.SearchItems(searchstring, "boardgame", exact)
		switch results.Total {
		case "0":
			session.ChannelMessageSend(m.ChannelID, "No results found.")
		case "1":
			gameID := results.Items[0].ID
			meta := bgg.GetItemPage("https://boardgamegeek.com/boardgame/" + gameID)
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
			complexMessage := discordgo.MessageEmbed{

				Title: "Found " + results.Total + " results.",
				URL:   searchurl,
			}
			session.ChannelMessageSendEmbed(m.ChannelID, &complexMessage)
		}

	case "exact":
		searchstring := strings.Join(parts[2:], "+")
		exact := true
		results, _ := bgg.SearchItems(searchstring, "boardgame", exact)
		gameID := results.Items[0].ID

		meta := bgg.GetItemPage("https://boardgamegeek.com/boardgame/" + gameID)
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

// RandomPlayer generates a random number
func RandomPlayer(session *discordgo.Session, m *discordgo.MessageCreate) {

	parts := strings.Split(m.Content, " ")
	if len(parts) != 2 {
		session.ChannelMessageSend(m.ChannelID, m.Author.Username+", you should have a single number e.g. !random 3 ")
		return
	}
	max, _ := strconv.Atoi(parts[1])
	// Choose random number between 1 and max
	randomNumber := rand.Intn(max-1) + 1
	log.Println("The secret number is", randomNumber)
	session.ChannelMessageSend(m.ChannelID, m.Author.Username+"'s random number: "+strconv.Itoa(randomNumber))
}
