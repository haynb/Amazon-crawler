package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type verify struct {
	username string
	password string
	userKey  string
}

func (v *verify) Login() error {
	loginUrl := "http://api.damagou.top/apiv1/login.html?username=" + v.username + "&password=" + v.password
	resp, err := http.Get(loginUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	v.userKey = string(data)
	return nil
}
func (v *verify) GetCode(base64Code string) (string, error) {
	if v.userKey == "" {
		return "", errors.New("userKey is empty")
	}
	getCodeUrl := "http://api.damagou.top/apiv1/recognize.html"
	// 构建URL encoded表单数据
	data := url.Values{}
	data.Set("image", url.QueryEscape(base64Code)) // 对Base64字符串进行编码，防止特殊字符干扰
	data.Set("userkey", v.userKey)
	data.Set("type", "1003")
	// 构建请求
	client := &http.Client{}
	req, err := http.NewRequest("POST", getCodeUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	// 设置Content-Type为表单数据
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body), nil
}
