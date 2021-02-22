package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/aymond/hive.discordbot/internal/pkg/bgg"
	"github.com/bwmarrin/discordgo"
	"github.com/dyatlov/go-htmlinfo/htmlinfo"
)

var gamestatus string

// TODO Make bot command prefix setting configurable i.e. default "!", but may be overridden..

// MessageCreate is called every time a new message is created on any channel that the authenticated bot has access to.
func MessageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages that the bot itself creates
	if m.Author.ID == session.State.User.ID {
		return
	}

	// Check just the first character for the bot command prefix. If bot command prefix is used, then check which command.
	if strings.HasPrefix(m.Content, "!") {

		log.Println("Bot Command:", m.Content)

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
			message := randomPlayer(session, m)
			session.ChannelMessageSend(m.ChannelID, message)
		}
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

	// TODO Normalise the BGG Commands to an intent e.g. find, get, should map to search action
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
			u := "https://boardgamegeek.com/boardgame/" + gameID
			resp, err := http.Get(u)
			if err != nil {
				panic(err)
			}

			defer resp.Body.Close()

			info := htmlinfo.NewHTMLInfo()
			// if url can be nil too, just then we won't be able to fetch (and generate) oembed information
			err = info.Parse(resp.Body, &u, nil)

			if err != nil {
				panic(err)
			}

			complexMessage := discordgo.MessageEmbed{

				Title:       info.OGInfo.Title,
				Description: info.OGInfo.Description,
				URL:         info.OGInfo.URL,
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: info.OGInfo.Images[0].URL},
			}

			channelMessageSendEmbed(session, m.ChannelID, &complexMessage)

		case "2":
			session.ChannelMessageSend(m.ChannelID, "Found "+results.Total+" results.")
		default:
			complexMessage := discordgo.MessageEmbed{

				Title: "Found " + results.Total + " results.",
				URL:   searchurl,
			}
			session.ChannelMessageSendEmbed(m.ChannelID, &complexMessage)
		}
	default:
		session.ChannelMessageSend(m.ChannelID, "Hello "+m.Author.Username)
	}
}

// RandomPlayer generates a random number
func randomPlayer(session *discordgo.Session, m *discordgo.MessageCreate) string {
	var message string
	parts := strings.Split(m.Content, " ")
	if len(parts) != 2 {
		message = (m.Author.Username + ", you should have a single number e.g. !random 3 ")
	} else {
		max, _ := strconv.Atoi(parts[1])
		// Choose random number between 1 and the given range
		randomNumber := rand.Intn(max-1) + 1
		message = m.Author.Username + "'s random number: " + strconv.Itoa(randomNumber)
	}
	return message
}

func channelMessageSend(session *discordgo.Session, m *discordgo.MessageCreate, message string) {
	session.ChannelMessageSend(m.ChannelID, message)
}

func channelMessageSendEmbed(session *discordgo.Session, m string, message *discordgo.MessageEmbed) {
	response, err := session.ChannelMessageSendEmbed(m, message)

	if err != nil {
		fmt.Println("Error sending embedded message:", response, err)
	}
}
