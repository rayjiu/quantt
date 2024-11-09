package crawler

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/panjf2000/ants"
	"github.com/playwright-community/playwright-go"
	"github.com/rayjiu/quantt/data/internal/db/model"
	"github.com/rayjiu/quantt/data/internal/db/service"
	log "github.com/sirupsen/logrus"
)

var (
	fund_flow_url = `https://push2his.eastmoney.com/api/qt/stock/fflow/daykline/get?cb=jQuery112309704473719868352_1731164189783&lmt=0&klt=500&fields1=f1&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65&ut=b2884a393a59ad64002292a3e90d46a5&secid=90.%v&_=1731164189784`
)

// funfFlowCrawler 历史资金流向爬取
type funflowCrawler struct {
	// cfg *config.Config
	urlChain          chan fundFlowUrlInfo
	stopDataWriteChan chan int
	fundFlowChan      chan []*model.SecFundFlow
}

type fundFlowUrlInfo struct {
	url     string
	secCode string
	action  int // 0 表示忽略，1表示停止
}

var ffCrawler *funflowCrawler = &funflowCrawler{
	urlChain:          make(chan fundFlowUrlInfo),
	stopDataWriteChan: make(chan int),
	fundFlowChan:      make(chan []*model.SecFundFlow),
}

// startFundFlowData 开始抓取板块历史资金流向
func (c *funflowCrawler) startHistoryFundFlowData() {
	c.startReceiveData()
	log.Infof("Start to start crawer.")

	service.SecotorService.RefreshCache()
	c.startBrowser()
	var sectorMap = service.SecotorService.GetCachedSector()
	for _, v := range sectorMap {
		var finalUrl = fmt.Sprintf(fund_flow_url, v.SecCode)
		c.urlChain <- fundFlowUrlInfo{url: finalUrl, secCode: v.SecCode}
	}

	c.urlChain <- fundFlowUrlInfo{action: 1}
}

func (c *funflowCrawler) startBrowser() {

	go func() {
		pw, err := playwright.Run()
		if err != nil {
			log.Errorf("could not start Playwright: %v", err)
		}
		defer pw.Stop()

		pools, _ := ants.NewPool(5)

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
			urlInfo := <-c.urlChain
			var secCode = urlInfo.secCode
			if urlInfo.action == 1 {
				if receivedUrl {
					wg.Wait()
				}

				log.Info("Received stop action. Exiting goroutine.")
				c.stopDataWriteChan <- 1
				break // Exit the loop and end the goroutine
			} else {
				receivedUrl = true
				pools.Submit(func() {
					c.sendPgeRequest(&wg, secCode, browser, urlInfo.url)
				})
			}
		}
	}()
}

func (c *funflowCrawler) sendPgeRequest(wg *sync.WaitGroup, secCode string, browser playwright.Browser, url string) {
	wg.Add(1)
	defer wg.Done()
	page, err := browser.NewPage()
	if err != nil {
		log.Errorf("err:%v", err)
	}
	defer page.Close()

	var finalUrl = url
	if response, err := page.Request().Get(finalUrl); err != nil {
		log.Errorf("err:%+v", err)
	} else {

		if content, err := response.Text(); err == nil {

			re := regexp.MustCompile(`\((.*)\)`)

			// 使用 FindStringSubmatch 提取括号中的内容
			matches := re.FindStringSubmatch(content)

			if len(matches) > 1 {
				// matches[1] 是括号内的内容
				var content = matches[1]

				var resp FundFlowResponse
				err := json.Unmarshal([]byte(content), &resp)
				if err != nil {
					log.Error("Error unmarshaling JSON:", content)
					return
				}

				if resp.Data != nil {
					var parsedKlines, err = ParseKlines(secCode, resp.Data.Klines)
					if err != nil {
						panic(err)
					}

					c.fundFlowChan <- parsedKlines
				}

			} else {
				log.Error("No match found")
			}

		} else {
			log.Errorf("err:%+v", err)
		}
	}
}

func (c *funflowCrawler) startReceiveData() {
	var targetFundFlowDatas []*model.SecFundFlow
	go func() {
		for {
			select {
			case <-c.stopDataWriteChan:
				log.Infof("接受数据等待写入的channel关闭, targetSecotrs.Len:%v", len(targetFundFlowDatas))
				if len(targetFundFlowDatas) > 0 {
					service.SecFundFlowService.BatchUpsert(targetFundFlowDatas)
				}

				return

			case data := <-c.fundFlowChan:
				log.Infof("DeleteME -> received:%+v", data)
				targetFundFlowDatas = append(targetFundFlowDatas, data...)
				if len(targetFundFlowDatas) >= 500 {
					service.SecFundFlowService.BatchUpsert(targetFundFlowDatas)
					targetFundFlowDatas = nil
				}

			}

		}
	}()
}
