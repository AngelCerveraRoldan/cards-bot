package main

import (
	"fmt"
	"github.com/angelcerveraroldan/cards-bot/cmd/commands"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const prefix = "!h"

var discordToken string

func init() {
	// Load discord token from docker env variables
	discordToken = os.Getenv("TOKEN")

	if discordToken == "" {
		panic("Discord token needed")
	}
}

func main() {
	ds, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("Error creating discord session, err: ", err)
	}

	// This function will be run every time a message is sent into any channel that the commands can read
	ds.AddHandler(messageHandler)

	err = ds.Open()
	if err != nil {
		fmt.Println("Error opening discord commands: ", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	ds.Close()
}

// messageHandler
//
// This function will be run when any message is sent.
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	messages := strings.Fields(m.Message.Content)

	// Ignore messages sent by a bot
	if m.Author.Bot {
		return
	}

	// Ignore any message that wasn't intended for the bot
	if messages[0] == prefix {
		commands.RunCommand(messages[1:], s, m)
	}
}
