package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mainjzb/go.num/v2/zh"
	"github.com/mattn/go-ieproxy"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type QQBindCharacter struct {
	gorm.Model
	Character string
	QQ        int `gorm:"primarykey"`
	Lock      bool
}

type MapleInfo struct {
	Message           string  `json:"message"`
	AchievementPoints int     `json:"AchievementPoints"`
	AchievementRank   int     `json:"AchievementRank"`
	CharacterImageURL string  `json:"CharacterImageURL"`
	Class             string  `json:"Class"`
	ClassRank         int     `json:"ClassRank"`
	EXP               int64   `json:"EXP"`
	EXPPercent        float64 `json:"EXPPercent"`
	GlobalRanking     int     `json:"GlobalRanking"`
	GraphData         []struct {
		AvatarURL        string `json:"AvatarURL"`
		ClassID          int    `json:"ClassID"`
		ClassRankGroupID int    `json:"ClassRankGroupID"`
		CurrentEXP       int64  `json:"CurrentEXP"`
		DateLabel        string `json:"DateLabel"`
		EXPDifference    int64  `json:"EXPDifference"`
		EXPToNextLevel   int64  `json:"EXPToNextLevel"`
		ImportTime       int    `json:"ImportTime"`
		Level            int    `json:"Level"`
		Name             string `json:"Name"`
		ServerID         int    `json:"ServerID"`
		ServerMergeID    int    `json:"ServerMergeID"`
		TotalOverallEXP  int64  `json:"TotalOverallEXP"`
	} `json:"GraphData"`
	Guild              string `json:"Guild"`
	LegionCoinsPerDay  int    `json:"LegionCoinsPerDay"`
	LegionLevel        int    `json:"LegionLevel"`
	LegionPower        int    `json:"LegionPower"`
	LegionRank         int    `json:"LegionRank"`
	Level              int    `json:"Level"`
	Name               string `json:"Name"`
	Server             string `json:"Server"`
	ServerClassRanking int    `json:"ServerClassRanking"`
	ServerRank         int    `json:"ServerRank"`
	ServerSlug         string `json:"ServerSlug"`
}

type OfficialWorldRank []struct {
	Rank            int         `json:"Rank"`
	CharacterImgURL string      `json:"CharacterImgUrl"`
	PetImgURL       string      `json:"PetImgUrl"`
	WorldName       string      `json:"WorldName"`
	JobName         string      `json:"JobName"`
	IsSearchTarget  bool        `json:"IsSearchTarget"`
	TimeSum         int         `json:"TimeSum"`
	Stage           int         `json:"Stage"`
	StarSum         int         `json:"StarSum"`
	LegionLevel     int         `json:"LegionLevel"`
	RaidPower       int         `json:"RaidPower"`
	TierName        string      `json:"TierName"`
	GuildName       interface{} `json:"GuildName"`
	AccountID       interface{} `json:"AccountId"`
	Score           interface{} `json:"Score"`
	CombatPower     interface{} `json:"CombatPower"`
	MatchSN         interface{} `json:"MatchSN"`
	Percentile      interface{} `json:"Percentile"`
	SeasonNo        int         `json:"SeasonNo"`
	CharacterID     int         `json:"CharacterID"`
	CharacterName   string      `json:"CharacterName"`
	Exp             int         `json:"Exp"`
	Gap             int         `json:"Gap"`
	JobDetail       int         `json:"JobDetail"`
	JobID           int         `json:"JobID"`
	Level           int         `json:"Level"`
	StartRank       int         `json:"StartRank"`
	TranferStatus   int         `json:"TranferStatus"`
	WorldID         int         `json:"WorldID"`
}

