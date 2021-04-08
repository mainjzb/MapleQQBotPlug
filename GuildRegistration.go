package main

import (
	"database/sql"
	"github.com/mainjzb/MapleQQBotPlug/config"
	"sort"
	"strconv"
	"time"
)

type person struct {
	QQ    int
	Score int
}

func GuildCheck(loginQQ, fromGroup, fromQQ int, groupMessage string) (result int) {
	isGuildFlagRace := false
	for _, v := range config.Instance.GuildFlagRaceQQGroup {
		if v == fromGroup {
			isGuildFlagRace = true
			break
		}
	}
	if isGuildFlagRace == false {
		return 0
	}

	//date := dt.Format("2006-01-02")

	AllCommand := []struct {
		Function func(loginQQ, fromGroup, fromQQ int, groupMessage string) bool
		Pre      []string
	}{
		{recordScore, []string{"跑旗录入", "跑旗记录", "记录跑旗", "录入跑旗", "登记跑旗", "跑旗登记"}},
		{removeMyScore, []string{"跑旗删除", "删除跑旗"}},
		{queryMyScore, []string{"跑旗查询", "查询跑旗", "跑旗查看", "查看跑旗", "我的跑旗", "跑旗我的"}},
		{queryTodayScore, []string{"跑旗今天", "跑旗本日", "跑旗今日", "本日跑旗", "今日跑旗", "今天跑旗"}},
		{queryThisWeekScore, []string{"跑旗本周", "跑旗这周", "本周跑旗", "这周跑旗"}},
		{queryLastWeekScore, []string{"跑旗上周", "上周跑旗"}},
		{queryLastLastWeekScore, []string{"跑旗上上周", "上上周跑旗"}},
		{query7DayScore, []string{"跑旗七天", "七天跑旗", "7天跑旗", "跑旗7天"}},
		{recordSumScore, []string{"跑旗总分录入", "录入跑旗总分", "跑旗录入总分", "跑旗总分登记", "登记跑旗总分", "跑旗登记总分", "跑旗总分记录", "记录跑旗总分", "跑旗记录总分"}},
		{removeSumScore, []string{"跑旗总分删除", "删除跑旗总分", "跑旗删除总分"}},
		{queryThisRoundScore, []string{"本轮跑旗", "这轮跑旗", "跑旗本轮", "跑旗这轮"}},
		{queryLastRoundScore, []string{"上轮跑旗", "跑旗上轮"}},
	}
	for _, command := range AllCommand {
		groupMessage, ok := IsPrefix(groupMessage, command.Pre...)
		if ok && command.Function(loginQQ, fromGroup, fromQQ, groupMessage) {
			return 1
		}
	}
	/*
		if strings.HasPrefix(groupMessage, "查询跑旗20") {
			if len(groupMessage) < 20 {
				return 1
			}
			groupMessage = strings.TrimSpace(string(msgRune[4:]))
			date1 := strings.TrimSpace(groupMessage[:10])
			date2 := strings.TrimSpace(groupMessage[10:])
			sumScore, persions, _ := guildFlagRaceQuery(loginQQ, fromGroup, fromQQ, date1, date2)
			resultContent := date1 + " " + date2 + "\n玩家总分：" + strconv.Itoa(sumScore)
			resultContent1 := resultContent
			resultContent2 := ""
			MedianScore := persions[len(persions)/2].Score
			for _, per := range persions {
				if per.Score > MedianScore {
					resultContent1 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
				} else {
					resultContent2 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
				}
			}
			SendGroupMsg(loginQQ, fromGroup, resultContent1)
			SendGroupMsg(loginQQ, fromGroup, resultContent2)
		}
	*/
	return 0
}

