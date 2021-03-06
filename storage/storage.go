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
	Correct  int
	Total    int
	Easy     Answered
	Medium   Answered
	Hard     Answered
	name     string
	id       string
	answered map[int]bool
}

type Answered struct {
	Correct int
	Total   int
}

type Score struct {
	Name  string
	Score int
}

func (r *Result) Score() Score {
	score := r.Easy.Correct*1 + r.Medium.Correct*3 + r.Hard.Correct*5
	return Score{Name: r.name, Score: score}
}

func (r *Result) Answered(index int) bool {
	_, exists := r.answered[index]
	return exists
}

func (r *Result) Add(q api.Question, isCorrect bool, index int) {
	switch q.Difficulty {
	case "easy":
		if isCorrect {
			r.Easy.Correct++
			r.Correct++
		}
		r.Easy.Total++
	case "medium":
		if isCorrect {
			r.Medium.Correct++
			r.Correct++
		}
		r.Medium.Total++
	case "hard":
		if isCorrect {
			r.Hard.Correct++
			r.Correct++
		}
		r.Hard.Total++
	}
	r.Total++
	r.answered[index] = isCorrect
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
	//init
	g.result[id].answered = make(map[int]bool)

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
	return NewGame(id)
}

func NewGame(id string) *GameState {
	games[id] = &GameState{
		questions: make([]api.Question, 0),
		result:    make(map[string]*Result),
		progress:  false,
	}
	return games[id]
}
