package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/hay-kot/yal"
)

type DiscordHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type command struct {
	name string
	desc string
	fn   DiscordHandler
}

func (c *command) ToAppCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.name,
		Description: c.desc,
	}
}

type app struct {
	s        *discordgo.Session
	commands []command
}

// NewApp constructs a new instance of the app struct and returns a pointer to it
// with the slash commands added to the app instance
func NewApp(s *discordgo.Session) *app {
	a := app{
		s: s,
	}

	a.commands = []command{
		{
			name: "mealie-default-credentials",
			desc: "Show default Credentials",
			fn:   a.HandlerStaticMessage(DefaultCredentials),
		},
		{
			name: "mealie-migration-links",
			desc: "Show helpful migration links",
			fn:   a.HandlerStaticMessage(V1MigrationLinks),
		},
		{
			name: "mealie-docker-faq",
			desc: "Show helpful Docker Problems and Solutions",
			fn:   a.HandlerStaticMessage(DockerFAQ),
		},
		{
			name: "mealie-docker-tags",
			desc: "Show helpful Docker Tags",
			fn:   a.HandlerStaticMessage(DockerTags),
		},
		{
			name: "mealie-token-time",
			desc: "Show information about the TOKEN_TIME variable",
			fn:   a.HandlerStaticMessage(TokenTime),
		},
		{
			name: "mealie-bot-debug",
			desc: "Show helpful debugging information",
			fn:   a.HandlerStaticMessage(fmt.Sprintf(BotDebug, AppVersion)),
		},
		{
			name: "mealie-feature-request",
			desc: "Show information on how to request a feature",
			fn:   a.HandlerStaticMessage(FeatureRequest),
		},
	}

	return &a
}

func (a *app) run() error {
	// ======================================================
	// Add Handlers to Discord Bot
	hdlrs := a.AppHandlers()

	a.s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := hdlrs[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// ======================================================
	// Start Session
	// Open a websocket connection to Discord and begin listening.
	err := a.s.Open()
	if err != nil {
		return err
	}

	// ======================================================
	// Add Commands to Bot
	yal.Info("Adding commands...")
	commands := a.AppCommands()

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := a.s.ApplicationCommandCreate(a.s.State.User.ID, *GuildID, v)
		if err != nil {
			yal.Fatalf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// Ensure that the commands are removed on shutdown
	defer func() {
		if !(*RemoveCommands) {
			yal.Info("Not removing commands...")
			return
		}

		yal.Info("Removing commands...")
		for _, v := range registeredCommands {
			err := a.s.ApplicationCommandDelete(a.s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				yal.Fatalf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}()

	// Wait here until CTRL-C or other term signal is received.
	yal.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	return a.s.Close()
}

func (a *app) AppCommands() []*discordgo.ApplicationCommand {
	results := make([]*discordgo.ApplicationCommand, len(a.commands))
	for i, cmd := range a.commands {
		results[i] = cmd.ToAppCommand()
	}
	return results
}

func (a *app) AppHandlers() map[string]DiscordHandler {
	results := make(map[string]DiscordHandler, len(a.commands))
	for _, cmd := range a.commands {
		results[cmd.name] = cmd.fn
	}
	return results
}

// HandleStaticMessage is a function that will construct a message handler that will write
// a static message. Useful for creating simple command that only produce static content
func (a *app) HandlerStaticMessage(msg string) DiscordHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		yal.Info("Sending message...")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: WrapMessage(msg),
			},
		})
	}
}
