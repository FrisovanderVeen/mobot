package plugins

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	onReadyHelp = Help{
		View:        false,
		Explanation: "sets the discord status of the bot",
	}
	_ = Register("OnReady", onReadyHelp, onReady)
)

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, fmt.Sprintf("%shelp", Prefix))
}
