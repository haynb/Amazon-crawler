package main

import (
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
		Title:          "iPhone 4",
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
	page = broswer.MustPage("https://www.amazon.com/s?k=ebike&s=review-rank&ds=v1%3AUCC6LSAQC0jlyYsdmZz%2BjmXBC%2FHsAaS%2FCezPctwbvRM&crid=2K5E7W838GY1E&qid=1703559079&sprefix=%2Caps%2C332&ref=sr_st_review-rank")
	page.MustWaitDOMStable()
	data, err := page.Screenshot(false, nil)
	if err != nil {
		panic(err)
	}
	// 保存截图
	utils.OutputFile("screenshot.png", data)
}
