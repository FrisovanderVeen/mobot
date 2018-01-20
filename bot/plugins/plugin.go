package plugins

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

var (
	// Plugins is a map of all plugins for the bot
	Plugins = make(map[string]Plugin)
	// Prefix is the prefix of all commands for the bot
	Prefix string

	// Channels for information and errors
	InfChan   = make(chan string)
	WarnChan  = make(chan error)
	ErrChan   = make(chan error)
	FatalChan = make(chan error)
)

// Plugin is a wrapper for a discordgo handler
type Plugin struct {
	Action interface{}

	Help Help
}

// Help is a collection of some useful fields for users
type Help struct {
	// A map of commands and their uses
	Commands map[string]string

	// If true the commands defined will be showed when help is called
	View bool

	// A explanation of the plugin
	Explanation string
}

// Register registers the plugin with the bot
func Register(name string, help Help, action interface{}) interface{} {
	switch action := action.(type) {
	case func(*discordgo.Session, *discordgo.MessageCreate), func(*discordgo.Session, *discordgo.Ready):
		for plugname, _ := range Plugins {
			if plugname == name {
				color.Yellow("[WARNING]: 2 or more plugins with the same name detected this may cause unwanted effects: %s", name)
			}
		}
		Plugins[name] = Plugin{
			Action: action,
			Help:   help,
		}
	default:
		log.Fatalf("Unknown plugin type: %T", action)
		return action
	}
	return action
}