type OfficialJobRank []struct {
	AccountID       interface{} `json:"AccountId"`
	CharacterID     int64       `json:"CharacterID"`
	CharacterImgURL string      `json:"CharacterImgUrl"`
	CharacterName   string      `json:"CharacterName"`
	CombatPower     interface{} `json:"CombatPower"`
	Exp             int64       `json:"Exp"`
	Gap             int64       `json:"Gap"`
	GuildName       interface{} `json:"GuildName"`
	IsSearchTarget  bool        `json:"IsSearchTarget"`
	JobDetail       int64       `json:"JobDetail"`
	JobID           int64       `json:"JobID"`
	JobName         string      `json:"JobName"`
	LegionLevel     int64       `json:"LegionLevel"`
	Level           int64       `json:"Level"`
	MatchSN         interface{} `json:"MatchSN"`
	Percentile      interface{} `json:"Percentile"`
	PetImgURL       string      `json:"PetImgUrl"`
	RaidPower       int64       `json:"RaidPower"`
	Rank            int64       `json:"Rank"`
	Score           interface{} `json:"Score"`
	SeasonNo        int64       `json:"SeasonNo"`
	Stage           int64       `json:"Stage"`
	StarSum         int64       `json:"StarSum"`
	StartRank       int64       `json:"StartRank"`
	TierName        string      `json:"TierName"`
	TimeSum         int64       `json:"TimeSum"`
	TranferStatus   int64       `json:"TranferStatus"`
	WorldID         int64       `json:"WorldID"`
	WorldName       string      `json:"WorldName"`
}

type OfficialLegionRank []struct {
	AccountID       interface{} `json:"AccountId"`
	CharacterID     int64       `json:"CharacterID"`
	CharacterImgURL string      `json:"CharacterImgUrl"`
	CharacterName   string      `json:"CharacterName"`
	CombatPower     interface{} `json:"CombatPower"`
	Exp             int64       `json:"Exp"`
	Gap             int64       `json:"Gap"`
	GuildName       interface{} `json:"GuildName"`
	IsSearchTarget  bool        `json:"IsSearchTarget"`
	JobDetail       int64       `json:"JobDetail"`
	JobID           int64       `json:"JobID"`
	JobName         string      `json:"JobName"`
	LegionLevel     int64       `json:"LegionLevel"`
	Level           int64       `json:"Level"`
	MatchSN         interface{} `json:"MatchSN"`
	Percentile      interface{} `json:"Percentile"`
	PetImgURL       string      `json:"PetImgUrl"`
	RaidPower       int64       `json:"RaidPower"`
	Rank            int64       `json:"Rank"`
	Score           interface{} `json:"Score"`
	SeasonNo        int64       `json:"SeasonNo"`
	Stage           int64       `json:"Stage"`
	StarSum         int64       `json:"StarSum"`
	StartRank       int64       `json:"StartRank"`
	TierName        string      `json:"TierName"`
	TimeSum         int64       `json:"TimeSum"`
	TranferStatus   int64       `json:"TranferStatus"`
	WorldID         int64       `json:"WorldID"`
	WorldName       string      `json:"WorldName"`
}

type CharInfoResult struct {
	imageURL, info string
}

var checkNumberOfTimes map[int]int
var charactersCatch map[string]CharInfoResult
var charactersLevelRank map[string]string
var catchTime int
var classLevelRank map[string]string

