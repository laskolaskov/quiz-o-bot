package main

import (
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

	godotenv.Overload()

	//create discord connection
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	defer discord.Close()

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
	log.Println("Quiz-o-bot started! Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
