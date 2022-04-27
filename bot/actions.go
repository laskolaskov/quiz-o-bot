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
		replay(s, m, "Game is already running! Try again after the game is finished.")
		return
	}

	//clear the game and init a new one, removing the current result
	game = storage.NewGame(channelID)

	input := api.ParseInput(strings.Split(args, " "))

	questions, err := api.GetQuestions(input.Amount, input.Category, input.Difficulty)

	if err != nil {
		log.Fatal(err)
		s.ChannelMessageSend(channelID, "Something went wrong and the questions could not be loaded. Try again a bit later maybe it will fix itself...")
		return
	}

	//set/replace the game questions and start it
	game.SetQuestions(prepareAnswers(questions))
	game.SetProgress(true)

	//starting the game
	s.ChannelMessageSendEmbed(channelID, print.Starting())
	time.Sleep(5 * time.Second)

	//asking the questions
	for i, q := range game.Questions() {

		s.ChannelMessageSendEmbed(channelID, print.Question(i, q))
		game.SetCurrent(i)

		time.Sleep(time.Duration(input.Time) * time.Second)
	}

	s.ChannelMessageSendEmbed(channelID, print.Ended())
	result(s, m)
	game.SetProgress(false)
}

func processAnswer(s *discordgo.Session, m *discordgo.MessageCreate, questionIndex int, answerIndex int) {

	game := storage.Game(m.ChannelID)
	questions := game.Questions()
	current := game.Current()

	if !game.Progress() {
		replay(s, m, "Game has ended.")
		return
	}
	//validate if answer is for the currently asked question
	if questionIndex != current {
		replay(s, m, fmt.Sprintf("The current question is #%v, not #%v.", current, questionIndex))
		return
	}
	//validate indices are in bounds
	if questionIndex >= len(questions) || answerIndex >= len(questions[questionIndex].Incorrect_answers) {
		replay(s, m, fmt.Sprintf("Your answer #%v-%v is not matching the available answers. Please check the question and choose correct option.", questionIndex, answerIndex))
		return
	}

	//get player result
	r := game.PlayerResult(m.Author.ID, m.Author.Username)

	//validate if the question has been answered - allow only one attempt to answer each question
	if r.Answered(questionIndex) {
		replay(s, m, "You already answered this question!")
		return
	}

	q := questions[questionIndex]
	isCorrect := isCorrect(questions, questionIndex, answerIndex)

	//set the ansewer
	q.Answer = q.Incorrect_answers[answerIndex]

	//process result
	r.Add(q, isCorrect, questionIndex)

	if isCorrect {
		replayEmbed(s, m, print.Correct(q))
	} else {
		replayEmbed(s, m, print.Incorrect(q))
	}
}

func result(s *discordgo.Session, m *discordgo.MessageCreate) {
	lb := storage.Game(m.ChannelID).LB()
	s.ChannelMessageSendEmbed(m.ChannelID, print.LB(lb))
}

func myResult(s *discordgo.Session, m *discordgo.MessageCreate) {
	r := storage.Game(m.ChannelID).PlayerResult(m.Author.ID, m.Author.Username)
	cn, _ := s.Channel(m.ChannelID)
	replayEmbed(s, m, print.Result(r, cn.Name))
}

func testAction(s *discordgo.Session, m *discordgo.MessageCreate) {
	replayEmbed(s, m, print.TestEmbed())
}
