package incanGold

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var 遊戲狀態 int = 0
var 遊戲狀態表 map[int]string = map[int]string{
	0: "準備狀態",
	1: "報名",
	2: "回合抽牌",
	3: "探險動作",
	4: "回合結算",
	5: "神殿結算",
	6: "探險總結算",
	7: "重置遊戲",
}

type 探險者 struct {
	帳篷   int
	口袋   int
	探險狀態 int
	神器背包 int
	死亡次數 int
}

var 探險隊 map[string]探險者 = make(map[string]探險者)
var 探險者狀態表 map[int]string = map[int]string{
	0: "準備狀態",
	1: "待選擇",
	2: "繼續探險",
	3: "打包回家",
}

var 當前互動訊息 string

func 準備印加寶藏(d *DiscordBotService, 神殿數 int) {

	探險隊 = map[string]探險者{}
	回合卡池 = []string{}

	// 初始化牌組
	for 單一災難 := range 災難 {
		回合卡池 = append(回合卡池, 單一災難)
	}
	for 單一寶藏 := range 寶藏 {
		回合卡池 = append(回合卡池, 單一寶藏)
	}

	神殿 = 1
	神殿總回合數 = 神殿數
	神器庫存 = 神器
	回合卡池 = append(回合卡池, 抽神器())

	回合 = 0
	遊戲狀態 = 1
	是否神殿結算 = false
	是否探險總結算 = false
	回合台面 = []string{}

	發送訊息到頻道("發現藏寶圖，誰印啦?  🗺")
}

func 抽神器() (選中神器 string) {
	var 抽第幾張 int = 隨機數字(1, 7-神殿)
	var 序數 int = 1
	var 神器暫存 map[string]int = make(map[string]int)
	for 單一神器 := range 神器庫存 {
		if 序數 == 抽第幾張 {
			選中神器 = 單一神器
		} else {
			神器暫存[單一神器] = 神器[單一神器]
		}
		序數++
	}

	神器庫存 = 神器暫存
	return
}

func 報名參加(探險者名稱 string) (報名結果 string) {

	if _, 已報名 := 探險隊[探險者名稱]; 已報名 {
		報名結果 = 探險者名稱 + "重複報名了  💦\n"
		return
	}

	if 遊戲狀態 == 1 {
		探險隊[探險者名稱] = 探險者{
			帳篷:   0,
			口袋:   0,
			探險狀態: 0,
			神器背包: 0,
			死亡次數: 0,
		}
		報名結果 = "貪婪的探險者 " + 探險者名稱 + " ! 🧙‍♂️"
	} else {
		報名結果 = "報名失敗，目前遊戲狀態:" + 遊戲狀態表[遊戲狀態] + " 💦\n"
	}
	return
}

func 探險者動作(探險者名稱, 動作 string) (探險者動作錯誤 string) {
	探險者動作錯誤 = ""
	if 遊戲狀態 != 3 {
		探險者動作錯誤 = "動作失敗，目前遊戲狀態:" + 遊戲狀態表[遊戲狀態] + "  💦\n"
		log.Error(探險者名稱, " ", 探險者動作錯誤)
		return
	}

	if 探險隊[探險者名稱].探險狀態 == 1 {
		var 探險者異動 探險者 = 探險隊[探險者名稱]
		if 動作 == "探險" {
			探險者異動.探險狀態 = 2
			探險中人數++
		} else { //"撤退"
			探險者異動.探險狀態 = 3
			回合撤退人數++
		}

		探險隊[探險者名稱] = 探險者異動
		log.Debug("收到動作: ", 探險者名稱, " ", 動作, " ", 探險隊[探險者名稱])
		已出牌人數++
	} else {
		探險者動作錯誤 = "動作失敗，探險者-探險狀態:" + 探險者狀態表[探險隊[探險者名稱].探險狀態] + "  💦\n"
		log.Error(探險者名稱, " ", 探險者動作錯誤)
		return
	}

	if 已出牌人數 == 本次探險人數 {
		是否需要自動出牌 = false
		編輯互動訊息(fmt.Sprintf("已收到 %s (%s)", 計算等待人數(), 探險者名稱), true)
		go 回合結算()
	} else {
		是否需要自動出牌 = true
		等待時間 = 5
		編輯互動訊息(fmt.Sprintf("已收到 %s (%s)", 計算等待人數(), 探險者名稱), false)
	}

	return
}