func guildThisWeek(loginQQ, fromGroup, fromQQ int, weekNumber time.Duration) (sumScore int, persons []person, personsMap map[int]int) {

	persons = make([]person, 0, 30)
	personsMap = make(map[int]int)

	var historyScore [5]int

	dt := time.Now().Add(-8*time.Hour - weekNumber*24*time.Hour)
	date := dt.Format("2006-01-02")

	rows, err := db.Query(`SELECT Score1, Score2, Score3, Score4, Score5, QQ from Log where Date>=? `, date)
	if err != nil || rows == nil {
		SendGroupMsg(loginQQ, fromGroup, "数据库错误1")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var qq int
		err := rows.Scan(&historyScore[0], &historyScore[1], &historyScore[2], &historyScore[3], &historyScore[4], &qq)
		if err != nil {
			SendGroupMsg(loginQQ, fromGroup, "数据库错误2:"+err.Error())
		}

		qqScore := historyScore[0] + historyScore[1] + historyScore[2] + historyScore[3] + historyScore[4]

		personsMap[qq] += qqScore

		sumScore += qqScore
	}

	err = rows.Err()
	if err != nil {
		SendGroupMsg(loginQQ, fromGroup, "数据库错误3:"+err.Error())
	}

	for key, value := range personsMap {
		persons = append(persons, person{key, value})
	}
	sort.Slice(persons, func(i, j int) bool {
		return persons[i].Score > persons[j].Score
	})

	return

}

func guildLastWeek(loginQQ, fromGroup, fromQQ int, weekNumber time.Duration) (sumScore int, persons []person, personsMap map[int]int) {

	persons = make([]person, 0, 30)
	personsMap = make(map[int]int)

	var historyScore [5]int

	dt1 := time.Now().Add(-8*time.Hour - (weekNumber+7)*24*time.Hour).Format("2006-01-02")
	dt2 := time.Now().Add(-8*time.Hour - (weekNumber)*24*time.Hour).Format("2006-01-02")

	rows, err := db.Query(`SELECT Score1, Score2, Score3, Score4, Score5, QQ from Log where Date>=? And Date <? `, dt1, dt2)
	if err != nil || rows == nil {
		SendGroupMsg(loginQQ, fromGroup, "数据库错误1:")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var qq int
		err := rows.Scan(&historyScore[0], &historyScore[1], &historyScore[2], &historyScore[3], &historyScore[4], &qq)
		if err != nil {
			SendGroupMsg(loginQQ, fromGroup, "数据库错误2:"+err.Error())
		}

		qqScore := historyScore[0] + historyScore[1] + historyScore[2] + historyScore[3] + historyScore[4]

		personsMap[qq] += qqScore

		sumScore += qqScore
	}

	err = rows.Err()
	if err != nil {
		SendGroupMsg(loginQQ, fromGroup, "数据库错误3:"+err.Error())
	}

	for key, value := range personsMap {
		persons = append(persons, person{key, value})
	}
	sort.Slice(persons, func(i, j int) bool {
		return persons[i].Score > persons[j].Score
	})

	return

}

func guildFlagRaceQuery(loginQQ, fromGroup, fromQQ int, dt1, dt2 string) (sumScore int, persons []person, personsMap map[int]int) {

	persons = make([]person, 0, 30)
	personsMap = make(map[int]int)

	var historyScore [5]int

	//dt1 := time.Now().Add( -8 * time.Hour - (weekNumber + 7) * 24 * time.Hour ).Format("2006-01-02")
	//dt2 := time.Now().Add( -8 * time.Hour - (weekNumber) * 24 * time.Hour ).Format("2006-01-02")
	layout := "2006-01-02"
	t1, err := time.Parse(layout, dt1)
	if err != nil {
		return
	}
	t2, err := time.Parse(layout, dt2)
	if err != nil {
		return
	}
	if t2.Sub(t1).Hours() > 15*24 {
		return
	}

	rows, err := db.Query(`SELECT Score1, Score2, Score3, Score4, Score5, QQ from Log where Date>=? And Date <= ? `, dt1, dt2)
	if err != nil || rows == nil {
		SendGroupMsg(loginQQ, fromGroup, "数据库错误1:")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var qq int
		err := rows.Scan(&historyScore[0], &historyScore[1], &historyScore[2], &historyScore[3], &historyScore[4], &qq)
		if err != nil {
			SendGroupMsg(loginQQ, fromGroup, "数据库错误2:"+err.Error())
			return
		}

		qqScore := historyScore[0] + historyScore[1] + historyScore[2] + historyScore[3] + historyScore[4]

		personsMap[qq] += qqScore

		sumScore += qqScore
	}

	err = rows.Err()
	if err != nil {
		SendGroupMsg(loginQQ, fromGroup, "数据库错误3:"+err.Error())
	}

	for key, value := range personsMap {
		persons = append(persons, person{key, value})
	}
	sort.Slice(persons, func(i, j int) bool {
		return persons[i].Score > persons[j].Score
	})

	return
}

