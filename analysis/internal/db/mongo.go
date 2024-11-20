package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type MongoDB struct {
	DBURI string
}

func (m *MongoDB) InitDB() error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().
		ApplyURI(m.DBURI).
		SetMaxPoolSize(100). // 最大 100 个连接
		SetMinPoolSize(10)   // 最少保留 10 个连接
	client, err = mongo.Connect(ctx, opts)

	return err
}

func (m *MongoDB) GetQuotaDB() *mongo.Database {
	return client.Database("quota")
}
