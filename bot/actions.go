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
		//TODO: send message
		return
	}

	//TODO: clear/initialize the game, removing the current result

	input := api.ParseInput(strings.Split(args, " "))

	questions, err := api.GetQuestions(input.Amount, input.Category, input.Difficulty)

	if err != nil {
		log.Fatal(err)
		//TODO: send message
		return
	}

	//set/replace the game questions and start it
	game.SetQuestions(prepareAnswers(questions))
	game.SetProgress(true)

	//TODO: send message that the game starts, maybe wait for a few seconds before the first question is asked

	//asking the questions
	for i, q := range game.Questions() {

		embed := print.Question(i, q)

		s.ChannelMessageSendEmbed(channelID, embed)
		game.SetCurrent(i)

		time.Sleep(5 * time.Second) //TODO: make the time between questions a command argument for !start ?
	}

	result(s, m)
	game.SetProgress(false)
}

func processAnswer(s *discordgo.Session, m *discordgo.MessageCreate, questionIndex int, answerIndex int) {

	game := storage.Game(m.ChannelID)
	questions := game.Questions()

	//validate indices are in bounds
	//TODO: validate if answer is for the currently asked question
	if questionIndex >= len(questions) || answerIndex >= len(questions[questionIndex].Incorrect_answers) {
		replay(s, m, "out of bounds !")
		return
	}

	q := questions[questionIndex]
	isCorrect := isCorrect(questions, questionIndex, answerIndex)

	//process result
	r := game.PlayerResult(m.Author.ID, m.Author.Username)
	r.Add(q, isCorrect)

	if isCorrect {
		replay(s, m, "correct !!")
	} else {
		replay(s, m, "wrong ! correct is: "+questions[questionIndex].Correct_answer)
	}
	fmt.Println(game.PlayerResult(m.Author.ID, m.Author.Username))
}

func result(s *discordgo.Session, m *discordgo.MessageCreate) {
	//TODO implement simple sorted result leaderboard
	lb := storage.Game(m.ChannelID).LB()
	s.ChannelMessageSendEmbed(m.ChannelID, print.LB(lb))
}

func myResult(s *discordgo.Session, m *discordgo.MessageCreate) {
	r := storage.Game(m.ChannelID).PlayerResult(m.Author.ID, m.Author.Username)
	fmt.Println("--- printing player result ---")
	fmt.Println(r)
	cn, _ := s.Channel(m.ChannelID)
	replay(s, m, "Here will be your result from channel '"+cn.Name+"' !")
}
