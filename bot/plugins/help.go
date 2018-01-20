package plugins

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	helpHelp = Help{
		Commands: map[string]string{
			fmt.Sprintf("%shelp", Prefix): "sends all commands and their uses to the text channel",
		},
		View:        true,
		Explanation: "gives help to the discord users",
	}
	_ = Register("Help", helpHelp, help)
)

func help(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, fmt.Sprintf("%shelp", Prefix)) {
		for _, plugin := range Plugins {
			if plugin.Help.View {
				for com, exp := range plugin.Help.Commands {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s%s: %s", Prefix, com, exp))
				}
			}
		}
	}
}
