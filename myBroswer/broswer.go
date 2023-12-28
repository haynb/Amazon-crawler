package myBroswer

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

//	func GetDevices() devices.Device {
//		return devices.Device{
//			Title:          "iPhone 15",
//			Capabilities:   []string{"touch", "mobile"},
//			UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
//			AcceptLanguage: "en",
//			Screen: devices.Screen{
//				DevicePixelRatio: 2,
//				Horizontal: devices.ScreenSize{
//					Width:  480,
//					Height: 320,
//				},
//				Vertical: devices.ScreenSize{
//					Width:  320,
//					Height: 480,
//				},
//			},
//		}
//	}
func GetDevices() devices.Device {
	return devices.Device{
		Title:          "Windows 11 PC",
		Capabilities:   []string{"keyboard", "mouse", "desktop"},
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		AcceptLanguage: "en-US,en;q=0.9",
		Screen: devices.Screen{
			DevicePixelRatio: 1, // PC通常是1
			Horizontal: devices.ScreenSize{
				Width:  1920, // 常见的全高清分辨率宽度
				Height: 1080, // 常见的全高清分辨率高度
			},
			Vertical: devices.ScreenSize{
				Width:  1080, // 当屏幕垂直时（假设）
				Height: 1920, // 当屏幕垂直时（假设）
			},
		},
	}
}
