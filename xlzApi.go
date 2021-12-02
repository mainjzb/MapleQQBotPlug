package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type RetMessage struct {
	Ret string `json:"ret"`
}

type JsonGoupList struct {
	Ret  string `json:"ret"`
	List []struct {
		GroupID                int64  `json:"GroupID"`
		GIN                    int    `json:"GIN"`
		CFlag                  int    `json:"cFlag"`
		GroupInfoSeq           int    `json:"GroupInfoSeq"`
		DwGroupFlagExt         int    `json:"dwGroupFlagExt"`
		DwGroupRankSeq         int    `json:"dwGroupRankSeq"`
		DwCertificationType    int    `json:"dwCertificationType"`
		DwShutupTimestamp      int    `json:"dwShutupTimestamp"`
		DwMyShutupTimestamp    int    `json:"dwMyShutupTimestamp"`
		DwCmdUinUinFlag        int    `json:"dwCmdUinUinFlag"`
		DwAdditionalFlag       int    `json:"dwAdditionalFlag"`
		DwGroupTypeFlag        int    `json:"dwGroupTypeFlag"`
		DwGroupSecType         int    `json:"dwGroupSecType"`
		DwGroupSecTypeInfo     int    `json:"dwGroupSecTypeInfo"`
		DwGroupClassExt        int    `json:"dwGroupClassExt"`
		DwAppPrivilegeFlag     int    `json:"dwAppPrivilegeFlag"`
		DwSubscriptionUin      int    `json:"dwSubscriptionUin"`
		DwMemberNum            int    `json:"dwMemberNum"`
		DwMemberNumSeq         int    `json:"dwMemberNumSeq"`
		DwMemberCardSeq        int    `json:"dwMemberCardSeq"`
		DwGroupFlagExt3        int    `json:"dwGroupFlagExt3"`
		DwGroupOwnerUin        int    `json:"dwGroupOwnerUin"`
		CIsConfGroup           int    `json:"cIsConfGroup"`
		CIsModifyConfGroupFace int    `json:"cIsModifyConfGroupFace"`
		CIsModifyConfGroupName int    `json:"cIsModifyConfGroupName"`
		DwCmduinJoinTime       int    `json:"dwCmduinJoinTime"`
		StrGroupName           string `json:"strGroupName"`
		StrGroupMemo           string `json:"strGroupMemo"`
	} `json:"List"`
}

type Session struct {
	SessionID int `json:"session_id"`
}

func Sendprivatemsg(LoginQQ, toQQ int, text string) {
	apiUrl := "http://127.0.0.1:10429"
	resource := "/sendgroupmsg"
	data := url.Values{}
	data.Set("fromqq", strconv.Itoa(LoginQQ))
	data.Set("toqq", strconv.Itoa(toQQ))
	data.Set("text", text)

	dateCover := data.Encode()
	dateCover = strings.ReplaceAll(dateCover, "+", "%20")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"
	fmt.Println(u.String())
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(dateCover)) // URL-encoded payload

	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//resp, err := client.Do(r)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
}

func SendGroupMsg(LoginQQ, formGroup int, text string) {
	apiURL := "http://127.0.0.1:10429"
	resource := "/sendgroupmsg"
	data := url.Values{}
	data.Set("fromqq", strconv.Itoa(LoginQQ))
	data.Set("togroup", strconv.Itoa(formGroup))
	data.Set("text", text)

	dateCover := data.Encode()
	dateCover = strings.ReplaceAll(dateCover, "+", "%20")

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"
	fmt.Println(u.String())
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(dateCover)) // URL-encoded payload

	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//resp, err := client.Do(r)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	//fmt.Println(resp.Status)
}

func GetGropCard(LoginQQ, formGrop, toQQ int) string {
	apiURL := "http://127.0.0.1:10429"
	resource := "/getgroupcard"
	data := url.Values{}
	data.Set("fromqq", strconv.Itoa(LoginQQ))
	data.Set("group", strconv.Itoa(formGrop))
	data.Set("toqq", strconv.Itoa(toQQ))

	dateCover := data.Encode()
	dateCover = strings.ReplaceAll(dateCover, "+", "%20")

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"
	fmt.Println(u.String())
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(dateCover)) // URL-encoded payload

	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//resp, err := client.Do(r)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	card := RetMessage{}
	resultByte, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resultByte))
	err = json.Unmarshal(resultByte, &card)
	if err != nil {
		fmt.Println(err.Error())
	}

	return card.Ret
}

