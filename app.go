package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mainjzb/Golang-Bot/Translation"
	"github.com/mainjzb/Golang-Bot/calc"
	"github.com/mattn/go-ieproxy"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/websocket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type QQMessage struct {
	Status string `json:"status"`
	Events []struct {
		Type   string `json:"Type"`
		FromQQ struct {
			UIN       int    `json:"UIN"`
			Card      string `json:"Card"`
			SpecTitle string `json:"SpecTitle"`
			Pos       struct {
				Lo int `json:"Lo"`
				La int `json:"La"`
			} `json:"Pos"`
		} `json:"FromQQ"`
		LogonQQ   int `json:"LogonQQ"`
		TimeStamp struct {
			Recv int `json:"Recv"`
			Send int `json:"Send"`
		} `json:"TimeStamp"`
		FromGroup struct {
			GIN  int    `json:"GIN"`
			Name string `json:"name"`
		} `json:"FromGroup"`
		Msg struct {
			Req       int    `json:"Req"`
			Random    int64  `json:"Random"`
			SubType   int    `json:"SubType"`
			AppID     int    `json:"AppID"`
			Text      string `json:"Text"`
			TextReply string `json:"Text_Reply"`
			BubbleID  int    `json:"BubbleID"`
		} `json:"Msg"`
		File struct {
			ID   string `json:"ID"`
			MD5  string `json:"MD5"`
			Name string `json:"Name"`
			Size int64  `json:"Size"`
		} `json:"File"`
	} `json:"events"`
}

var db *sql.DB
var gdb *gorm.DB
var dbQA *sql.DB
var banLists []string
var adminQQs []int64
var qqGroupNumber int64 = 856067978
var loginQQ int = 2137511870

func UnescapeUnicode(raw string) string {
	str, _ := strconv.Unquote(strings.Replace(strconv.Quote(raw), `\\u`, `\u`, -1))
	return str
}

func test() {

	//QueryClassRanking("213", 1, "")
	CheckClassRank(1, 1, 1, "夜光第一")
}

//go:generate cqcfg -c .
// cqp: 名称: MapleRank
// cqp: 版本: 1.0.0:1
// cqp: 作者: mao9
// cqp: 简介: 一个超棒的Go语言插件Demo，它会回复你的私聊消息~
func main() {
	//test()
	//0----------------------------------------------

	ws, err := websocket.Dial("ws://127.0.0.1:10429", "", "ws://127.0.0.1:10429")
	if err != nil {
		log.Fatal(err)
	}

	Data := make(map[string]string)
	Data["sessid"] = strconv.Itoa(Allocsession())
	var msg = make([]byte, 512)

	go func() {
		for {
			_, err := ws.Write([]byte("123"))

			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second * 20)
		}
	}()

	for {
		m, err := ws.Read(msg)

		if err != nil {
			//log.Fatal(err)
		}

		if string(msg[:m]) == "NewEvent" {
			client := &http.Client{}
			//post要提交的数据
			DataUrlVal := url.Values{}
			for key, val := range Data {
				DataUrlVal.Add(key, val)
			}
			req, err := http.NewRequest("POST", "http://127.0.0.1:10429/geteventv2", strings.NewReader(DataUrlVal.Encode()))
			if err != nil {
				return
			}
			//提交请求
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal("error: get request")
				log.Fatal(err)
			}
			//读取返回值
			resultByte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("error: return value")
				log.Fatal(err)
			}
			resp.Body.Close()
			//fmt.Println(string(resultByte))
			//resultString := string(resultByte)
			jsonList := bytes.Split(resultByte, []byte("\r"))
			for _, jsonText := range jsonList {

				if len(string(jsonText)) > 5 {
					var message QQMessage

					cleanJson := strings.Map(func(r rune) rune {
						if unicode.IsGraphic(r) {
							return r
						}
						return -1
					}, string(jsonText))
					cleanJson = strings.ReplaceAll(cleanJson, `\'`, `'`)
					fmt.Println("start" + string(cleanJson) + "end")
					err = json.Unmarshal([]byte(cleanJson), &message)
					if err != nil {
						log.Fatal("error: json")
						log.Fatal(err)
						fmt.Println("答应失败")
						continue
					}

					if message.Status == "OK" {
						for i, _ := range message.Events {
							if message.Events[i].Type != "GroupMsg" {
								continue
							}

							//fmt.Println(message.Events[i].Msg.Text)
							groupMessage := UnescapeUnicode(message.Events[i].Msg.Text)
							groupMessage = strings.TrimSpace(groupMessage)
							gropFromGroup := message.Events[i].FromGroup.GIN
							gropFromQQ := message.Events[i].FromQQ.UIN
							loginQQ := message.Events[i].LogonQQ
							if message.Events[i].FromQQ.UIN != loginQQ {
								route(loginQQ, gropFromGroup, gropFromQQ, groupMessage)
							}
						}
					}
				} else {
					fmt.Println("jsontestt:" + string(jsonText))
				}
			}

		}

	}

}

