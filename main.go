package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/hay-kot/yal"
	"github.com/joho/godotenv"
)

// Bot parameters
var (
	LogLevel       = flag.String("log", "info", "log level (debug, info, warn, error)")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
)

func init() {
	flag.Parse()
}

func SetLogger(value string) {
	switch strings.ToLower(value) {
	case "debug":
		yal.Log.Level = yal.LevelDebug
	case "info":
		yal.Log.Level = yal.LevelInfo
	case "warn":
		yal.Log.Level = yal.LevelWarn
	case "error":
		yal.Log.Level = yal.LevelError
	default:
		yal.Log.Level = yal.LevelInfo
	}
}

func main() {
	flag.Parse()

	SetLogger(*LogLevel)

	err := godotenv.Load()
	if err != nil {
		yal.Fatal("Error loading .env file")
	}

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		yal.Fatal("Error creating Discord session", err)
	}

	err = run(discord)
	if err != nil {
		yal.Fatal("Error running bot", err)
	}
}

func run(dg *discordgo.Session) error {
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		yal.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// =====================================================
	// Setup The Bot
	var (
	// integerOptionMinValue          = 1.0
	// dmPermission                   = false
	// defaultMemberPermissions int64 = discordgo.PermissionManageServer
	)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// ======================================================
	// Mount '/' Commands

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "mealie-default-credentials",
			Description: "Show default Credentials",
		},
		{
			Name:        "mealie-migration-links",
			Description: "Show helpful migration links",
		},
		{
			Name:        "mealie-docker-faq",
			Description: "Show helpful Docker Problems and Solutions",
		},
		{
			Name:        "mealie-docker-tags",
			Description: "Show Docker Tag Description",
		},
		{
			Name:        "mealie-token-time",
			Description: "Show information about the TOKEN_TIME env variable",
		},
	}

	msgCommandFunc := func(result string) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: result,
				},
			})
		}
	}

	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"mealie-default-credentials": msgCommandFunc(msg(DefaultCredentials)),
		"mealie-migration-links":     msgCommandFunc(msg(V1MigrationLinks)),
		"mealie-docker-faq":          msgCommandFunc(msg(DockerFAQ)),
		"mealie-docker-tags":         msgCommandFunc(msg(DockerTags)),
		"mealie-token-time":          msgCommandFunc(msg(TokenTime)),
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// Open a websocket connection to Discord and begin listening.
	err := dg.Open()
	if err != nil {
		return err
	}

	// ======================================================
	// Add Commands to Bot
	yal.Info("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, *GuildID, v)
		if err != nil {
			yal.Fatalf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// Wait here until CTRL-C or other term signal is received.
	yal.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// ======================================================
	// Cleanup Commands

	if *RemoveCommands {
		yal.Info("Removing commands...")
		for _, v := range registeredCommands {
			err := dg.ApplicationCommandDelete(dg.State.User.ID, *GuildID, v.ID)
			if err != nil {
				yal.Fatalf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	return dg.Close()
}
