package cmd

import (
	"log"

	"github.com/FrisovanderVeen/mobot/bot/bot"
	"github.com/urfave/cli"
)

var version = "0.1"
var helpTemplate = `NAME:
{{.Name}} - {{.Usage}}
DESCRIPTION:
{{.Description}}
USAGE:
{{.Name}} {{if .Flags}}[flags] {{end}}command{{if .Flags}}{{end}} [arguments...]
COMMANDS:
	{{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
	{{end}}{{if .Flags}}
FLAGS:
	{{range .Flags}}{{.}}
	{{end}}{{end}}
VERSION:
` + version +
	`{{ "\n"}}`

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "config, conf",
		Value: "config.toml",
		Usage: "The TOML settings file",
	},
}

// Cmd is a command-line application
type Cmd struct {
	*cli.App
}

// NewApp creates a new client
func NewApp() *Cmd {
	app := cli.NewApp()
	app.Name = "Crep Bot"
	app.Author = ""
	app.Usage = "Crep Bot"
	app.Description = "Another discord bot"
	app.Flags = globalFlags

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.Action = func(c *cli.Context) error {
		conf := c.String("config")

		bot1, err := bot.NewBot(conf)
		if err != nil {
			log.Printf("Could not create bot: %v", err)
			return err
		}
		if err := bot1.Run(); err != nil {
			log.Fatalf("Could not run bot: %v", err)
			return err
		}

		return nil
	}

	return &Cmd{App: app}
}
