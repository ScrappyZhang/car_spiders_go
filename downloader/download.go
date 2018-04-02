package downloader

import (
	"io"
	"net/http"
	"github.com/axgle/mahonia"
	"log"
	"routinego/car_prices/fake"
)

//获取页面信息
func Get(url string) io.Reader {
	//1.创建一个http请求客户端
	client := &http.Client{}
	//2. GET请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("http.NewRequest err: %v", err)
	}
	// 3. 增加请求头进行伪装
	req.Header.Add("User-Agent", fake.GetUserAgent())
	req.Header.Add("Referer", "https://car.autohome.com.cn")
	// 4.发送请求，获得响应
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("client.Do err: %v", err)
	}
	// 5.响应BODY解码并返回 <meta charset="gbk" />
	mah := mahonia.NewDecoder("gbk")
	return mah.NewReader(resp.Body)
}