func getTotalScore() (nowScore, yesterdayScore int, err error) {

	//读取今天游戏总分
	currentDate := time.Now().Add(-8 * time.Hour).Format("2006-01-02")
	if err = db.QueryRow(`SELECT score FROM SumScore WHERE date_time = ?`, currentDate).Scan(&nowScore); err != nil {
		if err == sql.ErrNoRows {
			nowScore = -1
			err = nil
		}
		return
	}

	//如果今天不是星期1，则获取昨天
	if time.Duration(time.Now().Add(-8*time.Hour).Weekday()) != 1 {
		//读取昨天游戏总分
		yesterdayDate := time.Now().Add(-(8 + 24) * time.Hour).Format("2006-01-02")
		if err = db.QueryRow(`SELECT score FROM SumScore WHERE date_time = ?`, yesterdayDate).Scan(&yesterdayScore); err != nil {
			if err == sql.ErrNoRows {
				yesterdayScore = -1
				err = nil
			}
			return
		}
	}

	return
}

func recordScore(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	score, err := strconv.Atoi(groupMessage)
	if err != nil {
		SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n录入失败: 不要输入无关文字")
		return true
	}
	if score < 10 || score%5 != 0 || score > 100 {
		SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n录入失败: 请确认自己的分数,不要作弊！")
		return true
	}

	dt := time.Now().Local().Add(-8 * time.Hour)
	date := dt.Format("2006-01-02")

	sqlStmt := `SELECT Score1, Score2, Score3, Score4, Score5  FROM Log WHERE Date = ? AND QQ = ?`
	var historyScore [5]int
	if err := db.QueryRow(sqlStmt, date, fromQQ).Scan(&historyScore[0], &historyScore[1], &historyScore[2], &historyScore[3], &historyScore[4]); err != nil {
		if err != sql.ErrNoRows {
			SendGroupMsg(loginQQ, fromGroup, "数据库错误"+err.Error())
			return true
		} else {
			stmt, _ := db.Prepare("INSERT INTO Log(QQ, Date, Score1, Score2, Score3, Score4, Score5) values(?, ?, ?, ?, ?, ?, ?)")
			_, err := stmt.Exec(fromQQ, date, score, 0, 0, 0, 0)
			if err != nil {
				SendGroupMsg(loginQQ, fromGroup, "数据插入库错误"+err.Error())
				return true
			}
			SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n录入成功！\n今日积分："+groupMessage+"\n请将截图加入群相册，方便管理查询")
			return true
		}
	} else {
		currentScore := ""
		for index, number := range historyScore {
			if number == 0 {
				historyScore[index] = score
				stmt, _ := db.Prepare("UPDATE Log SET Score1 = ?, Score2 = ?, Score3 = ?, Score4 = ?, Score5 = ?  WHERE Date = ? AND QQ = ?  ")
				_, err := stmt.Exec(historyScore[0], historyScore[1], historyScore[2], historyScore[3], historyScore[4], date, fromQQ)
				if err != nil {
					SendGroupMsg(loginQQ, fromGroup, "数据插入库错误"+err.Error())
					return true
				}
			}
			if index != 0 {
				currentScore += "+"
			}
			currentScore += strconv.Itoa(historyScore[index])
			if number == 0 {
				SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n录入成功："+groupMessage+"\n今日积分："+currentScore+"\n请将截图加入群相册，方便管理查询")
				return true
			}
		}
		SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n你已经录入超过5次，请不要作弊")
		return true
	}
}

