package audio

import (
	"fmt"
	"strings"

	"github.com/FrisovanderVeen/mobot/bot/plugins"
	"github.com/bwmarrin/discordgo"
)

var (
	audioHelp = plugins.Help{
		Commands: map[string]string{
			fmt.Sprintf("%sskip", plugins.Prefix): "Skips the currently playing song.",
		},
		View:        true,
		Explanation: "Provides helper commands for audio",
	}
	_ = plugins.Register("Audio", audioHelp, audio)
)

func audio(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch {
	case strings.HasPrefix(m.Content, fmt.Sprintf("%sskip", plugins.Prefix)):
		plugins.SkipSig <- 1
	}
}
