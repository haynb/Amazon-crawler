package main

import (
	"amson/database"
	myBroswer "amson/myBroswer"
	"amson/myUtils"
	"amson/verify"
	"context"
	"fmt"
	"github.com/go-rod/rod"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"sync"
	"time"
)

type BikeData struct {
	Title    string    `bson:"title"`
	Url      string    `bson:"url"`
	Brand    string    `bson:"brand"`
	Price    string    `bson:"price"`
	ImgLink  string    `bson:"imgLink"`
	Comments []Comment `bson:"comments"`
}
type Comment struct {
	Stars int    `bson:"stars"`
	Text  string `bson:"text"`
}

func DataPage(broswer *rod.Browser, url string, data *BikeData, wg *sync.WaitGroup, limiter *chan struct{}) {
	page := broswer.MustPage(url).MustEmulate(myBroswer.GetDevices())
	defer page.MustClose()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//start := time.Now() // 记录当前时间为开始时间
	page.Context(ctx).WaitLoad()
	//fmt.Println(time.Since(start))
	verify.CheckWeb(page)
	if exist, _, _ := page.HasX("/html/body/div[1]/div/div[9]/div[3]/div[4]/div[39]/div//table/tbody/tr"); exist {
		brands := page.MustElementsX("/html/body/div[1]/div/div[9]/div[3]/div[4]/div[39]/div//table/tbody/tr")
		for _, brand := range brands {
			left := brand.MustElementX("td[1]/span").MustText()
			if left != "Brand" {
				continue
			}
			right := brand.MustElementX("td[2]/span").MustText()
			data.Brand = right
		}
	}
	if exists, moreButton, _ := page.Has("#reviews-medley-footer > div.a-row.a-spacing-medium > a"); exists {
		moreButton.MustClick()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		page.Context(ctx).WaitLoad()
		verify.CheckWeb(page)
		GetCommantsDetail(page, data)
	} else {
		if exists, _, _ := page.HasX("/html/body/div[1]/div/div[9]/div[35]/div/div/div/div/div[2]/div/div[2]/span[2]/div/div/div[3]/div[3]/div/div"); exists {
			comments := page.MustElementsX("/html/body/div[1]/div/div[9]/div[35]/div/div/div/div/div[2]/div/div[2]/span[2]/div/div/div[3]/div[3]/div/div")
			for _, comment := range comments {
				var star int
				if exists, starX, _ := comment.HasX("div//i/span"); exists {
					star, _ = strconv.Atoi(starX.MustText()[:1])
					//star = starX.MustText()
				}
				var msg string
				if exists, msgX, _ := comment.HasX("div[4]/span/span"); exists {
					msg = msgX.MustText()
				}
				data.Comments = append(data.Comments, Comment{Stars: star, Text: msg})
				fmt.Print("...")
			}
		}
		fmt.Println("")
	}
	database.InsertOne(data)
	//fmt.Println("commants: " + strconv.Itoa(len(data.Comments)))
	//fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	defer func() {
		<-*limiter
	}()
	defer wg.Done()
}
func GetCommantsDetail(page *rod.Page, data *BikeData) {
	comments := page.MustElementsX("/html/body/div[1]/div[2]/div/div[1]/div/div[1]/div[5]/div[3]/div/div/div/div")
	for _, comment := range comments {
		var star int
		if exists, starX, _ := comment.HasX("div//i/span"); exists {
			star, _ = strconv.Atoi(starX.MustText()[:1])
			//star = starX.MustText()
		}
		var msg string
		if exists, msgX, _ := comment.HasX("div[4]/span/span"); exists {
			msg = msgX.MustText()
		}
		data.Comments = append(data.Comments, Comment{Stars: star, Text: msg})
		fmt.Print("***")
	}
	fmt.Println("")
	if exists, nextPage, _ := page.Has("#cm_cr-pagination_bar > ul > li.a-last > a"); exists {
		//start := time.Now() // 记录当前时间为开始时间
		nextPage.MustClick()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		page.Context(ctx).WaitStable(1 * time.Second)
		//fmt.Println(time.Since(start))
		verify.CheckWeb(page)
		GetCommantsDetail(page, data)
	}
}
func ChangeCountry(page *rod.Page) {
	// 切换地区
	if exists, changeCountry, _ := page.HasX("//*[@id=\"nav-global-location-popover-link\"]"); exists {
		changeCountry.MustClick()
		page.MustWaitDOMStable()
		input := page.MustElementX("//*[@id=\"GLUXZipUpdateInput\"]")
		input.MustInput("20001")
		button := page.MustElementX("//*[@id=\"GLUXZipUpdate\"]/span/input")
		button.MustClick()
		page.MustWaitDOMStable()
		myUtils.TakeScreenShot(page, "切换地区")
		doneButton := page.MustElementX("/html/body/div[6]/div/div/div[2]/span/span/input")
		doneButton.MustClick()
		page.MustWaitDOMStable()
		verify.CheckWeb(page)
	}
	fmt.Println("已经切换了地区")
}
func main() {
	thisPage := 0
	thisCount := 0
	broswer := myBroswer.GetBrowser()
	defer broswer.MustClose()
	pageNum := 1
	var wg sync.WaitGroup
	limiter := make(chan struct{}, 8)
	count := 0
	for pageNum < 21 {
		if pageNum < thisPage {
			pageNum++
			continue
		}
		page := broswer.MustPage()
		page.MustEmulate(myBroswer.GetDevices())
		baseUrl := "https://www.amazon.com/s?k=ebike&s=review-rank&page=" +
			strconv.Itoa(pageNum) +
			"&crid=2CR8AL3A9I0TN&qid=1703656735&sprefix=ebike%2Caps%2C294&ref=sr_pg_18"
		fmt.Println("新一页： " + baseUrl)
		page = broswer.MustPage(baseUrl)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		//start := time.Now() // 记录当前时间为开始时间
		page.Context(ctx).WaitLoad()
		//fmt.Println(time.Since(start))
		verify.CheckWeb(page)
		if count == 0 {
			ChangeCountry(page)
			count++
		}
		if thisPage != 0 && count == 1 {
			page = broswer.MustPage(baseUrl)
			ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
			//start := time.Now() // 记录当前时间为开始时间
			page.Context(ctx2).WaitLoad()
			//fmt.Println(time.Since(start))
			cancel2()
			verify.CheckWeb(page)
			count++
		}
		myUtils.TakeScreenShot(page, "第"+strconv.Itoa(pageNum)+"页")
		bikes := page.MustElementsX("/html/body/div[1]/div[1]/div[1]/div[1]/div/span[1]/div[1]/div/div/div/span/div/div/div")
		for count, bike := range bikes {
			if count < thisCount {
				continue
			}
			limiter <- struct{}{}
			fmt.Println("=========================================     page：" + strconv.Itoa(pageNum) + "   NUM: " + strconv.Itoa(count))
			bikeData := BikeData{}
			link := "https://www.amazon.com/" + bike.MustElementX("div[2]//h2/a/@href").MustText()
			fmt.Println("link:   " + link)
			bikeData.Url = link
			title := bike.MustElementX("div[2]//h2/a/span").MustText()
			bikeData.Title = title
			img := bike.MustElementX("div[1]//img/@src").MustText()
			bikeData.ImgLink = img
			if exists, priceBlock, _ := bike.HasX("div[2]//a/span/span[1]"); exists {
				price := priceBlock.MustText()
				//fmt.Println("price:   " + price)
				bikeData.Price = price
			}
			wg.Add(1)
			go DataPage(broswer, link, &bikeData, &wg, &limiter)
		}
		wg.Wait()
		cancel()
		pageNum++
	}
	defer database.Disconnect()
}