func recordSumScore(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	adminList := GetAdminList(loginQQ, fromGroup)
	flag := false
	for _, qq := range adminList {
		qqNumber, _ := strconv.Atoi(qq)
		if qqNumber == fromQQ {
			flag = true
			break
		}
	}
	if flag == false {
		return false
	}

	dt := time.Now().Local().Add(-8 * time.Hour)
	date := dt.Format("2006-01-02")

	score, err := strconv.Atoi(groupMessage)
	if err != nil {
		SendGroupMsg(loginQQ, fromGroup, "录入失败: 不要输入无关文字")
		return true
	}

	yesterdyDate := time.Now().Add(-32 * time.Hour).Format("2006-01-02")

	stmsRsult, err := db.Exec(`INSERT OR IGNORE INTO SumScore (date_time, score, qq) VALUES (?, ?, ?)`, date, score, fromQQ)
	if err != nil {
		SendGroupMsg(loginQQ, fromGroup, err.Error())
		return true
	}
	if size, _ := stmsRsult.RowsAffected(); size == 0 {
		socore := 0
		sqlStmt := `SELECT score FROM SumScore WHERE date_time = ?`
		if err := db.QueryRow(sqlStmt, date).Scan(&socore); err == nil {
			SendGroupMsg(loginQQ, fromGroup, "总分录入失败：已有管理员录入。\n当前游戏总分："+strconv.Itoa(socore))
		} else {
			SendGroupMsg(loginQQ, fromGroup, "数据库错误："+err.Error())
		}
		return true
	} else {
		var totalScore int
		// 今天不是星期1，则查询昨天的记录
		if time.Duration(dt.Weekday()) != 1 {
			sqlStmt := `SELECT score  FROM SumScore WHERE date_time = ?`
			if err := db.QueryRow(sqlStmt, yesterdyDate).Scan(&totalScore); err != nil {
				if err == sql.ErrNoRows {
					sumScore, _, _ := guildThisWeek(loginQQ, fromGroup, fromQQ, time.Duration(0))
					SendGroupMsg(loginQQ, fromGroup, "昨日无管理录入，目前游戏总分："+strconv.Itoa(score)+"  今日玩家总分："+strconv.Itoa(sumScore))
				} else {
					SendGroupMsg(loginQQ, fromGroup, "数据库错误："+err.Error())
				}
			}
		}

		sumScore, _, _ := guildThisWeek(loginQQ, fromGroup, fromQQ, time.Duration(0))
		SendGroupMsg(loginQQ, fromGroup, "录入成功\n今日总分："+strconv.Itoa(score-totalScore)+"  今日玩家总分："+strconv.Itoa(sumScore))
		return true
	}
}

func removeMyScore(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	dt := time.Now().Local().Add(-8 * time.Hour)
	date := dt.Format("2006-01-02")
	stmt, _ := db.Prepare("DELETE FROM Log WHERE Date = ? AND QQ = ?")
	stmtResult, _ := stmt.Exec(date, fromQQ)
	if size, _ := stmtResult.RowsAffected(); size == 0 {
		SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n您今日未登记过。\n提示：晚上8点前的数据都算为昨天。")
	} else {
		SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n删除记录成功")

	}
	return true
}

func removeSumScore(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	//管理员权限检测
	adminList := GetAdminList(loginQQ, fromGroup)
	flag := false
	for _, qq := range adminList {
		qqNumber, _ := strconv.Atoi(qq)
		if qqNumber == fromQQ {
			flag = true
			break
		}
	}
	if flag == false {
		return false
	}
	dt := time.Now().Local().Add(-8 * time.Hour)
	date := dt.Format("2006-01-02")

	stmt, _ := db.Prepare("DELETE FROM sumScore WHERE date_time = ?")
	stmsRsult, _ := stmt.Exec(date)
	if size, _ := stmsRsult.RowsAffected(); size == 0 {
		SendGroupMsg(loginQQ, fromGroup, "今日未登记过总分。\n提示：晚上8点前的数据都算为昨天。")
	} else {
		SendGroupMsg(loginQQ, fromGroup, "删除总分成功")
	}
	return true
}

