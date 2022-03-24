package print

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/laskolaskov/quiz-o-bot/api"
)

var red = color.New(color.FgRed).SprintfFunc()

func Help() string {
	return fmt.Sprintf("\nHelp message will be here! Someday ...")
}

func ListCategories(categories []api.Category) string {
	msg := "Here is the list of categories. Use the number in the command to start the game.\nExample: !start -c=19 for \"Science: Mathematics\"\n\n"
	for _, c := range categories {
		msg += fmt.Sprintf("%v - %v\n", c.Id, c.Name)
	}
	return msg
}
