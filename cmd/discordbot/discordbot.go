package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"github.com/gorilla/mux"

	"github.com/aymond/hive.discordbot/pkg/cmd"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	token      string
	gamestatus string
	results    []string
	Session    *discordgo.Session
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/post", postHandler)
	r.HandleFunc("/messagepost", messagepostHandler)

	log.Printf("Post Handler listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

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
	Session = session
}

func messageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {
	cmd.MessageCreate(session, m)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var dir string = "../../web/static/"
	t, _ := template.ParseFiles(dir + "index.html") //parse the html file index.html

	t.Execute(w, "") //execute the template
}

func messagepostHandler(w http.ResponseWriter, r *http.Request) {
	//curl -d "channel=1234&message=Hello,%20World%21" -X POST http://localhost:3000/messagepost
	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostForm = %v\n", r.PostForm)
		cmd.ChannelMessageSend(Session, r.PostFormValue("channel"), r.PostFormValue("message"))
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")

	}
}

// postHandler converts post request body to string
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		results = append(results, string(body))

		log.Println(w, "POST done")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
