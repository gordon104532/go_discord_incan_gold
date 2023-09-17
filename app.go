package incanGold

import (
	"math/rand"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

type DiscordBotService struct {
	dg              *discordgo.Session
	botToken        string
	applicationID   string
	guildID         string
	textChannelID   string
	isSessionImport bool
	debug           bool
}

var DiscordBotSrv *DiscordBotService
var discordFuncMap = map[string]string{
	"!印啦":   "印加寶藏-參加",
	"!說明印啦": "印加寶藏-初始化",
	"!準備印啦": "印加寶藏-初始化",
	"!開始印啦": "印加寶藏-開始",
	"!等等印啦": "儲存遊戲狀態",
	"!繼續印啦": "讀取遊戲狀態",
	"!印到底啦": "直接跳到總結算",
}

func NewDiscordBotService(dg *discordgo.Session, botToken, applicationID, guildID, textChannelID string, debug bool) *DiscordBotService {
	var isSessionImport bool
	if dg == nil {
		var err error
		dg, err = discordgo.New("Bot " + botToken)
		if err != nil {
			log.Fatal("DiscordBot new session error")
		}
		isSessionImport = true
	}

	DiscordBotSrv = &DiscordBotService{
		dg:              dg,
		botToken:        botToken,
		applicationID:   applicationID,
		guildID:         guildID,
		textChannelID:   textChannelID,
		isSessionImport: isSessionImport,
		debug:           debug,
	}

	return DiscordBotSrv
}

func (d *DiscordBotService) Run() {
	rand.Seed(time.Now().UnixNano())
	if d.debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.Println("DiscordBot Init")

	// Register the messageCreate func as a callback for MessageCreate events.
	d.dg.AddHandler(d.discordMessageHandle)

	d.dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionMessageComponent:
			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		default:
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "指令未實作",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	})

	d.dg.Identify.Intents = discordgo.IntentsGuildMessages

	if d.isSessionImport {
		// 開啟連線
		err := d.dg.Open()
		if err != nil {
			log.Fatal("DiscordBot error opening connection,", err)
			return
		}

		// Cleanly close down the Discord session.
		defer d.dg.Close()
	}
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

	var content string // 回傳的訊息內容

	if strings.Contains(m.Content, "!") {
		switch m.Content {
		case "!說明印啦":
			for v, i := range discordFuncMap {
				content = content + v + " : " + i + "\n"
			}
		case "!準備印啦":
			準備印加寶藏(&d)
		case "!印啦":
			content = 報名參加(m.Author.Username)
		case "!開始印啦":
			回合初始化()
		case "!等等印啦":
			儲存印加進度()
			content = "儲存印加進度"
		case "!繼續印啦":
			讀取印加進度()
			content = "讀取印加進度"
		case "!印到底啦":
			探險總結算()
		}
	}

	if len(content) > 0 {
		// 送出訊息
		d.SendMsgToDiscord(content)
	}
}

// 外部訊息傳入discord
func (d DiscordBotService) SendMsgToDiscord(content string) {
	d.dg.ChannelMessageSend(d.textChannelID, content)
}

// 外部訊息傳入discord
func (d DiscordBotService) SendButtonMsgToDiscord(content string) {
	d.dg.ChannelMessageSendComplex(d.textChannelID, &discordgo.MessageSend{
		Content:    content,
		Components: buttonComponent,
	})
}
