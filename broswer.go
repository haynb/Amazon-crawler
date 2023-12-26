package main

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
)

func GetBrowser() *rod.Browser {
	browserPath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	u := launcher.New().Bin(browserPath).MustLaunch()
	broswer := rod.New().ControlURL(u).MustConnect()
	return broswer
}
func GetDevices() devices.Device {
	return devices.Device{
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
	}
}
