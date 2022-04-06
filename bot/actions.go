package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/laskolaskov/quiz-o-bot/api"
	"github.com/laskolaskov/quiz-o-bot/print"
	"github.com/laskolaskov/quiz-o-bot/storage"
)

func start(s *discordgo.Session, m *discordgo.MessageCreate, args string) {
	channelID := m.ChannelID
	game := storage.Game(channelID)

	if game.Progress() {
		fmt.Println("Game started for channel: ", channelID)
	}

	input := api.ParseInput(strings.Split(args, " "))

	questions, err := api.GetQuestions(input.Amount, input.Category, input.Difficulty)

	if err != nil {
		log.Fatal(err)
	}

	game.SetQuestions(prepareAnswers(questions))
	//storage.SetGame(game, channelID)

	fmt.Println(game)
	fmt.Println(storage.Game(channelID))

	//start asking the questions
	for i, q := range game.Questions() {

		embed := print.QuestionEmbed(i, q)

		//s.ChannelMessageSend(channelID, q.Question)
		s.ChannelMessageSendEmbed(channelID, embed)
		time.Sleep(1 * time.Second)
		//s.ChannelMessageSendEmbed()
	}
}

func result(s *discordgo.Session, m *discordgo.MessageCreate, questionIndex int, answerIndex int) {
	fmt.Println(questionIndex, answerIndex)
	fmt.Println(m.ChannelID, m.Content, m.Author.ID, m.Author.Username)

	game := storage.Game(m.ChannelID)
	fmt.Println(game)
	questions := game.Questions()

	if isCorrect(questions, questionIndex, answerIndex) {
		replay(s, m, "correct !!")
	} else {
		replay(s, m, "wrong ! correct is: "+questions[questionIndex].Correct_answer)
	}
}
