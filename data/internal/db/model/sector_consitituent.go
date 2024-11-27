package model

import (
	"time"

	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

// Sector 表示 sec_base_info 表的模型
type SectorConstituent struct {
	SecConstID uuid.UUID `gorm:"column:sec_cons_id;type:uuid;default:gen_random_uuid();primaryKey"`    // 主键字段
	SecCode    string    `gorm:"column:sec_code;type:varchar(255);index;uniqueIndex:sec_const_unique"` // sec_code 字段，建立索引
	StockCode  string    `gorm:"column:stock_code;type:varchar;uniqueIndex:sec_const_unique"`
	MarketType int16     `gorm:"column:market_type;type:int2;uniqueIndex:sec_const_unique"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamptz;"`
}

func (SectorConstituent) TableName() string {
	return "sec_constituent"
}
