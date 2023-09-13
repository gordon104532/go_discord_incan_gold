package DiscordBot

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DiscordBotService struct {
	dg            *discordgo.Session
	botToken      string
	applicationID string
	guildID       string
	textChannelID string
}

var discordFuncMap = map[string]string{
	"!準備印加": "印加寶藏-初始化",
	"!印啦":   "印加寶藏-參加",
	"!準備印啦": "印加寶藏-開始",
	"!等等印啦": "儲存遊戲狀態",
	"!繼續印啦": "讀取遊戲狀態",
	"!印到底啦": "直接跳到總結算",
}

func NewDiscordBotService(botToken, applicationID, guildID, textChannelID string) *DiscordBotService {
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal("DiscordBot new session error")
	}

	return &DiscordBotService{
		dg:            dg,
		botToken:      botToken,
		applicationID: applicationID,
		guildID:       guildID,
		textChannelID: textChannelID,
	}
}

func (d DiscordBotService) Run() {
	rand.Seed(time.Now().UnixNano())
	log.Println("DiscordBot Init")

	log.Println("DiscordBot Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := d.dg.ApplicationCommandCreate(d.applicationID, d.guildID, v)
		if err != nil {
			log.Printf("DiscordBot Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	d.dg.AddHandler(d.discordMessageHandle)

	d.dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	d.dg.Identify.Intents = discordgo.IntentsGuildMessages

	// 開啟連線
	err := d.dg.Open()
	if err != nil {
		log.Fatal("DiscordBot error opening connection,", err)
		return
	}

	// Cleanly close down the Discord session.
	defer d.dg.Close()
}

// 接收訊息與處理
func (d DiscordBotService) discordMessageHandle(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 跳過機器人自身發話
	if m.Author.ID == s.State.User.ID {
		return
	}

	// 指定頻道才繼續
	if m.ChannelID != d.textChannelID {
		return
	}

	var context string // 回傳的訊息內容

	if strings.Contains(m.Content, "!") {
		switch m.Content {
		case "!h":
			for v, i := range discordFuncMap {
				context = context + v + " : " + i + "\n"
			}
		case "!準備印啦":
			準備印加寶藏(&d)
		case "!印啦":
			context = 報名參加(m.Author.Username)
		case "!開始印啦":
			回合初始化()
		case "!等等印啦":
			儲存印加進度()
			context = "儲存印加進度"
		case "!繼續印啦":
			讀取印加進度()
			context = "讀取印加進度"
		case "!印到底啦":
			探險總結算()
		}
	}

	if len(context) > 0 {
		// 送出訊息
		d.SendMsgToDiscord(context)
	}
}

// 外部訊息傳入discord
func (d DiscordBotService) SendMsgToDiscord(context string) {
	d.dg.ChannelMessageSend(d.textChannelID, context)
}
