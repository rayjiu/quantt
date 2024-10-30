package crawler

import (
	"github.com/playwright-community/playwright-go"
	"github.com/rayjiu/quantt/data/internal/config"
	log "github.com/sirupsen/logrus"
)

type Crawler struct {
	cfg *config.Config
}

func NewCrawler(cfg *config.Config) *Crawler {
	return &Crawler{
		cfg: cfg,
	}
}

func (*Crawler) Start() {
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

	// 打开一个新页面并导航到网址
	page, err := browser.NewPage()
	if err != nil {
		log.Errorf("could not create page: %v", err)
	}
	if response, err := page.Request().Get("https://api.ipify.org?format=json"); err != nil {
		log.Errorf("err:%+v", err)
	} else {
		if body, err := response.Text(); err == nil {
			log.Infof("Http response:%v", body)
		} else {
			log.Errorf("err:%+v", err)
		}
	}
}
