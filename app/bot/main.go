package main

import (
	"flag"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/hay-kot/yal"
	"github.com/joho/godotenv"
)

var (
	AppVersion     = "0.0.1"
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
		if !os.IsNotExist(err) {
			yal.Fatal("Error loading .env file")
		}
		yal.Warn("No .env file found")
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

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	app := NewApp(dg)
	return app.run()
}
