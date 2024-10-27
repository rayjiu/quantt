package crawler

import "github.com/rayjiu/quantt/data/internal/config"

type Crawler struct {
	cfg *config.Config
}

func NewCrawler(cfg *config.Config) *Crawler {
	return &Crawler{
		cfg: cfg,
	}
}

func (*Crawler) Start() {

}