func GetAdminList(LoginQQ, formGrop int) []string {
	//{"ret":"212427942\r\n404346435\r\n451782246\r\n517682597\r\n565436128\r\n594900769\r\n1021268160\r\n1732351842\r\n2637020248\r\n"}
	apiUrl := "http://127.0.0.1:10429"
	resource := "/getgroupmgrlist"
	data := url.Values{}
	data.Set("fromqq", strconv.Itoa(LoginQQ))
	data.Set("group", strconv.Itoa(formGrop))

	dateCover := data.Encode()
	dateCover = strings.ReplaceAll(dateCover, "+", "%20")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"
	fmt.Println(u.String())
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(dateCover)) // URL-encoded payload

	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	retMessage := RetMessage{}
	resultByte, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resultByte))
	err = json.Unmarshal(resultByte, &retMessage)
	if err != nil {
		fmt.Println(err.Error())
	}
	list := strings.Split(retMessage.Ret, "\r\n")
	newList := []string{}
	for _, v := range list {
		if v != "" {
			newList = append(newList, v)
		}
	}
	return newList
}

func Allocsession() int {
	client := &http.Client{}
	DataUrlVal := url.Values{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:10429/allocsession", strings.NewReader(DataUrlVal.Encode()))
	if err != nil {
		return 0
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("error: get request")
		log.Fatal(err)
	}
	defer resp.Body.Close()

	resultByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error: return value")
		log.Fatal(err)
	}
	fmt.Println(string(resultByte))

	sess := Session{}
	err = json.Unmarshal(resultByte, &sess)

	return sess.SessionID
}

func GetGroupImage(loginQQ, formGrop, fromtype int, ImagePath string) string {
	apiURL := "http://127.0.0.1:10429"
	resource := "/sendgrouppic"
	data := url.Values{}
	data.Set("fromqq", strconv.Itoa(loginQQ))
	data.Set("togroup", strconv.Itoa(formGrop))
	if fromtype == 2 {
		// 网络图片
		data.Set("fromtype", "2")
		data.Set("url", ImagePath)
	} else if fromtype == 1 {
		// 本地图片
		data.Set("fromtype", "1")
		data.Set("path", ImagePath)
	}

	dateCover := data.Encode()
	dateCover = strings.ReplaceAll(dateCover, "+", "%20")

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"
	fmt.Println(u.String())
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(dateCover)) // URL-encoded payload

	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	//resp, err := client.Do(r)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer resp.Body.Close()

	card := RetMessage{}
	resultByte, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resultByte))
	err = json.Unmarshal(resultByte, &card)
	if err != nil {
		fmt.Println(err.Error())
	}

	return card.Ret
}

func GetPhotoURL(loginQQ, formGrop int, pic string) string {
	apiURL := "http://127.0.0.1:10429"
	resource := "/getphotourl"

	data := url.Values{}
	data.Set("fromqq", strconv.Itoa(loginQQ))
	data.Set("group", strconv.Itoa(formGrop))
	data.Set("photo", pic)

	dateCover := data.Encode()
	dateCover = strings.ReplaceAll(dateCover, "+", "%20")

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"
	fmt.Println(u.String())
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(dateCover)) // URL-encoded payload

	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	card := RetMessage{}
	resultByte, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resultByte))
	err = json.Unmarshal(resultByte, &card)
	if err != nil {
		fmt.Println(err.Error())
	}

	return card.Ret
}

func GetAllGroupList(LoginQQ int) []int {
	//{"ret":"212427942\r\n404346435\r\n451782246\r\n517682597\r\n565436128\r\n594900769\r\n1021268160\r\n1732351842\r\n2637020248\r\n"}
	apiUrl := "http://127.0.0.1:10429"
	resource := "/getgrouplist"
	data := url.Values{}
	data.Set("logonqq", strconv.Itoa(LoginQQ))

	dateCover := data.Encode()
	dateCover = strings.ReplaceAll(dateCover, "+", "%20")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"
	fmt.Println(u.String())
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(dateCover)) // URL-encoded payload

	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	goupList := JsonGoupList{}
	resultByte, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resultByte))
	err = json.Unmarshal(resultByte, &goupList)
	if err != nil {
		fmt.Println(err.Error())
	}

	newList := []int{}
	for _, v := range goupList.List {
		newList = append(newList, v.GIN)
	}
	return newList
}
