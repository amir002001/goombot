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

	standupCommand := createStandupCommand()

	ccmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, guildID, &standupCommand)
	if err != nil {
		log.Panicln(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	cleanUp(ccmd, guildID, discord)

	discord.Close()
}

func cleanUp(ccmd *discordgo.ApplicationCommand, guildID string, discord *discordgo.Session) {
	log.Println("Removing commands...")

	if err := discord.ApplicationCommandDelete(discord.State.User.ID, guildID, ccmd.ID); err != nil {
		log.Panicf("Cannot delete '%v' command: %v", ccmd.Name, err)
	}
}

func handleStandup(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "standup" {
		embed := createEmbed()
		if _, err := discord.ChannelMessageSendEmbed("1217489307431075991", &embed); err != nil {
			panic(err)
		}
	}
}

func createEmbed() discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Title:       "hello goombi",
		Type:        "rich",
		Description: "goombster",
		Color:       0,
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://amir.day",
			Name:    "Amir Azizafshari",
			IconURL: "https://picsum.photos/200",
		},
		Fields: []*discordgo.MessageEmbedField{{Name: "goombi field", Value: "goombi value", Inline: false}},
	}

	return embed
}

func createStandupCommand() discordgo.ApplicationCommand {
	standupCommand := discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "standup",
		Description: "what are you upto today!?",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "plans",
				Description: "semicolon (;) separated list of things you want to do",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "blockers",
				Description: "semicolon (;) separated list of things blocking you",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "shoutout",
				Description: "give someone a shoutout",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "fun fact",
				Description: "something you recently learned that's fun!",
				Required:    false,
			},
		},
	}

	return standupCommand
}
