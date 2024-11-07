package service

import (
	"fmt"
	"os"
	"testing"

	"github.com/rayjiu/quantt/data/internal/db"
	"github.com/rayjiu/quantt/data/internal/db/model"
	"github.com/rayjiu/quantt/data/internal/db/repository"
)

func TestBatchUpsert(t *testing.T) {
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
	var service = *NewSectorService(repository.NewSectorRepository(db))
	service.BatchUpsert([]model.Sector{
		{
			SecName: "A",
			SecCode: "ABC",
			SecType: 1,
		},
	})

}
