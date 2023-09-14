package incanGold

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "incan-gold",
			Description: "Incan Gold command",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "action",
					Description: "可選動作",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "forward",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.ChineseCN: "探險",
							},
							Value: 1,
						},
						{
							Name: "retreat",
							NameLocalizations: map[discordgo.Locale]string{
								discordgo.ChineseCN: "撤退",
							},
							Value: 2,
						},
					},
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"incan-gold": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			context := "已收到"
			if option, ok := optionMap["action"]; ok {
				switch option.IntValue() {
				case 1:
					探險者動作(i.Member.User.Username, "探險")
				case 2:
					探險者動作(i.Member.User.Username, "撤退")
				default:
					context = context + "無效的指令"
				}

				if 探險者動作錯誤 == "" {
					context = context + 計算等待人數()
				} else {
					context = 探險者動作錯誤
				}
			} else {
				context = context + "無效的指令"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: context,
				},
			})
		},
	}
)

func (d DiscordBotService) MakeGuildCommand() {
	url := fmt.Sprintf("https://discord.com/api/v10/applications/%s/guilds/%s/commands", d.applicationID, d.guildID)

	command := discordgo.ApplicationCommand{}

	body, _ := json.Marshal(command)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Discord MakeGuildCommand NewRequest err:", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bot "+d.botToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Discord MakeGuildCommand DefaultClient Do err:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Println("Discord MakeGuildCommand http response code : ", resp.StatusCode, "\n body: ", string(bodyBytes))
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Discord MakeGuildCommand Api ioutil.ReadAll err:", err)
		return
	}

	fmt.Println(string(bodyBytes))
}

func (d DiscordBotService) DeleteGuildCommand(commandID string) {
	url := fmt.Sprintf("https://discord.com/api/v10/applications/%s/guilds/%s/commands/%s", d.applicationID, d.guildID, commandID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Println("Discord DeleteGuildCommand NewRequest err:", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bot "+d.botToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Discord DeleteGuildCommand DefaultClient Do err:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Println("Discord DeleteGuildCommand http response code : ", resp.StatusCode, "\n body: ", string(bodyBytes))
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Discord DeleteGuildCommand Api ioutil.ReadAll err:", err)
		return
	}

	fmt.Println(string(bodyBytes))
}
