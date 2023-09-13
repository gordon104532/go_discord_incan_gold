package main

import (
	"fmt"
	"log"
	"main/DiscordBot"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Server Start")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("BOT_TOKEN")
	applicationID := os.Getenv("APPLICATION_ID")
	guildID := os.Getenv("GUILD_ID")
	textChannelID := os.Getenv("TEXT_CHANNEL_ID")

	fmt.Println("🎈", botToken, applicationID, guildID, textChannelID)
	discordBotService := DiscordBot.NewDiscordBotService(botToken, applicationID, guildID, textChannelID)
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
