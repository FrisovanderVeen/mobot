package plugins

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/FrisovanderVeen/mobot/bot/plugins"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

var (
	youtubeHelp = plugins.Help{
		Commands:    map[string]string{},
		View:        true,
		Explanation: "plays audio of youtube videos in the users voice channel",
	}
	_ = plugins.Register("Youtube", youtubeHelp, youtubeFunc)

	key string

	queries = make(map[string]*youtubeQuery)
)

type youtubeQuery struct {
	Videos []*youtubeVideo
}

type youtubeVideo struct {
	ID          string
	Title       string
	Author      string
	Description string
}

func youtubeFunc(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	switch {
	case strings.HasPrefix(m.Content, fmt.Sprintf("%splay ", plugins.Prefix)):
		if key == "" {
			key = plugins.Config.Youtube.Key
		}

		client := &http.Client{
			Transport: &transport.APIKey{Key: key},
		}

		service, err := youtube.New(client)
		if err != nil {
			plugins.ErrChan <- fmt.Errorf("Could not create new youtube client: %v", err)
			return
		}
		call := service.Search.List("id, snippet").
			Q(strings.TrimPrefix(m.Content, fmt.Sprintf("%splay ", plugins.Prefix))).
			MaxResults(5).
			Type("video")

		resp, err := call.Do()
		if err != nil {
			plugins.ErrChan <- fmt.Errorf("Could not do youtube api call: %v", err)
			return
		}

		yq := &youtubeQuery{
			Videos: []*youtubeVideo{},
		}

		embed := &discordgo.MessageEmbed{}

		i := 1

		for _, video := range resp.Items {
			yv := &youtubeVideo{
				ID:          video.Id.VideoId,
				Title:       video.Snippet.Title,
				Author:      video.Snippet.ChannelTitle,
				Description: video.Snippet.Description,
			}
			yq.Videos = append(yq.Videos, yv)
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("%v. %s - %s", i, yv.Title, yv.Author),
				Value: yv.Description,
			})
			i++
		}
		queries[m.Author.ID] = yq
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("Use %schoose <number> to choose", plugins.Prefix),
			Value: fmt.Sprintf("For example %schoose 3", plugins.Prefix),
		})
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	case strings.HasPrefix(m.Content, fmt.Sprintf("%schoose ", plugins.Prefix)):
		yq, ok := queries[m.Author.ID]
		if !ok {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Please use %splay first, for example: %splay yee", plugins.Prefix, plugins.Prefix))
			return
		}

		number, err := strconv.Atoi(strings.TrimPrefix(m.Content, fmt.Sprintf("%schoose ", plugins.Prefix)))
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Please use whole actual numbers, for example: %schoose 4", plugins.Prefix))
			return
		}
		number--

		if number < 0 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Please use numbers above 0, for example: %schoose 3", plugins.Prefix))
			return
		} else if number > len(queries[m.Author.ID].Videos)-1 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Please use numbers less than the number of videos returned, for example: %schoose 1", plugins.Prefix))
			return
		}

		vid := yq.Videos[number]
		options := dca.StdEncodeOptions
		options.RawOutput = true
		options.Bitrate = 96
		options.Application = "lowdelay"

		videoInfo, err := ytdl.GetVideoInfo(fmt.Sprintf("https://youtube.com/watch?v=%s", vid.ID))
		if err != nil {
			plugins.ErrChan <- fmt.Errorf("Could not get video info:%v", err)
			return
		}

		format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
		downloadURL, err := videoInfo.GetDownloadURL(format)
		if err != nil {
			plugins.ErrChan <- fmt.Errorf("Could not get download url:%v", err)
			return
		}

		encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
		if err != nil {
			plugins.ErrChan <- fmt.Errorf("Could not encode video:%v", err)
			return
		}
		defer encodingSession.Cleanup()
		filetemp, err := ioutil.TempFile("", "")
		if err != nil {
			plugins.ErrChan <- fmt.Errorf("Could not create temp file: %v", err)
			return
		}
		io.Copy(filetemp, encodingSession)
		if err = filetemp.Close(); err != nil {
			plugins.ErrChan <- fmt.Errorf("Could not close temp file: %v", err)
			os.Remove(filetemp.Name())
			return
		}
		defer os.Remove(filetemp.Name())
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
				err = playSound(s, vid, filetemp.Name(), g.ID, vs.ChannelID, m.ChannelID)
				if err != nil {
					plugins.ErrChan <- fmt.Errorf("Could not play sound: %v")
				}

				return
			}
		}

		delete(queries, m.Author.ID)
	case strings.HasPrefix(m.Content, fmt.Sprintf("%scancel", plugins.Prefix)):
		_, ok := queries[m.Author.ID]
		if !ok {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Please use %splay first, for example: %splay yee", plugins.Prefix, plugins.Prefix))
			return
		}
		delete(queries, m.Author.ID)
	}
}

func loadSound(fileloc string) ([][]byte, error) {
	file, err := os.Open(fileloc)
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
	var sound = make([][]byte, 0)

	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return sound, nil
		}
		if err != nil {
			return nil, fmt.Errorf("Could not read file: %v", err)
		}

		inBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &inBuf)
		if err != nil {
			return nil, fmt.Errorf("Could not read file: %v", err)
		}

		sound = append(sound, inBuf)
	}
}

func playSound(s *discordgo.Session, vid *youtubeVideo, fileloc, guildID, voiceChannelID, textChannelID string) error {
	buffer, err := loadSound(fileloc)
	if err != nil {
		return fmt.Errorf("load airhorn sound: %v", err)
	}

	plugins.SoundQueue <- &plugins.Sound{
		View:           true,
		Content:        buffer,
		Session:        s,
		GuildID:        guildID,
		VoiceChannelID: voiceChannelID,
		TextChannelID:  textChannelID,
		Name:           vid.Title,
		Author:         vid.Author,
	}

	return nil
}
