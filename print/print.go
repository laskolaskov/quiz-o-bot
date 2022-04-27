package print

import (
	"fmt"
	"html"
	"math"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/laskolaskov/quiz-o-bot/api"
	"github.com/laskolaskov/quiz-o-bot/storage"
)

var red = color.New(color.FgRed).SprintfFunc()

func Help() *discordgo.MessageEmbed {
	embed := NewEmbed().
		SetTitle("Usage").
		AddField("Command: !start <arg1 arg2...> : starts a new game in the guild channel you type into.", "example: !start, !start -a=15 -c=21 -d=easy").
		AddField("-c, --category", "Set the category of the questions. Not using this defaults to ALL categories. To see the category IDs, use !categories. Example: !start -c=21 : setscategory with ID 21, which, at the time of writing this, happens to be 'Sports' "+html.UnescapeString("&#128512;.")).
		AddField("-a, --amount", "Set the amount of questions the game will have. Not using this defaults to 5. Example: !start --amount=10 : the game will have 10 questions.").
		AddField("-d, --difficulty", "Sets the difficulty of the questions. Possible values: easy|medium|hard. Not using this defaults to ALL difficulties. Example: !start -d=easy.").
		AddField("-t, --time", "Sets the time (in seconds) players have to answer each question. Not using this defaults to 25 seconds.").
		AddField("Command: !categories : Send this (also usable in DM) to get the list of the question categories available and their IDs in a DM replay.", "Example: !categories "+html.UnescapeString("&#128512;.")).
		AddField("Command: !result : Show the standings for the current/last game.", "Example: !result "+html.UnescapeString("&#128518;.")).
		AddField("Command: !myResult : Shows you details about your result in the current/last game in a DM replay.", "Example: !myResult "+html.UnescapeString("&#129315;.")).
		SetColor(0xffffff).MessageEmbed

	return embed
}

func ListCategories(categories []api.Category) string {
	msg := "Here is the list of categories. Use the number in the command to start the game.\nExample: !start -c=19 for \"Science: Mathematics\"\n\n"
	for _, c := range categories {
		msg += fmt.Sprintf("%v - %v\n", c.Id, html.UnescapeString(c.Name))
	}
	return msg
}

func LB(lb []storage.Score) *discordgo.MessageEmbed {
	embed := NewEmbed().
		SetTitle(fmt.Sprintf("Standings")).
		SetColor(0x0000ff)

	for i, s := range lb {
		if len(strings.TrimSpace(s.Name)) > 0 {
			embed.AddField(fmt.Sprintf("#%v %v : %v", i+1, s.Name, s.Score), html.UnescapeString("\u200b"))
		}
	}

	return embed.MessageEmbed
}

