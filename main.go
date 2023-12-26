package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	browserPath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	u := launcher.New().Bin(browserPath).MustLaunch()
	broswer := rod.New().ControlURL(u).MustConnect()
	defer broswer.MustClose()
	page := broswer.MustPage()
	page.MustEmulate(devices.Device{
		Title:          "iPhone 14",
		Capabilities:   []string{"touch", "mobile"},
		UserAgent:      "Mozilla/5.0 (iPhone; CPU iPhone OS 7_1_2 like Mac OS X)",
		AcceptLanguage: "en",
		Screen: devices.Screen{
			DevicePixelRatio: 2,
			Horizontal: devices.ScreenSize{
				Width:  480,
				Height: 320,
			},
			Vertical: devices.ScreenSize{
				Width:  320,
				Height: 480,
			},
		},
	})
	page = broswer.MustPage("https://www.amazon.com/-/zh/product-reviews/B08B1PV8N1/ref=cm_cr_dp_d_show_all_btm?ie=UTF8&reviewerType=all_reviews")
	page.MustWaitDOMStable()
	data, err := page.Screenshot(true, nil)
	if err != nil {
		panic(err)
	}
	// 保存截图
	utils.OutputFile("screenshot0.png", data)
	if exists, imgLink, _ := page.HasX("/html/body/div/div[1]/div[3]/div/div/form/div[1]/div/div/div[1]/img/@src"); exists {
		//fmt.Println(imgLink)
		img := Img{}
		img.ImgLink = imgLink.MustText()
		err := img.GetImgLink()
		if err != nil {
			panic(err)
		}
		img.EncodeToBase64()
		//fmt.Println("Base:   " + img.Base)
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
	}
	comments := page.MustElements(" div.a-row.a-spacing-small.review-data > span > span.cr-original-review-content")
	for _, comment := range comments {
		fmt.Println(comment.MustText())
	}
}
