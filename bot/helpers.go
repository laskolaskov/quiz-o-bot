package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/laskolaskov/quiz-o-bot/api"
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
	s := strings.SplitN(m, " ", 2)
	command := s[0]
	args := ""
	if len(s) > 1 {
		args = s[1]
	}
	return command, args
}

//replay with DM string message
func replay(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	ch, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		fmt.Println("Error while creating DM channel:", err)
	}
	s.ChannelMessageSend(ch.ID, msg)
}

//replay with DM embed
func replayEmbed(s *discordgo.Session, m *discordgo.MessageCreate, e *discordgo.MessageEmbed) {
	ch, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		fmt.Println("Error while creating DM channel:", err)
	}
	s.ChannelMessageSendEmbed(ch.ID, e)
}

func checkAnswer(m string) (int, int, error) {
	if "#" != m[:1] {
		return 0, 0, errors.New("not an answer command: missing '#' prefix")
	}

	command := strings.SplitN(m[1:], "-", 2)

	if len(command) < 2 {
		return 0, 0, errors.New("not an answer command: must have two numbers splitted by '-'")
	}

	q, err := strconv.Atoi(command[0])

	if err != nil {
		return 0, 0, errors.New("not an answer command: first param is not integer")
	}

	a, err := strconv.Atoi(command[1])

	if err != nil {
		return 0, 0, errors.New("not an answer command: second param is not integer")
	}

	return q, a, nil
}

func prepareAnswers(questions []api.Question) []api.Question {
	for i, q := range questions {
		a := q.Incorrect_answers
		a = append(a, q.Correct_answer)
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(a), func(i, j int) {
			a[i], a[j] = a[j], a[i]
		})
		questions[i].Incorrect_answers = a
	}
	return questions
}

func isCorrect(questions []api.Question, qIndex int, aIndex int) bool {
	q := questions[qIndex]
	check := q.Correct_answer == q.Incorrect_answers[aIndex]
	return check
}