func Question(index int, q api.Question) *discordgo.MessageEmbed {
	embed := NewEmbed().
		SetTitle(fmt.Sprintf("Question #%v : %v", index, html.UnescapeString(q.Question))).
		SetDescription(fmt.Sprintf("%v - %v", html.UnescapeString(q.Category), q.Difficulty)).
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

func Ended() *discordgo.MessageEmbed {
	embed := NewEmbed().
		SetTitle(fmt.Sprintf("Game ended!")).
		SetColor(0x0000ff)

	return embed.MessageEmbed
}

func Correct(q api.Question) *discordgo.MessageEmbed {
	embed := NewEmbed().
		SetTitle(fmt.Sprintf("%v", html.UnescapeString(q.Question))).
		SetColor(0x00ff00).
		AddField(fmt.Sprintf("Your answer '%v' is CORRECT!", html.UnescapeString(q.Answer)), html.UnescapeString("\u200b"))

	return embed.MessageEmbed
}

func Incorrect(q api.Question) *discordgo.MessageEmbed {
	embed := NewEmbed().
		SetTitle(fmt.Sprintf("%v", html.UnescapeString(q.Question))).
		SetColor(0xff0000).
		AddField(fmt.Sprintf("Your answer '%v' is WRONG!", html.UnescapeString(q.Answer)), html.UnescapeString("\u200b")).
		AddField(fmt.Sprintf("The correct answer is: '%v'", html.UnescapeString(q.Correct_answer)), html.UnescapeString("\u200b"))

	return embed.MessageEmbed
}

func Result(r *storage.Result, channelName string) *discordgo.MessageEmbed {

	//calculate percentages
	e := math.Floor((float64(r.Easy.Correct) / float64(r.Easy.Total)) * 100)
	if math.IsNaN(e) {
		e = 0
	}

	m := math.Floor((float64(r.Medium.Correct) / float64(r.Medium.Total)) * 100)
	if math.IsNaN(m) {
		m = 0
	}

	h := math.Floor((float64(r.Hard.Correct) / float64(r.Hard.Total)) * 100)
	if math.IsNaN(h) {
		h = 0
	}

	t := math.Floor((float64(r.Correct) / float64(r.Total)) * 100)
	if math.IsNaN(t) {
		t = 0
	}

	embed := NewEmbed().
		SetTitle("Your result:").
		SetDescription(fmt.Sprintf("Game in channel '%v'.", channelName)).
		AddField(html.UnescapeString("\u200b"), fmt.Sprintf("easy: %v / %v ( %.2f%% )", r.Easy.Correct, r.Easy.Total, e)).
		AddField(html.UnescapeString("\u200b"), fmt.Sprintf("medium: %v / %v ( %.2f%% )", r.Medium.Correct, r.Medium.Total, m)).
		AddField(html.UnescapeString("\u200b"), fmt.Sprintf("hard: %v / %v ( %.2f%% )", r.Hard.Correct, r.Hard.Total, h)).
		AddField(html.UnescapeString("\u200b"), fmt.Sprintf("total: %v / %v ( %.2f%% )", r.Correct, r.Total, t)).
		SetColor(0x00ffff).MessageEmbed

	return embed
}

func TestEmbed() *discordgo.MessageEmbed {
	/*
		embed := NewEmbed().
			SetTitle("I am an embed").
			SetDescription("This is a discordgo embed").
			AddField("I am a field", "I am a value").
			AddField("I am a second field", "I am a value").
			SetImage("https://cdn.discordapp.com/avatars/119249192806776836/cc32c5c3ee602e1fe252f9f595f9010e.jpg?size=2048").
			SetThumbnail("https://cdn.discordapp.com/avatars/119249192806776836/cc32c5c3ee602e1fe252f9f595f9010e.jpg?size=2048").
			SetColor(0x00ff00).MessageEmbed
	*/

	embed := NewEmbed().
		SetTitle("Usage").
		AddField("Command: !start <arg1 arg2...> : starts a new game in the guild channel you type into.", "example: !start, !start -a=15 -c=21 -d=easy").
		AddField("-c, --category", "Set the category of the questions. Not using this defaults to ALL categories. To see the category IDs, use !categories. Example: !start -c=21 : setscategory with ID 21, which, at the time of writing this, happens to be 'Sports' "+html.UnescapeString("&#128512;.")).
		AddField("-a, --amount", "Set the amount of questions the game will have. Not using this defaults to 5. Example: !start --amount=10 : the game will have 10 questions.").
		AddField("-d, --difficulty", "Sets the difficulty of the questions. Possible values: easy|medium|hard. Not using this defaults to ALL difficulties. Example: !start -d=easy.").
		AddField("-t, --time", "Sets the time (in seconds) players have to answer each question. Not using this defaults to 25 seconds.").
		AddField("Command: !categories : Send this (also usable in DM) to get the list of the question categories available and their IDs in a DM replay.", "Example: !categories "+html.UnescapeString("&#128518;.")).
		AddField("Command: !result : Show the standings for the current/last game.", "Example: !result :).").
		AddField("Command: !myResult : Shows you details about your result in the current/last game in a DM replay.", "Example: !myResult "+html.UnescapeString("&#129315;.")).
		SetColor(0xffffff).MessageEmbed

	return embed
}
