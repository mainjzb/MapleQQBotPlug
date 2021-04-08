package main

import (
	"github.com/mainjzb/MapleQQBotPlug/Translation"
	"github.com/mainjzb/MapleQQBotPlug/calc"
	"github.com/mainjzb/MapleQQBotPlug/config"
	"github.com/shiguanghuxian/txai"
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
	AllCommand := []struct {
		Function func(loginQQ, fromGroup, fromQQ int, groupMessage string) bool
		Pre      []string
	}{
		{Translate, []string{"翻译", "翻译"}},
		{GetMaplestoryVersionInfo, []string{"版本内容", "版本活动", "版本"}},
		{GetMaplestoryMaintainInfo, []string{"维护"}}}

	for _, command := range AllCommand {
		groupMessage, ok := IsPrefix(groupMessage, command.Pre...)
		if ok && command.Function(loginQQ, fromGroup, fromQQ, groupMessage) {
			return true
		}
	}

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
	} else {
		txAi := txai.New(config.Instance.QQChatID, config.Instance.QQChatKey, true)
		response, err := txAi.NlpTextchatForText("123456", groupMessage)
		if err != nil {
			return true
		}
		SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]"+response.Data.Answer)
	}
	return true
}
