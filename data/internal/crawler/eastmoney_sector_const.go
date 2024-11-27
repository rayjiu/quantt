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

type secConstituentCrawler struct {
	urlChain          chan secConsUrlInfo
	stopDataWriteChan chan int
	consChan          chan *model.SectorConstituent
}

type secConsUrlInfo struct {
	url     string
	secCode string
	action  int // 0 表示忽略 1表示停止
}

var secCCrrawler secConstituentCrawler = secConstituentCrawler{
	urlChain:          make(chan secConsUrlInfo),
	stopDataWriteChan: make(chan int),
	consChan:          make(chan *model.SectorConstituent),
}

var consRawUlr = `https://push2.eastmoney.com/api/qt/clist/get?cb=jQuery1123028842927578359423_1732548402765&fid=f174&po=1&pz=200&pn=1&np=1&fltt=2&invt=2&ut=b2884a393a59ad64002292a3e90d46a5&fs=b:%v&fields=f12,f14,f13`

func (k *secConstituentCrawler) startCrawAllConsData() {
	k.startReceiveData()
	log.Infof("Start to start crawer.")

	service.SecotorService.RefreshCache()
	k.startBrowser()

	var sectorMap = service.SecotorService.GetCachedSector()
	for _, v := range sectorMap {
		var finalUrl = fmt.Sprintf(consRawUlr, v.SecCode)
		k.urlChain <- secConsUrlInfo{url: finalUrl, secCode: v.SecCode}
	}

	k.urlChain <- secConsUrlInfo{action: 1}
}

func (k *secConstituentCrawler) startReceiveData() {
	var targetDatas []*model.SectorConstituent
	go func() {
		for {
			select {
			case <-k.stopDataWriteChan:
				log.Infof("接受数据等待写入的channel关闭, targetSecotrs.Len:%v", len(targetDatas))
				if len(targetDatas) > 0 {
					service.SectorConsService.BatchUpsert(targetDatas)
				}
				return

			case data := <-k.consChan:
				targetDatas = append(targetDatas, data)
				if len(targetDatas) >= 500 {
					service.SectorConsService.BatchUpsert(targetDatas)
					targetDatas = nil
				}
			}
		}
	}()
}

func (k *secConstituentCrawler) startBrowser() {
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
				go k.sendPgeRequest(&wg, browser, urlInfo.url, urlInfo.secCode)
			}
		}
	}()
}

func (k *secConstituentCrawler) sendPgeRequest(wg *sync.WaitGroup, browser playwright.Browser, url, secCode string) {
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

				var resp Response
				err := json.Unmarshal([]byte(content), &resp)
				if err != nil {
					log.Error("Error unmarshaling JSON:", content)
					return
				}

				if resp.Data != nil {
					for _, cons := range resp.Data.Diff {
						k.consChan <- &model.SectorConstituent{
							SecCode:    secCode,
							StockCode:  cons.F12,
							MarketType: int16(cons.F13),
						}
					}

				}

			} else {
				log.Error("No match found")
			}

		}
	}
}
