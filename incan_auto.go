package incanGold

import (
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

// 等待時間倒數
var 等待時間 int = 0
var 是否需要自動出牌 bool

func init() {
	背景倒數()
}

func 背景倒數() {
	c := cron.New()
	c.AddFunc("0 */1 * * *", func() {
		if 等待時間 == 0 && 是否需要自動出牌 {
			自動出牌()
		}
		if 等待時間 > 0 {
			等待時間--
		}
	})
	c.Start()
}

func 自動出牌() {
	var 自動出牌報告 string = "自動出牌報告:\n"
	for 探險者名稱 := range 探險隊 {
		if 探險隊[探險者名稱].冒險狀態 == 1 {
			var 探險者異動 探險者 = 探險隊[探險者名稱]
			var 動作 string
			if 隨機數字(1, 11)%2 == 0 {
				探險者異動.冒險狀態 = 2
				探險中人數++
				動作 = "探險"
			} else {
				探險者異動.冒險狀態 = 3
				回合撤退人數++
				動作 = "撤退"
			}
			探險隊[探險者名稱] = 探險者異動
			log.Info("自動出牌: ", 探險者名稱, " ", 動作, " ", 探險隊[探險者名稱])
			自動出牌報告 = 自動出牌報告 + 探險者名稱 + ": " + 動作 + "\n"
			已出牌人數++
		}
	}

	DiscordBotSrv.SendMsgToDiscord(自動出牌報告)
	是否需要自動出牌 = false
	go 回合結算()
}