func CheckOfficialRank(name string, gropFromQQ int) (result CharInfoResult) {
	name = strings.TrimSpace(name)
	resetCacheEveryday()

	if checkNumberOfTimes[gropFromQQ] >= 12 {
		return CharInfoResult{"", "今日查询已达上限"}
	}

	// map缓存查找
	if chara, ok := charactersCatch[name]; ok {
		checkNumberOfTimes[gropFromQQ]++
		return chara
	}

	jobURL := "https://maplestory.nexon.net/api/ranking?id=job&id2=&rebootIndex=0&page_index=1&character_name=" + name
	body, err := webGetRequest(jobURL)
	if err != nil {
		return CharInfoResult{"", "查询失败"}
	}

	rankInfo := OfficialJobRank{}
	err = json.Unmarshal(body, &rankInfo)
	if err != nil {
		return CharInfoResult{"", "查询失败：网络错误"}
	}
	if len(rankInfo) == 0 {
		return CharInfoResult{"", "无此角色"}
	}

	result.imageURL = rankInfo[0].CharacterImgURL
	outName := strconv.QuoteToASCII(rankInfo[0].CharacterName)
	result.info += "角色:" + outName[1:len(outName)-1] + "  （" + serverName2Chinese[rankInfo[0].WorldName] + "）\n"
	exp := float64(rankInfo[0].Exp) / float64(LevelExp[rankInfo[0].Level]) * 100
	result.info += fmt.Sprintf("等级:%v - %.2f%% \n", rankInfo[0].Level, exp)

	rRankInfo := OfficialJobRank{}
	var jobName = fmt.Sprintf("%v%v", rankInfo[0].JobName, rankInfo[0].JobDetail)
	if jName, ok := class2Chinese[jobName]; ok {
		jobName = jName
	}
	if rankInfo[0].Level > 219 {
		var rRankURL = "https://maplestory.nexon.net/api/ranking?id=job&id2=&rebootIndex=2&page_index=1&character_name=" + name
		if rankInfo[0].WorldID == 45 {
			rRankURL = "https://maplestory.nexon.net/api/ranking?id=job&id2=&rebootIndex=1&page_index=1&character_name=" + name
		}

		body, err := webGetRequest(rRankURL)
		if err != nil {
			result.info += fmt.Sprintf("职业:%v（全排%v）\n\n", jobName, rankInfo[0].Rank)
		} else {
			err = json.Unmarshal(body, &rRankInfo)
			if err != nil || len(rRankInfo) == 0 {
				result.info += fmt.Sprintf("职业:%v（全排%v）\n\n", jobName, rankInfo[0].Rank)
			} else {
				result.info += fmt.Sprintf("职业:%v（单排%v 全排%v）\n\n", jobName, rRankInfo[0].Rank, rankInfo[0].Rank)
			}
		}
	} else {
		result.info += fmt.Sprintf("职业:%v（全排%v）\n\n", jobName, rankInfo[0].Rank)
	}

	legionURL := "https://maplestory.nexon.net/api/ranking?id=legion&id2=" + fmt.Sprintf("%v", rankInfo[0].WorldID) + "&page_index=1&character_name=" + name
	legionByte, _ := webGetRequest(legionURL)
	legion := OfficialLegionRank{}
	if err := json.Unmarshal(legionByte, &legion); err == nil && len(legion) > 0 {
		result.info += fmt.Sprintf("联盟等级:%v（排名%v）\n", legion[0].LegionLevel, legion[0].Rank)
		result.info += fmt.Sprintf("联盟战斗力:%v（每日%.0f币）\n", legion[0].RaidPower, float64(legion[0].RaidPower)/1157407.41)

		var legionCoinCap = 200
		switch {
		case legion[0].LegionLevel >= 8000:
			legionCoinCap = 700
		case legion[0].LegionLevel >= 5500:
			legionCoinCap = 500
		case legion[0].LegionLevel >= 3000:
			legionCoinCap = 300
		}
		result.info += fmt.Sprintf("联盟币上限:%v\n", legionCoinCap)
	}

	// 保存数据进缓存
	checkNumberOfTimes[gropFromQQ]++
	//charactersLevelRank[strconv.Itoa(rankInfo[0].)] = name
	if rankInfo[0].Level > 219 && rankInfo[0].WorldID == 45 { // 大于219且在R区才有排名记录缓存
		jobWithRank := fmt.Sprintf("%v%v", rRankInfo[0].JobName, rRankInfo[0].Rank)
		classLevelRank[jobWithRank] = rankInfo[0].CharacterName
	}
	charactersCatch[name] = result
	return result
}

