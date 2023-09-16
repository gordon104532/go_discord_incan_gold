package incanGold

import (
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

	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"fd_forward": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			探險者動作(i.Member.User.Username, "探險")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
			})
		},
		"fd_retreat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			探險者動作(i.Member.User.Username, "撤退")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
			})
		},
	}
)
