package storage

import "github.com/laskolaskov/quiz-o-bot/api"

type Data struct {
	questions  []api.Question
	categories []api.Category
	progress   bool
}

var data = Data{}

func Questions() []api.Question {
	return data.questions
}

func SetQuestions(q []api.Question) {
	data.questions = q
}

func Categories() []api.Category {
	return data.categories
}

func SetCategories(c []api.Category) {
	data.categories = c
}

func Progress() bool {
	return data.progress
}

func SetProgress(p bool) {
	data.progress = p
}