func 計算等待人數() (等待人數 string) {
	等待人數 = fmt.Sprint(已出牌人數, "/", 本次探險人數)
	log.Debug("回傳等待人數: ", 等待人數)
	return
}

// 回合變數
var 回合 int = 0
var 已出牌人數 int = 0
var 本次探險人數 int = 0
var 探險中人數 int = 0
var 回合撤退人數 int = 0
var 場上閒置鑽石 int = 0
var 移除災難 string
var 場上閒置神器 string
var 回合卡池 []string = make([]string, 0)
var 回合台面 []string = make([]string, 0)
var 是否神殿結算 bool

func 回合初始化() {
	// 遊戲狀態 = 1
	if len(探險隊) == 0 {
		發送訊息到頻道("無人參加")
		return
	}

	if 回合 == 0 {
		探險中人數 = len(探險隊)
	}

	本次探險人數 = 探險中人數

	抽出卡號 := 隨機數字(0, len(回合卡池))
	抽出卡片 := 回合卡池[抽出卡號]

	log.Debug("===回合初始化=== ", 回合)
	log.Debugf("回合卡池: %d, 卡號: %d", len(回合卡池), 抽出卡號)
	log.Debug(回合卡池)
	鑽石回報 := ""
	// 公告抽中的牌 分給幾個人 場上剩餘幾顆
	if strings.Contains(抽出卡片, "寶藏") || strings.Contains(抽出卡片, "神器") {
		if 神器鑽石數, 有神器 := 神器[抽出卡片]; 有神器 {

			if 探險中人數 == 1 {
				// 找出是誰Solo拿走神器
				for 探險者名稱 := range 探險隊 {
					var 探險者異動 探險者 = 探險隊[探險者名稱]
					if 探險者異動.探險狀態 == 1 || 回合 == 0 {
						探險者異動.探險狀態 = 1
						探險者異動.神器背包 = 神器鑽石數
						探險隊[探險者名稱] = 探險者異動
					}
				}
			} else {
				for 探險者名稱 := range 探險隊 {
					var 探險者異動 探險者 = 探險隊[探險者名稱]
					if 探險者異動.探險狀態 == 1 || 回合 == 0 {
						探險者異動.探險狀態 = 1
						探險隊[探險者名稱] = 探險者異動
					}
				}
				場上閒置神器 = 抽出卡片
			}

		}

		if 寶藏鑽石數, 有寶藏 := 寶藏[抽出卡片]; 有寶藏 {
			回合總鑽石數 := 場上閒置鑽石 + 寶藏鑽石數
			每人可分鑽石 := 回合總鑽石數 / 探險中人數

			for 探險者名稱 := range 探險隊 {
				var 探險者異動 探險者 = 探險隊[探險者名稱]
				if 探險者異動.探險狀態 == 1 || 回合 == 0 {
					探險者異動.探險狀態 = 1
					探險者異動.口袋 = 探險者異動.口袋 + 每人可分鑽石
					探險隊[探險者名稱] = 探險者異動

					回合總鑽石數 = 回合總鑽石數 - 每人可分鑽石
				}
			}
			場上閒置鑽石 = 回合總鑽石數
			鑽石回報 = fmt.Sprintf("共 %d 人平分, 桌面剩 %d 鑽\n", 探險中人數, 場上閒置鑽石)
		}
	} else {
		for 探險者名稱 := range 探險隊 {
			var 探險者異動 探險者 = 探險隊[探險者名稱]
			if 探險者異動.探險狀態 == 1 || 回合 == 0 {
				探險者異動.探險狀態 = 1
				探險隊[探險者名稱] = 探險者異動
			}
		}

		// 災難檢查
		災難類型 := string([]rune(抽出卡片)[:2])
		for 檢查序數 := range 回合台面 {
			if strings.Contains(回合台面[檢查序數], 災難類型) {
				鑽石回報 = "💥這個穴不行了🔥\n"

				// 下局移除此災難
				移除災難 = 抽出卡片
				是否神殿結算 = true
				break
			}
		}

		if !是否神殿結算 {
			鑽石回報 = fmt.Sprintf("桌面剩 %d 鑽\n", 場上閒置鑽石)
		}
	}

	// 整理卡池
	回合台面 = append(回合台面, 抽出卡片)
	var 卡池異動 []string = make([]string, 0)

	for 檢查序數 := range 回合卡池 {
		if 回合卡池[檢查序數] != 抽出卡片 {
			卡池異動 = append(卡池異動, 回合卡池[檢查序數])
		}
	}
	回合卡池 = 卡池異動
	探險中人數 = 0
	遊戲狀態 = 3

	// 戰況回報
	log.Debug("回合台面", 回合台面)
	log.Debugf("探險隊 %+v", 探險隊)
	if 場上閒置神器 != "" {
		鑽石回報 = 鑽石回報 + "場上神器: " + 場上閒置神器 + "\n"
	}

	if 是否神殿結算 {
		發送按鈕訊息到頻道(抽出卡片, strings.Join(回合台面, " "), 鑽石回報, true)
		神殿結算()
	} else {
		發送按鈕訊息到頻道(抽出卡片, strings.Join(回合台面, " "), 鑽石回報, false)
	}
}

