package storage

import "github.com/laskolaskov/quiz-o-bot/api"

var games = make(map[string]*GameState)
var categories = make([]api.Category, 0, 50)

//var data = Data{games, categories}
//var dataMap = Data{games, categories}

/* type Data struct {
	games      map[string]GameState
	categories []api.Category
} */

type GameState struct {
	questions []api.Question
	result    map[string]*Result
	progress  bool
}

type Result struct {
}

func (g *GameState) Questions() []api.Question {
	return g.questions
}

func (g *GameState) SetQuestions(q []api.Question) {
	g.questions = q
}

func (g *GameState) Result() map[string]*Result {
	return g.result
}

func (g *GameState) Progress() bool {
	return g.progress
}

func (g *GameState) SetProgress(p bool) {
	g.progress = p
}

func Categories() []api.Category {
	return categories
}

func SetCategories(c []api.Category) {
	categories = c
}

func Game(id string) *GameState {
	_, exists := games[id]
	if exists {
		return games[id]
	}
	games[id] = &GameState{}
	return games[id]
}
