package Translation


import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Md5(s string) string {//计算md5的值
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
func u2s(form string) (to string, err error) {//unicode转字符串
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}
func translate(word string) []byte {//调用api进行翻译
	salt := strconv.Itoa(time.Now().Second())
	appid := "20181019000221834"
	pwd := "JJmkoylJbdoYgs0KkkRF"
	data := make(url.Values)
	data["q"] = []string{word}
	data["from"] = []string{"auto"}
	data["to"] = []string{"auto"}
	data["appid"] = []string{appid}
	data["salt"] = []string{salt}
	s := appid + word + salt + pwd
	sign := Md5(s)
	data["sign"] = []string{sign}
	res, err := http.PostForm("http://api.fanyi.baidu.com/api/trans/vip/translate", data)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	//str := string(body)
	return body
}
func Trans(words string) string {//翻译函数
	body := translate(words)
	js, err := simplejson.NewJson(body)
	if err != nil {
		panic(err.Error())
	}
	dst := js.Get("trans_result").GetIndex(0).Get("dst").MustString()
	return dst
}
/*
func main() {
	words := "abandon"
	dst := trans(words)
	fmt.Println(dst)//打印结果
}*/
