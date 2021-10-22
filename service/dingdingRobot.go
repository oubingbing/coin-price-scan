package service

import (
"bytes"
"crypto/hmac"
"crypto/sha256"
"encoding/base64"
"encoding/json"
"fmt"
"io/ioutil"
"net/http"
"time"
)

type DingDingReq struct {
	MsgType  string           `json:"msgtype"`
	Text     DingDingText     `json:"text"`
	Markdown DingDingMarkdown `json:"markdown"`
	At       At               `json:"at"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
}

type DingDingText struct {
	Content string `json:"content"`
}

type DingDingMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func SendDingDing(accessToken, secret string, req DingDingReq) {
	t := time.Now().UnixNano() / 1e6
	sign := sign(t, secret)

	url := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%d&sign=%s", accessToken, t, sign)
	contentByte, _ := json.Marshal(&req)
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(contentByte))
	if err != nil {
		return
	}

	b, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(b))
}

func sign(t int64, secret string) string {
	strToHash := fmt.Sprintf("%d\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToHash))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}
