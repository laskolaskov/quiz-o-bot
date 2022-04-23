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
		fmt.Println("timer :::", time.Duration(input.Time)*time.Second)

		time.Sleep(time.Duration(input.Time) * time.Second)
	}

	//TODO: game ended message
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

	//process result
	r.Add(q, isCorrect, questionIndex)

	if isCorrect {
		replay(s, m, "correct !!") //TODO: the actual message
	} else {
		replay(s, m, "wrong ! correct is: "+questions[questionIndex].Correct_answer) //TODO: the actual message
	}
	fmt.Println(game.PlayerResult(m.Author.ID, m.Author.Username))
}

func result(s *discordgo.Session, m *discordgo.MessageCreate) {
	lb := storage.Game(m.ChannelID).LB()
	s.ChannelMessageSendEmbed(m.ChannelID, print.LB(lb))
}

func myResult(s *discordgo.Session, m *discordgo.MessageCreate) {
	r := storage.Game(m.ChannelID).PlayerResult(m.Author.ID, m.Author.Username)
	fmt.Println("--- printing player result ---")
	fmt.Println(r)
	cn, _ := s.Channel(m.ChannelID)
	replay(s, m, "Here will be your result from channel '"+cn.Name+"' !") //TODO the actual message
}