func 回合結算() {
	// 探險隊結算
	log.Debug("===回合結算=== ", 回合)
	log.Debug("探險中人數:", 探險中人數, "回合撤退人數", 回合撤退人數)
	遊戲狀態 = 4

	var 繼續探險名單 []string = make([]string, 0)
	var 打包回家名單 []string = make([]string, 0)
	var 撤退可分鑽石 int
	if 回合撤退人數 > 0 {
		撤退可分鑽石 = 場上閒置鑽石 / 回合撤退人數
		場上閒置鑽石 = 場上閒置鑽石 - (撤退可分鑽石 * 回合撤退人數)
	}

	for 探險者名稱 := range 探險隊 {
		var 探險者異動 探險者 = 探險隊[探險者名稱]
		switch 探險者異動.探險狀態 {
		case 2: // 繼續探險
			{
				探險者異動.探險狀態 = 1
				探險隊[探險者名稱] = 探險者異動
				繼續探險名單 = append(繼續探險名單, fmt.Sprintf("%s (%d)", 探險者名稱, 探險隊[探險者名稱].口袋))
			}
		case 3: // 打包回家
			{

				if 回合撤退人數 == 1 && 場上閒置神器 != "" {
					探險者異動.神器背包 = 神器[場上閒置神器]
					場上閒置神器 = ""
				}

				var 本次收穫鑽石 int = 0
				if 探險者異動.神器背包 > 0 {
					本次收穫鑽石 = 本次收穫鑽石 + 探險者異動.神器背包
					探險者異動.神器背包 = 0
				}

				本次收穫鑽石 = 本次收穫鑽石 + 探險者異動.口袋 + 撤退可分鑽石
				探險者異動.口袋 = 0

				探險者異動.帳篷 = 探險者異動.帳篷 + 本次收穫鑽石
				探險者異動.探險狀態 = 0
				探險隊[探險者名稱] = 探險者異動
				回合撤退人數--
				打包回家名單 = append(打包回家名單, fmt.Sprintf("%s (%d)", 探險者名稱, 本次收穫鑽石))
			}
		}
	}

	回合++
	已出牌人數 = 0
	log.Debug("探險隊結算", 探險隊)

	if 探險中人數 == 0 {
		神殿結算()
	} else {
		time.Sleep(time.Second * 1)
		回合回報 := "回合結算中 🎲\n"
		if len(繼續探險名單) > 0 {
			回合回報 = 回合回報 + fmt.Sprintf("繼續探險: %s\n", strings.Join(繼續探險名單, ", "))
		}
		if len(打包回家名單) > 0 {
			回合回報 = 回合回報 + fmt.Sprintf("打包回家: %s\n", strings.Join(打包回家名單, ", "))
		}

		編輯互動訊息(回合回報, true)

		time.Sleep(time.Second * 3)
		回合初始化()
	}
}

// 神殿變數
var 神殿 int = 0
var 神殿總回合數 int = 0
var 神器庫存 map[string]int = make(map[string]int)

func 神殿初始化() {
	log.Debug("===神殿初始化=== ", 神殿)

	// 重建回合卡池
	for 單一災難 := range 災難 {
		if 單一災難 != 移除災難 {
			回合卡池 = append(回合卡池, 單一災難)
		}
	}
	for 單一寶藏 := range 寶藏 {
		回合卡池 = append(回合卡池, 單一寶藏)
	}
	回合卡池 = append(回合卡池, 抽神器())

	回合台面 = []string{}

	回合初始化()
}

