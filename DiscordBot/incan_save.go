package DiscordBot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type IncanSave struct {
	GameStatus             int                        `json:"遊戲狀態"`
	Temple                 int                        `json:"神殿"`
	ActedNum               int                        `json:"已出牌人數"`
	AdventureNumThisTime   int                        `json:"本次探險人數"`
	AdventureNumInProgress int                        `json:"探險中人數"`
	IdleDiamond            int                        `json:"場上閒置鑽石"`
	RemoveDisaster         string                     `json:"移除災難"`
	IdleTool               string                     `json:"場上閒置神器"`
	RoundCardPool          []string                   `json:"回合卡池"`
	RoundCardTable         []string                   `json:"回合台面"`
	StockTool              map[string]int             `json:"神器庫存"`
	IsTempleSettle         bool                       `json:"是否神殿結算"`
	IsAdventureSettle      bool                       `json:"是否探險總結算"`
	AdventureTeam          map[string]IncanSavePlayer `json:"探險隊"`
}
type IncanSavePlayer struct {
	Tent         int `json:"帳篷"`
	Pocket       int `json:"口袋"`
	Status       int `json:"冒險狀態"`
	ToolBackpack int `json:"神器背包"`
	Death        int `json:"死亡次數"`
}

func 讀取印加進度() {
	var 存檔 IncanSave
	var tempStr string
	// open the file
	file, err := os.Open("incanSave.txt")

	//handle errors while opening
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)

	// read line by line
	for fileScanner.Scan() {
		tempStr = tempStr + fileScanner.Text()
	}

	err = json.Unmarshal([]byte(tempStr), &存檔)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s", err)
	}

	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	遊戲狀態 = 存檔.GameStatus
	神殿 = 存檔.Temple
	已出牌人數 = 存檔.ActedNum
	本次探險人數 = 存檔.AdventureNumThisTime
	探險中人數 = 存檔.AdventureNumInProgress
	場上閒置鑽石 = 存檔.IdleDiamond
	移除災難 = 存檔.RemoveDisaster
	場上閒置神器 = 存檔.IdleTool
	回合卡池 = 存檔.RoundCardPool
	回合台面 = 存檔.RoundCardTable
	神器庫存 = 存檔.StockTool
	是否神殿結算 = 存檔.IsTempleSettle
	是否探險總結算 = 存檔.IsAdventureSettle

	for 探險者名稱 := range 存檔.AdventureTeam {
		var 暫存探險者 探險者
		暫存探險者.帳篷 = 存檔.AdventureTeam[探險者名稱].Tent
		暫存探險者.口袋 = 存檔.AdventureTeam[探險者名稱].Pocket
		暫存探險者.冒險狀態 = 存檔.AdventureTeam[探險者名稱].Status
		暫存探險者.神器背包 = 存檔.AdventureTeam[探險者名稱].ToolBackpack
		暫存探險者.死亡次數 = 存檔.AdventureTeam[探險者名稱].Death

		探險隊[探險者名稱] = 暫存探險者
	}

	file.Close()
}

func 儲存印加進度() {
	var 存檔 IncanSave
	var 探險隊存檔 map[string]IncanSavePlayer = make(map[string]IncanSavePlayer)
	存檔.GameStatus = 遊戲狀態
	存檔.Temple = 神殿
	存檔.ActedNum = 已出牌人數
	存檔.AdventureNumThisTime = 本次探險人數
	存檔.AdventureNumInProgress = 探險中人數
	存檔.IdleDiamond = 場上閒置鑽石
	存檔.RemoveDisaster = 移除災難
	存檔.IdleTool = 場上閒置神器
	存檔.RoundCardPool = 回合卡池
	存檔.RoundCardTable = 回合台面
	存檔.StockTool = 神器庫存
	存檔.IsTempleSettle = 是否神殿結算
	存檔.IsAdventureSettle = 是否探險總結算

	for 探險者名稱 := range 探險隊 {
		var 暫存探險者 IncanSavePlayer
		暫存探險者.Tent = 探險隊[探險者名稱].帳篷
		暫存探險者.Pocket = 探險隊[探險者名稱].口袋
		暫存探險者.Status = 探險隊[探險者名稱].冒險狀態
		暫存探險者.ToolBackpack = 探險隊[探險者名稱].神器背包
		暫存探險者.Death = 探險隊[探險者名稱].死亡次數
		探險隊存檔[探險者名稱] = 暫存探險者
	}

	存檔.AdventureTeam = 探險隊存檔
	msgJSON, err := json.Marshal(存檔)
	if err != nil {
		log.Printf("儲存印加進度 Marshal失敗 err:\n %v", err)
		return
	}
	writeErr := os.WriteFile("incanSave.txt", msgJSON, 0644)
	if writeErr != nil {
		log.Printf("儲存印加進度失敗 請備份內容如下:\n %v, writeErr:%v\n", 存檔, writeErr)
	} else {
		fmt.Println("儲存印加進度完成")
	}
}
