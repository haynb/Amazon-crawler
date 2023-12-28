package myUtils

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/utils"
)

func TakeScreenShot(page *rod.Page, name string) {
	data, err := page.Screenshot(true, nil)
	if err != nil {
		panic(err)
	}
	// 保存截图
	utils.OutputFile(name+".png", data)
}
