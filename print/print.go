package print

import (
	"fmt"
	"html"

	"github.com/bwmarrin/discordgo"
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

func QuestionEmbed(index int, q api.Question) *discordgo.MessageEmbed {
	fmt.Printf("%v\n", q.Incorrect_answers)
	//prepareAnswers(&q)
	//fmt.Printf("%v\n", q.Incorrect_answers)
	embed := NewEmbed().
		SetTitle(fmt.Sprintf("Question #%v : %v", index, html.UnescapeString(q.Question))).
		SetDescription(fmt.Sprintf("%v - %v", q.Category, q.Difficulty)).
		SetColor(0x0000ff)

	for i, a := range q.Incorrect_answers {
		embed.AddField(fmt.Sprintf("#%v-%v", index, i), html.UnescapeString(a))
	}

	return embed.MessageEmbed
}