func CheckMapleGG(name string, gropFromQQ int) CharInfoResult {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return CharInfoResult{"", "查询失败"}
	}
	resetCacheEveryday()
	if checkNumberOfTimes[gropFromQQ] >= 12 {
		return CharInfoResult{"", "今日查询已达上限"}
	}

	// map缓存查找
	if chara, ok := charactersCatch[name]; ok {
		checkNumberOfTimes[gropFromQQ]++
		return chara
	}

	body, err := webGetRequest("https://api.maplestory.gg/v1/public/character/gms/" + name)
	if err != nil {
		return CharInfoResult{"", "查询失败"}
	}

	var mapleInfo MapleInfo
	err = json.Unmarshal(body, &mapleInfo)
	if err != nil || len(mapleInfo.Name) == 0 {
		return CharInfoResult{"", "查询失败2"}
	}

	if mapleInfo.Message == "Unable to find character" {
		// 查询的角色不存在
		checkNumberOfTimes[gropFromQQ]++
		charactersCatch[name] = CharInfoResult{"", "查询的角色不存在"}
		return charactersCatch[name]
	}

	result := CharInfoResult{}
	result.imageURL = mapleInfo.CharacterImageURL
	outName := strconv.QuoteToASCII(mapleInfo.Name)
	result.info += "角色:" + outName[1:len(outName)-1] + "  （" + serverName2Chinese[mapleInfo.Server] + "）\n"
	result.info += fmt.Sprintf("等级:%v - %v%% （排名%v）\n", mapleInfo.Level, mapleInfo.EXPPercent, mapleInfo.ServerRank)
	if mapleInfo.ServerClassRanking > 0 {
		result.info += fmt.Sprintf("职业:%v（排名%v）\n\n", mapleInfo.Class, mapleInfo.ServerClassRanking)
	} else {
		result.info += fmt.Sprintf("职业:%v\n\n", mapleInfo.Class)
	}

	if mapleInfo.AchievementPoints != 0 && mapleInfo.AchievementRank != 0 {
		result.info += fmt.Sprintf("成就值:%v（排名%v）\n", mapleInfo.AchievementPoints, mapleInfo.AchievementRank)
	}
	if mapleInfo.LegionLevel != 0 && mapleInfo.LegionRank != 0 && mapleInfo.LegionPower != 0 && mapleInfo.LegionCoinsPerDay != 0 {
		result.info += fmt.Sprintf("联盟等级:%v（排名%v）\n", mapleInfo.LegionLevel, mapleInfo.LegionRank)
		if mapleInfo.LegionPower > 1000000 {
			result.info += fmt.Sprintf("联盟战斗力:%vm（每日%v币）\n", mapleInfo.LegionPower/1000000, mapleInfo.LegionCoinsPerDay)
		} else {
			result.info += fmt.Sprintf("战五渣（每日%v币）\n", mapleInfo.LegionCoinsPerDay)
		}
		var legionCoinCap = 200
		switch {
		case mapleInfo.LegionLevel >= 8000:
			legionCoinCap = 700
		case mapleInfo.LegionLevel >= 5500:
			legionCoinCap = 500
		case mapleInfo.LegionLevel >= 3000:
			legionCoinCap = 300
		}
		result.info += fmt.Sprintf("联盟币上限:%v\n", legionCoinCap)
	}

	result.info += "---------------------------------\n"
	result.info += calcExpDay(int64(mapleInfo.Level), mapleInfo)
	fmt.Println(result.info)

	// 保存数据进缓存
	checkNumberOfTimes[gropFromQQ]++
	charactersLevelRank[strconv.Itoa(mapleInfo.ServerRank)] = name
	if (mapleInfo.Class != "Thief" && mapleInfo.Class != "Dual Blade") || mapleInfo.ServerClassRanking != 0 {
		classAndRank := mapleInfo.Class + strconv.Itoa(mapleInfo.ServerClassRanking)
		classLevelRank[classAndRank] = name
	}
	if mapleInfo.LegionRank != 0 {
		classLevelRank["联盟"+strconv.Itoa(mapleInfo.LegionRank)] = name
	}
	charactersCatch[name] = result
	return result
}

