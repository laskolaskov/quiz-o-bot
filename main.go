package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/laskolaskov/quiz-o-bot/api"
	"github.com/laskolaskov/quiz-o-bot/bot"
	"github.com/laskolaskov/quiz-o-bot/storage"
)

func main() {

	godotenv.Load()

	//create discord connection
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	//load data
	categories, err := api.GetCategories()
	if err != nil {
		log.Fatal(err)
	}
	storage.SetCategories(categories)

	//handlers
	discord.AddHandler(bot.MessageCreateListener)
	discord.AddHandler(bot.GuildCreateListener)

	//open websocket
	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Quiz-o-bot started! Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()

	//apiArgs := api.ParseInput(os.Args[1:len(os.Args)])

	/* if apiArgs.Help {
		fmt.Print(prn.Help())
		os.Exit(0)
	} */

	/* if apiArgs.ListCategories {
		fmt.Print(prn.ListCategories(categories))
		os.Exit(0)
	} */

	/* questions, err := api.GetQuestions(apiArgs.Amount, apiArgs.Category, apiArgs.Difficulty)
	if err != nil {
		log.Fatal(err)
	} */

	/* for k, c := range categories {
		//fmt.Println(c.Id, c.Name)
		fmt.Sprintln(k, c)
	} */

	/* for _, q := range questions {
		utils.PrintStruct(q)
	} */
}
