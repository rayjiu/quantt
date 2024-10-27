package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/rayjiu/quantt/data/internal/config"
	"github.com/rayjiu/quantt/data/internal/crawler"
	"github.com/rayjiu/quantt/data/internal/logging"
)

func main() {
	logging.InitLogger()
	cfg := config.LoadConfig()
	crawler := crawler.NewCrawler(cfg)

	crawler.Start()
	log.Info("crawler main start complete")
	for {
	}
}
