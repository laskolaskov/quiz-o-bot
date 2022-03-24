package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const apiUrl string = "https://opentdb.com/api.php"
const categoryUrl string = "https://opentdb.com/api_category.php"

const defaultAmount int = 5

var inputMap = map[string]string{
	"-c":                "category",
	"--category":        "category",
	"-a":                "amount",
	"--amount":          "amount",
	"-d":                "difficulty",
	"--difficulty":      "difficulty",
	"-h":                "help",
	"--help":            "help",
	"?":                 "help",
	"-lc":               "listCategories",
	"--list-categories": "listCategories"}

type Category struct {
	Id   int
	Name string
}

type CategoryResponse struct {
	Trivia_categories []Category
}

type Question struct {
	Category          string
	Type              string
	Difficulty        string
	Question          string
	Correct_answer    string
	Incorrect_answers []string
	Answer            string
	Correct           bool
}

type QuestionResponse struct {
	Response_code int
	Results       []Question
}

type ApiArgs struct {
	Amount         int
	Category       int
	Difficulty     string
	Help           bool
	ListCategories bool
}

func GetCategories() ([]Category, error) {
	body, err := get(categoryUrl)
	if err != nil {
		return []Category{}, err
	}

	var categories CategoryResponse

	if e := json.Unmarshal(body, &categories); e != nil {
		return []Category{}, e
	}

	return categories.Trivia_categories, nil
}

func GetQuestions(amount int, category int, difficulty string) ([]Question, error) {
	url := fmt.Sprintf("%v?type=multiple&amount=%v", apiUrl, getAmount(amount))

	if category > 0 {
		url += fmt.Sprintf("&category=%v", category)
	}
	if validDifficulty(difficulty) {
		url += fmt.Sprintf("&difficulty=%v", difficulty)
	}

	body, err := get(url)
	if err != nil {
		return []Question{}, err
	}

	var questions QuestionResponse

	if e := json.Unmarshal(body, &questions); e != nil {
		return []Question{}, e
	}

	return questions.Results, nil
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return []byte{}, fmt.Errorf("%v : %v\n", url, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func validDifficulty(d string) bool {
	return d == "easy" || d == "medium" || d == "hard"
}

func getAmount(a int) int {
	if a <= 0 {
		return defaultAmount
	}
	return a
}

func ParseInput(input []string) ApiArgs {

	apiArgs := ApiArgs{}

	for _, v := range input {
		parameter := strings.Split(v, "=")
		name, e := inputMap[parameter[0]]
		if e {
			switch name {
			case "category":
				c, e := strconv.Atoi(parameter[1])
				if e == nil {
					apiArgs.Category = c
				}
			case "difficulty":
				apiArgs.Difficulty = parameter[1]
			case "amount":
				a, e := strconv.Atoi(parameter[1])
				if e == nil {
					apiArgs.Amount = a
				}
			case "help":
				apiArgs.Help = true
			case "listCategories":
				apiArgs.ListCategories = true
			}
		}
	}

	return apiArgs
}
