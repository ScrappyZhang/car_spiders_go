package main

import (
	"log"
	"github.com/PuerkitoBio/goquery" //网页解析
	"routinego/car_prices/spiders"
	"routinego/car_prices/downloader"
	"routinego/car_prices/scheduler"
	"routinego/car_prices/model"
	"time"
	"fmt"
)

var (
	StartUrl = "/2sc/%s/a0_0msdgscncgpi1ltocsp1exb4/"
	BaseUrl  = "https://car.autohome.com.cn"
	maxPage int = 99
	cars []spiders.QcCar
)
//开始抓取
func Start(url string, ch chan []spiders.QcCar) {
	body := downloader.Get(BaseUrl + url) //获取页面内容
	doc, err := goquery.NewDocumentFromReader(body) //页面内容生成goquery对象
	if err != nil {
		log.Printf("Downloader.Get err: %v", err)
	}

	currentPage := spiders.GetCurrentPage(doc) //获取当前页页码
	nextPageUrl, _ := spiders.GetNextPageUrl(doc) //从页面获取下一页链接
	//当前页页码未达到99页时，开始解析网页内容
	if currentPage > 0 && currentPage <= maxPage {
		// 获取当前页车列表信息
		cars := spiders.GetCars(doc)
		log.Println(cars)
		ch <- cars
		//添加下一页URL到URL下载队列
		if url := nextPageUrl; url != "" {
			scheduler.AppendUrl(url)
		}

		log.Println(url)
	} else {
		log.Println("Max page !!!")
	}
}

func main() {
	//1.获取所有城市信息，并组织相应的url信息
	citys := spiders.GetCitys()
	for _, v := range citys {
		scheduler.AppendUrl(fmt.Sprintf(StartUrl, v.Pinyin))
	}
	//2.配置6秒无响应则跳出协程
	start := time.Now()
	delayTime := time.Second * 6
	//3. 创建通道
	ch := make(chan []spiders.QcCar)

Loop:
	for {
		//4.1 开始抓取
		if url := scheduler.PopUrl(); url != "" {
			go Start(url, ch)
		}
		select {
		// 4.2 读取通道内信息，将其添加到cars
		case r := <-ch:
			cars = append(cars, r...)
			//4.2.1 开启协程抓取下一页
			go Start(scheduler.PopUrl(), ch)
		// 4.2 若超时6秒，则抛弃当前页
		case <-time.After(delayTime):
			log.Println("Timeout...")
			break Loop
		}
	}
	//将数据写入数据库
	if len(cars) > 0 {
		model.AddCars(cars)
	}
	//记录消耗时间
	log.Printf("Time: %s", time.Since(start) - delayTime)
}