func calcExpDay(level int64, mapleInfo MapleInfo) string {
	if level == 300 {
		return "大佬啊~ 带我打Boss！ 我混车超稳的！"
	}

	var totalExp int64
	var nextLevel int64
	for _, nextLevelExp := range TotalExpCollection {
		if level < nextLevelExp.Level {
			totalExp = nextLevelExp.Exp
			nextLevel = nextLevelExp.Level
			break
		}
	}

	var day float64
	var dayOne float64
	dateLen := len(mapleInfo.GraphData)
	if dateLen > 1 {
		avgDayExp := float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP - mapleInfo.GraphData[0].TotalOverallEXP)
		day = float64(dateLen-1) * float64(totalExp-mapleInfo.GraphData[dateLen-1].TotalOverallEXP) / avgDayExp
		lastDayExp := float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP - mapleInfo.GraphData[dateLen-2].TotalOverallEXP)
		if lastDayExp != 0 {
			dayOne = float64(totalExp-mapleInfo.GraphData[dateLen-1].TotalOverallEXP) / lastDayExp
		}
	}

	switch {
	case 0 < day && day < 10000 && 0 < dayOne && dayOne < 10000:
		return fmt.Sprintf("参照最近1天还有%.1f天就到%d级！\n", dayOne, nextLevel) +
			fmt.Sprintf("参照最近%d天还有%.1f天就到%d级！", dateLen-1, day, nextLevel)
	case 0 < dayOne && dayOne < 10000:
		return fmt.Sprintf("参照最近1天还有%.1f天就到%d级！", dayOne, nextLevel)
	case 0 < day && day < 10000:
		return fmt.Sprintf("参照最近%d天还有%.1f天就到%d级！", dateLen-1, day, nextLevel)
	default:
		return fmt.Sprintf("这辈子没希望肝到%d了~", nextLevel)
	}
}

func QueryClassRanking(rankNumber string, groupFromQQ int, url string) (result CharInfoResult) {
	resetCacheEveryday()

	if checkNumberOfTimes[groupFromQQ] >= 12 {
		return CharInfoResult{"", "今日查询已达上限"}
	}

	if characterName, ok := classLevelRank[rankNumber]; ok {
		fmt.Println("职业缓存")
		return charactersCatch[characterName]
	}

	ieproxy.OverrideEnvWithStaticProxy()
	http.DefaultTransport.(*http.Transport).Proxy = ieproxy.GetProxyFunc()
	client := http.Client{
		Timeout: 4 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CharInfoResult{"", "查询失败:client.Do(req)获取失败"}
	}

	req.Header.Set("authority", "maplestory.nexon.net")
	req.Header.Set("scheme", "https")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36 Edg/89.0.774.54")
	req.Header.Set("dnt", "1")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en,zh-CN;q=0.9,zh;q=0.8,en-GB;q=0.7,en-US;q=0.6,zh-TW;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		return CharInfoResult{"", "查询失败:client.Do(req)获取失败"}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CharInfoResult{"", "查询失败:ReadAll获取失败"}
	}

	var officialRank OfficialWorldRank
	err = json.Unmarshal(body, &officialRank)
	if err != nil || len(officialRank) == 0 {
		return CharInfoResult{"", "查询失败:json解析失败"}
	}

	return CheckMapleGG(officialRank[0].CharacterName, groupFromQQ)
}

func QueryRanking(rankNumber string, gropFromQQ int, url string) (result CharInfoResult) {
	resetCacheEveryday()

	if checkNumberOfTimes[gropFromQQ] >= 12 {
		return CharInfoResult{"", "今日查询已达上限"}
	}

	if characterName, ok := charactersLevelRank[rankNumber]; ok {
		return charactersCatch[characterName]
	}

	request, err := http.NewRequest(http.MethodGet, url+rankNumber, nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36 Edg/88.0.705.81")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code:", resp.StatusCode)
		return CharInfoResult{"", "查询失败"}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Print(err)
	}

	ss := doc.Find(".ranking-container table tbody tr td").Eq(2)
	if !ss.Is("td") {
		return CharInfoResult{"", "查询失败"}
	}

	characterName := strings.TrimSpace(ss.Text())
	fmt.Println("characterName:" + characterName)

	return CheckMapleGG(characterName, gropFromQQ)
}

