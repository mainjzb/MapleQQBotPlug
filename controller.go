package main

import (
	"github.com/mainjzb/Golang-Bot/Translation"
	"github.com/mainjzb/Golang-Bot/calc"
	"strconv"
	"strings"
)

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

func Translate(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	SendGroupMsg(loginQQ, fromGroup, Translation.Trans(groupMessage))
	return true
}

func Wiki(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {

	if GuildCheck(loginQQ, fromGroup, fromQQ, groupMessage) != 0 {
		//跑旗相关问答
	} else if QAReply(loginQQ, fromGroup, fromQQ, groupMessage) {
		//问题数据库查询

	} else if IsDigitCalc(groupMessage) {
		//计算器
		answer, err := calc.Calc(groupMessage)
		if err == nil {
			SendGroupMsg(loginQQ, fromGroup, strconv.FormatFloat(answer, 'g', 12, 64))
		}
	} else if IsEnglish(groupMessage) {
		//翻译
		SendGroupMsg(loginQQ, fromGroup, Translation.Trans(groupMessage))
	}
	return true
}
