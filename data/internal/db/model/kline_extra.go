package model

import "time"

type KlineExtra struct {
	StockCode   string    `gorm:"column:stock_code;type:varchar;uniqueIndex:volume_ratio_unique"`
	MarketType  int16     `gorm:"column:market_type;type:int2;uniqueIndex:volume_ratio_unique"`
	MarketDate  int32     `gorm:"column:market_date;type:int4;uniqueIndex:volume_ratio_unique"`
	VolumeRatio float64   `gorm:"column:volume_ratio;type:float8"`
	UpdateTime  time.Time `gorm:"column:update_time;type:timestamptz;default:now()"`
}

// TableName specifies the table name for KlineExtra
func (KlineExtra) TableName() string {
	return "kline_extra"
}
