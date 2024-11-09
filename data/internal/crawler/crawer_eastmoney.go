package crawler

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/rayjiu/quantt/data/internal/constants"
	"github.com/rayjiu/quantt/data/internal/db/model"
	"github.com/rayjiu/quantt/data/internal/db/service"
	log "github.com/sirupsen/logrus"
)

var (
	industry_url = `http://65.push2.eastmoney.com/api/qt/clist/get?cb=jQuery1124012922987576750145_1730305104306&pn=%v&pz=30&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&dect=1&wbp2u=|0|0|0|web&fid=f3&fs=m:90+t:2&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152,f133,f104,f105&_=1730305104307`
	area_url     = `http://90.push2.eastmoney.com/api/qt/clist/get?cb=jQuery112405339790134151299_1730305902747&pn=%v&pz=30&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&dect=1&wbp2u=|0|0|0|web&fid=f3&fs=m:90+t:1+f:!50&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f26,f22,f33,f11,f62,f128,f136,f115,f152,f124,f107,f104,f105,f140,f141,f207,f208,f209,f222&_=1730305902748`
	concept_url  = `http://5.push2.eastmoney.com/api/qt/clist/get?cb=jQuery112405271695738408106_1730305793916&pn=%v&pz=20&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&dect=1&wbp2u=|0|0|0|web&fid=f3&fs=m:90+t:3+f:!50&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f26,f22,f33,f11,f62,f128,f136,f115,f152,f124,f107,f104,f105,f140,f141,f207,f208,f209,f222&_=1730305793917`
)

type crawler struct {
	// cfg *config.Config
	urlChain           chan urlInfo
	stopDataWriteChan  chan int
	sectorBaseInfoChan chan *model.Sector
	sectorQuoteChan    chan *model.SecQuote
}

type urlInfo struct {
	url     string
	secType int
	action  int // 0 表示忽略，1表示停止
}

var eastmoney *crawler = &crawler{
	urlChain:           make(chan urlInfo),
	stopDataWriteChan:  make(chan int),
	sectorBaseInfoChan: make(chan *model.Sector),
	sectorQuoteChan:    make(chan *model.SecQuote),
}

// startCrawSectorInfo 开始爬取板块列表和基本信息
func (c *crawler) startCrawSectorInfo() {
	c.startReceiveSectorData()
	log.Infof("Start to start crawer.")

	c.startBrowser()
	c.urlChain <- urlInfo{url: industry_url, secType: constants.SectorIndustry}
	c.urlChain <- urlInfo{url: concept_url, secType: constants.SectorConcept}
	c.urlChain <- urlInfo{url: area_url, secType: constants.SectorArea}
	c.urlChain <- urlInfo{action: 1}
}

func (c *crawler) startBrowser() {

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
			urlInfo := <-c.urlChain

			if urlInfo.action == 1 {
				if receivedUrl {
					wg.Wait()
				}

				log.Info("Received stop action. Exiting goroutine.")
				c.stopDataWriteChan <- 1
				break // Exit the loop and end the goroutine
			} else {
				receivedUrl = true
				go c.sendPgeRequest(&wg, browser, urlInfo.url, urlInfo.secType)
			}
		}
	}()
}

func (c *crawler) sendPgeRequest(wg *sync.WaitGroup, browser playwright.Browser, url string, sectionType int) {
	wg.Add(1)
	defer wg.Done()
	page, err := browser.NewPage()
	if err != nil {
		log.Errorf("err:%v", err)
	}

	for i := 1; i < 100; i++ {
		var finalUrl = fmt.Sprintf(url, i)
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

					var resp Response
					err := json.Unmarshal([]byte(content), &resp)
					if err != nil {
						fmt.Println("Error unmarshaling JSON:", content)
						return
					}

					if resp.Data == nil {
						log.Infof("[%v]finished crawl for %v pages.", sectionType, (i))
						// stopChan <- 0
						break
					} else {
						var dataDiffs = resp.Data.Diff
						for _, diff := range dataDiffs {
							c.sectorBaseInfoChan <- &model.Sector{
								SecCode:    diff.F12,
								SecName:    diff.F14,
								SecType:    int16(sectionType),
								UpdateTime: time.Now(),
							}

							c.sectorQuoteChan <- &model.SecQuote{
								SecCode:     diff.F12,
								LastPrice:   diff.F2,
								ChgRatio:    diff.F4,
								ExchgRatio:  diff.F8, // 给出的数据是原始数据*100的
								TotalMktCap: float64(diff.F20),
								MarketType:  constants.SecMarketType,
							}
						}
					}

				} else {
					fmt.Println("No match found")
				}

			} else {
				log.Errorf("err:%+v", err)
			}
		}
	}
}

func (c *crawler) startReceiveSectorData() {
	var targetSectors []*model.Sector
	var targetSecQuotes []*model.SecQuote
	go func() {
		for {
			select {
			case <-c.stopDataWriteChan:
				log.Infof("接受数据等待写入的channel关闭, targetSecotrs.Len:%v", len(targetSectors))
				if len(targetSectors) > 0 {
					service.SecotorService.BatchUpsert(targetSectors)
				}
				if len(targetSecQuotes) > 0 {
					service.SecQuoteService.BatchUpsert(targetSecQuotes)
				}
				return

			case data := <-c.sectorBaseInfoChan:
				targetSectors = append(targetSectors, data)
				if len(targetSectors) >= 100 {
					service.SecotorService.BatchUpsert(targetSectors)
					targetSectors = nil
				}

			case data := <-c.sectorQuoteChan:
				targetSecQuotes = append(targetSecQuotes, data)
				if len(targetSectors) >= 100 {
					service.SecQuoteService.BatchUpsert(targetSecQuotes)
					targetSecQuotes = nil
				}

			}

		}
	}()
}
