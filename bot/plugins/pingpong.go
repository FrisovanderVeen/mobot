package plugins

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	pingPongHelp = Help{
		Commands: map[string]string{
			fmt.Sprintf("%sping", Prefix): "sends Pong! to the text channel",
			fmt.Sprintf("%spong", Prefix): "sends Ping! to the text channel",
		},
		View:        true,
		Explanation: "replies to pong with Ping! and to ping with Pong!",
	}
	_ = Register("PingPong", pingPongHelp, pingPong)
)

func pingPong(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == fmt.Sprintf("%sping", Prefix) {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == fmt.Sprintf("%spong", Prefix) {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
