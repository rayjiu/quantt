package main

import (
	"github.com/rayjiu/quantt/analysis/internal/logging"
	log "github.com/sirupsen/logrus"
)

func main() {
	logging.InitLogger()
	log.Infof("analysis module started.")

	for {
	}
}
