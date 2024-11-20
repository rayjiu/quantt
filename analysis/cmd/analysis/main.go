package main

import (
	"os"

	"github.com/rayjiu/quantt/analysis/internal"
	"github.com/rayjiu/quantt/analysis/internal/db"
	"github.com/rayjiu/quantt/analysis/internal/db/repository"
	"github.com/rayjiu/quantt/analysis/internal/db/service"
	"github.com/rayjiu/quantt/analysis/internal/logging"
	log "github.com/sirupsen/logrus"
)

func main() {
	logging.InitLogger()
	log.Infof("analysis module started12.")

	initDBService()

	internal.StartAnalysis()

	for {
	}
}

func initDBService() {
	var mongDBURI = os.Getenv("MongoURI")
	var mongoDB = db.MongoDB{
		DBURI: mongDBURI,
	}
	mongoDB.InitDB()

	var db = mongoDB.GetQuotaDB()
	service.TrenddService = service.NewTrendService(repository.NewTrendRepository(db, "kline_trend"))

}
