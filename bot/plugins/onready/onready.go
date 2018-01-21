package plugins

import (
	"fmt"

	"github.com/FrisovanderVeen/mobot/bot/plugins"
	"github.com/bwmarrin/discordgo"
)

var (
	onReadyHelp = plugins.Help{
		View:        false,
		Explanation: "sets the discord status of the bot",
	}
	_ = plugins.Register("OnReady", onReadyHelp, onReady)
)

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, fmt.Sprintf("%shelp", plugins.Prefix))
}
