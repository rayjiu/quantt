package crawler

import (
	"encoding/json"
	"fmt"
	"regexp"

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
}

var eastmoney *crawler = &crawler{}

// startCrawSectorInfo 开始爬取板块列表和基本信息
func (c *crawler) startCrawSectorInfo() {
	c.startReceiveSectorData()
	log.Infof("Start to start crawer.")
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

	eastmoney.sendPgeRequest(browser, concept_url, constants.SectorConcept)
	eastmoney.sendPgeRequest(browser, area_url, constants.SectorArea)
	eastmoney.sendPgeRequest(browser, industry_url, constants.SectorIndustry)
}

func (*crawler) sendPgeRequest(browser playwright.Browser, url string, sectionType int) {

	var page, err = browser.NewPage()
	if err != nil {
		log.Error(err)
	}

	for i := 1; i < 100; i++ {
		log.Infof("Start to crawl page:%v", (i))
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
						log.Infof("finished crawl for %v pages.", (i))
						// stopChan <- 0
						break
					} else {
						var dataDiffs = resp.Data.Diff
						for _, diff := range dataDiffs {
							sectorDataChan <- model.Sector{
								SecCode: diff.F12,
								SecName: diff.F14,
								SecType: int16(sectionType),
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

var stopChan chan int = make(chan int)
var sectorDataChan chan model.Sector = make(chan model.Sector)

func (*crawler) startReceiveSectorData() {
	var targetSectors []model.Sector
	go func() {
		for {
			select {
			case <-stopChan:
				log.Infof("接受数据等待写入的channel关闭, targetSecotrs.Len:%v", len(targetSectors))
				if len(targetSectors) > 0 {
					service.SecotorService.BatchUpsert(targetSectors)
				}
				return

			case data := <-sectorDataChan:
				targetSectors = append(targetSectors, data)
				if len(targetSectors) >= 20 {
					service.SecotorService.BatchUpsert(targetSectors)
					targetSectors = nil
				}
			}

		}
	}()
}