func resetCacheEveryday() {
	if now := time.Now().Add(time.Duration(-4) * time.Hour).Day(); now != catchTime {
		catchTime = now
		charactersLevelRank = nil
		charactersLevelRank = make(map[string]string)
		classLevelRank = nil
		classLevelRank = make(map[string]string)
		charactersCatch = nil
		charactersCatch = make(map[string]CharInfoResult)
		checkNumberOfTimes = nil
		checkNumberOfTimes = make(map[int]int)
	}
}

func CheckClassRank(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	generalNameCreate := func(names []string) (result []string) {
		for index, name := range names {
			if index == 0 {
				result = append(result, strings.ToLower(name)+"第")
			}
			result = append(result, name+"第")
		}
		return
	}

	// 查询职业排名
	for _, class := range AllClass {
		groupMessage, ok := IsPrefix(groupMessage, generalNameCreate(class.Names)...)
		if ok {
			if len(groupMessage) == 0 {
				return true
			}
			level, err := strconv.Atoi(groupMessage)
			if err != nil {
				var num zh.Uint64
				_, err := fmt.Sscan(groupMessage, &num)
				if err != nil {
					return true
				}
				level = int(num)
				groupMessage = strconv.Itoa(level)
			}

			if level > 10000 {
				return true
			}
			character := QueryClassRanking(class.Names[0]+groupMessage, fromQQ, classURL[class.Names[0]]+groupMessage)
			if character.imageURL == "" {
				SendGroupMsg(loginQQ, fromGroup, character.info)
			} else {
				imageURL := GetGroupImage(loginQQ, fromGroup, 2, character.imageURL)
				SendGroupMsg(loginQQ, fromGroup, imageURL+"\n"+character.info)
			}
			return true
		}
	}

	// 查询我
	if groupMessage == "me" || groupMessage == "wo" || groupMessage == "我" {
		user := QQBindCharacter{QQ: fromQQ}
		result := gdb.First(&user, "QQ = ?", fromQQ)
		if result.RowsAffected > 0 {
			groupMessage = user.Character
		} else {
			SendGroupMsg(loginQQ, fromGroup, "请先绑定角色\n例如：查询绑定PaperWang")
		}
	}

	// 查询{角色名}
	for _, v := range groupMessage {
		if unicode.Is(unicode.Han, v) {
			return false
		}
	}
	character := CheckMapleGG(groupMessage, fromQQ)
	if character.imageURL == "" {
		SendGroupMsg(loginQQ, fromGroup, character.info)
	} else {
		imageURL := GetGroupImage(loginQQ, fromGroup, 2, character.imageURL)
		SendGroupMsg(loginQQ, fromGroup, imageURL+"\n"+character.info)
	}
	return true
}

func bindCharacter(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	for _, v := range groupMessage {
		if unicode.Is(unicode.Han, v) {
			return false
		}
	}
	if len(groupMessage) == 0 {
		return false
	}
	groupMessage = strings.ToLower(groupMessage) // 角色名无视大小写，全部小写录入可锁定角色,防止大小写不同重复录入
	_ = gdb.AutoMigrate(&QQBindCharacter{})
	user := QQBindCharacter{QQ: fromQQ}

	// 角色名拆失败
	if gdb.First(&user, "(character = ? OR qq = ?) AND lock = true ", groupMessage).RowsAffected == 0 {
		// 更新数据库
		if gdb.Model(user).Update("Character", groupMessage).RowsAffected == 0 {
			user.Character = groupMessage
			gdb.Create(&user)
		}
		SendGroupMsg(loginQQ, fromGroup, "绑定成功")
	} else {
		SendGroupMsg(loginQQ, fromGroup, "绑定失败")
	}
	return true
}
