package repository

import (
	"context"
	"fmt"

	"github.com/rayjiu/quantt/analysis/internal/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SnapshotRepository struct {
	Collection *mongo.Collection
}

// NewSnapshotRepository 初始化 TrendRepository
func NewSnapshotRepository(database *mongo.Database, collectionName string) *SnapshotRepository {
	return &SnapshotRepository{
		Collection: database.Collection(collectionName),
	}
}

// GetTrendsByStock 查询符合条件的 Trend 列表
func (r *SnapshotRepository) GetLatestSnapshot(ctx context.Context, stockId string, marketType uint32, marketDate uint32) (*model.StockSnapshot, error) {
	filter := bson.M{
		"stock_id":    stockId,
		"market_type": marketType,
		"market_date": marketDate,
	}

	fmt.Printf("filter: %+v \n", filter)

	// 查询选项
	opts := options.Find().SetSort(bson.M{"trade_time": -1})
	// 执行查询
	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 遍历查询结果
	var snaps []model.StockSnapshot
	if err := cursor.All(ctx, &snaps); err != nil {
		return nil, err
	}
	if len(snaps) > 0 {
		return &snaps[0], nil
	}

	return nil, nil

}
