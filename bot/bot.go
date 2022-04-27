package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/laskolaskov/quiz-o-bot/print"
	"github.com/laskolaskov/quiz-o-bot/storage"
)

func MessageCreateListener(s *discordgo.Session, m *discordgo.MessageCreate) {
	//do not process self authored messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	command, args := process(m.Content)

	if isDM, _ := isDM(s, m); isDM {
		switch command {
		case "!categories":
			replay(s, m, print.ListCategories(storage.Categories()))
		case "!test":
			testAction(s, m)
		default:
			replayEmbed(s, m, print.Help())
		}
	} else {
		switch command {
		case "!categories":
			replay(s, m, print.ListCategories(storage.Categories()))
		case "!start":
			start(s, m, args)
		case "!result":
			result(s, m)
		case "!myResult":
			myResult(s, m)
		default:
			q, a, err := checkAnswer(m.Content)
			if err != nil {
				log.Printf("\nUnknown command - CMD: %v ARGS: %v\n\n", command, args)
				return
			}
			processAnswer(s, m, q, a)
		}
	}
}

func GuildCreateListener(s *discordgo.Session, event *discordgo.GuildCreate) {
	/* if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if isTextChannel(channel) {
			_, err := s.ChannelMessageSend(channel.ID, "Quiz-o-bot is ready for some trivia games! Send it a DM (with any content) to see the instructions about how to start one.\nMany thanks to https://opentdb.com/ for the great free database.")
			if err != nil {
				log.Println(err)
			}
		}
	} */
}
