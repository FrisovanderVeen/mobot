package plugins

import (
	"fmt"

	"github.com/FrisovanderVeen/mobot/bot/plugins"
	"github.com/bwmarrin/discordgo"
)

var (
	pingPongHelp = plugins.Help{
		Commands: map[string]string{
			fmt.Sprintf("%sping", plugins.Prefix): "Sends Pong! to the text channel.",
			fmt.Sprintf("%spong", plugins.Prefix): "Sends Ping! to the text channel.",
		},
		View:        true,
		Explanation: "replies to pong with Ping! and to ping with Pong!",
	}
	_ = plugins.Register("PingPong", pingPongHelp, pingPong)
)

func pingPong(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == fmt.Sprintf("%sping", plugins.Prefix) {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == fmt.Sprintf("%spong", plugins.Prefix) {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
