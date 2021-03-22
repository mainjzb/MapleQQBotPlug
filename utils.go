package main

import (
	"github.com/mattn/go-ieproxy"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func webGetRequest(url string) (result []byte, err error) {
	ieproxy.OverrideEnvWithStaticProxy()
	http.DefaultTransport.(*http.Transport).Proxy = ieproxy.GetProxyFunc()
	client := http.Client{
		Timeout: 4 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func IsPrefix(groupMessage string, prefixs ...string) (string, bool) {
	for _, value := range prefixs {
		if strings.HasPrefix(groupMessage, value) {
			return strings.TrimSpace(groupMessage[len(value):]), true
		}
	}

	return groupMessage, false
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
