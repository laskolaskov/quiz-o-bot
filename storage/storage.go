package storage

import (
	"sort"

	"github.com/laskolaskov/quiz-o-bot/api"
)

var games = make(map[string]*GameState)
var categories = make([]api.Category, 0, 50)

type GameState struct {
	questions []api.Question
	result    map[string]*Result
	progress  bool
	current   int
}

type Result struct {
	correct int
	total   int
	easy    Answered
	medium  Answered
	hard    Answered
	name    string
	id      string
}

type Answered struct {
	correct int
	total   int
}

type Score struct {
	Name  string
	Score int
}

func (r *Result) Score() Score {
	score := r.easy.correct*1 + r.medium.correct*3 + r.hard.correct*5
	return Score{Name: r.name, Score: score}
}

func (r *Result) Add(q api.Question, isCorrect bool) {
	switch q.Difficulty {
	case "easy":
		if isCorrect {
			r.easy.correct++
			r.correct++
		}
		r.easy.total++
	case "medium":
		if isCorrect {
			r.medium.correct++
			r.correct++
		}
		r.medium.total++
	case "hard":
		if isCorrect {
			r.hard.correct++
			r.correct++
		}
		r.hard.total++
	}
	r.total++
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

func (g *GameState) LB() []Score {
	lb := make([]Score, len(g.result))
	for _, r := range g.result {
		lb = append(lb, r.Score())
	}
	sort.Slice(lb, func(a, b int) bool { return lb[a].Score > lb[b].Score })
	return lb
}

func (g *GameState) PlayerResult(id string, name string) *Result {
	r, exists := g.result[id]
	if exists {
		return r
	}
	//init
	g.result[id] = &Result{}
	//set player name and ID
	g.result[id].name = name
	g.result[id].id = id

	return g.result[id]
}

func (g *GameState) Progress() bool {
	return g.progress
}

func (g *GameState) SetProgress(p bool) {
	g.progress = p
}

func (g *GameState) Current() int {
	return g.current
}

func (g *GameState) SetCurrent(i int) {
	g.current = i
}

func Categories() []api.Category {
	return categories
}

func SetCategories(c []api.Category) {
	categories = c
}

func Game(id string) *GameState {
	g, exists := games[id]
	if exists {
		return g
	}
	games[id] = &GameState{
		questions: make([]api.Question, 0),
		result:    make(map[string]*Result),
		progress:  false,
	}
	return games[id]
}