func Baike(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {

	if GuildCheck(loginQQ, fromGroup, fromQQ, groupMessage) != 0 {
		//跑旗相关问答
	} else if QAReply(loginQQ, fromGroup, fromQQ, groupMessage) {
		//问题数据库查询

	} else if IsDigitCalc(groupMessage) {
		//计算器
		answer, error := calc.Calc(groupMessage)
		if error == nil {
			SendGroupMsg(loginQQ, fromGroup, strconv.FormatFloat(answer, 'g', 12, 64))
		}
	} else if IsEnglish(groupMessage) {
		//翻译
		SendGroupMsg(loginQQ, fromGroup, Translation.Trans(groupMessage))
	}
	return true
}

func Translate(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	SendGroupMsg(loginQQ, fromGroup, Translation.Trans(groupMessage))
	return true
}

func AddQuestion(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {

	messageStr := strings.Split(groupMessage, "\r答")
	if len(messageStr) == 2 {
		QAAddMatch(loginQQ, fromGroup, fromQQ, strings.TrimSpace(messageStr[0]), strings.TrimSpace(messageStr[1]))
		return true
	}

	messageStr = strings.Split(groupMessage, "\n答")
	if len(messageStr) == 2 {
		QAAddMatch(loginQQ, fromGroup, fromQQ, strings.TrimSpace(messageStr[0]), strings.TrimSpace(messageStr[1]))
		return true
	}

	messageStr = strings.Split(groupMessage, " 答")
	if len(messageStr) == 2 {
		QAAddMatch(loginQQ, fromGroup, fromQQ, strings.TrimSpace(messageStr[0]), strings.TrimSpace(messageStr[1]))
		return true
	}
	return false
}

func init() {
	// 初始化缓存
	charactersCatch = make(map[string]charInfoResult)
	charactersLevelRank = make(map[string]string)
	checkNumberOfTimes = make(map[int]int)
	/*
		classUrl = map[string]string{
			"Warrior":         "https://maplestory.nexon.net/rankings/job-ranking/explorer/warrior?rebootIndex=1&pageIndex=",
			"Magician":        "https://maplestory.nexon.net/rankings/job-ranking/explorer/magician?rebootIndex=1&pageIndex=",
			"Bowman":          "https://maplestory.nexon.net/rankings/job-ranking/explorer/bowman?rebootIndex=1&pageIndex=",
			"Thief":           "https://maplestory.nexon.net/rankings/job-ranking/explorer/thief?rebootIndex=1&pageIndex=",
			"Pirate":          "https://maplestory.nexon.net/rankings/job-ranking/explorer/pirate?rebootIndex=1&pageIndex=",
			"Evan":            "https://maplestory.nexon.net/rankings/job-ranking/evan/null?rebootIndex=1&pageIndex=",
			"Aran":            "https://maplestory.nexon.net/rankings/job-ranking/aran/aran?rebootIndex=1&pageIndex=",
			"Mercedes":        "https://maplestory.nexon.net/rankings/job-ranking/mercedes/null?rebootIndex=1&pageIndex=",
			"Luminous":        "https://maplestory.nexon.net/rankings/job-ranking/luminous/null?rebootIndex=1&pageIndex=",
			"Jett":            "https://maplestory.nexon.net/rankings/job-ranking/jett/null?rebootIndex=1&pageIndex=",
			"Mihile":          "https://maplestory.nexon.net/rankings/job-ranking/mihile/null?rebootIndex=1&pageIndex=",
			"Kaiser":          "https://maplestory.nexon.net/rankings/job-ranking/kaiser/null?rebootIndex=1&pageIndex=",
			"Angelic Buster":  "https://maplestory.nexon.net/rankings/job-ranking/angelic-buster/null?rebootIndex=1&pageIndex=",
			"Phantom":         "https://maplestory.nexon.net/rankings/job-ranking/phantom/null?rebootIndex=1&pageIndex=",
			"Hayato":          "https://maplestory.nexon.net/rankings/job-ranking/sengoku/hayato?rebootIndex=1&pageIndex=",
			"Kanna":           "https://maplestory.nexon.net/rankings/job-ranking/sengoku/kanna?rebootIndex=1&pageIndex=",
			"Xenon":           "https://maplestory.nexon.net/rankings/job-ranking/xenon/null?rebootIndex=1&pageIndex=",
			"Zero":            "https://maplestory.nexon.net/rankings/job-ranking/zero/null?rebootIndex=1&pageIndex=",
			"Beast Tamer":     "https://maplestory.nexon.net/rankings/job-ranking/beast-tamer/null?rebootIndex=1&pageIndex=",
			"Shade":           "https://maplestory.nexon.net/rankings/job-ranking/shade/null?rebootIndex=1&pageIndex=",
			"Kinesis":         "https://maplestory.nexon.net/rankings/job-ranking/kinesis/null?rebootIndex=1&pageIndex=",
			"Blaster":         "https://maplestory.nexon.net/rankings/job-ranking/blaster/null?rebootIndex=1&pageIndex=",
			"Cadena":          "https://maplestory.nexon.net/rankings/job-ranking/cadena/null?rebootIndex=1&pageIndex=",
			"Illium":          "https://maplestory.nexon.net/rankings/job-ranking/illium/null?rebootIndex=1&pageIndex=",
			"Ark":             "https://maplestory.nexon.net/rankings/job-ranking/ark/null?rebootIndex=1&pageIndex=",
			"Pathfinder":      "https://maplestory.nexon.net/rankings/job-ranking/pathfinder/null?rebootIndex=1&pageIndex=",
			"Hoyoung":         "https://maplestory.nexon.net/rankings/job-ranking/hoyoung/null?rebootIndex=1&pageIndex=",
			"Adele":           "https://maplestory.nexon.net/rankings/job-ranking/adele/null?rebootIndex=1&pageIndex=",
			"Dawn Warrior":    "https://maplestory.nexon.net/rankings/job-ranking/cygnus-knights/dawn-warrior?rebootIndex=1&pageIndex=",
			"Blaze Wizard":    "https://maplestory.nexon.net/rankings/job-ranking/cygnus-knights/blaze-wizard?rebootIndex=1&pageIndex=",
			"Wind Archer":     "https://maplestory.nexon.net/rankings/job-ranking/cygnus-knights/wind-archer?rebootIndex=1&pageIndex=",
			"Night Walker":    "https://maplestory.nexon.net/rankings/job-ranking/cygnus-knights/night-walker?rebootIndex=1&pageIndex=",
			"Thunder Breaker": "https://maplestory.nexon.net/rankings/job-ranking/cygnus-knights/thunder-breaker?rebootIndex=1&pageIndex=",
			"Demon Slayer":    "https://maplestory.nexon.net/rankings/job-ranking/resistance/demon-slayer?rebootIndex=1&pageIndex=",
			"Battle Mage":     "https://maplestory.nexon.net/rankings/job-ranking/resistance/battle-mage?rebootIndex=1&pageIndex=",
			"Wild Hunter":     "https://maplestory.nexon.net/rankings/job-ranking/resistance/wild-hunter?rebootIndex=1&pageIndex=",
			"Mechanicr":       "https://maplestory.nexon.net/rankings/job-ranking/resistance/mechanic?rebootIndex=1&pageIndex=",
			"Demon Avenger":   "https://maplestory.nexon.net/rankings/job-ranking/resistance/demon-avenger?rebootIndex=1&pageIndex=",
			"联盟":              "https://maplestory.nexon.net/rankings/legion/reboot-(na)?rebootIndex=0&pageIndex=",
		}
	*/
	classUrl = map[string]string{
		"Warrior":         "https://maplestory.nexon.net/api/ranking?id=job&id2=1&rebootIndex=1&page_index=",
		"Magician":        "https://maplestory.nexon.net/api/ranking?id=job&id2=2&rebootIndex=1&page_index=",
		"Bowman":          "https://maplestory.nexon.net/api/ranking?id=job&id2=3&rebootIndex=1&page_index=",
		"Thief":           "https://maplestory.nexon.net/api/ranking?id=job&id2=4&rebootIndex=1&page_index=",
		"Pirate":          "https://maplestory.nexon.net/api/ranking?id=job&id2=5&rebootIndex=1&page_index=",
		"Aran":            "https://maplestory.nexon.net/api/ranking?id=job&id2=20&rebootIndex=1&page_index=",
		"Evan":            "https://maplestory.nexon.net/api/ranking?id=job&id2=22&rebootIndex=1&page_index=",
		"Mercedes":        "https://maplestory.nexon.net/api/ranking?id=job&id2=23&rebootIndex=1&page_index=",
		"Phantom":         "https://maplestory.nexon.net/api/ranking?id=job&id2=24&rebootIndex=1&page_index=",
		"Jett":            "https://maplestory.nexon.net/api/ranking?id=job&id2=201&rebootIndex=1&page_index=",
		"Mihile":          "https://maplestory.nexon.net/api/ranking?id=job&id2=202&rebootIndex=1&page_index=",
		"Luminous":        "https://maplestory.nexon.net/api/ranking?id=job&id2=203&rebootIndex=1&page_index=",
		"Kaiser":          "https://maplestory.nexon.net/api/ranking?id=job&id2=204&rebootIndex=1&page_index=",
		"Angelic Buster":  "https://maplestory.nexon.net/api/ranking?id=job&id2=205&rebootIndex=1&page_index=",
		"Hayato":          "https://maplestory.nexon.net/api/ranking?id=job&id2=206&rebootIndex=1&page_index=",
		"Kanna":           "https://maplestory.nexon.net/api/ranking?id=job&id2=207&rebootIndex=1&page_index=",
		"Xenon":           "https://maplestory.nexon.net/api/ranking?id=job&id2=208rebootIndex=1&page_index=",
		"Zero":            "https://maplestory.nexon.net/api/ranking?id=job&id2=210&rebootIndex=1&page_index=",
		"Beast Tamer":     "https://maplestory.nexon.net/api/ranking?id=job&id2=211&rebootIndex=1&page_index=",
		"Shade":           "https://maplestory.nexon.net/api/ranking?id=job&id2=212&rebootIndex=1&page_index=",
		"Kinesis":         "https://maplestory.nexon.net/api/ranking?id=job&id2=214&rebootIndex=1&page_index=",
		"Blaster":         "https://maplestory.nexon.net/api/ranking?id=job&id2=215&rebootIndex=1&page_index=",
		"Cadena":          "https://maplestory.nexon.net/api/ranking?id=job&id2=216&rebootIndex=1&page_index=",
		"Illium":          "https://maplestory.nexon.net/api/ranking?id=job&id2=217&rebootIndex=1&page_index=",
		"Ark":             "https://maplestory.nexon.net/api/ranking?id=job&id2=218&rebootIndex=1&page_index=",
		"Pathfinder":      "https://maplestory.nexon.net/api/ranking?id=job&id2=219&rebootIndex=1&page_index=",
		"Hoyoung":         "https://maplestory.nexon.net/api/ranking?id=job&id2=220&rebootIndex=1&page_index=",
		"Adele":           "https://maplestory.nexon.net/api/ranking?id=job&id2=221&rebootIndex=1&page_index=",
		"Dawn Warrior":    "https://maplestory.nexon.net/api/ranking?id=job&id2=11&rebootIndex=1&page_index=",
		"Blaze Wizard":    "https://maplestory.nexon.net/api/ranking?id=job&id2=12&rebootIndex=1&page_index=",
		"Wind Archer":     "https://maplestory.nexon.net/api/ranking?id=job&id2=13&rebootIndex=1&page_index=",
		"Night Walker":    "https://maplestory.nexon.net/api/ranking?id=job&id2=14&rebootIndex=1&page_index=",
		"Thunder Breaker": "https://maplestory.nexon.net/api/ranking?id=job&id2=15&rebootIndex=1&page_index=",
		"Demon Slayer":    "https://maplestory.nexon.net/api/ranking?id=job&id2=31&rebootIndex=1&page_index=",
		"Battle Mage":     "https://maplestory.nexon.net/api/ranking?id=job&id2=32&rebootIndex=1&page_index=",
		"Wild Hunter":     "https://maplestory.nexon.net/api/ranking?id=job&id2=33&rebootIndex=1&page_index=",
		"Mechanicr":       "https://maplestory.nexon.net/api/ranking?id=job&id2=34&rebootIndex=1&page_index=",
		"Demon Avenger":   "https://maplestory.nexon.net/api/ranking?id=job&id2=209&rebootIndex=1&page_index=",
		"联盟":              "https://maplestory.nexon.net/api/ranking?id=legion&id2=45&page_index=",
	}
	onEnable()
}

func onEnable() int32 {

	//// 初始化数据库
	var err error
	dir, _ := os.Getwd()
	fmt.Println(dir)

	db, err = sql.Open("sqlite3", dir+"\\MapleMiao.db")
	if err != nil {
		fmt.Println(err.Error())
	}

	dbQA, err = sql.Open("sqlite3", dir+"\\QA.db")
	if err != nil {
		fmt.Println(err.Error())
	}

	gdb, err = gorm.Open(sqlite.Open("MapleMiao.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
	}

	/*
		////读取禁言黑名单语句
		for rows, _ := db.Query("SELECT * FROM banList"); rows.Next(); {
			var rowString string
			rows.Scan(&rowString)
			banLists = append(banLists, rowString)
		}

		////读取管理员QQ
		for rows, _ := db.Query("SELECT * FROM Admin"); rows.Next(); {
			var rowString int64
			rows.Scan(&rowString)
			adminQQs = append(adminQQs, rowString)
		}
	*/

	go CheckMaplestoryInfo()

	return 0
}

func FindDB() (result string) {
	rows, err := db.Query("SELECT * FROM Stack;")
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var QQ int64
		err = rows.Scan(&QQ)
		if err != nil {
			fmt.Println(err)
		}
		result += fmt.Sprintf("[CQ:at,qq=%v]\n", QQ)
	}
	fmt.Println(result)
	return
}

func GetMaplestoryVersionInfo(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	var result string
	resp, err := http.Get("https://bbs.gjfmxd.com/forum-1-1.htm?tagids=2")
	if err != nil {
		// handle error
		SendGroupMsg(loginQQ, fromGroup, "获取失败，可能网络存在问题")
		return true
	}
	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//content := string(body)
	//content = content[strings.Index(content,"</head>") + 7:]

	//doc, err := html.ParseWithOptions(strings.NewReader(content), html.ParseOptionEnableScripting(false))
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".card-threadlist .card-body li .subject").Each(func(i int, s *goquery.Selection) {
		//fmt.Println(i, s.Text())
		// For each item found, get the band and titl
		band, ok := s.Last().Find("a").First().Attr("href")
		//title := s.Find("i").Text()
		if ok {
			result = "https://bbs.gjfmxd.com/" + band
		}
	})

	SendGroupMsg(loginQQ, fromGroup, result)
	return true
}

func GetMaplestoryMaintainInfo(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	var result string
	resp, err := http.Get("https://bbs.gjfmxd.com/forum-1-1.htm?tagids=5")
	if err != nil {
		// handle error
		SendGroupMsg(loginQQ, fromGroup, "获取失败")
		return true
	}
	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//content := string(body)
	//content = content[strings.Index(content,"</head>") + 7:]

	//doc, err := html.ParseWithOptions(strings.NewReader(content), html.ParseOptionEnableScripting(false))
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
		return true
	}

	ss := doc.Find(".card-threadlist .card-body li .media-body .subject").Last().Find("a").First()
	band, ok := ss.Attr("href")
	if ss.Find("span").Size() != 0 {
		if ok {
			result = "https://bbs.gjfmxd.com/" + band
		} else {
			result = "查询失败"
		}
	} else {
		result = "暂无维护信息"
	}

	SendGroupMsg(loginQQ, fromGroup, result)
	return true
}

