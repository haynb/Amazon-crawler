package main

import (
	myBroswer "amson/myBroswer"
	"amson/verify"
	"fmt"
	"github.com/go-rod/rod"
	_ "github.com/go-sql-driver/mysql"
)

func DataPage(broswer *rod.Browser, url string) {
	page := broswer.MustPage()
	defer page.MustClose()
	page.MustEmulate(myBroswer.GetDevices())
	page = broswer.MustPage(url)
	page.MustWaitDOMStable()
	verify.CheckWeb(page)
	if exists, moreButton, _ := page.Has("#reviews-medley-footer > div.a-row.a-spacing-medium > a"); exists {
		moreButton.MustClick()
		page.MustWaitDOMStable()
		verify.CheckWeb(page)
		GetCommantsDetail(page)
	} else {
		comments := page.MustElementsX("/html/body/div[1]/div/div[9]/div[35]/div/div/div/div/div[2]/div/div[2]/span[2]/div/div/div[3]/div[3]/div/div")
		if len(comments) == 0 {
			return
		}
		for _, comment := range comments {
			star := comment.MustElementX("div//a/i/span").MustText()
			fmt.Println("star:   " + star)
			msg := comment.MustElementX("div//span/span").MustText()
			fmt.Println("msg:   " + msg)
		}
	}
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
}
func GetCommantsDetail(page *rod.Page) {
	comments := page.MustElementsX("/html/body/div[1]/div[2]/div/div[1]/div/div[1]/div[5]/div[3]/div/div/div/div")
	fmt.Println(len(comments))
	for _, comment := range comments {
		star := comment.MustElementX("div//i/span").MustText()
		fmt.Println("star:   " + star)
		msg := comment.MustElementX("div//span/span").MustText()
		fmt.Println("msg:   " + msg)
	}
	if exists, nextPage, _ := page.Has("#cm_cr-pagination_bar > ul > li.a-last > a"); exists {
		nextPage.MustClick()
		page.MustWaitDOMStable()
		verify.CheckWeb(page)
		GetCommantsDetail(page)
	}
}
func main() {
	//resp, err := http.Get("https://images-na.ssl-images-amazon.com/captcha/wxvwzfzh/Captcha_yrqbxysuih.jpg")
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//// 检查响应状态码
	//if resp.StatusCode != http.StatusOK {
	//	panic(err)
	//}
	//// 读取响应体
	//data, err := io.ReadAll(resp.Body)
	//baseStr := base64.StdEncoding.EncodeToString(data)
	//code, err := verify.GetCode(baseStr)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Code:   " + code)
	baseUrl := "https://www.amazon.com"
	broswer := myBroswer.GetBrowser()
	defer broswer.MustClose()
	page := broswer.MustPage()
	page.MustEmulate(myBroswer.GetDevices())
	page = broswer.MustPage(baseUrl)
	page.MustWaitDOMStable()
	verify.CheckWeb(page)
	//bikes := page.MustElementsX("/html/body/div[1]/div[1]/div[1]/div[1]/div/span[1]/div[1]/div/div/div/span/div/div/div")
	//for _, bike := range bikes {
	//	fmt.Println("=========================================")
	//	link := "https://www.amazon.com/" + bike.MustElementX("div[2]//h2/a/@href").MustText()
	//	fmt.Println("link:   " + link)
	//	title := bike.MustElementX("div[2]//h2/a/span").MustText()
	//	fmt.Println("title:   " + title)
	//	img := bike.MustElementX("div[1]//img/@src").MustText()
	//	fmt.Println("img:   " + img)
	//	DataPage(broswer, link)
	//}
}
