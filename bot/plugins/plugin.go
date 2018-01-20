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
	case func(*discordgo.Session, interface{}), func(*discordgo.Session, *discordgo.ChannelCreate), func(*discordgo.Session, *discordgo.ChannelDelete), func(*discordgo.Session, *discordgo.ChannelPinsUpdate), func(*discordgo.Session, *discordgo.ChannelUpdate), func(*discordgo.Session, *discordgo.Connect), func(*discordgo.Session, *discordgo.Disconnect), func(*discordgo.Session, *discordgo.Event), func(*discordgo.Session, *discordgo.GuildBanAdd), func(*discordgo.Session, *discordgo.GuildBanRemove), func(*discordgo.Session, *discordgo.GuildCreate), func(*discordgo.Session, *discordgo.GuildDelete), func(*discordgo.Session, *discordgo.GuildEmojisUpdate), func(*discordgo.Session, *discordgo.GuildIntegrationsUpdate), func(*discordgo.Session, *discordgo.GuildMemberAdd), func(*discordgo.Session, *discordgo.GuildMemberRemove), func(*discordgo.Session, *discordgo.GuildMemberUpdate), func(*discordgo.Session, *discordgo.GuildMembersChunk), func(*discordgo.Session, *discordgo.GuildRoleCreate), func(*discordgo.Session, *discordgo.GuildRoleDelete), func(*discordgo.Session, *discordgo.GuildRoleUpdate), func(*discordgo.Session, *discordgo.GuildUpdate), func(*discordgo.Session, *discordgo.MessageAck), func(*discordgo.Session, *discordgo.MessageCreate), func(*discordgo.Session, *discordgo.MessageDelete), func(*discordgo.Session, *discordgo.MessageDeleteBulk), func(*discordgo.Session, *discordgo.MessageReactionAdd), func(*discordgo.Session, *discordgo.MessageReactionRemove), func(*discordgo.Session, *discordgo.MessageReactionRemoveAll), func(*discordgo.Session, *discordgo.MessageUpdate), func(*discordgo.Session, *discordgo.PresenceUpdate), func(*discordgo.Session, *discordgo.PresencesReplace), func(*discordgo.Session, *discordgo.RateLimit), func(*discordgo.Session, *discordgo.Ready), func(*discordgo.Session, *discordgo.RelationshipAdd), func(*discordgo.Session, *discordgo.RelationshipRemove), func(*discordgo.Session, *discordgo.Resumed), func(*discordgo.Session, *discordgo.TypingStart), func(*discordgo.Session, *discordgo.UserGuildSettingsUpdate), func(*discordgo.Session, *discordgo.UserNoteUpdate), func(*discordgo.Session, *discordgo.UserSettingsUpdate), func(*discordgo.Session, *discordgo.UserUpdate), func(*discordgo.Session, *discordgo.VoiceServerUpdate), func(*discordgo.Session, *discordgo.VoiceStateUpdate):
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
