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

type OfficialRank []struct {
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

type CharInfoResult struct {
	imageUrl, info string
}

var classUrl map[string]string
var checkNumberOfTimes map[int]int
var charactersCatch map[string]CharInfoResult
var charactersLevelRank map[string]string
var catchTime int
var classLevelRank map[string]string

func CheckMapleGG(name string, gropFromQQ int) (result CharInfoResult) {
	name = strings.TrimSpace(name)
	resetCacheEveryday()

	if checkNumberOfTimes[gropFromQQ] >= 12 {
		return CharInfoResult{"", "????????????????????????"}
	}

	//map????????????
	if chara, ok := charactersCatch[name]; ok {
		checkNumberOfTimes[gropFromQQ] += 1
		return chara
	}

	body, err := webGetRequest("https://api.maplestory.gg/v1/public/character/gms/" + name)
	if err != nil {
		return CharInfoResult{"", "????????????"}
	}

	var mapleInfo MapleInfo
	err = json.Unmarshal(body, &mapleInfo)
	if err != nil {
		return CharInfoResult{"", "????????????2"}
	}

	if mapleInfo.Message == "Unable to find character" {
		//????????????????????????
		checkNumberOfTimes[gropFromQQ] += 1
		charactersCatch[name] = CharInfoResult{"", "????????????????????????"}
		return charactersCatch[name]
	}

	serverAbbr := make(map[string]string)
	serverAbbr["Reboot (NA)"] = "R???"
	serverAbbr["Aurora"] = "A???"
	serverAbbr["Bera"] = "B???"
	serverAbbr["Scania"] = "S???"

	result.imageUrl = mapleInfo.CharacterImageURL
	outName := strconv.QuoteToASCII(mapleInfo.Name)
	result.info += "??????:" + outName[1:len(outName)-1] + "  ???" + serverAbbr[mapleInfo.Server] + "???\n"
	result.info += fmt.Sprintf("??????:%v - %v%% ?????????%v???\n", mapleInfo.Level, mapleInfo.EXPPercent, mapleInfo.ServerRank)
	if mapleInfo.ServerClassRanking > 0 {
		result.info += fmt.Sprintf("??????:%v?????????%v???\n\n", mapleInfo.Class, mapleInfo.ServerClassRanking)
	} else {
		result.info += fmt.Sprintf("??????:%v\n\n", mapleInfo.Class)
	}

	if mapleInfo.AchievementPoints != 0 && mapleInfo.AchievementRank != 0 {
		result.info += fmt.Sprintf("?????????:%v?????????%v???\n", mapleInfo.AchievementPoints, mapleInfo.AchievementRank)
	}
	if mapleInfo.LegionLevel != 0 && mapleInfo.LegionRank != 0 && mapleInfo.LegionPower != 0 && mapleInfo.LegionCoinsPerDay != 0 {
		result.info += fmt.Sprintf("????????????:%v?????????%v???\n", mapleInfo.LegionLevel, mapleInfo.LegionRank)
		result.info += fmt.Sprintf("???????????????:%v?????????%v??????\n", mapleInfo.LegionPower, mapleInfo.LegionCoinsPerDay)

		var legionCoinCap = 200
		switch {
		case mapleInfo.LegionLevel >= 8000:
			legionCoinCap = 700
		case mapleInfo.LegionLevel >= 5500:
			legionCoinCap = 500
		case mapleInfo.LegionLevel >= 3000:
			legionCoinCap = 300
		}
		result.info += fmt.Sprintf("???????????????:%v\n", legionCoinCap)
	}

	result.info += "---------------------------------\n"
	dateLen := len(mapleInfo.GraphData)
	if mapleInfo.Level < 250 {
		var day float64
		var dayOne float64
		if dateLen > 1 {
			day = float64(dateLen-1) * float64(9652268252132-mapleInfo.GraphData[dateLen-1].TotalOverallEXP) / float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP-mapleInfo.GraphData[0].TotalOverallEXP)
			if float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP-mapleInfo.GraphData[dateLen-2].TotalOverallEXP) != 0 {
				dayOne = float64(9652268252132-mapleInfo.GraphData[dateLen-1].TotalOverallEXP) / float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP-mapleInfo.GraphData[dateLen-2].TotalOverallEXP)
			}
		}
		if 0 < day && day < 2500 && 0 < dayOne && dayOne < 2500 {
			result.info += fmt.Sprintf("????????????1??????????????? ??????%.1f?????????250??????\n", dayOne)
			result.info += fmt.Sprintf("????????????%d??????????????? ??????%.1f?????????250??????", dateLen-1, day)
		} else {
			result.info += fmt.Sprintf("????????????????????????250???~")
		}
	} else if mapleInfo.Level < 275 {
		var day float64
		var dayOne float64
		if dateLen > 1 {
			day = float64(dateLen-1) * float64(86471480104001-mapleInfo.GraphData[dateLen-1].TotalOverallEXP) / float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP-mapleInfo.GraphData[0].TotalOverallEXP)
			if float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP-mapleInfo.GraphData[dateLen-2].TotalOverallEXP) != 0 {
				dayOne = float64(86471480104001-mapleInfo.GraphData[dateLen-1].TotalOverallEXP) / float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP-mapleInfo.GraphData[dateLen-2].TotalOverallEXP)
			}
		}
		if 0 < day && day < 2500 && 0 < dayOne && dayOne < 2500 {
			result.info += fmt.Sprintf("????????????1??????????????? ??????%.1f?????????275??????\n", dayOne)
			result.info += fmt.Sprintf("????????????%d??????????????? ??????%.1f?????????275??????", dateLen-1, day)
		} else if 0 < dayOne && dayOne < 2500 {
			result.info += fmt.Sprintf("????????????1??????????????? ??????%.1f?????????275??????", dayOne)
		} else if 0 < day && day < 2500 {
			result.info += fmt.Sprintf("????????????%d??????????????? ??????%.1f?????????275??????", dateLen-1, day)
		} else {
			result.info += fmt.Sprintf("????????????????????????275???~")
		}
	} else {
		result.info += fmt.Sprintf("?????????~ ?????????Boss??? ?????????????????????")
	}

	fmt.Println(result.info)

	// ?????????????????????
	checkNumberOfTimes[gropFromQQ] += 1
	charactersLevelRank[strconv.Itoa(mapleInfo.ServerRank)] = name
	if (mapleInfo.Class != "Thief" && mapleInfo.Class != "Dual Blade") || mapleInfo.ServerClassRanking != 0 {
		classAndRank := mapleInfo.Class + strconv.Itoa(mapleInfo.ServerClassRanking)
		classLevelRank[classAndRank] = name
	}
	if mapleInfo.LegionRank != 0 {
		classLevelRank["??????"+strconv.Itoa(mapleInfo.LegionRank)] = name
	}
	charactersCatch[name] = result
	return
}

