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
	QQ        int
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

type charInfoResult struct {
	imageUrl, info string
}

var classUrl map[string]string
var checkNumberOfTimes map[int]int
var charactersCatch map[string]charInfoResult
var charactersLevelRank map[string]string
var catchTime int
var classLevelRank map[string]string

func CheckMapleGG(name string, gropFromQQ int) (result charInfoResult) {
	name = strings.TrimSpace(name)
	resetCacheEveryday()

	if checkNumberOfTimes[gropFromQQ] >= 12 {
		return charInfoResult{"", "今日查询已达上限"}
	}

	if chara, ok := charactersCatch[name]; ok {
		checkNumberOfTimes[gropFromQQ] += 1
		fmt.Println("CheckMapleGG cache")
		return chara
	}

	//var PTransport = & http.Transport {Proxy: http.ProxyFromEnvironment}

	//proxyUrl, err := url.Parse("http://proxyIp:proxyPort")
	//ieproxy.OverrideEnvWithStaticProxy()
	//http.DefaultTransport.(*http.Transport).Proxy = http.ProxyFromEnvironment
	//client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

	//req, err := http.NewRequest("GET", "https://api.maplestory.gg/v1/public/character/gms/" + name, nil)
	//resp, err := client.Do(req)
	ieproxy.OverrideEnvWithStaticProxy()
	http.DefaultTransport.(*http.Transport).Proxy = ieproxy.GetProxyFunc()
	client := http.Client{
		Timeout: 4 * time.Second,
	}

	resp, err := client.Get("https://api.maplestory.gg/v1/public/character/gms/" + name)

	if err != nil {
		// handle error
		resp2, err2 := client.Get("https://api.maplestory.gg/v1/public/character/gms/" + name)
		if err2 != nil {
			return charInfoResult{"", "查询失败"}
		}
		defer resp2.Body.Close()
		resp, resp2 = resp2, resp
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//s := string(body)
	var mapleInfo MapleInfo
	err = json.Unmarshal(body, &mapleInfo)
	if err != nil {
		return charInfoResult{"", "查询失败2"}
	}

	if mapleInfo.Message == "Unable to find character" {
		//查询的角色不存在
		checkNumberOfTimes[gropFromQQ] += 1
		charactersCatch[name] = charInfoResult{"", "查询的角色不存在"}
		return charactersCatch[name]
	}

	serverAbbr := make(map[string]string)
	serverAbbr["Reboot (NA)"] = "R区"
	serverAbbr["Aurora"] = "A区"
	serverAbbr["Bera"] = "B区"
	serverAbbr["Scania"] = "S区"

	result.imageUrl = mapleInfo.CharacterImageURL
	outName := strconv.QuoteToASCII(mapleInfo.Name)
	result.info += "角色:" + outName[1:len(outName)-1] + "  （" + serverAbbr[mapleInfo.Server] + "）\n"
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
		result.info += fmt.Sprintf("联盟战斗力:%v（每日%v币）\n", mapleInfo.LegionPower, mapleInfo.LegionCoinsPerDay)

		var legionCoinCap int = 200
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
	dateLen := len(mapleInfo.GraphData)
	if mapleInfo.Level < 250 {
		var day float64
		var dayOne float64
		if dateLen > 1 {
			day = float64(dateLen-1) * float64(9654369842607-mapleInfo.GraphData[dateLen-1].TotalOverallEXP) / float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP-mapleInfo.GraphData[0].TotalOverallEXP)
			if float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP-mapleInfo.GraphData[dateLen-2].TotalOverallEXP) != 0 {
				dayOne = float64(9654369842607-mapleInfo.GraphData[dateLen-1].TotalOverallEXP) / float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP-mapleInfo.GraphData[dateLen-2].TotalOverallEXP)
			}
		}
		if 0 < day && day < 2500 && 0 < dayOne && dayOne < 2500 {
			result.info += fmt.Sprintf("按照最近1天的肝度！ 还有%.1f天就到250级！\n", dayOne)
			result.info += fmt.Sprintf("按照最近%d天的肝度！ 还有%.1f天就到250级！", dateLen-1, day)
		} else {
			result.info += fmt.Sprintf("这辈子没希望肝到250了~")
		}
	} else if mapleInfo.Level < 275 {
		var day float64
		var dayOne float64
		if dateLen > 1 {
			day = float64(dateLen-1) * float64(86473581694476-mapleInfo.GraphData[dateLen-1].TotalOverallEXP) / float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP-mapleInfo.GraphData[0].TotalOverallEXP)
			if float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP-mapleInfo.GraphData[dateLen-2].TotalOverallEXP) != 0 {
				dayOne = float64(86473581694476-mapleInfo.GraphData[dateLen-1].TotalOverallEXP) / float64(mapleInfo.GraphData[dateLen-1].TotalOverallEXP-mapleInfo.GraphData[dateLen-2].TotalOverallEXP)
			}
		}
		if 0 < day && day < 2500 && 0 < dayOne && dayOne < 2500 {
			result.info += fmt.Sprintf("按照最近1天的肝度！ 还有%.1f天就到275级！\n", dayOne)
			result.info += fmt.Sprintf("按照最近%d天的肝度！ 还有%.1f天就到275级！", dateLen-1, day)
		} else if 0 < dayOne && dayOne < 2500 {
			result.info += fmt.Sprintf("按照最近1天的肝度！ 还有%.1f天就到275级！", dayOne)
		} else if 0 < day && day < 2500 {
			result.info += fmt.Sprintf("按照最近%d天的肝度！ 还有%.1f天就到275级！", dateLen-1, day)
		} else {
			result.info += fmt.Sprintf("这辈子没希望肝到275了~")
		}
	} else {
		result.info += fmt.Sprintf("大佬啊~ 带我打Boss！ 我混车超稳的！")
	}

	fmt.Println(result.info)

	// 保存数据进缓存
	checkNumberOfTimes[gropFromQQ] += 1
	charactersLevelRank[strconv.Itoa(mapleInfo.ServerRank)] = name
	if (mapleInfo.Class != "Thief" || mapleInfo.Class != "Dual Blade") && mapleInfo.ServerClassRanking != 0 {
		classAndRank := mapleInfo.Class + strconv.Itoa(mapleInfo.ServerClassRanking)
		classLevelRank[classAndRank] = name
	}
	if mapleInfo.LegionRank != 0 {
		classLevelRank["联盟"+strconv.Itoa(mapleInfo.LegionRank)] = name
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

func QueryClassRanking(rankNumber string, gropFromQQ int, url string) (result charInfoResult) {
	resetCacheEveryday()

	if checkNumberOfTimes[gropFromQQ] >= 12 {
		return charInfoResult{"", "今日查询已达上限"}
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
	//

	req.Header.Set("authority", "maplestory.nexon.net")
	req.Header.Set("scheme", "https")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36 Edg/89.0.774.54")
	req.Header.Set("dnt", "1")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	//req.Header.Set("referer", "https://maplestory.nexon.net/rankings/job-ranking/explorer/warrior?pageIndex=&rebootIndex1=6&page_index=1")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en,zh-CN;q=0.9,zh;q=0.8,en-GB;q=0.7,en-US;q=0.6,zh-TW;q=0.5")

	//resp, err := client.Get("https://maplestory.nexon.net/api/ranking?id=job&id2=1&page_index=1")
	resp, err := client.Do(req)
	if err != nil {
		return charInfoResult{"", "查询失败:client.Do(req)获取失败"}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return charInfoResult{"", "查询失败:ReadAll获取失败"}
	}

	var officialRank OfficialRank
	err = json.Unmarshal(body, &officialRank)
	if err != nil || len(officialRank) <= 0 {
		return charInfoResult{"", "查询失败:json解析失败"}
	}

	//fmt.Println(officialRank)
	return CheckMapleGG(officialRank[0].CharacterName, gropFromQQ)
}

func QueryRanking(rankNumber string, gropFromQQ int, url string) (result charInfoResult) {
	resetCacheEveryday()

	if checkNumberOfTimes[gropFromQQ] >= 12 {
		return charInfoResult{"", "今日查询已达上限"}
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
		return charInfoResult{"", "查询失败"}
	}
	/*
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", all)

	*/

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	ss := doc.Find(".ranking-container table tbody tr td").Eq(2)
	if !ss.Is("td") {
		return charInfoResult{"", "查询失败"}
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
		charactersCatch = make(map[string]charInfoResult)
		checkNumberOfTimes = nil
		checkNumberOfTimes = make(map[int]int)
	}
}

func CheckClassRank(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {

	generalNameCreate := func(names []string) (result []string) {
		for index, name := range names {
			if index == 0 {
				result = append(result, strings.ToLower(name)+"排名")
				result = append(result, strings.ToLower(name)+"第")
			}
			result = append(result, name+"排名")
			result = append(result, name+"第")
		}
		return
	}

	AllClass := []struct {
		Names []string
	}{
		{[]string{"联盟"}},
		{[]string{"Warrior", "战士", "英雄", "黑骑", "圣骑"}},
		{[]string{"Magician", "法师", "冰雷", "火毒", "主教"}},
		{[]string{"Bowman", "弓手", "弓箭手", "神射手", "神射", "弩手"}},
		{[]string{"Thief", "飞侠", "双刀", "刀飞", "标飞"}},
		{[]string{"Pirate", "海盗", "船长", "队长", "火炮"}},
		{[]string{"Aran", "战神"}},
		{[]string{"Evan", "龙神"}},
		{[]string{"Mercedes", "双弩"}},
		{[]string{"Phantom", "幻影"}},
		{[]string{"Jett", "杰特"}},
		{[]string{"Mihile", "米哈尔", "米哈哈"}},
		{[]string{"Luminous", "夜光", "夜光法师"}},
		{[]string{"Kaiser", "凯撒", "狂龙"}},
		{[]string{"Angelic Buster", "天使", "ab", "AB"}},
		{[]string{"Hayato", "剑豪"}},
		{[]string{"Kanna", "阴阳师", "娜神"}},
		{[]string{"Xenon", "煎饼", "尖兵"}},
		{[]string{"Zero", "神之子", "神子"}},
		{[]string{"Beast Tamer", "BT", "bt", "林志玲", "林之灵", "lzl"}},
		{[]string{"Shade", "隐月"}},
		{[]string{"Kinesis", "超能"}},
		{[]string{"Blaster", "爆破", "爆破者"}},
		{[]string{"Cadena", "卡姐", "卡德娜"}},
		{[]string{"Illium", "黑皮", "圣经使徒", "圣晶使徒"}},
		{[]string{"Ark", "牙科", "亚克"}},
		{[]string{"Pathfinder", "pf", "开拓者", "古迹猎人", "PF"}},
		{[]string{"Hoyoung", "虎影"}},
		{[]string{"Adele", "阿呆", "阿黛尔"}},
		{[]string{"Dawn Warrior", "dw", "DW", "魂骑"}},
		{[]string{"Blaze Wizard", "BW", "bw", "炎术士"}},
		{[]string{"Wind Archer", "WA", "wa", "风铃"}},
		{[]string{"Night Walker", "NW", "nw", "夜行", "夜行者"}},
		{[]string{"Thunder Breaker", "TB", "tb", "奇袭", "奇袭者"}},
		{[]string{"Demon Slayer", "DS", "ds", "红毛"}},
		{[]string{"Battle Mage", "BM", "bm", "战法"}},
		{[]string{"Wild Hunter", "WH", "wh", "豹弩"}},
		{[]string{"Mechanicr", "轮椅", "机械"}},
		{[]string{"Demon Avenger", "白毛", "DA"}},
	}

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
	//查询我
	if groupMessage == "me" || groupMessage == "wo" || groupMessage == "我" {
		user := QQBindCharacter{QQ: fromQQ}
		//gdb.Delete(user)
		//gdb.Create(user)
		result := gdb.Where("QQ = ?", fromQQ).First(&user)

		if result.RowsAffected > 0 {
			groupMessage = user.Character
		} else {
			SendGroupMsg(loginQQ, fromGroup, "请先绑定角色\n例如：查询绑定PaperWang")
		}
	}

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

	user := QQBindCharacter{QQ: fromQQ}

	result := gdb.Where("QQ = ?", fromQQ).First(&user)
	user.Character = groupMessage
	if result.Error != gorm.ErrRecordNotFound && result.Error != nil {
		SendGroupMsg(loginQQ, fromGroup, result.Error.Error())
	}
	if result.RowsAffected > 0 {
		gdb.Model(&user).Where("QQ = ?", fromQQ).Update("character", groupMessage)
	} else {
		gdb.Create(&user)
	}
	SendGroupMsg(loginQQ, fromGroup, "绑定成功")
	return true
}

func getOfficialRank() {
	ieproxy.OverrideEnvWithStaticProxy()
	http.DefaultTransport.(*http.Transport).Proxy = ieproxy.GetProxyFunc()
	client := http.Client{
		Timeout: 4 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://maplestory.nexon.net/api/ranking?id=job&id2=1&page_index=1", nil)

	req.Header.Set("authority", "maplestory.nexon.net")
	req.Header.Set("scheme", "https")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36 Edg/89.0.774.54")
	req.Header.Set("dnt", "1")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", "https://maplestory.nexon.net/rankings/job-ranking/explorer/warrior?pageIndex=&rebootIndex1=6&page_index=1")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en,zh-CN;q=0.9,zh;q=0.8,en-GB;q=0.7,en-US;q=0.6,zh-TW;q=0.5")

	//resp, err := client.Get("https://maplestory.nexon.net/api/ranking?id=job&id2=1&page_index=1")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var officialRank OfficialRank
	err = json.Unmarshal(body, &officialRank)
	fmt.Println(officialRank)
}
