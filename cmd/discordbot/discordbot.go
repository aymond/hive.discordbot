package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aymond/hive.discordbot/pkg/cmd"
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

	if err := cmd.RunBot(token, gamestatus); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
