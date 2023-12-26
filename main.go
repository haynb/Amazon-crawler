package main

import (
	"fmt"
	"github.com/go-rod/rod"
	_ "github.com/go-sql-driver/mysql"
)

func DataPage(broswer *rod.Browser, url string) {
	page := broswer.MustPage()
	defer page.MustClose()
	page.MustEmulate(GetDevices())
	page = broswer.MustPage(url)
	page.MustWaitDOMStable()
	CheckWeb(page)
	if exists, moreButton, _ := page.Has("#reviews-medley-footer > div.a-row.a-spacing-medium > a"); exists {
		moreButton.MustClick()
		page.MustWaitDOMStable()
		CheckWeb(page)
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
		CheckWeb(page)
		GetCommantsDetail(page)
	}
}
func main() {
	baseUrl := "https://www.amazon.com/s?k=ebike&s=review-rank&crid=1J6CMDGWJYCJU&qid=1703580152&sprefix=ebike%2Caps%2C275&ref=sr_st_review-rank&ds=v1%3Atp8QjG3ie1Cx7QyakoY5tNNvBLYKjLGpoI0%2B4Zqf7TU"
	broswer := GetBrowser()
	defer broswer.MustClose()
	page := broswer.MustPage()
	page.MustEmulate(GetDevices())
	page = broswer.MustPage(baseUrl)
	page.MustWaitLoad()
	CheckWeb(page)
	bikes := page.MustElementsX("/html/body/div[1]/div[1]/div[1]/div[1]/div/span[1]/div[1]/div/div/div/span/div/div/div")
	for _, bike := range bikes {
		fmt.Println("=========================================")
		link := "https://www.amazon.com/" + bike.MustElementX("div[2]//h2/a/@href").MustText()
		fmt.Println("link:   " + link)
		title := bike.MustElementX("div[2]//h2/a/span").MustText()
		fmt.Println("title:   " + title)
		img := bike.MustElementX("div[1]//img/@src").MustText()
		fmt.Println("img:   " + img)
		DataPage(broswer, link)
	}
}
