# go_discord_incan_gold
- 整理之前寫的桌遊印加寶藏 discord版
- 就覺得用中文寫很酷
- 嘗試改版成可獨立運作

## local run
- 先有discord頻道
- 再有機器人
- 機器人加入頻道
- 新增.\cmd\\.env
- go run .\cmd\main.go

## 現有功能
- "!說明印啦" => "在頻道說明功能"
- "!準備印啦" => "印加寶藏-初始化 可帶回合數(!準備印啦3)，預設5"
- "!印啦" =>     "印加寶藏-參加"
- "!開始印啦" => "印加寶藏-開始"
- "!等等印啦" => "儲存遊戲狀態"
- "!繼續印啦" => "讀取遊戲狀態"
- "!印到底啦" => "直接跳到總結算"

## TODO
- 補文件說明