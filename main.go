package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	guildID, exists := os.LookupEnv("GUILD_ID")
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

	log.Println("opening connection...")

	if err = discord.Open(); err != nil {
		log.Panicln(err)
	}

	discord.AddHandler(handleStandup)

	standupCommand := discordgo.ApplicationCommand{
		ID:                       "",
		ApplicationID:            "",
		GuildID:                  guildID,
		Version:                  "",
		Type:                     0,
		Name:                     "standup",
		NameLocalizations:        &map[discordgo.Locale]string{},
		DefaultPermission:        new(bool),
		DefaultMemberPermissions: new(int64),
		DMPermission:             new(bool),
		NSFW:                     new(bool),
		Description:              "what are you upto today!?",
		DescriptionLocalizations: &map[discordgo.Locale]string{},
		Options:                  []*discordgo.ApplicationCommandOption{},
	}

	ccmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, guildID, &standupCommand)
	if err != nil {
		log.Panicln(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Removing commands...")

	if err := discord.ApplicationCommandDelete(discord.State.User.ID, guildID, standupCommand.ID); err != nil {
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
			Timestamp:   "",
			Color:       0,
			Footer: &discordgo.MessageEmbedFooter{
				Text:         "",
				IconURL:      "",
				ProxyIconURL: "",
			},
			Image: &discordgo.MessageEmbedImage{
				URL:      "",
				ProxyURL: "",
				Width:    0,
				Height:   0,
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL:      "",
				ProxyURL: "",
				Width:    0,
				Height:   0,
			},
			Video: nil,
			Provider: &discordgo.MessageEmbedProvider{
				URL:  "",
				Name: "",
			},
			Author: &discordgo.MessageEmbedAuthor{
				URL:          "https://amir.day",
				Name:         "Amir Azizafshari",
				IconURL:      "https://picsum.photos/200",
				ProxyIconURL: "",
			},
			Fields: []*discordgo.MessageEmbedField{{Name: "goombi field", Value: "goombi value", Inline: false}},
		}

		if _, err := discord.ChannelMessageSendEmbed("1217489307431075991", &embed); err != nil {
			panic(err)
		}
	}
}
