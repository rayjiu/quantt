package repository

import (
	"context"

	"github.com/rayjiu/quantt/analysis/internal/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TrendRepository struct {
	Collection *mongo.Collection
}

// NewTrendRepository 初始化 TrendRepository
func NewTrendRepository(database *mongo.Database, collectionName string) *TrendRepository {
	return &TrendRepository{
		Collection: database.Collection(collectionName),
	}
}

// GetTrendsByStock 查询符合条件的 Trend 列表
func (r *TrendRepository) GetTrendsByStock(ctx context.Context, stockId string, marketType uint32, marketDate uint32) ([]model.TrendItem, error) {
	filter := bson.M{
		"stock_id":    stockId,
		"market_type": marketType,
		"market_date": marketDate,
	}

	// 查询选项
	opts := options.Find().SetSort(bson.M{"market_date": 1})

	// 执行查询
	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 遍历查询结果
	var trends []*model.TrendKData
	if err := cursor.All(ctx, &trends); err != nil {
		return nil, err
	}
	if len(trends) > 0 {
		return trends[0].TrendItem, nil
	}

	return nil, nil

}