func queryMyScore(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	var historyScore [5]int
	sqlStmt := `SELECT Score1, Score2, Score3, Score4, Score5  FROM Log WHERE Date = ? AND QQ = ?`
	dt := time.Now().Local().Add(-8 * time.Hour)
	date := dt.Format("2006-01-02")
	if err := db.QueryRow(sqlStmt, date, fromQQ).Scan(&historyScore[0], &historyScore[1], &historyScore[2], &historyScore[3], &historyScore[4]); err != nil {
		if err != sql.ErrNoRows {
			SendGroupMsg(loginQQ, fromGroup, "数据库错误"+err.Error())
			return true
		} else {
			weekNumber := time.Duration(time.Now().Add(-8*time.Hour).Weekday()) - 1
			//美国人周日是第一天 weekNumber == -1
			if weekNumber == -1 {
				weekNumber = 6
			}
			_, _, persionsMap := guildThisWeek(loginQQ, fromGroup, fromQQ, weekNumber)
			SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n今日暂无跑旗记录\n本周累计积分："+strconv.Itoa(persionsMap[fromQQ]))
			return true
		}
	} else {
		currentScore := ""
		var index, number int
		for index, number = range historyScore {
			if number == 0 {
				weekNumber := time.Duration(time.Now().Add(-8*time.Hour).Weekday()) - 1
				//美国人周日是第一天 weekNumber == -1
				if weekNumber == -1 {
					weekNumber = 6
				}
				_, _, persionsMap := guildThisWeek(loginQQ, fromGroup, fromQQ, weekNumber)

				SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n今日积分："+currentScore+"\n本周累计积分："+strconv.Itoa(persionsMap[fromQQ]))
				return true
			}
			if index != 0 {
				currentScore += "+"
			}
			currentScore += strconv.Itoa(historyScore[index])
		}
		SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n今日已录入5次，请不要作弊！\n今日积分："+currentScore)
		return true
	}
}

func queryTodayScore(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	sumScore, persions, _ := guildThisWeek(loginQQ, fromGroup, fromQQ, time.Duration(0))
	resultContent := "今天玩家总分：" + strconv.Itoa(sumScore) + "\n"
	for _, per := range persions {
		//persions[index]
		/*
			info := cqp.GetGroupMemberInfo(fromGroup, per.QQ, false)
			if len([]rune(info.Card)) > 8{
				info.Card = string([]rune(info.Card)[:8])
			}
		*/

		resultContent += strconv.Itoa(per.Score) + "  " + string([]rune(GetGropCard(loginQQ, config.Instance.GuildFlagRaceQQGroup[0], per.QQ))[:6]) + "\n"
	}
	//读取游戏总分
	nowScore, beforeScore, err := getTotalScore()
	if err != nil {
		SendGroupMsg(loginQQ, fromGroup, err.Error())
		return true
	}
	if nowScore == -1 {
		SendGroupMsg(loginQQ, fromGroup, "今天游戏总分还未录入！\n"+resultContent)
		return true
	}
	if beforeScore == -1 {
		SendGroupMsg(loginQQ, fromGroup, "昨天没有管理录入游戏总分！\n当前游戏总分："+strconv.Itoa(nowScore)+"\n"+resultContent)
		return true
	}

	SendGroupMsg(loginQQ, fromGroup, "今天游戏总分："+strconv.Itoa(nowScore-beforeScore)+"\n"+resultContent)
	return true
}

func queryThisWeekScore(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	dt := time.Now().Local().Add(-8 * time.Hour)
	weekNumber := time.Duration(dt.Weekday()) - 1
	//美国人周日是第一天 weekNumber == -1
	if weekNumber == -1 {
		weekNumber = 6
	}
	sumScore, persions, _ := guildThisWeek(loginQQ, fromGroup, fromQQ, weekNumber)
	resultContent := "本周玩家总分：" + strconv.Itoa(sumScore) + "\n"
	for _, per := range persions {
		resultContent += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, config.Instance.GuildFlagRaceQQGroup[0], per.QQ))[:6]) + "\n"
	}
	SendGroupMsg(loginQQ, fromGroup, resultContent)
	return true
}

