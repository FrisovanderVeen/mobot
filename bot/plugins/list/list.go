package plugins

import (
	"fmt"

	"github.com/FrisovanderVeen/mobot/bot/plugins"
	"github.com/bwmarrin/discordgo"
)

var (
	listHelp = plugins.Help{
		Commands: map[string]string{
			fmt.Sprintf("%slist", plugins.Prefix): "lists all plugins in a text channel",
		},
		View:        true,
		Explanation: "lists all plugins for discord users",
	}
	_ = plugins.Register("List", listHelp, list)
)

func list(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == fmt.Sprintf("%slist", plugins.Prefix) {
		for name, plugin := range plugins.Plugins {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %s", name, plugin.Help.Explanation))
		}
	}
}
