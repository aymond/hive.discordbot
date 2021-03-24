package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"html"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/aymond/hive.discordbot/internal/pkg/bgg"
	"github.com/bwmarrin/discordgo"
	"github.com/dyatlov/go-htmlinfo/htmlinfo"
)

// TODO Make bot command prefix setting configurable i.e. default "!", but may be overridden..

// MessageCreate is called every time a new message is created on any channel that the authenticated bot has access to.
func MessageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages that the bot itself creates
	if m.Author.ID == session.State.User.ID {
		return
	}

	// Load config defaults each time there is a bot command.
	viper.SetDefault("hello", "Hi!")
	viper.SetDefault("ping", "pong!")
	// Load the standard responses mounted at /data/discordbot-config/
	viper.SetConfigName("config.yaml")              // name of config file (without extension)
	viper.SetConfigType("yaml")                     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/data/discordbot-config/") // path to look for the config file in
	viper.AddConfigPath("./configs/")               // optionally look for config in the working directory
	err := viper.ReadInConfig()                     // Find and read the config file
	if err != nil {                                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	log.Printf(viper.GetString("version"))
	viper.SetConfigName("version.yaml")
	viper.AddConfigPath("/data/discordbot-config/") // path to look for the config file in
	viper.AddConfigPath("./configs/")
	err = viper.MergeInConfig() // Merge in the additional version config
	log.Printf(viper.GetString("version"))

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
			message := randomPlayer(m)
			ChannelMessageSend(session, m.ChannelID, message)
		}
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		ChannelMessageSend(session, m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		ChannelMessageSend(session, m.ChannelID, "Ping!")
	}
}

// AnswerHello answers the Hello command
func AnswerHello(session *discordgo.Session, m *discordgo.MessageCreate) {

	response, err := session.ChannelMessageSend(m.ChannelID, "Hello "+m.Author.Username)
	if err != nil {
		log.Println("Error sending embedded message:", response, err)
	}
}

// AnswerBgg answers the BGG command
func AnswerBgg(session *discordgo.Session, m *discordgo.MessageCreate) {

	parts := strings.Split(m.Content, " ")
	if len(parts) < 2 {
		ChannelMessageSend(session, m.ChannelID, "Should have had three parts "+m.Author.Username)
		return
	}

	// TODO Normalise the BGG Commands to an intent e.g. find, get, should map to search action
	switch parts[1] {
	case "find":
		ChannelMessageSend(session, m.ChannelID, "Finding "+m.Author.Username)
	case "search":
		searchString := strings.Join(parts[2:], "+")
		results, searchURL := bgg.SearchItems(searchString, "boardgame", true)

		resultsCount, _ := strconv.Atoi(results.Total)

		// If no results, then try again with exact set to false.
		if resultsCount == 0 {
			log.Println("No results found with exact search for " + searchString + ". Using non-exact search.")
			results, searchURL = bgg.SearchItems(searchString, "boardgame", false)
			log.Printf(searchURL)
			resultsCount, _ = strconv.Atoi(results.Total)
		}

		if resultsCount == 0 {
			ChannelMessageSend(session, m.ChannelID, "No results found.")
		}

		if resultsCount >= 1 && resultsCount <= 6 {
			ChannelMessageSend(session, m.ChannelID, "Found "+results.Total+" results.")
			for i, v := range results.Items {
				log.Printf(strconv.Itoa(i))
				gameID := v.ID
				u := "https://boardgamegeek.com/boardgame/" + gameID
				resp, err := http.Get(u)
				if err != nil {
					panic(err)
				}

				info := htmlinfo.NewHTMLInfo()
				// if url can be nil too, just then we won't be able to fetch (and generate) oembed information
				err = info.Parse(resp.Body, &u, nil)

				if err != nil {
					panic(err)
				}

				err = resp.Body.Close()
				if err != nil {
					log.Println("Error retrieving boardgame by gameID:", err)
				}

				description := SplitLines(info.OGInfo.Description)
				descriptionEscaped := html.UnescapeString(description[0])

				complexMessage := discordgo.MessageEmbed{

					Title:       info.OGInfo.Title,
					Description: descriptionEscaped,
					URL:         info.OGInfo.URL,
					Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: info.OGInfo.Images[0].URL},
				}

				channelMessageSendEmbed(session, m.ChannelID, &complexMessage)
			}
		}

		if resultsCount > 6 {
			ChannelMessageSend(session, m.ChannelID, "Found "+results.Total+" results.")
		}

	default:
		ChannelMessageSend(session, m.ChannelID, "Hello "+m.Author.Username)
	}
}

// RandomPlayer generates a random number
func randomPlayer(m *discordgo.MessageCreate) string {
	var message string
	parts := strings.Split(m.Content, " ")
	if len(parts) != 2 {
		message = m.Author.Username + ", you should have a single number e.g. !random 3 "
	} else {
		max, _ := strconv.Atoi(parts[1])
		// Choose random number between 1 and the given range
		randomNumber := rand.Intn(max-1) + 1
		message = m.Author.Username + "'s random number: " + strconv.Itoa(randomNumber)
	}
	return message
}

func ChannelMessageSend(session *discordgo.Session, channelID string, message string) {
	response, err := session.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Println("Error sending embedded message:", response, err)
	}
}

func channelMessageSendEmbed(session *discordgo.Session, channelID string, message *discordgo.MessageEmbed) {
	response, err := session.ChannelMessageSendEmbed(channelID, message)

	if err != nil {
		log.Println("Error sending embedded message:", response, err)
	}
}

// SplitLines splits lines
func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}
