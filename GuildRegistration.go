package main

import (
	"database/sql"
	"sort"
	"strconv"
	"strings"
	"time"
)

type person struct {
	QQ    int
	Score int
}

func GuildCheck(loginQQ, fromGroup, fromQQ int, msg string) (result int) {
	msgRune := []rune(msg)
	if fromGroup != 732888280 && fromGroup != 667082876 {
		return 0
	}
	checkGroup := 732888280
	dt := time.Now().Local().Add(-8 * time.Hour)
	date := dt.Format("2006-01-02")

	if strings.HasPrefix(msg, "跑旗录入") || strings.HasPrefix(msg, "跑旗记录") ||
		strings.HasPrefix(msg, "记录跑旗") || strings.HasPrefix(msg, "录入跑旗") ||
		strings.HasPrefix(msg, "登记跑旗") || strings.HasPrefix(msg, "跑旗登记") {

		scoreStr := strings.TrimSpace(string(msgRune[4:]))
		score, err := strconv.Atoi(scoreStr)
		if err != nil {
			SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n录入失败: 不要输入无关文字")
			return 1
		}
		if score < 10 || score%5 != 0 || score > 100 {
			SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n录入失败: 请确认自己的分数,不要作弊！")
			return 1
		}

		//rows, _ := db.Query("SELECT * FROM Log where Data=\""+date + "\" And QQ=" + fromQQ)
		sqlStmt := `SELECT Score1, Score2, Score3, Score4, Score5  FROM Log WHERE Date = ? AND QQ = ?`
		var historyScore [5]int
		if err := db.QueryRow(sqlStmt, date, fromQQ).Scan(&historyScore[0], &historyScore[1], &historyScore[2], &historyScore[3], &historyScore[4]); err != nil {
			if err != sql.ErrNoRows {
				SendGroupMsg(loginQQ, fromGroup, "数据库错误"+err.Error())
				return 1
			} else {
				stmt, _ := db.Prepare("INSERT INTO Log(QQ, Date, Score1, Score2, Score3, Score4, Score5) values(?, ?, ?, ?, ?, ?, ?)")
				_, err := stmt.Exec(fromQQ, date, score, 0, 0, 0, 0)
				if err != nil {
					SendGroupMsg(loginQQ, fromGroup, "数据插入库错误"+err.Error())
					return 0
				}
				SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n录入成功！\n今日积分："+scoreStr+"\n请将截图加入群相册，方便管理查询")
				return 1
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
						return 0
					}
				}
				if index != 0 {
					currentScore += "+"
				}
				currentScore += strconv.Itoa(historyScore[index])
				if number == 0 {
					SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n录入成功："+scoreStr+"\n今日积分："+currentScore+"\n请将截图加入群相册，方便管理查询")
					return 1
				}
			}
			SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n你已经录入超过5次，请不要作弊")

		}
	} else if msg == "跑旗查询" || msg == "查询跑旗" || msg == "跑旗查看" || msg == "查看跑旗" || msg == "我的跑旗" || msg == "跑旗我的" {
		var historyScore [5]int
		sqlStmt := `SELECT Score1, Score2, Score3, Score4, Score5  FROM Log WHERE Date = ? AND QQ = ?`
		if err := db.QueryRow(sqlStmt, date, fromQQ).Scan(&historyScore[0], &historyScore[1], &historyScore[2], &historyScore[3], &historyScore[4]); err != nil {
			if err != sql.ErrNoRows {
				SendGroupMsg(loginQQ, fromGroup, "数据库错误"+err.Error())
				return 1
			} else {
				weekNumber := time.Duration(time.Now().Add(-8*time.Hour).Weekday()) - 1
				//美国人周日是第一天 weekNumber == -1
				if weekNumber == -1 {
					weekNumber = 6
				}
				_, _, persionsMap := guildThisWeek(loginQQ, fromGroup, fromQQ, weekNumber)
				SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n今日暂无跑旗记录\n本周累计积分："+strconv.Itoa(persionsMap[fromQQ]))
				return 1
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
					return 1
				}
				if index != 0 {
					currentScore += "+"
				}
				currentScore += strconv.Itoa(historyScore[index])
			}
			SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n今日已录入5次，请不要作弊！\n今日积分："+currentScore)
			return 1
		}
	} else if msg == "跑旗今天" || msg == "跑旗本日" || msg == "跑旗今日" || msg == "本日跑旗" || msg == "今日跑旗" || msg == "今天跑旗" {
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

			resultContent += strconv.Itoa(per.Score) + "  " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
		}
		//读取游戏总分
		nowScore, beforeScore, err := getTotalScore()
		if err != nil {
			SendGroupMsg(loginQQ, fromGroup, err.Error())
			return 1
		}
		if nowScore == -1 {
			SendGroupMsg(loginQQ, fromGroup, "今天游戏总分还未录入！\n"+resultContent)
			return 1
		}
		if beforeScore == -1 {
			SendGroupMsg(loginQQ, fromGroup, "昨天没有管理录入游戏总分！\n当前游戏总分："+strconv.Itoa(nowScore)+"\n"+resultContent)
			return 1
		}

		SendGroupMsg(loginQQ, fromGroup, "今天游戏总分："+strconv.Itoa(nowScore-beforeScore)+"\n"+resultContent)
		return 1
	} else if msg == "跑旗本周" || msg == "跑旗这周" || msg == "本周跑旗" || msg == "这周跑旗" {
		weekNumber := time.Duration(dt.Weekday()) - 1
		//美国人周日是第一天 weekNumber == -1
		if weekNumber == -1 {
			weekNumber = 6
		}
		sumScore, persions, _ := guildThisWeek(loginQQ, fromGroup, fromQQ, weekNumber)
		resultContent := "本周玩家总分：" + strconv.Itoa(sumScore) + "\n"
		for _, per := range persions {
			resultContent += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
		}
		SendGroupMsg(loginQQ, fromGroup, resultContent)
		return 1
	} else if msg == "跑旗上周" || msg == "上周跑旗" {
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
			//persions[index]
			if per.Score > 50 {
				resultContent1 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
			} else {
				resultContent2 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
			}
		}
		SendGroupMsg(loginQQ, fromGroup, resultContent1)
		SendGroupMsg(loginQQ, fromGroup, resultContent2)
		return 1
	} else if msg == "跑旗七天" || msg == "七天跑旗" || msg == "7天跑旗" || msg == "跑旗7天" ||
		msg == "最近7天跑旗" || msg == "最近七天跑旗" || msg == "跑旗最近7天" || msg == "跑旗最近七天" {
		sumScore, persons, _ := guildThisWeek(loginQQ, fromGroup, fromQQ, time.Duration(6))
		resultContent := "最近7天玩家总分：" + strconv.Itoa(sumScore) + "\n"
		for _, per := range persons {
			resultContent += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
		}
		SendGroupMsg(loginQQ, fromGroup, resultContent)
		return 1
	} else if msg == "跑旗删除" || msg == "删除跑旗" {
		stmt, _ := db.Prepare("DELETE FROM Log WHERE Date = ? AND QQ = ?")
		stmtResult, _ := stmt.Exec(date, fromQQ)
		if size, _ := stmtResult.RowsAffected(); size == 0 {
			SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n您今日未登记过。\n提示：晚上8点前的数据都算为昨天。")
		} else {
			SendGroupMsg(loginQQ, fromGroup, "[@"+strconv.Itoa(fromQQ)+"]\n删除记录成功")

		}
		return 1
	} else if strings.HasPrefix(msg, "跑旗总分录入") || strings.HasPrefix(msg, "录入跑旗总分") || strings.HasPrefix(msg, "跑旗录入总分") ||
		strings.HasPrefix(msg, "跑旗总分登记") || strings.HasPrefix(msg, "登记跑旗总分") || strings.HasPrefix(msg, "跑旗登记总分") ||
		strings.HasPrefix(msg, "跑旗总分记录") || strings.HasPrefix(msg, "记录跑旗总分") || strings.HasPrefix(msg, "跑旗记录总分") {
		//管理员检测
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
			return 0
		}

		//admin := cqp.GetGroupMemberInfo(fromGroup, fromQQ, false)
		//if admin.Auth == 1{
		//	SendGroupMsg(loginQQ, fromGroup, "[@"+ strconv.Itoa(int(fromQQ)) + "]\n权限不足，请让管理员帮忙操作，普通成员不要进行此操作" )
		//	return 1
		//}

		scoreStr := strings.TrimSpace(string(msgRune[6:]))
		score, err := strconv.Atoi(scoreStr)
		if err != nil {
			SendGroupMsg(loginQQ, fromGroup, "录入失败: 不要输入无关文字")
			return 1
		}

		yesterdyDate := time.Now().Add(-32 * time.Hour).Format("2006-01-02")

		stmsRsult, err := db.Exec(`INSERT OR IGNORE INTO SumScore (date_time, score, qq) VALUES (?, ?, ?)`, date, score, fromQQ)
		if err != nil {
			SendGroupMsg(loginQQ, fromGroup, err.Error())
			return 1
		}
		if size, _ := stmsRsult.RowsAffected(); size == 0 {
			socore := 0
			sqlStmt := `SELECT score FROM SumScore WHERE date_time = ?`
			if err := db.QueryRow(sqlStmt, date).Scan(&socore); err == nil {
				SendGroupMsg(loginQQ, fromGroup, "总分录入失败：已有管理员录入。\n当前游戏总分："+strconv.Itoa(socore))
			} else {
				SendGroupMsg(loginQQ, fromGroup, "数据库错误："+err.Error())
			}
			return 1
		} else {
			var totalScore int
			// 今天不是星期1，则查询昨天的记录
			if time.Duration(dt.Weekday()) != 1 {
				sqlStmt := `SELECT score  FROM SumScore WHERE date_time = ?`
				if err := db.QueryRow(sqlStmt, yesterdyDate).Scan(&totalScore); err != nil {
					if err == sql.ErrNoRows {
						sumScore, _, _ := guildThisWeek(loginQQ, fromGroup, fromQQ, time.Duration(0))
						SendGroupMsg(loginQQ, fromGroup, "昨日无管理录入，目前游戏总分："+strconv.Itoa(score)+"  今日玩家总分："+strconv.Itoa(sumScore))
						return 1
					} else {
						SendGroupMsg(loginQQ, fromGroup, "数据库错误："+err.Error())
						return 1
					}
				}
			}

			sumScore, _, _ := guildThisWeek(loginQQ, fromGroup, fromQQ, time.Duration(0))
			SendGroupMsg(loginQQ, fromGroup, "录入成功\n今日总分："+strconv.Itoa(score-totalScore)+"  今日玩家总分："+strconv.Itoa(sumScore))
			return 1
		}
	} else if msg == "跑旗总分删除" || msg == "删除跑旗总分" || msg == "跑旗删除总分" {
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
			return 0
		}

		stmt, _ := db.Prepare("DELETE FROM sumScore WHERE date_time = ?")
		stmsRsult, _ := stmt.Exec(date)
		if size, _ := stmsRsult.RowsAffected(); size == 0 {
			SendGroupMsg(loginQQ, fromGroup, "今日未登记过总分。\n提示：晚上8点前的数据都算为昨天。")
		} else {
			SendGroupMsg(loginQQ, fromGroup, "删除总分成功")
		}
		return 1
	} else if msg == "本轮跑旗" {
		weekNumber := (time.Duration(dt.Weekday()) + 3) % 7 // 周四是一轮的第一天 0
		//美国人周日是0第一天 weekNumber == -1
		sumScore, persions, _ := guildThisWeek(loginQQ, fromGroup, fromQQ, weekNumber)
		resultContent := "周四是第一天\n本轮玩家总分：" + strconv.Itoa(sumScore) + "\n"
		resultContent1 := resultContent
		resultContent2 := ""
		for _, per := range persions {
			//persions[index]
			if per.Score > 50 {
				resultContent1 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
			} else {
				resultContent2 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
			}
		}

		SendGroupMsg(loginQQ, fromGroup, resultContent1)
		SendGroupMsg(loginQQ, fromGroup, resultContent2)
		return 1
	} else if msg == "上轮跑旗" {
		weekNumber := (time.Duration(dt.Weekday()) + 3) % 7
		sumScore, persions, _ := guildLastWeek(loginQQ, fromGroup, fromQQ, weekNumber)
		resultContent := "周四是第一天\n上轮玩家总分：" + strconv.Itoa(sumScore) + "\n"
		resultContent1 := resultContent
		resultContent2 := ""
		for _, per := range persions {
			//persions[index]
			if per.Score > 50 {
				resultContent1 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
			} else {
				resultContent2 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
			}
		}
		SendGroupMsg(loginQQ, fromGroup, resultContent1)
		SendGroupMsg(loginQQ, fromGroup, resultContent2)
		return 1
	} else if msg == "跑旗上上周" || msg == "上上周跑旗" {
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
				resultContent1 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
			} else {
				resultContent2 += strconv.Itoa(per.Score) + "   " + string([]rune(GetGropCard(loginQQ, checkGroup, per.QQ))[:6]) + "\n"
			}
		}
		SendGroupMsg(loginQQ, fromGroup, resultContent1)
		SendGroupMsg(loginQQ, fromGroup, resultContent2)
		return 1
	} else if strings.HasPrefix(msg, "查询跑旗20") {
		if len(msg) < 20 {
			return 1
		}
		msg = strings.TrimSpace(string(msgRune[4:]))
		date1 := strings.TrimSpace(msg[:10])
		date2 := strings.TrimSpace(msg[10:])
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
	return 0
}

