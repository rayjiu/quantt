package main

import (
	"fmt"
	"os"

	"github.com/rayjiu/quantt/data/internal/calculator"
	"github.com/rayjiu/quantt/data/internal/db"
	"github.com/rayjiu/quantt/data/internal/db/repository"
	"github.com/rayjiu/quantt/data/internal/db/service"
	"github.com/rayjiu/quantt/data/internal/logging"
)

func main() {
	logging.InitLogger()
	// cfg := config.LoadConfig()

	initDBService()
	// crawler.Start()
	calculator.StartCalHistoryVolumeRate()

	for {
	}
}

func initDBService() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db.InitDB(dsn)
	var db = db.GetDB()
	service.SecotorService = *service.NewSectorService(repository.NewSectorRepository(db))
	service.SecQuoteService = *service.NewSecQuoteService(repository.NewSecQuoteRepository(db))
	service.SecFundFlowService = *service.NewSecFundFlowService(repository.NewSecFundFlowRepository(db))
	service.KlineService = *service.NewklineDayService(repository.NewKlineDayRepository(db))
	service.SectorConsService = *service.NewSectorConsService(repository.NewSectorConsitutuentRepository(db))
	service.KlineExtraService = *service.NewKlineExtraService(repository.NewKlineExtraRepository(db), &service.KlineService)
}
