package plugins

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	airhornHelp = Help{
		Commands: map[string]string{
			fmt.Sprintf("%sairhorn", Prefix): "plays airhorn in the users current voice channel",
		},
		View:        true,
		Explanation: "annoys people with airhorns",
	}
	_ = Register("Airhorn", airhornHelp, airhorn)
)

func airhorn(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, fmt.Sprintf("%sairhorn", Prefix)) {
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			WarnChan <- fmt.Errorf("Could not find channel: %v", err)
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			WarnChan <- fmt.Errorf("Could not find guild: %v", err)
			return
		}

		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = playAirhornSound(s, g.ID, vs.ChannelID)
				if err != nil {
					ErrChan <- fmt.Errorf("Could not play sound: %v")
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
			ErrChan <- fmt.Errorf("Could not close file: %v", err)
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

func playAirhornSound(s *discordgo.Session, guildID string, channelID string) error {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}
	defer func() {
		vc.Disconnect()
	}()

	time.Sleep(250 * time.Millisecond)

	vc.Speaking(true)

	buffer, err := loadAirhornSound()
	if err != nil {
		return fmt.Errorf("load airhorn sound: %v", err)
	}

	for _, buf := range buffer {
		vc.OpusSend <- buf
	}

	vc.Speaking(false)

	time.Sleep(250 * time.Millisecond)

	return nil
}
