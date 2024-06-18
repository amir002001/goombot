package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	botToken, exists := os.LookupEnv("BOT_AUTH_TOKEN")
	if !exists {
		log.Fatalln("BOT_AUTH_TOKEN missing from environment")
	}
	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalln(err)
	}

	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentGuildMessages

	err = discord.Open()

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("opening connection...")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			log.Fatalln(err)
		}
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Ping!")
		if err != nil {
			log.Fatalln(err)
		}
	}
}
