package print

import (
	"fmt"
	"html"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/laskolaskov/quiz-o-bot/api"
	"github.com/laskolaskov/quiz-o-bot/storage"
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

func LB(lb []storage.Score) *discordgo.MessageEmbed {
	embed := NewEmbed().
		SetTitle(fmt.Sprintf("Standings")).
		SetColor(0x0000ff)

	for i, s := range lb {
		if len(strings.TrimSpace(s.Name)) > 0 {
			embed.AddField(fmt.Sprintf("#%v %v : %v", i+1, s.Name, s.Score) /* strconv.Itoa(s.Score) */, html.UnescapeString("\u200b"))
		}
	}

	return embed.MessageEmbed
}

func Question(index int, q api.Question) *discordgo.MessageEmbed {
	embed := NewEmbed().
		SetTitle(fmt.Sprintf("Question #%v : %v", index, html.UnescapeString(q.Question))).
		SetDescription(fmt.Sprintf("%v - %v", q.Category, q.Difficulty)).
		SetColor(0x0000ff)

	for i, a := range q.Incorrect_answers {
		embed.AddField(fmt.Sprintf("#%v-%v", index, i), html.UnescapeString(a))
	}

	return embed.MessageEmbed
}

func Starting() *discordgo.MessageEmbed {
	embed := NewEmbed().
		SetTitle(fmt.Sprintf("Game starts!")).
		SetDescription(fmt.Sprintf("Prepare your brains...")).
		SetColor(0x0000ff)

	return embed.MessageEmbed
}
