package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	prn "github.com/laskolaskov/quiz-o-bot/print"
	"github.com/laskolaskov/quiz-o-bot/storage"
)

func MessageCreateListener(s *discordgo.Session, m *discordgo.MessageCreate) {
	//do not process self authored messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	t := m.ChannelID
	//var msg string
	/* for _, g := range g {
		//sGuild, _ := s.Guild(g.ID)
		//msg += sGuild.Name + "\n"
		fmt.Println(g.Name)
	} */
	fmt.Printf("%v %T\n", t, t)
	fmt.Println("test")

	command, args := process(m.Content)

	if isDM, _ := isDM(s, m); isDM {
		switch command {
		case "!categories":
			replay(s, m, prn.ListCategories(storage.Categories()))
		default:
			replay(s, m, prn.Help())
		}
	} else {
		switch command {
		case "!categories":
			replay(s, m, prn.ListCategories(storage.Categories()))
		case "!help":
			replay(s, m, prn.ListCategories(storage.Categories()))
		default:
			fmt.Printf("\nCMD: %v ARGS: %v\n", command, args)
		}
	}
}

func GuildCreateListener(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	fmt.Printf("%v %T", event.Guild.Name, event.Guild.Name)
	fmt.Println("")
	fmt.Printf("%v %T", event.Guild.ID, event.Guild.ID)
	fmt.Println("\n\nChannels:")
	fmt.Println("")

	for _, channel := range event.Guild.Channels {
		fmt.Printf("%v %T", channel.Name, channel.Name)
		fmt.Println("")
		fmt.Printf("%v %T", channel.ID, channel.ID)
		fmt.Println("")
		fmt.Printf("%v %T", channel.Type, channel.Type)
		fmt.Println("")
		fmt.Println("")

		if isTextChannel(channel) {
			_, err := s.ChannelMessageSend(channel.ID, "Quiz-o-bot is ready for some trivia games! Send it a DM, or type '!help' in the channel to see how to start one.\nMany thanks to https://opentdb.com/ for the great free database.")
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
