package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mainjzb/MapleQQBotPlug/config"
	"github.com/mainjzb/MapleQQBotPlug/service"
	"github.com/mattn/go-ieproxy"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
var loginQQ int

var configFile = flag.String("config", "./config.yaml", "配置文件路径")

func UnescapeUnicode(raw string) string {
	str, _ := strconv.Unquote(strings.Replace(strconv.Quote(raw), `\\u`, `\u`, -1))
	return str
}

func test() {
	res := CheckMapleGG("TigerZY", 0)
	fmt.Println(res)
}

func main() {
	// test()
	// 0----------------------------------------------

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
				log.Println("WS Read Error ", err.Error())
			}
			time.Sleep(time.Second * 20)
		}
	}()

	for {
		m, err := ws.Read(msg)
		if err != nil {
			logrus.Error(err)
			continue
		}

		if string(msg[:m]) == "NewEvent" {
			client := &http.Client{}
			// post要提交的数据
			DataUrlVal := url.Values{}
			for key, val := range Data {
				DataUrlVal.Add(key, val)
			}
			req, err := http.NewRequest("POST", "http://127.0.0.1:10429/geteventv2", strings.NewReader(DataUrlVal.Encode()))
			if err != nil {
				return
			}
			// 提交请求
			resp, err := client.Do(req)
			if err != nil {
				logrus.Error(err)
				continue
			}
			// 读取返回值
			resultByte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logrus.Error(err)
				continue
			}
			err = resp.Body.Close()
			if err != nil {
				logrus.Error(err)
				continue
			}
			// fmt.Println(string(resultByte))
			// resultString := string(resultByte)
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
					logrus.Info(cleanJson)
					err = json.Unmarshal([]byte(cleanJson), &message)
					if err != nil {
						logrus.Error(err)
						continue
					}

					if message.Status == "OK" {
						for i, v := range message.Events {
							if v.Type != "GroupMsg" {
								continue
							}

							// fmt.Println(message.Events[i].Msg.Text)
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

func init() {
	// 初始化缓存
	charactersCatch = make(map[string]CharInfoResult)
	charactersLevelRank = make(map[string]string)
	checkNumberOfTimes = make(map[int]int)

	// config
	conf := config.Init(*configFile)
	loginQQ = int(conf.LoginQQ)

	// 初始化日志
	writer, err := rotatelogs.New(
		conf.LogFile+".%Y%m%d%H",
		rotatelogs.WithLinkName(conf.LogFile),
		rotatelogs.WithMaxAge(time.Duration(3)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(36)*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(writer)
	logrus.SetLevel(logrus.InfoLevel)

	/*
		if file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
			logrus.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			logrus.SetOutput(os.Stdout)
			logrus.Error(err)
		}*/

	// 初始化数据库

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

	service.NewCasbin(gdb)
	//
	cronClient := cron.New()
	cronClient.AddFunc("10 18 * * ?", sendADs)
	errID, err := cronClient.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	if err != nil {
		logrus.Error(err, errID)
	}

	cronClient.Start()

	go CheckMaplestoryInfo()
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
	// body, err := ioutil.ReadAll(resp.Body)
	// content := string(body)
	// content = content[strings.Index(content,"</head>") + 7:]

	// doc, err := html.ParseWithOptions(strings.NewReader(content), html.ParseOptionEnableScripting(false))
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".card-threadlist .card-body li .subject").Each(func(i int, s *goquery.Selection) {
		// fmt.Println(i, s.Text())
		// For each item found, get the band and titl
		band, ok := s.Last().Find("a").First().Attr("href")
		// title := s.Find("i").Text()
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
	// body, err := ioutil.ReadAll(resp.Body)
	// content := string(body)
	// content = content[strings.Index(content,"</head>") + 7:]

	// doc, err := html.ParseWithOptions(strings.NewReader(content), html.ParseOptionEnableScripting(false))
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

// 查询官网信息
func CheckMaplestoryInfo() {
	var content = ""
	for {
		http.DefaultTransport.(*http.Transport).Proxy = ieproxy.GetProxyFunc()
		resp, err := http.Get("http://maplestory.nexon.net/news/all")
		if err != nil {
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			continue
		}
		doc.Find(".news-container li").First().Each(func(i int, s *goquery.Selection) {
			band, ok := s.Find(".text h3 a").Attr("href")
			band = "https://maplestory.nexon.net" + band
			// title := s.Find("i").Text()
			if ok {
				if content == "" {
					content = band
				} else if content != band {
					content = band
					for _, v := range config.Instance.OfficialNoticeQQGroup {
						SendGroupMsg(loginQQ, int(v), "停一下！ 都停一下！ 百科有话说！ 官网发布新内容了！")
						SendGroupMsg(loginQQ, int(v), content)
					}
				}
			}
		})
		err = resp.Body.Close()
		if err != nil {
			log.Println("CheckMaplestoryInfo: ", err.Error())
			continue
		}

		time.Sleep(5 * time.Minute)
	}
}

func sendADs() {
	dir, _ := os.Getwd()

	image1 := GetGroupImage(loginQQ, 698931513, 1, dir+"\\Botimage\\elemAd1.jpg") // elemAd1.jpg
	image2 := GetGroupImage(loginQQ, 698931513, 1, dir+"\\Botimage\\elemAd2.jpg")
	SendGroupMsg(loginQQ, 698931513, "嘀嘀嘀！干饭时间到了！冲冲冲!")
	SendGroupMsg(loginQQ, 698931513, image1+image2)
	/*
		groupList := GetAllGroupList(loginQQ)
		for _, group := range groupList {
			image1 := GetGroupImage(loginQQ, group, 1, dir+"\\Botimage\\elemAd1.jpg") //elemAd1.jpg
			image2 := GetGroupImage(loginQQ, group, 1, dir+"\\Botimage\\elemAd2.jpg")
			SendGroupMsg(loginQQ, group, "嘀嘀嘀！干饭时间到了！冲冲冲!")
			SendGroupMsg(loginQQ, group, image1+image2)
			time.Sleep(time.Second * 3)
		}
	*/
}

func route(loginQQ, fromGroup, fromQQ int, groupMessage string) {
	if IsBanQQ(fromQQ) {
		return
	}
	AllCommand := []struct {
		Function func(loginQQ, fromGroup, fromQQ int, groupMessage string) bool
		Pre      []string
	}{
		{QueryContainQuestion, []string{"模糊查询问题", "模糊搜索问题"}},
		{QueryQuestion, []string{"文本查询问题", "查询问题"}},
		{DeleteAnswer, []string{"删除答案"}},
		{bindCharacter, []string{"查询绑定", "绑定查询", "绑定"}},
		{CheckClassRank, []string{"查询"}},
		{AddQuestion, []string{"问"}},
		{ChangeAuth, []string{"修改权限"}},
		{DeleteAuth, []string{"删除权限"}},
		{QADeleteQuestion, []string{"文本删除问题"}},
		{Wiki, []string{"百科"}},
	}

	for _, command := range AllCommand {
		groupMessage, ok := IsPrefix(groupMessage, command.Pre...)
		if ok && command.Function(loginQQ, fromGroup, fromQQ, groupMessage) {
			break
		}
	}
}
