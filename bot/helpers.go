package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func isDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)

	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}

func isTextChannel(c *discordgo.Channel) bool {
	return c.Type == discordgo.ChannelTypeGuildText
}

func process(m string) (string, string) {
	s := strings.SplitAfterN(m, " ", 2)
	command := s[0]
	args := ""
	if len(s) > 1 {
		args = s[1]
	}
	return command, args
}

//replay to DM with string message
func replay(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	ch, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		fmt.Println("Error while creating DM channel:", err)
	}
	s.ChannelMessageSend(ch.ID, msg)
}
