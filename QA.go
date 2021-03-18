package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func QAReply(loginQQ, fromGroup, fromQQ int, msg string) bool {

	//普通查询
	sqlStmt := `SELECT answer_template.text FROM question_template INNER JOIN answer_template ON  question_template.id = answer_template.parent   WHERE gc=? AND question_template.text=? AND question_template.match!=10;`
	//正则查询  match=10是正则
	sqlStmtReg := `SELECT question_template.text, answer_template.text FROM question_template INNER JOIN answer_template ON  question_template.id = answer_template.parent  WHERE gc=? AND question_template.match=10;`
	var answers = make([]string, 0)
	rows1, err := dbQA.Query(sqlStmt, fromGroup, msg)
	if err != nil {
		return false
	}
	rows2, err := dbQA.Query(sqlStmtReg, fromGroup)
	if err != nil {
		return false
	}
	defer rows1.Close()
	defer rows2.Close()

	for rows1.Next() {
		var ans string
		err := rows1.Scan(&ans)
		answers = append(answers, ans)
		if err != nil {
			return false
		}
	}
	for rows2.Next() {
		var qus, ans string
		err := rows2.Scan(&qus, &ans)
		if err != nil {
			return false
		}
		reg := regexp.MustCompile(qus)
		if reg.MatchString(msg) {
			answers = append(answers, ans)
		}
	}

	rand.Seed(time.Now().Unix())
	if len(answers) > 0 {
		answer := answers[rand.Intn(len(answers))]
		dir, _ := os.Getwd()
		reg, _ := regexp.Compile(`(\[pic,[^\[\]]*\])`)
		picList := reg.FindStringSubmatch(answer)
		for _, pic := range picList {
			image := GetGroupImage(loginQQ, fromGroup, 1, dir+"\\Botimage\\"+pic)
			answer = strings.ReplaceAll(answer, pic, image)
		}

		SendGroupMsg(loginQQ, fromGroup, answer)
		return true
	}

	return false
}

