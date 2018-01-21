package plugins

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/FrisovanderVeen/mobot/bot/plugins"
	"github.com/bwmarrin/discordgo"
)

var (
	airhornHelp = plugins.Help{
		Commands: map[string]string{
			fmt.Sprintf("%sairhorn", plugins.Prefix): "Plays a airhorn sound in the users current voice channel.",
		},
		View:        true,
		Explanation: "annoys people with airhorns",
	}
	_ = plugins.Register("Airhorn", airhornHelp, airhorn)
)

func airhorn(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, fmt.Sprintf("%sairhorn", plugins.Prefix)) {
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			plugins.WarnChan <- fmt.Errorf("Could not find channel: %v", err)
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			plugins.WarnChan <- fmt.Errorf("Could not find guild: %v", err)
			return
		}

		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = playAirhornSound(s, g.ID, vs.ChannelID)
				if err != nil {
					plugins.ErrChan <- fmt.Errorf("Could not play sound: %v")
				}

				return
			}
		}
	}
}

func loadAirhornSound() ([][]byte, error) {
	file, err := os.Open("recources/sound/airhorn.dca")
	if err != nil {
		return nil, fmt.Errorf("Could not open file: %v", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			plugins.ErrChan <- fmt.Errorf("Could not close file: %v", err)
		}
	}()

	var opuslen int16
	var airhornsound = make([][]byte, 0)

	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return airhornsound, nil
		}
		if err != nil {
			return nil, fmt.Errorf("Could not read file: %v", err)
		}

		inBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &inBuf)
		if err != nil {
			return nil, fmt.Errorf("Could not read file: %v", err)
		}

		airhornsound = append(airhornsound, inBuf)
	}
}

func playAirhornSound(s *discordgo.Session, guildID string, voiceChannelID string) error {
	buffer, err := loadAirhornSound()
	if err != nil {
		return fmt.Errorf("load airhorn sound: %v", err)
	}

	plugins.SoundQueue <- &plugins.Sound{
		View:           false,
		Content:        buffer,
		Session:        s,
		GuildID:        guildID,
		VoiceChannelID: voiceChannelID,
	}

	return nil
}
