package incanGold

import (
	"math/rand"
	"strconv"
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

		if strings.Contains(m.Content, "!準備印啦") {
			str := strings.Replace(m.Content, "!準備印啦", "", -1)
			rounds, err := strconv.Atoi(str)
			if err != nil {
				log.Error("!準備印啦 回合數 parse err:", err)
				rounds = 5
			}

			準備印加寶藏(&d, rounds)
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
func (d DiscordBotService) SendButtonMsgToDiscord(draw, table, diamond string, removeButton bool) {
	setComponent := buttonComponent
	if removeButton {
		setComponent = buttonComponentDisable
	}

	msg, err := d.dg.ChannelMessageSendComplex(d.textChannelID, &discordgo.MessageSend{
		Components: setComponent,
		Embed: &discordgo.MessageEmbed{
			Type:  "rich",
			Title: "本回合抽出: [" + draw + "]",
			Color: 2552136,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "桌面:", Value: "[" + table + "]", Inline: false},
				{Name: "鑽石:", Value: diamond, Inline: false},
				{Name: "收到:", Value: "", Inline: false},
			},
		},
	})

	if err != nil {
		log.Error("SendButtonMsg err:", err)
		return
	}

	當前互動訊息 = msg.ID
}

// 外部訊息傳入discord
func (d DiscordBotService) EditButtonMsg(addContent string, removeButton bool) {
	interactMsg, err := d.dg.ChannelMessage(d.textChannelID, 當前互動訊息)
	if err != nil {
		log.Error("EditButtonMsg get msg err:", err)
		return
	}

	setComponent := buttonComponent
	if removeButton {
		setComponent = buttonComponentDisable
	}

	if interactMsg.Embeds[0].Fields[2].Value != "" {
		addContent = "\n" + addContent
	}

	_, err = d.dg.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         interactMsg.ID,
		Channel:    interactMsg.ChannelID,
		Content:    &interactMsg.Content,
		Components: setComponent,
		Embed: &discordgo.MessageEmbed{
			Type:  interactMsg.Embeds[0].Type,
			Title: interactMsg.Embeds[0].Title,
			Color: interactMsg.Embeds[0].Color,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "桌面:", Value: interactMsg.Embeds[0].Fields[0].Value, Inline: false},
				{Name: "鑽石:", Value: interactMsg.Embeds[0].Fields[1].Value, Inline: false},
				{Name: "探險者", Value: interactMsg.Embeds[0].Fields[2].Value + addContent, Inline: false},
			},
		},
	})
	if err != nil {
		log.Error("EditButtonMsg EditComplex err:", err)
		return
	}
}