func queryLastWeekScore(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	dt := time.Now().Local().Add(-8 * time.Hour)
	weekNumber := time.Duration(dt.Weekday()) - 1
	//美国人周日是第一天 weekNumber == -1
	if weekNumber == -1 {
		weekNumber = 6
	}
	sumScore, persions, _ := guildLastWeek(loginQQ, fromGroup, fromQQ, weekNumber)
	resultContent := "上周玩家总分：" + strconv.Itoa(sumScore) + "\n"
	resultContent1 := resultContent
	resultContent2 := ""
	for _, per := range persions {
		if per.Score > 50 {
			resultContent1 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, config.Instance.GuildFlagRaceQQGroup[0], per.QQ))[:6]) + "\n"
		} else {
			resultContent2 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, config.Instance.GuildFlagRaceQQGroup[0], per.QQ))[:6]) + "\n"
		}
	}
	SendGroupMsg(loginQQ, fromGroup, resultContent1)
	SendGroupMsg(loginQQ, fromGroup, resultContent2)
	return true
}

func queryLastLastWeekScore(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	dt := time.Now().Local().Add(-8 * time.Hour)
	weekNumber := time.Duration(dt.Weekday()) - 1
	//美国人周日是第一天 weekNumber == -1
	if weekNumber == -1 {
		weekNumber = 6
	}
	sumScore, persions, _ := guildLastWeek(loginQQ, fromGroup, fromQQ, weekNumber+7)
	resultContent := "上上周玩家总分：" + strconv.Itoa(sumScore) + "\n"
	resultContent1 := resultContent
	resultContent2 := ""
	for _, per := range persions {
		//persions[index]
		if per.Score > 50 {
			resultContent1 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, config.Instance.GuildFlagRaceQQGroup[0], per.QQ))[:6]) + "\n"
		} else {
			resultContent2 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, config.Instance.GuildFlagRaceQQGroup[0], per.QQ))[:6]) + "\n"
		}
	}
	SendGroupMsg(loginQQ, fromGroup, resultContent1)
	SendGroupMsg(loginQQ, fromGroup, resultContent2)
	return true
}

func query7DayScore(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	sumScore, persons, _ := guildThisWeek(loginQQ, fromGroup, fromQQ, time.Duration(6))
	resultContent := "最近7天玩家总分：" + strconv.Itoa(sumScore) + "\n"
	for _, per := range persons {
		resultContent += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, config.Instance.GuildFlagRaceQQGroup[0], per.QQ))[:6]) + "\n"
	}
	SendGroupMsg(loginQQ, fromGroup, resultContent)
	return true
}

func queryThisRoundScore(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	dt := time.Now().Local().Add(-8 * time.Hour)
	weekNumber := (time.Duration(dt.Weekday()) + 3) % 7 // 周四是一轮的第一天 0
	//美国人周日是0第一天 weekNumber == -1
	sumScore, persions, _ := guildThisWeek(loginQQ, fromGroup, fromQQ, weekNumber)
	resultContent := "周四是第一天\n本轮玩家总分：" + strconv.Itoa(sumScore) + "\n"
	resultContent1 := resultContent
	resultContent2 := ""
	for _, per := range persions {
		//persions[index]
		if per.Score > 50 {
			resultContent1 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, config.Instance.GuildFlagRaceQQGroup[0], per.QQ))[:6]) + "\n"
		} else {
			resultContent2 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, config.Instance.GuildFlagRaceQQGroup[0], per.QQ))[:6]) + "\n"
		}
	}

	SendGroupMsg(loginQQ, fromGroup, resultContent1)
	SendGroupMsg(loginQQ, fromGroup, resultContent2)
	return true
}

func queryLastRoundScore(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	dt := time.Now().Local().Add(-8 * time.Hour)
	weekNumber := (time.Duration(dt.Weekday()) + 3) % 7
	sumScore, persions, _ := guildLastWeek(loginQQ, fromGroup, fromQQ, weekNumber)
	resultContent := "周四是第一天\n上轮玩家总分：" + strconv.Itoa(sumScore) + "\n"
	resultContent1 := resultContent
	resultContent2 := ""
	for _, per := range persions {
		//persions[index]
		if per.Score > 50 {
			resultContent1 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, config.Instance.GuildFlagRaceQQGroup[0], per.QQ))[:6]) + "\n"
		} else {
			resultContent2 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, config.Instance.GuildFlagRaceQQGroup[0], per.QQ))[:6]) + "\n"
		}
	}
	SendGroupMsg(loginQQ, fromGroup, resultContent1)
	SendGroupMsg(loginQQ, fromGroup, resultContent2)
	return true
}
