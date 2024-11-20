package model

import (
	"time"
)

// KlineDay represents the structure for the kline_day table
type KlineDay struct {
	StockCode  string    `gorm:"column:stock_code;type:varchar;uniqueIndex:kline_day_unique"`
	MarketType int16     `gorm:"column:market_type;type:int2;uniqueIndex:kline_day_unique"`
	MarketDate int32     `gorm:"column:market_date;type:int4;uniqueIndex:kline_day_unique"`
	Open       float64   `gorm:"column:open;type:float8"`
	Close      float64   `gorm:"column:close;type:float8"`
	Low        float64   `gorm:"column:low;type:float8"`
	High       float64   `gorm:"column:high;type:float8"`
	Volume     int64     `gorm:"column:volume;type:int8"`
	Amount     float64   `gorm:"column:amount;type:float8"`
	Amplitude  float64   `gorm:"column:amplitude"`
	Change     float64   `gorm:"column:change"`
	ChangePct  float64   `gorm:"column:changepct"`
	Turnover   float64   `gorm:"column:turnover"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamptz;default:now()"`
}

// TableName specifies the table name for KlineDay struct
func (KlineDay) TableName() string {
	return "kline_day"
}
