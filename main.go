package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

var config ConfigSchemaJson

var standUpOptions = struct {
	Plans    string
	Blockers string
	ShoutOut string
	FunFact  string
}{
	Plans:    "plans",
	Blockers: "blockers",
	ShoutOut: "shoutout",
	FunFact:  "funfact",
}

func init() {
	configTmp, err := createConfig()
	if err != nil {
		log.Panicln(err)
	}

	config = *configTmp
}

func main() {
	// TODO clean every time errors
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

		embed, err := createEmbed(goombi, commandData.Options)
		if err != nil {
			log.Panicln(err)
		}

		if _, err := discord.ChannelMessageSendEmbed(config.Goombot.StandupChannelId, embed); err != nil {
			log.Panicln(err)
		}
	}
}

func createEmbed(goombi ConfigSchemaJsonGoombisElem, commandDataOptions []*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.MessageEmbed, error) {
	fields := []*discordgo.MessageEmbedField{}

	for _, option := range commandDataOptions {
		stringBuilder := strings.Builder{}
		field := discordgo.MessageEmbedField{}

		switch option.Name {
		case standUpOptions.Plans:
			field.Name = "Plans"
			plans := strings.Split(option.StringValue(), ";")

			for _, plan := range plans {
				stringBuilder.WriteString(fmt.Sprintf("* %s\n", plan))
			}

			field.Value = stringBuilder.String()

		case standUpOptions.Blockers:
			field.Name = "Blockers"
			blockers := strings.Split(option.StringValue(), ";")

			for _, blocker := range blockers {
				stringBuilder.WriteString(fmt.Sprintf("* %s\n", blocker))
			}

			field.Value = stringBuilder.String()

		case standUpOptions.FunFact:
			field.Name = "Fun Fact"
			field.Value = option.StringValue()
		case standUpOptions.ShoutOut:
			field.Name = "Shoutout"
			field.Value = option.StringValue()
		default:
			return nil, fmt.Errorf("cannot find valid option name %s", option.Name)
		}

		fields = append(fields, &field)
	}

	colorStr := strings.TrimPrefix(goombi.EmbedColor, "#")

	colorInt, err := strconv.ParseInt(colorStr, 16, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse goombi color: %w", err)
	}

	embed := discordgo.MessageEmbed{
		Title: "Daily Standup",
		Type:  "rich",
		Color: int(colorInt),
		Author: &discordgo.MessageEmbedAuthor{
			URL:     goombi.Url,
			Name:    goombi.Name,
			IconURL: goombi.ThumbnailUrl,
		},
		Fields: fields,
	}

	return &embed, nil
}

func createStandupCommand() discordgo.ApplicationCommand {
	standupCommand := discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "standup",
		Description: "what are you up to today!?",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        standUpOptions.Plans,
				Description: "semicolon (;) separated list of things you want to do",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        standUpOptions.Blockers,
				Description: "semicolon (;) separated list of things blocking you",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        standUpOptions.ShoutOut,
				Description: "give someone a shoutout",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        standUpOptions.FunFact,
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
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	byteValue, err := io.ReadAll(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("failed reading byte value: %w", err)
	}

	yamlFile.Close()

	var config ConfigSchemaJson
	if err := yaml.Unmarshal(byteValue, &config); err != nil {
		log.Fatalln(err)
	}

	return &config, nil
}