func QueryWorldRanking(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
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

	character := QueryRanking(groupMessage, fromQQ, "https://maplestory.nexon.net/rankings/world-ranking/reboot-(na)?rebootIndex=0&pageIndex=")
	if character.imageUrl == "" {
		SendGroupMsg(loginQQ, fromGroup, character.info)
	} else {
		imageUrl := GetGroupImage(loginQQ, fromGroup, 2, character.imageUrl)
		SendGroupMsg(loginQQ, fromGroup, imageUrl+"\n"+character.info)
	}
	return true
}

func QueryClassRanking(rankNumber string, gropFromQQ int, url string) (result CharInfoResult) {
	resetCacheEveryday()

	if checkNumberOfTimes[gropFromQQ] >= 12 {
		return CharInfoResult{"", "????????????????????????"}
	}

	if characterName, ok := classLevelRank[rankNumber]; ok {
		fmt.Println("????????????")
		return charactersCatch[characterName]
	}

	ieproxy.OverrideEnvWithStaticProxy()
	http.DefaultTransport.(*http.Transport).Proxy = ieproxy.GetProxyFunc()
	client := http.Client{
		Timeout: 4 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CharInfoResult{"", "????????????:client.Do(req)????????????"}
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
		return CharInfoResult{"", "????????????:client.Do(req)????????????"}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CharInfoResult{"", "????????????:ReadAll????????????"}
	}

	var officialRank OfficialRank
	err = json.Unmarshal(body, &officialRank)
	if err != nil || len(officialRank) <= 0 {
		return CharInfoResult{"", "????????????:json????????????"}
	}

	//fmt.Println(officialRank)
	return CheckMapleGG(officialRank[0].CharacterName, gropFromQQ)
}

