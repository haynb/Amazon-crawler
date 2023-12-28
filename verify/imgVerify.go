package verify

import (
	"amson/myUtils"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/go-rod/rod"
	querystring "github.com/google/go-querystring/query"
	"io"
	"net/http"
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

type Verify struct {
	Username string
	Password string
	userKey  string
}

func (v *Verify) Login() error {
	loginUrl := "http://api.damagou.top/apiv1/login.html?username=" + v.Username + "&Password=" + v.Password
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
func GetCode(base64Code string) (string, error) {
	type RecognizeRequest struct {
		Image   string `url:"image"`
		Userkey string `url:"userkey"`
		Type    string `url:"type,omitempty"`
		IsJson  string `url:"isJson,omitempty"`
	}
	userKey := "878144361c457012af865f138bb0a984"
	getCodeUrl := "http://api.damagou.top/apiv1/recognize.html"
	// 设置请求参数
	params := RecognizeRequest{
		Image:   base64Code,
		Userkey: userKey, // 用你的userkey替换这里
		Type:    "1003",  // 假设你需要的Type参数是"1001"
		IsJson:  "1",     // 假设你需要以JSON格式返回
	}
	formData, err := querystring.Values(params)
	if err != nil {
		fmt.Println("Form数据编码错误:", err)
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", getCodeUrl, bytes.NewBufferString(formData.Encode()))
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
		code, err := GetCode(img.Base)
		if err != nil {
			panic(err)
		}
		fmt.Println("Code:   " + code)
		page.MustElement("#captchacharacters").MustInput(code)
		myUtils.TakeScreenShot(page, "验证码")
		page.MustElementX("/html/body/div/div[1]/div[3]/div/div/form/div[2]/div/span/span/button").MustClick()
		page.MustWaitLoad()
		myUtils.TakeScreenShot(page, "验证码后")
	}
}
