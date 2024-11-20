package crawler

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/playwright-community/playwright-go"
	"github.com/rayjiu/quantt/data/internal/db/model"
	"github.com/rayjiu/quantt/data/internal/db/service"
	log "github.com/sirupsen/logrus"
)

type klineCrawler struct {
	urlChain          chan klineUrlInfo
	stopDataWriteChan chan int
	klineChan         chan *model.KlineDay
}

type klineUrlInfo struct {
	url        string
	stockCode  string
	marketType int
	action     int // 0 表示忽略 1表示停止
}

var kCrawler *klineCrawler = &klineCrawler{
	urlChain:          make(chan klineUrlInfo),
	stopDataWriteChan: make(chan int),
	klineChan:         make(chan *model.KlineDay),
}

var rawUrl = `https://push2his.eastmoney.com/api/qt/stock/kline/get?cb=jQuery35106153870917858113_1731491914077&secid=%v.%v&ut=fa5fd1943c7b386f172d6893dbfba10b&fields1=f1,f2,f3,f4,f5,f6&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61&klt=101&fqt=1&end=20500101&lmt=10000&_=1731491914105`

func (k *klineCrawler) startCrawlKlineData(stockCode string, marketType uint32) {
	k.startReceiveData()
	log.Infof("Start to start crawer.")

	k.startBrowser()

	var finalUrl = fmt.Sprintf(rawUrl, marketType, stockCode)

	k.urlChain <- klineUrlInfo{
		url: finalUrl,
	}

	k.urlChain <- klineUrlInfo{action: 1}
}

func (k *klineCrawler) startReceiveData() {
	var targetDatas []*model.KlineDay
	go func() {
		for {
			select {
			case <-k.stopDataWriteChan:
				log.Infof("接受数据等待写入的channel关闭, targetSecotrs.Len:%v", len(targetDatas))
				if len(targetDatas) > 0 {
					service.KlineService.BatchUpsert(targetDatas)
				}
				return

			case data := <-k.klineChan:
				targetDatas = append(targetDatas, data)
				if len(targetDatas) >= 500 {
					service.KlineService.BatchUpsert(targetDatas)
					targetDatas = nil
				}
			}
		}
	}()
}

func (k *klineCrawler) startBrowser() {
	go func() {
		pw, err := playwright.Run()
		if err != nil {
			log.Errorf("could not start Playwright: %v", err)
		}
		defer pw.Stop()
		// 启动 Chromium 浏览器
		browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(true),
		})
		if err != nil {
			log.Errorf("could not launch browser: %v", err)
		}
		defer browser.Close()
		var wg sync.WaitGroup
		var receivedUrl bool
		for {
			urlInfo := <-k.urlChain

			if urlInfo.action == 1 {
				if receivedUrl {
					log.Info("Received stop action. waiting.")
					wg.Wait()
				}

				log.Info("Received stop action. Exiting goroutine.")
				k.stopDataWriteChan <- 1
				break // Exit the loop and end the goroutine
			} else {
				receivedUrl = true
				go k.sendPgeRequest(&wg, browser, urlInfo.url)
			}
		}
	}()
}

func (k *klineCrawler) sendPgeRequest(wg *sync.WaitGroup, browser playwright.Browser, url string) {
	log.Infof("start crawl ur:%v", url)
	wg.Add(1)
	defer wg.Done()

	page, err := browser.NewPage()
	if err != nil {
		log.Errorf("err:%v", err)
		panic(err)
	}

	if response, err := page.Request().Get(url); err != nil {
		log.Errorf("err:%+v", err)
	} else {
		if content, err := response.Text(); err == nil {

			re := regexp.MustCompile(`\((.*)\)`)

			// 使用 FindStringSubmatch 提取括号中的内容
			matches := re.FindStringSubmatch(content)

			if len(matches) > 1 {
				// matches[1] 是括号内的内容
				var content = matches[1]

				var resp KlineResponse
				err := json.Unmarshal([]byte(content), &resp)
				if err != nil {
					log.Error("Error unmarshaling JSON:", content)
					return
				}

				if resp.Data != nil {
					for _, kline := range resp.Data.Klines {
						log.Infof("kline:%+v", kline)
						var parsedKline, err = ParseKline(resp.Data.Code, uint32(resp.Data.Market), kline)
						if err != nil {
							panic(err)
						}
						log.Infof("kline:%+v", parsedKline)
						k.klineChan <- parsedKline
					}

				}

			} else {
				log.Error("No match found")
			}

		}
	}
}