func QueryRanking(rankNumber string, gropFromQQ int, url string) (result CharInfoResult) {
	resetCacheEveryday()

	if checkNumberOfTimes[gropFromQQ] >= 12 {
		return CharInfoResult{"", "????????????????????????"}
	}

	if characterName, ok := charactersLevelRank[rankNumber]; ok {
		return charactersCatch[characterName]
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
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
		return CharInfoResult{"", "????????????"}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	ss := doc.Find(".ranking-container table tbody tr td").Eq(2)
	if !ss.Is("td") {
		return CharInfoResult{"", "????????????"}
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
				result = append(result, strings.ToLower(name)+"??????")
				result = append(result, strings.ToLower(name)+"???")
			}
			result = append(result, name+"??????")
			result = append(result, name+"???")
		}
		return
	}

	AllClass := []struct {
		Names []string
	}{
		{[]string{"??????"}},
		{[]string{"Warrior", "??????", "??????", "??????", "??????"}},
		{[]string{"Magician", "??????", "??????", "??????", "??????"}},
		{[]string{"Bowman", "??????", "?????????", "?????????", "??????", "??????"}},
		{[]string{"Thief", "??????", "??????", "??????", "??????"}},
		{[]string{"Pirate", "??????", "??????", "??????", "??????"}},
		{[]string{"Aran", "??????"}},
		{[]string{"Evan", "??????"}},
		{[]string{"Mercedes", "??????"}},
		{[]string{"Phantom", "??????"}},
		{[]string{"Jett", "??????"}},
		{[]string{"Mihile", "?????????", "?????????"}},
		{[]string{"Luminous", "??????", "????????????"}},
		{[]string{"Kaiser", "??????", "??????"}},
		{[]string{"Angelic Buster", "??????", "ab", "AB"}},
		{[]string{"Hayato", "??????"}},
		{[]string{"Kanna", "?????????", "??????"}},
		{[]string{"Xenon", "??????", "??????"}},
		{[]string{"Zero", "?????????", "??????"}},
		{[]string{"Beast Tamer", "BT", "bt", "?????????", "?????????", "lzl"}},
		{[]string{"Shade", "??????"}},
		{[]string{"Kinesis", "??????"}},
		{[]string{"Blaster", "??????", "?????????"}},
		{[]string{"Cadena", "??????", "?????????"}},
		{[]string{"Illium", "??????", "????????????", "????????????"}},
		{[]string{"Ark", "??????", "??????"}},
		{[]string{"Pathfinder", "pf", "?????????", "????????????", "PF"}},
		{[]string{"Hoyoung", "??????"}},
		{[]string{"Adele", "??????", "?????????"}},
		{[]string{"Dawn Warrior", "dw", "DW", "??????"}},
		{[]string{"Blaze Wizard", "BW", "bw", "?????????"}},
		{[]string{"Wind Archer", "WA", "wa", "??????"}},
		{[]string{"Night Walker", "NW", "nw", "??????", "?????????"}},
		{[]string{"Thunder Breaker", "TB", "tb", "??????", "?????????"}},
		{[]string{"Demon Slayer", "DS", "ds", "??????"}},
		{[]string{"Battle Mage", "BM", "bm", "??????"}},
		{[]string{"Wild Hunter", "WH", "wh", "??????"}},
		{[]string{"Mechanicr", "??????", "??????"}},
		{[]string{"Demon Avenger", "??????", "DA"}},
	}

	//????????????
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
			character := QueryClassRanking(class.Names[0]+groupMessage, fromQQ, classUrl[class.Names[0]]+groupMessage)
			if character.imageUrl == "" {
				SendGroupMsg(loginQQ, fromGroup, character.info)
			} else {
				imageUrl := GetGroupImage(loginQQ, fromGroup, 2, character.imageUrl)
				SendGroupMsg(loginQQ, fromGroup, imageUrl+"\n"+character.info)
			}
			return true
		}
	}

	//?????????
	if groupMessage == "me" || groupMessage == "wo" || groupMessage == "???" {
		user := QQBindCharacter{QQ: fromQQ}
		//gdb.Delete(user)
		//gdb.Create(user)
		result := gdb.First(&user, "QQ = ?", fromQQ)

		if result.RowsAffected > 0 {
			groupMessage = user.Character
		} else {
			SendGroupMsg(loginQQ, fromGroup, "??????????????????\n?????????????????????PaperWang")
		}
	}

	//??????{?????????}
	for _, v := range groupMessage {
		if unicode.Is(unicode.Han, v) {
			return false
		}
	}
	character := CheckMapleGG(groupMessage, fromQQ)
	if character.imageUrl == "" {
		SendGroupMsg(loginQQ, fromGroup, character.info)
	} else {
		imageUrl := GetGroupImage(loginQQ, fromGroup, 2, character.imageUrl)
		SendGroupMsg(loginQQ, fromGroup, imageUrl+"\n"+character.info)
	}

	return true
}

func bindCharacter(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	for _, v := range groupMessage {
		if unicode.Is(unicode.Han, v) {
			return false
		}
	}

	//gdb.AutoMigrate(&QQBindCharacter{})

	user := QQBindCharacter{QQ: fromQQ}
	groupMessage = strings.ToLower(groupMessage) //????????????????????????????????????????????????????????????,?????????????????????????????????

	//??????????????????
	if gdb.First(&user, "(character = ? OR qq = ?) AND lock = true ", groupMessage).RowsAffected == 0 {
		//???????????????
		if gdb.Model(user).Update("Character", groupMessage).RowsAffected == 0 {
			gdb.Create(&user)
		}
		SendGroupMsg(loginQQ, fromGroup, "????????????")
	} else {
		SendGroupMsg(loginQQ, fromGroup, "????????????")
		return true
	}

	return true
}
