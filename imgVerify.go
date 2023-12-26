package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Img struct {
	ImgLink string
	data    []byte
	Base    string
}

func (i *Img) GetImgLink() error {
	// 发起GET请求获取图片数据
	resp, err := http.Get(i.ImgLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	// 读取响应体
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	i.data = data
	return nil
}
func (i *Img) EncodeToBase64() {
	i.Base = base64.StdEncoding.EncodeToString(i.data)
}

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

func CheckWeb(page *rod.Page) {
	if exists, imgLink, _ := page.HasX("/html/body/div/div[1]/div[3]/div/div/form/div[1]/div/div/div[1]/" +
		"img/@src"); exists {
		img := Img{}
		img.ImgLink = imgLink.MustText()
		err := img.GetImgLink()
		if err != nil {
			panic(err)
		}
		img.EncodeToBase64()
		verify := verify{
			username: "heanyang",
			password: "heanyang",
		}
		err = verify.Login()
		if err != nil {
			panic(err)
		}
		code, err := verify.GetCode(img.Base)
		if err != nil {
			panic(err)
		}
		fmt.Println("Code:   " + code)
		page.MustElement("#captchacharacters").MustInput(code)
		data, err := page.Screenshot(true, nil)
		if err != nil {
			panic(err)
		}
		// 保存截图
		utils.OutputFile("验证码.png", data)
		page.MustElement("body > div > div.a-row.a-spacing-double-large > div.a-section > div > div " +
			"> form > div.a-section.a-spacing-extra-large > div > span > span > button").MustClick()
		page.MustWaitLoad()
	}
}
