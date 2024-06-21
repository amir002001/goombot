package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	guildId, exists := os.LookupEnv("GUILD_ID")
	if !exists {
		log.Fatalln("GUILD_ID missing from environment")
	}
	botToken, exists := os.LookupEnv("BOT_AUTH_TOKEN")
	if !exists {
		log.Fatalln("BOT_AUTH_TOKEN missing from environment")
	}
	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalln(err)
	}

	discord.AddHandler(handleStandup)

	standupCommand := discordgo.ApplicationCommand{
		Name:        "standup",
		Description: "what are you upto today!?",
	}
	ccmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, guildId, &standupCommand)

	if err != nil {
		log.Panicln(err)
	}

	discord.Identify.Intents = discordgo.IntentGuildMessages

	err = discord.Open()

	if err != nil {
		log.Panicln(err)
	}
	log.Println("opening connection...")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Removing commands...")

	if err := discord.ApplicationCommandDelete(discord.State.User.ID, guildId, standupCommand.ID); err != nil {
		log.Panicf("Cannot delete '%v' command: %v", ccmd.Name, err)
	}

	discord.Close()
}

func handleStandup(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "standup" {
		embed := discordgo.MessageEmbed{
			URL:         "https://goombi.com",
			Type:        "rich",
			Title:       "hello goombi",
			Description: "goombster",
			Timestamp:   time.Now().String(),
			Color:       0,
			Footer:      &discordgo.MessageEmbedFooter{},
			Image:       &discordgo.MessageEmbedImage{},
			Thumbnail:   &discordgo.MessageEmbedThumbnail{},
			Video:       &discordgo.MessageEmbedVideo{},
			Author: &discordgo.MessageEmbedAuthor{
				URL:     "https://amir.day",
				Name:    "Amir Azizafshari",
				IconURL: "https://picsum.photos/200",
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "goombi field",
					Value:  "goombi value",
					Inline: false,
				},
			},
		}

		if _, err := discord.ChannelMessageSendEmbed("", &embed); err != nil {
			log.Fatalf(err.Error())
		}
	}
}
