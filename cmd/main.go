package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	incan "github.com/gordon104532/go_discord_incan_gold"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Server Start")

	err := godotenv.Load("./cmd/.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	applicationID := os.Getenv("APPLICATION_ID")
	guildID := os.Getenv("GUILD_ID")
	textChannelID := os.Getenv("TEXT_CHANNEL_ID")

	discordBotService := incan.NewDiscordBotService(botToken, applicationID, guildID, textChannelID)
	discordBotService.Run()

	wg := &sync.WaitGroup{}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-c
		_ = sig
		wg.Done()
	}()
	wg.Add(1)
	wg.Wait()
	log.Println("Server Shutdown")
}