func CheckMaplestoryInfo() {
	var content string
	for {
		//cqp.AddLog(cqp.Info,"查询官网更新","info")
		http.DefaultTransport.(*http.Transport).Proxy = ieproxy.GetProxyFunc()
		resp, err := http.Get("http://maplestory.nexon.net/news/all")
		if err != nil {
			// handle error
			//cqp.AddLog(cqp.Info,"查询官网更新","失败1")
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			//cqp.AddLog(cqp.Info,"查询官网更新","失败2")
			continue
		}
		doc.Find(".news-container li").First().Each(func(i int, s *goquery.Selection) {
			fmt.Println(i)
			// For each item found, get the band and titl

			band, ok := s.Find(".text h3 a").Attr("href")
			band = "http://maplestory.nexon.net/" + band
			//title := s.Find("i").Text()
			if ok {
				if content == "" {
					content = band
				} else if content != band {
					content = band

					Sendprivatemsg(loginQQ, 212427942, content)

					SendGroupMsg(loginQQ, 318497715, "停一下！ 都停一下！ 百科有话说！ 官网发布新内容了！")
					SendGroupMsg(loginQQ, 318497715, content)

					SendGroupMsg(loginQQ, 732888280, "停一下！ 都停一下！ 百科有话说！ 官网发布新内容了！")
					SendGroupMsg(loginQQ, 732888280, content)

					//cqp.SendPrivateMsg(212427942, content)

				}
				//fmt.Println(content)
				//result = "https://bbs.gjfmxd.com/"+band
			}
		})
		time.Sleep(5 * time.Minute)
	}
}