func guildThisWeek(loginQQ, fromGroup, fromQQ int, weekNumber time.Duration) (sumScore int, persions []person, persionsMap map[int]int) {

	persions = make([]person, 0, 30)
	persionsMap = make(map[int]int)

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

		persionsMap[qq] += qqScore

		sumScore += qqScore
	}

	err = rows.Err()
	if err != nil {
		SendGroupMsg(loginQQ, fromGroup, "数据库错误3:"+err.Error())
	}

	for key, value := range persionsMap {
		persions = append(persions, person{key, value})
	}
	sort.Slice(persions, func(i, j int) bool {
		return persions[i].Score > persions[j].Score
	})

	return

}

func guildLastWeek(loginQQ, fromGroup, fromQQ int, weekNumber time.Duration) (sumScore int, persions []person, persionsMap map[int]int) {

	persions = make([]person, 0, 30)
	persionsMap = make(map[int]int)

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

		persionsMap[qq] += qqScore

		sumScore += qqScore
	}

	err = rows.Err()
	if err != nil {
		SendGroupMsg(loginQQ, fromGroup, "数据库错误3:"+err.Error())
	}

	for key, value := range persionsMap {
		persions = append(persions, person{key, value})
	}
	sort.Slice(persions, func(i, j int) bool {
		return persions[i].Score > persions[j].Score
	})

	return

}

func guildFlagRaceQuery(loginQQ, fromGroup, fromQQ int, dt1, dt2 string) (sumScore int, persions []person, persionsMap map[int]int) {

	persions = make([]person, 0, 30)
	persionsMap = make(map[int]int)

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

		persionsMap[qq] += qqScore

		sumScore += qqScore
	}

	err = rows.Err()
	if err != nil {
		SendGroupMsg(loginQQ, fromGroup, "数据库错误3:"+err.Error())
	}

	for key, value := range persionsMap {
		persions = append(persions, person{key, value})
	}
	sort.Slice(persions, func(i, j int) bool {
		return persions[i].Score > persions[j].Score
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