func 神殿結算() {
	log.Debug("===神殿結算=== ", 神殿)
	遊戲狀態 = 5

	// 清空回合數據
	for 探險者名稱 := range 探險隊 {
		var 探險者異動 探險者 = 探險隊[探險者名稱]
		if 探險者異動.探險狀態 == 1 {
			探險者異動.死亡次數++
		}
		探險者異動.口袋 = 0
		探險者異動.探險狀態 = 0
		探險者異動.神器背包 = 0
		探險隊[探險者名稱] = 探險者異動
	}

	回合 = 0
	場上閒置鑽石 = 0
	場上閒置神器 = ""
	回合卡池 = []string{}
	是否神殿結算 = false

	if 神殿 == 神殿總回合數 {
		編輯互動訊息(fmt.Sprint("神殿", 神殿, "結束，探險結算中 🎲"), true)
		time.Sleep(time.Second * 3)
		探險總結算()
	} else {
		神殿++
		編輯互動訊息(fmt.Sprint("神殿", 神殿, "初始化中 🎲"), true)
		time.Sleep(time.Second * 5)
		神殿初始化()
	}
}

var 是否探險總結算 bool = false

func 探險總結算() {
	log.Debug("===探險總結算===")
	遊戲狀態 = 6
	是否探險總結算 = true

	var 結算表 map[string]int = make(map[string]int, 0)
	for 探險者名稱 := range 探險隊 {
		結算表[探險者名稱] = 探險隊[探險者名稱].帳篷
	}

	排序名單 := make([]string, 0, len(結算表))

	for 探險者名稱 := range 結算表 {
		排序名單 = append(排序名單, 探險者名稱)
	}
	sort.SliceStable(排序名單, func(i, j int) bool {
		return 結算表[排序名單[i]] > 結算表[排序名單[j]]
	})

	var 探險結果 string = "🥇探險總結算🥇\n"
	for 排名, 探險者名稱 := range 排序名單 {
		探險結果 = 探險結果 + fmt.Sprintf("第%d名 %s: %d \n", 排名+1, 探險者名稱, 結算表[探險者名稱])
	}

	死亡領主, 死亡次數 := 計算死亡領主()
	探險結果 = 探險結果 + fmt.Sprintf("死亡領主: %s, 死%d次 \n", 死亡領主, 死亡次數)
	探險結果 = 探險結果 + "遊戲結束 🤖"
	發送訊息到頻道(探險結果)
}

func 計算死亡領主() (死亡領主 string, 次數 int) {

	var 結算表 map[string]int = make(map[string]int, 0)
	for 探險者名稱 := range 探險隊 {
		結算表[探險者名稱] = 探險隊[探險者名稱].死亡次數
	}

	排序名單 := make([]string, 0, len(結算表))

	for 探險者名稱 := range 結算表 {
		排序名單 = append(排序名單, 探險者名稱)
	}
	sort.SliceStable(排序名單, func(i, j int) bool {
		return 結算表[排序名單[i]] > 結算表[排序名單[j]]
	})

	死亡領主 = ""
	for 排名, 探險者名稱 := range 排序名單 {
		if 排名 == 0 {
			死亡領主 = 死亡領主 + 探險者名稱
			次數 = 探險隊[探險者名稱].死亡次數
		} else if 探險隊[排序名單[排名]].死亡次數 == 探險隊[排序名單[排名-1]].死亡次數 {
			// 與前一名同次數
			死亡領主 = 死亡領主 + ", " + 探險者名稱
		} else if 探險隊[排序名單[排名]].死亡次數 != 探險隊[排序名單[排名-1]].死亡次數 {
			return
		}
	}
	return
}

func 隨機數字(最小, 最大 int) int {
	return 最小 + rand.Intn(最大-最小)
}

func 發送訊息到頻道(內容 string) {
	DiscordBotSrv.SendMsgToDiscord(內容)
}

func 發送按鈕訊息到頻道(抽牌, 桌面, 鑽石 string, 是否移除按鈕 bool) {
	DiscordBotSrv.SendButtonMsgToDiscord(抽牌, 桌面, 鑽石, 是否移除按鈕)
}

func 編輯互動訊息(內容 string, 是否移除按鈕 bool) {
	DiscordBotSrv.EditButtonMsg(內容, 是否移除按鈕)
}
