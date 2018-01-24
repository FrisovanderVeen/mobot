package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FrisovanderVeen/mobot/bot/config"
	"github.com/FrisovanderVeen/mobot/bot/plugins"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"

	// the airhorn plugin being registered
	_ "github.com/FrisovanderVeen/mobot/bot/plugins/airhorn"
	// The help plugin being registered
	_ "github.com/FrisovanderVeen/mobot/bot/plugins/help"
	// The onready plugin being registered
	_ "github.com/FrisovanderVeen/mobot/bot/plugins/onready"
	// The pingpong plugin being registered
	_ "github.com/FrisovanderVeen/mobot/bot/plugins/pingpong"
	// The youtube plugin being registered
	_ "github.com/FrisovanderVeen/mobot/bot/plugins/youtube"
	// The audio plugin being registered
	_ "github.com/FrisovanderVeen/mobot/bot/plugins/audio"
)

// Bot is a wrapper for a discordgo session
type Bot struct {
	Session *discordgo.Session
	Prefix  string

	Exit chan error
}

// NewBot creates a new bot based on the settings in the configuration file
func NewBot(confloc string) (*Bot, error) {
	conf, err := config.GetConfig(confloc)
	if err != nil {
		return nil, err
	}

	token := conf.Discord.Token
	if token == "" {
		return nil, fmt.Errorf("No token specified")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("Could not create session: %v", err)
	}

	bot := &Bot{
		Session: dg,
		Prefix:  conf.Discord.Prefix,
		Exit:    make(chan error),
	}

	go func() {
		c := color.New(color.FgRed, color.Bold)
		for {
			select {
			case inf := <-plugins.InfChan:
				color.Green("[INFO]: %s", inf)
			case err := <-plugins.WarnChan:
				color.Yellow("[WARNING]: %v", err)
			case err := <-plugins.ErrChan:
				color.Red("[ERROR]: %v", err)
			case err := <-plugins.FatalChan:
				c.Printf("[FATAL]: %v\n", err)
				bot.Exit <- err
			}
		}
	}()

	for name, opts := range conf.Plugins {
		addPlugins(name, opts, bot, conf)
	}

	plugins.Prefix = conf.Discord.Prefix
	plugins.Config = conf

	return bot, nil
}

func addPlugins(name string, opts config.Plugin, bot *Bot, conf *config.TomlConfig) {
	for _, plugin := range plugins.Plugins {
		if name == plugin.Name {
			if opts.Enabled {
				bot.Session.AddHandler(plugin.Action)
				plugin.Enabled = true
				return
			} else if !opts.Enabled {
				return
			}
		}
	}

	if opts.Necessary {
		plugins.FatalChan <- fmt.Errorf("Plugin necessary but not present: %s", name)
	} else {
		plugins.WarnChan <- fmt.Errorf("Plugin specified in config but not present: %s", name)
	}
}

// Run runs the bot and exits if CTRL-C is pressed or if there is a fatal error
func (b *Bot) Run() error {
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("Could not open session: %v", err)
	}
	defer func() {
		if err := b.Session.Close(); err != nil {
			log.Printf("Could not close session: %v\n", err)
		}
	}()

	plugins.InfChan <- "Bot is now running.  Press CTRL-C to exit."
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	select {
	case <-sc:
	case err := <-b.Exit:
		return err
	}

	return nil
}