func QAAddMatch(loginQQ, fromGroup, fromQQ int, queMsg, ansMsg string) bool {
	/*
		1、查询这个问题是否存在，不存在则创建问题
		2、添加回答，父ID是问题ID
	*/

	if AuthorityQuery(fromGroup, fromQQ) < 2 {
		SendGroupMsg(loginQQ, fromGroup, "权限不足")
		return false
	}

	// 1、查询这个问题是否存在，不存在则创建问题
	parentID, lock := getQueMsgID(fromGroup, queMsg)
	if lock == 1 {
		SendGroupMsg(loginQQ, fromGroup, "添加失败，词条已锁定")
		return false
	}

	if parentID < 0 {
		//添加问题到数据库
		//sqlStmtQueryInsert := `INSERT INTO question_template VALUES (NULL,732888280,'dcw', 3, 0, 0,0,0,1139035718,1597033184);`
		//rows1, err := dbQA.Query(sqlStmtQueryInsert, fromGroup, msg)
		stmt, err := dbQA.Prepare(`INSERT INTO question_template VALUES (NULL, ?, ?, ?, 0, 0, 0, 0, ?, ?);`)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		defer stmt.Close()

		res, err := stmt.Exec(fromGroup, queMsg, len(queMsg), fromQQ, time.Now().Unix())
		if err != nil {
			return false
		}

		parentID, err = res.LastInsertId()
		if err != nil {
			return false
		}
		//SendGroupMsg(loginQQ, fromGroup, "问题已添加")
	}

	//理论检测 assert
	if parentID == -1 {
		return false
	}

	//2.1  获取答案里的图片pic

	reg, _ := regexp.Compile(`(\[pic,[^\[\]]*\])`)
	picList := reg.FindStringSubmatch(ansMsg)

	//2.2 下载图片存到本地
	for _, picHash := range picList {
		imagPath := GetPhotoUrl(loginQQ, fromGroup, picHash)
		dir, _ := os.Getwd()
		name := picHash[:42] + ";time=" + strconv.FormatInt(time.Now().Unix(), 10) + "]"
		pathName := dir + "\\Botimage\\" + name
		resp, err := http.Get(imagPath)
		if err != nil {
			return false
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		out, err := os.Create(pathName)
		if err != nil {
			return false
		}
		io.Copy(out, bytes.NewReader(body))
		out.Close()

		//2.3 修改答案
		ansMsg = strings.ReplaceAll(ansMsg, picHash, name)
	}

	//3、添加答案
	//INSERT INTO answer_template (id, parent, text,  created ,createdTime) VALUES (NULL, 500, '不知道', 99999999, 111111);
	stmt2, err := dbQA.Prepare(`INSERT INTO answer_template (id, parent, text,  created ,createdTime) VALUES (NULL, ?, ?, ?, ?);`)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt2.Close()

	res, err := stmt2.Exec(parentID, ansMsg, fromQQ, time.Now().Unix())
	if err != nil {
		return false
	}

	instertID, err := res.LastInsertId()
	if err != nil {
		return false
	}

	SendGroupMsg(loginQQ, fromGroup, "添加成功问题ID:"+strconv.FormatInt(parentID, 10)+", 答案ID:"+strconv.FormatInt(instertID, 10))

	//sqlStmtInsert := `INSERT INTO question_template VALUES (NULL,732888280,'dcw', 3, 0, 0,0,0,1139035718,1597033184);`

	return true
}

func QADeleteQuestion(loginQQ, fromGroup, fromQQ int, msg string) bool {

	if AuthorityQuery(fromGroup, fromQQ) < 2 {
		SendGroupMsg(loginQQ, fromGroup, "权限不足")
		return false
	}

	parentID, lock := getQueMsgID(fromGroup, msg)
	if parentID < 0 || lock == 1 {
		SendGroupMsg(loginQQ, fromGroup, "删除失败，问题不存在或已锁定")
		return false
	}

	//delete question
	stmtQue, err := dbQA.Prepare("delete from question_template where text=? AND lock = 0 AND match=0 AND gc = ?")
	if err != nil {
		return false
	}
	defer stmtQue.Close()

	resQue, err := stmtQue.Exec(msg, fromGroup)
	if err != nil {
		return false
	}

	_, err = resQue.RowsAffected()
	if err != nil {
		return false
	}

	// 删除答案包含的图片
	deleteAnsPic(parentID)

	//delete answer
	stmtAns, err := dbQA.Prepare("delete from answer_template where parent=? AND lock = 0")
	if err != nil {
		return false
	}
	defer stmtAns.Close()

	_, err = stmtAns.Exec(parentID)
	if err != nil {
		return false
	}

	SendGroupMsg(loginQQ, fromGroup, "成功删除问题ID："+strconv.FormatInt(parentID, 10))

	return true
}

func AuthorityQuery(fromGroup, fromQQ int) (right int) {
	sqlStmtQuery := `SELECT right FROM group_member WHERE gc=? AND qq=? ;`
	rows1, err := dbQA.Query(sqlStmtQuery, fromGroup, fromQQ)
	if err != nil {
		return 0
	}
	defer rows1.Close()
	if rows1.Next() {
		err := rows1.Scan(&right)
		if err != nil {
			return 0
		}
		return
	} else {
		return 0
	}
}

func getQueMsgID(fromGroup int, msg string) (parentID, lock int64) {
	parentID = -1
	//query parent ID
	sqlStmtQuery := `SELECT id, lock FROM question_template WHERE gc=? AND text=? AND question_template.match = 0;`
	rows1, err := dbQA.Query(sqlStmtQuery, fromGroup, msg)
	if err != nil {
		parentID = -1
		return
	}

	defer rows1.Close()
	if rows1.Next() {
		err := rows1.Scan(&parentID, &lock)
		if err != nil {
			parentID = -1
			return
		}
	} else {
		parentID = -1
		return
	}
	return
}

func getFuzzyQueMsgIDs(fromGroup int, msg string) (questions map[int64]string) {
	msg = "%" + msg + "%"
	questions = make(map[int64]string)
	//query parent ID
	sqlStmtQuery := `SELECT id, text FROM question_template WHERE gc=? AND match = 0 AND text LIKE ?;`
	rows1, err := dbQA.Query(sqlStmtQuery, fromGroup, msg)
	if err != nil {
		return
	}

	defer rows1.Close()
	for rows1.Next() {
		var parentID int64
		var answer string
		err := rows1.Scan(&parentID, &answer)
		if err != nil {
			return
		}
		questions[parentID] = answer
	}

	return
}

func deleteAnsPic(parentID int64) {
	sqlStmtQuery := `SELECT text FROM answer_template WHERE parent=? AND lock = 0;`
	rows1, err := dbQA.Query(sqlStmtQuery, parentID)
	if err != nil {
		return
	}

	defer rows1.Close()
	for rows1.Next() {
		var answer string
		err := rows1.Scan(&answer)
		if err != nil {
			return
		}
		reg, _ := regexp.Compile(`(\[pic,[^\[\]]*\])`)
		picList := reg.FindStringSubmatch(answer)
		for _, pic := range picList {
			dir, _ := os.Getwd()
			os.Remove(dir + "\\Botimage\\" + pic)
		}
	}
	return
}

func ChangeAuth(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	strList := strings.Fields(groupMessage)
	if len(strList) != 2 {
		return false
	}

	AuthQQ, err := strconv.Atoi(strings.TrimSpace(strList[0]))
	if err != nil {
		return true
	}
	AuthLevel, err := strconv.Atoi(strings.TrimSpace(strList[1]))

	if AuthorityQuery(fromGroup, fromQQ) >= 5 && AuthLevel < 5 || AuthorityQuery(fromGroup, fromQQ) >= 10 {

	} else {
		SendGroupMsg(loginQQ, fromGroup, "权限不足")
		return false
	}

	//删除权限
	stmtAns, err := dbQA.Prepare("delete from group_member where gc = ? AND qq = ?")
	if err != nil {
		return false
	}
	defer stmtAns.Close()

	_, err = stmtAns.Exec(fromGroup, AuthQQ)
	if err != nil {
		return false
	}

	//修改权限
	stmt2, err := dbQA.Prepare(`INSERT INTO group_member("id", "gc", "qq", "right", "silentState", "silentEndTime", "point", "remark", "createdTime") VALUES (NULL, ?, ?, ?, 0, 0, 0, '', ?);
`)
	if err != nil {
		fmt.Println(err.Error())
		return true
	}
	defer stmt2.Close()

	res, err := stmt2.Exec(fromGroup, AuthQQ, AuthLevel, time.Now().Unix())
	if err != nil {
		return true
	}

	if rowNum, err := res.RowsAffected(); rowNum == int64(1) && err == nil {
		SendGroupMsg(loginQQ, fromGroup, "权限修改成功")
	} else {
		SendGroupMsg(loginQQ, fromGroup, "权限修改失败")
		return true
	}

	return true
}

func DeleteAuth(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	AuthQQ, err := strconv.Atoi(groupMessage)
	if err != nil {
		return false
	}

	if AuthorityQuery(fromGroup, fromQQ) < 5 {
		SendGroupMsg(loginQQ, fromGroup, "权限不足")
		return true
	}

	stmtAns, err := dbQA.Prepare("delete from group_member where gc = ? AND qq = ?")
	if err != nil {
		return true
	}
	defer stmtAns.Close()

	_, err = stmtAns.Exec(fromGroup, AuthQQ)
	if err != nil {
		return true
	}

	SendGroupMsg(loginQQ, fromGroup, "成功删除权限:"+strconv.Itoa(AuthQQ))
	return true
}

func QueryQuestion(loginQQ, fromGroup, fromQQ int, Quemsg string) bool {
	if AuthorityQuery(fromGroup, fromQQ) < 2 {
		SendGroupMsg(loginQQ, fromGroup, "权限不足")
		return false
	}

	result := ""
	parentID, _ := getQueMsgID(fromGroup, Quemsg)
	if parentID < 0 {
		SendGroupMsg(loginQQ, fromGroup, "问题不存在")
		return false
	}

	result += "问题[" + strconv.FormatInt(parentID, 10) + "]：" + Quemsg + "\n"

	// 查询答案

	sqlStmt := `SELECT id, text FROM answer_template  WHERE parent=?;`
	rows1, err := dbQA.Query(sqlStmt, parentID)
	if err != nil {
		return false
	}
	defer rows1.Close()
	index := 0
	for rows1.Next() {
		index += 1
		if index == 10 {
			result += "\n\n" + "答案过多，其余答案已省略……"
			break
		}
		var ans string
		var id int
		err := rows1.Scan(&id, &ans)
		if err != nil {
			return false
		}
		result += "\n" + "答案[" + strconv.Itoa(id) + "]：" + ans
	}

	// 图片处理
	/*
		dir, _ := os.Getwd()
		reg, _ := regexp.Compile(`(\[pic,[^\[\]]*\])`)
		picList := reg.FindStringSubmatch(result)
		for _, pic := range picList{
			image := GetGroupImage(loginQQ, fromGroup,1, dir + "\\Botimage\\" + pic)
			result = strings.ReplaceAll(result, pic, image)
		}
	*/

	SendGroupMsg(loginQQ, fromGroup, result)

	return true
}

func QueryContainQuestion(loginQQ, fromGroup, fromQQ int, groupMessage string) bool {
	if AuthorityQuery(fromGroup, fromQQ) < 2 {
		SendGroupMsg(loginQQ, fromGroup, "权限不足")
		return false
	}

	result := ""
	questions := getFuzzyQueMsgIDs(fromGroup, groupMessage)

	if len(questions) == 0 {
		SendGroupMsg(loginQQ, fromGroup, "未搜索到任何问题")
		return false
	}

	for parentID, queMsg := range questions {

		result += "问题[" + strconv.FormatInt(parentID, 10) + "]：" + queMsg + "\n"

		// 查询答案

		sqlStmt := `SELECT id, text FROM answer_template  WHERE parent=?;`
		rows1, err := dbQA.Query(sqlStmt, parentID)
		if err != nil {
			return false
		}

		index := 0
		for rows1.Next() {
			index += 1
			if index == 10 {
				result += "\n\n" + "答案过多，其余答案已省略……"
				break
			}
			var ans string
			var id int
			err := rows1.Scan(&id, &ans)
			if err != nil {
				return false
			}
			result += "答案[" + strconv.Itoa(id) + "]：" + ans + "\n"
		}
		result += "\n"

		rows1.Close()

		// 图片处理
		/*
			dir, _ := os.Getwd()
			reg, _ := regexp.Compile(`(\[pic,[^\[\]]*\])`)
			picList := reg.FindStringSubmatch(result)
			for _, pic := range picList{
				image := GetGroupImage(loginQQ, fromGroup,1, dir + "\\Botimage\\" + pic)
				result = strings.ReplaceAll(result, pic, image)
			}
		*/
	}

	SendGroupMsg(loginQQ, fromGroup, strings.TrimSpace(result))

	return true
}

func DeleteAnswer(loginQQ, fromGroup, fromQQ int, msg string) bool {
	if AuthorityQuery(fromGroup, fromQQ) < 2 {
		SendGroupMsg(loginQQ, fromGroup, "权限不足")
		return false
	}
	parentID, err := strconv.Atoi(msg)
	if err != nil {
		return false
	}

	//delete answer
	stmtAns, err := dbQA.Prepare("delete from answer_template where id=? AND lock = 0")
	if err != nil {
		return false
	}
	defer stmtAns.Close()

	res, err := stmtAns.Exec(parentID)
	if err != nil {
		return false
	}

	row, err := res.RowsAffected()
	if err != nil {
		return false
	}

	if row == 0 {
		SendGroupMsg(loginQQ, fromGroup, "删除失败")
	} else {
		SendGroupMsg(loginQQ, fromGroup, "成功删除答案["+msg+"]")
	}

	return true
}
