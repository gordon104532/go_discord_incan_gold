package incanGold

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	buttonComponent = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "探險",
					Style:    discordgo.SuccessButton,
					Disabled: false,
					CustomID: "fd_forward",
				},
				discordgo.Button{
					Label:    "撤退",
					Style:    discordgo.DangerButton,
					Disabled: false,
					CustomID: "fd_retreat",
				},
			},
		},
	}

	buttonComponentDisable = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "探險",
					Style:    discordgo.SuccessButton,
					Disabled: true,
					CustomID: "fd_forward",
				},
				discordgo.Button{
					Label:    "撤退",
					Style:    discordgo.DangerButton,
					Disabled: true,
					CustomID: "fd_retreat",
				},
			},
		},
	}

	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"fd_forward": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			res := 探險者動作(i.Member.User.Username, "探險")
			var content string
			if res != "" {
				content = fmt.Sprintf("%s %s", i.Member.User.Username, res)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
		"fd_retreat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			res := 探險者動作(i.Member.User.Username, "撤退")
			var content string
			if res != "" {
				content = fmt.Sprintf("%s %s", i.Member.User.Username, res)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
	}
)