func IsDigitCalc(data string) bool {
	digit := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", " ", "(", ")", "+", "-", "*", "/", "."}
	flag := false
	for _, i := range data {
		flag = false
		for _, item := range digit {
			if string(i) == item {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

func IsEnglish(data string) bool {
	for i := range data {
		if !(31 < data[i] && data[i] < 123) {
			return false
		}
	}
	return true
}

func IsPrefix(groupMessage string, prefixs ...string) (string, bool) {
	for _, value := range prefixs {
		if strings.HasPrefix(groupMessage, value) {
			return strings.TrimSpace(groupMessage[len(value):]), true
		}
	}

	return groupMessage, false
}

func route(loginQQ, fromGroup, fromQQ int, groupMessage string) {
	AllCommand := []struct {
		Function func(loginQQ, fromGroup, fromQQ int, groupMessage string) bool
		Pre      []string
	}{
		{QueryContainQuestion, []string{"模糊查询问题", "模糊搜索问题"}},
		{QueryQuestion, []string{"文本查询问题"}},
		{DeleteAnswer, []string{"删除答案", "删除问题"}},
		{QueryWorldRanking, []string{"查询排名", "查询第", "查询等级第"}},
		{bindCharacter, []string{"查询绑定", "绑定查询", "绑定"}},
		{CheckClassRank, []string{"查询"}},
		{AddQuestion, []string{"问"}},
		{ChangeAuth, []string{"修改权限"}},
		{DeleteAuth, []string{"删除权限"}},
		{QADeleteQuestion, []string{"文本删除问题"}},
		{Translate, []string{"百科翻译", "百度翻译"}},
		{GetMaplestoryVersionInfo, []string{"百科版本内容", "百科版本活动", "百科版本", "百科活动"}},
		{GetMaplestoryMaintainInfo, []string{"百科维护"}},
		{Baike, []string{"百科"}},
	}

	for _, command := range AllCommand {
		groupMessage, ok := IsPrefix(groupMessage, command.Pre...)
		if ok && command.Function(loginQQ, fromGroup, fromQQ, groupMessage) {
			break
		}
	}
}
