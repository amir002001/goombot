package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"slices"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

var config ConfigSchemaJson

func init() {
	configTmp, err := createConfig()
	if err != nil {
		log.Panicln(err)
	}

	config = *configTmp
}

func main() {
	discord, err := discordgo.New("Bot " + config.Goombot.BotAuthToken)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("opening connection...")

	if err = discord.Open(); err != nil {
		log.Panicln(err)
	}

	discord.AddHandler(handleStandup)

	standupCommand := createStandupCommand()

	ccmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, config.Goombot.GuildId, &standupCommand)
	if err != nil {
		log.Panicln(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	cleanUp(ccmd, config.Goombot.GuildId, discord)

	discord.Close()
}

func cleanUp(ccmd *discordgo.ApplicationCommand, guildID string, discord *discordgo.Session) {
	log.Println("Removing commands...")

	if err := discord.ApplicationCommandDelete(discord.State.User.ID, guildID, ccmd.ID); err != nil {
		log.Panicf("Cannot delete '%v' command: %v", ccmd.Name, err)
	}
}

func handleStandup(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {
	commandData := interaction.ApplicationCommandData()
	if commandData.Name == "standup" {
		goombiIdx := slices.IndexFunc(
			config.Goombis,
			func(c ConfigSchemaJsonGoombisElem) bool {
				return c.Id == interaction.Member.User.ID
			},
		)
		goombi := config.Goombis[goombiIdx]

		if interaction.Member.User.ID == config.Goombis[1].Id {
			log.Println("eureka")
		}

		embed := createEmbed(goombi, commandData.Options)

		if _, err := discord.ChannelMessageSendEmbed(config.Goombot.StandupChannelId, &embed); err != nil {
			panic(err)
		}
	}
}

func createEmbed(goombi ConfigSchemaJsonGoombisElem, commandData []*discordgo.ApplicationCommandInteractionDataOption) discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Title:       "Daily Standup",
		Type:        "rich",
		Description: "What I'm up to today",
		Color:       0,
		Author: &discordgo.MessageEmbedAuthor{
			URL:     goombi.Url,
			Name:    goombi.Name,
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
				Name:        "funfact",
				Description: "something you recently learned that's fun!",
				Required:    false,
			},
		},
	}

	return standupCommand
}

func createConfig() (*ConfigSchemaJson, error) {
	yamlFile, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}

	byteValue, err := io.ReadAll(yamlFile)
	if err != nil {
		return nil, err
	}

	yamlFile.Close()

	var config ConfigSchemaJson
	if err := yaml.Unmarshal(byteValue, &config); err != nil {
		log.Fatalln(err)
	}
	return &config, nil
}
